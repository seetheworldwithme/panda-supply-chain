// bill - excelfields.go 票据域 Excel 字段表(Java @ExcelColumn 声明序)与导入导出流程。
package bill

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

type impField struct {
	header  string
	col     string
	kind    string // s/i/d/t/b/dd
	sys     bool
	dict    bool
	notNull bool
}

var billFields = []impField{
	{"id", "id", "i", false, false, false},
	{"导入批次号", "im_batch_number", "s", false, false, false},
	{"企业id", "company_id", "i", false, false, false},
	{"系统编号", "system_code", "s", true, false, false},
	{"合同id", "contract_id", "i", false, false, false},
	{"合同编号", "contract_code", "s", false, false, false},
	{"合同名称", "contract_name", "s", false, false, false},
	{"进项类型", "in_type", "s", false, false, false},
	{"发票代码", "bill_code", "s", false, false, false},
	{"发票号码", "bill_number", "s", false, false, false},
	{"开票日期", "open_date", "dd", false, false, false},
	{"税号", "duty_paragraph", "s", false, false, false},
	{"企业名称", "company_name", "s", false, false, false},
	{"商品名称", "goods_name", "s", false, false, false},
	{"单位", "uit", "s", false, false, false},
	{"数量", "count", "d", false, false, false},
	{"单价(不含税)", "unit_price", "d", false, false, false},
	{"金额", "base_money", "d", false, false, false},
	{"税率", "tax_rate", "s", false, false, false},
	{"税额", "tax_amount", "d", false, false, false},
	{"含税金额", "includ_money", "d", false, false, false},
	{"规格型号", "modle", "s", false, false, false},
	{"所属月份", "month", "s", false, false, false},
	{"类型(1、进项发票2、销项发票)", "type", "i", false, false, false},
	{"是否已关联", "relation_flag", "b", false, false, false},
	{"关联id", "relation_id", "i", false, false, false},
	{"创建时间", "create_time", "t", false, false, false},
	{"修改时间", "update_time", "t", false, false, false},
	{"备注字段1", "remark1", "s", false, false, false},
	{"备注字段2", "remark2", "s", false, false, false},
	{"备注字段3", "remark3", "s", false, false, false},
	{"备注字段4", "remark4", "s", false, false, false},
	{"备注字段5", "remark5", "s", false, false, false},
}

var rateFields = []impField{
	{"id", "id", "i", false, false, false},
	{"导入批次号", "im_batch_number", "s", false, false, false},
	{"企业id", "company_id", "i", false, false, false},
	{"系统编号", "system_code", "s", false, false, false},
	{"所属月份", "month", "s", false, false, false},
	{"不含税金额", "base_money", "d", false, false, false},
	{"销项税额", "tax_amount", "d", false, false, false},
	{"销项含税金额", "includ_money", "d", false, false, false},
	{"进项税额", "jx_money", "d", false, false, false},
	{"增值税额", "zz_money", "d", false, false, false},
	{"城建税", "cj_money", "d", false, false, false},
	{"教育费附加税额", "jy_money", "d", false, false, false},
	{"地方教育费附加", "dfjy_money", "d", false, false, false},
	{"印花税额", "yh_money", "d", false, false, false},
	{"企业所得税", "qy_money", "d", false, false, false},
	{"滞纳金", "znj_money", "d", false, false, false},
	{"合计", "sum_money", "d", false, false, false},
	{"备注", "remark", "s", false, false, false},
	{"创建时间", "create_time", "t", false, false, false},
	{"修改时间", "update_time", "t", false, false, false},
	{"备注字段1", "remark1", "s", false, false, false},
	{"备注字段2", "remark2", "s", false, false, false},
	{"备注字段3", "remark3", "s", false, false, false},
	{"备注字段4", "remark4", "s", false, false, false},
	{"备注字段5", "remark5", "s", false, false, false},
}

var settRecFields2 = []impField{
	{"id", "id", "i", false, false, false},
	{"企业id", "company_id", "i", false, false, false},
	{"结算单id", "statement_id", "i", false, false, false},
	{"票据id", "bill_id", "i", false, false, false},
	{"创建时间", "create_time", "t", false, false, false},
	{"修改时间", "update_time", "t", false, false, false},
	{"备注字段1", "remark1", "s", false, false, false},
	{"备注字段2", "remark2", "s", false, false, false},
	{"备注字段3", "remark3", "s", false, false, false},
	{"备注字段4", "remark4", "s", false, false, false},
	{"备注字段5", "remark5", "s", false, false, false},
}

