徐先生，评审完成。文件均在 `/Volumes/code/GoProject/panda-supply-chain/frontend/src/views/Mindex/`；简写 APB=AddPurchaseBill、EPB=EditPurchaseBill、PB=PurchaseBill、ASB=AddSalesBill、ESB=EditSalesBill、SB=SalesBill、MP=MainPage（均 .vue）。

## 范围
7 个文件的 template+script 全部读完（MP 至 2590 行，style 跳过）；另读 http.js、filters.js、main.js 与两个结算 api 模块，并核对后端 capital/settlement.go:224-332、bill/stats.go:148-394。仓内无 node_modules，el-upload、el-pagination 源码未读。

## 发现

### [🔴] 编辑页离开不清 queryAll，残留被「新建合同」页当编辑保存，覆盖无关合同
- 位置：EPB:792、923、934；ESB:714、834、845
- 证据：EPB 只在取消或保存成功时 `removeItem('queryAll')`，无 beforeDestroy。读方 AddPurchaseContract.vue:1040-1042、1547-1550：`if(sessionStorage.getItem('queryAll')){ this.AddData = JSON.parse(…)`，保存走 `this.api.cgContract.edit(this.AddData)`；入口 PurchaseContract.vue:1409 不清键。
- 场景：编辑采购结算后经侧栏或后退离开，再到采购合同列表点「新建」：表单被结算单记录替换（`singers.length` 抛错中断 mounted），补全后保存即按结算单 id 更新同 id 的合同。AddSalesContract.vue:987、1495 同。
- 修法：两个 Edit 页加 `beforeDestroy(){ sessionStorage.removeItem('queryAll') }`；根因需全局处理。
- 把握：high。

### [🔴] 结算单「保存」无防重复，双击建两单
- 位置：APB:8、801-831；ASB:8、753-783；弹窗新建单位 APB:445、1218
- 证据：`<el-button size="small" @click="AddDatasubmitForm('AddData')">保存</el-button>` 无 loading/disabled；后端 stAdd 直接 InsertNoFill，无 systemCode 唯一校验。
- 场景：慢网双击，两次 add，生成两条同编码结算单，合计金额翻倍。
- 修法：提交置 saving 标志并 `:loading`，finally 复位。
- 把握：high；库上若有唯一索引则降级。

### [🟡] MP 的 echarts 实例从不 dispose
- 位置：MP:801、1149、1576、1772、1968、2159、2336
- 证据：`var myChartsettlement = echarts.init(settlement);` 局部变量，全文件无 beforeDestroy。
- 场景：进出仪表盘 N 次泄漏 7N 个实例（含地图飞线动画帧）；echarts 4 下每次筛选再 init 一个。
- 修法：实例存非响应式字段，init 前判空，beforeDestroy 逐个 dispose。
- 把握：high。

### [🟡] 查询先发请求再重置页码；失败后 loading 不复位
- 位置：PB:1044-1048、1050-1054、1164-1166、910-934；SB:867-868、1008-1018、1128-1130、874-898；Add/Edit 四页的 contractsearch、searchtwo
- 证据：先 `this.getpurcon()` 后 `this.current = 1;`；`this.loading = true;` 只在 then 内复位，无 catch。
- 场景：停在第 3 页搜索或删末页，请求仍带 current:3，分页器显示第 1 页而表里是第 3 页或空表；接口 500 被 reject 后「加载中」遮罩永久盖住整页。
- 修法：先置 current=1 再查询；补 catch 复位 loading。
- 把握：high。

### [🟡] 导入假进度 setInterval 未清理
- 位置：MP:2491-2498；PB:858-865；SB:828-835
- 证据：`let f = setInterval(()=>{ this.perstep += 3; if(this.perstep == 99){ this.$refs.upload.submit(); clearInterval(f); } },150)`
- 场景：点「导入」后 5 秒内关弹窗，定时器照常 submit，一键导入照跑；此时切路由，`$refs.upload` 为 undefined 抛错，clearInterval 执行不到，定时器永续。
- 修法：句柄存实例，在关闭函数与 beforeDestroy 里 clear。
- 把握：high。

