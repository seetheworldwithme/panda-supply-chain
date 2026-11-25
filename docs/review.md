# 全仓代码评审（前端 + Go 后端）

> 日期：2026-11-25　基线：main `215e420`　方法：`code-review-skill`（Go / Vue / 安全 / 通用质量 / 常见缺陷清单）
> 原始分组报告（含逐条证据行）见 `docs/review-2026-11-25/`。

## 0. 结论

| 项 | 结果 |
| --- | --- |
| 范围 | 后端 121 个 Go 源文件 + 52 个测试文件全部逐行评审；前端应用外壳（api / router / store / 顶层页）、`setup` + `Emain`、购销合同 4 个页面逐行评审；其余页面只做了「接口调用 ↔ Go 路由」机械比对（见 §5） |
| 修复 | **后端 53 组、前端 15 组已修复并验证**（§2、§3，每组含 1 到数十处同类站点），全部是 Go 迁移自身引入的问题，或**不改变合法业务行为**的加固 |
| 待决策 | **21 项原 Java 遗留的设计级安全问题**（§4）。按 `docs/plans/2026-11-19-java-go-business-parity.md` 第 4 条「保留 Java 问题，不擅自修复」，本次**没有动**，但其中多条可直接导致账号接管、任意出款、清库 |
| 未修 | 7 个端点的 Java 定制写入逻辑在 Go 里缺失（§6.1，切换后端前必须补）；其余前端页面问题已定位未修（§6.2） |
| 验证 | `go vet ./...` 干净；`go test ./...` 与 `PANDA_LOCAL_AUDIT=1 go test ./...`（本机 MySQL，库从后端机备份导入）全绿；前端 20 个改动文件通过 `node --check` |

严重度：🔴 阻断（可被利用 / 资金或数据损失 / 常见路径崩溃）　🟡 重要　🟢 轻微。

## 1. 先看这里：需要您拍板的 5 件事

下面都是 **Java 原系统就有、Go 版照原样保留** 的问题，不是迁移引入的。修任何一条都会偏离 Java 行为，所以只列不改。

1. **任何登录用户可以把自己变成平台管理员。** `POST /saas-auth/chainAuthUser/resetPwd {"id":<管理员id>}` 把任意账号密码重置为 `123456`；`add` 可直接建 `accountType=5` 的账号；`chainAuthUserList/updateById` 可自审批。`authadmin/user_compat.go:263-287、216-261`，`authadmin/user.go:222-243、485-500`。
2. **知道一个手机号就能登录。** `/saas-auth/authentication/mobile`（ELogin 页）不校验短信码；`/saas-auth/chainAuthChat/loginByAccount|buslogin` 传账号即发 token；`loginChatApp` 仅凭 phone；`loginChatPhone` 对固定手机号免码；前端 `Sponsor.vue:155-180` 甚至支持 `?account=xxx` 直接登录。位置：`auth/compat.go:79`，`commonrest/chat_local.go:62-84、183、490-552`。
3. **红包出款的金额和收款人取自请求体。** 任意登录用户自建一条流水后调 `/busRedEnvelopeRecord/receiveAmount` 即可向自己发起微信商家转账。`warehouse/red_transfer.go:79-103、160-161`，`content_create.go:187-210、262-279`。
4. **空条件清库。** `POST /saas-contract/chainContract/delOnce` 传 `{}` 会无 WHERE 删除合同、结算、收付款、票据、运单、磅单 7 张表——这一条**已修**（§2.1），但 `{"im_batch_number":null}` 仍会删光所有非导入数据，是原功能，请确认是否保留。
5. **前端内置的静态 `Signature/AuthorizationId` 等于全站匿名凭据。** 带上 `http.js` 里写死的两个请求头，就能匿名读平台账号列表、角色、菜单、e签宝配置，甚至 `chainAuthSecret/list` 读出全部 HMAC 密钥。`auth/middleware.go:68-83`。

建议：在做小程序 / 前端切换到 Go 之前，单独立一个「鉴权收口」批次处理 §4.1–4.4（登录族、管理员边界、出款、SignedAnon 白名单）。需要我先出方案再实施的话说一声。

## 2. 后端：已修复

### 2.1 🔴 安全（可被利用）

