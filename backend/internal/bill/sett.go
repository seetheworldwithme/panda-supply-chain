// bill - sett.go /chainBillSettlement 发票/结算单对应表(9 端点)。
// importExcel:selectSystemCode 查 chain_bill_settlement.system_code——该表无此列 →
// SQL 报错恒 999(同四期 tw/lc 实证)。listByPage 的 ${searchStr} null 内插空串。
package bill

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func (s *Service) handleSett(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	d := s.pick(r)
	switch {
	case strings.HasSuffix(p, "/list"):
		m := bodyMap(r)
		cond, ok := allEqCond(m, entSettRow.cols)
		if !ok {
			platform.Unknown(w)
			return
		}
		rows, err := platform.ListWhere(d, "chain_bill_settlement", cond)
		if err != nil {
			platform.Unknown(w)
			return
		}
		if rows == nil {
			rows = []map[string]any{}
		}
		platform.Success(w, entSettRow.rows(rows))
	case strings.HasSuffix(p, "/listByPage"):
		m := bodyMap(r)
		if _, ok := m["page"].(map[string]any); !ok {
			platform.Unknown(w)
			return
		}
		ss := "%%" // MyBatis ${null} 输出空串
		if v, ok := m["searchStr"]; ok && v != nil {
			ss = "%" + str(m, "searchStr") + "%"
		}
		current, size := pageOf(m)
		q := "select * from chain_bill_settlement where CONCAT(id) like '" + ss + "' order by create_time desc"
		countQ := "SELECT count(1) FROM chain_bill_settlement where CONCAT(id) like '" + ss + "'"
		ip, err := paged(d, q, countQ, nil, current, size)
		if err != nil {
			platform.Unknown(w)
			return
		}
		recs, _ := ip["records"].([]map[string]any)
		ip["records"] = entSettRow.rows(recs)
		platform.Success(w, ip)
	case strings.HasSuffix(p, "/getById"):
		m := bodyMap(r)
		var row map[string]any
		if id, ok := int64Of(m, "id"); ok {
			var err error
			row, err = platform.GetByID(d, "chain_bill_settlement", id)
			if err != nil {
				platform.Unknown(w)
				return
			}
		}
		if row == nil {
			platform.Success(w, nil)
			return
		}
		platform.Success(w, entSettRow.rows([]map[string]any{row})[0])
	case strings.HasSuffix(p, "/add"):
		m := bodyMap(r)
		if !validateDT(m, []string{"create_time"}, []string{}, []string{"update_time"}) {
			platform.Unknown(w)
			return
		}
		vals := writeVals(&entSettRow, m)
		delete(vals, "id")
		delete(vals, "create_time")
		delete(vals, "update_time")
		id, err := platform.InsertNoFill(d, "chain_bill_settlement", vals)
		if err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, echoOf(&entSettRow, m, map[string]any{
			"id": json.Number(strconv.FormatInt(id, 10)),
		}))
	case strings.HasSuffix(p, "/deleteById"):
		m := bodyMap(r)
		id, ok := int64Of(m, "id")
		if !ok {
			platform.Unknown(w)
			return
		}
		if err := platform.DeleteByID(d, "chain_bill_settlement", id); err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, "删除成功，id为："+str(m, "id"))
	case strings.HasSuffix(p, "/updateById"):
		m := bodyMap(r)
		if !validateDT(m, []string{"create_time"}, []string{}, []string{"update_time"}) {
			platform.Unknown(w)
			return
		}
		id, ok := int64Of(m, "id")
		if !ok {
			platform.Unknown(w)
			return
		}
		vals := writeVals(&entSettRow, m)
		delete(vals, "id")
		delete(vals, "create_time")
		delete(vals, "update_time")
		if _, err := platform.UpdateByIDNoFill(d, "chain_bill_settlement", id, vals); err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, echoOf(&entSettRow, m, nil))
	case strings.HasSuffix(p, "/importExcel"):
		// selectSystemCode 表无 system_code 列 → SQL 报错恒 999
		platform.Unknown(w)
	case strings.HasSuffix(p, "/exportExcel"):
		s.exportHandler(w, r, "chain_bill_settlement", entSettRow.fields)
	case strings.HasSuffix(p, "/delBatch"):
		m := bodyMap(r)
		ids, ok := m["ids"].([]any)
		if !ok {
			platform.Success(w, "删除成功") // MP removeByIds null 短路
			return
		}
		if len(ids) == 0 {
			platform.Success(w, "删除成功")
			return
		}
		numIds := make([]any, 0, len(ids))
		for _, it := range ids {
			n, ok := toID(it)
			if !ok {
				platform.Unknown(w)
				return
			}
			numIds = append(numIds, json.Number(strconv.FormatInt(n, 10)))
		}
		if err := platform.DeleteByIDs(d, "chain_bill_settlement", numIds); err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, "删除成功")
	default:
		platform.Unknown(w)
	}
}
