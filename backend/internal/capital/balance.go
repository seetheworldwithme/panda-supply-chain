// capital - balance.go 金额核销链路(balancingAccounts/cancelBalancing)。
// Java:ChainCapitalDetailedServiceImpl.balancingAccounts/cancelBalancing。
// 无事务:中途异常留半写,照实复刻(999 前已落库部分不回滚)。
package capital

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// dtBalancingAccounts 核销:hxbh=HXJL. 编号;details/sments 批更 khxje-=sjhxje、
// hxCode+=","+hxbh、hxStatus 3/2;detail_settlement 按 (contract,statement,detail) 累加;
// counts 填 hxCode 后批插;最后 save 核销记录。BigDecimal null → NPE 999。
func (s *Service) dtBalancingAccounts(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	u := currentUser(r)
	hxbh, ok := s.countSysCode(d, u)
	if !ok {
		platform.Unknown(w)
		return
	}
	details, _ := m["details"].([]any)
	sments, _ := m["sments"].([]any)
	counts, _ := m["counts"].([]any)
	if details == nil || sments == nil {
		platform.Unknown(w) // getDetails()/getSments() null → for NPE → 999
		return
	}
	var sumPrice decimal
	// 批量改流水
	for _, it := range details {
		row, ok := it.(map[string]any)
		if !ok {
			platform.Unknown(w)
			return
		}
		id, ok := int64Of(row, "id")
		if !ok {
			platform.Unknown(w) // MP updateBatchById 无 id → 999
			return
		}
		khxje, ok1 := decOf(row, "khxje") // null → NPE → 999
		sjhxje, ok2 := decOf(row, "sjhxje")
		if !ok1 || !ok2 {
			platform.Unknown(w)
			return
		}
		rest := decSub(khxje, sjhxje)
		hxCode := "null"
		if v, has := row["hxCode"]; has && v != nil {
			hxCode = str(row, "hxCode")
		}
		vals := map[string]any{
			"khxje": rest.s, "hx_code": hxCode + "," + hxbh,
			"hx_status": json.Number("2"),
		}
		if decCmp(rest, "0") == 0 {
			vals["hx_status"] = json.Number("3")
		}
		if _, err := platform.UpdateByIDNoFill(d, "chain_capital_detailed", id, vals); err != nil {
			platform.Unknown(w)
			return
		}
	}
	// 批量改结算单
	for _, it := range sments {
		row, ok := it.(map[string]any)
		if !ok {
			platform.Unknown(w)
			return
		}
		id, ok := int64Of(row, "id")
		if !ok {
			platform.Unknown(w)
			return
		}
		khxje, ok1 := decOf(row, "khxje")
		sjhxje, ok2 := decOf(row, "sjhxje")
		if !ok1 || !ok2 {
			platform.Unknown(w)
			return
		}
		rest := decSub(khxje, sjhxje)
		hxCode := "null"
		if v, has := row["hxCode"]; has && v != nil {
			hxCode = str(row, "hxCode")
		}
		vals := map[string]any{
			"khxje": rest.s, "hx_code": hxCode + "," + hxbh,
			"hx_status": json.Number("2"),
		}
		if decCmp(rest, "0") == 0 {
			vals["hx_status"] = json.Number("3")
		}
		if _, err := platform.UpdateByIDNoFill(d, "chain_capital_settlement", id, vals); err != nil {
			platform.Unknown(w)
			return
		}
		sumPrice = decAdd(sumPrice, sjhxje)
	}
	// 合同/结算单/流水对应
	for _, sit := range sments {
		st, ok := sit.(map[string]any)
		if !ok {
			platform.Unknown(w)
			return
		}
		stID, ok1 := int64Of(st, "id")
		stContract, ok2 := int64Of(st, "contractId")
		if !ok1 || !ok2 {
			platform.Unknown(w) // contractId null → 构造器 Integer 拆箱/插入语义 → 999
			return
		}
		for _, dit := range details {
			dt, ok := dit.(map[string]any)
			if !ok {
				platform.Unknown(w)
				return
			}
			dtID, ok1 := int64Of(dt, "id")
			dtSjhxje, ok2 := decOf(dt, "sjhxje")
			if !ok1 || !ok2 {
				platform.Unknown(w)
				return
			}
			var existID int64
			var existMoney sql.NullString
			err := d.QueryRow("SELECT id, money FROM chain_capital_detail_settlement WHERE contract_id = ? AND statement_id = ? AND detail_id = ?",
				stContract, stID, dtID).Scan(&existID, &existMoney)
			if err == sql.ErrNoRows {
				if _, err := platform.InsertNoFill(d, "chain_capital_detail_settlement", map[string]any{
					"contract_id":  json.Number(strconv.FormatInt(stContract, 10)),
					"statement_id": json.Number(strconv.FormatInt(stID, 10)),
					"detail_id":    json.Number(strconv.FormatInt(dtID, 10)),
					"money":        dtSjhxje.s, "hx_code": hxbh,
				}); err != nil {
					platform.Unknown(w)
					return
				}
			} else if err != nil {
				platform.Unknown(w)
				return
			} else {
				newMoney := decAdd(NewDecimal(existMoney.String), dtSjhxje)
				if _, err := platform.UpdateByIDNoFill(d, "chain_capital_detail_settlement", existID, map[string]any{
					"money": newMoney.s,
				}); err != nil {
					platform.Unknown(w)
					return
				}
			}
		}
	}
	// 核销明细批插
	for _, it := range counts {
		mx, ok := it.(map[string]any)
		if !ok {
			platform.Unknown(w)
			return
		}
		vals := pick(mx, entCountMx.cols...)
		delete(vals, "id")
		delete(vals, "create_time")
		delete(vals, "update_time")
		vals["hx_code"] = hxbh
		if _, err := platform.InsertNoFill(d, "chain_capital_count_mx", vals); err != nil {
			platform.Unknown(w)
			return
		}
	}
	// 核销记录
	decreaseType := m["decreaseType"]
	recVals := map[string]any{
		"system_code": hxbh,
		"user_id":     strconv.FormatInt(u.UserID, 10),
		"user_name":   u.Username,
		"sment_num":   json.Number(strconv.Itoa(len(sments))),
		"detail_num":  json.Number(strconv.Itoa(len(details))),
		"money":       sumPrice.s,
	}
	if decreaseType != nil {
		recVals["type"] = decreaseType
	}
	if _, err := platform.InsertNoFill(d, "chain_capital_count", recVals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "核销成功")
}