| # | 位置 | 问题 | 修法 |
| --- | --- | --- | --- |
| S1 | `contract/contract.go` `delOnce` | body 为 `{}` 时无条件 `DELETE` 7 张业务表 | 无条件直接拒绝 |
| S2 | `commonrest/external.go` `downById` / `bigUpload` | `fileName=../../etc/passwd` 任意文件读；`fileName=../../x` 把分片写到上传目录之外 | 只接受纯文件名 |
| S3 | `commonrest/template.go`、`external_commonfile.go`、`authadmin/common.go:444` | `chainCommonFile/add`、`updateById`、`updBthById`、企业更新都允许客户端改写 `file_path`，配合预览接口读任意文件（含 JWT 私钥、`.env`） | 新增时 `file_path` 必须在上传根内、`file_name` 必须是纯文件名；更新时忽略 `file_path`；读取时目录为空的行不再解析 |
| S4 | `commonrest/external_commonfile.go` `delByRelId`、`external.go` `preFIle` / `getImgByRelId` | 空条件 = 全表：逐行删光附件及磁盘文件 / 全表读盘转 base64 | 空条件直接返回 |
| S5 | `capital/payment.go:379` + `auth/middleware.go` | 招行回调只看请求头有没有 `CMB-BodySign`：带有效 JWT 不验签即可把订单改成已支付 | 必须是 RSA 验签通过的请求（`CurrentUser.BankSigned`） |
| S6 | `warehouse/paynotify.go:91` | 未配置商户密钥时密钥退化为字面 `"null"`，支付通知可伪造 | 无密钥直接按验签失败处理 |
| S7 | `platform/esign.go` | e签宝回调验签没有时间窗口，可重放 | 时间戳偏离超过 10 分钟拒绝（需服务器时钟同步） |
| S8 | `platform/wechat_cert.go` | 免登录的 `appPayNotify` 带随机 `wechatpay-serial` 即可无限触发平台证书下载，且全程占全局锁 | 未知序列号一分钟内只下载一次 |
| S9 | `platform/uploadrender.go` | 几 KB 的「解压炸弹」图片声明 60000×60000，解码即 OOM 整个进程 | 解码前读图片头，超过 5000 万像素不解码 |
| S10 | 5 份 `excel.go`（commonrest/bill/capital/logistics/warehouse） | `<c r="ZZZZZZZZ1"/>` 触发无上限补列、zip 炸弹无上限解压 | 列坐标 ≤16384 且与前一列间隔 ≤256；单部件 64MB、总量 256MB |
| S11 | `warehouse/product_bp.go:327`、`product.go:29` | 请求体**键名**被当成列名拼进 SQL（`{"id>0 or sleep(3) or id":"x"}`） | 只接受实体字段 / 列名（Java `getSql` 同样忽略未知键） |
| S12 | `warehouse/product_write.go:502` `bpDelBatch` | `ids` 元素未校验即拼进 `WHERE id=` | 先过 `validateBody` |
| S13 | SQL 字面拼接（Java `${}` 照搬）| `warehouse`: `template.go:566`、`noticeorder.go:221/558/567/1014`、`msgdomain.go:695-703`、`prolog.go:139-148/243-249/400`、`product_list.go:232/466/483`、`product_bp.go:320/346`、`product.go:64/96/115`、`remaining.go:193/209`；`commonrest`: `template.go:583-598`、`bespoke.go:102`；`contract`: `stats.go:132-143/177/195`、`stats2.go:256/387`、`sub.go:64/577-580`；`bill/sett.go:46`；`capital/small.go:342`；`authadmin/biz.go:54-56`；`logistics`: `bill.go:69/392/530/545`、`misc.go:99/794-799`、`waybill.go:56-92` | 值一律改参数绑定或转义；良性输入结果不变 |
| S14 | `commonrest/external.go` `extStream` | 上传 `.html` / `.svg` 后按扩展名内联返回，同域下可盗登录 token | 仅图片与 PDF 内联，其余强制下载并加 `nosniff` |
| S15 | `authadmin/user.go` `shapeUser`、`warehouse/msgdomain.go:385`、`product_get.go:27`、`commonrest/chat_local.go:39` | 任意登录用户可读到全表口令散列、AppSecret | 响应中 `password` / `secret` 置空（保留键，形状不变） |
| S16 | `auth/handler.go`、`compat.go` | 验证码缓存只增不删（键含客户端可控的 `X-Forwarded-For`）；图形码画布最大 4096×4096；新密码可为空 | 写入时清理过期项；画布封顶 512；空密码拒绝 |
| S17 | `authadmin/*` | bcrypt 对超过 72 字节的口令返回空串并被落库；上传目录取自可自改的 `phone` | 空哈希拒绝落库；目录取 `filepath.Base` |
| S18 | `auth/bank.go`、`platform/body.go`、`main.go`、各回调 / 上传 | 鉴权判定前无上限读 body；JSON 路由、multipart 导入均无上限 | 全局 512MB、JSON 32MB、回调 1MB、`uploadBatch` 100MB |
| S19 | `contract/esign.go:69` | `eqbFlowId` 未校验即拼进出站路径 | 仅允许字母数字与 `-_` |

