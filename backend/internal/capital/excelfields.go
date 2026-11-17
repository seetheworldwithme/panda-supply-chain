// capital - excelfields.go 各实体 Excel 列(Java @ExcelColumn 声明序,表头原文)与
// Detailed/Settlement 的 importAlone/importExcel/exportExcel 流程。
package capital

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

type impField struct {
	header  string
	col     string // snake 列
	kind    string // s=String i=Integer d=BigDecimal t=LocalDateTime b=Boolean
	sys     bool   // SystemCodeColumn
	dict    bool   // DictColumn
	notNull bool   // NotBankColumn
}

var dtFields = []impField{
	{"id", "id", "i", false, false, false},
	{"导入批次号", "im_batch_number", "s", false, false, false},
	{"系统编号", "system_code", "s", true, false, false},
	{"企业id", "company_id", "i", false, false, false},
	{"账户名称", "account_name", "s", false, false, true},
	{"账户号码", "account", "s", false, false, true},
	{"交易日期", "trans_time", "dd", false, false, true},
	{"金额", "money", "d", false, false, true},
	{"对方账户", "opposite_name", "s", false, false, true},
	{"对方账号", "opposite_account", "s", false, false, false},
	{"摘要", "abs_content", "s", false, false, false},
	{"用途", "purpose", "s", false, false, false},
	{"收付标识(1、收款单2、付款单)", "detype", "s", false, false, true},
	{"可核销金额", "khxje", "d", false, false, false},
	{"实际核销金额", "sjhxje", "d", false, false, false},
	{"核销状态(1、未核销2、核销中3、已核销)", "hx_status", "i", false, false, false},
	{"是否已关联", "relation_flag", "b", false, false, false},
	{"银行名称", "bank_name", "s", false, false, true},
	{"核销编号", "hx_code", "s", false, false, false},
	{"创建时间", "create_time", "t", false, false, false},
	{"修改时间", "update_time", "t", false, false, false},
	{"备注字段1", "remark1", "s", false, false, false},
	{"备注字段2", "remark2", "s", false, false, false},
	{"备注字段3", "remark3", "s", false, false, false},
	{"备注字段4", "remark4", "s", false, false, false},
	{"备注字段5", "remark5", "s", false, false, false},
}

var stFields = []impField{
	{"id", "id", "i", false, false, false},
	{"导入批次号", "im_batch_number", "s", false, false, false},
	{"系统编号", "system_code", "s", false, false, false},
	{"企业id", "company_id", "i", false, false, false},
	{"合同id", "contract_id", "i", false, false, false},
	{"合同编号", "contract_code", "s", false, false, false},
	{"合同名称", "contract_name", "s", false, false, false},
	{"货物名称", "goods_name", "s", false, false, true},
	{"结算开始时间", "set_start_time", "dd", false, false, true},
	{"结算结束时间", "set_end_time", "dd", false, false, true},
	{"结算数量", "settlement_count", "d", false, false, true},
	{"结算单价", "unit_price", "d", false, false, true},
	{"结算总价", "sum_price", "d", false, false, true},
	{"可核销金额", "khxje", "d", false, false, false},
	{"实际核销金额", "sjhxje", "d", false, false, false},
	{"核销状态(1、未核销2、核销中3、已核销)", "hx_status", "i", false, false, false},
	{"结算类型（1、采购结算 2、销售结算）", "type", "i", false, false, true},
	{"客户/供应商id", "supplier_id", "i", false, false, false},
	{"客户/供应商名称", "supplier_name", "s", false, false, false},
	{"是否已关联流水", "rel_detail_flag", "b", false, false, false},
	{"是否已关联票据", "rel_bill_flag", "b", false, false, false},
	{"核销编号", "hx_code", "s", false, false, false},
	{"所属月份", "month", "s", false, false, false},
	{"关联id", "relation_id", "i", false, false, false},
	{"创建时间", "create_time", "t", false, false, false},
	{"修改时间", "update_time", "t", false, false, false},
	{"备注字段1", "remark1", "s", false, false, false},
	{"备注字段2", "remark2", "s", false, false, false},
	{"备注字段3", "remark3", "s", false, false, false},
	{"备注字段4", "remark4", "s", false, false, false},
	{"备注字段5", "remark5", "s", false, false, false},
}

