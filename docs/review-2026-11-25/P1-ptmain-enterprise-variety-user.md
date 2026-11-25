徐先生，Ptmain 第一组 11 个页面评审完毕（全程只读，未改任何文件）。

## 范围
`frontend/src/views/Ptmain/` 11 个 .vue 的 template+script 全部读完（PtIndex、EditQycensus、varietyList 分段读至 `</script>`），`<style>` 跳过；demos.vue 仅 2 行占位且未注册路由。另核对 `api/http.js`、`main.js`、`src/api` 相关方法，并定点看了后端 `busProduct/updateById|updateBatch` 的落库方式。

## 发现（未标把握者均为 high）

### [🔴] 整行回传：保存即用旧值覆盖并发修改
- 位置：EditQycensus.vue:1917；varietyList.vue:2015、2113-2115；userasp.vue:1074、1154
- 证据：`this.api.busrodect.updateBatch(this.tableDatas).then(res=>{});`
- 场景：打开企业编辑页后企业在小程序改了价格/库存，管理员只改「业务员」点保存 → 挂载时载入的整批煤种行（旧价、旧库存）被写回，响应还被丢弃。煤种「保存/重新审核」同样回写整条记录。后端 `warehouse/template.go:782`、`product_write.go:422` 都是 `pick(m, dbcols)`，带哪些列写哪些列（用户表走同族 `tplUpdateById`，未逐行核）。
- 修法：评星只发 `{id,isContractQualified,isQualityQualified}`；其余只带改动字段。

### [🟡] 煤种审核失败即永久遮罩；无煤棚照片必崩
- 位置：varietyList.vue:1979、2008、1943-1957、2114-2128、1922-1927
- 证据：`this.mzloading = true;` … `this.editData.coverId = this.listtwo[0].compressedId;`
- 场景：审核没有煤棚照片的煤种（导入/老数据）→ `listtwo` 为空抛 TypeError，请求发不出而遮罩已亮；审核拒绝/重新审核遇 code≠200 或 HTTP 500 时 `mzloading` 不复位，关闭函数也不清，之后每个表单都是「更新中…」。
- 修法：请求前判 `listtwo.length`；catch 与关闭函数里复位。

### [🟡] loading 不复位，HTTP 500 后整页一直「加载中」
- 位置：labeltag.vue:150-163、Produlist.vue:162-175、coaltoper.vue:163-175、ordergoods.vue:111-124、procureList.vue:195-210（reject 才卡）；Qycensus.vue:396-413、userasp.vue:1180-1200、varietyList.vue:2328-2443（code≠200 也卡）；EditQycensus.vue:1370、1424、1486、1570 的 OCR `uploading`
- 证据：`this.loading = true; … .then(res=>{ if(res.data.code == 200){ this.loading = false;`
- 场景：接口 5xx 后遮罩不消失；EditQycensus 的 `uploading` 盖住含「保存/取消」的整张表单。
- 修法：then 内先复位再判断，补 `.catch`。

### [🟡] sessionStorage 的 queryAll：后退进入崩溃，离开不清理
- 位置：EditQycensus.vue:1162-1164（写入方 Qycensus.vue:434）
- 证据：`this.AddData = JSON.parse(sessionStorage.getItem('queryAll'));`
- 场景：保存/取消后浏览器后退回本页，key 已删 → null，mounted 抛错、页面残缺。经侧边栏或后退离开则残留：另有 19 个页面写同名 key，Mindex/AddPurchaseBill.vue:792、AddPurchaseContract.vue:1040-1042 直接把它当 `AddData`，新建页会读到企业对象。
- 修法：判空则跳回 Qycensus；`beforeRouteLeave` 清 key。
- 把握：崩溃 high；串数据 medium（只看了 grep 行）。

### [🟡] 先请求再重置页码；空数组被当成「无数据」
- 位置：labeltag.vue:182-186、Qycensus.vue:301-304/440-443/460-461、userasp.vue:879-881/915-917/1251-1252、varietyList.vue:1584-1586 等（9 个文件约 45 处）；userasp.vue:1236-1239、1268-1271 不重置
- 证据：`this.getpurcon();` `this.current = 1;`；`if(res.data.data.records != ""){`（Qycensus.vue:414）
- 场景：停在第 2 页再搜索，请求仍带 `current:2`，结果只有一页时显示「暂无数据」；`[] != ""` 为 false 走 else，`tableData`/`total` 不更新，翻到空页仍显示旧行。
- 修法：先置页码再请求；判空用 `Array.isArray`，else 清空。

### [🟡] 附件删除与表单值脱节，且保存前就物理删文件
- 位置：varietyList.vue:1725、1759-1764；EditQycensus.vue:1307、1342、1354
- 证据：`this.listtwo.map((item,key)=>{` `if(item == file.response.data[0].id){ this.listtwo.splice(key, 1);`
- 场景：`listtwo` 元素是对象，`item == id` 恒 false → 煤棚照片点 × 后文件已删，id 仍在 `imgList`，保存即裂图；删除接口失败时缩略图已消失而 `listarr` 仍有该 id；EditQycensus 官网图/视频点 × 即删服务器文件，随后「取消」则库里引用已删文件。
- 修法：比较 `item.id`；失败时回滚缩略图；旧文件改为保存成功后再删。

### [🟡] 「企业是否核实」筛选被静默丢弃
- 位置：varietyList.vue:2330-2383
- 证据：第 1、2 分支只有 `status:this.carriageName,` `isRelationEtc:this.caeName,`，`isRealCheck` 仅在 2403、2430
- 场景：同时选「待审核」+「已核实企业」→ 请求不带 `isRealCheck`，未核实企业混入结果。
- 修法：四个分支合成一个参数对象。

