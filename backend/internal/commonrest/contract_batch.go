package commonrest

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"
	"unicode/utf16"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func (s *Service) saveContractRows(r *http.Request, rows []map[string]any) error {
	if len(rows) == 0 {
		return nil
	}
	return platform.WithTx(s.pick(r), func(tx *sql.Tx) error { return s.saveContractRowsTx(r, tx, rows) })
}

func (s *Service) saveContractRowsTx(r *http.Request, tx *sql.Tx, rows []map[string]any) error {
	for _, row := range rows {
		name, ok := row["sales_unit"].(string)
		if !ok || row["contract_type"] == nil {
			return errors.New("missing sales unit or contract type")
		}
		cust := 2
		if len(utf16.Encode([]rune(name))) <= 4 {
			cust = 1
		}
		partner := 1
		if num(row["contract_type"]) == 2 {
			name, _ = row["buyer_unit"].(string)
			partner = 2
		}
		var id int64
		var stored sql.NullInt64
		err := tx.QueryRow("SELECT id,cust_type FROM chain_contract_supplier WHERE enterprise_name=? AND partner_type=? LIMIT 1", name, partner).Scan(&id, &stored)
		if err == sql.ErrNoRows {
			id, err = platform.TxInsertNoFill(tx, "chain_contract_supplier", map[string]any{"enterprise_name": name, "partner_type": partner, "cust_type": cust})
			stored = sql.NullInt64{Int64: int64(cust), Valid: true}
		}
		if err != nil {
			return err
		}
		row["supplier_id"] = id
		if stored.Valid {
			row["cust_type"] = stored.Int64
		} else {
			row["cust_type"] = nil
		}
		if uid := currentUser(r).UserID; uid > 0 {
			if row["create_user_id"] == nil {
				row["create_user_id"] = uid
			}
			if row["update_user_id"] == nil {
				row["update_user_id"] = uid
			}
		}
	}
	for _, row := range rows {
		vals := map[string]any{}
		for _, col := range authRestEnts["ChainContract"].dbcols {
			if row[col] != nil {
				vals[col] = row[col]
			}
		}
		delete(vals, "id")
		id, err := platform.TxInsert(tx, "chain_contract", insertVals(vals))
		if err != nil {
			return err
		}
		row["id"] = id
		for _, col := range []string{"create_time", "update_time"} {
			if row[col] == nil {
				row[col] = nowSpace()
			}
		}
	}
	return nil
}

func (s *Service) contractDictionary(r *http.Request) (map[string]string, error) {
	rows, err := queryRows(s.pick(r), "SELECT d_key,d_value FROM chain_common_dict")
	if err != nil {
		return nil, err
	}
	out := map[string]string{}
	for _, row := range rows {
		// MyBatis 对全 NULL 行返回 null 元素,Java dictMap 的 m.get 直接 NPE→500(8765 实测,
		// chain 库 id 106/107 两行)。1:1 复刻该失败;单列 NULL 行 Java 有守卫,跳过。
		if row["d_key"] == nil && row["d_value"] == nil {
			return nil, errors.New("java dictMap NPE on all-null dict row")
		}
		if row["d_key"] != nil && row["d_value"] != nil {
			out[btsStr(row["d_key"])] = btsStr(row["d_value"])
		}
	}
	return out, nil
}
func (s *Service) contractCode(r *http.Request, kind string) (string, error) {
	left := "HTCG."
	if kind == "2" {
		left = "HTXS."
	}
	left += currentUser(r).EnterpriseCode + time.Now().Format("20060102")
	var code string
	if err := s.pick(r).QueryRow("SELECT IFNULL(max(system_code),'') FROM chain_contract WHERE system_code LIKE ?", left+"%").Scan(&code); err != nil {
		return "", err
	}
	if code == "" {
		return left + "0000001", nil
	}
	return nextContractCode(code)
}
func nextContractCode(code string) (string, error) {
	if len(code) < 20 {
		return "", errors.New("invalid Java system-code prefix")
	}
	n, err := strconv.Atoi(code[20:])
	if err != nil {
		return "", err
	}
	return code[:20] + autoGenericCode(strconv.Itoa(n+1), 7), nil
}

