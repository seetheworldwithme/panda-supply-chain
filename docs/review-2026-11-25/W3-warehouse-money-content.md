徐先生，以下是本评审员（warehouse 定制资金/会员/内容处理器）的最终报告。全程只读，未运行构建/测试。

## 范围
- 全文读完 13 个源码文件 + 7 个测试文件，共 20 个：noticeorder.go、paynotify.go、pay_v3.go、red_transfer.go、member2.go、leasing_write.go、remaining.go、videos.go、msgdomain.go、prolog.go、content_create.go、content_external.go、hanlp_native.go 及对应测试。
- 为核对行为另读了 platform/{resp,dbutil,maps,body,wechatpay,wechat_cert,wechat_http,redis,tencent,wxconfig}.go、auth/{middleware,member_expiry}.go、warehouse/{warehouse,register,template}.go。
- 未复核：template.go、excel.go、product_* 的内部实现（他人范围）；本机没有 Java 源码，parity 标签依据代码注释里的“实测/原样保留”。

## 发现（按严重度降序）

### [🔴] 红包出款的金额、收款人、流水归属都取自请求，任意登录用户可给自己出款 (parity:true)
- 位置：red_transfer.go:79-103、160-161；content_create.go:187-210、262-279；member2.go:207-294
- 证据：
  ```go
  amount := m["receiveAmount"] // red_transfer.go:160，下一行直接作 "total_amount": amount
  uid := m["userId"]           // :103，转账到 user["openid"]
  withdrawal = map[string]any{"clogId": clogID, "receiveType": opinionReward, "receiveAmount": amount, "userId": body["userId"], ...} // content_create.go:270
  ```
- 场景：
  1. `/busRedEnvelopeRecord/generateAmount {receiveType:X}` 自建流水，再调 `/receiveAmount {clogId, receiveType:X, userId:自己, receiveAmount:任意分}`，即触发微信商家转账。只有类型 1/6/11/12 有次数上限，其余类型仅受“当日已出款总额”约束，且该上限先查后转、单笔可超额。
  2. 自建意见后调 `/busOpinion/updateById {id, status:1, userId:自己}`，逐条出款。
  3. `/busVideos/updateById {id, status:1}` 任意用户可审任意视频。奖励次数按“调用者 uid”查（:200），记录却写视频主 uid，两个小号互刷不受 3 次/日限制。
- 注释已写“原 Java 不核对客户端金额…保留该顺序与问题”，但 docs 3.2 例外只登记了 payNotify，未登记此项。
- 修法：`redReceive` 以 `clog.create_amount` 为金额，要求 `clog.user_id == currentUser(r).UserID`（内部调用改传显式受信参数）；审核端点加角色判断。
- 把握：high。能否定：`bus_red_envelope_range` 仅存在有上限的类型（无法查库确认）。

### [🔴] 请求值直拼 SQL，照 Java `${}` 注入 (parity:true)
- 位置：noticeorder.go:69、:221、:564-567、:1014；msgdomain.go:696、:703；prolog.go:400；remaining.go:193、:209
- 证据：
  ```go
  q := "select * from bus_notice " + ss + " and CONCAT(id) like ? ..."        // noticeorder.go:69，ss = body.searchSql
  cs = append(cs, k+" = '"+v+"'")                                             // :221
  q := "select *, ... from bus_visitor_log " + sqlStr + " and CONCAT(id) like '%" + inner + "%' ..." // msgdomain.go:702-703
  "SELECT code FROM chain_area_info WHERE name LIKE '"+city+"'"               // remaining.go:193，Excel 单元格
  ```
- 场景：登录用户 `POST /busVisitorLog/list {"searchStr":"' OR 1=1 -- "}`，或导入一行 city 为 `x' UNION SELECT ...` 的 Excel。可 UNION 或时间盲注读全库（DSN 为 admin 账号）。
- 修法：searchStr、userPhone、city/area 改 `?` 参数；`searchSql` 整段用白名单正则校验，否则 999。template.go:571 的 searchLit（busOrder、busVisitorLog 模板）同病，归模板负责人。
- 把握：high。除 remaining.go 外，注释均自述“字面注入”。

### [🔴] 公网回调无上限读 body (parity:unknown)
- 位置：paynotify.go:81、:103；同类（需登录）videos.go:56、msgdomain.go:44
- 证据：`body, err := io.ReadAll(r.Body)`（:81）
- 场景：匿名向 `/saas-warehouse/busMyBill/payNotify` POST 数 GB body，内存被读满。payNotify 不在 body_contracts_gen.go 内，body.go 也不兜底；后端机仅 8G。
- 修法：`io.LimitReader(r.Body, 1<<20)`。
- 把握：high。能否定：前置网关有 body 上限（未核实）。

