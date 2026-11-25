徐先生，以下是 setup/ 与 Emain/ 的前端评审报告。路径均相对仓库根 `/Volumes/code/GoProject/panda-supply-chain/`。

## 范围
18 个文件全部读完（setup 14 + Emain 4，读到末尾 `<style>` 前），并对照 api 模块、store 菜单结构和后端对应 handler。项目无 keep-alive、无 node_modules（element-ui 源码未能核对）。

## 发现

### [🔴] 仅凭 URL 里的 account 就能换到登录 token
- 位置：frontend/src/views/Emain/Sponsor.vue:155-180；backend/internal/commonrest/chat_local.go:52-87
- 证据：`this.api.wxChat.wxChat({account:theRequest.account})` → `sessionStorage.setItem('Authorization','Bearer '+res.data.data.token)`；后端按 `phone=? OR account=?` 查到用户即签发 JWT，无密码/验证码
- 场景：打开 `…/#/Sponsor?account=<他人手机号或账号>&level=1` 即以该用户身份登录
- 修法：后端改为校验小程序侧签名或一次性 ticket；前端不再从 URL 取账号
- 把握：high

### [🔴] “查询签署”用 HTML 字符串渲染签署人姓名（存储型 XSS）
- 位置：Emain/EsignIndex.vue:253-261、Emain/Sponsor.vue:320-328
- 证据：`arr.push('<div>' + res.data.data[key].signerName + ' : ' + day + '</div>')` → `dangerouslyUseHTMLString: true`
- 场景：发起签署时“签署方姓名”填 `<img src=x onerror=…>`，对方点“查询签署”即执行脚本，可读走 token
- 修法：去掉 `dangerouslyUseHTMLString`，改纯文本或 VNode
- 把握：high（若 e签宝会清洗姓名则不可利用，代码仍不安全）

### [🔴] “发 起 签 署”无防重，双击发起两条签署流程
- 位置：Emain/InitiateStart.vue:164、459-544
- 证据：`this.api.eqbflows.launchFlow(this.AddData)`，按钮无 loading；后端 contract/esign.go:147-165、221-225 每次都插 flow+signers 并调 `createFlowOneStep`，无幂等
- 修法：提交标志 + `:loading`，finally 复位
- 把握：high

### [🟡] 角色权限页：回显被覆盖、切换会串写
- 位置：setup/Roleperm.vue:171-199、306-311、372-373；PtRoleperm.vue:160-189、292-297、360-361；QyRoleperm.vue:160-188、286-291、351-352
- 证据：`this.setmenudata(this.treedata); this.$refs.addtree.setCheckedNodes(newArr);`，而 `handleNodeClick` 对刚勾选的叶子 `data.button.push('新增')…`
- 场景：①回显勾选触发 check-change，新勾上的叶子被追加全部可用操作，盖掉已保存的选择；管理员勾一个新菜单就保存，该角色所有菜单操作权限被写大（后端 saveEtcAuths 先删后插）。②`this.roleID = id` 先改再请求回显，请求失败或快速切换乱序返回时，树上是 A 的权限而 roleID 已是 B，保存即整体串写
- 修法：①回显前置标志位、`$nextTick` 后清除，`handleNodeClick` 遇标志 return（或改用 `@check`）；②丢弃过期响应，失败时 `roleID = null`
- 把握：① Roleperm high，另两页 medium（依赖 el-tree 异步派发）；② medium

### [🟡] 签署方按 id 分组，删除中间签署方后他人被合并或丢失
- 位置：Emain/InitiateStart.vue:375-411、480-518
- 证据：`for(let i=1;i<=this.numbers;i++){ … if(i === this.AddData.domains[key].id){ Object.assign(a,…) } arrays.push(a) }`，而 `removeDomain` 只做 `this.numbers--`
- 场景：加 A(1)、B(2) → 删 A → 再加 C（id 仍是 2）→ `singers=[{}, B与C合并]`；3 个删第 2 个 → 第 3 个整组丢弃；企业/个人切换后隐藏的另一半也照常提交
- 修法：id 用只增计数器并按实际 id 收集；切换类型时清空另一数组
- 把握：high

### [🟡] 删除已传文件后点“发起签署”静默无反应
- 位置：Emain/InitiateStart.vue:462、520、631-645
- 证据：`fileRemove` 成功后未清 `this.AddData.field`；提交时 `this.AddData.contractName = this.enclosurefileList[0].name`
- 场景：上传→删除→发起：`field` 仍非空绕过校验，`enclosurefileList[0]` 为 undefined 抛 TypeError，既不发起也不提示
- 修法：删除成功后 `this.AddData.field = ''`
- 把握：high

### [🟡] 企业轮播图删除时误删相邻图片
- 位置：setup/EnterpriseArchives.vue:478-489
- 证据：`this.enfileList.forEach((value,index,array)=>{ if(value.id == file.id){ array.splice(value,1) }else{ arr.push(value) } })`
- 场景：3 张图删第 1 张 → `splice(对象)` 起始位被转成 0，遍历错位，第 2 张也被挤出，点“更新”时 `ids` 不含它
- 修法：`this.enfileList = this.enfileList.filter(v => v.id != file.id)`。同文件 593 行 `/1024/1024/1024 < 3` 实为 3GB（提示写 3MB）；627 行 `areaCode = []` 会把数组提交给字符串字段
- 把握：high（前端必丢；后端是否解绑未核）

