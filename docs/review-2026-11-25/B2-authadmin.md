徐先生，authadmin 包评审报告如下。

## 范围
通读 `backend/internal/authadmin/` 全部 11 个文件（3061 行，含 3 个测试），并读上下文：`main.go`、`platform/{resp,dbutil,maps,body,uploadrender}.go`、`auth/{middleware,handler}.go`、`auth/compat.go:100-176`、`platform/tencent.go:55-109`、`commonrest/external.go:172-330`、`docs/validation/java-go-20261118/缺失接口与兼容差异.md`。无跳过。Java 源码本机不存在，parity 只能依据代码内“照 Java”注释和该文档判断。评审期间 `auth/handler.go` 行号在变（有人在并发改动），下列行号以我读取时为准。

## 发现

### [🔴] 用户域写接口没有任何权限边界，任意登录用户可接管平台管理员 (parity:true)
- 位置：`user_compat.go:263-287`（updPwd/resetPwd）、`:374-378`（joinCpy）、`:216-261`（add/updCarAccount）、`:317-340`（updBatch）、`:28`（updHeadImg）；`user.go:222-243`（add）、`:258-263`、`:377-378`、`:485-500`（三个 updateById）、`:386-403`（delete）、`:417-441`（logOffUser）、`:508-530`（assignRole）
- 证据（`user_compat.go:280-282`）：
  ```go
  vals := pick(m, userWriteCols...)
  vals["password"] = auth.HashPassword("123456") // original Java resetPwd contract
  if _, err := platform.UpdateByIDNoFill(s.DB.Chain, "chain_auth_user", m["id"], vals); err != nil {
  ```
- 场景：这些 handler 只按 body 的 id 写库，不查归属也不查调用者 account_type，且不检查 `UserID==0`（持 HMAC 的 SignedAnon 同样可达）。普通企业账号发 `resetPwd {"id":<管理员 id>}` 后用 123456 登录即管理员；或 `add {"account":"x","password":"x","accountType":5,"level":2,"handleResult":1}` 自造平台账号；`chainAuthUserList/updateById {"id":自己,"handleResult":1}` 自审批；`logOffUser {"id":X}` 删任意用户及其 `bus_product/bus_videos/bus_business_card/bus_visitor_log/bus_user_care`。管理员账号名可由 `/chainAuthUser/ptlist` 读到。
- 修法：需产品决策。若收口，加包内小函数：`UserID==0` 拒绝；调用者 account_type∉{2,5} 时，密码、重置、删除类接口强制 `m["id"]=cu.UserID`，并从 vals 剔除 `password/account_type/level/handle_result/company_id`。
- 把握：Go 行为 high。Java 同样无校验为 medium——各 handler 注释都写明是 MP `save/updateById/removeById(s)` 原样，且 docs 3.2 例外清单未列；若 Java 控制器从 JWT 取 id，应改 parity:false。

### [🔴] 任意文件读取：企业更新接口允许 body 改写任意文件行的物理路径 (parity:true 写入侧 / unknown 读取链)
- 位置：`common.go:445-454`（commonFileUpdate），经 `common.go:164-206` 触发；读取出口在 `commonrest/external.go:234,264,288`
- 证据（`common.go:445-454`）：
  ```go
  vals := pick(em, "file_real_name", "file_name", "file_description", "file_path", "f_type", ...)
  vals["relation_id"] = m["id"]
  _, err := platform.UpdateByIDNoFill(s.DB.Chain, "chain_common_file", em["id"], vals)
  ```
- 场景：任意登录用户先 `/saas-common/upload/batch` 得到文件行 F，再发 `/saas-common/chainBasicEnterprise/updateById {"id":"999999999","filelist":[{"id":F,"filePath":"/etc","fileName":"passwd"}]}`（id 取不存在的值可避开 removeFileByRel 误删），然后 `GET /saas-common/upload/previewStream?fileName=passwd` 得到 `os.ReadFile("/etc/passwd")`。指向 KeysDir 可读 JWT 私钥伪造任意用户，指向 .env 可得库口令。
- 修法：pick 去掉 `"file_path", "file_name"`。
- 把握：high（Go 两侧代码都已读到）。commonrest 的 `chainCommonFile/updateById` 若同样放行 filePath，则根因在通用 CRUD（他人范围），此条是其可达实例。