var settRecFields = []impField{
	{"id", "id", "i", false, false, false},
	{"企业id", "company_id", "i", false, false, false},
	{"结算单数量", "sett_num", "i", false, false, false},
	{"发票数量", "bill_num", "i", false, false, false},
	{"操作人", "operater", "s", false, false, false},
	{"结算类型（1、采购结算 2、销售结算）", "type", "i", false, false, false},
	{"创建时间", "create_time", "t", false, false, false},
	{"修改时间", "update_time", "t", false, false, false},
	{"备注字段1", "remark1", "s", false, false, false},
	{"备注字段2", "remark2", "s", false, false, false},
	{"备注字段3", "remark3", "s", false, false, false},
	{"备注字段4", "remark4", "s", false, false, false},
	{"备注字段5", "remark5", "s", false, false, false},
}

var accountsFields = []impField{
	{"id", "id", "i", false, false, false},
	{"企业id", "company_id", "i", false, false, false},
	{"账户名称", "account_name", "s", false, false, false},
	{"账户号码", "account_number", "s", false, false, false},
	{"所属类型(1、自己2、对方)", "owner_type", "i", false, false, false},
	{"创建时间", "create_time", "t", false, false, false},
	{"修改时间", "update_time", "t", false, false, false},
	{"备注字段1", "remark1", "s", false, false, false},
	{"备注字段2", "remark2", "s", false, false, false},
	{"备注字段3", "remark3", "s", false, false, false},
	{"备注字段4", "remark4", "s", false, false, false},
	{"备注字段5", "remark5", "s", false, false, false},
}

var bkcallFields = []impField{
	{"", "id", "i", false, false, false},
	{"消息号", "msg_nbr", "s", false, false, false},
	{"消息类型", "msg_type", "s", false, false, false},
	{"消息键值", "msg_key", "s", false, false, false},
	{"平台编号", "plt_nbr", "s", false, false, false},
	{"商户编号", "mch_nbr", "s", false, false, false},
	{"订单号", "order_nbr", "s", false, false, false},
	{"订单交易序号", "order_sub_nbr", "s", false, false, false},
	{"商户订单号", "mch_order_nbr", "s", false, false, false},
	{"商户交易流水号", "mch_trans_nbr", "s", false, false, false},
	{"订单类型：PAY：消费类BCK：退货类", "order_type", "s", false, false, false},
	{"订单子类：DRT：直接支付GRT：担保支付BCK：退货退款", "order_subtype", "s", false, false, false},
	{"订单推送时间", "order_post_time", "s", false, false, false},
	{"订单金额", "order_amt", "d", false, false, false},
	{"交易金额", "trans_amt", "d", false, false, false},
	{"订单币种", "order_ccy_nbr", "s", false, false, false},
	{"订单状态", "status", "s", false, false, false},
	{"请求结果", "rtn_flag", "s", false, false, false},
	{"退款标志：Y：已退款N：未退款I：处理中", "bck_flag", "s", false, false, false},
	{"退款金额", "bck_amt", "d", false, false, false},
	{"订单开始时间", "order_str_time", "s", false, false, false},
	{"订单完成时间", "order_end_time", "s", false, false, false},
	{"付款方账号", "payer_acc_nbr", "s", false, false, false},
	{"付款方账户名", "payer_acc_name", "s", false, false, false},
	{"付款方联行号", "payer_bank_nbr", "s", false, false, false},
	{"付款方联行名", "payer_bank_name", "s", false, false, false},
	{"收款方账号", "payee_acc_nbr", "s", false, false, false},
	{"收款方账户名", "payee_acc_name", "s", false, false, false},
	{"收款方联行号", "payee_bank_nbr", "s", false, false, false},
	{"收款方联行名", "payee_bank_name", "s", false, false, false},
	{"消息描述", "msgText", "s", false, false, false},
	{"创建时间", "create_time", "t", false, false, false},
	{"修改时间", "update_time", "t", false, false, false},
	{"备注字段1", "remark1", "s", false, false, false},
	{"备注字段2", "remark2", "s", false, false, false},
	{"备注字段3", "remark3", "s", false, false, false},
	{"备注字段4", "remark4", "s", false, false, false},
	{"备注字段5", "remark5", "s", false, false, false},
}

