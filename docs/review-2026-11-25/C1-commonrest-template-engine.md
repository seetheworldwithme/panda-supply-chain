徐先生，commonrest 通用 CRUD 模板引擎评审完成（6 个回合，只读，未建任何文件）。

## 范围
范围內 15 个文件全部读完，meta_gen.go / auth_meta_gen.go 也读了全文。为核对路由归属与守卫，补读了 platform/{dbutil,maps,resp,body,legacy_excel}.go、auth/middleware.go、authadmin/{authadmin,user,role}.go、main.go:56-142。按协调指令未读 skills 指南；同包 chat_*、external*、tlzf_*、contract_batch、enterprise、ocr、qcc 不在范围。未连库、未构建，表是否真有 `system_code` 列只能按实体列推断。

## 发现

### [🔴] 权限、密钥、支付配置表经模板 CRUD 暴露给任意登录账号 (parity:true)
- 位置：backend/internal/commonrest/auth_rest.go:31-78、register.go:34-45
- 证据：
```go
reg(base+"/add", func(w http.ResponseWriter, r *http.Request) { s.tplAdd(w, r, e) })
reg(base+"/updateById", ...); reg(base+"/deleteById", ...); reg(base+"/delBatch", ... tplDelBatch(w, r, e, true))
```
- 场景：`reg` 只补专用路由的空缺（main.go:123 最后注册），下列都是模板路由，只有 authGuard。引擎对 `company_id` 既不自动填充也不按调用者过滤，读写都信 body。
  - `/saas-auth/chainAuthSecret/{list,getById,add,updateById,deleteById,delBatch,exportExcel}`：`list {}` 读出全部 HMAC secret；`add` 自造一条后按 middleware.go:132 用 `HMAC(secret, id+"hyqc"+secret)` 得 SignedAnon=true，无账号通过全站 authGuard。
  - `/saas-auth/chainAuth{UseRole,MenuUser,AuthorityUser,AuthorityRole,MenuRole,Authority}/*`：`chainAuthUseRole/add {"userId":自己,"roleId":平台角色}` 自提权。
  - `/saas-common/chainTlzf{Order,Config,ContractLog,Log,SaleEtc}/{add,updateById,deleteById,delBatch,updBatch}`：改他企对公账号 `acct_num`、改支付单 `trans_status` 和金额。
  - `/saas-common/chainBasicEnterprise/{deleteById,delBatch,updBatch}`、`chainGlobalConfiguration/{add,deleteById,delBatch}`、`chainEqbConfig/{deleteById,list,exportExcel}`。
  - `/saas-contract/chain{,Cg,Xs}Contract/exportExcel`：全量导出所有企业合同。
- parity 依据：body_contracts_gen.go:14-34、55-68、88-101、531-540、557-572 由 JAR 的 @RequestBody 生成，说明 Java 有同款模板端点。专用 handler 同样无守卫（authadmin/user.go:27 `userWriteCols` 含 `account_type`、`level`），所以是全链路缺口，不是模板绕过专用守卫。
- 修法：两处 `reg` 包装里对上述实体加平台管理员守卫（复用 authadmin/user.go:556 的 `account_type IN (2,5)` 查询），至少拦住写端点和导出。会偏离 Java，需产品确认。
- 把握：high。否定条件：Java 网关另有按 URL 的角色校验（仓内无此证据）。

