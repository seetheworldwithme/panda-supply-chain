徐先生，以下是应用外壳层（API / 路由 / store / 全局辅助 / 顶层页面）的审查报告。

## 范围
- 通读 104 个文件：
  - `frontend/src/api/` 全部 71 个；`path.js.orig-backup` 只做了 diff，与 `path.js` 仅差 env 和 baseUrl/coludUrl 三行。
  - `router/index.js`、`store/index.js`、`main.js`、`App.vue`、`filters.js`、`upload.js`、`customTranslate.js` 及 `customTranslate/*` 2 个、`components/*` 4 个。
  - `views/` 顶层 18 个页面的 template 和 script。
- 跳过：所有 `<style>` 块。
- 为核实契约只读查阅了后端 `auth/middleware.go`、`auth/compat.go`、`auth/handler.go`、`platform/resp.go`、`menu/handler.go`、`authadmin/*`、`commonrest/enterprise.go`、`main.go`。
- 限制：仓库没有 node_modules，element-ui 2.15.6 与 vue-router 3.5.1 的运行时行为按已知源码推断，没有实跑。
- 过程说明：交叉核对时在 /tmp 建过两个临时清单并已删除，仓库内零写入。

## 发现

### [🔴] ELogin「手机号+验证码」登录，后端不校验验证码（跨范围，修在后端）
- 位置：`frontend/src/views/ELogin.vue:304-307`、`api/mobilelogin.js:8` → `backend/internal/auth/compat.go:79-103`、`handler.go:73-76`
- 证据：
```go
phone := strings.TrimSpace(r.FormValue("phone"))   // 全程未读 code
u, err := s.selectByUserName2(phone)
// WHERE account_type != 4 AND handle_result = 1 AND (phone = ? OR account = ?) ORDER BY level DESC LIMIT 1
platform.Success(w, tokens)
```
- 场景：`/ELogin` 免登录。填任意已启用账号的手机号或账号，加任意 6 位数字，就拿到该用户（含企业管理员）的 JWT。
- 修法：发 token 前调用 `VerifyPhoneCode(phone, code)`；前端已在提交 code，无需改动。
- 把握：high；若另有路径校验则否定，但 `main.go:78` 只注册了这一个 handler。

### [🔴] 硬编码的 Signature/AuthorizationId 是后端认可的通用匿名凭据（跨范围）
- 位置：`frontend/src/api/http.js:13-14`、`views/Register.vue:359-360,562` → `backend/internal/auth/middleware.go:68-83`、`authadmin/authadmin.go:30-33`
- 证据：
```js
config.headers.Signature = '328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f';
config.headers.AuthorizationId = 'wxda0b1c468b83bc99';
```
```go
user.SignedAnon = checkSignature(r, secretDB)
if user.UserID == 0 && !user.SignedAnon && !publicRequest(r.URL.Path) { platform.NotAuth(w) }
```
- 场景：签名随 JS 包公开。带上这两个头就能匿名调用 `/saas-auth/chainAuthUser/ptlist`（返回平台账号列表，`user_compat.go:18-25`），以及角色、菜单、e签宝配置等管理接口，这些 handler 都不检查 UserID。
- 修法：需要登录的路由要求 `UserID != 0`，SignedAnon 只放行登录前的白名单接口。
- 把握：high；前提是 `chain_auth_secret` 里有该 authorization_id，后端自己的 imageShow 也依赖它。

### [🟡] body code 403（HTTP 200）无人处理，token 失效后用户被搁浅
- 位置：`http.js:29-36`、`router/index.js:1241-1251`、`store/index.js:314`；`backend/internal/platform/resp.go:38-44`（NotAuth 返回 HTTP 200）
- 证据：
```js
if (response.status === 200) { …; return Promise.resolve(response); }
if(res.data.code == 200){ … }     // asynSetRoutes：非 200 什么都不做
```
- 场景：账号被禁用、删除或 token 过期后，所有接口返回 `code:403`。页面只弹绿色「请求未授权，访问失败」，不跳登录。此时按 F5，`selUserMenus` 返回 403，不注册路由，页面白屏。
- 修法：拦截器识别 `response.data.code === 403`，清除 Authorization 并跳 Login。注意 `authadmin/user.go:562` 也用 403 表示无权限，需后端换码，或只对引导接口处理。
- 把握：high。

