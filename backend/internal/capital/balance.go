// capital - balance.go 金额核销链路(balancingAccounts/cancelBalancing)。
// Java:ChainCapitalDetailedServiceImpl.balancingAccounts/cancelBalancing。
// 经 TxAdviceInterceptor("*"→REQUIRED) 整方法一个事务,中途异常整体回滚。
package capital

import (
	"database/sql"
	"encoding/json"
	"errors"
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
	details, _ := m["details"].([]any)
	sments, _ := m["sments"].([]any)
	counts, _ := m["counts"].([]any)
	if details == nil || sments == nil {
		platform.Unknown(w) // getDetails()/getSments() null → for NPE → 999
		return
	}
	// TxAdviceInterceptor 整方法一个事务。
	err := platform.WithTx(d, func(tx *sql.Tx) error {
		hxbh, ok := s.countSysCodeTx(tx, u)
		if !ok {
			return errors.New("system code")
		}
		var sumPrice decimal
		// 批量改流水(updateBatchById:SET 实体全部非 null 列+计算覆盖)
		for _, it := range details {
			row, ok := it.(map[string]any)
			if !ok {
				return errors.New("bad detail")
			}
			khxje, ok1 := decOf(row, "khxje") // null → NPE → 999
			sjhxje, ok2 := decOf(row, "sjhxje")
			if !ok1 || !ok2 {
				return errors.New("khxje/sjhxje null NPE")
			}
			rest := decSub(khxje, sjhxje)
			hxCode := "null"
			if v, has := row["hxCode"]; has && v != nil {
				hxCode = str(row, "hxCode")
			}
			vals := insertVals(pick(row, entDetailed.cols...))
			delete(vals, "id")
			vals["khxje"] = rest.s
			vals["hx_code"] = hxCode + "," + hxbh
			vals["hx_status"] = json.Number("2")
			if decCmp(rest, "0") == 0 {
				vals["hx_status"] = json.Number("3")
			}
			if id, ok := int64Of(row, "id"); ok {
				if err := platform.TxUpdateByIDNoFill(tx, "chain_capital_detailed", id, vals); err != nil {
					return err
				}
			} // id null → WHERE id=null 命中 0 行,不报错
		}
		// 批量改结算单
		for _, it := range sments {
			row, ok := it.(map[string]any)
			if !ok {
				return errors.New("bad sment")
			}
			khxje, ok1 := decOf(row, "khxje")
			sjhxje, ok2 := decOf(row, "sjhxje")
			if !ok1 || !ok2 {
				return errors.New("khxje/sjhxje null NPE")
			}
			rest := decSub(khxje, sjhxje)
			hxCode := "null"
			if v, has := row["hxCode"]; has && v != nil {
				hxCode = str(row, "hxCode")
			}
			vals := insertVals(pick(row, entSettlement.cols...))
			delete(vals, "id")
			vals["khxje"] = rest.s
			vals["hx_code"] = hxCode + "," + hxbh
			vals["hx_status"] = json.Number("2")
			if decCmp(rest, "0") == 0 {
				vals["hx_status"] = json.Number("3")
			}
			if id, ok := int64Of(row, "id"); ok {
				if err := platform.TxUpdateByIDNoFill(tx, "chain_capital_settlement", id, vals); err != nil {
					return err
				}
			}
			sumPrice = decAdd(sumPrice, sjhxje)
		}
		// 合同/结算单/流水对应
		for _, sit := range sments {
			st, ok := sit.(map[string]any)
			if !ok {
				return errors.New("bad sment")
			}
			stID, ok1 := int64Of(st, "id")
			stContract, ok2 := int64Of(st, "contractId")
			if !ok1 || !ok2 {
				return errors.New("contractId null")
			}
			for _, dit := range details {
				dt, ok := dit.(map[string]any)
				if !ok {
					return errors.New("bad detail")
				}
				dtID, ok1 := int64Of(dt, "id")
				dtSjhxje, ok2 := decOf(dt, "sjhxje")
				if !ok1 || !ok2 {
					return errors.New("detail id/sjhxje null")
				}
				var existID int64
				var existMoney sql.NullString
				err := tx.QueryRow("SELECT id, money FROM chain_capital_detail_settlement WHERE contract_id = ? AND statement_id = ? AND detail_id = ?",
					stContract, stID, dtID).Scan(&existID, &existMoney)
				if err == sql.ErrNoRows {
					if _, err := platform.TxInsertNoFill(tx, "chain_capital_detail_settlement", map[string]any{
						"contract_id":  json.Number(strconv.FormatInt(stContract, 10)),
						"statement_id": json.Number(strconv.FormatInt(stID, 10)),
						"detail_id":    json.Number(strconv.FormatInt(dtID, 10)),
						"money":        dtSjhxje.s, "hx_code": hxbh,
					}); err != nil {
						return err
					}
				} else if err != nil {
					return err
				} else {
					newMoney := decAdd(NewDecimal(existMoney.String), dtSjhxje)
					if err := platform.TxUpdateByIDNoFill(tx, "chain_capital_detail_settlement", existID, map[string]any{
						"money": newMoney.s,
					}); err != nil {
						return err
					}
				}
			}
		}
		// 核销明细批插(saveBatch:实体非 null 列全写,hxCode 覆盖)
		for _, it := range counts {
			mx, ok := it.(map[string]any)
			if !ok {
				return errors.New("bad count")
			}
			vals := insertVals(pick(mx, entCountMx.cols...))
			delete(vals, "id")
			vals["hx_code"] = hxbh
			if _, err := platform.TxInsertNoFill(tx, "chain_capital_count_mx", vals); err != nil {
				return err
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
		if _, err := platform.TxInsertNoFill(tx, "chain_capital_count", recVals); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
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
	// Java cancelBalancing 经 TxAdviceInterceptor 整方法一个事务;allEq(parm) 单参
	// → null2IsNull=true(显式 null → IS NULL,MP 3.4.2 字节码核实)。
	err := platform.WithTx(d, func(tx *sql.Tx) error {
		for _, it := range list {
			parm, ok := it.(map[string]any)
			if !ok {
				return errors.New("bad item")
			}
			cond, ok := allEqCondNull(parm, entCountMx.cols) // 保留 null;非列键 → SQL 异常 999
			if !ok {
				return errors.New("unknown column")
			}
			where, wargs := whereClause(cond)
			rows, err := tx.Query("SELECT * FROM chain_capital_count_mx"+where, wargs...)
			if err != nil {
				return err
			}
			mxes, err := platform.ScanMaps(rows)
			if err != nil {
				return err
			}
			for _, mx := range mxes {
				if mx["detail_id"] != nil {
					dtID := toInt64(mx["detail_id"])
					var khxje, detailMoney sql.NullString
					if err := tx.QueryRow("SELECT khxje FROM chain_capital_detailed WHERE id = ?", dtID).Scan(&khxje); err != nil {
						return err // selectById 查无 → MyBatis null → khxje NPE → 999
					}
					detailMoney.String = btsStr(mx["detail_money"])
					if detailMoney.String == "" || khxje.String == "" {
						return errors.New("add(null) NPE")
					}
					newKhxje := decAdd(NewDecimal(khxje.String), NewDecimal(detailMoney.String))
					if err := platform.TxUpdateByIDNoFill(tx, "chain_capital_detailed", dtID, map[string]any{
						"khxje": newKhxje.s, "hx_code": "2",
					}); err != nil {
						return err
					}
				}
				if mx["statement_id"] != nil {
					stID := toInt64(mx["statement_id"])
					var khxje, settMoney sql.NullString
					if err := tx.QueryRow("SELECT khxje FROM chain_capital_settlement WHERE id = ?", stID).Scan(&khxje); err != nil {
						return err
					}
					settMoney.String = btsStr(mx["sett_money"])
					if settMoney.String == "" || khxje.String == "" {
						return errors.New("add(null) NPE")
					}
					newKhxje := decAdd(NewDecimal(khxje.String), NewDecimal(settMoney.String))
					if err := platform.TxUpdateByIDNoFill(tx, "chain_capital_settlement", stID, map[string]any{
						"khxje": newKhxje.s, "hx_code": "2",
					}); err != nil {
						return err
					}
				}
				if _, err := tx.Exec("DELETE FROM chain_capital_count_mx WHERE id = ?", toInt64(mx["id"])); err != nil {
					return err
				}
			}
			// countParms 恒含 system_code 键,值=parm.get("hx_code")(可 null → IS NULL)。
			if err := platform.TxDeleteByMap(tx, "chain_capital_count", map[string]any{"system_code": parm["hx_code"]}); err != nil {
				return err
			}
			if err := platform.TxDeleteByMap(tx, "chain_capital_detail_settlement", cond); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "取消成功")
}

// allEqCondNull MP allEq(Map) 单参语义:键原样作列名,显式 null 保留(→IS NULL)。
func allEqCondNull(m map[string]any, cols []string) (map[string]any, bool) {
	allow := make(map[string]bool, len(cols))
	for _, c := range cols {
		allow[c] = true
	}
	cond := map[string]any{}
	for k, v := range m {
		if !allow[k] {
			return nil, false
		}
		cond[k] = v
	}
	return cond, true
}

// whereClause cond → WHERE(null → IS NULL),空 map 无 WHERE。键须经白名单校验。
func whereClause(cond map[string]any) (string, []any) {
	if len(cond) == 0 {
		return "", nil
	}
	cols := platform.SortedKeys(cond)
	var conds []string
	var args []any
	for _, k := range cols {
		if cond[k] == nil {
			conds = append(conds, k+" IS NULL")
		} else {
			conds = append(conds, k+" = ?")
			args = append(args, cond[k])
		}
	}
	return " WHERE " + strings.Join(conds, " AND "), args
}

// countSysCodeTx 核销编号(HXJL. 前缀,balancingAccounts 用,事务内)。
func (s *Service) countSysCodeTx(tx *sql.Tx, u *auth.CurrentUser) (string, bool) {
	like := "HXJL." + u.EnterpriseCode + strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", "")
	var cur sql.NullString
	if err := tx.QueryRow("SELECT IFNULL(max(system_code),'') FROM chain_capital_count WHERE system_code LIKE ?", like+"%").Scan(&cur); err != nil {
		return "", false
	}
	if cur.String == "" {
		return like + autoGenericCode("1", 7), true
	}
	if len(cur.String) < 20 {
		return "", false
	}
	n, err := strconv.Atoi(cur.String[20:])
	if err != nil {
		return "", false
	}
	return cur.String[:20] + autoGenericCode(strconv.Itoa(n+1), 7), true
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
