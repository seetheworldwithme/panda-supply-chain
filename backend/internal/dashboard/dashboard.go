// Package dashboard 仪表盘统计接口(企业 Main/平台 PtIndex/政府大屏 Xzmain)。
// 接口对照 Java saas-common/saas-bill/saas-capital/saas-contract/saas-logistics 的 count 系列,
// SQL 逐条翻译(见各 handler 注释的 Java mapper 行号)。
//
// 租户过滤说明:Java TenantLineInnerInterceptor 的 ignoreTable 白名单(coreConfig.getIgnoreTables,
// Nacos 配置)在本环境覆盖了全部仪表盘业务表——实测 hyqcbjhd(企业,company_id=405)统计接口返回的
// 均为全平台数据(2026-11-16 与 8765 Java 网关逐接口对比证实)。故 Go 版统计 SQL 一律不加
// company_id 租户条件,仅保留 Java service 层手工拼的 companyId 语义(如 selectBusCount 的 jwt 分支)。
package dashboard

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

type Service struct {
	DB *db.DB
}

func NewService(database *db.DB) *Service { return &Service{DB: database} }

// pick 按请求头 applytype 选库(默认 chain,zjjg 切监管库),对照 Java DynamicDataSource
func (s *Service) pick(r *http.Request) *sql.DB {
	if r.Header.Get("applytype") == "zjjg" {
		return s.DB.Zjjg
	}
	return s.DB.Chain
}

// Register 挂载仪表盘路由(路径与 Java 网关后的服务路径一致)
func Register(mux *http.ServeMux, s *Service) {
	// ── saas-common:chainBasicMainPage / chainBasicEnterprise / chainAreaInfo / upload ──
	post(mux, "/saas-common/chainBasicMainPage/selEtcMsg", s.selEtcMsg)
	post(mux, "/saas-common/chainBasicMainPage/countActBus", s.countActBus)
	post(mux, "/saas-common/chainBasicMainPage/countMonActBus", s.countMonActBus)
	post(mux, "/saas-common/chainBasicMainPage/countMonData", s.countMonData)
	post(mux, "/saas-common/chainBasicMainPage/countNetLogist", s.countNetLogist)
	post(mux, "/saas-common/chainBasicMainPage/countEtcTax", s.countEtcTax)
	post(mux, "/saas-common/chainBasicMainPage/countLogistArea", s.countLogistArea)
	post(mux, "/saas-common/chainBasicMainPage/countPriceTrend", s.countPriceTrend)
	post(mux, "/saas-common/chainBasicEnterprise/ptlist", s.ptlist)
	post(mux, "/saas-common/chainBasicEnterprise/selectCount", s.entSelectCount)
	// getById(getEtc)由 authadmin 包注册:企业档案页需要嵌套 pubFlows/urcs 等完整实体形状,
	// 此处一期简化版(snake 裸行)已删除,避免同路径重复注册
	post(mux, "/saas-common/chainAreaInfo/selByParent", s.selByParent)
	post(mux, "/saas-common/upload/getImgByRelId", s.getImgByRelId)

	// ── saas-bill:chainBill / chainBillRate ──
	post(mux, "/saas-bill/chainBill/selJxzb", s.selJxzb)
	post(mux, "/saas-bill/chainBill/selZeroTax", s.selZeroTax)
	post(mux, "/saas-bill/chainBill/selectBusCount", s.selectBusCount)
	post(mux, "/saas-bill/chainBill/countTaxRate", s.countTaxRate)
	post(mux, "/saas-bill/chainBill/countBuses", s.billCountBuses)
	post(mux, "/saas-bill/chainBill/countMonBuses", s.billCountMonBuses)
	post(mux, "/saas-bill/chainBill/countPtTax", s.countPtTax)
	post(mux, "/saas-bill/chainBillRate/getSumPrice", s.rateGetSumPrice)
	post(mux, "/saas-bill/chainBillRate/countMonBuses", s.rateCountMonBuses)

	// ── saas-capital:三条统计已移交 capital 包精确版(五期),dashboard 不再注册 ──

	// ── saas-contract ──
	post(mux, "/saas-contract/chainContract/selectQuantity", s.cttSelectQuantity)
	post(mux, "/saas-contract/chainContract/countPro", s.cttCountPro)

	// ── saas-logistics ──
	post(mux, "/saas-logistics/chainLogisticsBill/countLogist", s.loBillCountLogist)
	post(mux, "/saas-logistics/chainLogistics/countBuses", s.loCountBuses)
	post(mux, "/saas-logistics/chainLogistics/countMonBuses", s.loCountMonBuses)
}

func post(mux *http.ServeMux, path string, h http.HandlerFunc) {
	mux.HandleFunc("POST "+path, h)
}

// ── 通用小工具 ──

// bodyMap 解析 JSON body 为宽松 map(前端传参字段类型不统一,数字/字符串都可能)
func bodyMap(r *http.Request) map[string]any {
	var m map[string]any
	_ = json.NewDecoder(r.Body).Decode(&m)
	if m == nil {
		return map[string]any{}
	}
	return m
}