func init() {
	entBill.fields = billFields
	entRate.fields = rateFields
	entSettRow.fields = settRecFields2
}

// ── Bill 导入(Jx/Xx):readExcel 4 参(idt 编号仅组装给 readExcel,但注入仅在 5 参版本
// 发生——8765 实测导入行 system_code 为 NULL?否:Bill 有 @SystemCodeColumn 且走
// service.importExcel(idt) 组装编号后 readExcel(class,file,systemCode,0)。四参版本不注入
// (同 Detailed 实证),saveIrt 再按 type 重新分配 system_code → 最终入库编号来自 saveIrt) ──

func (s *Service) billImportExcel(w http.ResponseWriter, r *http.Request, p string) {
	idt := "PJJX."
	if strings.HasSuffix(p, "/importExcelXx") {
		idt = "PJXX."
	}
	d := s.pick(r)
	u := currentUser(r)
	prefix := idt + u.EnterpriseCode + time.Now().Format("20060102")
	_, ok := systemCodeNext(d, "chain_bill", prefix, "0")
	if !ok {
		platform.Unknown(w)
		return
	}
	fh, ok2 := formFile(r)
	if !ok2 {
		platform.Unknown(w)
		return
	}
	irt := readExcelFile(fh, billFields, "", true, emptyDict())
	s.billSaveIrtImport(w, r, irt)
}

// billSaveIrtImport 对照 BillService.saveIrt(ImportResult):查重+回填+编号+批次+日志 type2;
// 响应完整 ImportResult(无 fail 分支)
func (s *Service) billSaveIrtImport(w http.ResponseWriter, r *http.Request, irt *importResult) {
	d := s.pick(r)
	u := currentUser(r)
	var errInfos []string
	var sucInfos []any
	if len(irt.ErrInfos) > 0 {
		for _, e := range irt.ErrInfos {
			errInfos = append(errInfos, "发票sheet中，"+e)
		}
	} else {
		batch := uuidV4()
		for _, row := range irt.DataList {
			row["im_batch_number"] = batch
		}
		cg, ok1 := s.billSysCode(d, u, "PJJX.")
		xs, ok2 := s.billSysCode(d, u, "PJXX.")
		if !ok1 || !ok2 {
			platform.Unknown(w)
			return
		}
		var kept []map[string]any
		for _, row := range irt.DataList {
			var dup int
			_ = d.QueryRow("SELECT count(*) FROM chain_bill WHERE bill_code = ? AND bill_number = ?",
				nullOrNil(row["bill_code"]), nullOrNil(row["bill_number"])).Scan(&dup)
			if dup > 0 {
				errInfos = append(errInfos, "代码为"+btsStr(row["bill_code"])+"，号码为"+btsStr(row["bill_number"])+"的发票已存在，保存中止")
				continue
			}
			if v, has := row["contract_code"]; has && v != nil {
				var cid sqlNullInt64
				var cname sqlNullString
				if err := d.QueryRow("SELECT id, contract_name FROM chain_contract WHERE contract_code = ?",
					btsStr(v)).Scan(&cid, &cname); err == nil {
					row["contract_id"] = json.Number(strconv.FormatInt(cid.Int64, 10))
					row["contract_name"] = cname.String
				}
			}
			switch btsStr(row["type"]) {
			case "1":
				row["system_code"] = cg
				if n, okb := nextSystemCodeInMem(cg); okb {
					cg = n
				}
			case "2":
				row["system_code"] = xs
				if n, okb := nextSystemCodeInMem(xs); okb {
					xs = n
				}
			}
			kept = append(kept, row)
		}
		if len(errInfos) == 0 {
			for i, row := range kept {
				vals := insertVals(row)
				delete(vals, "id")
				newID, err := platform.InsertNoFill(d, "chain_bill", vals)
				if err != nil {
					platform.Unknown(w)
					return
				}
				irt.DataList[i] = echoOf(&entBill, platform.Camelize(row), map[string]any{
					"id": json.Number(strconv.FormatInt(newID, 10)),
				})
			}
			logOp(d, u, "导入成功，批次号为："+batch+"；导入发票"+strconv.Itoa(len(kept))+"条；", 2)
			sucInfos = []any{"票据管理数据成功 " + strconv.Itoa(len(irt.DataList))}
		}
	}
	out := importResult{DataList: irt.DataList, SucInfos: sucInfos, Errs: []any{}, ErrInfos: errInfos}
	if out.DataList == nil {
		out.DataList = []map[string]any{}
	}
	if out.SucInfos == nil {
		out.SucInfos = []any{}
	}
	if out.ErrInfos == nil {
		out.ErrInfos = []string{}
	}
	platform.Success(w, out)
}