var transFields = []impField{
	{"id", "id", "i", false, false, false},
	{"系统编号", "system_code", "s", false, false, false},
	{"订单号", "mch_order_nbr", "s", false, false, false},
	{"流水号", "mch_trans_nbr", "s", false, false, false},
	{"退款订单号", "bck_order_nbr", "s", false, false, false},
	{"退款流水号", "bck_trans_nbr", "s", false, false, false},
	{"企业id", "company_id", "i", false, false, false},
	{"付款单位id", "pay_unit_id", "i", false, false, false},
	{"收款单位id", "payee_unit_id", "i", false, false, false},
	{"付款单位名称", "pay_unit_name", "s", false, false, false},
	{"收款单位名称", "payee_unit_name", "s", false, false, false},
	{"支付方式（DRT：直接支付，GRT:担保支付）", "orderSubtype", "s", false, false, false},
	{"合同id", "contract_id", "i", false, false, false},
	{"合同编号", "contract_code", "s", false, false, false},
	{"付款户名", "pay_account_name", "s", false, false, false},
	{"付款账号", "pay_account", "s", false, false, false},
	{"付款金额", "pay_money", "d", false, false, false},
	{"收款户名", "payee_account_name", "s", false, false, false},
	{"收款账号", "payee_account", "s", false, false, false},
	{"收款商户编号", "mch_nbr", "s", false, false, false},
	{"状态（付款方：0、未付款1、已付款冻结中2、付款完成。收款方：1、待收款已冻结2、已收款。平台：0、未付款1、已付款冻结中2、付款完成）", "status", "i", false, false, false},
	{"审核状态（0、合同待确认1、合同已确认2、已付款待审核3、已付款已审核4、提现待审核5、提现已审核）", "approve_status", "i", false, false, false},
	{"招商订单状态(STR：发起OPR：已受理AUT：待付款WCF：待确认ING：处理中（退款）DBT：已扣款（退款）FIN：完成)", "bus_status", "s", false, false, false},
	{"行内外标志", "payer_bank_flag", "s", false, false, false},
	{"创建时间", "create_time", "t", false, false, false},
	{"修改时间", "update_time", "t", false, false, false},
	{"备注字段1", "remark1", "s", false, false, false},
	{"备注字段2", "remark2", "s", false, false, false},
	{"备注字段3", "remark3", "s", false, false, false},
	{"备注字段4", "remark4", "s", false, false, false},
	{"备注字段5", "remark5", "s", false, false, false},
}

var countFields = []impField{
	{"id", "id", "i", false, false, false},
	{"系统编号", "system_code", "s", false, false, false},
	{"企业id", "company_id", "i", false, false, false},
	{"用户id", "user_id", "s", false, false, false},
	{"用户名称", "user_name", "s", false, false, false},
	{"结算单数量", "sment_num", "i", false, false, false},
	{"银行流水数量", "detail_num", "i", false, false, false},
	{"类型（1、应收2、应付）", "type", "i", false, false, false},
	{"核销金额", "money", "d", false, false, false},
	{"创建时间", "create_time", "t", false, false, false},
	{"修改时间", "update_time", "t", false, false, false},
	{"备注字段1", "remark1", "s", false, false, false},
	{"备注字段2", "remark2", "s", false, false, false},
	{"备注字段3", "remark3", "s", false, false, false},
	{"备注字段4", "remark4", "s", false, false, false},
	{"备注字段5", "remark5", "s", false, false, false},
}

