package warehouse

// Information-library controllers do not use StringHelper.getSql: their service
// delegates to QueryWrapperHelper, including range filters and implicit ordering.
import (
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

var infoEntities = map[string]*ent{"/xmmcCoalInfo": entXmmcCoalInfo, "/xmmcCokingInfo": entXmmcCokingInfo, "/xmmcRailwayInfo": entXmmcRailwayInfo, "/xmmcWashingInfo": entXmmcWashingInfo}

func init() {
	for i := range tplCtls {
		if _, ok := infoEntities[tplCtls[i].Path]; ok {
			tplCtls[i].Skip = append(tplCtls[i].Skip, "/listByPage", "/updateById")
		}
	}
	addExtra(func(mux *http.ServeMux, s *Service) {
		for path, entity := range infoEntities {
			e := entity
			reg := func(method, suffix string, h http.HandlerFunc) {
				mux.HandleFunc(method+" /saas-warehouse"+path+suffix, authGuard(h))
			}
			reg("POST", "/listByPage", func(w http.ResponseWriter, r *http.Request) { s.infoList(w, r, e) })
			reg("POST", "/updateById", func(w http.ResponseWriter, r *http.Request) { s.infoUpdate(w, r, e) })
			reg("POST", "/importExcel", func(w http.ResponseWriter, r *http.Request) { s.infoImport(w, r, e) })
			reg("GET", "/generateAndDownloadExcel", func(w http.ResponseWriter, r *http.Request) { infoTemplate(w, e) })
		}
		mux.HandleFunc("POST /saas-warehouse/wechat/transfer", authGuard(func(w http.ResponseWriter, r *http.Request) { platform.Success(w, nil) }))
		mux.HandleFunc("GET /saas-warehouse/busVr/getVr", authGuard(s.infoVR))
	})
}

// LocalDateTimeDeserializer returns null for invalid strings; a date-only Start
// expands to midnight, and End to LocalTime.MAX. These are carrier fields only.
func normalizeInfoRanges(m map[string]any) {
	for key, value := range m {
		if !(strings.HasSuffix(key, "Start") || strings.HasSuffix(key, "End")) {
			continue
		}
		text, ok := value.(string)
		if !ok {
			m[key] = nil
			continue
		}
		if day, err := time.Parse("2006-01-02", text); err == nil {
			suffix := "T00:00:00"
			if strings.HasSuffix(key, "End") {
				suffix = "T23:59:59.999999999"
			}
			m[key] = day.Format("2006-01-02") + suffix
			continue
		}
		if t, err := time.Parse("2006-01-02T15:04:05", text); err == nil {
			m[key] = t.Format("2006-01-02T15:04:05.999999999")
		} else {
			m[key] = nil
		}
	}
}

func infoWhere(e *ent, m map[string]any) (string, []any) {
	excel := map[string]bool{}
	for _, f := range e.q {
		excel[f.Name] = true
	}
	cond := []string{}
	args := []any{}
	for _, f := range e.types {
		val := m[f.Name]
		if val == nil {
			continue
		}
		v, ok := bindVal(f.Kind, val)
		if !ok || v == nil {
			continue
		}
		col := platform.Underscore(map[string]any{f.Name: v})
		var name string
		for k := range col {
			name = k
		}
		op := ""
		switch {
		case f.Kind == "t" || f.Kind == "T":
			if strings.HasSuffix(f.Name, "Start") {
				name = strings.TrimSuffix(name, "_start")
				op = ">="
			} else if strings.HasSuffix(f.Name, "End") {
				name = strings.TrimSuffix(name, "_end")
				op = "<="
			}
		case f.Kind == "d":
			if strings.HasSuffix(f.Name, "Min") {
				name = strings.TrimSuffix(name, "_min")
				op = ">="
			} else if strings.HasSuffix(f.Name, "Max") {
				name = strings.TrimSuffix(name, "_max")
				op = "<="
			}
		case f.Kind == "s" && excel[f.Name] && v != "":
			op = "LIKE"
			v = "%" + v.(string) + "%"
		}
		if op != "" {
			if f.Kind == "t" || f.Kind == "T" {
				v = strings.Replace(v.(string), "T", " ", 1)
			}
			cond = append(cond, name+" "+op+" ?")
			args = append(args, v)
		}
	}
	if len(cond) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(cond, " AND "), args
}

func (s *Service) infoList(w http.ResponseWriter, r *http.Request, e *ent) {
	m := bodyMap(r)
	normalizeInfoRanges(m)
	if !e.validateBody(m) {
		platform.Unknown(w)
		return
	}
	where, args := infoWhere(e, m)
	q := "SELECT " + strings.Join(e.dbcols, ",") + " FROM " + e.table + where + " ORDER BY update_time DESC, area ASC"
	if _, ok := m["page"].(map[string]any); !ok {
		rows, cols, err := s.queryEnt(r, q, args...)
		if err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, e.shapeRows(rows, cols))
		return
	}
	current, size := pageOf(m)
	var total int64
	if err := s.pick(r).QueryRow("SELECT COUNT(*) FROM "+e.table+where, args...).Scan(&total); err != nil {
		platform.Unknown(w)
		return
	}
	if total == 0 {
		platform.Success(w, ipageOf([]map[string]any{}, 0, current, size))
		return
	}
	if size >= 0 {
		q += " LIMIT ? OFFSET ?"
		args = append(args, size, (current-1)*size)
	}
	rows, cols, err := s.queryEnt(r, q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, ipageOf(e.shapeRows(rows, cols), total, current, size))
}