// callBatch 直接复用已有批量业务处理函数，沿用当前请求身份，不发起网络调用。
func (s *Service) callBatch(r *http.Request, path, entity string, rows []map[string]any) (any, error) {
	if s.BatchMux == nil {
		return nil, errors.New("batch business handlers not connected")
	}
	e := authRestEnts[entity]
	body := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		body = append(body, e.shapeOne(row))
	}
	return s.callJSON(r, path, body)
}

// callJSON 沿用原 Feign 的身份和选库，复用单体内已注册的业务入口。
func (s *Service) callJSON(r *http.Request, path string, body any) (any, error) {
	if s.BatchMux == nil {
		return nil, errors.New("业务处理入口未连接")
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req := r.Clone(r.Context())
	url := *r.URL
	url.Path = path
	url.RawQuery = ""
	req.URL = &url
	req.Method = "POST"
	req.Body = io.NopCloser(bytes.NewReader(b))
	req.ContentLength = int64(len(b))
	req.Header = r.Header.Clone()
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.BatchMux.ServeHTTP(w, req)
	if w.Code >= 400 {
		return nil, fmt.Errorf("batch %s failed with HTTP %d", path, w.Code)
	}
	var result map[string]any
	d := json.NewDecoder(w.Body)
	d.UseNumber()
	if err := d.Decode(&result); err != nil {
		return nil, err
	}
	return result["data"], nil
}

type importSheet struct {
	entity, label, path string
	index               int
}

var allImportSheets = []importSheet{
	{"ChainCapitalSettlement", "结算单", "/saas-capital/chainCapitalSettlement/saveBatchSments", 3},
	{"ChainCapitalDetailed", "流水", "/saas-capital/chainCapitalDetailed/saveBatchDetails", 4},
	{"ChainBill", "进项发票", "/saas-bill/chainBill/saveBatch", 5},
	{"ChainBill", "销项发票", "/saas-bill/chainBill/saveBatch", 6},
	{"ChainBillRate", "税收", "/saas-bill/chainBillRate/saveBatch", 7},
	{"ChainLogistics", "运单", "/saas-logistics/chainLogistics/saveBatchLogistics", 1},
	{"ChainLogisticsWeighingList", "过磅单", "/saas-logistics/chainLogisticsWeighingList/saveBatchWeighingList", 2},
}

func (s *Service) contractImportAll(w http.ResponseWriter, r *http.Request) {
	fh, ok := formFile(r)
	if !ok {
		platform.Unknown(w)
		return
	}
	dict, err := s.contractDictionary(r)
	if err != nil {
		platform.Unknown(w)
		return
	}
	irt, err := s.importAllContractSheets(r, fh, dict)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, irt)
}

