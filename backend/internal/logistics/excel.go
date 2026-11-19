// logistics - excel.go Excel 导入导出(标准库 zip+xml 实现,无三方依赖)。
// 导出对照 ExcelUtils.writeExcel:全部 @ExcelColumn 字段按声明序,标题行+数据行,
// 单元格一律字符串(LocalDateTime.toString() ISO T、BigDecimal 保小数位、null 空)。
// 导入对照 ExcelUtils.readExcel:表头包含匹配、getCellValue 数字/日期/百分比转换、
// NotBankColumn 必填(notNullMap 首行后耗尽的原版 bug)、DictColumn 校验(翻译实际不生效,
// 非数字值落 0)、SystemCodeColumn 逐行注入、4 参版本遍历全部 sheet、5 参版本只取 sheet0。
// 行级异常 Java 在 forEach 内吞掉(字段留空行保留),非整行丢弃——照实复刻。
// ponytail: .xls(HSSF 二进制)不支持——Go 版 .xls 走"上传文件格式不正确"失败语义,
// 需要 .xls 时引入三方库再补;表头包含匹配用声明序近似 Java HashMap 迭代序(同标题多命中场景)。
package logistics

import (
	"archive/zip"
	"crypto/rand"
	"errors"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// impField Excel 列映射(header 与 Java @ExcelColumn value 逐字一致)
type impField struct {
	header  string
	col     string // snake 列
	kind    string // s=String i=Integer d=BigDecimal t=LocalDateTime b=Boolean
	sys     bool   // SystemCodeColumn
	dict    bool   // DictColumn
	notNull bool   // NotBankColumn
}

var (
	loFields = []impField{
		{"id", "id", "i", false, false, false},
		{"导入批次号", "im_batch_number", "s", false, false, false},
		{"批次号", "batch_number", "s", false, false, false},
		{"企业id", "company_id", "i", false, false, false},
		{"订单id", "order_id", "i", false, false, false},
		{"系统编号", "system_code", "s", true, false, false},
		{"运单编号", "bill_number", "s", false, false, true},
		{"时间", "trans_time", "t", false, false, true},
		{"车号", "car_number", "s", false, false, true},
		{"装货地", "loading_address", "s", false, false, true},
		{"卸货地", "unloading_point", "s", false, false, true},
		{"标的物", "subject_matter", "s", false, false, true},
		{"运费单价", "freight_unit", "d", false, false, true},
		{"运费金额", "freight_amount", "d", false, false, true},
		{"原发吨数", "original_tonnage", "d", false, false, true},
		{"实收吨数", "tonnage_paid", "d", false, false, true},
		{"应付金额", "pay_money", "d", false, false, true},
		{"承运平台", "carriage_platform", "i", false, true, false},
		{"附件ids", "file_ids", "s", false, false, false},
		{"运单创建时间", "lo_create_time", "t", false, false, false},
		{"装货地省份", "loading_province", "s", false, false, false},
		{"卸货地省份", "unloading_province", "s", false, false, false},
		{"装货地区县", "loading_area", "s", false, false, false},
		{"卸货地区县", "unloading_area", "s", false, false, false},
		{"装货地区县经纬度", "loading_lnglog", "s", false, false, false},
		{"卸货地区县经纬度", "unloading_lnglog", "s", false, false, false},
		{"装货地省份经纬度", "loading_pro_lnglog", "s", false, false, false},
		{"卸货地省份经纬度", "unloading_pro_lnglog", "s", false, false, false},
		{"装货地经纬度", "loading_lng", "s", false, false, false},
		{"卸货地经纬度", "unloading_lng", "s", false, false, false},
		{"开票状态（1、待开票2、已开票3、发票作废4、发票冲红）", "invoicing_status", "i", false, false, false},
		{"可开票金额", "invoicing_amount", "d", false, false, false},
		{"可开票吨数", "invoicing_tonnage", "d", false, false, false},
		{"创建时间", "create_time", "t", false, false, false},
		{"修改时间", "update_time", "t", false, false, false},
		{"备注字段1", "remark1", "s", false, false, false},
		{"备注字段2", "remark2", "s", false, false, false},
		{"备注字段3", "remark3", "s", false, false, false},
		{"备注字段4", "remark4", "s", false, false, false},
		{"备注字段5", "remark5", "s", false, false, false},
	}
	billFields = []impField{
		{"id", "id", "i", false, false, false},
		{"批次号", "batch_number", "s", false, false, false},
		{"企业id", "company_id", "i", false, false, false},
		{"订单昵称", "order_name", "s", false, false, false},
		{"订单创建时间", "order_create_time", "t", false, false, false},
		{"系统编号", "system_code", "s", false, false, false},
		{"订单编号", "order_number", "s", false, false, false},
		{"运输模式", "trans_method", "i", false, false, false},
		{"承运单位", "carriage_name", "i", false, false, false},
		{"承运平台", "carriage_platform", "i", false, false, false},
		{"物流单位", "logistics_company", "s", false, false, false},
		{"运输方式", "trans_type", "i", false, false, false},
		{"合同编号", "contract_number", "s", false, false, false},
		{"销货单位", "base_sales_unit", "s", false, false, false},
		{"购货单位", "base_buyer_unit", "s", false, false, false},
		{"发货地", "base_delivery_place", "s", false, false, false},
		{"标的物", "subject_matter", "s", false, false, false},
		{"吨数", "subject_count", "s", false, false, false},
		{"装货地点", "loading_address", "s", false, false, false},
		{"运费金额", "freight_amount", "d", false, false, false},
		{"卸货地点", "unloading_point", "s", false, false, false},
		{"车数", "vehicles", "i", false, false, false},
		{"原发吨数", "original_tonnage", "d", false, false, false},
		{"实收吨数", "tonnage_paid", "d", false, false, false},
		{"附件ids", "file_ids", "s", false, false, false},
		{"订单开始时间", "order_start_time", "s", false, false, false},
		{"订单结束时间", "order_end_time", "s", false, false, false},
		{"创建时间", "create_time", "t", false, false, false},
		{"修改时间", "update_time", "t", false, false, false},
		{"备注字段1", "remark1", "s", false, false, false},
		{"备注字段2", "remark2", "s", false, false, false},
		{"备注字段3", "remark3", "s", false, false, false},
		{"备注字段4", "remark4", "s", false, false, false},
		{"备注字段5", "remark5", "s", false, false, false},
	}
	lcFields = []impField{
		{"id", "id", "i", false, false, false},
		{"企业id", "company_id", "i", false, false, false},
		{"合同id", "contract_id", "i", false, false, false},
		{"运单id", "logistics_id", "i", false, false, false},
		{"创建时间", "create_time", "t", false, false, false},
		{"修改时间", "update_time", "t", false, false, false},
		{"备注字段1", "remark1", "s", false, false, false},
		{"备注字段2", "remark2", "s", false, false, false},
		{"备注字段3", "remark3", "s", false, false, false},
		{"备注字段4", "remark4", "s", false, false, false},
		{"备注字段5", "remark5", "s", false, false, false},
	}
	twFields = []impField{
		{"id", "id", "i", false, false, false},
		{"批次号", "batch_number", "s", false, false, false},
		{"企业标识（信用代码）", "credit_code", "s", false, false, false},
		{"订单昵称", "order_name", "s", false, false, false},
		{"订单创建时间", "order_create_time", "t", false, false, false},
		{"订单编号", "order_number", "s", false, false, false},
		{"订单装车点", "order_loading_point", "s", false, false, false},
		{"订单卸货点", "order_unloading_point", "s", false, false, false},
		{"订单开始时间", "order_start_time", "s", false, false, false},
		{"订单结束时间", "order_end_time", "s", false, false, false},
		{"运单编号", "bill_number", "s", false, false, false},
		{"运单创建时间", "trans_time", "t", false, false, false},
		{"运单发货地址", "wb_loading_address", "s", false, false, false},
		{"运单卸货地址", "wb_unloading_address", "s", false, false, false},
		{"接单车辆", "car_number", "s", false, false, false},
		{"标的物", "subject_matter", "s", false, false, false},
		{"运费单价", "freight_unit", "d", false, false, false},
		{"运费金额", "freight_amount", "d", false, false, false},
		{"原发吨数", "original_tonnage", "d", false, false, false},
		{"实收吨数", "tonnage_paid", "d", false, false, false},
		{"应付金额", "pay_money", "d", false, false, false},
		{"承运平台", "carriage_platform", "i", false, false, false},
		{"发货地经纬度", "loading_pro_lnglog", "s", false, false, false},
		{"卸货地经纬度", "unloading_pro_lnglog", "s", false, false, false},
		{"开票状态", "invoicing_status", "i", false, false, false},
		{"可开票金额", "invoicing_amount", "d", false, false, false},
		{"可开票吨数", "invoicing_tonnage", "d", false, false, false},
		{"是否已清洗", "clean_flag", "b", false, false, false},
		{"创建时间", "create_time", "t", false, false, false},
		{"修改时间", "update_time", "t", false, false, false},
		{"备注字段1", "remark1", "s", false, false, false},
		{"备注字段2", "remark2", "s", false, false, false},
		{"备注字段3", "remark3", "s", false, false, false},
		{"备注字段4", "remark4", "s", false, false, false},
		{"备注字段5", "remark5", "s", false, false, false},
	}
	wlFields = []impField{
		{"id", "id", "i", false, false, false},
		{"导入批次号", "im_batch_number", "s", false, false, false},
		{"企业id", "company_id", "i", false, false, false},
		{"订单id", "order_id", "i", false, false, false},
		{"系统编号", "system_code", "s", false, false, false},
		{"票号", "bill_number", "s", false, false, false},
		{"开票单位", "issuing_office", "s", false, false, false},
		{"原煤供煤方", "raw_coal_supplier", "s", false, false, false},
		{"实收套打单位", "sstd_company", "s", false, false, false},
		{"车牌号码", "license_plate", "s", false, false, false},
		{"进厂时间", "in_time", "t", false, false, false},
		{"实收进厂时间", "ss_in_time", "t", false, false, false},
		{"合同号", "contract_no", "s", false, false, false},
		{"购买方", "purchaser", "s", false, false, false},
		{"煤种", "coal_type", "s", false, false, false},
		{"皮重", "tare", "d", false, false, false},
		{"实收皮重", "ss_tare", "d", false, false, false},
		{"皮重时间", "tare_time", "t", false, false, false},
		{"毛重", "rough_weight", "d", false, false, false},
		{"实收毛重", "ss_rough_weight", "d", false, false, false},
		{"毛重时间", "rough_weight_time", "t", false, false, false},
		{"装煤司机", "coal_loading_driver", "s", false, false, false},
		{"装煤时间", "coal_loading_time", "t", false, false, false},
		{"净重", "net_weight", "d", false, false, false},
		{"实收净重", "ss_net_weight", "d", false, false, false},
		{"开票时间", "Invoicing_time", "t", false, false, false},
		{"实收套打时间", "sstd_time", "t", false, false, false},
		{"上传时间", "upload_time", "t", false, false, false},
		{"实收状态", "ss_status", "s", false, false, false},
		{"打印状态", "print_status", "s", false, false, false},
		{"时间间隔（分钟）", "time_interval_minute", "i", false, false, false},
		{"是否补录", "add_or_not", "s", false, false, false},
		{"回皮人", "repatriator", "s", false, false, false},
		{"称重人", "weigher", "s", false, false, false},
		{"实收回皮人", "ss_repatriator", "s", false, false, false},
		{"实收称重人", "ss_weigher", "s", false, false, false},
		{"作废原因", "reason_cancellation", "s", false, false, false},
		{"补录原因", "reason_supplementary", "s", false, false, false},
		{"创建时间", "create_time", "t", false, false, false},
		{"修改时间", "update_time", "t", false, false, false},
		{"备注字段1", "remark1", "s", false, false, false},
		{"备注字段2", "remark2", "s", false, false, false},
		{"备注字段3", "remark3", "s", false, false, false},
		{"备注字段4", "remark4", "s", false, false, false},
		{"备注字段5", "remark5", "s", false, false, false},
	}
)

// ── 导出端点(5) ──

// writeXlsx 最小 xlsx(inline string 单元格;标题行+数据行,文件名导出数据.xlsx)
func writeXlsx(w http.ResponseWriter, grid [][]string) {
	var sheet strings.Builder
	sheet.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>`)
	for ri, row := range grid {
		sheet.WriteString(`<row r="` + strconv.Itoa(ri+1) + `">`)
		for ci, v := range row {
			ref, _ := colRef(ci)
			sheet.WriteString(`<c r="` + ref + strconv.Itoa(ri+1) + `" t="inlineStr"><is><t xml:space="preserve">` + xmlEscape(v) + `</t></is></c>`)
		}
		sheet.WriteString(`</row>`)
	}
	sheet.WriteString(`</sheetData></worksheet>`)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment;filename="+url.QueryEscape("导出数据.xlsx"))
	zw := zip.NewWriter(w)
	mk := func(name, body string) {
		f, _ := zw.Create(name)
		_, _ = f.Write([]byte(body))
	}
	mk("[Content_Types].xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`+
		`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">`+
		`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>`+
		`<Default Extension="xml" ContentType="application/xml"/>`+
		`<Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>`+
		`<Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>`+
		`</Types>`)
	mk("_rels/.rels", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`+
		`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`+
		`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>`+
		`</Relationships>`)
	mk("xl/workbook.xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`+
		`<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">`+
		`<sheets><sheet name="Sheet1" sheetId="1" r:id="rId1"/></sheets></workbook>`)
	mk("xl/_rels/workbook.xml.rels", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`+
		`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`+
		`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/>`+
		`</Relationships>`)
	mk("xl/worksheets/sheet1.xml", sheet.String())
	_ = zw.Close()
}

// colIndex "D5" → 3(0 基列号)
func colIndex(ref string) (int, error) {
	n := 0
	for i := 0; i < len(ref); i++ {
		c := ref[i]
		if c < 'A' || c > 'Z' {
			return n - 1, nil
		}
		n = n*26 + int(c-'A') + 1
	}
	return n - 1, nil
}

func colRef(i int) (string, error) {
	s := ""
	for i >= 0 {
		s = string(rune('A'+i%26)) + s
		i = i/26 - 1
	}
	return s, nil
}

func xmlEscape(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&#34;", "'", "&#39;")
	return r.Replace(s)
}

func (s *Service) loExportExcel(w http.ResponseWriter, r *http.Request) {
	s.exportHandler(w, r, "chain_logistics", loFields)
}

func (s *Service) billExportExcel(w http.ResponseWriter, r *http.Request) {
	s.exportHandler(w, r, "chain_logistics_bill", billFields)
}

func (s *Service) lcExportExcel(w http.ResponseWriter, r *http.Request) {
	s.exportHandler(w, r, "chain_logistics_contract", lcFields)
}

func (s *Service) twExportExcel(w http.ResponseWriter, r *http.Request) {
	s.exportHandler(w, r, "chain_logistics_third_waybill", twFields)
}

func (s *Service) wlExportExcel(w http.ResponseWriter, r *http.Request) {
	s.exportHandler(w, r, "chain_logistics_weighing_list", wlFields)
}

// exportHandler 通用导出:SELECT * 全表 → grid → xlsx(标题行+字符串单元格)
func (s *Service) exportHandler(w http.ResponseWriter, r *http.Request, table string, fields []impField) {
	rows, err := s.pick(r).Query("SELECT * FROM " + table)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, err := platform.ScanMaps(rows)
	if err != nil {
		platform.Unknown(w)
		return
	}
	grid := make([][]string, 0, len(ms)+1)
	head := make([]string, 0, len(fields))
	for _, f := range fields {
		head = append(head, f.header)
	}
	grid = append(grid, head)
	for _, m := range ms {
		row := make([]string, 0, len(fields))
		for _, f := range fields {
			row = append(row, exportCell(m[f.col], f.kind))
		}
		grid = append(grid, row)
	}
	writeXlsx(w, grid)
}

// exportCell 对照 writeExcel 的 value.toString():LocalDateTime.toString 秒为 0 时省略
// ("2023-03-02T17:30")、Boolean → true/false、DECIMAL 保小数位、null 空
func exportCell(v any, kind string) string {
	if v == nil {
		return ""
	}
	s := btsStr(v)
	if kind == "b" {
		return map[string]string{"1": "true", "0": "false"}[s]
	}
	if kind == "t" && len(s) == 19 && s[10] == ' ' {
		s = s[:10] + "T" + s[11:]
		if s[17:] == "00" {
			return s[:16] // 秒为 0(纳秒恒 0)→ LocalDateTime.toString 省略 :ss
		}
		return s
	}
	return s
}

// ── 导入端点(5) ──

// formFile 取 multipart 文件(Java @RequestParam("file") 缺失 → 999)
func formFile(r *http.Request) (*multipart.FileHeader, bool) {
	_, fh, err := r.FormFile("file")
	if err != nil {
		return nil, false
	}
	return fh, true
}

func emptyDict() map[string]string { return map[string]string{} }

// getDictMap 对照 CmnServiceImpl.getDictMap:chain_common_dict 全量 d_key→d_value。
// 全 NULL 行 MyBatis 返回 null 元素 → m.get NPE → 500(与 dictMap 同,1:1 复刻)。
func (s *Service) getDictMap(r *http.Request) (map[string]string, error) {
	rows, err := s.pick(r).Query("SELECT d_key, d_value FROM chain_common_dict")
	if err != nil {
		return nil, err
	}
	ms, err := platform.ScanMaps(rows)
	if err != nil {
		return nil, err
	}
	out := map[string]string{}
	for _, m := range ms {
		if m["d_key"] == nil && m["d_value"] == nil {
			return nil, errors.New("java getDictMap NPE on all-null dict row")
		}
		if m["d_key"] != nil && m["d_value"] != nil {
			out[btsStr(m["d_key"])] = btsStr(m["d_value"])
		}
	}
	return out, nil
}

// loImportExcel 5 参 readExcel(单 sheet+字典)+saveIrt:错误包装"运单sheet中，"→
// errs=[导入失败 0];成功→imBatchNumber=UUID、systemCode 顺序、saveBatch(fill)、
// sucInfos、日志(op5),响应=完整 ImportResult
func (s *Service) loImportExcel(w http.ResponseWriter, r *http.Request) {
	s.importWithIrt(w, r, "chain_logistics", loFields, "WLYD.", "运单sheet中，", "运单管理数据成功",
		"导入运单", 5, true)
}

// wlImportExcel 同款:"过磅单sheet中，"+WLGBD.+"过磅单管理数据成功"+op9
func (s *Service) wlImportExcel(w http.ResponseWriter, r *http.Request) {
	s.importWithIrt(w, r, "chain_logistics_weighing_list", wlFields, "WLGBD.", "过磅单sheet中，",
		"过磅单管理数据成功", "导入过磅单", 9, true)
}

// importWithIrt 对照 saveIrt(ImportResult, request) 全流程;useDict=5 参 readExcel(单 sheet)
func (s *Service) importWithIrt(w http.ResponseWriter, r *http.Request, table string, fields []impField,
	idt, wrap, sucLabel, logLabel string, opType int, useDict bool) {
	fh, ok := formFile(r)
	if !ok {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	u := currentUser(r)
	prefix := idt + u.EnterpriseCode + time.Now().Format("20060102")
	syscode, ok2 := nextSystemCode(d, table, prefix, "0") // service.importExcel 起始 0
	if !ok2 {
		platform.Unknown(w)
		return
	}
	dict := emptyDict()
	if useDict {
		var err error
		dict, err = s.getDictMap(r)
		if err != nil {
			platform.Unknown(w)
			return
		}
	}
	irt := readExcelFile(fh, fields, syscode, false, dict)
	if len(irt.ErrInfos) > 0 {
		wrapped := make([]string, 0, len(irt.ErrInfos))
		for _, s2 := range irt.ErrInfos {
			wrapped = append(wrapped, wrap+s2)
		}
		irt.ErrInfos = wrapped
		irt.Errs = []any{"导入失败 0"}
		irt.SucInfos = []any{}
		// 错误路径同样按实体全键整形(Java 序列化的是内存实体,不落库故无 id/fill)
		for i2, row := range irt.DataList {
			irt.DataList[i2] = platform.EntityEcho(platform.Camelize(row), entOfTable(table).echoKeys(), nil)
		}
		platform.Success(w, irt)
		return
	}
	irt.Errs = []any{}
	batch := uuidV4()
	curcode, ok3 := nextSystemCode(d, table, prefix, "1") // saveIrt getSystemCode 起始 1
	if !ok3 {
		platform.Unknown(w)
		return
	}
	for _, row := range irt.DataList {
		row["im_batch_number"] = batch
		row["system_code"] = curcode
		next, okb := bumpSystemCode(curcode)
		if !okb {
			platform.Unknown(w)
			return
		}
		curcode = next
	}
	for i2, row := range irt.DataList {
		vals := insertVals(row)
		delete(vals, "id")
		fill := table == "chain_logistics" || table == "chain_logistics_weighing_list"
		var newID int64
		var err error
		if fill {
			newID, err = platform.Insert(d, table, vals)
		} else {
			newID, err = platform.InsertNoFill(d, table, vals)
		}
		if err != nil {
			platform.Unknown(w)
			return
		}
		// 回显行:解析键驼峰化后按实体全键整形(Java 内存实体 Jackson 全键序列化;
		// fill 判定须用解析行——EntityEcho 会预置全键 null)
		echo := platform.EntityEcho(platform.Camelize(row), entOfTable(table).echoKeys(), map[string]any{
			"id": json.Number(strconv.FormatInt(newID, 10)),
		})
		if _, has := row["create_time"]; !has && fill {
			echo["createTime"] = nowSpace()
		}
		if _, has := row["update_time"]; !has && fill {
			echo["updateTime"] = nowTNanos()
		}
		irt.DataList[i2] = echo
	}
	irt.SucInfos = []any{sucLabel + " " + strconv.Itoa(len(irt.DataList))}
	logOp(d, u, "导入成功，批次号为："+batch+"；"+logLabel+strconv.Itoa(len(irt.DataList))+"条；", opType)
	platform.Success(w, irt)
}

// billImportExcel 4 参 readExcel(全部 sheet,无字典):有错 → fail(errInfos);成功 → saveBatch+文本
func (s *Service) billImportExcel(w http.ResponseWriter, r *http.Request) {
	s.importPlain(w, r, "chain_logistics_bill", billFields, true)
}

func (s *Service) lcImportExcel(w http.ResponseWriter, r *http.Request) {
	// service.importExcel 先查 selectSystemCode,而 chain_logistics_contract 无 system_code
	// 列 → SQL 报错 → 恒 999(8765 实测),照实复刻
	fh, ok := formFile(r)
	if !ok {
		platform.Unknown(w)
		return
	}
	if _, ok2 := nextSystemCode(s.pick(r), "chain_logistics_contract", "WLDD."+currentUser(r).EnterpriseCode+time.Now().Format("20060102"), "0"); !ok2 {
		platform.Unknown(w)
		return
	}
	irt := readExcelFile(fh, lcFields, "", true, emptyDict())
	if len(irt.ErrInfos) > 0 {
		platform.FailData(w, 999, "失败", irt.ErrInfos)
		return
	}
	d := s.pick(r)
	for _, row := range irt.DataList {
		vals := insertVals(row)
		delete(vals, "id")
		if _, err := platform.InsertNoFill(d, "chain_logistics_contract", vals); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, "导入成功，共"+strconv.Itoa(len(irt.DataList))+"条")
}

func (s *Service) twImportExcel(w http.ResponseWriter, r *http.Request) {
	// 同 lc:chain_logistics_third_waybill 无 system_code 列 → 恒 999(8765 实测)
	fh, ok := formFile(r)
	if !ok {
		platform.Unknown(w)
		return
	}
	if _, ok2 := nextSystemCode(s.pick(r), "chain_logistics_third_waybill", "WLDD."+currentUser(r).EnterpriseCode+time.Now().Format("20060102"), "0"); !ok2 {
		platform.Unknown(w)
		return
	}
	irt := readExcelFile(fh, twFields, "", true, emptyDict())
	if len(irt.ErrInfos) > 0 {
		platform.FailData(w, 999, "失败", irt.ErrInfos)
		return
	}
	d := s.pick(r)
	for _, row := range irt.DataList {
		vals := insertVals(row)
		delete(vals, "id")
		if _, err := platform.InsertNoFill(d, "chain_logistics_third_waybill", vals); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, "导入成功，共"+strconv.Itoa(len(irt.DataList))+"条")
}

// importPlain 对照 Bill/Contract/ThirdWaybill importExcel:错误→999(data=列表);
// 成功→saveBatch(fill 单独处理)→"导入成功，共N条"
func (s *Service) importPlain(w http.ResponseWriter, r *http.Request, table string, fields []impField, billFill bool) {
	fh, ok := formFile(r)
	if !ok {
		platform.Unknown(w)
		return
	}
	irt := readExcelFile(fh, fields, "", true, emptyDict())
	if len(irt.ErrInfos) > 0 {
		platform.FailData(w, 999, "失败", irt.ErrInfos)
		return
	}
	d := s.pick(r)
	for _, row := range irt.DataList {
		vals := insertVals(row)
		delete(vals, "id")
		if billFill {
			if _, has := vals["update_time"]; !has {
				vals["update_time"] = nowSpace()
			}
			if _, err := platform.InsertNoFill(d, table, vals); err != nil {
				platform.Unknown(w)
				return
			}
		} else if _, err := platform.InsertNoFill(d, table, vals); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, "导入成功，共"+strconv.Itoa(len(irt.DataList))+"条")
}

// insertVals 解析行 → 入库列(时间 T 格式转空格;Boolean 转 1/0;nil 剔除交给库默认)
func insertVals(row map[string]any) map[string]any {
	out := make(map[string]any, len(row))
	for k, v := range row {
		if v == nil {
			continue
		}
		switch t := v.(type) {
		case string:
			if len(t) == 19 && t[10] == 'T' {
				out[k] = t[:10] + " " + t[11:]
			} else {
				out[k] = t
			}
		case bool:
			if t {
				out[k] = 1
			} else {
				out[k] = 0
			}
		default:
			out[k] = v
		}
	}
	return out
}

func uuidV4() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// ── readExcel 主体 ──

type importResult struct {
	DataList []map[string]any `json:"dataList"`
	SucInfos []any            `json:"sucInfos"`
	Errs     []any            `json:"errs"`
	ErrInfos []string         `json:"errInfos"`
}

// readExcelFile 对照 readExcel。allSheets=false 为 5 参版本(只取 sheet0+编号注入)。
func readExcelFile(fh *multipart.FileHeader, fields []impField, systemCode string, allSheets bool, dictMap map[string]string) *importResult {
	spec := newSpec(fields)
	irt := &importResult{DataList: []map[string]any{}, ErrInfos: []string{}}
	if fh == nil {
		irt.ErrInfos = []string{"上传文件为空"}
		return irt
	}
	if !strings.HasSuffix(strings.ToLower(fh.Filename), ".xlsx") {
		irt.ErrInfos = []string{"上传文件格式不正确"}
		return irt
	}
	f, err := fh.Open()
	if err != nil {
		irt.ErrInfos = []string{"解析 Excel 异常，请检查文件格式和内容"}
		return irt
	}
	defer f.Close()
	xl, perr := parseXlsx(f)
	if perr != nil {
		irt.ErrInfos = []string{"解析 Excel 异常，请检查文件格式和内容"}
		return irt
	}
	// notNull 待查集(声明序;Java HashMap 迭代序的确定性近似)
	var notNullLeft []impField
	for _, fl := range fields {
		if fl.notNull {
			notNullLeft = append(notNullLeft, fl)
		}
	}
	sheets := xl.sheets
	if !allSheets {
		if len(sheets) == 0 {
			return irt
		}
		sheets = sheets[:1]
	}
	rowIdx := 0 // SystemCode 注入序号(Java i-1,数据行从 0 起)
	for _, sh := range sheets {
		first := true
		mapped := map[int]impField{}
		for _, row := range sh.rows {
			if first {
				first = false
				for ci, cell := range row.cells {
					header := strings.ReplaceAll(strings.TrimSpace(cellText(cell, xl)), "\n", "")
					if fl, ok := spec.matchHeader(header); ok {
						mapped[ci] = fl
					}
				}
				continue
			}
			if len(row.cells) == 0 || isEmptyRow(row.cells) {
				continue
			}
			ent := map[string]any{}
			allBlank := true
			for ci, cell := range row.cells {
				fl, ok := mapped[ci]
				if !ok {
					continue
				}
				val := cellText(cell, xl)
				if strings.TrimSpace(val) != "" {
					allBlank = false
				}
				// validateField(NotBank 必填 / Dict 字典)
				if fl.notNull && strings.TrimSpace(val) == "" {
					irt.ErrInfos = append(irt.ErrInfos, sh.name+"sheet中，第"+strconv.Itoa(row.n)+"行第"+strconv.Itoa(ci+1)+"列字段为空")
				}
				if fl.dict {
					if _, err := strconv.Atoi(strings.TrimSpace(val)); err != nil {
						if _, has := dictMap[val]; !has {
							irt.ErrInfos = append(irt.ErrInfos, sh.name+"sheet中，第"+strconv.Itoa(row.n)+"行第"+strconv.Itoa(ci+1)+"列无相关字典")
						}
					}
				}
				// handleField(行级异常 Java 吞掉:字段留空,行保留)
				if strings.TrimSpace(val) != "" {
					setField(ent, fl, val)
				}
				// notNullMap.remove(无论单元格是否为空,列被访问即移除)
				for i2, nf := range notNullLeft {
					if nf.col == fl.col {
						notNullLeft = append(notNullLeft[:i2], notNullLeft[i2+1:]...)
						break
					}
				}
			}
			// SystemCodeColumn 注入(4 参版本 systemCode 为空 → 不注入)
			if systemCode != "" && len(systemCode) >= 20 {
				for _, fl := range fields {
					if fl.sys {
						if n, err := strconv.Atoi(systemCode[20:]); err == nil {
							ent[fl.col] = systemCode[:20] + autoGenericCode(strconv.Itoa(n+rowIdx), 7)
						}
					}
				}
			}
			rowIdx++
			if !allBlank {
				irt.DataList = append(irt.DataList, ent)
			}
			// notNullMap 耗尽(Java 原版 bug:仅首行可触发"不能为空");
			// 输出序 = JAR commons-collections4 HashedMap.keySet 迭代序(见 platform.JavaHashedMapOrder)。
			if len(notNullLeft) > 0 {
				keep := map[string]bool{}
				for _, nf := range notNullLeft {
					keep[platform.JavaFieldName(nf.col)] = true
				}
				var all []string
				title := map[string]string{}
				for _, fl := range fields {
					if fl.notNull {
						k := platform.JavaFieldName(fl.col)
						all = append(all, k)
						title[k] = fl.header
					}
				}
				for _, k := range platform.JavaHashedMapOrder(all, keep) {
					irt.ErrInfos = append(irt.ErrInfos, sh.name+"sheet中，"+title[k]+"不能为空")
				}
				notNullLeft = nil
			}
		}
	}
	return irt
}

// javaHash 对照 Java String.hashCode
func javaHash(s string) uint32 {
	var h uint32
	for _, c := range s { // Java 按 UTF-16 code unit,与 rune( BMP)一致
		h = h*31 + uint32(c)
	}
	return h
}

// javaMapOrder 模拟 JDK8 HashMap<String,?> 的迭代序(插入序建桶、按阈值翻倍扩容、
// resize 保相对序)——ExcelUtils.classMap 的 contrainKey 匹配次序由此决定
func javaMapOrder(keys []string) []string {
	// 精确模拟 JDK11 HashMap.computeIfAbsent 迭代序(ExcelUtils.buildClassMap 逐标题
	// computeIfAbsent 建表):入口 size>threshold 先扩容;冲突时头插(tab[i]=newNode(...,first),
	// 与 putVal 尾插不同);扩容按 spread hash & oldCap 拆 lo/hi 保相对序。8765 实测一致
	capacity, threshold, size := 16, 12, 0
	type entry struct {
		key string
		sp  uint32
	}
	table := make([][]entry, capacity)
	for _, k := range keys {
		h := javaHash(k)
		sp := h ^ (h >> 16)
		if size > threshold {
			old := capacity
			capacity *= 2
			threshold = capacity * 3 / 4
			nt := make([][]entry, capacity)
			for j := 0; j < old; j++ {
				var lo, hi []entry
				for _, e := range table[j] {
					if e.sp&uint32(old) == 0 {
						lo = append(lo, e)
					} else {
						hi = append(hi, e)
					}
				}
				nt[j] = lo
				if len(hi) > 0 {
					nt[j+old] = hi
				}
			}
			table = nt
		}
		idx := sp & uint32(capacity-1)
		table[idx] = append([]entry{{key: k, sp: sp}}, table[idx]...) // 头插
		size++
	}
	out := make([]string, 0, len(keys))
	for _, b := range table {
		for _, e := range b {
			out = append(out, e.key)
		}
	}
	return out
}

// excelSpec 实体导入规格:字段表 + Java classMap 迭代序
type excelSpec struct {
	fields []impField
	order  []string // 迭代序的标题
	byHdr  map[string]impField
}

func newSpec(fields []impField) *excelSpec {
	hdrs := make([]string, 0, len(fields))
	seen := map[string]bool{}
	by := map[string]impField{}
	for _, f := range fields {
		if !seen[f.header] {
			hdrs = append(hdrs, f.header)
			seen[f.header] = true
		}
		by[f.header] = f
	}
	return &excelSpec{fields: fields, order: javaMapOrder(hdrs), byHdr: by}
}

// matchHeader contrainKey:按 Java classMap 迭代序取首个包含命中的标题
func (sp *excelSpec) matchHeader(cell string) (impField, bool) {
	if cell == "" {
		return impField{}, false
	}
	for _, h := range sp.order {
		if strings.Contains(cell, h) {
			return sp.byHdr[h], true
		}
	}
	return impField{}, false
}

func isEmptyRow(cells []xlsxCell) bool {
	for _, c := range cells {
		if c.kind != "blank" && strings.TrimSpace(c.text) != "" {
			return false
		}
	}
	return true
}

// setField 对照 handleField 类型转换;行级异常吞掉(返回不中断)
func setField(ent map[string]any, fl impField, val string) {
	switch fl.kind {
	case "i":
		n, err := strconv.Atoi(strings.TrimSpace(val))
		if err != nil {
			n = 0 // NumberUtils.toInt 默认 0(含字典翻译失效落 0 的原版行为)
		}
		ent[fl.col] = json.Number(strconv.Itoa(n))
	case "d":
		if _, err := strconv.ParseFloat(strings.TrimSpace(val), 64); err != nil {
			return // new BigDecimal 异常被 forEach 内 catch 吞掉,字段留空
		}
		ent[fl.col] = json.Number(strings.TrimSpace(val))
	case "t":
		if t, ok := formatTimeAny(val); ok {
			ent[fl.col] = t // LocalDateTime 内存态 ISO T;失败 → 字段 null 不报错
		}
	case "b":
		ent[fl.col] = boolToBool(val)
	default:
		ent[fl.col] = val
	}
}

// formatTimeAny 对照 formatTime 六种格式;失败 ok=false(字段 null)
func formatTimeAny(s string) (string, bool) {
	layouts := []string{"2006-01-02 15:04:05", "2006-1-2 15:04:05", "2006-1-2 15:4:05",
		"2006/01/02 15:04:05", "2006/1/2 15:04:05", "2006/1/2 15:4:05"}
	for _, l := range layouts {
		if t, err := time.ParseInLocation(l, s, time.Local); err == nil {
			return t.Format("2006-01-02T15:04:05"), true
		}
	}
	return "", false
}

func boolToBool(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "true", "yes", "on", "y":
		return true
	}
	return false
}

// ── xlsx 解析 ──

type xlsxCell struct {
	kind    string // s/inline/str/num/bool/blank
	text    string // 字符串原始文本或数值文本
	str     bool   // 单元格本身是字符串类型(t=s/inlineStr/str)
	percent bool   // 样式为百分比格式
	dateFmt bool   // 样式为日期格式
}

type xlsxRow struct {
	n     int // 1-based 实际行号(错误消息用)
	cells []xlsxCell
}

type xlsxSheet struct {
	name string
	rows []xlsxRow
}

type xlsx struct {
	sheets   []xlsxSheet
	shared   []string
	percents map[int]bool // style idx → 百分比/日期格式
	dates    map[int]bool
}

type sheetRef struct {
	name string
	rid  string
}

func parseXlsx(f multipart.File) (*xlsx, error) {
	size, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	zr, err := zip.NewReader(f.(io.ReaderAt), size)
	if err != nil {
		return nil, err
	}
	out := &xlsx{percents: map[int]bool{}, dates: map[int]bool{}}
	rels := map[string]string{}
	var refs []sheetRef
	var sheetXMLs []struct{ name, data string }
	for _, zf := range zr.File {
		switch {
		case zf.Name == "xl/sharedStrings.xml":
			data, err := readZip(zf)
			if err == nil {
				out.shared = parseSharedStrings(string(data))
			}
		case zf.Name == "xl/_rels/workbook.xml.rels":
			data, err := readZip(zf)
			if err == nil {
				rels = parseRels(string(data))
			}
		case zf.Name == "xl/workbook.xml":
			data, err := readZip(zf)
			if err == nil {
				refs = parseWorkbookSheets(string(data))
			}
		case zf.Name == "xl/styles.xml":
			data, err := readZip(zf)
			if err == nil {
				parseStyles(string(data), out)
			}
		default:
			if strings.HasPrefix(zf.Name, "xl/worksheets/sheet") {
				data, err := readZip(zf)
				if err == nil {
					sheetXMLs = append(sheetXMLs, struct{ name, data string }{zf.Name, string(data)})
				}
			}
		}
	}
	byPath := map[string]string{}
	for _, sx := range sheetXMLs {
		byPath[sx.name] = sx.data
	}
	for _, ref := range refs {
		target := strings.TrimPrefix(rels[ref.rid], "/")
		if !strings.HasPrefix(target, "xl/") {
			target = "xl/" + target
		}
		data, ok := byPath[target]
		if !ok {
			continue
		}
		out.sheets = append(out.sheets, xlsxSheet{name: ref.name, rows: parseSheetRows(data, out)})
	}
	return out, nil
}

func readZip(zf *zip.File) ([]byte, error) {
	rc, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

func parseSharedStrings(data string) []string {
	var out []string
	dec := xml.NewDecoder(strings.NewReader(data))
	inT := false
	cur := strings.Builder{}
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local == "t" {
				inT = true
				cur.Reset()
			}
		case xml.CharData:
			if inT {
				cur.Write(t)
			}
		case xml.EndElement:
			if t.Name.Local == "t" {
				out = append(out, cur.String())
				inT = false
			}
		}
	}
	return out
}

func parseRels(data string) map[string]string {
	out := map[string]string{}
	dec := xml.NewDecoder(strings.NewReader(data))
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		if se, ok := tok.(xml.StartElement); ok && se.Name.Local == "Relationship" {
			var id, target string
			for _, a := range se.Attr {
				if a.Name.Local == "Id" {
					id = a.Value
				}
				if a.Name.Local == "Target" {
					target = a.Value
				}
			}
			out[id] = target
		}
	}
	return out
}

func parseWorkbookSheets(data string) []sheetRef {
	var out []sheetRef
	dec := xml.NewDecoder(strings.NewReader(data))
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		if se, ok := tok.(xml.StartElement); ok && se.Name.Local == "sheet" {
			var name, rid string
			for _, a := range se.Attr {
				if a.Name.Local == "name" {
					name = a.Value
				}
				if (a.Name.Space != "" && strings.HasSuffix(a.Name.Space, "relationships") && a.Name.Local == "id") || a.Name.Local == "id" {
					rid = a.Value
				}
			}
			out = append(out, sheetRef{name: name, rid: rid})
		}
	}
	return out
}

// parseStyles cellXfs 的 numFmtId → 每样式是否百分比/日期(内置+自定义粗判)
func parseStyles(data string, out *xlsx) {
	dec := xml.NewDecoder(strings.NewReader(data))
	section := ""
	custom := map[string]string{} // numFmtId → code
	var numFmtIDs []string
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "numFmts", "cellXfs":
				section = t.Name.Local
			case "numFmt":
				var id, code string
				for _, a := range t.Attr {
					if a.Name.Local == "numFmtId" {
						id = a.Value
					}
					if a.Name.Local == "formatCode" {
						code = a.Value
					}
				}
				custom[id] = code
			case "xf":
				if section == "cellXfs" {
					for _, a := range t.Attr {
						if a.Name.Local == "numFmtId" {
							numFmtIDs = append(numFmtIDs, a.Value)
						}
					}
				}
			}
		case xml.EndElement:
			if t.Name.Local == "numFmts" || t.Name.Local == "cellXfs" {
				section = ""
			}
		}
	}
	for idx, id := range numFmtIDs {
		code := custom[id]
		builtin, _ := strconv.Atoi(id)
		if builtin == 9 || builtin == 10 || strings.Contains(code, "%") {
			out.percents[idx] = true
		}
		isDate := (builtin >= 14 && builtin <= 22) || (builtin >= 27 && builtin <= 36) ||
			(builtin >= 45 && builtin <= 47) || (builtin >= 50 && builtin <= 58)
		low := strings.ToLower(code)
		if isDate || (code != "" && !strings.Contains(code, "%") &&
			(strings.Contains(low, "y") || strings.Contains(low, "d") || strings.Contains(low, "h"))) {
			out.dates[idx] = true
		}
	}
}

// parseSheetRows 行/单元格展平(列序按引用字母;n=行号)
func parseSheetRows(data string, xl *xlsx) []xlsxRow {
	var rows []xlsxRow
	dec := xml.NewDecoder(strings.NewReader(data))
	inRow, inCell, inV, inT := false, false, false, false
	var row xlsxRow
	var cell xlsxCell
	var vBuf, tBuf strings.Builder
	flushCell := func() {
		if !inCell {
			return
		}
		switch cell.kind {
		case "s": // shared string
			idx, _ := strconv.Atoi(strings.TrimSpace(vBuf.String()))
			if idx >= 0 && idx < len(xl.shared) {
				cell.text = xl.shared[idx]
			}
			cell.str = true
		case "inline":
			cell.text = tBuf.String()
			cell.str = true
		case "str":
			cell.text = vBuf.String()
			cell.str = true
		case "num":
			cell.text = vBuf.String() // 数值格:<v> 文本即值
		case "bool":
			cell.kind = "blank" // Java 布尔单元格 getStringCellValue 抛异常 → 行级跳过,等价空
		}
		row.cells = append(row.cells, cell)
		inCell, inV, inT = false, false, false
		vBuf.Reset()
		tBuf.Reset()
	}
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "row":
				inRow = true
				row = xlsxRow{n: len(rows) + 1}
				for _, a := range t.Attr {
					if a.Name.Local == "r" {
						if n, err := strconv.Atoi(a.Value); err == nil {
							row.n = n
						}
					}
				}
			case "c":
				if inRow {
					flushCell()
					inCell = true
					cell = xlsxCell{kind: "num"}
					var style, ref string
					for _, a := range t.Attr {
						switch a.Name.Local {
						case "t":
							switch a.Value {
							case "s":
								cell.kind = "s"
							case "inlineStr":
								cell.kind = "inline"
							case "str":
								cell.kind = "str"
							case "b":
								cell.kind = "bool"
							case "n":
								cell.kind = "num"
							}
						case "s":
							style = a.Value
						case "r":
							ref = a.Value
						}
					}
					// openpyxl 等写入方会省略空单元格 → 必须按 r 属性回填真实列位
					if ref != "" {
						if ci, err := colIndex(ref); err == nil {
							for len(row.cells) < ci {
								row.cells = append(row.cells, xlsxCell{kind: "blank"})
							}
						}
					}
					if style != "" {
						if idx, err := strconv.Atoi(style); err == nil {
							cell.percent, cell.dateFmt = xl.percents[idx], xl.dates[idx]
						}
					}
				}
			case "v":
				if inCell {
					inV = true
					vBuf.Reset()
				}
			case "t":
				if inCell {
					inT = true
					tBuf.Reset()
				}
			}
		case xml.CharData:
			if inV {
				vBuf.Write(t)
			}
			if inT {
				tBuf.Write(t)
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "c":
				flushCell()
			case "row":
				flushCell()
				if inRow {
					rows = append(rows, row)
					inRow = false
				}
			case "v":
				inV = false
			case "t":
				inT = false
			}
		}
	}
	return rows
}

// cellText 对照 getCellValue:字符串单元格直接返回(trim);数值单元格
// 百分比→*100%、日期→yyyy-MM-dd、其余整数化(小数保留,g 格式)
func cellText(c xlsxCell, xl *xlsx) string {
	if c.kind == "blank" {
		return ""
	}
	if c.str {
		return strings.TrimSpace(c.text)
	}
	db, err := strconv.ParseFloat(strings.TrimSpace(c.text), 64)
	if err != nil {
		return strings.TrimSpace(c.text)
	}
	if c.percent {
		return strconv.FormatFloat(db*100, 'f', -1, 64) + "%"
	}
	if c.dateFmt {
		return excelDateToDay(db)
	}
	aa := int64(db)
	if db >= 0 {
		if db > float64(aa) {
			return strconv.FormatFloat(db, 'g', -1, 64)
		}
		return strconv.FormatInt(aa, 10)
	}
	if db < float64(aa) {
		return strconv.FormatFloat(db, 'g', -1, 64)
	}
	return strconv.FormatInt(aa, 10)
}

// excelDateToDay Excel 序列日期(1900 系统) → yyyy-MM-dd(POI getCellValue 日期兜底格式)
func excelDateToDay(db float64) string {
	sec := int64((db - 25569) * 86400)
	return time.Unix(sec, 0).UTC().Format("2006-01-02")
}

// entOfTable 表名 → 实体元数据(导入回显用)
func entOfTable(table string) entShape {
	switch table {
	case "chain_logistics":
		return entLo
	case "chain_logistics_weighing_list":
		return entWl
	}
	return entShape{}
}
