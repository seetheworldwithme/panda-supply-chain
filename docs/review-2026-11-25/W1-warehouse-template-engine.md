徐先生，以下是 warehouse 通用 CRUD 模板引擎的评审报告（全程只读，未运行 build/test）。

## 范围
- **全文读完（7 个范围内文件）**：`meta_gen.go`（2566 行，分 10 段读到底）、`template.go`、`register.go`、`routes_gen.go`、`warehouse.go`、`jobs.go`、`template_test.go`。
- **为追调用另读**：`platform/{resp,dbutil,maps,body,jobs}.go`、`auth/middleware.go`、`main.go:30-129`、`excel.go:145-224`（仅 exportHandler / importTemplate 入口）、`product_sync.go:1-110`、`noticeorder.go:536-575`、`product_list.go:455-488`；证据文件 `docs/validation/java-go-20261119/server-java-routes.json`、`java-go-20261118/{java-source-routes,route-coverage}.json` 做了 grep 抽查。
- **未读**：`~/.claude/skills` 指南（协调方要求跳过）；`excel.go` 其余部分（他人范围）。
- **无法核对**：本地没有 Java 源码或 jar，方法体只能凭签名与生成器结论推断。

## 发现

### [🔴] `${maps.searchStr}` 字面拼接，`/busVisitorLog/listByPage` 可 SQL 注入 (parity:true)
- 位置：`backend/internal/warehouse/template.go:566-572`；`meta_gen.go:1845-1847`；`routes_gen.go:67`（只 skip `/list`）；`register.go:47-49`。
- 证据：
  ```go
  inner := "%null%"
  if s, ok2 := search.(string); ok2 {
      inner = s
  }
  q = strings.Replace(q, "${maps.searchStr}", inner, 1)
  ```
  meta 原文：``like '%${maps.searchStr}%'``，`searchLit: true`（全仓仅 busVisitorLog、busOrder 两处）。
- 场景：任意登录用户 `POST {"searchStr":"x%' or 1=1 or '","page":{"current":1,"size":10}}`，得到 `CONCAT(id) like '%%x%' or 1=1 or '%%'`，返回全表。UNION 对齐 16 列即可读 chain 库任意表（含 `chain_auth_user.password`）；`x%' and sleep(5) and '` 按行 sleep，可占满 20 连接的池。searchStr 含 ` order by ` 时 `stripOrderBy`（:578, :647）会截错，COUNT 报 999。
- 修法：`searchLit` 分支把 `'%${maps.searchStr}%'` 替换为 `?` 并 `append(args, search)`，良性输入行为不变；同步修正 `tools/migration/warehouse-meta.py`。
- 把握：high。若 Java mapper 原文并非 `${}`（即 meta 生成有误），则改为 parity:false，漏洞本身不变。

### [🔴] 模板读写无归属/租户校验，且全列批量赋值 (parity:true；包注释已声明有意保留)
- 位置：`template.go:677-679`（tplAdd）、`:694-696`（tplUpdateById）、`:762-765`（tplUpdBatch）、`:662`（getById）、`:710`/`:726`（删除）、`:476-492`（list）、`:798` → `excel.go:148`（全表导出）；`warehouse.go:20`「租户拦截器失效…一律不加 company_id」。
- 证据：
  ```go
  vals := pick(m, e.dbcols)
  delete(vals, "id")
  if _, err := platform.Insert(s.pick(r), e.table, vals); err != nil {
  ```
- 场景：
  - `POST /saas-warehouse/busUserMember/add {"userId":自己,"memberId":1,"memberType":1,"expirationTime":"2099-12-31T00:00:00","isEffect":true}` 可自封黄金会员。
  - `/busUserScore/updateById {"id":X,"hjjNum":999999}` 可自加积分。
  - `/busMyBill/updateById {"id":X,"status":2}` 可改账单状态。
  - getById/deleteById/delBatch 可按 id 跨企业读、改、删任意行。
  - `/busCustManage/exportExcel` 全表导出手机号。