### 2.2 🔴 / 🟡 正确性（Go 迁移引入）

| # | 位置 | 问题 | 修法 |
| --- | --- | --- | --- |
| C1 | `capital/detailed2.go:194` `countSeries` | 起点为 29–31 日时 `AddDate` 进位跳过结束月，循环永不退出直到 OOM（任何登录用户可触发，真实用户选「1 月 31 日到 2 月」也会中） | 按年月推进并限次 |
| C2 | `warehouse/product_write.go:368` `bpUpdateById` | 用驼峰键读下划线原始行，`associationId` / `createPhone` / `lastCheckTime` 等全部取空：每次更新重生成 `association_id`（附件脱离煤种），带 `createPhone` 时恒 999 | 行先驼峰化 |
| C3 | `warehouse/register.go` `Wrap2` | `busGlobalConfig/listByPage` 内层 500/999 被包成 HTTP 200 `success:true` | 非 200 原样透传 |
| C4 | `warehouse/template.go:565` | 无 `page` 时用了 `listSQL` 而非 `listNoSQL`（`busProLog` 少 `visitorHeadimg` 语义） | 改用 `listNoSQL` |
| C5 | `contract/stats.go:98-140` | 前端未选筛选时发空串，Go 当成 `false` / 空日期拼条件：开启过签章或在线支付的合同默认查不到 | 空串视为未传（Java 绑定为 null） |
| C6 | `warehouse/prolog.go:139-148` | `startTime:""` 对 nil 断言 panic | 判空跳过 |
| C7 | `logistics/excel.go:195`、`logistics.go:178` | 字段名 `Invoicing_time` 带大写：导入遇「开票时间」必 panic；接口传入的 `invoicingTime` 永远被丢弃 | 映射列小写；`pick` 大小写不敏感 |
| C8 | `logistics/waybill.go:406` `tencentSearch` | `count` 按 `float64` 断言恒得 0，「更新运单数据」永远是 0 条；第三方异常结构即 panic | 按 `json.Number` 取值并加守卫 |
| C9 | `logistics` / `bill` / `capital` 的 `parseSharedStrings` | 富文本单元格被按 `<t>` 拆成多条，其后所有字符串列索引后移，导入「成功」但整列错位 | 复用 `platform.ExcelSharedStrings` |
| C10 | `warehouse/member2.go:315` | 已签日期读错键名（`hjjNum` → `hjj_num`），只显示「+」 | 改键名 |
| C11 | `warehouse/product_write.go` `bpLog` / `bpSaveDelLog` | 日志写失败后调用方继续执行并回「成功」 | 返回 bool，失败即中止 |
| C12 | `warehouse/product_write.go:602` `bpImportExcel`、`contract/contract.go` `removeContract` | 逐行自动提交，中途失败半落库；删合同逐条吞错 | 收进事务 |
| C13 | `warehouse/product_bp.go:346` | `searchStr` 为数字时类型断言 panic | 走 `bindVal` |
| C14 | `capital/detailed2.go:274` `bucketSeries` | 补月循环每轮都取起点月，中间空月不补 | 按月推进 |
| C15 | `capital/detailed2.go:161`、`settlement.go:579` | `ScanMaps` 错误被吞，断连时返回全 0 的「成功」 | 判错回 999 |
| C16 | `contract/sub.go:284` `supplierMoney` | `float64` 相加再回写，大额丢精度、并发丢更新 | 交给 MySQL 做十进制加减 |
| C17 | `contract/esign.go:131、245` | 免费签章份数并发可扣成负数，之后「份数不足」校验永久失效 | 判断改 `<= 0`，自减加条件 |
| C18 | `contract/esign_file.go:186` | 文件状态 `for {}` 无限轮询，且在事务内 | 最多 5 分钟 |
| C19 | `commonrest/tlzf_schedule.go:23` | 每次支付 / 确认都常驻一个每 5 秒查库的 goroutine，只增不减 | 覆盖前取消旧任务 |
| C20 | `commonrest/external.go`、`authadmin/common.go`、`contract/esign_file.go` | 预览 PDF 用毫秒时间戳命名，并发同毫秒互相覆盖 | 时间戳 + 随机后缀 |
| C21 | `platform/jobs.go` | 处理器 panic 带崩整个进程且重启后同一条再崩；坏消息每 250ms 重读刷日志；`.seen` 只增不删 | `recover`、隔离 `.failed`、过期清理 |
| C22 | `platform/dbutil.go` `WithTx` | 回调 panic 时不回滚，连接池（20 条）被耗干 | `defer Rollback` |
| C23 | `authadmin/authadmin.go:193` `pageOf` | `size:0.5` 取整为 0，`IPage` 整数除零 panic | `size >= 1` 才覆盖 |
| C24 | `main.go` | `http.ListenAndServe` 无任何超时（慢连接可拖垮服务） | 请求头 10s、空闲 2min |
| C25 | `main.go` `maskDSN` | 用户名 ≥12 字符时密码明样打进日志 | 用驱动自带解析打码 |
| C26 | `auth/member_expiry.go:67` | Redis 键名内容拼进 SQL | 参数化 |
| C27 | `dashboard/logistics.go` `countLogist` | 月份区间无上限（0001–9999 年即 12 万次循环） | 1200 封顶 |
| C28 | `bill/bill.go`、`capital/shares.go` `writeVals` | body 全部键当列写入：载体键导致 1054，非法键名触发 `CheckCols` panic | 过滤载体键与非法键名 |
| C29 | `warehouse/paynotify.go:326` | `http.DefaultClient` 无超时，微信接口挂起时占住事务连接 | 15s 超时 |
| C30 | `platform/legacy_excel.go`、`uploadrender.go` | LibreOffice 子进程并发无上限 | 并发闸门（2） |
| C31 | `warehouse/template.go:742`、`commonrest/template.go:824` | `updBatch` 解码未用 `UseNumber`，金额 / 大整数经 `float64` 丢精度 | `UseNumber` |
| C32 | `authadmin/user.go:444` | 邀请码缓存只增不删 | 生成时清理过期项 |
| C33 | `logistics/misc.go:354/492` | 缺 `BatchNumber` 头写空串致重复校验失效；`LIKE` 定位可被 `%` 命中任意行 | 落 NULL；改等值 |
| C34 | `warehouse/product_get.go:143` | 遗留调试 `println` 每次详情页写日志 | 删除 |