type sqlNullInt64 = struct {
	Int64 int64
	Valid bool
}

type sqlNullString = struct {
	String string
	Valid  bool
}

func (s *Service) billExportExcel(w http.ResponseWriter, r *http.Request) {
	s.exportHandler(w, r, "chain_bill", billFields)
}

// ── BillRate 导入:readExcel+saveIrt(ImportResult):month 查重+逐条 saveBill ──

func (s *Service) rateImportExcel(w http.ResponseWriter, r *http.Request) {
	fh, ok := formFile(r)
	if !ok {
		platform.Unknown(w)
		return
	}
	irt := readExcelFile(fh, rateFields, "", true, emptyDict())
	s.rateSaveIrtImport(w, r, irt)
}

func (s *Service) rateSaveIrtImport(w http.ResponseWriter, r *http.Request, irt *importResult) {
	d := s.pick(r)
	u := currentUser(r)
	var errInfos []string
	var sucInfos []any
	errs := []any{}
	if len(irt.ErrInfos) > 0 {
		for _, e := range irt.ErrInfos {
			errInfos = append(errInfos, "税收sheet中，"+e)
		}
	} else {
		batch := uuidV4()
		for _, row := range irt.DataList {
			row["im_batch_number"] = batch
			var dup int
			_ = d.QueryRow("SELECT count(*) FROM chain_bill_rate WHERE month = ?", btsStr(row["month"])).Scan(&dup)
			if dup > 0 {
				errInfos = append(errInfos, "税收sheet中，所属月份为"+btsStr(row["month"])+"的数据重复")
				continue
			}
		}
		if len(errInfos) == 0 {
			code, ok1 := s.rateSysCode(d, u)
			if !ok1 {
				platform.Unknown(w)
				return
			}
			for i, row := range irt.DataList {
				row["system_code"] = code
				n, okb := nextSystemCodeInMem(code)
				if !okb {
					platform.Unknown(w)
					return
				}
				code = n
				vals := insertVals(row)
				delete(vals, "id")
				newID, err := platform.InsertNoFill(d, "chain_bill_rate", vals)
				if err != nil {
					platform.Unknown(w)
					return
				}
				irt.DataList[i] = echoOf(&entRate, platform.Camelize(row), map[string]any{
					"id": json.Number(strconv.FormatInt(newID, 10)),
				})
				logOp(d, u, "新增税收数据\n对比结果为：\n", 7)
			}
			sucInfos = []any{"缴税管理数据成功 " + strconv.Itoa(len(irt.DataList))}
			logOp(d, u, "导入成功，批次号为："+batch+"；导入缴税明细"+strconv.Itoa(len(irt.DataList))+"条；", 6)
		} else {
			errs = []any{"导入失败 0"}
		}
	}
	out := importResult{DataList: irt.DataList, SucInfos: sucInfos, Errs: errs, ErrInfos: errInfos}
	if out.DataList == nil {
		out.DataList = []map[string]any{}
	}
	if out.SucInfos == nil {
		out.SucInfos = []any{}
	}
	if out.ErrInfos == nil {
		out.ErrInfos = []string{}
	}
	// errInfos>0 时未插库行保持解析实体全键(未插库行仍是 snake 原始键 → 实体整形)
	for i, row := range out.DataList {
		if row != nil {
			if _, has := row["id"]; !has {
				if _, okSnake := row["month"]; okSnake {
					out.DataList[i] = echoOf(&entRate, platform.Camelize(row), nil)
				}
			}
		}
	}
	platform.Success(w, out)
}

func (s *Service) rateExportExcel(w http.ResponseWriter, r *http.Request) {
	s.exportHandler(w, r, "chain_bill_rate", rateFields)
}
