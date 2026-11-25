徐先生，评审完成。

## 范围
11 个文件的 template+script 全部读完（style 跳过），并对照 `frontend/src/api/*`、`path.js`、`backend/internal/contract`。下文 `M/` = `frontend/src/views/Mindex/`。

## 发现

### [🔴] 编辑页离开不清 queryAll，「新建」页读到后按 edit 覆盖原合同
- 位置：M/EditSalesContract.vue:1004、1476/1578、1589-1592
- 证据：`this.AddData = JSON.parse(sessionStorage.getItem('queryAll'))`；只有成功 onClose 和「取消」removeItem，没有 beforeRouteLeave。
- 场景：进编辑页后点菜单或后退离开，再去新建：AddSalesContract.vue:987-989 读到残留，1495-1498 走 `xsContract.edit`，原合同被覆盖。残留为 null 时 1005 行 `singers.length` 抛错。
- 修法：`beforeRouteLeave` 里 removeItem；1004 行后判空回列表。
- 把握：high。

### [🔴] 「发起签署 / 保存」无防重，后端不去重，产生双份签章流程、双扣份数
- 位置：M/AddSignContract.vue:8、508；M/EditSalesContract.vue:8、1467/1569
- 证据：`<el-button @click="AddDatasubmitForm('AddData')" type="primary">发起签署</el-button>`，无 loading/disabled。
- 场景：双击发两次请求；contract.go:765-770 → contractSign → esign.go:153 每次 insert 新 flow+singers。编辑页内联「保存」（1700/1825/1925）可建重复客户、地址。
- 修法：提交标志 + `:loading`，finally 复位。
- 把握：high；Go 端 e签宝端点已实现（esign.go），可达。

### [🟡] 签署方用可回退序号当 id：删后再加重号，签署方被合并丢失并多出空 signer
- 位置：M/EditSalesContract.vue:1149-1185、1404-1440、1506-1541；M/AddSignContract.vue:356-423、459-497
- 证据：`this.numbers--` / `this.numbers++ … id:this.numbers` / `for(let i=1;i<=this.numbers;i++){ a = {} … arrays.push(a)`
- 场景：加 A(1)、B(2)，删 A（numbers=1），加 C（又是 id=2）：i=1 无匹配 push `{}`，i=2 把 B、C 合成一条，B 收不到签署。回显分支 1012-1041 多个签署方同用 `singers.length-1`，同样重号。
- 修法：id 只增不减，按 id 分组提交。
- 把握：high。

### [🟡] 抽屉关闭保存补充协议的逻辑永不执行，附件列表跨合同累积
- 位置：M/SalesContract.vue:1415、1288-1323、1093-1105
- 证据：`if(seeData.createUserId == userIds){`
- 场景：ReferenceError 使 edit 和 1426-1427 的清理都不执行；`allone` 不清 `enclosurefileList`，看过 A 再看 B，B 的上传列表出现 A 的补充协议且逐次重复，点删除会 `addfile.del` 删掉 A 的文件；新上传的补充协议从未关联。
- 修法：`allone` 起始清空该列表。不要只改成 `this.`：激活后每次关抽屉都整包 updateById，已签章合同会再触发 contractSign；应仅在有增删时保存。
- 把握：high。

### [🟡] 删除合同不判 code，失败也提示「删除成功」
- 位置：M/SalesContract.vue:1445-1450
- 证据：`.then(res=>{ this.$message({ type: 'success', message: '删除成功!');`
- 修法：加 `if(res.data.code == 200)`。
- 把握：high。

### [🟡] 请求失败后 loading 不复位
- 位置：M/SignContract.vue:316-338（业务失败即永久转圈）；M/SalesContract.vue:1146-1173；M/Relation.vue:1490-1522、1715-1996；M/RelationBill.vue:166-186；M/EditRelationBill.vue:175-213
- 证据：`this.loading = true; … if(res.data.code == 200){ this.loading = false;`
- 场景：HTTP 500 被 reject 或 code≠200，遮罩一直「加载中」。
- 修法：`.finally` 复位。
- 把握：high。

### [🟡] 先请求后重置页码；切换合同不重置页码
- 位置：M/SalesContract.vue:1186-1227、989-991、1451-1453、1566-1567；M/SignContract.vue:342-352；M/Relation.vue:1555-1566、1605-1631、1739-2022、2388-2426；M/EditSalesContract.vue:1681-1685、1777-1784、1903-1910；M/RelationBill.vue:235-263；M/EditRelationBill.vue:233-261；M/Journal.vue:132-142
- 证据：`this.getpurcon(); this.current = 1;`
- 场景：停在第 3 页搜索或删除，仍查第 3 页，返回空，显示「暂无数据」而 total 仍是旧值。Relation 换合同沿用旧页码，四个 tab 同理。
- 修法：先置 1 再请求。
- 把握：high。