### [🟡] 登录成功后的跳转不等动态路由注册；无菜单用户白屏
- 位置：`Login.vue:248,259-266`、`store/index.js:311,389-391`
- 证据：
```js
this.$store.dispatch('asynSetRoutes');          // action 不 return promise
duration:1000, onClose:()=>{ this.$router.push({name:"Main"}); }
```
- 场景：
  - `selUserMenus` 比 toast 慢（约 1.3 秒）时，push 到尚未注册的命名路由，落到 `/`，随后 catch-all 生效显示 404。
  - `backMenus` 为空时根本不注册路由，同样白屏；原先的「未分配权限」提示已被注释。
- 修法：action 返回 promise，登录在 `.then` 里跳转；无菜单时给提示。
- 把握：medium-high，需要慢网络复现。

### [🟡] 路由表被原地裁剪、meta 权限标记不清除，同标签页换账号后继承上一账号的路由和按钮权限
- 位置：`store/index.js:321,329-352,361,369`；范围内的非整页登出点：`App.vue:60-61`、`Main.vue:490-491`、`NotFont.vue:16-17`、`Xzmain.vue:460-461`
- 证据：
```js
addRoutes = activeRouterMap;
routes[i].meta.add = 1;
routes[i].children = getRouteData(routes[i].children)
addRoutes[0].children = getRouteData(activeRouterMap[0].children);
```
- 场景：A（权限少）闲置 6 小时被登出（不 reload），B（权限多）在同一标签页登录。过滤在 A 裁剪过的树上再做一次，B 的侧栏和路由只剩 A 的子集。页面里 `$route.meta.add/edit`（55 和 48 处，用于显隐按钮）保留 A 的标记。
- 修法：过滤前深拷贝路由记录；或让这几处登出与 `Main.vue:568` 一致，执行 `this.$router.go(0)`。
- 把握：high。

### [🟡] 修改密码弹窗关闭时把表单模型换成空对象，字段失去响应性
- 位置：`Main.vue:505-508`
- 证据：`BuyerhandleClose(){ this.dialogFormVisible = false; this.ruleForm = {} }`
- 场景：点 X 关闭后再打开，输入后点「眼睛」图标，输入框被清空。再点过一次「取消」后，输入任何内容校验都报「请输入原密码」，只能 F5。
- 修法：改为 `this.$refs.ruleForm.resetFields()`。
- 把握：失去响应性是 high；具体表象 medium，没有实跑。

### [🟡] Forget.vue 使用未导入的 `Message`，验证码错误时抛 ReferenceError
- 位置：`Forget.vue:142-147`（文件没有任何 import）；`backend/internal/authadmin/user_compat.go:87-89` 返回 `FailMsg`（code 999）
- 证据：`if(res.data.code == 999){ Message({ message:'身份验证失败!', type:'error', }) }`
- 场景：忘记密码第一步输错短信码，页面层的红色提示不出现，只剩拦截器的绿色「失败」。
- 修法：改成 `this.$message({...})`。
- 把握：high。

### [🟡] 登录页「获取验证码」在账号不存在（2002）或请求失败后永久禁用
- 位置：`Login.vue:381-408`；`backend/internal/auth/handler.go:236-243`
- 证据：
```js
this.isDisabled = true;
if(res.data.code == 200){ … }
if(res.data.code == 999){ that.isDisabled = false; }
if(res.data.code == 2006){ that.isDisabled = false; }
```
- 场景：账号打错后点「获取验证码」，后端返回 2002「用户名不存在」，按钮一直灰。改对账号也无法再获取，只能 F5。
- 修法：`if (code == 200) {…} else { that.isDisabled = false; }`。
- 把握：high。

### [🟡] 营业执照上传的大小校验按 GB 算，上限实际是 3GB
- 位置：`Register.vue:578`
- 证据：`const isLt3M = file.size / 1024 / 1024 /1024 < 3;`
- 场景：提示写「不能超过 3MB」，实际 3GB 以内都放行，且这个上传接口免登录。
- 修法：去掉一个 `/ 1024`。
- 把握：high。