### [🔴] `${searchStr}` 字面替换，12 个 listByPage 端点可被登录用户 SQL 注入 (parity:true)
- 位置：template.go:583-585、593-599；meta_gen.go:65,82,99,133,150,167,456；auth_meta_gen.go:31,45,59,73,115
- 证据：
```go
if strings.Contains(q, "${searchStr}") {
    q = strings.ReplaceAll(q, "${searchStr}", javaStr(m["searchStr"]))
}
q = strings.Replace(q, "${maps.searchStr}", inner, 1)
```
- 场景：`searchStr` 只过 `bindVal("s")`，单引号原样进 SQL。`POST /saas-common/chainMsg/listByPage {"page":{"current":1,"size":10},"searchStr":"x' OR 1=1 OR '"}` 得全表；换成与列数对齐的 `UNION SELECT` 可读 `chain_auth_user.password` 等任意表（count 包装要求括号平衡，挡不住读取）。受影响端点：
  - `/saas-common/{chainMsg,chainAreaInfo,chainCommonAddress,chainCommonLog,chainCommonTemplate,chainCommonTempmx,chainCommonFile}/listByPage`
  - `/saas-auth/{chainAuthAuthority,chainAuthAuthorityRole,chainAuthAuthorityUser,chainAuthMenuRole,chainAuthUseRole}/listByPage`
  - chainCommonDict、chainContractExtend 已有专用路由，不受影响。
- 修法：两处替换前对值做 `strings.ReplaceAll(v, "'", "''")`，一行，正常输入语义不变。
- 把握：Go 侧 high；parity medium。Java mapper 原文同为 `${searchStr}`，MyBatis `${}` 不转义；若 Java 网关有全局 SQL 过滤器则 parity 不成立。

### [🔴] chain_auth_user 的 exportExcel 导出全体用户口令散列与 openid (parity:true)
- 位置：auth_rest.go:31、auth_meta_gen.go:142、excel.go:148
- 证据：`{"登陆密码", "password", "s", ...}`、`{"小程序openid", "openid", "s", ...}`
- 场景：任意登录账号 `POST /saas-auth/chainAuthUser/exportExcel`，得全量 xlsx（BCrypt 散列可离线爆破）。`/list`、`/add`、`/updateById` 已由 authadmin 接管，模板没有注册。
- 修法：auth_rest.go:31 对 `ChainAuthUser`、`ChainAuthSecret` 不注册 exportExcel。
- 把握：路由 high（别处未注册）；parity medium（表头取自 Java 实体 @ExcelColumn，未见导出时剔除的迹象）。

### [🟡] xlsx 解析按 `r` 属性补列无上限，几 KB 文件即可 OOM 整个进程 (parity:false)
- 位置：excel.go:926-931、121-131
- 证据：
```go
if ci, err := colIndex(ref); err == nil {
    for len(row.cells) < ci {
        row.cells = append(row.cells, xlsxCell{kind: "missing"})
```
- 场景：单元格写 `<c r="ZZZZZZZZ1"/>`，`colIndex` 得约 2×10^11 且永不报错，append 到内存耗尽。`/saas-contract/chainContract/importExcelCg`、`importExcelXs`（contract_import.go:36）均可达。同路径还有两个放大器：
  - `readZip` 对解压内容 `io.ReadAll` 无上限（excel.go:736），zip bomb。
  - multipart 路由不在 body contract 里，`ValidateJavaBodies` 直接放行（body.go:23），上传体积无上限。
- 修法：`colIndex` 结果超过 16384 时视为无效 ref、不补列；`readZip` 用 `io.LimitReader` 封顶。
- 把握：high。否定条件：前置网关限制了上传体积与内容（仓内无此证据）。

### [🟡] selPtAmt 的 companyId 未校验直接拼进 SQL (parity:unknown)
- 位置：bespoke.go:97-102
- 证据：`id := m["companyId"]` … `"SELECT sum(residue_freeze_amt) FROM chain_contract WHERE sales_unit_id=" + javaStr(id)`
- 场景：`{"companyId":"0 UNION SELECT ASCII(SUBSTR(password,1,1)) FROM chain_auth_user LIMIT 1"}` 可逐字节读任意表。同文件 selRelTraceFee（:114）先过 `validateBody`，此处漏了。
- 修法：函数开头补 `if !entChainTlzfOrder.validateBody(m) { platform.Unknown(w); return }`。
- 把握：Go 侧 high；Java 是否同为字符串拼接未核实。