// dtCancelBalancing 取消核销:List<Map<String,String>>;allEq(parm) 查明细还原金额,
// 删明细/核销记录/对应表。khxFlag=true 是载体键不入列(只有 khxje/hx_code 更新)。
func (s *Service) dtCancelBalancing(w http.ResponseWriter, r *http.Request) {
	list, ok := bodyList(r)
	if !ok {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	for _, it := range list {
		parm, ok := it.(map[string]any)
		if !ok {
			platform.Unknown(w)
			return
		}
		cond, ok := allEqCond(parm, entCountMx.cols) // 值 null 跳过;非列键 → 999
		if !ok {
			platform.Unknown(w)
			return
		}
		mxes, err := platform.ListWhere(d, "chain_capital_count_mx", cond)
		if err != nil {
			platform.Unknown(w)
			return
		}
		for _, mx := range mxes {
			if mx["detail_id"] != nil {
				dtID := toInt64(mx["detail_id"])
				var khxje, detailMoney sql.NullString
				if err := d.QueryRow("SELECT khxje FROM chain_capital_detailed WHERE id = ?", dtID).Scan(&khxje); err != nil {
					platform.Unknown(w) // selectById 查无 → MyBatis null → khxje NPE → 999
					return
				}
				detailMoney.String = btsStr(mx["detail_money"])
				if detailMoney.String == "" {
					platform.Unknown(w) // add(null) NPE → 999
					return
				}
				if khxje.String == "" {
					platform.Unknown(w)
					return
				}
				newKhxje := decAdd(NewDecimal(khxje.String), NewDecimal(detailMoney.String))
				if _, err := platform.UpdateByIDNoFill(d, "chain_capital_detailed", dtID, map[string]any{
					"khxje": newKhxje.s, "hx_code": "2",
				}); err != nil {
					platform.Unknown(w)
					return
				}
			}
			if mx["statement_id"] != nil {
				stID := toInt64(mx["statement_id"])
				var khxje, settMoney sql.NullString
				if err := d.QueryRow("SELECT khxje FROM chain_capital_settlement WHERE id = ?", stID).Scan(&khxje); err != nil {
					platform.Unknown(w)
					return
				}
				settMoney.String = btsStr(mx["sett_money"])
				if settMoney.String == "" || khxje.String == "" {
					platform.Unknown(w)
					return
				}
				newKhxje := decAdd(NewDecimal(khxje.String), NewDecimal(settMoney.String))
				if _, err := platform.UpdateByIDNoFill(d, "chain_capital_settlement", stID, map[string]any{
					"khxje": newKhxje.s, "hx_code": "2",
				}); err != nil {
					platform.Unknown(w)
					return
				}
			}
			if err := platform.DeleteByID(d, "chain_capital_count_mx", toInt64(mx["id"])); err != nil {
				platform.Unknown(w)
				return
			}
		}
		if hxCode, has := parm["hx_code"]; has && hxCode != nil {
			if err := platform.DeleteByMap(d, "chain_capital_count", map[string]any{"system_code": hxCode}); err != nil {
				platform.Unknown(w)
				return
			}
		} else {
			// allEq({system_code:null}) → MP 对 null 值键的 remove 语义:单参 allEq(Map) 默认跳过 null 键?
			// MP 源码:allEq(Map) = allEq(map, false) → null 值键跳过 → 全表删除!照实复刻(危险但保真)。
			if _, err := d.Exec("DELETE FROM chain_capital_count"); err != nil {
				platform.Unknown(w)
				return
			}
		}
		if err := platform.DeleteByMap(d, "chain_capital_detail_settlement", cond); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, "取消成功")
}

// countSysCode 核销编号(HXJL. 前缀,balancingAccounts 用)
func (s *Service) countSysCode(d *sql.DB, u *auth.CurrentUser) (string, bool) {
	like := "HXJL." + u.EnterpriseCode + strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", "")
	return systemCodeNext(d, "chain_capital_count", like, "1")
}

// ── BigDecimal 字面算术(字符串精度保真,add/sub/cmp) ──

type decimal struct{ s string }

func NewDecimal(s string) decimal { return decimal{s: strings.TrimSpace(s)} }

func (d decimal) String() string { return d.s }

func decOf(m map[string]any, k string) (decimal, bool) {
	v, ok := m[k]
	if !ok || v == nil {
		return decimal{}, false
	}
	switch t := v.(type) {
	case json.Number:
		return decimal{s: t.String()}, true
	case string:
		if strings.TrimSpace(t) == "" {
			return decimal{}, false
		}
		return decimal{s: t}, true
	default:
		return decimal{s: str(m, k)}, true
	}
}

func decScale(s string) (int, string) {
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = s[1:]
	}
	intPart, frac := s, ""
	if i := strings.Index(s, "."); i >= 0 {
		intPart, frac = s[:i], s[i+1:]
	}
	return len(frac), map[bool]string{true: "-", false: ""}[neg] + intPart + frac
}

