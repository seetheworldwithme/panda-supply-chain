徐先生，commonrest 外围评审完成，报告如下。

## 范围
- 全文读完 34 个文件：范围内 24 个（15 个源码 + `register.go` + 8 个测试），上下文 10 个（platform 的 `resp/dbutil/maps/allinpay/wechat_http/wxconfig`、auth 的 `middleware/compat`、`缺失接口与兼容差异.md`、`auth-login-current.json`）。部分读：`commonrest.go:28-202`、`template.go:440-514`、`excel.go:34-55/252-276`、`handler.go:236-295`。
- 未读（归模板引擎评审）：`template.go`/`excel.go` 其余、`bespoke.go`、`meta_gen.go`、`contract_import.go` 及 excel/routes/template 测试。`~/.claude/skills` 指南按协调者要求跳过。
- 下文“可调用”均指 `authGuard` 语义：持任意 JWT，或带 Signature 的匿名请求（`frontend/src/api/http.js` 含该签名逻辑，密钥是否内置未逐行核实）。

## 发现

### [🔴] downById 路径穿越任意文件读 (parity:unknown)
- 位置：backend/internal/commonrest/external.go:220-243
- 证据：
```go
name := r.URL.Query().Get("fileName")
p := os.Getenv("PANDA_TL_SFTP_DIR") + name
b, err := os.ReadFile(p)
```
- 场景：`GET /saas-common/upload/downById?fileName=../../../etc/passwd`，`os.Stat` 命中即跳过查库直接读出回写；环境变量未配置时 `p == name`，相对 cwd 同样成立。
- 修法：Stat 前加 `if name != filepath.Base(name) { platform.FailMsg(w, "未找到对应文件"); return }`。
- 把握：high；中间件链不清洗 query。

### [🔴] bigUpload 路径穿越任意文件写 (parity:unknown)
- 位置：external.go:575-634
- 证据：
```go
name := r.FormValue("fileName")
part := filepath.Join(dir, name+".part"+strconv.Itoa(chunk))
f, e := os.OpenFile(filepath.Join(dir, name), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
```
- 场景：`fileName=../../../tmp/x`（或进程用户可写的 `authorized_keys`），`totalChunks=1` → 分片落在 dir 外，随后以 APPEND 写入目标路径。
- 修法：`name = filepath.Base(name)`，含 `..` 或分隔符直接拒绝。
- 把握：high。

### [🔴] delByRelId 空 body 删除全部附件 (parity:unknown)
- 位置：external_commonfile.go:189-219（同类读：external.go:366-413、731-765）
- 证据：
```go
cond := "1=1"
rows, e := queryRows(s.pick(r), "SELECT id FROM chain_common_file WHERE "+cond, args...)
for _, row := range rows { s.removeCommonFile(r, row["id"]) }
```
- 场景：body `{}` 时 `validateBody` 对空 map 返回 true（template.go:445）→ 全表逐行删除，并清空企业 `advertising_video_id/website_img_id`、删磁盘预览文件。`upload/preFIle`、`getImgByRelId` 同样接受 `{}` → 全表读盘 base64，内存打满。
- 修法：`len(m) == 0` 直接返回空结果。
- 把握：high；Java 是否启用 MP BlockAttack 未知。

### [🔴] loginChatApp 仅凭 phone 即签发 token (parity:false)
- 位置：chat_weixin.go:490-552
- 证据：
```go
rows, err := queryRows(d, "SELECT * FROM chain_auth_user WHERE phone=? AND channel=2", m["phone"])
row = rows[0]
info, err := s.chatUserInfo(r, row)
```
- 场景：`POST /saas-auth/chainAuthChat/loginChatApp {"phone":"<他人手机号>"}`，无微信凭据、无短信码，直接得 token；未注册手机号则无验证建号（预注册劫持）。
- parity 说明：Java 源码意图如此，但运行 JAR 中 `hasPhoneLogin` 入口即抛错恒 500（文档 §3.2-1，Go 验收只覆盖“空体 200 分支”），Go 把线上从未可用的无验证登录激活了。
- 修法：要求 `encryptedData+iv`，复用 chatLoginChat:382 的 `wxDecryptData` 后再查用户。
- 把握：high。