## 3. 前端：已修复

| # | 位置 | 问题 | 修法 |
| --- | --- | --- | --- |
| F1 | `api/http.js` | 网关未授权是 HTTP 200 + body `code:403`：token 失效后只弹提示，用户被搁浅，刷新白屏；图形验证码是 Blob，也被当业务消息弹 | 识别 403（无 `data` 键）→ 清登录态整页回登录页；二进制响应不弹 |
| F2 | `Emain/EsignIndex.vue:253`、`Sponsor.vue:320` | 签署人姓名拼进 HTML 字符串渲染（`dangerouslyUseHTMLString`），存储型 XSS 可盗 token | 拼接前转义 |
| F3 | `Emain/InitiateStart.vue:459`、`Mindex/PurchaseContract.vue` 付款 | 「发起签署」「付款确定」无防重：双击发起两条签署流程 / 两笔担保支付订单 | 3–5 秒节流 |
| F4 | `App.vue`、`Main.vue`、`NotFont.vue`、`Xzmain.vue` | 非整页登出后同标签页换账号，继承上一个账号的路由与按钮权限 | 登出后整页刷新 |
| F5 | `filters.js` | 全局过滤器对 `null` 不设防，渲染期 `TypeError` 整表不刷新；`ToFiexd` 对空值显示 `NaN` | 入口空值返回空串 |
| F6 | `api/path.js:91-93` | 3 个前缀多一个 `/`，每次请求多一跳 307（原报告判断为 405，实测 Go 1.26 返回 307，功能可用） | 去掉多余斜杠 |
| F7 | `views/Main.vue:505` | 改密弹窗关闭时 `this.ruleForm = {}`，字段失去响应性 | `resetFields()` |
| F8 | `views/Login.vue:405` | 获取验证码在账号不存在（2002）或请求失败后按钮永久禁用 | 非 200 即恢复 |
| F9 | `views/Forget.vue` | 使用未导入的 `Message`，验证码错误时 `ReferenceError` | `this.$message` |
| F10 | `views/Register.vue:578`、`setup/EnterpriseArchives.vue:593` | 大小校验多除一次 1024，提示 3MB 实为 3GB；`Register.vue:442` 把明文密码打到控制台 | 修正；删除日志 |
| F11 | `setup/globalConfig.vue:159` | 红包金额最小值永远改不了 | 条件改用最小值本身 |
| F12 | `setup/EnterpriseArchives.vue:478` | `splice(value,1)` 起始位是对象，删一张图会挤掉相邻图 | `filter` |
| F13 | `setup/ContractDeploy.vue`、`setup/menuPower.vue` | 搜索先请求后重置页码，第 3 页搜索显示「暂无数据」 | 先重置再请求 |
| F14 | `setup/EditContractDeploy.vue:536` | 保存 / 取消后 `queryAll` 被清，后退到本页整页报错 | 为空时跳回列表 |
| F15 | `router/index.js:76` | 路由路径带尾空格 | 修正 |