### [🔴] setAccount 明文密码入库，与 docs 决定相反，并经用户查询接口回显 (parity:true 对 Java / 与文档相反)
- 位置：`user_compat.go:110-113`；回显点 `user.go:77-78,125`、`user_compat.go:129-134,295-300`、`common.go:134-147`
- 证据（`user_compat.go:112-113`）：
  ```go
  // Java chainAuthUserService.save(user) 原样入库(明文密码),1:1 复刻,不做加密。
  id, err := platform.InsertNoFill(s.DB.Chain, "chain_auth_user", vals)
  ```
- 场景：`shapeUser` 保留 `password` 键，任意登录用户发 `/saas-auth/chainAuthUser/list` body `{}` 即 `SELECT * FROM chain_auth_user` 全表返回。登录只做 bcrypt 比较，这类账号用密码永远登不进 `/saas-auth/login`。docs 3.2 第 5 条写的是“Go 使用既有 BCrypt”，代码与文档必有一处过期。
- 修法：恢复 `vals["password"] = auth.HashPassword(str(m,"password"))`，同步 docs。
- 把握：high。

### [🟡] 平台配置类写接口同样无管理员边界，机密字段回显 (parity:true)
- 位置：`common.go:382-385`（entUpdateTry）、`:282-299`（urcs）、`:544-551`（glconfig）、`:677-704`（eqb）；`biz.go:165-204`；`role.go:128-218`
- 证据：`vals := pick(m, entCols...)` 后 `UpdateByIDNoFill(s.DB.Chain, "chain_basic_enterprise", m["id"], vals)`
- 场景：按 id 改任意企业 90 多列：`is_authentication=1` 自认证、自送 `sign_free_count/service_expiration_time`；`urcs:[]` 把该公司全部用户 company_id 置空；全局支付、会员开关和 `red_uplimit_days` 任改；eqb 的 `redirecturl` 可改，`filepassword`、通联 `acct_num` 经 list/getById 明文返回。
- 修法：非平台账号拒绝这批路由；需产品决策。
- 把握：Go 行为 high；Java parity medium。

### [🟡] uploadBatch 无大小、类型、并发上限 (parity:unknown)
- 位置：`common.go:719-806`
- 证据：`r.ParseMultipartForm(32 << 20)`；`item["url"] = "data:image/jpg;base64," + base64.StdEncoding.EncodeToString(b)`（799-800）
- 场景：全仓无 `MaxBytesReader`，超过 32MB 的部分落盘但总量不限。图片被 `CompressFileRow` 和 799 行各整读一次再 base64，1GB 的“图片”约占 4~5GB 内存，单进程全站 OOM。xls、doc 每个文件起一次 soffice（各 60s 超时）且无并发闸门，docs 仍写“Go 复刻中断语义”，与代码不符。扩展名不限，.html、.svg 原样落盘。循环中途 break 时已落盘文件成孤儿（749-759）。
- 修法：`r.Body = http.MaxBytesReader(w, r.Body, N)`，加扩展名白名单。
- 把握：high。

### [🟡] 任意文件删除链：removeFileRow 按库内路径 os.Remove (parity:unknown)
- 位置：`common.go:429-432`
- 证据：`_ = os.Remove(filepath.Join(str(image, "file_path"), str(image, "file_name")))`
- 场景：commonrest 注册了 `/saas-common/chainCommonImg/add|updateById`（`meta_gen.go:106` 的 dbcols 含 file_path、file_name）。插入一条 img 行（file_id=F，路径任意），再让文件行 F 满足 `relation_id=X, file_type='3'`（同一个 entUpdateById 的 filelist 分支即可），第二次发 `{"id":X,"filelist":[]}` 即触发。
- 修法：Join 后校验以上传根为前缀，否则不删。
- 把握：medium，若通用 add 实际不落 file_path 则不可达。

### [🟢] HashPassword 失败返回 "" 被落库；另有 5 条路径把 body.password 原样写入 (parity:false / true)
- 位置：哈希后未判空 `user.go:235`、`user_compat.go:227,271,377`；原样写 `user.go:258,377,485`、`user_compat.go:240,326`
- 证据：`vals["password"] = auth.HashPassword(password)`，而 `HashPassword` 出错返回 `""`（x/crypto v0.59.0，>72 字节报 ErrPasswordTooLong）
- 场景：`add {"account":"a","password":"<73 字节>"}` 插入 password="" 的账号并返回成功，永久无法登录；`updPwd` 同理可清空受害者密码；`add` 不带 password 则 bcrypt("") 入库。
- 修法：四处调用后 `h == ""` 即 `platform.Unknown(w)`。
- 把握：high。Java BCrypt 抛异常走 500，故 parity:false。