### [🔴] loginByAccount/buslogin 无密码即签发任意账号 token (parity:true)
- 位置：chat_local.go:62-84
- 证据：
```go
query, args := "SELECT * FROM chain_auth_user WHERE account_type!=4 AND handle_result=1 AND (phone=? OR account=?)", []any{m["account"], m["account"]}
```
bcrypt 比对只出现在 `loginCarSystem` 分支（:71），其余 mode 查到即 `chatUserInfo` 发 token（:79）。
- 场景：任意登录用户 `{"account":"<他人手机号>"}` → 得到对方（含 level 2）token。
- 修法：若仅供内部 Signature 调用，则 `UserID != 0` 时拒绝；否则要求密码。
- 把握：high；`auth-login-current.json` 该路径 `success_required:true` 且差异为空。

### [🔴] updByOpenid/updByOpenId/reg 全字段批量赋值可自提权 (parity:true)
- 位置：chat_local.go:150-163、328-345；chat_weixin.go:771-776
- 证据：
```go
vals := insertVals(pick(m, e.dbcols))   // dbcols 含 level/account_type/company_id/handle_result/password
where := "openid IS NULL"                // m["openid"] == nil 时
```
- 场景：自注册用户 `updByOpenid {"openid":"<自己的openid>","level":2,"company_id":<任意企业>}` 后重新登录即持该企业 token；`reg` 可直接建 `level=2,handle_result=1` 账号。`insertVals` 只做类型转换不过滤列。文档 §3.2-6 写“缺 openid 全表更新 Go 阻止”，实际仍批量更新全部 `openid IS NULL AND channel=2` 用户。
- 修法：三处改列白名单（nickname/wxFaceImg/openid/officialOpenid/appClientId/isAgreeCall），缺 openid 直接失败。
- 把握：high。

### [🔴] 企业注册用 body 的 userId 覆写任意用户 (parity:unknown)
- 位置：enterprise.go:139-150（另 ：70、:121）
- 证据：
```go
vals := map[string]any{"company_id": cid, "account_type": "3", "account": m["account"], "handle_result": "0", "level": 1, "username": m["adminName"], "password": string(hashed)}
platform.UpdateByID(d, "chain_auth_user", m["userId"], vals)
```
- 场景：`chainBasicEnterprise/add {"account":"新","enterpriseName":"新企业","password":"x","adminName":"x","userId":<受害者id>}` → 受害者 account/password/company_id 被改写；`filelist[].id` 可改写任意 chain_common_file 行；新企业强制 `is_authentication=true`。
- 修法：`userId` 只取 `currentUser(r).UserID`；filelist 校验 `owner_id`。
- 把握：high（Go 行为），Java 同构未取证。

### [🟡] loginChatPhone 保留 15012345678 验证码旁路，与既定决策相悖 (parity:false)
- 位置：chat_local.go:183
- 证据：`if message := s.Auth.PhoneCodeError(phone, code); message != "" && !(chat && phone == "15012345678") {`
- 场景：`{"phone":"15012345678","code":"x"}` 免码得该账号 token；账号不存在则首个调用者建号。文档 §3.2-4 明确“Go 要求有效验证码，不复刻接管漏洞”。
- 修法：删掉 `&&` 之后的条件。
- 把握：high。

### [🟡] updByOpenId 可用任意 id 抽干全局注册积分池 (parity:true)
- 位置：chat_local.go:312-325 → chat_weixin.go:254-287
- 证据：
```go
if err := tx.QueryRow("SELECT count(*) FROM bus_score_log WHERE user_id=? AND score_get_type=1", m["id"]).Scan(&count)
if count == 0 { ... registScoreTx(tx, m["id"], m["phone"], false) }
"UPDATE bus_score_rule SET mo_get_sum_score=? WHERE id=1", pool-award
```
- 场景：遍历不存在的 `id`（phone 传数字串）→ 每次扣 `mo_get_sum_score` 并写垃圾 `bus_user_score/bus_score_log`，池耗尽后新用户无注册积分。
- 修法：发分前校验 `m["id"] == currentUser(r).UserID` 且用户存在。
- 把握：high。

