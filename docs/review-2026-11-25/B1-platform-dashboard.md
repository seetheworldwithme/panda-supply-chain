徐先生，以下是「共享基础设施 + 仪表盘统计」范围的只读审查报告（全程未改动任何文件、未运行 go build/test）。

## 范围
- 全文读完 29 个文件：`backend/internal/platform/` 的 allinpay/esign/jobs/legacy_excel/redis/tencent/uploadrender/wechat_cert/wechat_http/wechatpay/wxconfig 及 9 个对应测试（20）；`backend/internal/dashboard/` 4 个文件；为确认可达性另读了 `main.go`、`auth/middleware.go`、`auth/bank.go`、`auth/member_expiry.go`、`platform/{resp,dbutil,maps,body}.go`。
- 为核实调用方，部分读取：`warehouse/paynotify.go`、`warehouse/jobs.go`、`commonrest/{jobs,tlzf,tlzf_callback,chat_provider,external_convert}.go`、`commonrest/external.go`（470-579、640-749）、`authadmin/common.go`（640-829）、`contract/esign_file.go`。
- `body_contracts_gen.go` 按要求只看头尾并 grep：1256 条全是 `kind: object|array, required: true`，没有 `raw`/`ANY` 条目，所以 `body.go` 里这两个分支目前是死代码。

## 发现（按严重度降序）

### [🟡] 图片压缩可被解压炸弹打到整进程 OOM (parity:false)
- 位置：`backend/internal/platform/uploadrender.go:28,42`；入口 `authadmin/common.go:762`、`commonrest/external.go:530`（`/saas-common/upload/*`，需登录）
- 证据：
  ```go
  im, format, e := image.Decode(bytes.NewReader(b))
  ...
  dst := image.NewNRGBA(image.Rect(0, 0, ww, hh))
  ```
- 场景：登录用户上传一张声明 60000×60000 的纯色 PNG（几 MB，补零过 800KiB 即进入解码分支），按扩展名判为图片后直接 `image.Decode`。本机 GOROOT 的 `image/{png,jpeg,gif}` 解码器没有像素上限，`NewNRGBA` 会申请约 14GB，单个 Go 二进制被 OOM 杀掉，全部服务一起下线。仓库里没有任何 `MaxBytesReader`。Java 端是受 `-Xmx512m` 限制的 `OutOfMemoryError`，被外层 try 吞掉只影响单个请求，所以影响面是移植后放大的。
- 修法：解码前先 `cfg, _, err := image.DecodeConfig(bytes.NewReader(b))`，`int64(cfg.Width)*int64(cfg.Height)` 超过约 5000 万时原样返回 `b`。
- 把握：high。能否定它的证据：解码前另有尺寸校验，或上传入口有体积上限（已 grep，都没有）。

### [🟡] 公开回调可无限触发微信平台证书下载，且全程占全局锁 (parity:unknown)
- 位置：`backend/internal/platform/wechat_cert.go:31-32,64-68`；调用方 `warehouse/paynotify.go:121`（`/saas-warehouse/busMyBill/appPayNotify` 免登录）
- 证据：
  ```go
  wechatCertificates.Lock()
  defer wechatCertificates.Unlock()
  ...
  cert := selectCert()
  if cert != nil && (cfg.APIV3Key == "" || time.Since(wechatCertificates.refreshed) < time.Hour) {
  ```
- 场景：生产配置了 APIv3Key 时，攻击者对 appPayNotify 每次带一个随机 `wechatpay-serial`。`selectCert()` 返回 nil，必然走到 `downloadWechatCertificates`，这是对微信的一次 HTTPS 请求，客户端超时 30s，期间全局互斥锁不释放。下单、打款、正常回调的 `WechatV3JSON` 和 `WechatPlatformCertificate` 全部排在后面；`Lock()` 不响应 context，请求方断开也不会提前退出。持续几十 req/s 就能让微信支付路径基本不可用。
- 修法：对「未知序列号」做负缓存，例如记录上次下载尝试时间，1 分钟内直接返回错误。
- 把握：high（行为由代码直接可见）。注释称「未知序列号立即刷新」是沿用原 SDK，我没有核对 wechatpay-java 源码，所以 parity 标 unknown。

### [🟡] 后台任务 worker 无隔离：坏消息不隔离，处理器 panic 会变成崩溃循环 (parity:false)
- 位置：`backend/internal/platform/jobs.go:152-168`、`backend/main.go:49-51`
- 证据：
  ```go
  log.Printf("后台消息无法读取: %s", filepath.Base(path))
  continue
  ...
  handler(id, job.Payload)
  if err := os.Remove(path); err != nil { return err }
  ```
- 场景：
  - (a) `Run` 跑在独立 goroutine，没有 `recover`。任一处理器 panic 会带崩整个进程，而消息文件要等处理器返回才删除，重启后同一条消息再 panic。按 AGENTS.md 没有 systemd 拉起，结果就是起不来。
  - (b) 内容损坏的消息或没有注册处理器的消息名永远不会被移走，每 250ms 重读并打一行日志，约 4 行/秒，无限持续。
  - (c) `*.json` 但文件名不含 `-` 的 stray 文件会让 `SplitN(...)[1]` 越界 panic，效果同 (a)。
  - (d) `<id>.seen` 标记只写不删，目录无限增长，`Glob` 每 250ms 全量 readdir。
  - (e) 单 worker 串行，`warehouse.product` 每条固定 `Sleep(2s)`，`common.tlzf` 支付状态消息排在后面；批量发 500 条煤种后，支付状态同步要等十几分钟。
