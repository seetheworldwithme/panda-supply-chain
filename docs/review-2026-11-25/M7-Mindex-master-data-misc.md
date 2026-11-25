徐先生，Mindex 主数据与杂项页评审如下。

## 范围
28 个文件读完 template + script（22 个 `<style>` 按指示跳过）。另核对了 http.js、main.js、router 的 queryAll 守卫、相关 api 定义，以及后端 `settRecRelieve`、`stRelSmentsBills`。仓内无 node_modules，Element 行为未能查源码（见第 10 条）。

## 发现

### [🔴] 「解除关联」选中 id 只增不清且无确认，会解除已取消勾选的记录
- 位置：frontend/src/views/Mindex/Detailsment.vue:341-349
- 证据：
```js
purChange(val) {
    for(let i=0;i<val.length;i++){
        this.multipleSelection.push(val[i].id) 
    }
},
```
- 场景：勾选 A → 取消 A → 勾选 B → 点「解除关联」，请求 ids=[A,B]。后端 small.go:501-509 逐个清结算/发票的 relation 并删关联记录，未勾选的 A 被解除。成功后数组不清，下次不勾选直接点也会重发旧 id 并提示成功。
- 修法：`this.multipleSelection = val.map(r => r.id)`，成功后置空，加 `$confirm`。
- 把握：high。

### [🟡] 「关 联」无防重，双击建出两条关联记录
- 位置：Operation.vue:486-503、518-535
- 证据：`this.api.capital.relSmentsBills(this.selData)`（无 loading/disabled）
- 场景：路径不在 http.js 的 WRITE_RE 内。后端 settlement.go:691 每次调用插入一条 `chain_bill_sett_rec`，无事务/已关联校验，双击得到两条记录，后者覆盖 relation_id，前者成孤儿。
- 修法：WRITE_RE 加 `relSmentsBills|relieve|clearData`，或按钮加 loading。
- 把握：high。

### [🟡] 新建地址后跳回的页签与地址类型错位
- 位置：Addhair.vue:115；Address.vue:343-346、450、543
- 证据：`query:{sel:that.AddData.type}`；Address 中 `id: 1` 是「发货地」且请求 `type:2`，`id: 2` 是「收货地」且请求 `type:1`。
- 场景：新建「收货地」(type=1) 落到「发货地」页签，新记录不显示，易重复新建。
- 修法：传 1→2、2→1 映射后的值。把握：high。

### [🟡] 客户档案抽屉 5 个字段读的是供应商抽屉的数据
- 位置：Partner.vue:414-435
- 证据：`<span v-if="seeDatass.creditCode">{{seeDatass.creditCode}}</span>`（企业名称行用的是 `customerData`）
- 场景：先看供应商甲再看客户乙，乙抽屉显示甲的信用代码/经营地址/产品大类小类；没看过供应商则全是 ---。
- 修法：5 处改 `customerData`。把握：high。

### [🟡] 编辑页直接 JSON.parse(queryAll)，为 null 时整页报错；新守卫使 F5 必现
- 位置：EditAccount.vue:103、Editcustomer.vue:209、Editsupplier.vue:254、Editcollect 等四个地址编辑页：100；router/index.js:1244-1256（跨范围）
- 证据：`this.AddData = JSON.parse(sessionStorage.getItem('queryAll'));`
- 场景：`queryAllNavs` 初值 Infinity，编辑页按 F5 时首次导航即清掉 queryAll，`AddData=null`，模板取 `AddData.ownerType` 抛错（加守卫前 F5 可回填）。新标签页打开编辑 URL、离开后浏览器后退同理。
- 修法：created 判空则 `replace` 回列表页；守卫初值改 `0`。把握：high。

### [🟡] FileUpload：上传失败当成功、未传完可提交、无大小限制
- 位置：FileUpload.vue:85-93、290-311、334-373
- 证据：`<el-upload ... :on-success="resfile">`（无 on-error）；`if(response.code == 200){ ... this.enclosurefileList = fileList; }`
- 场景：① 服务端 code≠200 拒收时，文件仍带成功图标留在列表且无提示，确定时被静默漏掉；移除时 `file.response.data[0]` 遇 null 抛错。② 仍在上传时点确定，只提交已完成的文件，其余传成孤儿。③ 仅校验扩展名。
- 修法：code≠200 时移除并提示，补 on-error，拦截 `status==='uploading'`，加大小上限。把握：high。

