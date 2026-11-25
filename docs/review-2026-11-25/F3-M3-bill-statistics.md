徐先生，票据与缴税 9 个页面的评审如下（只读，未改任何文件）。

## 范围
9 个 .vue 的 template 与 script 全部读完（约 8500 行），`<style>` 跳过。另核对了对应 api 模块、filters.js，以及后端 internal/bill、platform、commonrest、contract/stats2.go 的相关端点。

## 发现

### [🟡] 编辑页离开不清 queryAll，残留被「新建合同」页当编辑数据
- 位置：EditInBillStatistics.vue:826、EditSalesStatistics.vue:798、EditCoststatistics.vue:290
- 证据：只有 `Addopen()` 和保存成功的 onClose 执行 `sessionStorage.removeItem('queryAll')`
- 场景：发票点「编辑」后用后退或侧边栏离开，再开「新建采购合同」。AddPurchaseContract.vue:1040-1043 见 queryAll 即 `this.AddData = JSON.parse(...)` 并读 `this.AddData.singers.length`，mounted 抛 TypeError，表单灌入发票字段；保存走 `cgContract.edit`（:1550），同 id 的合同可能被覆盖。
- 修法：三个编辑页加 `beforeRouteLeave(){ sessionStorage.removeItem('queryAll') }`
- 把握：high；覆盖错合同为 medium，取决于两表 id 是否撞号。

### [🟡] 编辑进项发票选供应商写错字段，并连带保存失败
- 位置：EditInBillStatistics.vue:1128
- 证据：`this.AddData.supplierName = this.dialogtextarea;`（Add 页 :1115、销项编辑页 :1114 都写 `companyName`）
- 场景：应用后「销方企业名称」不变。`supplierName` 随 payload 提交，后端 writeVals（billctl.go:330）不剔除非实体键，拼出 `UPDATE chain_bill SET supplier_name=?`，该表无此列，1054 后 Unknown() 回 HTTP 500，保存必败。Java 时代未知字段被忽略，属迁移后恶化。
- 修法：改为 `this.AddData.companyName`
- 把握：字段写错 high；1054 为 medium，需确认物理表确无该列。

### [🟡] 销项发票关联销售合同，购方被预填成本公司
- 位置：AddSalesStatistics.vue:908、EditSalesStatistics.vue:931
- 证据：`this.AddData.companyName = this.placecontent.salesUnit // 销售合同的客户`
- 场景：AddSalesContract.vue:1114 把销售合同的 `salesUnit` 设为本企业名，购货单位存在 `buyerUnit`（:1951）。关联后购方是自己公司，且能通过必填校验直接保存。
- 修法：改取 `placecontent.buyerUnit`。把握：high。

### [🟡] 缴税联动漏算：进项税额为 0 时不重算
- 位置：AddCoststatistics.vue:304-324、EditCoststatistics.vue:293-313
- 证据：`if(val.jxMoney != 0){ this.AddData.zzMoney = ...`
- 场景：进项税额保持默认 0、只填销项含税金额，增值税额和四项附加税全停在 0，缴税合计只含所得税和滞纳金，required 对 0 放行，可直接保存。把进项税额改回 `0` 时 `'0' != 0` 为假，增值税额停在上一次结果。另外编辑页 deep watch（:316-320）一打开就用公式覆盖库里的 sumMoney；页面规则写印花税 0.03，代码是 0.0003（:316）。
- 修法：联动抽成一个函数，两个 input 都调用，去掉 `!= 0`。把握：high。

### [🟡] 先发请求再重置页码
- 位置：Coststatistics.vue:558-562，其余见重复模式
- 证据：`this.getpurcon(); this.current = 1;`
- 场景：停在第 2 页时搜索，请求仍带 current=2，筛选结果只有 1 页则返回空，列表显示暂无数据而分页停在 1；删掉末页最后几条同理。
- 修法：先置 1 再请求。把握：high。

### [🟡] 请求失败后 loading 永不复位
- 位置：Coststatistics.vue:545-554；InBillStatistics.vue:763-774；SalesStatistics.vue:748-759
- 证据：`.then(res=>{ if(res.data.code == 200){ this.loading = false;`
- 场景：后端 Unknown() 直接回 HTTP 500（resp.go:48-52）被 reject，或 200 加 code≠200（Coststatistics），遮罩永久停在「加载中」，只能刷新。
- 修法：then 首行复位，并补 `.catch(() => this.loading = false)`。把握：high。