### [🟡] 附件列表与界面脱节
- 位置：M/AddSignContract.vue:592-631、499；M/EditSalesContract.vue:1621-1651
- 证据：`this.AddData.field = arr[0].id`（删除时不复位）/ `this.enclosurefileList[0].name`
- 场景：
  - 上传后删除再「发起签署」：field 仍为真、列表已空，499 行 TypeError，点了没反应。
  - 上传 code≠200：界面仍显示文件，列表和 field 未更新。
  - 编辑页删除回显附件即调 `addfile.del`，再点「取消」，contractIds 指向已删文件。
  - del 失败时被删附件仍随保存提交。
- 修法：删除回调同步清 field；上传失败移除该文件项；回显附件等保存成功后再删。
- 把握：high（前端逻辑）；el-upload on-remove 先移除后回调是按 2.15 源码记忆，仓内无 node_modules 可核。

### [🟡] 四流合一入口失败后卡在「path=Relation」模式
- 位置：M/Relation.vue:1490-1520
- 证据：`this.multiplesettlement = res.data.data.records[0];`
- 场景：ids 查不到合同时 TypeError，1519-1520 清 path 不执行，之后从菜单进本页仍按旧 ids 查，合同表恒空。
- 修法：进分支即清 path/pathID，判 `records.length`。
- 把握：high；触发需 ids 无匹配。

### [🟡] 合同抽屉图片附件预览必失败
- 位置：M/Relation.vue:2122
- 证据：`this.api.file.preFIle({ file_id: this.oneData.id })`
- 场景：`oneData` 是 `{name,response}` 包装对象，id 为 undefined，一直「文件加载中」。
- 修法：同 SalesContract.vue:1007，用 `response.data[0].id`。
- 把握：high。

### [🟡] 导入进度定时器未随页面销毁
- 位置：M/SalesContract.vue:1522-1529
- 证据：`let f = setInterval(()=>{ … this.$refs.upload.submit(); clearInterval(f);`
- 场景：点「导入」后 5 秒内离开页面，到 99 时 `$refs.upload` 已不存在，抛错且 clearInterval 不执行，每 150ms 报错直到刷新。
- 修法：定时器存到 this，beforeDestroy 与 ImporthandleClose 清理，先 clear 再 submit。
- 把握：high。

### [🟡] 下载/预览链接拼接固定签名与 AuthorizationId
- 位置：M/EditSalesContract.vue:1632；M/SalesContract.vue:563、1017、1029、1034；M/AddSignContract.vue:333、613；M/Relation.vue:2132、2163
- 场景：与 http.js 未登录分支是同一组静态值，会落入历史和代理日志；`fileName` 未编码，含 `&` 或 `#` 即失效。
- 修法：改走带 Authorization 头的接口。全局问题，数值未复述。
- 把握：high；能否越权下载取决于后端校验。

### [🟢] 零散
- M/EditSalesContract.vue:1232 `(unitPrice*1)*(contractNum*1)).toFixed(2)`：0.15×0.3 得 0.04，应按分相乘再 round。
- 死代码：PaymentRecord/MakeCollections/Reconciliation（router/index.js:344-361 已注释，无入口）；SignContract.vue:366-439；Relation.vue:2097-2214、2316-2320（`JshandleClick` 误调合同 `allone`）；EditSalesContract.vue:1790、`this.splice`（1099/1315/1354）。
- Relation.vue:2396-2453 弹窗重开不清旧数据，空页后分页永久隐藏；EditSalesContract.vue:1675 同。

## 重复模式
「请求后置页码」「loading 在 code==200 内复位」「`records != ""` 为假时不清表、不更新 total」遍布 9 个列表页，适合抽一个统一的分页请求函数。

## 已核对无问题
- `cgContract.one` 与 `xsContract.one` 是同一 handler 的别名（contract.go:59-62），用它查销售合同不是 bug。
- `eqbSingers.all` 的 `contractId` 和 `sys_flow_id` 后端都支持（sub.go:474-487）。
- EditRelationBill 回显可用：EditWaybill.vue:1471 进页前已写 queryAll。
- 提交里的 `domains`/`domainsone` 等展示字段，后端按列白名单落库（contract.go:723），不致污染。
- 无 v-html/eval；无 echarts/addEventListener；上传后端有 100MB 上限。
- 🎉 Relation.vue 未选合同时给明确提示；RelationBill/EditRelationBill 用 `row-key` + `reserve-selection` 保留跨页勾选。