// contractImportAliasOnce 对照 Cg/Xs 控制器 importOnce(:191):先全量导入(service 事务),
// 再按 sheet 2–6 以 settlement/detailed/bill/bill/rate 实体二次读取并经 feign 批量保存
// (各自独立事务)。1:1 复刻错位二次导入的重复写入、errs/sucInfos/errInfos 形状,
// 含票据双重计数 bug 与税收 errInfos 非空时 for-each 自增触发 ConcurrentModificationException→500。
func (s *Service) contractImportAliasOnce(w http.ResponseWriter, r *http.Request) {
	fh, ok := formFile(r)
	if !ok {
		platform.Unknown(w)
		return
	}
	dict, err := s.contractDictionary(r)
	if err != nil {
		platform.Unknown(w)
		return
	}
	irt, err := s.importAllContractSheets(r, fh, dict)
	if err != nil {
		platform.Unknown(w)
		return
	}
	sucInfos := []any{}
	errs := []any{}
	if len(irt.ErrInfos) > 0 {
		errs = append(errs, "合同管理数据失败 0")
	} else {
		sucInfos = append(sucInfos, "合同管理数据成功 "+strconv.Itoa(len(irt.DataList)))
	}
	// 结算单 sheet2
	sirt := readExcelAt(fh, authRestEnts["ChainCapitalSettlement"].excel, 2, dict)
	if len(sirt.ErrInfos) > 0 {
		for _, msg := range sirt.ErrInfos {
			irt.ErrInfos = append(irt.ErrInfos, "结算单sheet中，"+msg)
		}
		errs = append(errs, "结算单管理数据失败 0")
	} else {
		data, err := s.callBatch(r, allImportSheets[0].path, "ChainCapitalSettlement", sirt.DataList)
		if err != nil {
			platform.Unknown(w)
			return
		}
		list, ok := data.([]any)
		if !ok {
			platform.Unknown(w) // Java (List)jrt.getData() 为 null → addAll NPE
			return
		}
		for _, e := range list {
			irt.ErrInfos = append(irt.ErrInfos, btsStr(e))
		}
		if len(list) == 0 {
			sucInfos = append(sucInfos, "结算单管理数据成功 "+strconv.Itoa(len(sirt.DataList)))
		} else {
			errs = append(errs, "结算单管理数据失败 0")
		}
	}
	// 流水 sheet3(feign 返回的 errInfos 不检查)
	dirt := readExcelAt(fh, authRestEnts["ChainCapitalDetailed"].excel, 3, dict)
	if len(dirt.ErrInfos) > 0 {
		for _, msg := range dirt.ErrInfos {
			irt.ErrInfos = append(irt.ErrInfos, "流水sheet中，"+msg)
		}
		errs = append(errs, "流水管理数据失败 0")
	} else {
		if _, err := s.callBatch(r, allImportSheets[1].path, "ChainCapitalDetailed", dirt.DataList); err != nil {
			platform.Unknown(w)
			return
		}
		sucInfos = append(sucInfos, "流水管理数据成功 "+strconv.Itoa(len(dirt.DataList)))
	}
	// 票据 sheet4+5
	birt1 := readExcelAt(fh, authRestEnts["ChainBill"].excel, 4, dict)
	birt2 := readExcelAt(fh, authRestEnts["ChainBill"].excel, 5, dict)
	if len(birt1.ErrInfos) > 0 || len(birt2.ErrInfos) > 0 {
		for _, msg := range birt1.ErrInfos {
			irt.ErrInfos = append(irt.ErrInfos, "进项发票sheet中，"+msg)
		}
		for _, msg := range birt2.ErrInfos {
			irt.ErrInfos = append(irt.ErrInfos, "销项发票sheet中，"+msg)
		}
		errs = append(errs, "票据管理数据失败 0")
	} else {
		bills := append(append([]map[string]any{}, birt1.DataList...), birt2.DataList...)
		if _, err := s.callBatch(r, "/saas-bill/chainBill/saveBatch", "ChainBill", bills); err != nil {
			platform.Unknown(w)
			return
		}
		// Java birt1.addAll(birt2) 后再 size()+birt2.size() → 双重计数,照实复刻。
		sucInfos = append(sucInfos, "票据管理数据成功 "+strconv.Itoa(len(bills)+len(birt2.DataList)))
	}
	// 税收 sheet6
	rirt := readExcelAt(fh, authRestEnts["ChainBillRate"].excel, 6, dict)
	if len(rirt.ErrInfos) > 0 {
		// Java for(String s : rirt.getErrInfos()) rirt.getErrInfos().add(...) → CME → 全局 500。
		platform.Unknown(w)
		return
	}
	if _, err := s.callBatch(r, "/saas-bill/chainBillRate/saveBatch", "ChainBillRate", rirt.DataList); err != nil {
		platform.Unknown(w)
		return
	}
	sucInfos = append(sucInfos, "税收管理数据成功 "+strconv.Itoa(len(rirt.DataList)))
	irt.Errs = errs
	irt.SucInfos = sucInfos
	irt.DataList = nil
	platform.Success(w, irt)
}