### [🟢] 上传目录取自用户可自改的 phone (parity:unknown)
- 位置：`common.go:721-723`
- 证据：`apath = str(user, "phone")` → `dir := filepath.Join(uploadDir(), filepath.FromSlash(apath))`
- 场景：phone 在 `userWriteCols` 内且无校验。非企业账号先把 phone 改成 `../../tmp/x` 再上传，`MkdirAll` 和写文件落到上传根之外（文件名随机，扩展名受控）。
- 修法：`apath = filepath.Base(apath)`。
- 把握：medium（需 EnterpriseCode 为空）。

### [🟢] 事务内再取池连接并调外部 HTTP；userAuthEtc 三表写无事务 (parity:true / unknown)
- 位置：`user.go:262-317`；`user.go:343,370,378`
- 证据：WithTx 回调内 `cards, err := queryRows(s.DB.Chain, "SELECT * FROM bus_business_card WHERE msg_type=2")` 与 `platform.SendLegacySMS(...)`
- 场景：池上限 20，`queryRows` 用无 ctx 的 `d.Query`。20 个并发请求都持事务连接等第二个连接时互等至死锁。现网缺 msg_type 列使该查询秒失败，窗口极窄；但会员到期分支在持行锁期间同步调短信网关（30s 超时）。userAuthEtc 后两步失败会留下半写。
- 修法：cards 查询移到 WithTx 之前，短信放到提交之后；userAuthEtc 收进 WithTx。
- 把握：high。

### [🟢] 邀请码内存 map 只增不删 (parity:false)
- 位置：`user.go:444-479`
- 证据：`inviteCodes[key] = inviteCode{value: code, expire: time.Now().Add(time.Duration(secs) * time.Second)}`
- 场景：登录用户循环调 createCode，每次新增永不回收的条目（Java 用 Redis TTL）。effSeconds 很大时 Duration 溢出为负，码生成即过期。
- 修法：写入时顺手删过期项。
- 把握：high。

## 重复模式
- 按 id 直写、不校验归属与调用者身份：`user_compat.go:28,139,160,227,240,272,282,328,378`；`user.go:225,263,343,378,388,398,409,420,500,512`；`common.go:265,294,385`；`role.go:116,131,241,252,263`；`menu.go:83,99,112,131,180,215,226`；`biz.go:85,101,111,178`。Java 同款。
- 已注释声明刻意保留、仅知悉不计入发现：`biz.go:54,56` 是包内唯一拼接请求值的 SQL（`" and partner_type = " + str(m, "partnerType")`、`like '%"+searchStrOf(m)+"%'`，可注入，parity:true）；`user.go:194` 的 `phone==''`；`common.go:228-280` 循环内 clear+upd。

## 已核对无问题
- **SQL 键**：所有 Insert、Update 的键都来自 `pick`、`allEqCond` 的固定白名单，CheckCols panic 不可达；IN 列表全用占位符；ORDER BY 为常量。
- **panic**：body 为 nil map 时读取安全；`expirationTime[:10]`（DATETIME 经 parseTime 恒 19 字符）、`phone[7:]`、`parts[1]` 均有长度守卫；各 WithTx 回调内无可达 panic；IPage 的 size≥1。
- **资源**：rows 均经 ScanMaps 关闭；multipart 临时文件由 net/http 清理；两个渲染器的临时目录 defer RemoveAll；图片解码先读头部，超过 5000 万像素不解码。
- **并发**：`inviteCodes`、`auth.Service.sms` 均有 mutex。
- **测试**：三个测试都能失败且断言正确，后两个需 `PANDA_LOCAL_AUDIT=1`。
- **死代码**：`warehouseBrokerList`（`user.go:549`）全仓无注册点，它是包内唯一带 account_type 判断的 handler。
- **值得保留**：列白名单把 MP allEq 的“键即列名”彻底关死；Java 无事务的多表写（delRoleById、assignRole、saveEtcAuths、menuDelBatch、userLogOffUser）统一收进了 WithTx。