### [🟡] 附件路径字段可被改写 + 按扩展名内联回显 (parity:unknown)
- 位置：external_commonfile.go:108-120、242-251；external.go:263、292-296、469、492
- 证据：
```go
d.Exec("UPDATE chain_common_file SET "+strings.Join(sets, ",")+" WHERE file_name=?", args...)  // vals 含 file_path，无归属校验
typ := mime.TypeByExtension(filepath.Ext(name))
```
- 场景：① `chainCommonFile/updateById {"id":<他人文件id>,"filePath":"/etc","fileName":"passwd"}` → `previewStream?fileName=<id>` 读任意文件；② 上传不限扩展名，传 `x.html`/`x.svg` 后以 `text/html` 内联，152 上 `/saas/` 与前端同域 → 存储型 XSS 盗 JWT。
- 修法：写入时拒绝 `..`/绝对路径，读取时校验路径在 `externalRoot()` 内；非图片/PDF 强制 `attachment` + `nosniff`。
- 把握：medium-high；否定证据：nginx 对 `/saas-common/` 另加了下载头。

### [🟡] 通联 confirm/sync 无事务无行锁，先改本地再调通道 (parity:unknown)
- 位置：tlzf_schedule.go:224-236；tlzf_confirm.go:76-103、155-175
- 证据：
```go
d.Exec("UPDATE chain_contract SET contract_status=?,sum_freeze_amt=?,...")
platform.Insert(d, "chain_tlzf_contract_log", ...)
platform.UpdateByID(d, "chain_tlzf_order", row["id"], ...)
```
- 场景：(a) 合同 UPDATE 与订单 UPDATE 之间崩溃 → 重启重同步，冻结/解冻重复累加；(b) confirm 的 `sum` 校验与逐单读改写无锁无事务（`tlSyncMu` 只管 sync），并发两次请求各自通过校验 → 双倍 2090；(c) :168 先扣原单 `rel_operate_amt` 落库，:175 才调通联，通联拒绝时本地额度已少而无解冻单；(d) `amtJf<0` 时 `remainder.Cmp(need) < 0` 恒假，`rel_operate_amt` 反被放大。
- 修法：两段包进 `platform.WithTx`（通联调用前置或事务内 `SELECT ... FOR UPDATE` 原单）；`amtJf/amtTran <= 0` 直接 602。
- 把握：medium-high。

### [🟡] startTLTask goroutine 只增不减，start/stop 入口无角色校验 (parity:unknown)
- 位置：tlzf_schedule.go:17-40；chat_weixin.go:40-45
- 证据：
```go
s.tlTasks[tlTaskName] = cancel // 保留 Java 重复 start 覆盖缓存而未取消旧任务的行为。
go func() { for { ...syncTLOrders(r, s.DB.Zjjg, false)...
```
- 场景：每次 payMoney/confirm 起一个新 goroutine 并覆盖 cancel，旧 goroutine 永不取消，每 5s 查一次 zjjg 库，进程生命周期内单调累积；`/saas-common/start`、`/stopAll` 任意登录用户可调，stopAll 令待同步订单停滞。
- 修法：覆盖前 `if old := s.tlTasks[tlTaskName]; old != nil { old() }`；三入口删除或限管理员（文档 §3.3 称无前端调用证据）。
- 把握：high（Go 行为），注释声称照 Java 但未取证。

### [🟡] 登录事务跨内部调用与微信 HTTPS，可拖垮连接池 (parity:false)
- 位置：chat_local.go:197、224-242、280；chat_weixin.go:336、439、455、520
- 证据：
```go
tx, err := s.pick(r).BeginTx(r.Context(), nil)
has, err := s.callJSON(r, "/saas-warehouse/busBusinessCard/searchHasDefault", card)
care, err = s.chatOfficialCare(r, str(row, "official_openid"))
```
- 场景：20 个并发登录各持 1 条事务连接、再经 BatchMux 申请第 2 条（池上限 20）→ 互锁，全服务 DB 请求挂起直到客户端断开。验证码非一次性（auth/handler.go:267 只比对不删除），用自己的一个验证码并发重放即可触发。回滚时名片/Redis 已写，留孤儿数据。
- 修法：内部调用与微信调用移到事务外（Java 本就非事务），事务只包用户行写入。
- 把握：medium；若 BatchMux 内层使用独立连接池则不成立（未见此配置）。