### [🟡] MchSecret 缺失时密钥退化为字面 "null"，通知可伪造 (parity:false)
- 位置：paynotify.go:91-95
- 证据：
  ```go
  if key == "" { key = "null" }
  if !wxVerifyNotify(params, key) || params["return_code"] != "SUCCESS" {
  ```
- 场景：`WECHAT_LOCAL_CONFIG` 漏配或 JSON 缺 mchSecret，攻击者用 key=null 算 MD5，伪造自己待支付账单的已支付通知，免单开通会员。Java 密钥硬编码非空，不会落入此分支。
- 修法：`key == ""` 时直接回 200 空体。
- 把握：high。

### [🟡] 空字符串日期触发 panic (parity:false)
- 位置：prolog.go:139-141、:144-145、:243-248
- 证据：
  ```go
  st, _ := bindVal("D", v)   // v=="" 时返回 (nil, true)，validateBody 已放行
  proWhere += " and DATE_FORMAT(create_time, '%Y-%m-%d')>= " + st.(string)
  ```
- 场景：登录用户 `POST /busProLog/selBuserTeamData {"startTime":""}`（selVisisted 同）→ 对 nil 断言 panic。net/http 会 recover，但连接被断、堆栈日志刷屏。Java 对空串得 null 后跳过条件返回 200。
- 修法：`if st, ok := bindVal("D", v); !ok || st == nil { continue }`。
- 把握：high。能否定：Java 对空串也回 999（无源码未实测）。

### [🟡] 视频删除链先删文件行再查询，物理文件清理恒不执行 (parity:unknown)
- 位置：noticeorder.go:833→845；videos.go:747→759
- 证据：
  ```go
  d.Exec("DELETE FROM chain_common_file WHERE id=?", vid)      // noticeorder.go:833
  frow, err := platform.GetByID(d, "chain_common_file", vid)    // :845，恒为 nil
  ```
- 场景：任意一次视频删除后，视频文件与 `_img` 目录永久占盘，chain_common_img 只删文件不删行。
- 修法：`GetByID` 提到 DELETE 之前。
- 把握：Go 侧死码为 high；Java removeFile 是否同序未能核实。

### [🟡] miniAppPay 无超时，且在打开的事务内外呼 (parity:false)
- 位置：paynotify.go:326；content_external.go:188-221
- 证据：`resp, err := http.DefaultClient.Do(req)`（:326，Timeout 为 0）；tx 在 :188 Begin，:221 才 Commit
- 场景：微信下单接口挂起时，每个请求占住事务连接，再加事务内 `GetByID` 占的第二个，直到客户端断开。十余个并发即可耗尽 20 连接池，全站不可用。v3 路径同理，最长 30s×2。
- 修法：改用 `&http.Client{Timeout: 15 * time.Second}`。
- 把握：high。

### [🟡] payNotify 在事务内同步发短信，叠加无幂等 (parity:unknown)
- 位置：paynotify.go:191-196、:250；tencent.go:67
- 证据：`platform.WithTx(d, func(tx *sql.Tx) error { ... return s.billUpdUserInfo(r, tx, bill) })`，其内 `SendLegacySMS(...)`（超时 30s）
- 场景：短信通道慢，事务与行锁挂住最长 30s，微信等待超时后重试。账单已是 status=2 仍再次续期。例外清单 3.2#9 只登记“无金额/订单/幂等校验”，未登记这一放大因素。
- 修法：短信移到 Commit 之后；`UPDATE bus_my_bill SET status=2 WHERE id=? AND status<>2`，按 RowsAffected 决定是否做副作用。
- 把握：medium（取决于微信重试时机）。

### [🟡] searchByUser 向任意登录用户返回首行用户全字段 (parity:true)
- 位置：msgdomain.go:385、:394
- 证据：`"SELECT * FROM chain_auth_user ORDER BY id LIMIT 1"`，经 entChainAuthUserLive 整形，keys 含 password、openid、phone
- 场景：`{}` 请求即得 id 最小账号（通常是管理员）的密码哈希与手机号。
- 修法：不回传 password；userInfo 按 body userId 取。
- 把握：high（注释称 8765 实测 69 键，Java 同样）。