### [🟡] 用户页：编辑直绑行对象；清空日期提交 1970
- 位置：userasp.vue:1055、848
- 证据：`this.edituser = data;` `const today = new Date(time);`
- 场景：「编辑信息」改完点取消，列表行已被改写直到刷新；开通会员时清空到期时间（`@change` 传 null）→ `expirationTime` 成了 1970 年，提示开通成功却已过期。
- 修法：深拷贝后绑定；`changeTime` 开头 `if(!time) return`。

### [🟡] 企业编辑页：关联用户取消勾选无效，视频上传不回填
- 位置：EditQycensus.vue:909、1754-1756；1280-1282、1299-1301
- 证据：`if(!arr.includes(phone)){` `this.urcommontags.push(row)`；`catch (e) {` `this.uploading = false;` `return e;`
- 场景：`@select` 在取消勾选时也触发且只增不减 → 取消的用户点「应用」后仍被关联；官网视频 <5MB 分支不写 `advertisingVideoId`，保存后视频未更新；上传异常被 resolve，el-upload 把失败文件标成成功。另有 5 处成功回调不查 `response.code`（1373、1427、1489、1573、2033）。
- 修法：按 `selection` 增删；小文件分支补写并 `throw e`；回调首行判 code。

### [🟡] 首页图表：永久遮罩、时间轴错位、改年份不重算、实例泄漏
- 位置：PtIndex.vue:1366、1645、560、631-636、952-953
- 证据：缴税额图 `data:this.datapayment,`；`countMonDataTaxes` 里 `this.loading = true;` 之后无任何复位
- 场景：筛选只勾「客户数占比」(+「数据进度」)点确定 → 整页永久加载；取消勾选「业务量占比」再改区间 → 缴税额柱子对着旧月份；选「半年/季度/月份」后改年份 → 标题是新年份、查的是旧区间；`burdenLogins()` 每次挂载对已注释的 `#burdenLogins` 调 `removeAttribute` 抛错；每次刷新 `echarts.init` 新实例，无 dispose/`beforeDestroy`。
- 修法：该函数成功后复位 loading；横轴用 `paymentdatapayment`；`watch: yearValues` 重算；删 `burdenLoginsapi` 调用；init 前与销毁时 dispose。

### [🟡] 导入弹窗：定时器泄漏、类型校验无效
- 位置：varietyList.vue:2516-2523、2475-2490
- 证据：`let f = setInterval(()=>{` `this.perstep += 3;` `if(this.perstep == 99){`
- 场景：假进度期间点 × 或切路由 → 定时器不停，越过 99 后永不再提交；组件销毁后 `$refs.upload` 为 null，每 150ms 抛错。`on-change` 的返回值不拦截提交，非法文件照样传给 `importExcel`。
- 修法：句柄存 data，关闭函数与 `beforeDestroy` 清除；校验移到 `:before-upload`。

### [🟡] 「商品列表」页是标签页副本，增删改落在标签表
- 位置：Produlist.vue:37-53、172、207、250、293
- 证据：列是 `prop="userPhone"`、`prop="proName"` 等兑换字段，接口却是 `this.api.bustag.all(listByPageData)`
- 场景：表格每行全空；在此页「删除」删掉的是煤种正在用的标签。
- 修法：确认是否保留该页；保留则换积分商品接口，否则下线路由。

### [🟢] 静态签名硬编码进图片/下载/导出 URL
- 位置：Qycensus.vue:416；procureList.vue:142；EditQycensus.vue:1170、1224-1254、1292；varietyList.vue:917、1863、1884、2136（值不在此复述）
- 场景：`window.open` 导出不带登录态，只靠 URL 里的公开签名；后端 `auth/middleware.go:217` 又把它列入匿名白名单，全量煤种导出没有用户级鉴权。
- 修法：导出改走 axios 带 token；图片地址由后端下发；与 http.js 同一常量一并处理。
- 把握：存在性 high；可利用性 medium（取决于白名单范围）。

### [🟢] 零碎
- EditQycensus.vue:1073-1085：`/^[0-9]*[0-9][0-9]*$/` 至少要两位数，精煤仓数量填 `3` 报「数量格式不正确」→ 改 `/^\d+$/`。
- varietyList.vue:1702-1707：改大类不清 `threeProId`；1971-1974 只在取消时重置 `YuncangData`，下条煤种预填上条云仓号。
- labeltag.vue:252-257、Produlist.vue:264-269：`else if(res.data.code == 602)` 挂在 `if (valid)` 上，校验失败时 `res` 未定义；「确定」无防重，双击建两条标签。
- varietyList.vue:460、466：提货地/发站地误用 `v-if="scope.row.remark"`。

## 重复模式
1. `getpurcon(); current = 1;` 先请求后重置页码：9 个文件约 45 处。
2. loading 只在成功分支复位且无 `.catch`：全部有请求的文件。
3. 整行对象回传 updateById/updateBatch/edit：EditQycensus、varietyList、userasp、Qycensus（`Editisetc` 还先改行再请求）。
4. 静态签名拼 URL：4 个文件 13 处。

## 已核对无问题
- 页面调用的 api 方法在 `src/api` 均存在，均为 POST 对象并返回 axios 响应，与 `res.data.code` 读法一致。
- 无 `v-html`/`innerHTML`/`eval`/`javascript:`；PtIndex 自定义 tooltip 只拼日期与固定文案，含企业名的两张图用默认转义 tooltip。
- coaltoper、ordergoods、procureList 为只读列表，除上述模式外无异常；userasp「访问数据/兑换记录」抽屉入口已注释，属死代码。

🎉 值得保留：EditQycensus `ursearchtwo`（1825-1828）先置页码再请求，可作修模式 1 的范本；userasp `EditPurs`（1174）深拷贝隔离行对象。