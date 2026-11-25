徐先生，评审完成（只读，未改动任何文件）。

## 范围
完整读取 12 个范围内文件（excel.go、product.go、product_list.go、product_get.go、product_write.go、product_bp.go、product_sync.go、product_external.go 与 4 个测试）。另按需读了 platform/{resp,dbutil,maps}.go、auth/middleware.go、规格书全文、template.go 300-499/780-811、warehouse/jobs.go、platform/jobs.go 120-199、legacy_excel.go 18-53，并 grep 核对调用点。未读 HanLP 移植、WechatJSON、SendLegacySMS、red_transfer.go（非本范围）。已知的 colIndex 补列、readZip 无上限不重复报。

## 发现（按严重度降序）

### [🔴] 请求体"键名"被当作列名拼进 SQL (parity:false)
- 位置：backend/internal/warehouse/product_bp.go:327-343（selCheckPro）；product.go:29-49（mapGetSql3，allProductType :87）；根因 template.go:37-44
- 证据：
```go
if e.typeOf(k) == "x" || k == "page" || k == "ids" { continue }
bv, ok2 := bindVal(e.typeOf(k), v)
if !ok2 { platform.Unknown(w); return }
sqlStr += " and " + camelOf(e, k) + " like ?"
// template.go:43  return "s"   // 实体里查不到的键默认当字符串放行
```
- 场景：任意登录账号向 `/saas-warehouse/busProduct/selCheckPro` 发 `{"id>0 or if((select ascii(substr(password,1,1)) from chain_auth_user where id=1)>64,sleep(3),0) or id":"x"}`，生成 `where 1=1 and id>0 or if(...) or id like ?`，按响应延迟逐位拖库；`{"id>0 or id":"x"}` 直接绕过其余条件。Java getSql(Map,Class) 对实体里没有的键直接忽略（规格 1.3）；Go 下前端多传一个未知键也会 Unknown column → 999。
- 修法：两处循环只处理 `e.q` 中存在的键，其余 continue。
- 把握：high。meta_gen 全库仅 6 个 "x" 类型，覆盖不到任意键。

### [🔴] bpDelBatch 不校验 ids 元素就拼 SQL (parity:unknown，疑 Go 独有)
- 位置：backend/internal/warehouse/product_write.go:502-506
- 证据：
```go
ids, _ := m["ids"].([]any)          // 无 validateBody
id := str(map[string]any{"v": iv}, "v")
d.Exec("INSERT INTO bus_del_product SELECT * FROM bus_product WHERE id=" + id)
```
- 场景：`{"ids":["1 OR 1=1"]}` 把全部在售煤种复制进"已删除备份"表；`["0 OR IF(条件,1,0)"]` 之后 `GetByID` 是否命中决定响应是"删除成功"还是 999，可做布尔盲注。
- 修法：开头加 `if !entBusProduct.validateBody(m) { platform.Unknown(w); return }`，与 bpDeleteById 一致。
- 把握：Go 行为 high。parity 为 medium：validateBody 注释"ids 必须整数数组"暗示 Java 的 ids 是整数列表，字符串在 Jackson 反序列化即失败；若 Java 是 `List<String>` 则原版同样可注入。

### [🔴] bpUpdateById 用驼峰键读 snake_case 原始行，全部取空 (parity:false)
- 位置：backend/internal/warehouse/product_write.go:363-373、388-392、407-414、461-463；product_sync.go:157-162
- 证据：
```go
bpRows, _, err := s.queryEnt(r, "SELECT * FROM bus_product WHERE id=?", id) // ScanMapsCols 原始列名
bp := bpRows[0]
associationId := str(bp, "associationId")   // 行里只有 association_id，恒 ""
if associationId == "" { associationId = uuid32(); _, _ = d.Exec("UPDATE bus_product SET association_id=? ...
```
- 场景：与规格 3.6/2.2 相悖。
  - 每次 updateById 都重生成 association_id，未随 body 回传的化验单、图片、视频附件脱离该煤种。
  - body 带非空 createPhone 时 `cp != str(bp,"createPhone")` 恒真，落入 NPE 分支恒 999。
  - lastCheckTime/downTime 永不写。
  - 改价积分、首次发布标记、审核短信（orgStatus 恒空串）全部失效。
- 修法：`bp := platform.Camelize(bpRows[0])`（同 product_sync.go:27）。
- 把握：high。`str` 无 snake 回退；product_sync_test.go:88 手工传驼峰 bp，没覆盖。前端是否必传 createPhone 未核实，若必传则常规编辑整体不可用。