func decAlign(a, b decimal) (int, int, int) {
	sa, da := decScale(a.s)
	sb, db := decScale(b.s)
	if sa < sb {
		da += strings.Repeat("0", sb-sa)
		return sb, toInt(da), toInt(db)
	}
	if sb < sa {
		db += strings.Repeat("0", sa-sb)
	}
	return maxInt(sa, sb), toInt(da), toInt(db)
}

func decAdd(a, b decimal) decimal {
	sc, ia, ib := decAlign(a, b)
	sum := ia + ib
	return formatDec(sc, sum)
}

func decSub(a, b decimal) decimal {
	sc, ia, ib := decAlign(a, b)
	return formatDec(sc, ia-ib)
}

// decCmp a 与 b(字面)比较:-1/0/1
func decCmp(a decimal, b string) int {
	sc, ia, ib := decAlign(a, decimal{s: b})
	_ = sc
	if ia < ib {
		return -1
	}
	if ia > ib {
		return 1
	}
	return 0
}

func formatDec(scale, v int) decimal {
	neg := ""
	if v < 0 {
		neg = "-"
		v = -v
	}
	s := strconv.Itoa(v)
	for len(s) <= scale {
		s = "0" + s
	}
	if scale == 0 {
		return decimal{s: neg + s}
	}
	return decimal{s: neg + s[:len(s)-scale] + "." + s[len(s)-scale:]}
}

func toInt(s string) int {
	n, _ := strconv.Atoi(s)
	if s == "" {
		return 0
	}
	return n
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func toInt64(v any) int64 {
	switch t := v.(type) {
	case json.Number:
		n, _ := t.Int64()
		return n
	case int64:
		return t
	case string:
		n, _ := strconv.ParseInt(t, 10, 64)
		return n
	}
	return 0
}
