徐先生，capital 包只读评审已完成（6 个回合、20 次调用）。

## 范围
- 全文读完 14 个文件：`backend/internal/capital/` 全部 14 个（4,269 行），外加共享 helper（`platform/resp.go`、`dbutil.go`、`maps.go`、`body.go`、`auth/middleware.go`）和 `docs/validation/java-go-20261118/缺失接口与兼容差异.md`。
- `auth/bank.go` 只按 grep 看了 19-55 行；Java 源码不在仓库内，parity 均依据注释与 docs 判定。
- 以下各项均未登记在 docs 3.2「原 Java 问题例外」中。docs 3.3/4 标注招行生产成功路径「未实现、不是已完成迁移」。

## 发现

### [🔴] 1. countSeries 月份循环不退出，内存撑爆全进程 (parity:false，大概率)
- 位置：`backend/internal/capital/detailed2.go:194-200`（`stCountSett` 走此分支）、`:203-212`
- 证据：
```go
for i := 0; ; i++ {
    t := st.AddDate(0, i, 0)
    series = append(series, t.Format("2006-01"))
    if t.Format("2006-01") == et.Format("2006-01") {
        break
    }
```
- 场景：`POST /saas-capital/chainCapitalCgSettlement/countSett`（Xs 同）带 `{"startTime":"2026-01-31","endTime":"2026-02-10"}`。Go 的 `AddDate` 把「2 月 31 日」进位成 03-03，之后依次得到 "2026-03"、"2026-05"，"2026-02" 永不出现，循环不退出，`series` 一直追加直到 OOM。startTime 晚于 endTime 时同样不退出。任何登录账号都能触发，真实用户选「31 日到次月」区间也会触发。
- 修法：按年月算术生成并限定次数：`for i := 0; i <= jg; i++`，用 `mo := int(st.Month())+i-1` 和 `fmt.Sprintf("%04d-%02d", st.Year()+mo/12, mo%12+1)`；`jg<0` 返回空序列。
- 把握：high。Java 若用 `LocalDate.plusMonths`（月末日会夹到月末）则 Java 不挂；`bill/stats.go:258,515`、`logistics/bill.go:508` 有同款循环，建议对照。