### [🟡] 注册三步的「下一步」没有防重复提交，后端对应接口也非原子
- 位置：`Register.vue:46,90,191`；`backend/internal/authadmin/user_compat.go:100-119`、`commonrest/enterprise.go:37-103`
- 证据：`<el-button type="primary" @click="submitFormLogin('loginForm')" style="width:100%;">`（没有 `:loading`）
- 场景：
  - 双击第一步会发两次 `setAccount`。验证码校验后不消费，插入也不查重，结果是同一手机号两条账号。
  - 双击第三步会发两次 `enterpriseAdd`（先查后插），结果是重复企业和重复的 enterprise_code。
- 修法：三个按钮加 `:loading` 或提交标志；后端补唯一约束。
- 把握：medium；库上若有唯一索引，第二步会被拦，本地没有库无法核实。

### [🟡] path.js 三处路径多一个斜杠，直连 Go 时 301 → GET → 405
- 位置：`path.js:91-93`（busEnterpriseLeasing、busconfig、busbenfig）
- 证据：`var busEnterpriseLeasing = '/saas-warehouse//busEnterpriseLeasing/'`
- 场景：Go 的 `http.ServeMux` 会把非规范路径 301 到清洗后的路径，浏览器跟随时把 POST 改成 GET，于是命中不了 `POST /…` 模式，返回 405。`Ptmain/rentalSales.vue` 和 `vipment.vue` 的列表都会显示「连接出错(405)」。当前 `baseUrl` 正是直连 127.0.0.1:8765；经 nginx（merge_slashes）可能侥幸通过。
- 修法：删掉多余的斜杠。
- 把握：机制是 high；部署相关部分 medium。

### [🟡] 政府看板筛选没有请求序号，快切地区时旧响应覆盖新数据
- 位置：`Xzmain.vue:641-675`（市、县、乡、企业四处同构）、`970-1011`、`1099-1130`
- 证据：`this.countCapital({xzqhCode:value}); this.countBillss({xzqhCode:value}); … this.paymentCapital(invoiceDates);`
- 场景：连点「市 → 县」会各发 7 个请求，先发后到时，卡片和图表混入上一个地区的数据。`paymentCapital` 内嵌的 `.then` 与外层的数组重置交错，会让两组序列长度错位。
- 同文件另有两处：
  - `echarts.init` 的实例没有在 `beforeDestroy` 里销毁。
  - `selectItemtwo` 把二月末日写死成 `-02-28`（`:622`），闰年漏掉 29 日。
- 修法：每个方法带自增序号，回调里序号不等就丢弃；`beforeDestroy` 调 `dispose()`；二月末日用 `new Date(y,2,0).getDate()`。
- 把握：medium-high。

### [🟡] 已提交的 path.js 指向 127.0.0.1，备份文件一并入库（发版检查）
- 位置：`path.js:4-9`、`path.js.orig-backup`（自 fc11203 起如此）
- 证据：`var baseUrl = "http://127.0.0.1:8765"`
- 场景：按 HEAD 构建上线，所有接口都打到访问者本机。
- 修法：发版前切回域名，或改由 `.env` 注入；删除 `.orig-backup`。
- 把握：high；可能是有意为之，仅作提醒。

### [🟢] http.js 响应拦截器的三个细节
- 位置：`http.js:30-34,44`
- 证据：`if(response.data.msg != '成功' && response.data.code !=701){ Message({ message:response.data.msg, type:'success', }) }`
- 场景：
  - 业务错误用绿色 success 样式弹出，并和页面自己的红色提示重复。
  - 大屏图形码是 Blob（`api/imglogin.js:7`），`blob.msg` 是 undefined，仍会触发 `Message({message: undefined})`。按 element 2.15.6 的实现，要么弹空 toast，要么抛错使请求 reject、图形码不显示；本机无法核实是哪一种。
  - `case 602: err.message = err.response.status` 赋的是数字。
- 修法：加 `'code' in response.data` 守卫；错误改用 `warning`。
- 把握：第一、三点 high，第二点 medium。

### [🟢] 全局过滤器对空值和非数值不设防
- 位置：`filters.js:11-14,17-19,46-61,78-80`
- 证据：`time = value.substring(0,10);`；`value = value.toFixed(limit);`（默认值 `'0'` 是字符串）
- 场景：
  - 字段为 null 或 undefined 时，渲染期 TypeError，整表不刷新。
  - `ToFiexd` 对 undefined 显示「NaN」。
  - `newcurrency` 遇到字符串金额会抛错。
  - 现有 44 处调用里，只有 `DPmain/warehouse.vue:154` 没有 v-if 守卫。