func infoTemplate(w http.ResponseWriter, e *ent) {
	head := []string{}
	for _, f := range e.excel {
		head = append(head, f.header)
	}
	rec := httptest.NewRecorder()
	writeXlsx(rec, [][]string{head})
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment;filename="+url.QueryEscape("导入模版.xlsx"))
	_, _ = w.Write(rec.Body.Bytes())
}

func (s *Service) infoImport(w http.ResponseWriter, r *http.Request, e *ent) {
	fh, ok := formFile(r)
	if !ok {
		platform.Unknown(w)
		return
	}
	irt := readExcelFile(fh, e.excel, "", true, emptyDict())
	for _, row := range irt.DataList {
		row["province_code"] = "14"
		row["province"] = "山西省"
		row["import_batch_number"] = time.Now().Format("20060102150405")
		if city := btsStr(row["city"]); strings.TrimSpace(city) != "" {
			rows, err := queryRows(s.pick(r), "SELECT code FROM chain_area_info WHERE name LIKE ?", city)
			if err != nil || len(rows) > 1 {
				platform.Unknown(w)
				return
			}
			if len(rows) == 1 {
				row["city_code"] = rows[0]["code"]
			} else {
				delete(row, "city_code")
			}
		}
		if area := btsStr(row["area"]); strings.TrimSpace(area) != "" {
			code := row["city_code"]
			if code == nil {
				code = "null"
			}
			rows, err := queryRows(s.pick(r), "SELECT code FROM chain_area_info WHERE name LIKE ? AND parent_code=?", area, code)
			if err != nil || len(rows) > 1 {
				platform.Unknown(w)
				return
			}
			if len(rows) == 1 {
				row["area_code"] = rows[0]["code"]
			} else {
				delete(row, "area_code")
			}
		}
	}
	irt.SucInfos = []any{"导入成功 " + strconv.Itoa(len(irt.DataList))}
	irt.Errs = []any{}
	if len(irt.ErrInfos) > 0 {
		irt.Errs = []any{"导入失败 0"}
		platform.FailData(w, 999, "失败", irt.ErrInfos)
		return
	}
	for _, row := range irt.DataList {
		vals := insertVals(row)
		delete(vals, "id") // IdType.AUTO saveBatch allocates new ids, including imported id cells.
		now := time.Now().Format("2006-01-02 15:04:05")
		if vals["create_time"] == nil {
			vals["create_time"] = now
		}
		if vals["update_time"] == nil {
			vals["update_time"] = now
		}
		id, err := platform.Insert(s.pick(r), e.table, vals)
		if err != nil {
			platform.Unknown(w)
			return
		}
		row["id"] = strconv.FormatInt(id, 10)
		row["create_time"] = vals["create_time"]
		row["update_time"] = vals["update_time"]
	}
	// ImportResult contains Java entities: full keys, database fill mutations included.
	shaped := make([]map[string]any, 0, len(irt.DataList))
	for _, row := range irt.DataList {
		shaped = append(shaped, e.shapeOne(row))
	}
	platform.Success(w, map[string]any{"dataList": shaped, "sucInfos": irt.SucInfos, "errs": irt.Errs, "errInfos": irt.ErrInfos})
}

