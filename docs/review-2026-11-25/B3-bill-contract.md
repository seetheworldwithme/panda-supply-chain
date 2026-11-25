徐先生，bill + contract 两包的只读评审已完成。parity 标签以代码内“照 Java”注释和 docs 为据，未读 Java 源码。

## 范围
- 全文读完在范围 16 个文件（bill 8、contract 8，含 2 个测试），无跳过。
- 另读共享基线：platform 的 resp/dbutil/maps/body/esign/uploadrender、wxconfig.go:30-53、auth 的 middleware/bank、main.go；grep 了 body_contracts_gen.go、db.go、frontend/src。
- DSN（db.go:39）未开 multiStatements/interpolateParams，下文注入点不能堆叠语句。

## 发现（按严重度降序）

### [🔴] delOnce 收到空对象即无条件清空 7 张表 (parity:unknown)
- 位置：backend/internal/contract/contract.go:882-909
- 证据：
```go
where := ""
if len(conds) > 0 {
	where = " WHERE " + strings.Join(conds, " AND ")
}
if _, err := d.Exec("DELETE FROM "+t+where, args...); err != nil {
```
- 场景：任一登录用户 POST `/saas-contract/chainContract/delOnce`，body 为 `{}`（契约只要求非空对象，body_contracts_gen.go:603）。`conds` 为空而 count>0，于是对 chain_contract、chain_capital_settlement、chain_capital_detailed、chain_bill、chain_bill_rate、chain_logistics、chain_logistics_weighing_list 执行无 WHERE 的 DELETE，跨全部租户。`{"im_batch_number":null}` 则删光所有非导入行。
- 修法：`len(conds)==0` 时直接拒绝。
- 把握：high。能否定它的证据：Java 对空 maps 有显式拦截。文件头“行为保持要点”没有把这条列为有意保留。

### [🔴] 请求值直接拼进 SQL：5 个端点加 1 处二阶 (parity:true)
- 位置：
  - contract/stats2.go:389、399、401、408（getLogistics 的 `id`）
  - bill/sett.go:46-47（`searchStr`）
  - contract/sub.go:64（extend `searchStr`）、577、580（files `startTime`/`endTime`）
  - contract/stats.go:132、135（`ctStartTime`/`ctEndTime`）、143（`ids`）、177（`orderId`）
  - contract/stats2.go:256、258（二阶：`chain_capital_detailed.hx_code` 入库后再拼）
- 证据：
```go
" and l.order_id=c.order_id and c.id = " + cid + " and " + concat + " like ?"
q := "select * from chain_bill_settlement where CONCAT(id) like '" + ss + "' order by create_time desc"
sqlCond += " (hx_code like '%" + s + "%'"
```
- 场景：`getLogistics` 传 `{"id":"0 or 1=1 or ? -- "}`，占位符与实参个数仍一致，返回全平台运单。`chainBillSettlement/listByPage` 的 `searchStr` 传 `x' UNION SELECT …-- `，按 11 列对齐即可拖库。`searchFour` 的 detId 分支可用事先写入的 hx_code 做布尔或报错注入。
- 修法：全部改 `?` 参数化，参照同文件 sub.go:97 和 sub.go:411-415（nil 走 `like NULL`）。
- 把握：high。注释自述 Java 也是 `${}` 直拼，属保真复刻，建议列为显式例外。

### [🔴] 按 relation_id 删附件行不带 relation_type，跨实体误删 (parity:unknown)
- 位置：bill/billctl.go:361、391；contract/contract.go:797
- 证据：
```go
_, _ = d.Exec("DELETE FROM chain_common_file WHERE relation_id = ?", id)
"relation_id": strconv.FormatInt(id, 10), "relation_type": 1, "effective": true,
```
- 场景：第二行是 billctl.go:259 的写法，capital/shares.go:27、logistics/waybill.go:247 完全相同。发票、结算单、运单的附件都是 `relation_type=1` 加各自自增 id。删发票 #1000 会连带删掉结算单 #1000 和运单 #1000 的附件记录，`removeContract` 则反向删掉发票和结算单的附件。
- 修法：补 `AND relation_type=?`。三种实体现值都是 1，需先定业务取值。
- 把握：high（Go 行为）。能否定它的证据：Java 的 delByRelId 本就只按 relation_id 删。

