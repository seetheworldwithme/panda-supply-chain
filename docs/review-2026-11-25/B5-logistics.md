徐先生，logistics 包评审报告如下（只读，未改动任何文件；路径均相对 `/Volumes/code/GoProject/panda-supply-chain/backend/internal/`）。

## 范围
全文读完 logistics 包全部 8 个文件（bill.go 567、excel.go 1344、logistics.go 652、misc.go 1029、waybill.go 449、三个测试 12/24/33 行），另读 platform 的 resp/dbutil/maps/body/javamap 与 auth/middleware 作行为参照（不评审）。未读 dashboard 包里的 countBuses/countMonBuses/countLogist（不在范围）；按要求未读 skills 指南。colIndex 补列、readZip 上限未重复报，但补列的修复不完整，见第 3 条。

## 发现

### [🔴] 14 处把请求体文本原样拼进 SQL，可盲注可拖垮连接池 (parity:true)
- 位置：logistics/bill.go:69、392、530、545；logistics/misc.go:99、794、796；logistics/waybill.go:56、59、62、87、90、92、112
- 证据：
```go
where := "WHERE " + billListQuery() + searchStrOf(m) + "%'"                            // bill.go:69
sqlStr += " and carriage_platform = " + str(m, "carriagePlatform")                      // waybill.go:87
cond += " and (order_id is null or order_id='' or order_id=" + str(m, "orderId") + ")"  // misc.go:794
```
- 场景：任意登录账号调 `chainLogisticsBill/listByPage`，`searchStr` 传 `x' AND (SELECT level FROM chain_auth_user LIMIT 1)=2 AND '%'='%`，count 与列表两条 SQL 都合法，看返回的 total 即可布尔盲注当前库任意表；传 `' AND SLEEP(5) AND '` 则每行睡 5 秒，几条并发请求占满 20 连接池。bill.go:392/530/545 在 id、contractId 为 JSON 字符串时同样可注入。
- 修法：一律改 `like ?` / `= ?` 参数绑定，正常输入行为不变。waybill.go:52 注释写明「注入语义保留」，docs 3.2 例外清单未列，需产品决策。
- 把握：high（注释逐条对照 Java `${}` mapper，Java 同样存在）。

### [🔴] 按 id 读写与整表导出无租户校验 (parity:true)
- 位置：logistics/waybill.go:145–274、245；logistics/bill.go:82–324、202；logistics/misc.go:111–215、306–520、817–900、492、553；logistics/excel.go:312
- 证据：
```go
rows, err := s.pick(r).Query("SELECT * FROM " + table)     // excel.go:312
"file_type": "1", "owner_type": "1", "owner_id": strconv.FormatInt(u.UserID, 10),   // waybill.go:245–247
"relation_id": json.Number(strconv.FormatInt(id, 10)), "relation_type": 1, "effective": true,
```
- 场景：A 公司普通账号带 B 公司 id 调任一 getById/updateById/deleteById/delBatch 即可读改删；exportExcel 得到全平台运单/过磅单；updateById 的 filelist 可把任意 chain_common_file 行改挂到自己名下；订单 add 的 ids 可绑定任意合同；twClearData 在调用方无企业行时不带 credit_code 条件（misc.go:553），会处理全平台未清洗运单。除 twListByPage 按 JWT 信用代码过滤外均如此；注释显示 Java 为 MP getById/removeById，同样无公司条件。
- 修法：产品决策；收紧时复用 `u.CompanyID` 给 by-id 查询与导出加 `company_id` 条件（对接运单表用 credit_code）。

### [🟡] 稀疏列补齐仍可把内存放大上万倍，几 KB 文件即可 OOM 整个进程 (parity:false)
- 位置：logistics/excel.go:1246–1251（上限见 1023–1026）
- 证据：
```go
if ci, err := colIndex(ref); err == nil && ci >= 0 && ci < 16384 {
	for len(row.cells) < ci {
		row.cells = append(row.cells, xlsxCell{kind: "blank"})
	}
}
```
- 场景：工作表每行写 `<row r="1"><c r="A1"/><c r="XFD1"/></row>`（41 字节）即补 16384 个 40 字节的 xlsxCell，约 640KiB。readZip 允许单表 64MB，约 410KB 的高重复 XML（压缩后几 KB）对应 1 万行约 6GB，单二进制被 OOM kill、全模块不可用，登录账号即可触发。POI 只保存 XML 中存在的单元格，Java 无此放大。
- 修法：给 xlsxCell 加列号字段，readExcelFile 的两个循环（682、695 行）按单元格自带列号取值，不再补空。
- 把握：high；只有 main.go 存在全局内存上限才能否定，实查没有。

