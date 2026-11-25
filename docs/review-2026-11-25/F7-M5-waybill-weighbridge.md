## 范围
8 个文件通读 template + script（仅跳过 `<style>`）；用到的 api 模块（logistics/waybill/Weighing/listByPage/addfile/file）调用均存在、参数一致。为校准严重度抽查了后端 `internal/logistics/bill.go`、`misc.go`。本地无 node_modules，依赖库内部行为未实证处标 medium。

## 发现

### [🔴] `queryAll` 残留且全仓共用，运单会被别的订单夺走
- 位置：AddWaybill.vue:1713-1738、1743；EditWaybill.vue:1361-1369、1667-1669；Poundrder.vue:630；EditPoundrder.vue:472
- 证据：`this.tagsone = JSON.parse(sessionStorage.getItem('queryAll'));`（mounted 无条件采纳）
- 场景：编辑订单 A → 关联 → 返回 → 经菜单/后退离开（仅取消/成功才清键，全仓无 keep-alive）→ 点新建 → 空表单带出 A 的运单和合计 → 保存 → 后端 `billChildren`（bill.go:206-231）不校验归属，直接 `UPDATE chain_logistics SET order_id=新单`，A 丢运单。同键另一形态：磅单编辑页把磅单行写进 `queryAll` 且不清，之后编辑网络货运订单时 `tagsone.length` 为 undefined，保存被拦。
- 修法：关联结果改专用键，两页读取后立即 removeItem 并加 `Array.isArray` 守卫；EditWaybill 提交改用 `this.tagsone`；EditPoundrder 退出时清键。
- 把握：high。

### [🔴] 保存无防重复提交
- 位置：AddWaybill.vue:8、1913；AddPoundrder.vue:8、492；两个编辑页同模式
- 证据：`<el-button size="small" @click="AddDatasubmitForm('AddData')">保存</el-button>`，无 loading/disabled
- 场景：双击发两次 POST，后端 `billAdd`/`wlAdd` 无幂等 → 重复订单/磅单，重复磅单抬高累计净重。
- 修法：`submitting` 标志 + 按钮 `:loading`，响应后复位。
- 把握：high。

### [🟡] 编辑订单回写整行运单快照，重关联后合计不重算
- 位置：EditWaybill.vue:1667-1674、1361-1398
- 证据：`this.AddData.logistics = JSON.parse(sessionStorage.getItem('queryAll'));`
- 场景：① 发的是整行运单对象（新增页只发 `{id}`），快照后运单被改过的新值会被后端 `pick(e, entLo.cols...)` 整行覆盖；② 走 `queryAll` 分支跳过 else 里的运单统计，重新关联后车数/运费/吨数仍是旧值且输入框 disabled，与 AddWaybill.vue:1719-1738 不一致。
- 修法：提交前 `.map(v => ({id: v.id}))`，删 `contracts`；合计移到 if/else 外对 `this.tagsone` 统一算。
- 把握：high；①的实际损害取决于并发修改。

### [🟡] 运单管理抽屉附件跨运单累积并回写
- 位置：WaybillManag.vue:830-849、621-634、890-904
- 证据：`this.enclosurefileList.push(newarr[i]);`（从不清空）
- 场景：查看 A → 基本信息 → 关闭 → 查看 B → 基本信息 → 列表为 A+B 的附件，点「动态/展示图」即 `closedraws` 回写。跨层备注：后端 `twUpdateById`（misc.go:441-460）不处理 filelist，附件暂未落库，但也意味着抽屉的「上传磅单文件」实际不生效，请对照 Java 原行为。
- 修法：`allone` 先清空列表；切 tab 不触发保存。
- 把握：high（前端部分）。

### [🟡] 查询失败后 loading 不复位
- 位置：Waybill.vue:647-668；WaybillManag.vue:748-775；Poundrder.vue:637-658；WaybillDetails.vue:363-386
- 证据：`this.loading = true;` … `if(res.data.code == 200){ this.loading = false;`
- 场景：HTTP 5xx 被 reject 或 code≠200 → 「加载中」遮罩永久停留。
- 修法：`.catch(() => {}).finally(() => { this.loading = false })`。
- 把握：high。

### [🟡] 先请求后重置页码；合同弹窗空结果后分页消失
- 位置：Waybill.vue:683-687、784-785；WaybillManag.vue:778-812；Poundrder.vue:661-665、832-833；AddWaybill.vue:2467-2470、2442-2460；EditWaybill.vue:1887-1890、1862-1879
- 证据：`this.getpurcon(); this.current = 1;`
- 场景：停在第 3 页搜索仅 2 条 → 请求仍带 current=3 → 空响应 → 显示「暂无数据」而分页停在 1；删除末页同理。合同弹窗一次空结果后 `contractpaginationshow=false` 不再恢复、旧行不清。
- 修法：先置页码再请求；空结果清表、保留分页。
- 把握：medium-high；若 Element UI 2.15 对程序化改页会补发 current-change，列表页部分不成立。

