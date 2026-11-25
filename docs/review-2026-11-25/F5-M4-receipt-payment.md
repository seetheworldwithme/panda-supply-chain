徐先生，Mindex 资金页评审报告如下（位置均在 /Volumes/code/GoProject/panda-supply-chain/frontend/src/views/Mindex/）。

## 范围
8 个 .vue 的 template 与 script 全部读完（9,325 行，`<style>` 跳过）。另读 http.js、filters.js、9 个 api 模块，以及后端 /Volumes/code/GoProject/panda-supply-chain/backend/internal/capital/detailed.go、balance.go 校准严重度。node_modules 未安装，凭记忆判断 Element 行为处已在「把握」注明。

## 发现

### [🔴] 新建付款/收款单「保存」无防重，双击建两笔同额单据
- 位置：AddPaymentSlip.vue:8,654；AddReceipt.vue:8,670
- 证据：`<el-button size="small" @click="AddDatasubmitForm('AddData')">保存</el-button>` → `this.api.cgchainCap.add(this.AddData)`
- 场景：双击保存触发两次 add；后端 dtAdd 不查重 system_code，生成两条同额同编码单据，可核销金额翻倍。
- 修法：saving 标志 + `:loading`，请求返回后复位。
- 把握：high；库上 system_code 有唯一索引则降为重复报错。

### [🟡] 核销提交无防重，重复核销
- 位置：Receivable.vue:30,1058；Copewith.vue:30,1051
- 证据：`this.api.chainCap.writeoff(DataObj).then(res=>{ if(res.data.code == 200){`（前后无提交标志）
- 场景：双击「核 销」排队弹出两个确认框，各确定一次即发两次同 payload；后端按请求体里的 khxje 快照重算，两次都成功，核销记录多一条、detail_settlement.money 累加两次。慢请求期间再点同理。
- 修法：提交标志 + `:loading`，成功失败都复位。
- 把握：high；两个确认框依赖 MessageBox 排队，未核源码，慢请求路径不受影响。

### [🟡] 核销金额浮点连加后用 === 比较，合法核销被拦
- 位置：Receivable.vue:1014；Copewith.vue:1013
- 证据：`if(Math.abs(this.settlementmoeny) === Math.abs(this.contmoeny)){`
- 场景：结算 1000.02 加 2000.01 对收款单 3000.03，双精度和为 3000.0299999999997，提示「实际核销金额不一致」，怎么改都过不了，只能拆两次核销。
- 修法：`Math.abs(a - b) < 0.005`。
- 把握：high，已按 IEEE-754 逐位推算。

### [🟡] 应付单核销缺零金额拦截
- 位置：Copewith.vue:1013-1022（Receivable.vue:1015-1019 有）
- 证据：相等判断后直接 `let res; if(this.settlementmoeny){`
- 场景：只选供应商不勾行点「核 销」，确认「核销金额共计：0 元」后发出 details/sments 皆空的请求；后端仍生成 HXJL 编号并插入 money 为空的核销记录，或入库报错。
- 修法：补「核销金额不能为零」分支。
- 把握：high。

### [🟡] 列表加载：先请求后重置页码、空结果不清旧数据、reject 后 loading 不复位
- 位置：PaymentSlip.vue:654-678,834-844,956-958（Receipt.vue 同构）；Receivable.vue:1105-1160,1294-1303,1559-1560（Copewith.vue 同构）
- 证据：`search() { this.tableData = []; this.getpurcon(); this.current = 1; }`；`this.paginationshow = false;`
- 场景：停在第 2 页搜索，请求仍带 current=2，结果不足 2 页就显示「暂无数据」且 total 是旧值；核销页 paginationshow 置 false 后永不复位，之后翻不了页；删除末页、切换客户同理。请求 reject 时 `this.loading = true` 无 catch，表格一直「加载中」。
- 修法：先置 current=1 再请求；空分支清 tableData/total；补 catch。
- 把握：high；程序改页码不触发 current-change 凭 Element 2.15 行为，未核源码。