- 修法：`handler` 外包一层 `defer recover()`，出错时把文件改名为 `.failed`；读不出或无处理器的消息同样改名；`.seen` 过期后顺带删除。
- 把握：(b)(d)(e) high；(a)(c) medium，我没有逐一证明两个处理器内部存在可达 panic，需要负责 warehouse/commonrest 的 reviewer 确认。

### [🟢] e签宝回调验签没有时间戳窗口，可被重放 (parity:unknown)
- 位置：`backend/internal/platform/esign.go:49-56`；使用方 `contract/esign_file.go:85-110`、`auth/bank.go:48-55`
- 证据：
  ```go
  data.WriteString(r.Header.Get("X-Tsign-Open-TIMESTAMP"))
  ...
  return hmac.Equal([]byte(hex.EncodeToString(mac.Sum(nil))), []byte(r.Header.Get("X-Tsign-Open-SIGNATURE")))
  ```
- 场景：152 目前是 IP+HTTP 明文（AGENTS.md §5）。截获一条合法的 `SIGN_FLOW_FINISH`、`flowStatus=2` 回调后，事后重放会再次执行 `UPDATE chain_contract SET is_esign_complete=1`，把已重置的流程重新标成已签完；中间件同时为该请求放行 `SignedAnon`。
- 修法：验签前比较时间戳与当前时间，超过 5 分钟直接返回 false。
- 把握：medium。能否定它的证据：Java 端原本也不校验时间戳，且产品接受重放风险。

### [🟢] countLogist 月份区间无上限，登录用户可触发 CPU/内存放大 (parity:true)
- 位置：`backend/internal/dashboard/logistics.go:129-208`
- 证据：
  ```go
  num := jg + 1
  ...
  for i := 1; ; i++ {
  cur := start.AddDate(0, i-1, 0)
  if !strings.Contains(skdate, mo) {
  ```
- 场景：不传 timeType，`startTime=0001-01-01`、`endTime=9999-12-31`，jg=119999。结果是 `LIMIT 120000` 的派生表、chain_logistics 全表聚合、12 万次循环，`strings.Contains` 在不断变长的 `skdate` 上是 O(n²)，再加 24 万个 JSON 项，单请求要数秒 CPU。同文件的 `monthSeries` 已有 1200 上限，这里漏了。
- 修法：算出 jg 后，大于 1200 直接 `platform.Unknown(w)`。
- 把握：high。注释写明补月循环与 SQL 照 Java，所以 Java 有同样问题。

### [🟢] 外部响应读取无大小上限
- 位置：`allinpay.go:127`（`io.ReadAll(res.Body)`）、`wechatpay.go:71`；`esign.go:80`、`tencent.go:73`、`wechat_http.go:40` 是无上限的 JSON 解码
- 场景：endpoint 全部来自环境变量或常量，只有服务商被劫持或配置错误时才会出问题，风险低。
- 修法：`io.LimitReader(res.Body, 16<<20)`。

## 重复模式
- 预览文件名用毫秒时间戳且落在共享目录 `convertFiles`：`authadmin/common.go:788`、`commonrest/external.go:553`、`contract/esign_file.go:165`。同毫秒的并发上传会互相覆盖 PDF，属于跨租户串档；`RenderExcelUpload` 的输出目录按文件隔离，不受影响。修法：文件名改用 `newUUID()`。

## 跨范围线索（未深查，供对应 reviewer 核对）
- `warehouse/paynotify.go:326`：微信 v2 统一下单用 `http.DefaultClient`，没有超时。
- `contract/esign_file.go:158-205`：e签宝文件状态的轮询 `for {}` 包在 `WithTx` 里，转换不完成就一直占着一个连接池连接（池只有 20）。
- `contract/esign.go:70-94`：`EsignJSON` 的 path 直接拼请求字段 `eqbFlowId`，服务器会用 appKey 给用户选定的 e签宝路径签名。
- `main.go:56`：`http.ListenAndServe` 没有任何超时，全仓没有 `MaxBytesReader`。

## 已核对无问题
- **仪表盘 SQL**：`common.go`、`contract.go`、`logistics.go` 里凡是请求来源的值（companyId、startTime/endTime、xzqhCode、parentCode）都走了参数化；拼接进去的只有常量、`col` 字面量和重新格式化过的 `endD.Format`；`timeTypeCond` 返回常量。所有 `rows` 都进入 `scanMaps`（内部 Close），没有泄漏。不加租户条件是包注释写明的 Java 行为，不作为发现。
- **通联**：XML 键全是常量，值经 `xml.EscapeText`；openssl 和 iconv 的参数不来自请求，私钥只写进 0700 临时目录并 defer 删除；验签针对原始字节并去掉 sign_code。
- **腾讯云 TC3**：规范化请求、签名作用域和派生密钥都与协议一致，30s 超时，endpoint 只来自环境变量。
- **微信支付 v3**：响应验签覆盖 `timestamp\nnonce\nbody\n`，有 5 分钟窗口，未知序列号会拒绝；证书下载用 AES-256-GCM 并做财付通 CA 链校验；凭据从仓库外的私有文件读取；私钥类型断言都带 ok。
- **e签宝请求签名**：HMAC-SHA256 加常量时间比较，appSecret 为空直接拒绝；`EsignUpload` 的 URL 来自 e签宝响应，不是用户输入。
- **Redis**：每条命令新建连接并 defer Close；错误字符串不含凭据；RESP 解析对 bulk 16MB、数组 10 万元素设了上限；订阅断线由 `auth/member_expiry.go` 的循环 1 秒后重连。
- **LibreOffice/pdftoppm**：参数不含请求数据，独立 profile 目录，临时目录 defer 清理，超时 30s/60s；`OpenLegacyExcel` 输入限 32MB，两个调用方都 `defer cleanup()`。
- **测试**：9 个测试都能失败，也都断言了正确行为，没有空测试。