### [🟡] 附件上传业务失败被当成功，删除失败后界面与提交脱节
- 位置：AddWaybill.vue:2015-2043；EditWaybill.vue:1737-1767；Add/EditPoundrder、WaybillManag.vue:708-734 同模式
- 证据：`if(response.code == 200){ … this.enclosurefileList = fileList; }`
- 场景：el-upload 自带 XHR 不过 axios 拦截器，HTTP 200 + code≠200 时文件显示上传成功却未进列表 → 保存后附件缺失，点预览 `res.response.data[0]` 抛错；on-remove 删除接口失败时界面已移除、提交仍带该文件。
- 修法：else 分支提示 `response.msg` 并剔除该文件；删除失败则恢复列表。
- 把握：high。

### [🟡] 导入弹窗校验无效、错误无反馈、定时器不可取消
- 位置：WaybillManag.vue:402-412、953-960、986-1001；Poundrder.vue:382-392、762-769
- 证据：`let f = setInterval(()=>{ this.perstep += 3; if(this.perstep == 99){ … this.$refs.upload.submit(); clearInterval(f); } },150)`
- 场景：类型校验挂在 on-change 且返回值不拦截，选 .exe 也能点导入；无 on-error，HTTP 5xx 时停在「正在导入…」；关闭后 5 秒内重开再导入 → 两个定时器叠加，`perstep` 跳过 99，一个永不停；WaybillManag 导入完成不刷新列表。
- 修法：校验改 `before-upload`，加 `:on-error`，定时器挂 `this` 并在关闭/beforeDestroy 清。
- 把握：high。

### [🟡] 下载/预览 URL 拼接写死签名
- 位置：AddWaybill.vue:2026、Waybill.vue:613-631、WaybillManag.vue:417-670、Poundrder.vue:397-605 等 16 处
- 证据：`downById?Signature=<写死签名>&AuthorizationId=<写死ID>&fileName=' + res.response.data[0].fileName`
- 场景：静态签名进浏览器历史、代理日志、Referer；`fileName` 未编码。与 http.js 静态头同值，建议统一处置。
- 修法：改走带 token 的 blob 下载。把握：high。

### [🟢] 装货地弹窗单选框不清空（变量名写错）
- 位置：AddWaybill.vue:2582、2590；EditWaybill.vue:2002、2010
- 证据：`this.pretendradio = "";`，v-model 却是 `Placeradio`
- 场景：再打开仍勾着上次那行，直接点「应用」把已选装货地清成空串。
- 修法：改 `this.Placeradio`。把握：high。

### [🟢] 磅单净重无数字校验、半分边界
- 位置：AddPoundrder.vue:475-481；EditPoundrder.vue:498-504
- 证据：`this.AddData.netWeight = (this.AddData.roughWeight*1 - this.AddData.tare*1).toFixed(2)`
- 场景：输入字母得字符串 "NaN"；5.35 − 2.675 得 "2.67"（十进制应 2.68）。
- 修法：`Math.round((Math.round(r*1000) - Math.round(t*1000)) / 10) / 100` 并加数字校验。
- 把握：medium（边界值为推断，未运行）。

### [🟢] 其他
- `ids: [{ required: true }]` 对空数组不触发（缺 `type:'array'`），带星号的合同编号可不关联就保存（AddWaybill.vue:1468-1470、EditWaybill.vue:1219-1221）。medium。
- 扩展名校验区分大小写，`.JPG/.PNG` 被拒（AddWaybill.vue:1997-2006 等四处）。
- doc/docx 预览走 pdf.js，失败被 `.catch(() => {})` 吞掉 → 无限「文件加载中」（Waybill.vue:600-621、Poundrder.vue:574-595）。
- 死代码：AddWaybill.vue:1835-1891（含 `this.splice`、未声明字段）、WaybillDetails.vue:114-288、WaybillManag.vue:736-738/866-871。

## 重复模式
loading 不复位、先请求后置页码、上传回调、写死签名 URL 在范围内各重复 4~16 处，均已合并；全仓其他列表页大概率同样存在。

## 已核对无问题
- 无 v-html/innerHTML/eval，动态内容均经 `{{ }}` 转义；`window.open` 前缀固定，无 `javascript:` 注入面。
- 无明细行增删、序号当 id 的逻辑；除导入弹窗外无 setInterval/echarts/地图/addEventListener。
- 删除、取消均判断 `code == 200` 后才提示成功。