// str 取字符串字段;JSON 数字(前端 timeType/companyId 常为数字)转字符串
func str(m map[string]any, k string) string {
	v, ok := m[k]
	if !ok || v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case float64:
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10)
		}
		return strconv.FormatFloat(t, 'f', -1, 64)
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

// f64 取数值字段
func f64(m map[string]any, k string) float64 {
	if v, ok := m[k].(float64); ok {
		return v
	}
	s := str(m, k)
	n, _ := strconv.ParseFloat(s, 64)
	return n
}

// f64Ptr 取数值字段指针(区分"没传"和 0)
func f64Ptr(m map[string]any, k string) *float64 {
	if _, ok := m[k]; !ok || m[k] == nil {
		return nil
	}
	v := f64(m, k)
	return &v
}

// scanMaps 通用行扫描(转发 platform.ScanMaps,一期遗留别名)
func scanMaps(rows *sql.Rows) ([]map[string]any, error) {
	return platform.ScanMaps(rows)
}

// num 把 driver 返回的 DECIMAL([]byte)/int64/float64 统一转 float64
func num(v any) float64 {
	switch t := v.(type) {
	case nil:
		return 0
	case []byte:
		n, _ := strconv.ParseFloat(string(t), 64)
		return n
	case string:
		n, _ := strconv.ParseFloat(t, 64)
		return n
	case int64:
		return float64(t)
	case float64:
		return t
	default:
		return 0
	}
}

// numStr DECIMAL 值转 Java BigDecimal.toString 风格字符串(响应里字符串数字的场景,如 countPtTax)
func numStr(v any) string {
	if v == nil {
		return "0"
	}
	return strconv.FormatFloat(num(v), 'f', -1, 64)
}

// monthSeries 生成 [start月 .. end月] 的 "yyyy-MM" 序列,替代 Java 补月循环与 @cdate 月序列
func monthSeries(start, end string) []string {
	s, ok := parseDay(start)
	if !ok {
		return nil
	}
	e, ok := parseDay(end)
	if !ok {
		return nil
	}
	var out []string
	for t := s; !t.After(e); t = t.AddDate(0, 1, 0) {
		out = append(out, t.Format("2006-01"))
		if len(out) > 1200 { // 防御异常超大范围
			break
		}
	}
	return out
}

// daySeries 生成 [start .. end] 的 "yyyy-MM-dd" 序列(countDetail 日粒度用)
func daySeries(start, end string) []string {
	s, ok := parseDay(start)
	if !ok {
		return nil
	}
	e, ok := parseDay(end)
	if !ok {
		return nil
	}
	var out []string
	for t := s; !t.After(e); t = t.AddDate(0, 0, 1) {
		out = append(out, t.Format("2006-01-02"))
		if len(out) > 3700 {
			break
		}
	}
	return out
}

func parseDay(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if len(s) > 10 {
		s = s[:10]
	}
	t, err := time.Parse("2006-01-02", s)
	return t, err == nil
}

// monthDiff 对照 Java monthDiff:跨月数(同月=0)
func monthDiff(start, end string) int {
	s, ok1 := parseDay(start)
	e, ok2 := parseDay(end)
	if !ok1 || !ok2 {
		return 0
	}
	return (e.Year()-s.Year())*12 + int(e.Month()) - int(s.Month())
}

// xzqhCond 行政区划过滤条件(Java 手工拼 like 的参数化版本,防注入)。
// Java 判断 !"0".equals(xzqhCode);空/0 不加条件
func xzqhCond(m map[string]any) (string, string) {
	x := str(m, "xzqhCode")
	if x == "" || x == "0" {
		return "", ""
	}
	return " and CONCAT(province_code,',',city_code,',',area_code,',',county_code,',') like ?", "%" + x + ",%"
}

// timeCond concat(month,'-01') 与 startTime/endTime 的时间段条件(countTaxRate/countMonBuses 等共用)
func timeCond(m map[string]any, col string) (string, []any) {
	var cond []string
	var args []any
	if s := str(m, "startTime"); s != "" {
		cond = append(cond, " and concat("+col+", '-01') >= ?")
		args = append(args, s)
	}
	if s := str(m, "endTime"); s != "" {
		cond = append(cond, " and concat("+col+", '-01') <= ?")
		args = append(args, s)
	}
	return strings.Join(cond, ""), args
}

// timeTypeCond timeType 1本月 2上月 3近半年 4本年(无条件),对照各 Java service selectQuat
// (原 capital.go 迁出后留此供 contract.go 使用)
func timeTypeCond(m map[string]any) string {
	switch str(m, "timeType") {
	case "1":
		return " AND month = DATE_FORMAT(CURDATE(), '%Y-%m')"
	case "2":
		return " AND month = DATE_FORMAT(DATE_SUB(CURDATE(), INTERVAL 1 MONTH), '%Y-%m')"
	case "3":
		return " AND concat(month, '-01') BETWEEN date_sub(now(), INTERVAL 6 MONTH) AND now()"
	}
	return ""
}