### [🟡] 过磅单导入遇「开票时间」必 panic；接口写入永远丢 invoicingTime (parity:false)
- 位置：logistics/excel.go:195、884、473；logistics/logistics.go:178–187；logistics/misc.go:51；platform/dbutil.go:18、110
- 证据：
```go
{"开票时间", "Invoicing_time", "t", false, false, false},   // excel.go:195
ent[fl.col] = t                                             // excel.go:884
var colRe = regexp.MustCompile(`^[a-z_][a-zA-Z0-9_]*$`)     // dbutil.go:18
```
- 场景：wlImportExcel 中任一行开票时间是合法时间，行键即 `Invoicing_time`，`platform.Insert` 里的 `CheckCols` 因首字母大写而 panic；全仓只有 jobs.go:174 有 recover，net/http 打栈并断连，该行之前的行已提交，重试即重复。同一原因下 `pick()` 经 Underscore 只产小写键，`entWl.cols` 里的 `Invoicing_time` 永远匹配不上，wlAdd/wlUpdateById 传入的 invoicingTime 被静默丢弃（Java MP 生成 `invoicing_time`，列名不分大小写，可写）。
- 修法：`insertVals` 键统一 `strings.ToLower`；`pick` 对列名大小写不敏感匹配。
- 把握：high；commonrest/auth_meta_gen.go:261 同表写列已用小写 `invoicing_time`，可作旁证。

### [🟡] 共享字符串按 `<t>` 拆分，富文本单元格让整表串位 (parity:false)
- 位置：logistics/excel.go:1029–1057、1264–1268
- 证据：
```go
case xml.EndElement:
	if t.Name.Local == "t" {
		out = append(out, cur.String())
```
- 场景：Excel/WPS 保存的文件只要有一个富文本单元格（`<si><r><t>运</t></r><r><t>单</t></r></si>`，局部加粗或变色即产生），Go 记两条、POI 记一条，其后所有字符串单元格索引后移一位，车号、装货地整列读到相邻值并「导入成功」。内联富文本同理只留最后一段（`<t>` 起始处 `tBuf.Reset()`）。
- 修法：按 `<si>` 聚合（起始清空、结束追加一次）；内联 `<t>` 不重置 tBuf。
- 把握：high（Go 行为确定；POI XSSFRichTextString.getString 返回拼接文本）。

### [🟡] 删除物流订单会连带删除其他实体同 id 的附件记录 (parity:true)
- 位置：logistics/bill.go:298
- 证据：`_, _ = d.Exec("DELETE FROM chain_common_file WHERE relation_id = ?", id) // delByRelId`
- 场景：relation_id 是各实体共用号段，靠 relation_type 区分（字典 file_relation：1 合同 2 物流 3 结算 4 收付款 5 票据 6 企业；bill/billctl.go:259、capital/shares.go:27、warehouse/leasing_write.go:128 都写 1）。删除订单 #123，合同 #123、票据 #123 的附件行一并消失，物理文件成孤儿。其余包同样复刻 delByRelId，判为 Java 原行为。
- 修法：产品决策；若修，条件补 `relation_type`（需先核对物流应写的类型值，当前 updateFilelist 写的是 1 即「合同」）。

### [🟡] 导入端点无上传上限、无读超时 (parity:unknown)
- 位置：logistics/excel.go:361–367；main.go:60
- 证据：
```go
_, fh, err := r.FormFile("file")   // 5 个 importExcel 路由不在 body_contracts_gen.go，body.go:31 的 32MB 封顶不生效
```
- 场景：登录账号对 importExcel 持续推超大 multipart，每请求 32MB 驻留内存，其余落系统临时目录，时长不限（`http.Server` 只有 ReadHeaderTimeout），并发即占满磁盘或内存。临时文件在请求结束时由 net/http 清理，不算泄漏。
- 修法：`formFile()` 开头加 `r.Body = http.MaxBytesReader(w, r.Body, 32<<20)`。

### [🟡] updTransLocation 类型断言写错，功能恒返回 0 条 (parity:false)
- 位置：logistics/waybill.go:406、410、413、443–444
- 证据：
```go
cnt, _ := res["count"].(float64)   // 406；mapSearch 在 444 行 dec.UseNumber()，count 实为 json.Number
```
- 场景：腾讯地图正常返回 `count:10` 也被判 0，接口永远「更新运单数据0条」。修好后仍有隐患：`res["data"].([]any)`、`ad["district"].(string)`、`data[0]` 未校验，第三方返回异常即 panic；`maps` 只读不写（366、381 行），去重无效，每个地址组合 2 次外呼加 500ms 且扫全表；外呼不带 `r.Context()`，客户端断开后仍跑完。
- 修法：`num(res["count"]) > 0`；断言改逗号 ok 形式；`maps[addr] = true`；`NewRequestWithContext`。
- 把握：high。硬编码地图 Key（waybill.go:345）与 Java 同源，未计入。