- Java 对齐：`java-source-routes.json` 显示这些控制器在 Java 同样暴露 add/updateById/delBatch/getById/exportExcel；MP `save/updateById` 同样整实体落库、无租户过滤，是同一缺口。
- 修法：属产品决策，不建议单方面改。若要收紧，在 `pick` 后按实体强制覆盖 `user_id`/`company_id`，并把敏感列从 add/updateById 剔除。
- 把握：Go 侧 high；Java 侧依据生成器与路由清单，未见方法体。

### [🟡] 7 个 Java 方法签名带 `HttpServletRequest` 的写端点被当作模板 (parity:unknown)
- 位置：`routes_gen.go:12`（busBusinessCard）、`:16`（busContacts）、`:19`（busEnterpriseIntroduction）、`:51`（busScoreExchangeLog）、`:65`（busVerifyInfo，add 与 updateById）、`:67`（busVisitorLog）。它们的 Skip 均不含 `/add`，busVerifyInfo 也不含 `/updateById`；经 `register.go:56-59`、`:66-73` 进入模板。
- 证据（`docs/validation/java-go-20261119/server-java-routes.json` 行 15807、15957、16173、14991、18078、11676、11700）：
  ```
  "descriptor": "(Lcom/saas/warehouse/business/entity/BusScoreExchangeLog;Ljavax/servlet/http/HttpServletRequest;)Lcom/saas/cores/model/JsonResult;"
  ```
- 场景：同签名的 `busOpinion/add`、`busProCollection/add`、`busEnterpriseLeasing/add`、`busProduct/add`、`busVideos/add` 都被识别为 bespoke 并 skip，这 7 个却走 `tplAdd`/`tplUpdateById`，完全信任 body。
  - 若 Java 用 request 取 JWT 回填归属（名片、联系人、企业介绍），Go 可冒名替他人建记录。
  - 若 `busScoreExchangeLog/add` 在 Java 里扣积分，Go 只插日志，功能缺失。
  - `busVerifyInfo/updateById` 是唯一带 request 的 updateById。
- 修法：读或反编译这 7 个方法体确认；非 `service.save(entity)` 直通的补进 Skip 加 bespoke。
- 把握：medium。仅凭签名推断；若 Java 方法体与模板等价（request 未使用），则无差异。

### [🟢] Wrap2 把内层 999 也包成 HTTP 200 / success:true (parity:false)
- 位置：`register.go:107-117`、`:119-128`。
- 证据：
  ```go
  s.tplListByPage(rec, r, e, family)
  var inner any
  _ = json.Unmarshal(rec.buf, &inner)
  platform.Success(w, inner)
  ```
- 场景：`POST /saas-warehouse/busGlobalConfig/listByPage {"searchStr":{}}`。内层 `platform.Unknown` 写入 capWriter，`WriteHeader(500)` 被丢弃，外层回 200 `{"success":true,"code":200,"data":{"success":false,"code":999,…}}`。Java 的绑定/SQL 异常走 GlobalExceptionHandler，是未包一层的 500+999；前端会把报错当成功。
- 修法：`rec.code` 非 200 时把 `rec.header`/`rec.buf` 原样写回并 `WriteHeader(rec.code)`。
- 把握：medium-high。