### [🟡] 预览/下载 URL 拼入写死的签名凭据
- 位置：FileUpload.vue:219、376、381
- 证据：`'/saas-common/upload/previewStream?fileName=' + data.fileName + '&&AuthorizationId=…&&Signature=…'`（与 http.js 静态请求头同值，不复述）
- 场景：凭据进入浏览器历史与网关访问日志；fileName 未编码。
- 修法：随 http.js 静态凭据一并处置，至少 `encodeURIComponent`。把握：high。

### [🟡] 搜索/删除先发请求再重置页码
- 位置：Address.vue:472-476 等搜索与 4 个删除函数；Partner.vue:735-739、804-806；Contactsdetailed.vue:472-482、568-590；Operation.vue:656-672、722-738、886-902；FileUpload.vue:420-422
- 证据：
```js
this.hairtableData = [];
this.hairgetpurcon();
this.haircurrent = 1;
```
- 场景：停在第 2 页以后搜索一个只有 3 条结果的词，请求带旧页码返回空，页面显示「暂无数据」，分页器已到第 1 页却不重载。
- 修法：`current = 1` 提到调用前（Account.vue:141-145 写法正确）。把握：high。

### [🟢] 地址类型下拉标签错
- 位置：Editcollect.vue:87（`value: 4, label: '收货地'`）、Editpretend.vue:87（`label: '装货地'`）
- 场景：选第 4 项实际存成卸货地，记录从原列表消失。修法：改「卸货地」。把握：high。

### [🟢] current-change 回调未判 null
- 位置：Contactsdetailed.vue:427-444；Detailsment.vue:384-403
- 证据：`this.multiplesettlement = currentRow; if(this.multiplesettlement.type == 1){`
- 场景：列表刷新后原高亮行不在新数据中时，回调收到 null 即抛错，`multiplesettlement` 变 null，右侧列表再翻页时 `.systemCode` 再抛错、请求发不出。
- 修法：`if (!currentRow) return;`
- 把握：medium。否定证据：Element 2.15 `updateCurrentRowData` 若此时不发 null。

### [🟢] 产品大类→小类级联条件不严谨
- 位置：Addcustomer.vue:243、Editcustomer.vue:219-239、Addsupplier.vue:281、Editsupplier.vue:264-284
- 证据：`proArrs[i].indexOf(pro__middle_typeID) != -1`；`dvalue === this.AddData.pro1 * 1`；回显用 `id == 46`、新增用 `pID == 20`
- 场景：父 id 子串匹配会带出 11/12/21… 的小类；dvalue 为字符串时 `===` 永不成立，编辑页小类下拉为空；切换大类后 `pIDshow` 不复位。
- 修法：改 `==`，统一按 dvalue 判断。
- 把握：medium。字典 id 为长整型且 dvalue 为数字时，前两点不成立。

### [🟢] liveData「刷新」实为写操作
- 位置：liveData.vue:247-259
- 场景：`clearData` 不在 WRITE_RE，loading 只罩表格不罩按钮，可重复下发；成功后不调 `getpurcon`，清洗标记仍是旧值。
- 修法：开头判 `conversionloading`，成功后刷新列表。把握：high。

### [🟢] 零散
- Partner.vue:492 抽屉标题写死为某洗煤厂名称；:271 空态「新建」按钮无 `@click`；:260 `prop="createtime"` 拼错，排序无效。
- Contactsdetailed.vue:513、562 `paginationshow` 置 false 后不复位，右侧列表永久失去分页。
- 删除确认的 `.catch` 一律提示「已取消删除」，接口失败时也弹（Account/Address/Partner）。
- 死代码：Addcollect/Adddischarge/Addpretend 路由已注释，「保存」不发请求只提示成功；`this.splice(...)` 不可达；`this.tableData == []` 恒 false；`this.waterkeyupTime == null;` 无赋值。

## 重复模式
- 请求后置页码（上条 🟡）。
- 编辑页 `JSON.parse(getItem('queryAll'))` 7 处同构。
- `$confirm().then(api).catch('已取消')`。

## 已核对无问题
- 新增/编辑是独立页面，无共用表单残留；FileUpload 弹窗每次打开清空文件列表。
- add/updateById/delBatch 均为 POST+JSON 且命中 WRITE_RE，双击由网络层去重；loading 复位由 main.js mixin 兜底。
- 删除/取消核销/解除关联仅在 code==200 时提示成功；`hx_code` 与 `chain_capital_count.system_code` 字段对得上。
- 范围内无 v-html/innerHTML/eval；筛选器与依赖均已注册；liveData 无轮询/监听需清理；7 个占位页无逻辑。
- 🎉 列表页进编辑前都重新 `getById` 再写 queryAll，不回传列表行对象。
- 🎉 http.js 防重加 main.js 失败复位，使这批老页面免逐个补守卫，仅差上述 3 个未覆盖路径。