### [🟡] 编辑页 queryAll 清理不全，后退键重进时崩溃
- 位置：EditPaymentSlip.vue:611,689,700；EditReceipt.vue:627,705,716
- 证据：`this.AddData = JSON.parse(sessionStorage.getItem('queryAll'));`
- 场景：只在取消、保存成功时 removeItem。之后按浏览器后退回到编辑路由（无 keep-alive，组件重建），AddData 为 null，render 抛 TypeError，页面空白。经左侧菜单离开则残留，被 AddSalesContract.vue:987、AddPurchaseContract.vue:1040、AddWaybill.vue:1713 无条件读取（不在本次范围，供关联）。
- 修法：mounted 取不到值即 `$router.replace` 回列表；beforeDestroy 里 removeItem。
- 把握：high。

### [🟡] 附件未传完或上传失败仍可保存，附件静默丢失
- 位置：AddPaymentSlip.vue:193-202,641-645,717-745（其余三个 Add/Edit 页同构）
- 证据：`resfile(response, file, fileList){ if(response.code == 200){ ...; this.enclosurefileList = fileList; } }`
- 场景：选完附件立刻保存，上传未完成，filelist 不含该文件；上传 HTTP 200 但 code≠200 时列表仍显示成功、保存却不带；fileRemove 删除接口业务失败时不更新列表，已移除附件下次又被关联。
- 修法：保存前检查 uploadFiles 有无非 success 项；resfile 对非 200 提示并移除。
- 把握：high。

### [🟡] 编辑保存回写整条快照，覆盖并发核销结果
- 位置：EditPaymentSlip.vue:611,680；EditReceipt.vue:627,696
- 证据：`this.api.cgchainCap.edit(this.AddData)`，AddData 即 getById 整行
- 场景：打开编辑时 khxje=400，保存前别处又核销 300；后端 dtUpdateById 取实体全列，把快照的 khxje=400 与旧 hx_code 写回，可核销金额被恢复。
- 修法：只提交可编辑字段白名单。
- 把握：medium-high；若后端 pick 跳过未变列则不成立。

### [🟡] 静态签名与授权 ID 硬编码进 URL
- 位置：PaymentSlip.vue:410,620,633,638；Receipt.vue:412,623,636,641；Add/Edit 四页 openfile（如 AddPaymentSlip.vue:728）；Receivable.vue:289,682；Copewith.vue:288,681（兜底 http.js:13-14 归 lead）
- 证据：`window.open(url)`，url 为固定 Signature/AuthorizationId 参数拼 fileName
- 场景：产物中任何人可取到凭据，免登录按 fileName 调 downById/previewStream 取文件；凭据随 URL 进历史与日志。
- 修法：改走带登录态请求头的接口，不把签名拼进 URL。
- 把握：硬编码 high；可利用性 medium，未核后端校验。

### [🟢] 低风险合并项
- 位置：Add/Edit 四页账户弹窗（AddPaymentSlip.vue:244,754）及 Receivable.vue:78,1475；PaymentSlip.vue:736-743,619-628（Receipt.vue 同构）；AddPaymentSlip.vue:1015-1042（四页同构）；Receivable.vue:1167,1311
- 证据：`<el-radio :label="scope.$index" v-model="radio" ...>`；`this.src.promise.then(...).catch(() => {})`
- 场景：弹窗翻页不清 radio，新页同下标行呈选中，「应用」带出上一页那行；Add/Edit 页 search() 误写 `this.contracttotal = 1`；导入无 on-error，失败停在 99%，切路由后 setInterval 访问已销毁的 `$refs.upload`；PDF 预览失败后一直转圈；「结算关联」入口已注释，PlaceChoice* 不可达；detailsPur 定义两次后者覆盖前者。
- 修法：翻页时清 radio；加 on-error 与 beforeDestroy；删死代码。
- 把握：high。

## 重复模式
- 「先请求后重置页码」「无 catch」「保存无 loading」在 4 个列表页与 4 个编辑页成对出现，宜各抽一个加载函数和提交守卫。

## 已核对无问题
- api 调用与模块定义一致；Receipt.vue:863 的 chainCap.one 与 xschainCap.one 后端同一 handler。
- 🎉 核销输入 blur 校验 `val.khxje >= x*1 > 0` 加正则，超额、0、负数均被拦；列表恒返回非空 khxje/sjhxje，payload 不带 null 金额。
- 新建页不读 queryAll；成功提示均判 code==200；无 v-html/eval。