### [🟡] custNoticeCust 可向任意号码群发短信 (parity:true)
- 位置：msgdomain.go:485、:499-519
- 证据：`values, ok := ph.([]any)`，随后 `SendLegacySMSMany(r.Context(), phones, ...)`
- 场景：proId 不校验归属，phones 取自 body 且无条数上限。任意登录用户可让全站管理员/任意号码收到带货主名称与价格的短信，费用与骚扰无上限。
- 修法：proId 校验 create_user_id，phones 限制为该产品的客户表且设上限。
- 把握：high。

### [🟢] 签到两处小错 (parity:unknown)
- 位置：member2.go:315、:374-394
- 证据：`signMap[st[:10]] = str(sr, "hjjNum")`（原始行键是 hjj_num，恒 ""，已签日只显示 "+"）；`INSERT INTO bus_user_score (user_id, hjj_num) VALUES (?,1)` 写死 1，丢连续加成
- 场景：先 count 后 insert 无唯一约束或条件更新，并发双倍积分。
- 修法：改 `str(sr, "hjj_num")`；INSERT 用 addScore；签到行改唯一键加 `INSERT IGNORE`。
- 把握：显示问题 high；并发问题 medium（未查到唯一索引）。

## 重复模式
- **按 id 读改删不校验归属**（parity:true，包注释“租户拦截器失效”）：
  - videos.go:686 updByUserId 以 body userId 批量改写他人视频全列；:161 updateVideo 任意视频全列
  - noticeorder.go:369-373 删任意广告及关联 notice；:829、videos.go:740 删任意视频
  - leasing_write.go:51、:130 任意租售记录与任意 chain_common_file 行（含 file_path/file_name）
  - content_create.go:119 把任意文件 owner 改为自己；content_external.go:141 任意 userCouponId 抵扣，支付后该券被删
  - msgdomain.go:228 改任意 bus_msg；:207-209 selNoReadMsg 无 `um.user_id` 条件，可读全员消息（parity:unknown）；:249 oncClickRead 把所有用户该类型消息置已读
  - noticeorder.go:549 `isPtSearch:true` 看全量订单；prolog.go:483 body userId 等于货主即加分，可刷积分
- **计数读改写丢更新**（parity:true）：videos.go:141-157 整行回写；prolog.go:476-480（且 `_, _ =` 吞错）；msgdomain.go:494-495；red_transfer.go:62-72 日上限先查后转。
- **跨范围提示**（负责人文件）：public 的 appPayNotify 带任意 `wechatpay-serial` 即在全局锁内下载平台证书（wechat_cert.go:64-68），匿名可放大外呼并串行化全部 v3 调用。

## 已核对无问题
- 验签先于一切状态变更：v2 MD5（空值与 sign 不参与）；v3 RSA PKCS1v15 加 AES-GCM，校验 nonce 与 APIv3Key 长度。验签失败回 200 空体，落库失败回 500 让微信重试，账单与会员同事务。XML 用 encoding/xml 流式解析，无实体展开。
- 金额全程 big.Rat 或字符串十进制，无浮点运算；family C、rrTotoalMoney、elListByPage 的日期拼接经 `bindVal("D")` 校验，只能是日期；vidLit、strOrNull、javaNumStr、cardUidOf 拼接处均为已校验整数。
- syncBiddingInfo 出站 URL 取自环境变量，16MB LimitReader、25s 超时、defer Rollback；infoVR 固定 URL，响应上限 8MB，无 SSRF。
- hanlp_native：embed 模型经 sync.OnceValue 只加载一次，加载后只读、并发安全；下标访问均有界，分词路径恒含首尾节点；golden 测试为真实比对。
- vLcToggle、addVideo、updateOpinion、placeOrder 的事务各失败路径均 Rollback，范围内无 panic 落在事务内；placeOrder 无 defer 但逐路径回滚。
- pageOf 保证 size≥1，IPage 无除零；WithTx 内无 CheckCols 触达请求键。
- 测试无空洞：business_notify、red_transfer、leasing_write 均断言真实副作用。未核实项：hanlpNative 入参 `ss` 的长度上限（product_list.go，他人范围）。

## 🎉 值得保留
- 🎉 updateOpinion（content_create.go:239）：`SELECT status ... FOR UPDATE` 加事务内旧状态判定，是范围内唯一正确防并发双发的写路径。
- 🎉 支付金额与回调验签的纯标准库实现：big.Rat 精确十进制，v3 解密验签先于落库，失败交微信重试。