### [🟡] 红包“金额最小值”永远改不了
- 位置：setup/globalConfig.vue:159
- 证据：`if(eagstr.test(this.newSjhxje) && this.newmaxe*1 >= 10){`
- 场景：编辑最小值时 `newmaxe` 为 '' → 条件恒 false，一律提示“请输入正确金额或金额太小”并丢弃输入
- 修法：改为 `this.newSjhxje*1 >= 10`
- 把握：high

### [🟡] 弹窗重开后表单状态残留
- 位置：setup/menuPower.vue:371-381、483-488、507、515-518；setup/QyRoleperm.vue:110、411
- 证据：关闭/成功只做 `this.AdjustForm = {}`，随后 `this.AdjustForm.ids = this.checkedCities`；`that.RoleListData = {}` 丢掉 `roleType:3`
- 场景：①新建菜单 A 勾了查看+编辑，再新建 B 时仍勾着并随 `ids` 提交给 B。②QyRoleperm 第二次起新增的角色不带 roleType，后端缺省 `role_type=2`（authadmin/role.go:116-119），不进本页列表。③规则 `trigger:'chang'` 拼写错，“菜单类型”误用 `prop="pid"`，无必填校验
- 修法：关闭时 `checkedCities = []`；重置为 `{name:'', roleType:3}`；`prop="menuType"`
- 把握：① high；② medium（若 bodyMap 不做驼峰转下划线则首条起就异常，可否定）

### [🟡] 搜索/删除后先请求、后重置页码
- 位置：setup/ContractDeploy.vue:589-593、676-678；setup/menuPower.vue:445-449、616-618
- 证据：`this.getpurcon(); this.current = 1;`
- 场景：停在第 3 页搜索 → 请求仍带 page=3，结果在第 1 页时显示“暂无数据”
- 修法：先置 `current = 1` 再调用
- 把握：high

### [🟡] 编辑配置页后退/直接进入时整页渲染报错
- 位置：setup/EditContractDeploy.vue:536
- 证据：`this.AddData = JSON.parse(sessionStorage.getItem('queryAll'));`
- 场景：保存或取消会清掉 `queryAll`，之后浏览器后退到本路由（组件重建）→ `AddData` 为 null，模板取 `AddData.planName` 抛错白屏
- 修法：为 null 时 `this.$router.replace({name:'ContractDeploy'})`
- 把握：high

### [🟢] 其余防重/失败态
- 位置：setup/AddContractDeploy.vue:596-602、EditContractDeploy.vue:546-552；ContractDeploy.vue:565-586、PtUserList.vue:635-704、UserList.vue:324-396；menuPower.vue:504-525
- 场景：①提交前把 checkbox-group 的 v-model 数组原地 `.toString()`，失败留在页面后复选框全空、再点模型变 `true`；②`loading=false` 只写在 code==200 分支，业务失败后表格永久“加载中”；③配置“保存”双击建两条（后端 eqbAdd 无查重），菜单新建同样无防重
- 修法：提交用副本 `{...this.AddData, noticetype: arr.join(',')}`；`loading=false` 提到 then 首行；按钮加 loading
- 把握：high

### [🟢] 展示与杂项
- 位置：setup/ContractDeploy.vue:439、461-465、643-647；Emain/InitiateStart.vue:352、627
- 场景：查看抽屉“签署类型”误读 `signdatebeantype`；`'yyyy-MM-d'` 永不匹配；点“基本信息”页签调用不存在的 `this.allone` 报错；PDF 预览/下载 URL 硬编码 Signature/AuthorizationId 且 `fileName` 未编码。ContractTemplate.vue 是写死数据的占位页
- 修法：改字段名/删死代码；`encodeURIComponent`
- 把握：high

## 重复模式
- 三个角色权限页整份复制，上述两条问题三处同在；`confselSign` 的 XSS 在 EsignIndex 与 Sponsor 各一份
- “先请求后重置页码”在 EsignIndex/Sponsor 写法相同，但两页没有分页控件，暂不受影响

## 已核对无问题
- 18 页调用的 api 方法在对应模块均存在，path、method、参数包装一致
- 保存权限用 `getCheckedNodes(false,true)`：全选/半选父节点只进 `menuRoles`，未勾选节点不会被带上；按钮经 `new Set` 去重，未勾选的不进 `authoritys`
- PtUserList/UserList 的新增账号、分配角色、删除、恢复只传 id/handleResult，没有原样回传行对象；这两页没有编辑/重置密码功能
- 邀请码页只生成不提交；Application/Unknown 为静态页；18 页内无 v-html、无定时器或监听器残留
- 范围外顺带：backend/internal/commonrest/chat_local.go:185 `!(chat && phone == "15012345678")` 让该手机号在 loginChatPhone 下跳过验证码校验，请后端评审确认