- 修法：入口统一加 `if (value == null || value === '') return ''`；`newcurrency` 改成 `(value*1).toFixed`。
- 把握：high。

### [🟢] 死代码、错挂与清理类（合并）
- **错挂与无引用文件**
  - `views/Service.vue`：公开路由 `/Service` 挂的是 bpmn-js 建模器页，不是协议文本；`save()` 调用不存在的 `this.$apis`（`:61,78`）。
  - `Demo.vue`（script 块写成了 `type="text/wxtag-template"`）。
  - `Nodata.vue`。
  - `components/link.js`：是压缩版微信 JS-SDK，没有引用。
  - `src/customTranslate.js`：与 `customTranslate/customTranslate.js` 重复。
  - `components/NumberCounter.vue`：没有引用；value 为 NaN 时 `setInterval` 永不清除。
- **API 层死代码**
  - `api/file.js` 的 `pree`：`http.get(url, data)` 的第二参应是 config，没有调用；同文件的 `Header` 导入也没用。
  - `basics` 和 `basicslist` 的 `one/edit/del`、`getSignUrl`。
- **路由**：`router/index.js:76` 的路径 `'/InitiateStart '` 带尾空格。
- **Login.vue**
  - `:182` 的全局 `document.onkeydown` 只在登录成功时清除。经 router-link 去 Register 或 Forget 后按 Enter，会调用已销毁的实例，控制台报 TypeError。
  - 按钮聚焦时按 Enter 会触发 3 次登录请求（keydown、click、keyup 各一次）。
- **Forget.vue:187-190**：点「立即登录」把 `Jump` 置 0，倒计时 interval 减到 -1 后永不清除。
- **App.vue**
  - `:15-54` created 与 mounted 各跑同一段逻辑，刷新时 `selUserMenus` 和 `selMenus` 各请求两次。
  - `timeOut` 注释写 30 分钟，实际是 360 分钟。
- **CommonView.vue:79,85**：写的是未定义的 `this.MainPage`，`routeName` 被存成字符串 `"undefined"`。后端目前不读这个头，暂无影响。
- **upload.js:2**：Authorization 在模块加载时被快照，同标签页换账号后分片上传带的是旧 token。
- **Register.vue**
  - `:442` 的 `console.log(this.manageForm)` 输出明文密码。
  - `rule.username` 与 `prop="phone"` 对不上，手机号格式规则没有生效。

## 重复模式
- 「`code == 200` 才处理，其余静默」遍布全仓，是 403 搁浅、验证码按钮卡死、无菜单白屏的共同根因，适合在 http.js 集中兜底。
- 定时器和全局监听不随组件销毁：`Login.vue`、`Forget.vue`、`Register.vue`（两个 5 秒倒计时）、`Xzmain.vue`（echarts）。
- 范围内有 4 处非整页登出（App、改密、404 页、Xzmain），都会和「路由表原地裁剪」叠加；只有 `Main.vue` 的手动退出会 reload。

## 已核对无问题
- **路由守卫**：未登录访问动态路由一律回 Login。企业路由表和平台路由表按 level 互斥注册，每次都重置 `router.matcher`，没有重复注册。普通用户用 URL 直达对方路由得到 404。
- **API 层**：`path.js` 的键齐全，模块导出与 `index.js` 一致，313 处 `api.X.Y` 调用全部有定义。端点字符串除上面的 `//` 和死代码外，都在后端找到了注册。
- **XSS 与泄露**：范围内 `v-html` 只绑定内部常量字符串，没有 eval 或 innerHTML 接外部数据。`wxindex.vue` 和 `sergrent.vue` 没有 redirect_uri 或 redirect 跳转（后者是纯文本协议页）。token 只存 sessionStorage，密码走 form POST，不进 URL。
- **登录请求**：用 URLSearchParams 加 form 编码，与后端 `ParseForm` 一致；图形码的 Blob 处理正确。
- **过滤器调用点**：`orderStrTime`、`cardCode`、`Time` 系列都有 v-if 守卫（只有上面 1 处例外）；`currency` 没有 limit=0 的用法。
- **未能核实**：库内 id 是否超过 2^53。本地没有库数据；范围内也没看到对 id 做数值运算的写法。
- 🎉 store 按角色只注册对应路由并重置 matcher，路由级隔离干净。
- 🎉 登录和改密的凭据处理克制：form POST、sessionStorage、登出清理。