## 4. 后端：原 Java 遗留，需决策（本次未改）

### 4.1 认证与授权边界

| # | 位置 | 问题 |
| --- | --- | --- |
| B1 | `auth/compat.go:79`，`commonrest/chat_local.go:62-84、183、490-552` | 无密码 / 无验证码登录族：`authentication/mobile`、`loginByAccount`、`buslogin`、`loginChatApp`、固定手机号免码 |
| B2 | `authadmin/user_compat.go:216-287、317-340、374-378`，`user.go:222-243、377-500` | 用户域写接口无权限边界：改任意人密码、自建管理员、自审批、删除 / 注销任意用户 |
| B3 | `auth/middleware.go:68-83` + `frontend/src/api/http.js:13-14` | 前端静态签名 = 全站匿名凭据（SignedAnon 可过所有 `authGuard`） |
| B4 | `authadmin/common.go:282-299、382-385、544-551、677-704`，`role.go:128-218`，`menu.go`，`biz.go` | 平台配置、角色、菜单写接口无管理员边界；支付 / 通联配置字段明文回显 |
| B5 | `commonrest/chat_local.go:150-163、312-325、328-345`，`chat_weixin.go:771-776` | `updByOpenid` / `updByOpenId` / `reg` 全字段批量赋值可自提权；注册积分池可被抽干 |
| B6 | `commonrest/enterprise.go:139-150` | 企业注册用 body 的 `userId` 覆写任意用户 |
| B7 | `authadmin/user_compat.go:110-113` | `setAccount` 明文密码入库（与 `docs/validation` 3.2 第 5 条「Go 使用 BCrypt」相反，文档或代码有一处过期） |
| B8 | `auth/handler.go` `Login` | 无登录失败计数 / 锁定 |

### 4.2 租户与数据归属（全仓性）

| # | 位置 | 问题 |
| --- | --- | --- |
| B9 | `warehouse/template.go`、`commonrest/auth_rest.go`、`register.go` | 通用 CRUD 模板无归属 / 租户校验且全列批量赋值：自封会员、自加积分、改账单状态、自提权角色、伪造通联支付单 |
| B10 | 各包 `getById` / `updateById` / `deleteById` / `delBatch` / 导出 | 按 id 直取直改直删、导出全表，遍布 logistics、capital、bill、contract、warehouse |
| B11 | `bill/billctl.go:361、391`，`contract/contract.go:797`，`logistics/bill.go:298`，`capital/detailed.go:321` | 按 `relation_id` 删附件不带 `relation_type`，删发票 #1000 会连带删掉结算单 / 运单 #1000 的附件（各实体都写 `relation_type=1`，需先定业务取值） |
| B12 | `warehouse/red_transfer.go:79-103、160-161`，`content_create.go:187-210、262-279`，`member2.go:207-294` | 红包出款金额、收款人、流水归属取自请求；视频审核任意用户可做 |
| B13 | `warehouse/noticeorder.go:69` | `searchSql` 为客户端传入的原始 SQL 片段（需设计白名单） |
| B14 | `warehouse/msgdomain.go:207-209、228、249、385、485-519`，`videos.go:161、686`，`leasing_write.go:51、130`，`content_create.go:119`，`content_external.go:141`，`noticeorder.go:369、549、829` | 按 id 读写不校验归属；`custNoticeCust` 可向任意号码群发短信 |

