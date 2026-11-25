徐先生，评审完成（全程只读）。按上级约束不用表情符号，严重度用文字标注，对应简报四级。

## 范围
AP/AddPurchaseContract.vue（至 2229 行）、EP/EditPurchaseContract.vue（1957）、PC/PurchaseContract.vue（1746）、AS/AddSalesContract.vue（2191），四个文件的 template+script 全读，`<style>` 跳过；均在 /Volumes/code/GoProject/panda-supply-chain/frontend/src/views/Mindex/。另核对 api/*.js、filters.js、http.js、router，以及 backend/internal 下 contract/{contract,stats,esign_file}.go、capital/payment.go、authadmin/common.go 的对应段落。

## 发现

### [阻断] 付款“确定”无防重，后端每次调用都向银行下单
- 位置：PC:691、1249-1265；capital/payment.go:65-66、113
- 证据：`this.api.confirmPay.payment(this.AddData)`，按钮无 loading；后端每次生成新 `DDXX`+毫秒订单号并插一条 transinfo。
- 场景：银行接口慢时双击 → 两笔担保支付订单、两条流水。修法：按钮加 `:loading`，回调复位。把握：high。

### [重要] 列表查询参数与后端语义不符
- 位置：PC:1320-1323（初值 838、870 为 ''）；contract/stats.go:102-117、203-207
- 证据：`isEsign:this.signValue`；Go 端 `v != nil` 即拼 `" and c.is_esign is " + boolLiteral(v)`，而 `boolLiteral("")` 得 `false`。
- 场景：未选筛选时发的是空串，默认列表恒带 `is_esign is false and is_pay_online is false`，开启过签章或在线支付的合同默认不可见；清空日期后 `endTime:""` 同理查空。另外“签署状态”发的是 `isEsign`，表格却按 `isEsignComplete` 显示，筛“已完成”结果错。
- 修法：Go 端这几项加 `v != ""`；前端改发 `isEsignComplete`（遗留语义，改前确认）。把握：high；Java 把空串当未传是推断，其他列表页大概率同病。

### [重要] 编辑页离开不清 queryAll，“新建”变成覆盖保存
- 位置：EP:997、1953-1956（无 destroyed）；AP:1040-1043、1547-1550；AS:987-989、1495-1498
- 证据：`if(sessionStorage.getItem('queryAll')){ this.AddData.updateTime = ''; this.api.cgContract.edit(this.AddData)`
- 场景：点“编辑”后经菜单离开，残留的 queryAll 被新建采购/销售合同页读到，表单填满旧合同，保存走 edit（cg/xs 同一个 updateById）覆盖原合同。修法：EP 补 `destroyed(){ sessionStorage.removeItem('queryAll'); }`。把握：high。

### [重要] 签署方按 id 分组，删除后 id 重复丢人
- 位置：AP:1205-1241、1481-1519（EP、AS 同）
- 证据：`removeDomain(item) { this.numbers--; ...`，`addDomain() { this.numbers++; ... id:this.numbers`，提交按 `i === domains[key].id` 用 `Object.assign` 合并。
- 场景：加两个签署方、删第一个、再加一个 → 两组 id 相同，提交推出一个空签署方、两人合并成一个；草稿回显（1052-1079）所有签署方共用同一 id。修法：id 只增不减，循环跳过空对象。把握：high。

### [重要] 附件列表与界面脱节
- 位置：AP:1904-1917、1889-1897、1540-1543（AS 同）
- 证据：`fileRemove(file,fileList){ if(this.fileObj[file.name]){ ...`；`resfile` 仅 `code==200` 才更新列表；保存取 `itme.response.data[0].id`。
- 场景：草稿回显的附件删除后界面消失，但 `enclosurefileList` 未更新，保存仍带上。上传业务失败（或文档转 pdf 失败时后端返回 code 200 + 空 data，common.go:725-805）时文件仍显示已上传，保存却提示“请上传合同文件”或抛 `reading 'id'`。
- 修法：回显文件也走删除；以 `code==200 && data.length` 判成功。把握：high。

### [重要] 请求失败后 loading 不复位
- 位置：AP:1549-1587（AS:1497-1535 同）；PC:1313-1341
- 证据：`this.fullscreenLoading = true; this.api.cgContract.edit(this.AddData).then(res=>{`（无 catch）
- 场景：`platform.Unknown` 返回 HTTP 500，reject 后全屏锁永久锁死，只能刷新丢表单；列表页 `loading` 只在 code==200 复位。修法：补 `.catch` 复位。把握：high。

### [重要] 草稿路径两处静默失败
- 位置：AP:1792、1836、1317-1334（AS 同）；contract.go:427-444
- 证据：`this.AddData.eqbFlow.contractName = this.enclosurefileList[0].name;`
- 场景：开启签章、有签署方但未传附件，存草稿时 `[0]` 为 undefined 抛错，草稿没存也没提示，且无防重。打开草稿后焦点进出“合同编号”，`cgContract.list` 不按 fileStatus 过滤、查到自己，编号被清空。
- 修法：取 name 前判空并加 loading；有 queryAll 时跳过查重。把握：high。

### [轻微] 小问题合集
- PC:1537-1602：关闭抽屉或点“展示图”即 `cgContract.edit` 整体回写并记“修改合同”日志，已签章合同还会把附件重传签章服务（contract.go:766）。建议仅在补充协议增删后调用。
- PC:1354-1396：`search()` 先请求再 `current = 1`，带旧页码多发一次并竞态；删草稿（1093）只改页码不刷新。把握 medium。
- PC:1695-1702：导入进度 `setInterval` 离开未清；el-upload 未绑 on-error，上传失败停在 99%。
- AP:1532-1536（EP:1445-1449）：提示“签署方信息不能少于一个签署方”后无 return，仍继续提交。
- 三页 `clickTime`（AP:1293）清空签定时间后 `null.substr` 报错；总价 `((unitPrice*1)*(contractNum*1)).toFixed(2)`（AP:1288）半分边界差 1 分，如 1.15×0.5 得 0.57。

## 重复模式
- 下载/预览 URL 硬编码 `Signature`/`AuthorizationId`（值不在此复述）：PC:588、1120、1132、1137，AP:1900，EP:1625，AS:1849，与 http.js 同源。
- 签署方拼装、提交、草稿三段代码在 AP/EP/AS 各复制 2-3 份，修改需同步；`queryAll` 被 30 个页面共用，其他 Edit* 页是否清理建议一并核对。

## 已核对无问题
filters 全部已定义；api 方法均存在且均为 POST JSON；无生效的 v-html/eval，iframe 与 window.open 只用后端返回 URL；合同 id 自增，无 2^53 问题；`getById` 的 `singers` 空时为 `[]`；`updateById` 按原值整体回写，普通编辑不丢未展示字段。

[亮点] AP/AS 发请求前先置 `fullscreenLoading`，挡住了双击重复建合同。