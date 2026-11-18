// warehouse - excel.go Excel 导入导出(自四期 logistics/excel.go 复制适配,仓内惯例:
// 每服务包一份 excel 机器;六期新增 LocalDate "D" 与 Long "l" 两种列类型)。
// 导出对照 ExcelUtils.writeExcel:全部 @ExcelColumn 字段按声明序,标题行+数据行,
// 单元格一律字符串(LocalDateTime.toString() ISO T、BigDecimal 保小数位、null 空)。
// 导入对照 ExcelUtils.readExcel:表头包含匹配、getCellValue 数字/日期/百分比转换、
// NotBankColumn 必填(notNullMap 首行后耗尽的原版 bug)、DictColumn 校验(翻译实际不生效,
// 非数字值落 0)、SystemCodeColumn 逐行注入、4 参版本遍历全部 sheet、5 参版本只取 sheet0。
// 行级异常 Java 在 forEach 内吞掉(字段留空行保留),非整行丢弃——照实复刻。
// 关键怪癖:hyqc-core 4 参重载 readExcel(cls,file,systemCode,sheetIdx) 把 systemCode
// 丢成 "" 再委托 5 参——warehouse 模板导入实际从不注入编号,selectSystemCode 只是
// 先跑一次(表无 system_code 列时 SQL 报错 → 恒 999,bus_* 表全部无该列,已实测)。
// Binary XLS is normalized by LibreOffice before the shared OOXML reader.

package commonrest

import (
	"archive/zip"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// btsStr 数值/字节扫描值转字符串
func btsStr(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case []byte:
		return string(t)
	case string:
		return t
	default:
		return fmt.Sprintf("%v", t)
	}
}

// autoGenericCode 对照 StringHelper.autoGenericCode:parseInt 后 %0Nd 补零
func autoGenericCode(code string, width int) string {
	n, err := strconv.Atoi(strings.TrimSpace(code))
	if err != nil {
		return code // Java NumberFormatException 由调用方先拦,此处防御
	}
	s := strconv.Itoa(n)
	for len(s) < width {
		s = "0" + s
	}
	return s
}

// impField Excel 列映射(header 与 Java @ExcelColumn value 逐字一致)
type impField struct {
	header  string
	col     string // snake 列
	kind    string // s=String i=Integer l=Long d=BigDecimal t=LocalDateTime b=Boolean D=LocalDate
	sys     bool   // SystemCodeColumn
	dict    bool   // DictColumn
	notNull bool   // NotBankColumn
}

// ── 导出 ──

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
			return n, nil
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

func (s *Service) exportHandler(w http.ResponseWriter, r *http.Request, table string, fields []impField) {
	rows, err := s.pick(r).Query("SELECT " + strings.Join(entityForTable(table).dbcols, ",") + " FROM " + table)
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
	// l/D:Long/LocalDate toString 与扫描文本一致,无需转换
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

// getDictMap 对照 CmnServiceImpl.getDictMap:chain_common_dict 全量 d_key→d_value
// loImportExcel 5 参 readExcel(单 sheet+字典)+saveIrt:错误包装"运单sheet中，"→
// errs=[导入失败 0];成功→imBatchNumber=UUID、systemCode 顺序、saveBatch(fill)、
// sucInfos、日志(op5),响应=完整 ImportResult
// importPlain 对照 Bill/Contract/ThirdWaybill importExcel:错误→999(data=列表);
// 成功→saveBatch(fill 单独处理)→"导入成功，共N条"
// importTemplate 模板导入:4 参 readExcel(全 sheet,systemCode 传空——4 参重载丢弃编号的原版 bug)。
// bus_* 表全部无 system_code 列,调用方先以 selSystemCode 探测,探测失败走 999。
func (s *Service) importTemplate(w http.ResponseWriter, r *http.Request, table string, fields []impField) {
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
	for _, row := range irt.DataList {
		vals := insertVals(row)
		delete(vals, "id")
		if _, err := s.insert(r, entityForTable(table), vals); err != nil { // warehouse 实体全带 fill
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, "导入成功，共"+strconv.Itoa(len(irt.DataList))+"条")
}

// sysProbe 对照 service.importExcel 首行的 selectSystemCode:表无 system_code 列 → SQL 报错 999
func sysProbe(d *sql.DB, table, left string) bool {
	var s string
	return d.QueryRow("SELECT IFNULL(max(system_code),'') systemCode FROM "+table+" WHERE system_code LIKE ?", left+"%").Scan(&s) == nil
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
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if ext != ".xlsx" && ext != ".xls" {
		irt.ErrInfos = []string{"上传文件格式不正确"}
		return irt
	}
	f, err := fh.Open()
	if err != nil {
		irt.ErrInfos = []string{"解析 Excel 异常，请检查文件格式和内容"}
		return irt
	}
	defer f.Close()
	if ext == ".xls" {
		converted, cleanup, err := platform.OpenLegacyExcel(f)
		if err != nil {
			irt.ErrInfos = []string{"解析 Excel 异常，请检查文件格式和内容"}
			return irt
		}
		defer cleanup()
		f = converted
	}
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
			// notNullMap 耗尽(Java 原版 bug:仅首行可触发"不能为空")
			if len(notNullLeft) > 0 {
				for _, nf := range notNullLeft {
					irt.ErrInfos = append(irt.ErrInfos, sh.name+"sheet中，"+nf.header+"不能为空")
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
func setField(row map[string]any, fl impField, val string) {
	switch fl.kind {
	case "i":
		n, err := strconv.Atoi(strings.TrimSpace(val))
		if err != nil {
			n = 0 // NumberUtils.toInt 默认 0(含字典翻译失效落 0 的原版行为)
		}
		row[fl.col] = json.Number(strconv.Itoa(n))
	case "l":
		n, err := strconv.ParseInt(strings.TrimSpace(val), 10, 64)
		if err != nil {
			n = 0 // NumberUtils.toLong 默认 0
		}
		row[fl.col] = json.Number(strconv.FormatInt(n, 10))
	case "D":
		if _, err := time.Parse("2006-01-02", strings.TrimSpace(val)); err == nil {
			row[fl.col] = strings.TrimSpace(val) // LocalDate.parse ISO;失败被 forEach 吞掉字段留空
		}
	case "d":
		if _, err := strconv.ParseFloat(strings.TrimSpace(val), 64); err != nil {
			return // new BigDecimal 异常被 forEach 内 catch 吞掉,字段留空
		}
		row[fl.col] = json.Number(strings.TrimSpace(val))
	case "t":
		if t, ok := formatTimeAny(val); ok {
			row[fl.col] = t // LocalDateTime 内存态 ISO T;失败 → 字段 null 不报错
		}
	case "b":
		row[fl.col] = boolToBool(val)
	default:
		row[fl.col] = val
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
				out.shared = platform.ExcelSharedStrings(string(data))
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