### 4.3 资金与并发

| # | 位置 | 问题 |
| --- | --- | --- |
| B15 | `warehouse/paynotify.go:191-250` | `payNotify` 无金额 / 订单 / 幂等校验，且在事务内同步发短信（`docs/validation` 3.2 第 9 条已登记） |
| B16 | `capital/payment.go:175-216、289-344` | 招行退款 / 确认 / 查单只信请求体，不校验订单归属 |
| B17 | `commonrest/tlzf_confirm.go:76-103、155-175`，`tlzf_schedule.go:224-236` | 通联确认 / 同步无事务无行锁，先改本地再调通道；`start` / `stopAll` 无角色校验 |
| B18 | `capital/balance.go:46-99、307-324`，`warehouse/red_transfer.go:62-72`，`product_get.go:386-420`，`product_write.go:615-669` | 核销金额由前端给出且无行锁；红包日上限先查后转；计数器读改写丢更新 |
| B19 | `logistics/logistics.go:377-395` 等 | 系统编号 `max+1` 无锁，并发重复 |

### 4.4 其他

| # | 位置 | 问题 |
| --- | --- | --- |
| B20 | `warehouse/noticeorder.go:833`、`videos.go:747` | 删除视频先删文件行再查行，物理文件与 `_img` 目录永不清理（Java 同序） |
| B21 | `config.go:27` | 默认 DSN 为 `root` 无密码（本地开发缺省），部署必须显式配置环境变量 |

## 5. 覆盖范围与未完成项

### 5.1 后端：7 个端点缺 Java 定制写入逻辑（切换前必须补）

生成器把带 `HttpServletRequest` 的写端点当成模板直通。已用后端机 JAR 的字节码逐个确认（`javap`），Java 的 `add` 并不是 `service.save(entity)`：

| 端点 | Java 实际做的事（Go 目前只插一行 body） |
| --- | --- |
| `busScoreExchangeLog/add` | 保存兑换记录后**扣减积分**（`updUserScore`）、写站内消息与用户消息、按配置**发短信**。Go 下积分兑换不扣分 |
| `busVerifyInfo/add`、`updateById` | 生成 `associationId`、把 `filelist` 挂到当前用户并置 `effective`，**把煤种 `is_real_check` 置真** |
| `busBusinessCard/add` | 清除该用户默认名片，新名片 `isDefault=true`，`userId` 取登录人 |
| `busContacts/add` | 按登录人 + `cardId` 查重，`userId` 取登录人 |
| `busEnterpriseIntroduction/add` | `userId` 取登录人，logo / 图片置 `effective` |
| `busVisitorLog/add` | 按登录人 + `targetUserId` 存在则更新、否则新增 |

这些接口由小程序 / App 调用，Web 前端未使用。未修原因：需要按字节码补约 300 行业务并逐个验收，属独立迁移任务，不适合夹在评审修复里。字节码摘录已留在评审工作目录。

### 5.2 前端：已定位、未修（改动面大或需产品确认）