var countMxFields = []impField{
	{"id", "id", "i", false, false, false},
	{"企业id", "company_id", "i", false, false, false},
	{"结算单id", "statement_id", "i", false, false, false},
	{"结算单核销金额", "sett_money", "d", false, false, false},
	{"收付款单id", "detail_id", "i", false, false, false},
	{"收付款单核销金额", "detail_money", "d", false, false, false},
	{"核销编号", "hx_code", "s", false, false, false},
	{"创建时间", "create_time", "t", false, false, false},
	{"修改时间", "update_time", "t", false, false, false},
	{"备注字段1", "remark1", "s", false, false, false},
	{"备注字段2", "remark2", "s", false, false, false},
	{"备注字段3", "remark3", "s", false, false, false},
	{"备注字段4", "remark4", "s", false, false, false},
	{"备注字段5", "remark5", "s", false, false, false},
}

var detailSettFields = []impField{
	{"企业id", "company_id", "i", false, false, false},
	{"合同id", "contract_id", "i", false, false, false},
	{"结算单id", "statement_id", "i", false, false, false},
	{"收付款单id", "detail_id", "i", false, false, false},
	{"核销金额", "money", "d", false, false, false},
	{"核销编号", "hx_code", "s", false, false, false},
	{"修改时间", "update_time", "t", false, false, false},
	{"备注字段1", "remark1", "s", false, false, false},
	{"备注字段2", "remark2", "s", false, false, false},
	{"备注字段3", "remark3", "s", false, false, false},
	{"备注字段4", "remark4", "s", false, false, false},
	{"备注字段5", "remark5", "s", false, false, false},
}

// ── Detailed 导入/导出 ──

// dtImportAlone readExcel(file,"",0) → saveIrt(ImportResult):批次 UUID+systemCode 分配+
// 日志;响应完整 ImportResult
func (s *Service) dtImportAlone(w http.ResponseWriter, r *http.Request) {
	fh, ok := formFile(r)
	if !ok {
		irt := importResult{DataList: []map[string]any{}, ErrInfos: []string{"上传文件为空"}}
		platform.Success(w, irt)
		return
	}
	irt := readExcelFile(fh, dtFields, "", true, emptyDict())
	s.detailedSaveIrt(w, r, irt)
}

// dtImportExcel Sk/Fk:service.importExcel(idt) 组织编号 → 有错 fail(errInfos);
// 成功 saveBatch → "导入成功，共N条"
func (s *Service) dtImportExcel(w http.ResponseWriter, r *http.Request, p string) {
	idt := "SFSK."
	if p == "/saas-capital/chainCapitalDetailed/importExcelFk" ||
		p == "/saas-capital/chainCapitalCgDetailed/importExcelFk" ||
		p == "/saas-capital/chainCapitalXsDetailed/importExcelFk" {
		idt = "SFFK."
	}
	d := s.pick(r)
	u := currentUser(r)
	prefix := idt + u.EnterpriseCode + time.Now().Format("20060102")
	syscode, ok := systemCodeNext(d, "chain_capital_detailed", prefix, "0")
	if !ok {
		platform.Unknown(w)
		return
	}
	fh, ok2 := formFile(r)
	if !ok2 {
		platform.Unknown(w)
		return
	}
	// readExcel 4 参版本(class,file,systemCode,0):多 sheet 且不注入编号(8765 实测
	// 导入行 system_code 为 NULL,注入仅发生在 5 参版本);syscode 查询只为编号语义
	_, _ = syscode, ok
	irt := readExcelFile(fh, dtFields, "", true, emptyDict())
	if len(irt.ErrInfos) > 0 {
		platform.FailData(w, 999, "失败", irt.ErrInfos)
		return
	}
	for _, row := range irt.DataList {
		vals := insertVals(row)
		delete(vals, "id")
		if _, err := platform.InsertNoFill(d, "chain_capital_detailed", vals); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, "导入成功，共"+strconv.Itoa(len(irt.DataList))+"条")
}