func decimalEqual(a, b any) bool {
	if a == nil {
		a = "0"
	}
	if b == nil {
		b = "0"
	}
	aa, ok := new(big.Rat).SetString(btsStr(a))
	if !ok {
		return false
	}
	bb, ok := new(big.Rat).SetString(btsStr(b))
	return ok && aa.Cmp(bb) == 0
}

func infoChangeLog(e *ent, original, vals map[string]any) string {
	out := "对比结果为：\n"
	for _, f := range e.excel {
		v := vals[f.col]
		if v == nil {
			continue
		}
		bef := original[f.col]
		a, b := exportCell(bef, f.kind), exportCell(v, f.kind)
		if f.kind == "d" {
			if decimalEqual(bef, v) {
				continue
			}
			if bef == nil {
				a = "0"
			}
		} else if a == b {
			continue
		}
		label := strings.Split(f.header, "（")[0]
		out += label + "：" + a + "——》" + b + "\n"
	}
	return out
}

func (s *Service) infoUpdate(w http.ResponseWriter, r *http.Request, e *ent) {
	m := bodyMap(r)
	if !e.validateBody(m) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	original, err := platform.GetByID(d, e.table, m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	vals := pick(m, e.dbcols)
	delete(vals, "id")
	if e == entXmmcCoalInfo {
		if original == nil {
			platform.Unknown(w)
			return
		}
		if !decimalEqual(original["factory_tax_price"], m["factoryTaxPrice"]) {
			vals["price_upd_time"] = time.Now().Format("2006-01-02 15:04:05")
		}
	}
	vals["update_time"] = time.Now().Format("2006-01-02 15:04:05")
	if _, err := platform.UpdateByID(d, e.table, m["id"], vals); err != nil {
		platform.Unknown(w)
		return
	}
	// Other three services update first, then dereference originalInfo for logging.
	if original == nil {
		platform.Unknown(w)
		return
	}
	user := currentUser(r)
	phone := ""
	_ = d.QueryRow("SELECT phone FROM chain_auth_user WHERE id=?", user.UserID).Scan(&phone)
	if _, err := platform.Insert(d, "bus_common_log", map[string]any{"user_id": strconv.FormatInt(user.UserID, 10), "user_name": nil /* Java UserInfo.userName is not the inherited username property. */, "user_phone": phone, "pro_id": m["id"], "operation": infoChangeLog(e, original, vals), "operation_type": 2}); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// Java carries old vendor credentials inline. The Go migration takes the same
// vendor calls from environment so credentials never become repository assets.
func (s *Service) infoVR(w http.ResponseWriter, r *http.Request) {
	key, secret := os.Getenv("REALSEE_APP_KEY"), os.Getenv("REALSEE_APP_SECRET")
	if key == "" || secret == "" {
		platform.Unknown(w)
		return
	}
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.PostForm("https://app-gateway.realsee.com/auth/access_token", url.Values{"app_key": {key}, "app_secret": {secret}})
	if err != nil {
		platform.Unknown(w)
		return
	}
	defer resp.Body.Close()
	var token struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	if resp.StatusCode != 200 || json.NewDecoder(resp.Body).Decode(&token) != nil || token.Data.AccessToken == "" {
		platform.Unknown(w)
		return
	}
	req, err := http.NewRequestWithContext(r.Context(), "GET", "https://app-gateway.realsee.com/open/v1/entity/vr?vr_code=806gyq3v0ddj897JO1", nil)
	if err != nil {
		platform.Unknown(w)
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", token.Data.AccessToken)
	result, err := client.Do(req)
	if err != nil {
		platform.Unknown(w)
		return
	}
	defer result.Body.Close()
	b, err := io.ReadAll(io.LimitReader(result.Body, 8<<20))
	if err != nil || result.StatusCode != 200 || !json.Valid(b) {
		platform.Unknown(w)
		return
	}
	platform.Success(w, json.RawMessage(b))
}