### [🟡] 导入与批量写逐行提交无事务，中途失败留下半截数据 (parity:false)
- 位置：logistics/excel.go:466–493、522–529、550–557、575–590；logistics/misc.go:382–398、482–514、600–729；logistics/waybill.go:296–329
- 证据：
```go
if err != nil {
	platform.Unknown(w)   // 前 k-1 行已 autocommit 落库
	return
}
```
- 场景：第 k 行备注超长触发 1406，前 k-1 行已入库、响应 999，重试整批重复。Java 侧 saveBatch/updateBatchById 自带 MP 的 @Transactional 会整体回滚（据注释推断，把握 medium）。twClearData 更重：金额判空（misc.go:703–709）在订单/运单已写、`clean_flag=1`（691）之后才做，999 后重跑跳过这些行，车数/运费/吨数合计（721–729）永不回填，此顺序 Java 相同，parity:true。另 bill.go:195、202、208、230、295–298 与 waybill.go:245 的 `_, _ =` 吞错，关联失败仍回显成功。
- 修法：循环包进 `platform.WithTx`，改用 TxInsert/TxUpdateByID。

### [🟢] twAddBatch 缺 BatchNumber 头时写空串，重复校验失效 (parity:false)
- 位置：logistics/misc.go:354–356、361、391
- 证据：`q += "batch_number IS NULL"`（查重）对 `e["batchNumber"] = batchNumber`（空串入库）；361 行该分支查重出错也被吞。
- 场景：无 BatchNumber 头的批次重复提交永不报「批次号重复」。修法：batchNumber 为空时置 nil（autoTimeVals 会丢弃 nil，落列默认 NULL）。

### [🟢] twUpdateBatch 用未转义 LIKE 定位并吞查询错误 (parity:unknown)
- 位置：logistics/misc.go:492
- 证据：`_ = d.QueryRow("SELECT id FROM chain_logistics_third_waybill WHERE bill_number LIKE ?", str(e, "billNumber")).Scan(&id)`
- 场景：billNumber 含 `%` 即命中任意首行并被整行覆盖（无 credit_code 条件）；查询真出错时当作不存在而插入重复行。修法：`= ?` 并区分 `sql.ErrNoRows`。

### [🟢] 导入导出的三处边界 (parity:false / unknown)
- excel.go:878–881：`ParseFloat` 接受 NaN/Inf/Infinity，`json.Number("NaN")` 进 DECIMAL 列报 1292，整次导入 999（Java 按行吞掉、字段留空）。修法：加 `math.IsNaN/IsInf` 判断。
- excel.go:285–288：`xmlEscape` 不处理 XML 非法控制字符，备注含一个即整个导出损坏，且响应头已发出无法报错。修法：写入前过滤非法 rune。
- misc.go:639 `lobs[0]["id"].(int64)` 无 ok：驱动对 BIGINT UNSIGNED 返回 uint64（mysql@v1.10.1 packets.go:883–888），若该列为 UNSIGNED 则订单已存在时必 panic；仓库无 DDL，未能确认列型。misc.go:716–724 用 float64 累加后 `dec2` 强制两位小数，列小数位大于 2 时吨数被截断且为二进制舍入（2.675→2.67）。

## 重复模式
- SQL 拼接：14 个行号见第 1 条。
- 无租户校验：五个控制器的 getById/updateById/deleteById/delBatch/list/export 全部如此，见第 2 条。
- `r.FormFile` 无上限：本包 excel.go:362，另有 bill/excel.go:140、capital/excel.go:144、warehouse/excel.go:207、commonrest/excel.go:207（他组范围）。
- 逐行写无事务：见第 9 条。
- 系统编号 `max+1` 无锁（logistics.go:377–395 及三个调用点）：并发新增或导入会拿到相同编号，Java 同款，parity:true。

## 已核对无问题
- 分页：pageOf 默认 1/10 且仅 `num()>0` 时覆盖，size==0 除零不可达；LIMIT/OFFSET 溢出只报 999。
- rows：包内 13 处 Query 全部交给 ScanMaps，其余为 QueryRow+Scan，无提前返回泄漏；包内不用 WithTx，无事务或连接泄漏面。
- 请求键进 CheckCols 的两条路径都有白名单：allEqCond（列白名单）与 twClearData 的 validCol（大小写不敏感等值），后者拼接不可注入；bill.go:202 的 `cid.String()` 只可能是 JSON 数字字面。
- zip：总量 256MB、单条目 64MB 封顶，伪造声明大小会被单条目检查拦下。
- 导出为 inlineStr 字符串单元格（excel.go:229），`= + - @` 开头不会被当作公式。
- zz hash/spec：javaHash/javaMapOrder 只决定表头包含匹配的命中次序和必填错误信息的输出序，输入是 45 个以内的固定标题，无 panic 面；碰撞只影响同标题多命中时先匹配哪个。zz_spec_test、zz_hash_test 把 8765 实测序（如「创建时间」→trans_time）固化且断言具体列、可失败；excel_columns_test 断言 colIndex 三个值，有效；TestWlOrderDump 仅打印、不断言，无害。