### [🟡] getAppInfo 回显 AppSecret、getAccessToken 回显 access_token (parity:true)
- 位置：chat_local.go:24-40；chat_provider.go:19-27
- 证据：`platform.Success(w, authRestEnts["ChainAuthChatAuthorization"].shapeOne(rows[0]))`，`keys` 含 `secret`（auth_meta_gen.go:288），`shapeOne` 无脱敏。
- 场景：任意调用者拿到小程序/公众号 AppSecret。通用 CRUD 已对该实体 `continue`（auth_rest.go:21），本意是不外露。
- 修法：输出前 `delete(out, "secret")`；`getAccessToken` 限管理员。
- 把握：high；两侧一致 → Java 同样回显。

### [🟡] chainTlzfOrder 通用 CRUD 可伪造/改写支付订单 (parity:true)
- 位置：register.go:34；tlzf.go:177-186；tlzf_confirm.go:76、103
- 证据：
```go
row := pick(m, entChainTlzfOrder.dbcols)   // payMoney/confirm 的 body 可带 id/trans_status
```
- 场景：`chainTlzfOrder/add {"contract_id":X,"trans_type":1,"trans_status":1,"rel_operate_amt":99999999}` → confirm 的 `sum(rel_operate_amt)…trans_status=1` 被抬高；updateById/updBatch 可改他人订单。
- 修法：该实体不注册 add/updateById/updBatch；`tlPay` 强制 `trans_status=0` 并丢弃 `id`。
- 把握：high。

### [🟢] serveLocalDownload speed 溢出导致除零 panic (parity:false)
- 位置：external.go:170
- 证据：`wait := time.Duration(end)*time.Second/time.Duration(speed*1024) - time.Since(start)`
- 场景：`speed=18014398509481984`（2^54）时 `speed*1024` 回绕为 0 → 整数除零 panic（仅断开该连接）。前提 `D:/xmmc.apk` 存在，线上基本不可达。
- 修法：`int64(speed)*1024`，并限制上限。

## 重复模式
- **按 id 读写不校验归属**（Java getById/allEq 同，parity:true）：tlzf.go:65-70、99、247；external.go:189-218、658；external_commonfile.go:80、122、176、250；enterprise.go:121。建议在模板层统一加 `company_id/owner_id` 过滤。
- **io.ReadAll 无上限**：external.go:73、external_commonfile.go:16、:222，不在 body_contracts（body.go:31 的 32MB 上限不生效）；换 `http.MaxBytesReader` 即可。
- 顺带提示（非本范围）：`/saas-auth/chainAuthUser/list`（auth_rest.go:57）对任意登录用户返回全表用户行，含 `password` 哈希与手机号。

## 已核对无问题
- tlAsync：先 `VerifyAllinpay`（SM2/SM3）→ 解析 → `req_type=2` 日志查重 + 1h 内存锁 → 提交任务；`tlSyncStatus` 在 `tlSyncMu` 内重读订单、`trans_status==1` 即返回，回调重放幂等；未签名回调只记日志；`asyncTlzf` 匿名路径除验签回调外无写能力；日志不含私钥/密码。
- 金额全程 `big.Rat`，`FloatString(2)` 与 Java HALF_UP 一致；SQL 值全部参数化；动态列名经 `externalColumn` 白名单。
- qcc/ocr/esign/微信的 URL 与密钥全部取自 env/私有 JSON，无 SSRF；出站 HTTP 均有 10~30s 超时，通联响应 `LimitReader` 16MB。
- `queryRows` → `ScanMaps` 统一 Close；各事务均带 `defer Rollback`；`uploadFiles` 有 `defer RemoveAll`，`readExcel` 未显式清理但 net/http `finishRequest` 会代删临时文件。
- 短信码 6 位/60 秒/按手机号分键，暴力破解不可行；`contract_batch` 多批半写、`enterpriseAdd` 吞掉插入失败均有注释且属 Java 原行为，未报。
- 测试均为真实断言，无空测试。
- 提请复核（platform/allinpay.go，非本范围）：`VerifyAllinpay` 用 `tlCertPath`（同时是商户证书序列号来源）的公钥验通联签名，若线上回调由通联私钥签名则验签恒失败——功能问题而非安全问题。

## 🎉 值得保留
- `WithTx + defer Rollback` 统一防半写，通联金额全程 `big.Rat` 且测试用真实 openssl SM2 签名链验收。
- 外部适配器全部 env/私有 JSON 注入，源码零密钥。