// detailedSaveIrt 对照 DetailedService.saveIrt(ImportResult):错误包装"流水sheet中，";
// 成功→批次+分配+日志,响应 ImportResult
func (s *Service) detailedSaveIrt(w http.ResponseWriter, r *http.Request, irt *importResult) {
	d := s.pick(r)
	u := currentUser(r)
	var errInfos []string
	if len(irt.ErrInfos) > 0 {
		for _, e := range irt.ErrInfos {
			errInfos = append(errInfos, "流水sheet中，"+e)
		}
	}
	var sucInfos []any
	var errs []any
	if len(errInfos) > 0 {
		errs = []any{"导入失败 0"}
	} else {
		cg, ok1 := s.detailedSysCode(d, u, "SFCG.")
		xs, ok2 := s.detailedSysCode(d, u, "SFXS.")
		if !ok1 || !ok2 {
			platform.Unknown(w)
			return
		}
		batch := uuidV4()
		for _, row := range irt.DataList {
			row["im_batch_number"] = batch
			switch str(row, "detype") {
			case "1":
				row["system_code"] = xs
				if n, okb := nextSystemCodeInMem(xs); okb {
					xs = n
				}
			case "2":
				row["system_code"] = cg
				if n, okb := nextSystemCodeInMem(cg); okb {
					cg = n
				}
			}
		}
		for i, row := range irt.DataList {
			vals := insertVals(row)
			delete(vals, "id")
			newID, err := platform.InsertNoFill(d, "chain_capital_detailed", vals)
			if err != nil {
				platform.Unknown(w)
				return
			}
			irt.DataList[i] = echoOf(&entDetailed, platform.Camelize(row), map[string]any{
				"id": json.Number(strconv.FormatInt(newID, 10)),
			})
		}
		sucInfos = []any{"流水管理数据成功 " + strconv.Itoa(len(irt.DataList))}
		logOp(d, u, "导入成功，批次号为："+batch+"；导入流水"+strconv.Itoa(len(irt.DataList))+"条；", 3)
	}
	out := importResult{DataList: irt.DataList, SucInfos: sucInfos, Errs: errs, ErrInfos: errInfos}
	if out.DataList == nil {
		out.DataList = []map[string]any{}
	}
	if out.SucInfos == nil {
		out.SucInfos = []any{}
	}
	if out.Errs == nil {
		out.Errs = []any{}
	}
	if out.ErrInfos == nil {
		out.ErrInfos = []string{}
	}
	platform.Success(w, out)
}

// dtSaveBatch(List) 的重定向:saveIrt(List) 语义(无批次号)
// (dtSaveBatch 在 detailed.go)

// dtExportExcel list() 全表 → writeExcel
func (s *Service) dtExportExcel(w http.ResponseWriter, r *http.Request) {
	s.exportHandler(w, r, "chain_capital_detailed", dtFields)
}

// ── Settlement 导入/导出 ──

func (s *Service) stImportAlone(w http.ResponseWriter, r *http.Request) {
	fh, ok := formFile(r)
	if !ok {
		irt := importResult{DataList: []map[string]any{}, ErrInfos: []string{"上传文件为空"}}
		platform.Success(w, irt)
		return
	}
	irt := readExcelFile(fh, stFields, "", true, emptyDict())
	s.settlementSaveIrt(w, r, irt)
}

// stImportExcel service.importExcel(idt="") → 有错 fail;成功 saveIrt(List) → "导入成功，共N条"
func (s *Service) stImportExcel(w http.ResponseWriter, r *http.Request) {
	fh, ok := formFile(r)
	if !ok {
		platform.Unknown(w)
		return
	}
	irt := readExcelFile(fh, stFields, "", true, emptyDict())
	if len(irt.ErrInfos) > 0 {
		platform.FailData(w, 999, "失败", irt.ErrInfos)
		return
	}
	s.settlementSaveIrtRows(w, r, irt.DataList, true)
}