### [🟡] 新建页保存无防重
- 位置：AddInBillStatistics.vue:778-791，其余见重复模式
- 证据：`duration:500, onClose(){ that.$router.push(...) }`
- 场景：成功后 500ms 才跳转，按钮始终可点，双击会建出两条同 systemCode 的发票（billAdd 无唯一约束）。
- 修法：`saving` 标志加按钮 `:loading`。把握：high。

### [🟡] 导入：校验挂错钩子、失败无提示、定时器泄漏
- 位置：InBillStatistics.vue:487-497、:795-810、:836-843，其余见重复模式
- 证据：`:on-change="supbeforedemoUpload"`；`let f = setInterval(...)`
- 场景：el-upload 不使用 on-change 的返回值，importnav 又在校验前赋值，选 .exe 仍可点「导入」直传。没有 on-error，HTTP 失败时停在 99% 的「正在导入...」且无提示。进度 5 秒内离开页面，f 永不清除。
- 修法：校验移到 `:before-upload`，补 `:on-error`，定时器存到 this 并在 beforeDestroy 清。把握：high。

### [🟢] 附件列表与提交 payload 脱节
- 位置：AddInBillStatistics.vue:766-768、:841-869，其余见重复模式
- 证据：`arr.push({id:itme.response.data[0].id});`
- 场景：上传返回 code≠200 的文件仍留在 el-upload 列表，下一次成功上传会把它带进 enclosurefileList，保存时抛 TypeError，点保存无反应。`addfile.del` 失败时不更新列表，界面已删但 id 照样提交。编辑页删附件是立即物理删除，随后点「取消」附件也已丢。
- 修法：提交前过滤没有 `response.data` 的项，删除成功后再移出列表。把握：medium。

### [🟢] 弹窗残留状态
- 位置：Coststatistics.vue:667-677 等三页
- 证据：ImporthandleClose 不复位 `fourShow`
- 场景：导入成功后点 X 关闭再打开，步骤三结果面板与步骤一同屏。合同弹窗关闭不清 `placecontent`，重开直接点「应用」会套用上次的合同。
- 修法：关闭时一并复位。把握：high。

## 重复模式
- 先请求后置页码：Coststatistics.vue:561、:758；InBillStatistics.vue:881、:941、:947、:1062；SalesStatistics.vue:866、:926、:932、:1047
- 保存无防重：AddSalesStatistics.vue:744-757、AddCoststatistics.vue:283-296、各页 CompanyNewsubmitForm
- 导入：SalesStatistics.vue:473-483、:780-795、:821-828；Coststatistics.vue:350-360、:646-661、:687-694；弹窗残留在 InBillStatistics.vue:816-826、SalesStatistics.vue:801-811
- 附件：EditInBillStatistics.vue:852-882；AddSalesStatistics.vue:807-835；EditSalesStatistics.vue:824-854
- 下载和预览 URL 拼接写死的静态 Signature、AuthorizationId，fileName 未编码：InBillStatistics.vue:502、:716、:729、:734；SalesStatistics.vue:488、:701、:714、:719；Coststatistics.vue:365；AddInBillStatistics.vue:852、EditInBillStatistics.vue:863、AddSalesStatistics.vue:818、EditSalesStatistics.vue:835。属全仓模式，建议 lead 统一处理。

## 已核对无问题
- SalesStatistics.vue:948 用 `bill.one` 而非 `xxbill.one`：chainBill、chainCgBill、chainXsBill 共用 handleBill 和 chain_bill 表（bill.go:51-57），两者等价。
- 编辑页不带 contractId 时，后端只写 body 里出现的键，合同关联不会丢；getSments 空集返回 `[]`；searchFour 的 `contract_id` 键名与前端一致；currency 过滤器对 null 安全。
- 税额输入产生 NaN（如只输入 `.`）会被 pattern 规则拦在保存前；toFixed 半分舍入未见可静态证明的错误（本次未运行验证），未单列。
- 三个新建页不读 queryAll；无 v-html 等 XSS sink；本范围没有明细行、echarts 和事件监听；采购结算弹窗是注释掉的死代码，未计。