### [🔴] 2. listByPage 的 searchStr 原生内插，构成 SQL 注入 (parity:true，注释标明刻意保留)
- 位置：`backend/internal/capital/small.go:338-343`
- 证据：
```go
ss = "%" + str(m, "searchStr") + "%"
q = "select * from chain_capital_detail_settlement where CONCAT(id) like '" + ss + "' order by create_time desc"
countQ = "SELECT count(1) FROM chain_capital_detail_settlement where CONCAT(id) like '" + ss + "'"
```
- 场景：`{"page":{"current":1,"size":10},"searchStr":"' OR 1=1 -- "}` 返回全表（跨租户）；`' UNION SELECT … -- ` 可经 remark 列带出 `chain_auth_user` 等表（count 语句先执行，UNION 需单列）。
- 修法：改成 `like ?` 加同一个 `%v%` 参数。除含 `'` 或 `\` 的输入（Java 本就是 SQL 报错 999）外，行为逐字相同。是否保留由产品决策。
- 把握：high。

### [🔴] 3. 招行退款/确认/查单只信请求体，不校验订单归属 (parity:unknown，倾向 true)
- 位置：`backend/internal/capital/payment.go:289-325`（refund）、`:175-216`（confirmPay）、`:328-344`
- 证据：
```go
res, err := bankPost(r, "/order-bck", bankRefundBody(row)) // row 来自 refundList 元素
...
if row["bus_status"] != "WCF" {
res, err := bankPost(r, "/order-cfm", bankConfirmBody(m))    // mchNbr/订单号/金额取自请求体,非 row
```
- 场景：任意登录用户用自己的短信码（`bankCode` 只校验本人手机）POST refund，带上任意 `mchNbr`、`oriMchOrderNbr`、`oriOrderAmt`，平台就会用自家 appid 和私钥向招行发全额退款。全程不查本地流水。包内 `u.CompanyID` 只在 `small.go:289/292` 作可选筛选，没有任何归属比对。
- 修法：发银行请求前按 `mch_order_nbr` 取本地 transinfo，要求 `pay_unit_id = u.CompanyID`；confirm/cancel 改用库中订单号。这会改变 Java 行为，建议在招行对接上线前决策。
- 把握：high（Go 路径确定）；Java 是否同样缺失未见源码，但文件头称「按原 Service 保留」，没有删减校验的迹象。

### [🟡] 4. 招行回调：带 JWT 的请求不验签即可伪造 (parity:false/unknown)
- 位置：`payment.go:379-383`；`auth/middleware.go:71-73`
- 证据：
```go
if r.Header.Get("CMB-BodySign") == "" {
    platform.Unknown(w)
    return
}
m := bodyMap(r)
```
  中间件里 `if checkBankSignature(r) { user.SignedAnon = true }` 只影响匿名放行。
- 场景：任意登录用户带 Bearer、`CMB-BodySign: x` 和 `{"msgNbr":"9","mchOrderNbr":"他人订单","status":"FIN","orderType":"PAY","rtnFlag":"S"}`，即可把该订单写成 `bus_status='FIN', status=2`（付款完成）。
- 修法：中间件把验签结果放进 context 或 `CurrentUser`，handler 改为「未验签即 999」。
- 把握：high。

### [🟡] 5. Excel 导入三处无资源上限，登录用户可打死进程 (parity:false)
- 位置：`excel.go:142-148`、`:572-579`、`:797-804`
- 证据：
```go
_, fh, err := r.FormFile("file") // 无 RemoveAll、无大小上限
return io.ReadAll(rc)            // zip 条目解压无上限
for len(row.cells) < ci {        // ci 来自 <c r="…">
```
- 场景：
  - import 路由不在 `body_contracts_gen.go`，跳过了 32MB 封顶。超过 32MB 的部分落为 /tmp 临时文件且永不删除（`contract/esign_file.go:124`、`commonrest/external.go:479` 都有 `RemoveAll`，此处漏了），反复上传可写满磁盘。
  - 高压缩比 xlsx 经 `io.ReadAll` 解压到内存，直到 OOM。
  - 几百字节文件里放 `<c r="ZZZZZZZZZ1"/>`，ci≈5.4e12，回填循环直到 OOM。
- 修法：`formFile` 后加 `defer r.MultipartForm.RemoveAll()` 和 `MaxBytesReader`；`readZip` 套 `io.LimitReader`；`ci > 16384` 时忽略。
- 把握：high。同源解析器在 bill、commonrest、warehouse、logistics 各有一份。

### [🟡] 6. 小控制器 add/updateById 把全部 body 键当列 (parity:false)
- 位置：`small.go:392`、`:424` 调 `shares.go:59-64`
- 证据：
```go
u := platform.Underscore(m)
for k, v := range u {
    out[k] = v
}
```
- 场景：`/chainCapitalAccounts/add` 带 `{"accountName":"a","a-b":1}` 触发 `CheckCols` panic。HTTP 链路没有 recover（全仓只有 `jobs.go:174`），连接被掐、stderr 打栈，而不是返回 JSON 999。带 `{"id":5,"transType":"1"}` 时生成 `trans_type=?`，报 Unknown column；Java 的 MP 会忽略临时字段，更新成功。影响 7 个控制器 × 2 端点。
- 修法：改用 `pick(m, def.ent.cols...)`，再按 `ent.renames` 两行补驼峰列。
- 把握：high；前端是否带载体键未核实。

### [🟢] 7. 统计查询吞掉 ScanMaps 错误 (parity:false)
- 位置：`detailed2.go:161`、`settlement.go:579`
- 证据：`ms, _ := platform.ScanMaps(rows)`
- 场景：rows 迭代中断连时返回全 0 的成功响应。
- 修法：判 err 后走 `platform.Unknown(w)`。
- 把握：high。

### [🟢] 8. bucketSeries 补月循环不推进起点 (parity:unknown)
- 位置：`detailed2.go:274-282`
- 证据：`cur := stT.Format("2006-01")`（每轮都是起点月）
- 场景：timeType=4 且 startTime=2025-01-01、endTime=2025-03-01 时，尾部追加 3 条重复的 `{"2025-01",0}`，中间空月不补。
- 修法：用 `stT.AddDate(0, i, 0)` 的结果作 cur。
- 把握：high；若 Java 原码有 `startTime = startTime.plusMonths(1)`，则是移植遗漏。

## 重复模式
- **A 无归属/租户校验，实体全列直写（parity:unknown）**：getById、updateById、deleteById 按 id 直取直改直删（`detailed.go:228-326`、`settlement.go:254-396`、`small.go:154-190`）；`list` 和 `exportExcel` 全表（`excel.go:95`）；updateById 可直接改 `money/khxje/hx_status/company_id`，绕过核销流程。Java MP 是同款调用，倾向同样缺失，docs 无 TenantLine 记录，需产品确认。
- **B 多表写无事务（parity:true）**：`detailed.go:313-323`、`settlement.go:344-358`、`:687-731`（记录已插、后续失败不回滚）、逐行导入。包内只有 `balance.go:34/196` 和 `payment.go:385` 用了 `WithTx`，与 dbutil 注释「Go 统一收进事务」的意图不符。
- **C Java 原缺陷照搬（parity:true，仅登记）**：
  - 核销金额全由前端给出，没有正数或 sjhxje≤khxje 校验，没有行锁，hxbh 取 max 不加锁（`balance.go:46-99`、`:307-324`）。
  - `dtGetSumPrice` 按 `d.hx_code=s.hx_code` 连接，一单多结算时金额成倍放大（`detailed2.go:66-79`）。
  - deleteById 按 `relation_id` 删 `chain_common_file` 不带 relation_type，会删到其他模块的同号附件（`detailed.go:321`）。
  - cancelPay 状态不符仍发银行请求，成功计数照加（测试已固化）。
- **D** `bankCallback` 在事务内发起 30s 超时的银行查询（`payment.go:399`），占用 20 条连接的池。

## 已核对无问题
- 请求值一律参数化并限列（`allEqCond`/`pick`/`whereClause`），除发现 2 外没有请求值进 SQL 文本。json.Number 作参数在 go-sql-driver v1.10.1 可用。
- 金额全程用字符串十进制（`dec*`），聚合交给 SQL DECIMAL，没有 float 运算。`pageOf` 保证 size≥1，没有除零；OFFSET 溢出只会得到 999。
- 招行外呼：按发送字节做 SHA256 摘要和 RSA 签名；appsecret 只进签名源串，不进报文、日志或响应；30s 超时，Body 必关，host 只取环境变量，无 SSRF。
- 回调：msgNbr 单调判断挡住顺序重放，bkcall 与 transinfo 同事务回滚（测试用触发器验证）。
- 导出：单元格一律 inlineStr 并做 xmlEscape，无公式注入，文件名不落盘。
- rows 全部经 ScanMaps 关闭，没有迭代中嵌套查询；`WithTx` 内没有 panic 路径。
- 测试：三份测试都能失败、断言方向正确；`TestBankPaymentPersistsJavaTransaction` 默认 skip，未覆盖「JWT + 伪造 CMB-BodySign」。
- 🎉 可保留：`dec*` 字符串十进制；`entShape` 把列白名单、回显和日期格式收敛在一处。