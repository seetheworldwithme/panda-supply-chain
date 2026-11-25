徐先生，评审完成。先说明三点：11 个文件均完整读完 template 与 script（style 跳过）；Thaw、TradeContract、LogStat 三个页面的路由已被注释，用户到不了；没有发现 🔴 级问题，有 11 条 🟡。

路径基准：页面在 `/Volumes/code/GoProject/panda-supply-chain/frontend/src/views/Ptmain/`，路由是 `/Volumes/code/GoProject/panda-supply-chain/frontend/src/router/index.js`，后端在 `/Volumes/code/GoProject/panda-supply-chain/backend/internal/`。node_modules 未安装，Element 分页行为凭 2.15 源码记忆，相关处标 medium。

## 范围
- 读完全部 11 个页面，并核对 http.js、router:980-1140、20 个 api 模块定义及后端对应处理器（commonrest/template.go、warehouse/content_create.go、leasing_write.go、authadmin/common.go）。

## 发现

### [🟡] 意见审核回调引用未定义的 `that`，失败也显示"已通过"
- 位置：advocate.vue:98-130
- 证据：`data.status = 1; this.api.busopin.edit(data).then(res=>{ that.$message({`
- 场景：点"通过/拒绝"后回调抛 ReferenceError，无提示不刷新；后端失败时行已改成 `status=1`，显示"已通过"而红包未发。
- 修法：`then` 内判 `code==200` 再提示并刷新。
- 把握：high。

### [🟡] 云仓分类换上级不清下级 id，必填校验形同虚设
- 位置：Cloudcloud.vue:263-296
- 证据：`this.editData.fstProId = arr[0]*1;` … `this.editData.secName = null;`
- 场景：改一级大类后二级显示为空，但 `bigProId/threeProId/threeName` 仍是旧值，校验通过，存成"新一级 + 旧二三级 + secName=null"。新增同理；两弹窗共用 handler 互写，`addPurs()` 不重置 `AddData`。
- 修法：换上级时清空下级 id/名称，handler 只写当前弹窗对象。
- 把握：high；后端为通用整行更新。

### [🟡] 新建仓号无防重复提交
- 位置：Cloudcloud.vue:350-374
- 证据：`if (valid) { this.api.ware.add(this.AddData)`
- 场景：双击"确定"发两次 `busWarehouse/add`，同仓号两条。
- 修法：`submitting` 标志 + 按钮 `:loading`。
- 把握：high；后端 tplAdd 不对该表查重（template.go:713-723）。

### [🟡] 租赁售卖：分页恒回第 1 页；删除回调字段拼错
- 位置：rentalSales.vue:726-728；996
- 证据：`this.current = 1;`（`getpurcon` 开头）；`if(res.data.cdoe == 200){`
- 场景：点第 2 页触发的 `current-change` 又把页码置 1，超 100 条的记录无法查看审核；删除成功后条件恒假，行留在表里且无提示。
- 修法：`current = 1` 移到 `search/prologistics/changeType`；`cdoe` 改 `code`。
- 把握：high。

### [🟡] 编辑弹窗关闭不清附件列表，残留上一条图片
- 位置：rentalSales.vue:844-846、799-838、901-908
- 证据：`editShowClose() { this.editShow = false; }`
- 场景：看 A（有证书）后点 X 再开 B，上传框显示"A+B"而 `cerList` 只有 B 的，删 A 的残留图会按下标删掉 B 的附件 id；`checkList` 同样残留。
- 修法：`editInfo` 开头重置 `certificate/imgList/checkList`。
- 把握：残留显示 high；下标错位 medium。

### [🟡] `rules` 重复键被覆盖，导入未选文件也能点"导入"
- 位置：BusinessAdjust.vue:430-434、446-469
- 证据：`rules: { importnav: [` … `rules: { month: [`
- 场景：后者覆盖前者，`fileForm` 无必填规则，进度到 99% 后 `upload.submit()` 无文件，弹窗卡在"正在导入..."。
- 修法：导入规则改名 `fileRules`。
- 把握：high。