### [🟢] tplUpdBatch 解码未用 UseNumber (parity:false)
- 位置：template.go:824
- 证据：`if err := json.NewDecoder(r.Body).Decode(&items); err != nil {`
- 场景：数字变 float64，超过 15 位有效数字的 bigint/DECIMAL 被改写；`bodyMap` 走 UseNumber，单条接口没这个问题。
- 修法：加 `dec.UseNumber()`。
- 把握：high，影响面小。

### [🟢] importTemplate 逐行插入无事务（当前不可达）(parity:false)
- 位置：excel.go:234-241
- 证据：`for _, row := range irt.DataList { ... if _, err := s.insert(r, entityForTable(table), vals); err != nil {`
- 场景：第 k 行入库失败时前 k-1 行已落库并返回 999，重试即重复；Java saveBatch 是事务。现在走不到：`/saas-common/*/importExcel` 因全部 `syscol:false` 在 template.go:870 直接 999；auth/contract 实体按列推断没有 `system_code`，`sysProbe` 先失败。
- 修法：改成 `platform.WithTx` + `TxInsert`。
- 把握：缺陷 high；“不可达”medium（未连库）。

## 重复模式
- **Excel 机器仓内有 5 份**，上面补列与 zip bomb 缺陷每份都有：commonrest/excel.go、bill/excel.go、capital/excel.go、logistics/excel.go、warehouse/excel.go。
- **`pageOf`/`paged` 除零**：`size` 取 0.5 时 `int64(0.5)=0`，`platform.IPage` 整数除零 panic。commonrest.go:112-137 这份没有调用方；authadmin.go:187-231 的同名副本有调用方（userListByPage、userSupList），请 authadmin 评审确认。
- **`.xls` 走外部进程**：每个上传起一个 `soffice`（platform/legacy_excel.go:44，30s 超时，无并发闸门）。输入封顶 32MB、临时目录会清理，这两点没问题。

## 已核对无问题
- 写入列全部走 `pick(m, e.dbcols)` 白名单并剔除 `id`，`validateBody` 按实体类型校验；请求键到不了列名位置，`CheckCols` 不会被请求触发。
- `tplList` 的 `allEqCond` 白名单加参数化；`getSql` 的列名来自元数据，值参数化。
- family C 的 `startTime`/`endTime` 经 `time.Parse("2006-01-02")` 后才拼接；`chain_common_address.type` 是 Integer 校验。两处都注入不了。
- `getById`/`updateById`/`deleteById` 的 id 校验为整数；`delBatch` 的 ids 参数化，非数字元素只得到 999。
- 模板分页自带 `ipageOf` 防除零；`size` 为 0、负、非数字都不 panic，`current≤0` 得 OFFSET 负 → 999（照 Java）。`size` 很大或 `/list` 空 body 即全表载入内存，同样照 Java，docs 3.4 已列为未覆盖。
- count 与数据查询同源（`stripOrderBy` 只去尾部 ORDER BY）。
- 导出单元格全部是 inlineStr 字符串并过 `xmlEscape`，没有公式注入。
- `tplUpdBatch` 用 `WithTx` 包多行更新，panic 时 `defer tx.Rollback` 能还连接。
- `mapTypeCache`（commonrest.go:200-226）名字带 sync，但写共享 map 没加锁；本包无调用方，是死代码，暂不构成竞争。
- jobs.go：订单号参数化，`Jobs` 只在启动期注入，platform/jobs.go:174 有 recover。
- contract_import.go：`system_code` 查询参数化，`systemCode[20:]` 有长度守卫。
- 四个测试都能失败，断言对象正确。

🎉 值得保留的两处设计：
- `reg` 先查 `mux.Handler` 再注册，配合 routes_test.go 的全量注册防冲突测试，模板与专用 handler 可以安全共存。
- 导出全部用字符串单元格加 `xmlEscape`，天然免疫公式注入。