### [🟡] xlsx 解析可被极小文件打到 OOM (parity:false)
- 位置：bill/excel.go:570、790-796；入口 excelfields.go:125-130
- 证据：
```go
return io.ReadAll(rc)
if ci, err := colIndex(ref); err == nil {
	for len(row.cells) < ci {
		row.cells = append(row.cells, xlsxCell{kind: "blank"})
```
- 场景：`/saas-bill/chainBill/importExcelJx` 不在 body 契约内，没有 32MB 封顶。sheet 里放一个 `<c r="ZZZZZZZZZ1"/>`，`colIndex` 约 26^9，循环一直追加到内存耗尽。高压缩比的 sharedStrings.xml（zip 炸弹）走 `ReadAll` 同理。Go 单进程没有 -Xmx 兜底，整个后端会被 OOM kill。Java POI 有膨胀比保护、行用 TreeMap，不会这样。
- 修法：`readZip` 用 `io.LimitReader(rc, 64<<20)`，`ci` 上限 16384，上传加 `http.MaxBytesReader`。
- 把握：high。logistics/excel.go 是同源副本（范围外），建议同步。

### [🟡] esignUpload 无上限读文件，事务内无限轮询，失败留孤儿文件 (parity:unknown)
- 位置：contract/esign_file.go:132、153、158、186-205
- 证据：
```go
data, err := io.ReadAll(file)
if err = os.WriteFile(filepath.Join(dir, name), data, 0644); err != nil {
if str(state, "status") == wanted {
case <-time.After(3 * time.Second):
```
- 场景：`/saas-common/chainEqbFlow/uploadFile` 把整文件读进内存，没有大小上限。状态轮询没有次数上限且在 `WithTx` 内：e签宝返回失败态或转换卡住时，只要客户端不断开，就一直占着 1 条连接。20 个这样的请求就抽干连接池。落盘发生在事务之前，后续任何失败都留下没有记录的文件和预览 pdf。
- 修法：读取加 `io.LimitReader`，轮询设最大次数后 `return err`，失败分支删盘文件。
- 把握：high（代码路径）；Java 是否同样写法未知。

### [🟡] 签章免费份数可扣成负数，之后“份数不足”校验永久失效 (parity:unknown)
- 位置：contract/esign.go:131、244-245
- 证据：
```go
if num(etc["sign_free_count"])+num(etc["sign_charge_count"]) == 0 {
if num(etc["sign_free_count"]) > 0 {
_, err = tx.Exec("UPDATE chain_basic_enterprise SET sign_free_count=sign_free_count-1 WHERE id=?", u.CompanyID)
```
- 场景：免费份数为 1 时同一企业并发两次 launchFlow（双击）。REPEATABLE READ 下两个事务都读到 1 并通过校验，相对自减各执行一次，结果为 -1。此后 `-1+0==0` 恒假，可以无限发起签章（外部计费）。
- 修法：`… WHERE id=? AND sign_free_count>0`，校验改为 `<= 0`。
- 把握：medium-high。能否定它的证据：Java 也用相对自减（否则 Java 只会丢更新、不会变负）。

### [🟡] writeVals 把 body 全部键当列写入 (parity:false)
- 位置：bill/bill.go:484-502；调用点 billctl.go:294、330，rate.go:161、186，sett.go:78、113
- 证据：
```go
u := platform.Underscore(m)
for k, v := range u {
	out[k] = v
}
id, err := platform.InsertNoFill(d, "chain_bill", vals)
```
- 场景：本包自己在列表行里输出 `relSearch:false`（bill.go:443-447），还有 `page`、`searchStr`、`filelist`、`timeType` 等载体键（bill.go:30-32）。前端若把行对象原样提交到 add/updateById，就得到 `UPDATE … SET rel_search=0`，MySQL 1054，返回 999。键名含 `.`、空格或中文时触发 `platform.CheckCols` panic（dbutil.go:25），net/http 恢复后直接断开连接。Java 的 exist=false 字段被忽略，所以是成功的。
- 修法：`writeVals` 只保留 `ent.cols` 白名单键，做法同 contract.go:580 的 `pick(m, contractCols...)`。
- 把握：medium。frontend/ 里没有发票 CRUD 页面，无法证实前端载荷；前端不回传载体键即可否定。

### [🟡] 多表写无事务且吞错，失败仍报成功 (parity:unknown)
- 位置：contract/contract.go:792-801、594-640；bill/billctl.go:303、338
- 证据：
```go
_, _ = d.Exec("DELETE FROM chain_contract WHERE id = ?", id)
_, _ = d.Exec("UPDATE chain_capital_settlement SET contract_id=null, contract_code=null, contract_name=null WHERE contract_id = ?", id)
billIds(d, m, id)
```
- 场景：第一条 DELETE 锁等待超时被吞，接口仍回“删除成功”，后面的解绑照常执行，结果是合同还在但结算单和发票已与它解除关联。`billIds` 先删后插，插入失败时返回值被丢弃，关联丢失。这与 dbutil.go:202“Go 统一收进事务防半写”的声明不符。
- 修法：`removeContract` 包 `platform.WithTx` 并返回 error；`billIds` 判返回值。
- 把握：high（Go 行为）。