### [🟡] 业务核算行内编辑无数字校验，非法输入也提示成功
- 位置：BusinessAdjust.vue:598-698
- 证据：`this.editAdjustForm.busCount = this.newbusCount *1;`
- 场景：输入字母得 `NaN`，序列化成 `null`，后端 `validateBody` 放行，页面却提示"编辑成功"；空格存成 0，负数照存。
- 修法：提交前校验 `Number.isFinite` 且非负。
- 把握：NaN→null 为 high；落库结果 medium。

### [🟡] 移除 banner 即删服务器文件，点"取消"后配置仍引用它
- 位置：vipment.vue:680-693 等 4 处
- 证据：`yeRemove(file, fileList) { this.api.addfile.del({id:this.editData.banner1})`
- 场景：点图片删除图标后再点"取消"，文件已删而 `banner1` 未保存，轮播图失效。
- 修法：`on-remove` 只清本地，删除延后到保存成功后。
- 把握：medium-high。

### [🟡] 先发请求再重置页码（6 个页面）
- 位置：BusinessAdjust.vue:706-720、Cloudcloud.vue:256-260、advocate.vue:107-109、chdetails.vue:157-161、visit.vue:349-362、redenvecord.vue:189-242 等
- 证据：`this.getpurcon(); this.current = 1;`
- 场景：停在第 3 页时搜索，请求仍带 `current=3` 得空结果；页码随后被改成 1 但分页组件不再触发请求。
- 修法：先置 1 再请求。
- 把握：medium-high。

### [🟡] 加载失败后 loading 不复位，整页被遮罩锁死
- 位置：BusinessAdjust.vue:534-548（复位在 `if` 内，业务码非 200 也卡）；Cloudcloud.vue:246、rentalSales.vue:741、advocate.vue:143、chdetails.vue:148、redenvecord.vue:227、visit.vue:399
- 证据：`if(res.data.code == 200){ this.loading = false;`
- 场景：任一 502/超时被拦截器 reject 后无 `.catch`，`loading` 永远 true。
- 修法：补 `.catch(() => {}).then(() => this.loading = false)`。
- 把握：high。

### [🟢] 其余
- `validate` 回调里 `}else if(res.data.code == 602){` 引用未定义的 `res`：Cloudcloud.vue:425、vipment.vue:955、chdetails.vue:197。
- BusinessAdjust.vue:339-349 类型校验挂在 `on-change`，非 xlsx 照样可导入；`:801` 局部 `setInterval` 离开路由不清理。
- redenvecord.vue:37 累计金额只在改日期时请求，进页面、搜索后显示 0 元。
- 写死的签名/AuthorizationId 拼进 URL：BusinessAdjust.vue:354、rentalSales.vue:811/831、vipment.vue:846/896/907/919，同 http.js:13-14，归 lead 统一处理。
- http.js:38-43 把业务失败的 msg 以 success 样式弹出（lead 范围）。
- 死页面 LogStat、TradeContract、Thaw（router/index.js:993-999、1009-1031）：Thaw.vue:56 只传一个参数而 `purSee(index,row)` 取 `row.id`；TradeContract.vue:671-682 在 `function(){}` 内 `clearInterval(this.calearfin)` 清不掉倒计时；解冻/退款提交无防重点。建议删除。

## 重复模式
整行对象回传 `updateById`（BusinessAdjust、rentalSales、vipment、chdetails、advocate）：会用旧快照覆盖并发修改，改动面大，仅记录。

## 已核对无问题
- api 方法名与参数形态全部对得上；无 `v-html/eval`，无 query 驱动的 `href/location`，无 `sessionStorage.queryAll` 用法。
- BusinessAdjust 新增双击由后端 `(enterprise_name, month)` 查重兜住；advocate 重复点击由后端 `FOR UPDATE` + `before==0` 守卫兜住，不会重复发红包。
- `upload/batch` 的 id 是十进制字符串，`Number()` 对自增 id 安全；金额只做 `/100` 展示，无浮点累加。

🎉 vipment.vue:852/925、chdetails.vue:165 深拷贝后再编辑，不直接改表格行。