### [🟢] 无 page 的 B/C 族路径用 `listSQL` 而非 `listNoSQL` (parity:false)
- 位置：`template.go:565`（只取 `e.listSQL`）、`:596-602`；`meta_gen.go:2241-2242`。`listNoSQL` 全后端零引用。
- 证据（busProLog 两条 SQL 不同，其余 67 实体相同）：
  ```
  listSQL:   `select (select wx_face_img from chain_auth_user where id=user_id) visitorHeadimg, bus_pro_log.* from bus_pro_log ...`
  listNoSQL: `select * from bus_pro_log ${searchSql} and CONCAT(id) like #{maps.searchStr} order by create_time desc`,
  ```
- 场景：`POST /busProLog/listByPage` 不带 page。Java 的 listNoPage 返回 `visitorHeadimg=null`，Go 返回头像 URL。
- 修法：无 page 分支改用 `e.listNoSQL`。
- 把握：medium，取决于 Java 无 page 时确实调 listNoPage（包注释如此）。

### [🟢] tplUpdBatch 解码未用 UseNumber (parity:false)
- 位置：`template.go:742`；对比 `warehouse.go:50-51`。
- 证据：`if err := json.NewDecoder(r.Body).Decode(&items); err != nil {`
- 场景：23 个模板 updBatch 端点（含 busRedEnvelopeDetail、busUserCoupon 的金额列）。大于 2^53 的整数和多位小数被 float64 舍入，会更新到错误行或丢精度。`template_test.go:31` 的大整数用例只覆盖 bodyMap。实际 id 远小于 2^53，影响小。
- 修法：Decode 前 `dec.UseNumber()`。
- 把握：high。

### [🟢] 固定 sleep 2s 落在串行 worker 内 (parity:unknown)
- 位置：`jobs.go:36`；`platform/jobs.go:148-169`（单循环逐个执行 handler）。
- 场景：批量审核通过 N 条煤种，后续同步与通知排队不少于 2N 秒；RocketMQ 消费端是并发的。只延迟，不丢失。
- 修法：接受现状，或把业务移入 goroutine 并保留 MarkDelivery 先行。
- 把握：medium。

## 重复模式
1. **`${maps.searchStr}` 字面注入**（parity:true，同一 Java bug）。他人范围内、请对应评审确认是否已报：
   - `meta_gen.go:1119-1120`（busOrder）→ `noticeorder.go:560-567`（`/busOrder/listByPage`，searchStr 未包 `%` 直接拼）。
   - `product_list.go:466`、`:483`（`/busProduct/listByPage`）。
   - `product_bp.go:320`、`msgdomain.go:684`。
   - `noticeorder.go:558`（phone 来自库，字面拼）。
2. **body 全信任**：`pick(m, e.dbcols)` 对 67 实体一致生效，见上面第二条。

## 已核对无问题
- **过滤键白名单**：`/list` 走 `allEqCond`（`warehouse.go:155-171`），listByPage 走 `e.q` 元数据。68 实体的 q 字段下划线化后均在 dbcols 内，dbcols 全部匹配 CheckCols 正则。请求键到不了 CheckCols，无 panic 路径。
- **值参数化**：除上述 searchStr 外全部 `?` 绑定。C 族 startTime/endTime 经 `bindVal("D")` 严格校验（`time.Parse("2006-01-02")`）后才拼接，不可注入。
- **分页**：`ipageOf` 对 size≤0 防除零（`template.go:504-507`）。page 非对象回 999；current/size 负值或溢出由 MySQL 报负 OFFSET 回 999，与 Java 同。仅 size 为非数字串时 Go 回空页而 Java 回 999，属吹毛求疵。
- **资源与事务**：`queryEnt`/`GetByID`/`loadMapTypes` 都关闭 rows；WithTx 内无可达 panic。
- **tplDelBatch**：没调 validateBody，但 ids 只作参数绑定，嵌套对象由驱动报错回 999。
- **importExcel**：68 实体 syscol 全为 false，模板导入恒 999。
- **Skip 清单**：36 个 skip 端点逐一找到 bespoke 注册（`noticeorder.go:31-53`、`msgdomain.go:34-38`、`product_write.go:22-37`、`product.go:21-23`/`143-147`、`prolog.go:20-29`、`leasing_write.go:14-15`、`content_create.go:53-55`、`product_list.go:30`、`product_get.go:18`、`remaining.go:34-37`；`/busRedEnvelopeRange/list|updBatch` 归 authadmin），无 404 空洞。
- **并发**：ent 元数据只读，`mapTypeCache` 是 sync.Map，`s.Jobs` 启动期写入，无数据竞争。
- **jobs.go**：MarkDelivery 先行、失败只记日志，与注释声明的 Java 行为一致。fallback goroutine 只在 `jobs==nil`（离线路由核查）时可达，生产不可达。
- **template_test.go**：5 个用例都有可失败的断言，无空测试。

## 🎉
- `warehouse.go:48-57` 的 bodyMap 用 UseNumber，配合 `bindVal` 把 DECIMAL/整数以字符串绑定，模板引擎内全程没有 float 运算。
- `bindVal`/`validateBody` 复刻 Jackson 绑定面，类型不符统一回 999，且发生在拼 SQL 之前。