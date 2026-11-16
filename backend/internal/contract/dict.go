// dict.go 三期补充:chainCommonDict 基础维护字典(合同页面 filters 依赖 sessionStorage 字典,
// 来自 /saas-common/chainCommonDict/list)。Excel 导入导出 2 端点债务未实现。
// 实体 @TableId(AUTO);getDType 类 getter → Jackson 键 dtype/dkey/dvalue/dtypeName/...
// (与 chain_contract_extend 同款全小写);custFlag/dextend 为 Boolean,createTime 有 @JsonFormat 空格格式
package contract

import (
	"net/http"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

var dictCols = []string{"id", "parent_id", "d_type", "d_type_name", "child_key", "d_key",
	"d_value", "d_sel_type", "d_base_default", "d_template_default", "d_cancel", "d_extend",
	"cust_flag", "create_time", "update_time", "remark1", "remark2", "remark3", "remark4", "remark5"}

// dictRename 列 → JSON 键(getDXxx 全小写现象,实测键序:...dkey,dvalue,dcancel,dextend,dtype,dtypeName,...)
var dictRename = map[string]string{"d_type": "dtype", "d_type_name": "dtypeName",
	"d_key": "dkey", "d_value": "dvalue", "d_sel_type": "dselType",
	"d_base_default": "dbaseDefault", "d_template_default": "dtemplateDefault",
	"d_cancel": "dcancel", "d_extend": "dextend"}

var dictCarrier = []string{"searchStr", "page", "ids"}

var dictEchoKeys = func() []string {
	out := []string{"id", "parentId", "childKey", "custFlag", "createTime", "updateTime",
		"remark1", "remark2", "remark3", "remark4", "remark5"}
	return append(out, append(dictCarrier, "dkey", "dvalue", "dcancel", "dextend",
		"dtype", "dtypeName", "dselType", "dbaseDefault", "dtemplateDefault")...)
}()

// dictCamelRename Camelize 后的键 → Jackson 输出键(getDXxx 全小写现象)
var dictCamelRename = map[string]string{"dType": "dtype", "dTypeName": "dtypeName",
	"dKey": "dkey", "dValue": "dvalue", "dSelType": "dselType",
	"dBaseDefault": "dbaseDefault", "dTemplateDefault": "dtemplateDefault",
	"dCancel": "dcancel", "dExtend": "dextend"}

// shapeDict 行形状:rename 小写键 + Boolean 化(cust_flag/d_extend;NULL 保持 null)
func shapeDict(rows []map[string]any) []map[string]any {
	out := platform.CamelizeRows(rows)
	for _, r := range out {
		for from, to := range dictCamelRename {
			if v, ok := r[from]; ok {
				delete(r, from)
				r[to] = v
			}
		}
		for _, c := range []string{"custFlag", "dextend", "dcancel", "dbaseDefault", "dtemplateDefault"} {
			if v, ok := r[c]; ok && v != nil {
				r[c] = num(v) != 0
			}
		}
		fillNulls(r, dictCarrier)
	}
	return out
}

// registerDict 挂 /saas-common/chainCommonDict 端点(eq 由 Register 提供)
func (s *Service) registerDict(eq func(string, http.HandlerFunc)) {
	eq("/chainCommonDict/list", s.dictList)
	eq("/chainCommonDict/listByPage", s.dictListByPage)
	eq("/chainCommonDict/getById", s.dictGetById)
	eq("/chainCommonDict/add", s.dictAdd)
	eq("/chainCommonDict/deleteById", s.dictDeleteById)
	eq("/chainCommonDict/updateById", s.dictUpdateById)
	eq("/chainCommonDict/delBatch", s.dictDelBatch)
}

// dictList = findList:QueryWrapper.allEq(maps,false)
func (s *Service) dictList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cond, ok := allEqCond(m, dictCols)
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := platform.ListWhere(s.pick(r), "chain_common_dict", cond)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if rows == nil {
		rows = []map[string]any{}
	}
	platform.Success(w, shapeDict(rows))
}

// dictListByPage = mapper:CONCAT(id) like '%${searchStr}%'(null → 字面 "null" 永假空页)
func (s *Service) dictListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	current, size := pageOf(m)
	if v, ok := m["searchStr"]; !ok || v == nil {
		platform.Unknown(w) // ${searchStr} null 直拼 → Java 实测 999(行为保持)
		return
	}
	ss := "%" + str(m, "searchStr") + "%"
	q := "SELECT * FROM chain_common_dict WHERE CONCAT(id) like ? order by create_time desc"
	countQ := "SELECT count(*) FROM chain_common_dict WHERE CONCAT(id) like ?"
	page, err := paged(s.pick(r), q, countQ, []any{ss}, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	page["records"] = shapeDict(page["records"].([]map[string]any))
	platform.Success(w, page)
}

func (s *Service) dictGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	row, err := platform.GetByID(s.pick(r), "chain_common_dict", m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	if row == nil {
		platform.Success(w, nil)
		return
	}
	platform.Success(w, shapeDict([]map[string]any{row})[0])
}

func (s *Service) dictAdd(w http.ResponseWriter, r *http.Request) {
	s.writeAdd(w, r, "chain_common_dict", bodyMap(r), dictCols, dictRename, dictEchoKeys, false, nil)
}

func (s *Service) dictUpdateById(w http.ResponseWriter, r *http.Request) {
	s.writeUpdate(w, r, "chain_common_dict", bodyMap(r), dictCols, dictRename, dictEchoKeys, false)
}

func (s *Service) dictDeleteById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if _, err := s.pick(r).Exec("DELETE FROM chain_common_dict WHERE id = ?", m["id"]); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功，id为："+str(m, "id"))
}

func (s *Service) dictDelBatch(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	ids, _ := m["ids"].([]any)
	if err := platform.DeleteByIDs(s.pick(r), "chain_common_dict", ids); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功")
}