### [🟡] Excel 导入逐行插入无事务，中途失败半提交 (parity:false)
- 位置：backend/internal/warehouse/excel.go:234-241（importTemplate）；product_write.go:582-592（bpImportExcel）
- 证据：
```go
for _, row := range irt.DataList {
    if _, err := platform.Insert(s.pick(r), table, vals); err != nil {
        platform.Unknown(w)
        return
    }
}
```
- 场景：导入 500 行，第 300 行 varchar 超长或唯一键冲突，前 299 行已提交并返回 999；用户修正后重导，前 299 行翻倍。Java 的 saveBatch 自带 `@Transactional(rollbackFor=Exception.class)`，整体回滚。
- 修法：循环包进 `platform.WithTx`，改用 `platform.TxInsert`。
- 把握：high。

### [🟡] soffice 子进程并发无上限，分享图每次落盘不清理 (parity:false；落盘增长 parity:true)
- 位置：excel.go:314-321 → platform/legacy_excel.go:42-44；product_external.go:362-363、388-397
- 证据：
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
exec.CommandContext(ctx, bin, "-env:UserInstallation="+profile, "--headless", "--convert-to", "xlsx", ...)
name := strconv.FormatInt(time.Now().UnixMilli(), 10) + ".jpg"
```
- 场景：登录账号并发 30 次上传 .xls 或循环调 getDetailImg，30 个 soffice（各数百 MB）跑到 30s/60s 超时，8G 机器 OOM。getDetailImg 每次还在 PANDA_UPLOAD_DIR 永久落一张新图并插一行 chain_common_file，无去重。Java 用 POI 进程内解析，无子进程。
- 修法：包级信号量 channel（容量 1~2）供两个调用点共用；分享图按产品 id 加短期缓存。
- 把握：high（无任何限流）。

### [🟡] 公众号通知逐 openid 取 token，占住唯一 worker (parity:false)
- 位置：product_external.go:261-276、31-36；platform/wechat_http.go:46-63；platform/jobs.go:149-183
- 证据：
```go
for _, u := range users {               // 全部 official_openid
    ...
    if err = s.bpTemplate(r, body); err != nil {   // 每次都 bpWXToken
```
- 场景：WechatAccessToken 无缓存，每次查库并 POST stable_token。N=3000、每次 150ms，一个任务串行 2N 次 HTTPS 约 15 分钟。Jobs.Run 是单 worker 且每条 Sleep 2s，其间所有煤种的价格历史、审核消息、红包全部排队；stable_token 日配额打满后当天微信功能全挂。
- 修法：循环外取一次 token 传给 bpTemplate；WechatAccessToken 进程内缓存到过期前 60s。
- 把握：medium-high；若 WechatJSON 内已有缓存或熔断则影响减轻（未读）。

### [🟡] 异步同步任务整行回写，且 zjjg 请求落到 chain 库 (parity:true)
- 位置：product_sync.go:20、33-84；product_write.go:458-463；jobs.go:26-27
- 证据：
```go
job.Header.Del("applytype") // 原 MQ 消费者没有请求头，使用默认业务库。
...
if _, err = platform.UpdateByID(d, "bus_product", proID, pick(p, entBusProduct.dbcols)); err != nil {
```
- 场景：
  - 任务先快照、再发短信和微信、最后整行回写。期间审核下架（status=3）会被旧快照覆盖回 1。
  - zjjg 租户更新煤种后，任务在 chain 库按同一数字 id 命中另一租户的煤种，改写它并给其创建人发奖、发消息。前提是 zjjg 库有 bus_product，未核对。
- 修法：回写只 SET 实际改动列（is_certified_etc、group_id、sec_name、th_name、pre_*_price、*_upd_time、last_check_fac_price）。
- 把握：机制 high，影响 medium。

### [🟡] 点击记录、收藏、分享红包 check-then-insert 无锁 (parity:true)
- 位置：product_get.go:399-420、386-397；product.go:346-359；product_get.go:457-466
- 证据：
```go
if ucCnt == 0 { d.Exec("INSERT INTO bus_user_click_pro ...") }   // 并发两次首访各插一行
...
if ucCnt > 1 { platform.Unknown(w); return false }               // 此后该用户每次详情页恒 999
```
- 场景：新用户首次同时打开两个煤种详情，bus_user_click_pro 出现两行，之后 getById 永久 999（除非库有唯一索引，未核对）。收藏双击同理。新注册用户并发带不同 shareUserId 可多次领分享红包，且收款人由调用方任意指定。
- 修法：INSERT … ON DUPLICATE KEY UPDATE 配唯一索引，或 `UPDATE` 影响行数为 0 再 INSERT。
- 把握：medium。

### [🟢] selCheckPro 对 searchStr 直接断言 string (parity:false)
- 位置：product_bp.go:346-348
- 证据：`inner = "%" + v.(string) + "%"`
- 场景：`{"searchStr":123}` 得到 json.Number，断言 panic。该 handler 无 validateBody，连接被断开、客户端空响应；Java 的 Map 取值正常。
- 修法：改用 `bindVal("s", v)`，`!ok` 时 Unknown。
- 把握：high。

### [🟢] bpLog/bpSaveDelLog 失败后调用方不返回 (parity:false)
- 位置：product_write.go:186-193、169-173、453-464、515-526
- 证据：`s.bpSaveDelLog(w, r, d, bp)` 之后继续执行批量 DELETE 和 `platform.Success`
- 场景：日志插入失败时响应先 500 再追加一段 success JSON；delBatch 仍删除煤种，而 Java 在 saveDelLog 抛异常处中止。
- 修法：两个函数返回 error，调用处失败即 return。
- 把握：high。

### [🟢] 遗留调试 println (parity:false)
- 位置：product_get.go:143
- 证据：`println("DBGSTATUS raw=", ...)`
- 场景：每次详情页向 stderr 写一行，nohup 下日志无限增长。删掉即可，`fmt` 在 :162 仍在用。
- 把握：high。

### [🟢] inlineStr 富文本或注音只保留最后一个 `<t>` (parity:false)
- 位置：excel.go:935-939
- 证据：`case "t": if inCell { inT = true; tBuf.Reset() }`
- 场景：`<is><r><t>云仓</t></r><r><t>号</t></r></is>` 只得到"号"，表头匹配失败，整列静默丢弃。sharedStrings 路径正确，仅 inline 路径有问题。
- 修法：改为在 `<is>` 开始时 Reset，并跳过 `<rPh>`。
- 把握：high，属边缘路径。

### [🟢] 自定义数字格式日期启发式误判 (parity:unknown)
- 位置：excel.go:826-832、984-986
- 证据：`strings.Contains(low, "y") || strings.Contains(low, "d") || strings.Contains(low, "h")`
- 场景：`#,##0.00;[Red]-#,##0.00` 或 `0.00"USD"` 被当日期，5500 变 "1915-01-11"，入 decimal 列报错并触发上面的半提交。
- 修法：判定前剔除 `[...]` 与 `"..."` 片段。
- 把握：medium（Java 侧判定方式未见）。

## 重复模式
- **`${}` 裸拼接（parity:true，注释与规格 0.5 登记为刻意保留）**：
  - 无类型校验、任意字符串可注入：product.go:61、64、93、96（parent_id/fst_pro_id）、:112、115（groupId/wareIdentity）。
  - province/city/area_code 类型是 "s"：product_list.go:232-234，`"0 or 1=1"` 可注入。规格称"非数字→999"，不成立。
  - 其余裸拼点（product_list.go:210-219、273、286、327；product_bp.go:121-127）已被 bindVal 数字校验或来自 JWT/库。
  - 建议：语义为 id 或编码的键加 `^\d+$` 白名单。
- **按 id 读写不校验归属（parity:true，规格 0.2 明示 companyId 从未使用）**：
  - 涉及 updateById、deleteById、delBatch、shareTime、addViewTime、addPhoneTime、getDetailImg。
  - updateById 可用 body 的 createUserId/createPhone 换主。
  - add/updateById 的 status、isPlatformCertification、isRealCheck、effective 可由 body 直写（自审核）。
  - getById 的 userMap/brokerList 原样带出 chain_auth_user 全行，含 password 列（product_get.go:27-31）。请抽查 8765 响应确认 Java 是否一致。
- **exportHandler 全表导出、全量入内存、无租户过滤（parity:true，规格 3.8）**：excel.go:148，约 70 个控制器共用。
- **读-改-写计数器非原子（parity:true）**：sharing/view/phone_times（product_write.go:615、649、669）、mo_get_sum_score（:227、268、308）。
- **N+1（parity:true）**：bpAttachTags 每行一次查询（product_bp.go:35-44），无 page 分支为全表 ×1。串行执行，不耗尽连接池。

## 已核对无问题
- writeXlsx 全部 inlineStr 加 xmlEscape，无公式注入；文件名固定；列名和表名来自元数据。
- 列映射只写 impField.col 固定列，id 被剔除，CheckCols 收不到请求派生键。
- ipageOf/pageOf 对 size≤0 有保护，无除零；负 size 不分页照 MP。
- rows 均经 ScanMaps 关闭；无迭代中嵌套查询。
- 外呼 URL 固定或取自环境变量，openid/token 经 url.Values 编码，无 SSRF；subprocess 参数不含请求值；SVG 文本已转义。
- LibreOffice 转换有 32MB 上限、30s 超时、独立 profile、临时目录清理。👍
- 金额计算全部走 big.Rat（bpPriceDifference、decimalSubtractExact）。👍
- Jobs 已在 main.go:126 接线，platform/jobs.go:174 有 recover；裸 goroutine 仅 Jobs==nil 时走到。
- productBroadcast、redCachePut 有互斥锁；bodyMap 使用 UseNumber。
- 测试：两个集成测试默认 skip，启用后断言有效；两个 excel 小测试可失败且断言正确。
- multipart 体积上限取决于 platform/body.go:31 的 32MB MaxBytesReader 是否覆盖导入路由（lead 范围，未核实）。