### [🟡] 通用 CRUD 无租户/属主约束与 mass assignment (parity:true)
- 位置：contract/contract.go:787、580、723；contract/stats.go:86-93；bill/excel.go:92；contract/sub.go:266
- 证据：
```go
s.removeContract(s.pick(r), id)
isPtSearch, _ := m["isPtSearch"].(bool)
if !isPtSearch {
```
- 场景：任何登录用户带 `isPtSearch:true` 就能列出全平台合同。`deleteById` 绕过了 delBatch 的 is_pay 和创建者校验。`pick(m, contractCols...)` 允许写入 `company_id`、`create_user_id`、`is_pay`。导出是全表的 `SELECT *`。`supplierMoney` 可按 id 改任意供应商的应收应付。Java 同样缺失，列此备查。
- 修法：若要偏离 Java，服务端强制 `company_id=u.CompanyID`，并忽略 `isPtSearch`。

### [🟡]（范围外，auth 包）鉴权判定前无上限读 body
- 位置：auth/bank.go:23、52；调用顺序见 auth/middleware.go:71-75，先于 middleware.go:80 的放行判断
- 场景：未登录者带 `CMB-BodySign` 头，或访问路径含 `getSignCallBack`，就能让服务器读完整 body（main.go:60 没有读写超时）。esign_file.go:86 的 `io.ReadAll` 是第二次读，同一段 body 没有上限。
- 修法：两处都加 `io.LimitReader`。

### [🟢] 按 relation_id 连接结算单与发票，金额被放大 (parity:true)
- 位置：bill/billctl.go:557、581；bill/stats.go:643-644；contract/stats2.go:174
- 场景：一张发票对应 N 张结算单时，`sum(includ_money)` 重复计入 N 次，列表出现重复行。
- 修法：先按 relation_id 去重再求和。注释自述是 Java 原表切换，故仅备查。

### [🟢] supplierMoney 用 float64 算金额，且读改写无锁 (parity:false)
- 位置：contract/sub.go:284-287
- 证据：
```go
val := cur + sign*amt
if _, err = d.Exec("UPDATE chain_contract_supplier SET "+col+" = ? WHERE id = ?", val, m["id"]); err != nil {
```
- 场景：DECIMAL 经 ParseFloat 再回写，大额时丢精度；并发加款丢更新（后者 Java 也无锁）。
- 修法：`SET col=IFNULL(col,0)+?`。

### [🟢] e签宝出站路径由 body 拼接 (parity:unknown)
- 位置：contract/esign.go:20、69-70、82-84
- 场景：`eqbFlowId` 传 `../…`，可让服务器用自己的签名请求 e签宝主机的任意路径；:82 的 phone 未转义。
- 修法：校验 `^[A-Za-z0-9]+$`。

### [🟢] exportHandler 全表驻留内存 (parity:true)
- 位置：bill/excel.go:92、102-115
- 场景：`SELECT *` 没有租户条件，整表 grid 和整张 sheet 字符串都在内存里。
- 修法：流式写 sheet。

### [🎉] 值得保留的设计
- bill/stats.go:19-135 的 bdec 系用 math/big 精确复刻 BigDecimal 的 add/divide（HALF_UP），税负统计不落浮点。
- bill/excel.go:29、85-88 导出单元格一律 `inlineStr` 加 `xmlEscape`，没有公式注入面。

## 重复模式
- `${}` 直拼：见第二条 SQL 注入清单。
- 多表写吞错、无事务：contract.go:594-640、731-763、792-801；billctl.go:265-283、303、338、353-364、383-394。
- 全量 body 键入列：bill.go:484-502。
- 按 id 直操作、无属主校验：billctl.go:228、316、345；rate.go:136、175、197；sett.go:56、90、102；sub.go:75、97、254、266、441、530、648；dict.go:108、131；esign_file.go:26（`contractIds` 可引用任意文件 id）。

## 已核对无问题
- 除上述拼接点外，值参全部参数化；list 路径的 `allEqCond` 白名单挡住了 CheckCols panic。
- `pageOf` 保证 size≥1，`IPage` 不会除零；OFFSET 溢出只会得到 999。
- rows 都经 `ScanMaps`（内部 defer Close），没有提前 return 泄漏，也没有遍历 rows 中嵌套查询。
- `bparts` 遇到非法串退化为 0，不会 panic；除零返回 999，已标注与 Java 一致。
- `esignFlow` 的写入列走 `bodyVals` 白名单；`WithTx` 带 defer Rollback。
- 回调验签用 HMAC 常量时间比较、±10 分钟窗口，无 appSecret 即拒绝；`EsignJSON` 响应封顶 16MB、30s 超时、Body 必关。
- 测试无空断言：excel_columns_test.go 有效；esign_test.go 的主用例需 `PANDA_LOCAL_AUDIT=1` 和本机 MySQL，默认 skip。