// settlementSaveIrtRows saveIrt(List):type 1→JSCG. 2→JSXS.;withResp=true 输出
// "导入成功，共N条"(importExcel 路径)
func (s *Service) settlementSaveIrtRows(w http.ResponseWriter, r *http.Request, rows []map[string]any, withResp bool) {
	d := s.pick(r)
	u := currentUser(r)
	cg, ok1 := s.settSysCode(d, u, "JSCG.")
	xs, ok2 := s.settSysCode(d, u, "JSXS.")
	if !ok1 || !ok2 {
		platform.Unknown(w)
		return
	}
	for _, row := range rows {
		vals := insertVals(row)
		delete(vals, "id")
		switch str(row, "type") {
		case "1":
			vals["system_code"] = cg
			if n, okb := nextSystemCodeInMem(cg); okb {
				cg = n
			}
		case "2":
			vals["system_code"] = xs
			if n, okb := nextSystemCodeInMem(xs); okb {
				xs = n
			}
		}
		if _, err := platform.InsertNoFill(d, "chain_capital_settlement", vals); err != nil {
			platform.Unknown(w)
			return
		}
	}
	if withResp {
		platform.Success(w, "导入成功，共"+strconv.Itoa(len(rows))+"条")
	}
}

// settlementSaveIrt 对照 SettlementService.saveIrt(ImportResult):"结算单sheet中，"包装
func (s *Service) settlementSaveIrt(w http.ResponseWriter, r *http.Request, irt *importResult) {
	d := s.pick(r)
	u := currentUser(r)
	var errInfos []string
	if len(irt.ErrInfos) > 0 {
		for _, e := range irt.ErrInfos {
			errInfos = append(errInfos, "结算单sheet中，"+e)
		}
	}
	var sucInfos []any
	var errs []any
	if len(errInfos) > 0 {
		errs = []any{"导入失败 0"}
	} else {
		cg, ok1 := s.settSysCode(d, u, "JSCG.")
		xs, ok2 := s.settSysCode(d, u, "JSXS.")
		if !ok1 || !ok2 {
			platform.Unknown(w)
			return
		}
		batch := uuidV4()
		for _, row := range irt.DataList {
			row["im_batch_number"] = batch
			switch str(row, "type") {
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
		}
		for i, row := range irt.DataList {
			vals := insertVals(row)
			delete(vals, "id")
			newID, err := platform.InsertNoFill(d, "chain_capital_settlement", vals)
			if err != nil {
				platform.Unknown(w)
				return
			}
			irt.DataList[i] = echoOf(&entSettlement, platform.Camelize(row), map[string]any{
				"id": json.Number(strconv.FormatInt(newID, 10)),
			})
		}
		sucInfos = []any{"结算管理数据成功 " + strconv.Itoa(len(irt.DataList))}
		logOp(d, u, "导入成功，批次号为："+batch+"；导入结算"+strconv.Itoa(len(irt.DataList))+"条；", 4)
	}
	out := importResult{DataList: irt.DataList, SucInfos: sucInfos, Errs: errs, ErrInfos: errInfos}
	if out.DataList == nil {
		out.DataList = []map[string]any{}
	}
	if out.SucInfos == nil {
		out.SucInfos = []any{}
	}
	if out.Errs == nil {
		out.Errs = []any{}
	}
	if out.ErrInfos == nil {
		out.ErrInfos = []string{}
	}
	platform.Success(w, out)
}

func (s *Service) stExportExcel(w http.ResponseWriter, r *http.Request) {
	s.exportHandler(w, r, "chain_capital_settlement", stFields)
}

func init() {
	entDetailed.fields = dtFields
	entSettlement.fields = stFields
	entBillSettRec.fields = settRecFields
	entAccounts.fields = accountsFields
	entBkcall.fields = bkcallFields
	entTransinfo.fields = transFields
	entCount.fields = countFields
	entCountMx.fields = countMxFields
	entDetailSettlement.fields = detailSettFields
}