func (s *Service) importAllContractSheets(r *http.Request, fh *multipart.FileHeader, dict map[string]string) (*importResult, error) {
	cg, err := s.contractCode(r, "1")
	if err != nil {
		return nil, err
	}
	xs, err := s.contractCode(r, "2")
	if err != nil {
		return nil, err
	}
	irt := readExcelAt(fh, authRestEnts["ChainContract"].excel, 0, dict)
	batch := uuidV4()
	errs := append([]string{}, irt.ErrInfos...)
	contracts := map[string]map[string]any{}
	seen := map[string]bool{}
	good := []map[string]any{}
	if len(irt.ErrInfos) == 0 {
		for _, row := range irt.DataList {
			code := btsStr(row["contract_code"])
			contracts[code] = row
			if seen[code] {
				errs = append(errs, "合同sheet中，合同编号为“"+code+"”的数据重复")
				continue
			}
			var n int
			if err := s.pick(r).QueryRow("SELECT count(*) FROM chain_contract WHERE contract_code=?", row["contract_code"]).Scan(&n); err != nil {
				return nil, err
			}
			if n > 0 {
				errs = append(errs, "合同sheet中，合同编号为“"+code+"”的合同数据已存在，导入失败")
				continue
			}
			if row["contract_type"] == nil {
				return nil, errors.New("missing contract type")
			}
			switch int(num(row["contract_type"])) {
			case 1:
				row["system_code"] = cg
				cg, err = nextContractCode(cg)
			case 2:
				row["system_code"] = xs
				xs, err = nextContractCode(xs)
			}
			if err != nil {
				return nil, err
			}
			seen[code] = true
			row["im_batch_number"] = batch
			good = append(good, row)
		}
	}
	imports := make([]*importResult, len(allImportSheets))
	for i, sheet := range allImportSheets {
		result := readExcelAt(fh, authRestEnts[sheet.entity].excel, sheet.index, dict)
		imports[i] = result
		for _, msg := range result.ErrInfos {
			errs = append(errs, sheet.label+"sheet中，"+msg)
		}
		if len(result.ErrInfos) > 0 {
			continue
		}
		for _, row := range result.DataList {
			row["im_batch_number"] = batch
			switch sheet.entity {
			case "ChainCapitalSettlement":
				matches, err := queryRows(s.pick(r), "SELECT * FROM chain_contract WHERE contract_code=?", row["contract_code"])
				if err != nil || len(matches) > 1 {
					if err == nil {
						err = errors.New("multiple contracts for settlement")
					}
					return nil, err
				}
				if len(matches) == 0 {
					if contracts[btsStr(row["contract_code"])] == nil {
						errs = append(errs, "结算单sheet中，合同编号为“"+btsStr(row["contract_code"])+"”的合同数据不存在，导入失败")
					}
				} else {
					if err := bindSettlementContract(row, matches[0]); err != nil {
						return nil, err
					}
				}
			case "ChainBill":
				var n int
				if err := s.pick(r).QueryRow("SELECT count(*) FROM chain_bill WHERE bill_code=? AND bill_number=?", row["bill_code"], row["bill_number"]).Scan(&n); err != nil {
					return nil, err
				}
				if n > 0 {
					errs = append(errs, "代码为"+btsStr(row["bill_code"])+"，号码为"+btsStr(row["bill_number"])+"的发票已存在，保存中止")
				}
			case "ChainBillRate":
				var n int
				if err := s.pick(r).QueryRow("SELECT count(*) FROM chain_bill_rate WHERE month <=> ?", row["month"]).Scan(&n); err != nil {
					return nil, err
				}
				if n > 0 {
					errs = append(errs, "税收sheet中，所属月份为"+btsStr(row["month"])+"的数据重复")
				}
			}
		}
	}
	irt.ErrInfos = errs
	irt.Errs = []any{}
	irt.SucInfos = []any{}
	if len(errs) > 0 {
		irt.Errs = append(irt.Errs, "导入失败 0")
		irt.DataList = authRestEnts["ChainContract"].shape(irt.DataList)
		return irt, nil
	}
	// 保留 Java 本地合同事务边界；跨服务子处理通过独立连接读取已提交数据。
	if err := platform.WithTx(s.pick(r), func(tx *sql.Tx) error {
		if err := s.saveContractRowsTx(r, tx, good); err != nil {
			return err
		}
		for _, row := range imports[0].DataList {
			if row["contract_id"] == nil {
				if err := bindSettlementContract(row, contracts[btsStr(row["contract_code"])]); err != nil {
					return err
				}
			}
		}
		irt.SucInfos = append(irt.SucInfos, "合同管理数据成功 "+strconv.Itoa(len(irt.DataList)))
		bills := append(imports[2].DataList, imports[3].DataList...)
		jobs := []struct {
			path, entity, label string
			rows                []map[string]any
		}{
			{allImportSheets[0].path, "ChainCapitalSettlement", "结算单管理数据成功 ", imports[0].DataList},
			{allImportSheets[1].path, "ChainCapitalDetailed", "流水管理数据成功 ", imports[1].DataList},
			{allImportSheets[2].path, "ChainBill", "票据管理数据成功 ", bills},
			{allImportSheets[4].path, "ChainBillRate", "税收管理数据成功 ", imports[4].DataList},
			{allImportSheets[5].path, "ChainLogistics", "物流运单数据成功 ", imports[5].DataList},
			{allImportSheets[6].path, "ChainLogisticsWeighingList", "物流过磅单数据成功 ", imports[6].DataList},
		}
		for i, job := range jobs {
			if i >= 4 && len(job.rows) == 0 {
				continue
			}
			if _, err := s.callBatch(r, job.path, job.entity, job.rows); err != nil {
				return err
			}
			irt.SucInfos = append(irt.SucInfos, job.label+strconv.Itoa(len(job.rows)))
		}
		u := currentUser(r)
		operation := fmt.Sprintf("一键导入成功，批次号为：%s；导入合同%d条；导入结算单%d条；导入流水%d条；导入票据%d条；导入缴税%d条；", batch, len(irt.DataList), len(imports[0].DataList), len(imports[1].DataList), len(bills), len(imports[4].DataList))
		if _, err := platform.TxInsertNoFill(tx, "chain_common_log", map[string]any{"user_id": strconv.FormatInt(u.UserID, 10), "user_name": u.Username, "operation": operation, "operation_type": 8}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	irt.DataList = authRestEnts["ChainContract"].shape(irt.DataList)
	return irt, nil
}

func bindSettlementContract(row, contract map[string]any) error {
	if contract == nil || contract["contract_name"] == nil || row["type"] == nil {
		return errors.New("missing settlement contract data")
	}
	row["contract_id"] = contract["id"]
	row["contract_name"] = contract["contract_name"]
	row["supplier_id"] = contract["supplier_id"]
	row["supplier_name"] = contract["sales_unit"]
	if int(num(row["type"])) != 1 {
		row["supplier_name"] = contract["buyer_unit"]
	}
	row["khxje"] = row["sum_price"]
	return nil
}

func (s *Service) contractImportStep(w http.ResponseWriter, r *http.Request) {
	fh, ok := formFile(r)
	if !ok {
		platform.Unknown(w)
		return
	}
	dict, err := s.contractDictionary(r)
	if err != nil {
		platform.Unknown(w)
		return
	}
	irt := readExcelAt(fh, authRestEnts["ChainContract"].excel, 0, dict)
	cg, err := s.contractCode(r, "1")
	if err != nil {
		platform.Unknown(w)
		return
	}
	xs, err := s.contractCode(r, "2")
	if err != nil {
		platform.Unknown(w)
		return
	}
	good := []map[string]any{}
	seen := map[string]bool{}
	if len(irt.ErrInfos) == 0 && len(irt.DataList) > 0 {
		for _, row := range irt.DataList {
			code := btsStr(row["contract_code"])
			var n int
			if err = s.pick(r).QueryRow("SELECT count(*) FROM chain_contract WHERE contract_code=?", row["contract_code"]).Scan(&n); err != nil {
				platform.Unknown(w)
				return
			}
			if seen[code] {
				irt.ErrInfos = append(irt.ErrInfos, "合同sheet中，合同编号为“"+code+"”的数据重复")
				continue
			}
			if n > 0 {
				irt.ErrInfos = append(irt.ErrInfos, "合同sheet中，合同编号为“"+code+"”的合同数据已存在")
				continue
			}
			if row["contract_type"] == nil {
				platform.Unknown(w)
				return
			}
			if int(num(row["contract_type"])) == 1 {
				row["system_code"] = cg
				cg, err = nextContractCode(cg)
			} else if int(num(row["contract_type"])) == 2 {
				row["system_code"] = xs
				xs, err = nextContractCode(xs)
			}
			if err != nil {
				platform.Unknown(w)
				return
			}
			seen[code] = true
			good = append(good, row)
		}
		if len(irt.ErrInfos) == 0 {
			if err = s.saveContractRows(r, good); err != nil {
				platform.Unknown(w)
				return
			}
		}
	}
	irt.SucInfos = []any{}
	irt.Errs = []any{}
	if len(irt.ErrInfos) > 0 {
		irt.Errs = append(irt.Errs, "合同管理数据失败 0")
	} else {
		irt.SucInfos = append(irt.SucInfos, "合同管理数据成功 "+strconv.Itoa(len(irt.DataList)))
	}
	steps := []importSheet{{"ChainCapitalSettlement", "结算单", allImportSheets[0].path, 2}, {"ChainCapitalDetailed", "流水", allImportSheets[1].path, 3}, {"ChainBill", "进项发票", allImportSheets[2].path, 4}, {"ChainBill", "销项发票", allImportSheets[3].path, 5}, {"ChainBillRate", "税收", allImportSheets[4].path, 6}}
	for i := 0; i < len(steps); i++ {
		step := steps[i]
		parsed := readExcelAt(fh, authRestEnts[step.entity].excel, step.index, dict)
		label := step.label
		extraCount := 0
		if step.entity == "ChainBill" {
			other := readExcelAt(fh, authRestEnts["ChainBill"].excel, 5, dict)
			for _, msg := range other.ErrInfos {
				parsed.ErrInfos = append(parsed.ErrInfos, "销项发票sheet中，"+msg)
			}
			extraCount = len(other.DataList)
			parsed.DataList = append(parsed.DataList, other.DataList...)
			label = "票据"
			i++
		}
		if len(parsed.ErrInfos) > 0 {
			for _, msg := range parsed.ErrInfos {
				irt.ErrInfos = append(irt.ErrInfos, step.label+"sheet中，"+msg)
			}
			irt.Errs = append(irt.Errs, label+"管理数据失败 0")
			if step.entity == "ChainBillRate" {
				platform.Unknown(w)
				return
			}
			continue
		}
		data, err := s.callBatch(r, step.path, step.entity, parsed.DataList)
		if err != nil {
			platform.Unknown(w)
			return
		}
		if step.entity == "ChainCapitalSettlement" || step.entity == "ChainBill" {
			items, ok := data.([]any)
			if !ok {
				platform.Unknown(w)
				return
			}
			if len(items) > 0 {
				for _, v := range items {
					irt.ErrInfos = append(irt.ErrInfos, fmt.Sprint(v))
				}
				irt.Errs = append(irt.Errs, label+"管理数据失败 0")
				continue
			}
		}
		irt.SucInfos = append(irt.SucInfos, label+"管理数据成功 "+strconv.Itoa(len(parsed.DataList)+extraCount))
	}
	irt.DataList = nil
	platform.Success(w, irt)
}