### [🟡] 仪表盘统计口径与文案不一致
- 位置：MP:1512-1514；MP:1279-1290、1246-1251
- 证据：`salechange(name){ this.saletime = name; },`；`this.screeninvoice = this.topYear; this.TimeFilter = this.yearValues + '年 ~ ' + this.topYearSize;`
- 场景：「销售数据」选「本月」，数字仍是旧区间；筛选器点「上半年」后改年份再确定，请求旧年区间而按钮显示新年，重开弹窗区间被重置但高亮仍是「下半年」。
- 修法：salechange 同 settlemenchange 调两个统计接口；screenchange 按当前 yearValues 重算 topYear*。
- 把握：high。

### [🟡] 地图丢最后一个发货地分组
- 位置：MP:814-826
- 证据：`if(d.length>0 && d[d.length-1][0].name != res.data.data[i].loadingArea){ c.push([d[0][0].name,d])`，循环后的 d 未入 c。
- 场景：只有一个发货地的企业 series 为空，地图无线路；多地时末组缺失。
- 修法：循环后补 push。
- 把握：high。

### [🟡] 结算时间 form-item 的 prop 误写
- 位置：ASB:171；ESB:171；EPB:214（APB:214 正确）
- 证据：`<el-form-item prop="goodsName"> <el-date-picker v-model="AddData.Time"`
- 场景：Time 必填失效，新建销售结算可空结算时间提交；编辑页清空后保存，`this.AddData.Time[0]` 抛错，按钮无反应。
- 修法：改 `prop="Time"`。
- 把握：high。

### [🟡] 附件列表与提交内容脱节
- 位置：APB:895-923；ASB:846-874；EPB:960-990；ESB:871-901
- 证据：`if(response.code == 200){ … this.enclosurefileList = fileList; }`；fileRemove 仅在 del 返回 200 时同步。
- 场景：上传业务失败，界面仍显示成功，保存时不含该附件；删除接口失败，界面已移除而 id 仍被提交。
- 修法：失败分支剔除该文件并提示；fileRemove 先同步 fileList。
- 把握：medium-high，取决于 el-upload on-remove 是否等待 Promise。

### [🟢] 小项
- MP:811、1419、1464、1718、1912、2099、2272：`res.data.data != []` 判空恒真后取 `.length`；后端 selZeroTax 无告警时返回 nil 切片即 `data:null`，有票据无税负告警的企业每次进仪表盘 tskbAll 抛错。改 `(res.data.data || [])`。
- MP:2541 导入完成再调 tskbAll 未清空，提示重复追加；MP:1558 `this.credit = '良好'` 触发后不复位；MP:1393 二月固定 28 天。
- 关联弹窗未选行点「应用」：APB:1164 `this.placecontent.id` 抛错，EPB:1238 清空合同与供应商输入；EPB:1237、ESB:1136 换绑合同后 contractCode、supplierId 不随动，PB:72-74 仍显示旧合同编码。
- 下载 URL 的 fileName 未 encodeURIComponent；导入 on-change 返回值被忽略，非法类型仍可上传。

## 重复模式
- `queryAll` 为全站单键，约 15 个 Edit 页同样只在取消或成功时清除，建议统一 beforeDestroy 清理。
- 写死的 Signature、AuthorizationId 在这 7 个文件重复并拼进 downById、previewStream 的 URL；按名取文件是否校验归属、MP:1714 传入的 companyId 是否被后端采信，均未读后端，请核实。

## 已核对无问题
- 无 v-html、eval，服务端字段均走文本插值。
- 数量、单价、总价以字符串原样提交，前端不做乘除，toFixed 仅用于显示，过滤器已对 null 设防。
- 编辑页 supplierId 为 undefined 被丢弃：后端 stUpdateById 只更新请求体出现的列，不会清空。
- 图表容器无条件渲染，init 均在 mounted 或回调内，无 window 监听与轮询；提示语下标切分与后端格式一致。
- 付款单、发票关联弹窗与详情抽屉入口已注释，属死代码。