- **购销合同页**（`Mindex/Add|Edit|PurchaseContract.vue`、`AddSalesContract.vue`）：编辑页离开不清 `queryAll`，之后「新建」会覆盖保存（`EditPurchaseContract.vue:997、1953`，另两份同）；签署方按可回退的 `numbers` 编号，删除后再加会重号合并（`AddPurchaseContract.vue:1205-1241`，三页同）；草稿回显的附件删除后仍随保存提交（`:1904-1917`）；保存 / 查询失败后 `fullscreenLoading` 不复位（`:1549-1587`）；「签署方不能少于一个」提示后没有 `return`（`:1532-1536`）；「签署状态」筛选发的是 `isEsign`、表格按 `isEsignComplete` 显示（`PurchaseContract.vue:1320`）。
- **角色权限页 ×3**（`setup/Roleperm.vue:171-199`、`PtRoleperm.vue`、`QyRoleperm.vue`）：回显勾选触发 `handleNodeClick`，把新勾上的菜单操作权限写大；切换角色时 `roleID` 先改后请求，失败或乱序会整体串写。
- **签署发起页**（`Emain/InitiateStart.vue:375-411、631-645`）：签署方按 id 分组删除后重号；删除已传文件后 `field` 未清空，点发起静默无反应。
- `setup/menuPower.vue:371-381`、`QyRoleperm.vue:110、411` 弹窗重开后表单状态残留；`ContractDeploy.vue:565-586` 等业务失败后 `loading` 不复位。
- `Xzmain.vue:641-675` 看板筛选无请求序号，快切地区旧响应覆盖新数据；`echarts` 实例未销毁；二月末日写死 28。
- `Login.vue:182` 全局 `document.onkeydown` 只在登录成功时清除；`App.vue` `created` 与 `mounted` 各跑同一段逻辑；`CommonView.vue:79` 引用未定义的 `this.MainPage`；`upload.js` 在模块加载时快照 token。
- 死代码：`api/file.js` 等 12 个无页面调用的函数（其中 12 个指向后端不存在的路径）、`views/Service.vue`（公开路由挂的是 bpmn 建模器）、`views/Demo.vue`、`components/link.js`。
- 发版提醒：`api/path.js` 当前指向 `http://127.0.0.1:8765`（`env = "developmen_test"`），构建上线前须切回域名；`path.js.orig-backup` 是入库的备份文件，建议删除。

### 5.3 未逐行评审的范围

受账户 API 限流影响，并行评审只能 3 组同时跑。以下前端页面**没有逐行读**，只被机械比对覆盖：`views/Mindex` 其余 67 个页面（运单、票据、统计、收付款、客户等）、`views/Ptmain` 22 个页面、`views/DPmain` 8 个页面。机械比对结论：`src/api` 的 402 处调用里 384 处命中已注册的 Go 路由，0 处方法不匹配；未命中 18 处 = 2 处登录路由不在 Java 源路由表内（已注册）、3 处已注释、12 处死函数、1 处带查询串的下载（已注册）。

`tools/`、`scripts/` 下的 Python 校验脚本不在本次范围。`staticcheck` 与 Go 1.26 不兼容未能运行，以 `go vet` 全量代替。

## 6. 验证

- `cd backend && go vet ./...`：无告警。
- `go test ./... -count=1`：全绿。
- `PANDA_LOCAL_AUDIT=1 go test ./... -count=1`：全绿。本机 MySQL 8.0.46 导入后端机 `20261117` 备份的 `chain`（146 表）/ `zjjg`（117 表），隔离 fixture 用例与写验证用例全部执行。
- 新增回归测试：`dsn_mask_test.go`、`auth/sms_test.go`、`platform/robustness_test.go`、`commonrest/pathguard_test.go`、`warehouse/capital/logistics` 的 `review_fixes_test.go`；并扩充了 `platform/wechat_cert_test.go`、`contract/esign_test.go`（过期时间戳的回调被拒绝）。
- 前端：20 个改动文件的脚本块逐个 `node --check` 通过；行尾风格与各文件原样一致（仓库内 CRLF / LF 混用）。仓库没有 `node_modules`，未做构建与浏览器回归。
- 行为对齐：所有修复对良性输入保持原响应；会改变响应的只有「原本就是错误 / 攻击路径」的请求。

## 7. 🎉 值得保留

- `warehouse/warehouse.go` 的 `bodyMap` 统一 `UseNumber`，配合 `bindVal` 把 DECIMAL / 整数按字符串绑定，模板引擎全程没有浮点运算。
- 资金、票据、通联金额全部走字符串十进制（`big.Rat` / `dec*`），与 Java `BigDecimal` 语义一致。
- Excel 导出单元格一律 `inlineStr` + `xmlEscape`，天然免疫公式注入；外部适配器的地址与密钥全部来自环境变量，源码零密钥。
- RSA 密钥文件 `0600` / 目录 `0700`，JWT 校验显式限定 RS256；`reg` 先查 `mux.Handler` 再注册，配合路由防冲突测试。
- 前端 `store` 按角色只注册对应路由并重置 `matcher`，路由级隔离干净。
