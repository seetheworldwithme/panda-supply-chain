# BusProductController → Go 迁移规格书（saas-warehouse :8765 `/busProduct`）

> 来源：`saas-warehouse/src/main/java/com/saas/warehouse/business/controller/BusProductController.java`（344 行，26 个端点）、
> `service/impl/BusProductServiceImpl.java`（2768 行，全文精读）、`mapper/BusProductMapper.java`（225 行，全文）、
> 实体 `saas-warehouse-core/.../model/BusProduct.java`（907 行）、`BusProLog.java`、`BusTag.java`、`BusPriceLog.java` 及被引用的
> mapper/service/util（StringHelper、ResultTool/JsonResult/ResultCode、EntityUtils、ExcelUtils、RedisUtil、RocketMqUtils、
> WarehouseConfig、BusPriceLogMapper、BusProLogMapper、BusScoreLogMapper、BusDelProductMapper、BusUserClickProMapper、
> ChainAuthUserMapper、Feign IChainAuthServiceFeign / IChainCommonFileServiceFeign、BusRedEnvelopeRecordServiceImpl、
> BusVideosServiceImpl、GlobalExceptionHandler）。MyBatis-Plus 版本 3.4.2。
> 本文档为唯一迁移依据：Go 侧照本文实现，不再读 Java。

---

## 0. 全局约定

### 0.1 响应封装 JsonResult（所有端点，除 exportExcel）
```json
{ "success": true|false, "code": 200, "msg": "成功", "data": ... }
```
- 成功：`ResultTool.success()`/`success(data)` → `{success:true, code:200, msg:"成功", data:<见各端点>}`。
- `ResultTool.fail(code,msg)` → `{success:false, code:<code>, msg:<msg>, data:null}`。
- `ResultTool.fail(data)`（仅 judgeIsFollow 用到）→ `{success:false, code:999, msg:"失败", data:<data>}`。
- 任何未捕获异常 → `GlobalExceptionHandler`：**HTTP 500** + body `{success:false, code:999, msg:"失败", data:"系统发生未知错误，请稍后重试。"}`。下文所有"→999"均指该响应。
- 字段名 camelCase；`data:null` 时照常输出 `"data": null`（Jackson 默认含 null）。

### 0.2 jwtInfo（网关注入 `request.getAttribute("jwtInfo")`）
字段：account/username/userId/companyId/enterpriseCode（全部 String）。
**本控制器只用 `userId`**（作为 create_user_id、SQL 里的 create_user_id/business_id/area_buser_id、bus_pro_log.user_id、积分/日志归属、redis key `userId{userId}`）。
`companyId`、`enterpriseCode` 在 BusProduct 全部 26 个端点中从未使用（Go 侧不需要）。
`userId` 在 SQL 中一律字符串拼接（`${}` 注入点，见 0.5）。部分端点取出 jwtInfo 但不用，Go 侧可忽略。

### 0.3 实体序列化（BusProduct，DB 列 217 字段）
- 全字段 camelCase，null 字段照常输出 `"x": null`。
- `createTime`、`updateTime` 有 `@JsonFormat(pattern="yyyy-MM-dd HH:mm:ss", timezone="GMT+8")` → 字符串 `"2024-01-02 03:04:05"`。
- 其余 LocalDateTime（fpriceUpdTime/ppriceUpdTime/lastShareTime/lastCheckTime/downTime）**无** @JsonFormat → Jackson JavaTimeModule 默认 ISO：`"2024-01-02T03:04:05"`（带 `T`，不带时区）。
- LocalDate（lastTestDate/purchStartTime/purchEndTime）→ `"2024-01-02"`。
- BigDecimal → 数字原样（如 `550.00`）。Boolean → true/false。Integer → 数字。
- `@TableField(exist=false)` 的附加字段（searchStr/page/ids/warehouseVideo/…/prags 等约 60 个）是实体的一部分，非 null 时一并序列化；null 输出 null。
- `tags` 字段：仅当 mapper 方法带 `@Results(@Many selTags)` 时填充（List<BusTag{id,signName,createTime,updateTime,remark1..5}>，无标签时空数组 `[]`）；不带 @Results 的查询里恒为 null。

### 0.4 MyBatis-Plus 语义（Go 侧等价物对照）
- `getById(id)` → `SELECT * FROM bus_product WHERE id=?`（单行或 null）。
- `updateById(entity)` → `UPDATE bus_product SET <非null字段>=? , update_time=now() WHERE id=?`（**null 字段不更新**；update_time 由 INSERT_UPDATE fill 自动写）。
- `save(entity)` → INSERT 全部非null字段（createTime 由 INSERT fill、updateTime 由 INSERT_UPDATE fill 自动写）。
- `saveBatch/updateBatchById` → 批量 insert / 按id批量 update（updateBatchById 同 updateById 语义，按非null字段更新）。
- `removeById/removeByIds` → 物理 DELETE。
- `count(QueryWrapper)` → `SELECT COUNT(*)`。
- `list(QueryWrapper)` → `SELECT *` + 条件。
- `allEq(map,false)` → 每个 key `col = val`，val 为 null 的跳过。
- `like(col,v)` → `col LIKE '%<v>%'`；v 为 null 时字符串拼接成 `LIKE '%null%'`。
- `Page` 序列化（IPage 返回值）：`{"records":[...],"total":n,"size":n,"current":n,"pages":n,"orders":[],"optimizeCountSql":true,"searchCount":true,"maxLimit":null,"countId":null}`（MP 3.4.2 形状；`orders` 来自 body 里 `page.orders`，本项目前端不传时为 `[]`）。
- `@SqlParser(filter=true)`：跳过多租户 SQL 解析器；本服务实际未启用租户拦截，Go 侧可忽略该标记。

### 0.5 SQL 注入点与参数化点总表（${}=拼接注入，#{}=预编译参数）
- `${searchSql}`：getSql 产出的条件串。其中 like 条件全部 `#{maps.x}` 参数化；但服务层直接字符串拼接的条件是注入点：
  - listByPage：`min/max_calorific_value`、`max_water_content`、`max_ash_content`、`max_std`、`max_volatiles`、`min_csr`、`min_bond_index`、`min_y`（数值拼接）、`p.fst_pro_id/p.big_pro_id/p.three_pro_id`、`is_real_check`、`province_code/city_code/area_code`（**字符串无引号拼接**，非数字直接 SQL 语法错误→999）、`p.company_id/group_id`（集团查询）、`etc_rel_phone`（redis 取 phone 拼接）。
  - warelist：`fst_pro_id/big_pro_id/three_pro_id` 数值拼接。
  - `'%${maps.searchStr}%'`：searchStr 直接进 SQL（listByPage/listNoPage/selCheckPro 分支）。**注意 getMaps 已把 searchStr 包成 `%v%`，模板里再包一层 `'%...%'` → 实际生成 `LIKE '%%v%%'`，语义与 `LIKE '%v%'` 等价**（Go 直接实现 `'%'+v+'%'` 即可）。
  - `${maps.userId}`（relListByPage/relListNoPage 的 join 条件，值=jwt userId）、`${proId}`、`${priceType}`、`${warehouseNb}`、`${id}`、`${leftDate}`、`${proWhere}/${psWhere}/${etcWhere}`、`${where}`（selWarseCode）均为拼接。
- `#{maps.userId}/#{maps.searchStr}/#{maps.enterpriseName}/#{maps.thName}/#{maps.*StartDate}` 等为参数化。

### 0.6 写操作无事务
本服务无 `@Transactional`。多表写（如 delBatch 的 bus_del_product→bus_product→bus_pro_log，updateById 的 bus_product→chain_common_file(feign)→bus_pro_tag→bus_common_log→MQ）**部分失败即部分提交**，Go 单体保持等价（不需要补偿；顺序照抄）。

### 0.7 依赖设施
- Redis（单机，db0）：key 见 5 节。
- RocketMQ 5.x（5.1.4）：topic 由 `WarehouseConfig`（prefix `config.*`）提供：`syncTopicName`（tag `syncOtherInfo`）、`noticeTopicName`（tag `updPro`/`addPro`）。发送后写 `bus_mq_consume_log` 一条（ip,topic,tag,now,messageId,body）。
- 微信公众号/小程序 HTTP、unipush HTTP、短信（SendCommonMsg，模板号硬编码）→ 全部登记为【外部依赖】。
- Feign（Go 单体内同库直查）：`saas-auth`（/chainAuthUser/getById、/selByIds、/list、/chainAuthChat/getAccessToken）、`saas-common`（/chainCommonFile/updateById、/updBthById、/list、/add）。返回的 data 是实体 JSON（camelCase）或 DB 行 Map（snake_case，见各端点）。
- HanLP 分词 + 自定义词典（"主焦"、"1/3焦"）【外部依赖，见 3.2 分词分支】。
- 本地文件 `/home/ubuntu/uploadfiles/`（frame73.jpg 背景图、生成图片落盘）。

---

## 1. StringHelper.getSql / getMaps 精确语义（决定 listByPage/cuslist/warelist 条件面）

### 1.1 getSql(T t, String unstr)
遍历 **BusProduct 全部声明字段**，凡带 `@ExcelColumn`（≈215 个业务字段+id，注解 value 可为空串）且满足：
1) 字段值非 null；2) 字段名不在 unstr 串里；3) **值≠字符串"0"**（`!"0".equals(f.get(t))`——只排除 String 类型的 "0"；**Integer 0 不会被排除**，会生成 `%0%` like 条件）
→ 追加 ` and <snake(字段名)> like #{maps.<字段名>}`。
返回 `"where 1=1" + 上述条件`。

BusProduct 的 @ExcelColumn 字段全集=全部表列（id 也带 @ExcelColumn(value="")，**body 里 id 非 null 时会生成 `and id like #{maps.id}`**，照抄）。exist=false 的查询辅助字段（searchStr/page/ids/warehouseType/…）**没有** @ExcelColumn，不进 getSql，但进 getMaps。

### 1.2 getMaps(T t)
遍历**全部声明字段**（不筛注解），非 null 的放入 map（key=驼峰字段名）：
- `Integer` 或 `LocalDate` 类型 → **原值**；
- 其余一切类型（String/Boolean/BigDecimal/LocalDateTime/对象/集合）→ **`"%"+toString()+"%"` 包裹**。
即 `searchStr` 放进去时已是 `%v%`。

### 1.3 各端点的 unstr
- listByPage：`"fstProId,bigProId,threeProId,minCalorificValue,maxCalorificValue,minWaterContent,maxWaterContent,minAshContent,maxAshContent,minStd,maxStd,minVolatiles,maxVolatiles,minCsr,maxCsr,minBondIndex,maxBondIndex,minY,maxY,isPlatformCertification,isRealCheck,provinceCode,cityCode,areaCode,groupId"`（这 25 个不走 getSql，由服务层手工拼，见 3.2）。
- warelist：`"fstProId,bigProId,threeProId"`。
- selCheckPro：`getSql(Map<String,Object> params, Class)`——遍历**请求 map 的 key**，key 能在 BusProduct 中找到**同名且带 @ExcelColumn** 的字段 → ` and snake(key) like #{maps.<key>}`（值不包裹，LIKE 精确值）；否则忽略。`searchStr` 无 @ExcelColumn → 只进 `'%${maps.searchStr}%'` 拼接。

---

## 2. 共享内部逻辑（多端点引用）

### 2.1 updUserScore(busProduct, getType, scoreType, exchangeLog) —— 积分发放
表：`bus_score_rule`（list().get(0)，**空表→IndexOutOfBounds→999**）、`bus_user_score`（user_id 定位）、`bus_score_log`（插入流水）。
- 若 `busProduct.status==2`（审核拒绝）：取 `getFreezeScore(userId,proId,3)`（`select ifnull(score,0) from bus_score_log where user_id=? and pro_id=? and score_type=2 and score_get_type=3 limit 1`）→ `bus_score_rule.mo_get_sum_score += 该值` 并 update；`bus_user_score.dj_hjj_num -= 该值` 并 update；插一条 bus_score_log（content="审核未通过，返回冻结积分", score=0, rule_type=4, score_type=1）；**return 0**。
- user_score 不存在 → 新建（hjj_num=0, dj_hjj_num=0）稍后 insert。
- getType 分支（本控制器触达 2/3/5/6）：
  - 2"发布煤种"：updScore=min(rule.publish_coal_score, rule.mo_get_sum_score)（getUpdScore 上限截断）；`dj_hjj_num+=updScore`（冻结）。
  - 3"更新煤种价格"：同上用 rule.upd_price_score，`dj_hjj_num+=updScore`。
  - 5"审核通过发布煤种"：freezeScore=getFreezeScore(...,2)；为 0 则 =rule.publish_coal_score 并插一条"发布煤种"流水(rule_type 1, score_type 2)；`hjj_num+=freezeScore`；`dj_hjj_num>0 ? dj_hjj_num-=freezeScore : rule.mo_get_sum_score-=rule.publish_coal_score`。
  - 6"审核通过更新价格"：同 5，用 getFreezeScore(...,3)/rule.upd_price_score（注意此分支判 `null==freezeScore||freezeScore==0`）。
- 月度上限段：tdSumScore/moSumScore 查询执行但结果仅记日志；`if(rule.mo_get_sum_score>0){ if(rule.mo_get_sum_score>=updScore) rule.mo_get_sum_score-=updScore; } else return 0;`（**mo_get_sum_score<=0 时直接 return 0，不落 user_score/log**）。
- 末尾：insert/update bus_user_score；插 bus_score_log（userId,userPhone,content,updScore,原可用,终可用,原冻结,终冻结,proId,ruleType,scoreType,getType）；update bus_score_rule；**return updScore**。

### 2.2 getUpdPriceScore(busProduct新值, bp旧值) —— updateById 端点积分
- `oldFac=selNewPriceLog(proId,1)`、`oldPla=selNewPriceLog(proId,2)`（`SELECT price FROM bus_price_log where pro_id=${proId} and price_type=${pt} ORDER BY create_time DESC LIMIT 1`）。
- isUpdFac = oldFac!=null && 新exFactoryPrice!=null；isUpdPla 同理。
- hasUpdPrice = (isUpdFac && 新出厂价.compareTo(oldFac)!=0) || (isUpdPla &&新车板价.compareTo(oldPla)!=0)。
- hasUpdPrice 且 旧status==1→新status==0：`updUserScore(bp,3,2,null)`（改价给冻结积分）；旧status==0→新status==1：`updUserScore(bp,6,1,null)`。返回该 score（无则 0）。

### 2.3 saveDelLog(pro, jwtInfo) —— 删除日志
remark 拼接（各段前导"，"最终去掉）：灰分Ad%≤max_ash_content、，水分Mt%≤max_water_content、，全硫St.d%≤max_std、，挥发分Vdaf%≤max_volatiles、，粘结指数G≥min_bond_index（null 段跳过）。
content = `"id为"+jwt.userId+"的用户删除了产品名称为"+pro.productName+"的煤种，相关指标为："+remark`。
INSERT `bus_pro_log`(user_id, pro_id, operate_type=3, operate_content=content)（构造器 BusProLog(userId,proId,3,content)）。

### 2.4 createNewCode(busProduct) —— warse_code 生成（getBusProduct 内）
- companyId==null 或 ∉{44,420} → `selWarseCode(" where warse_code like '0%'")`； companyId==44 → `like 'HY%'`；420 → `like 'KJ%'`。
- selWarseCode = `SELECT warse_code from bus_product <where> order by id desc limit 1`。
- `String.format("%05d", Integer.parseInt(warseCode)+1)`。**怪癖**：HY/KJ 分支拿到的码带字母前缀 → `Integer.parseInt` NumberFormatException→999；无行→null→NPE→999。只有 '0%' 分支数学上可行，照抄。

### 2.5 sendAdminMsg（saveBusProduct 内）【外部依赖-短信】
`bus_config` where msg_type=1 列表非空时，对每行：若 `config.sendMsgNotice`（WarehouseConfig.sendMsgNotice 开关）→ 短信模板 **2068872**，params=[busProduct.createPhone, "企业名 品名，{出厂价}元（出厂价）"]，发给 config.adminPhone。列表为空仅 log.error。

### 2.6 MQ 消费侧 syncOtherInfo（updateBusProduct 发出的消息的处理器，Go 单体需同库等价实现）
消息体 `"{proId},{userId},{isFirst},{orgStatus},{ts}"`，tag=syncOtherInfo：
- status==0（二次编辑）：selNewPriceLog 取旧价 → sendUpdAdminMsg（价差文案短信给平台管理员，模板 2068879，bus_config msg_type=1 每行一条，受 sendMsgNotice 开关）。
- status==1：updHisPrices（见 2.7）；企业核实分：count(bus_product where is_real_check=1 and company_id=?) → updCheckScore(0/1, companyId)（`update chain_basic_enterprise set coal_verify_score=${v} where id=${id}`）；`update bus_videos set status=1 where pro_id=?`。
- 通用：按 chain_basic_enterprise.is_authentication 刷 bus_product.is_certified_etc、group_id；按 bus_big_prodect 补 sec_name/th_name（仅当 big_pro_id 非空且 sec_name 为空）；updateById(bus_product)；status!=0 时 saveCheckResult（orgStatus=="0"：status=1→bus_msg(msg_type=3,"煤种信息审核通过,快去查看吧!",remark=proId,jumpId=proId)+bus_user_msg(create_user_id,msgId)+短信模板 1767390（受 sendMsgNotice）；status=2→同上文案"未通过"+updUserScore(bp,5,1,null)；orgStatus=="1"且 status=3→短信模板 2119567 给平台管理员）。
- updEtcScore(proId,companyId)（见 3.25）。
- isFirst=="true"：updUserScore(bp,5,1,null)（发布积分）；handRedEnvelopes（红包 type2：selTodayPubCount(create_user) 当日<3 才发；generateAmount 写 bus_red_money_clog → receiveAmount 领取【写 bus_red_envelope_record，受 config.openRedAmount/redStartTime/redEndTime/每日上限/redis clogId 锁控制】）；saveAddBroadcast（redis `broadcastMsg`（environmentType==3 时 `zscbroadcastMsg`）追加 `"新增{云仓号};"`，expireAt 当天 23:59:59+08:00）；pushAppMsg(type1)【外部依赖-unipush】；MQ noticeTopic `addPro`。

### 2.7 updHisPrices（status==1 时改价同步，消费侧）
- 出厂价：oldFac==null 且新价非null → bus_price_log 初始行(pre_price=0)；两者非null且不等 → bus_product.pre_ex_factory_price=旧值、fprice_upd_time=now、last_check_fac_price=新价；updFacHisPrice：selByPro 收藏列表、savePriceMsg（bus_msg msg_type=1 文本"{云仓号}|价格|{新价}|上升了X|元/下降了X|元"+指标 remark；bus_user_msg 给收藏者）、sendShareMsg（短信 2165435，redis `smsPhone{userId}` 当日去重，200 个一批，号码前缀"86"，受 sendMsgNotice）、rewardDistribution（redis `{proId}t8` 1h 去重；selMaxPriceLog 本周无记录→红包 type8）、pushAppMsg(type2)【unipush】、bus_price_log 插行(pre旧价,price新价,type1)；sendOfficeNotice 开→MQ noticeTopic `"proId-旧-新-ts"` tag updPro。
- 车板价：对称（type2；last_check_fac_price 仍写 exFactoryPrice——照抄这个 bug）；无出厂价变动且 sendOfficeNotice 才发 MQ。
- savePriceBroadcast：redis broadcastMsg（env3→zscbroadcastMsg）追加 `"{云仓号} {旧}元/吨 上涨到 {新}元/吨（出厂价/（车板价）);" `（精确格式：`warehouseNb+" "+old+"元/吨 上涨到 "+new+"元/吨"+（出厂价?"（出厂价）":"（车板价）"）+";"`），expireAt 当天末。

---

## 3. 26 个端点逐条规格

### 3.1 POST `/busProduct/list`（cuslist）【只读】
- 请求 body：`Map<String,Object>`。**key 必须是 bus_product 列名（snake_case）**，否则 SQL 报错→999。body 首先被 `maps.put("effective",1)` 覆盖 effective 键。
- 逻辑：`SELECT COUNT/SELECT * WHERE <每个 key: col=val 非null才拼> ORDER BY update_time DESC`。
  - 若 `three_pro_id` 存在且 `"888"` → 整体替换为 `SELECT * FROM bus_product WHERE status=1 AND effective=1 AND (company_id=420 OR group_id=420)`（kjPros，无 ORDER BY）。
  - 若 `three_pro_id`=="999" → 替换为 `... AND is_platform_certification=1`（authPros）。
  - 888 判定在前。tags 不查（null）。
- 响应：data = `List<BusProduct>`（SELECT * 全实体序列化）。
- jwtInfo：不使用。
- 怪癖：body 里的 `effective` 被强制改 1；`three_pro_id=888/999` 时其他条件全部丢弃。

### 3.2 POST `/busProduct/listByPage`（listByPage，jwtInfo 三参版）【只读】
- 请求 body：BusProduct 实体（可含全部查询辅助字段）+ `page:{current,size[,orders]}`（可缺省→null）。
- jwtInfo：userId。null 时 maps.userId=0。
- **条件组装顺序**（拼进 `searchSql`，起始 `where 1=1`）：
  1. getSql(t, 25 字段 unstr)——每个非null且非"0"串的 @ExcelColumn 字段 → ` and <col> like #{maps.<field>}`（参数化；getMaps 提供值：Integer/LocalDate 原值，其余 `%v%`）。
  2. 指标区间（**注意 min/max 错列 bug 照抄**）：minCalorificValue→`min_calorific_value is not null and !=0 and >= {min}`；maxCalorificValue→`max_calorific_value ... <= {max}`；minWaterContent→`max_water_content >= {min}`；maxWaterContent→`max_water_content <= {max}`；minAshContent→`max_ash_content >= {min}`；maxAshContent→`max_ash_content <= {max}`；minStd→`max_std >= {min}`；maxStd→`max_std <= {max}`；minVolatiles→`max_volatiles >= {min}`；maxVolatiles→`max_volatiles <= {max}`；minCsr→`min_csr >= {min}`；maxCsr→**`min_csr <= {max}`（bug）**；minBondIndex→`min_bond_index >= {min}`；maxBondIndex→`min_bond_index <= {max}`；minY→`min_y >= {min}`；maxY→`min_y <= {max}`。全部数值直接拼接（注入点）。
  3. `fstProId→ and p.fst_pro_id = n`；`bigProId→ p.big_pro_id`；`threeProId→ p.three_pro_id`。
  4. `isRealCheck` 非 null → ` and is_real_check = {1|0}`；provinceCode/cityCode/areaCode 非 null → ` and xxx_code = {裸值}`（字符串裸拼，非数字→999）。
  5. 日期区间（参数化）：createStartTimeDate/End→`DATE_FORMAT(create_time,'%Y-%m-%d') >=|<= #{maps.createStartTimeDate}`；updateStartTimeDate/End→update_time；facPriceStartTimeDate/End→fprice_upd_time。
  6. `isRelationEtc`：true→` and p.company_id is not null`；false→` is null`。
  7. `warehouseType==0` → ` and status = 1 and (is_certified_etc=1 or ware_identity=4)`。
  8. `warehouseType==1` → redis `HGETALL userId{jwt.userId}` 取 phone；` and (create_user_id = {uid} or business_id= {uid} or area_buser_id= {uid}[ or etc_rel_phone={phone}])`。jwtInfo 为 null → **NPE→999**；redis 无记录 → `or etc_rel_phone=null`。
  9. `warehouseType==2` → otherUserId（body）→ redis `userId{otherUserId}` 取 phone；` and status != 3 and (create_user_id = {ouid} or business_id= {ouid} or area_buser_id= {ouid}[ or etc_rel_phone={phone}])`。
  10. `isUpDown==true` → ` and (status = 1 or status = 3)`。
  11. `isPlatformCertification==true` → ` and is_platform_certification=1`。
  12. **isDirect 判定**（5 个 count，全部 `col LIKE '%{searchStr}%'`，searchStr null → LIKE '%null%'）：enterprise_name、th_name、warehouse_nb、warse_code、create_phone 任一 count>0，或 `searchStr.contains("公司")`（**searchStr 为 null 且全部 count=0 → NPE→999**）→ isDirect=true。
  13. `isGroupSearch` 非 null：true→` and (p.company_id = {groupId} or group_id = {groupId})`；false→` and ((p.company_id is null or p.company_id !=420) and (group_id is null or group_id != 420))`。
  14. `isShared/isforwared/isCared/isVisited ==true` → 分别 ` and exists (select 1 from bus_pro_log where operate_obj=4 and operate_type={4|5} and pro_id=p.id)` / ` and exists (select 1 from bus_pro_collection where pro_id=p.id)` / `... operate_type=7 ...`。
  15. **HanLP 分词分支（!isDirect 时）**：词典加"主焦"、"1/3焦"后 `HanLP.segment(searchStr)`。逐词：
      - 未匹配过企业名：count(`enterprise_name LIKE '%word%'`)>0 → ` and enterprise_name like #{maps.enterpriseName}`（maps 值=`%word%`），hasEtcName=true，continue。
      - 未匹配过煤种名：count(`th_name LIKE '%word%'`)>0 → ` and th_name like #{maps.thName}`，hasCoalName=true，continue。
      - word 能 Double.parseDouble：含"硫"→` and max_std ={下一个词}`；含"挥发"→` max_volatiles =`；含"水"→` max_water_content =`；否则 ` and (max_ash_content={w} or max_volatiles={w} or max_std={w} or max_water_content={w} or min_bond_index={w} or min_csr={w} or min_y={w} or lithofacies={w} or clean_coal_recovery_rate={w} or gangue={w})`。**"硫/挥发/水"词是最后一个词时 termList.get(i+1) 越界→999**。
      - 循环后 `if(termList.size()>1)`：page==null→`listNoPFenCi`；否则 `listByFenCi`，**return**。⚠️该分支 SQL **不含** CONCAT(searchStr) 全字段过滤——搜索词只有命中的 like 条件；全未命中时=无条件全表。
  16. **推荐煤种分支**（threeProId==null 且 jwtInfo!=null 且 isClientSearch==true）：`maps.userId=jwt.userId`；`selLastClickPro`（`select *,sumScore,isShowEtcName from bus_product p, bus_user_click_pro cp where p.id=cp.last_pro_id and status=1 and publish_type=1 and is_last_click=1 and user_id=${userId}`）非空 → ` and p.id != {lastPros[0].id}`，否则 `page.setSize(10)`。返回 `resMap={"lastPros":[...], "prolist": relListByPage(searchSql,maps,page)}`（⚠️page==null 时先执行 relListNoPage（结果被丢弃）再以 page=null 调 relListByPage——实践上必须传 page，Go 侧登记差异：无 page 时走 relListNoPage 并直接返回）。relListByPage SQL：`select p.*,(sumScore 子查 chain_basic_enterprise by company_id),(isShowEtcName/isShowEtcPhone 子查 bus_member_benefits×chain_auth_user by create_user_id) from bus_user_click_pro cp RIGHT JOIN bus_product p on cp.three_pro_id=p.three_pro_id and cp.user_id=${maps.userId} {searchSql} and CONCAT(IFNULL(warehouse_nb,''),region,product_name,resource_number,send_addres) like #{maps.searchStr} order by cp.is_last_click desc, cp.three_pro_id DESC, p.update_time desc`（参数化 like，值已含 %v%；**无 tags**）。
  17. 常规分支：page==null→`listNoPage`；否则→`listByPage`。
      - **listNoPage**：`select (warehouseVideo),(warehouseImg),(collCount #{maps.userId}),(isShowEtcName),(isShowEtcPhone),p.* from bus_product p {searchSql} and CONCAT(IFNULL(enterprise_name,''),warehouse_nb,region,product_name,resource_number,send_addres,fst_name,sec_name,th_name,warse_code,coal_code,create_time,update_time,create_time,update_time,business_name,area_buser_name,coal_rel_name,create_phone) like '%${maps.searchStr}%' order by status asc, update_time desc` + tags @Many。返回 List。
      - **listByPage**：头部再追加常量与子查询列：`select false as isDivider, false as loading, (chain_common_file where association_id=p.association_id and f_type=4 and is_compressed_file=0 的 count) mzkfyq, (bus_pro_log obj4 type7 target=create_user) visitedTotal, (DISTINCT user_id 版) vtPeopCount, (type4) shareTotal, (type5) fwdTotal, (type6) interestTotal, warehouseVideo, warehouseImg, collCount, sumScore, isShowEtcName, isShowEtcPhone, (bus_verify_info where pro_id=p.id 的 count) isVerify, p.* from bus_product p {searchSql} and CONCAT(同上19槽) like '%${maps.searchStr}%' order by status asc, update_time desc` + tags @Many。返回 IPage。
- 响应：data = List（无分页/分词分支）或 IPage 对象或 `{"lastPros":[],"prolist":{IPage}}`（推荐分支）。
- 【只读】。

### 3.3 POST `/busProduct/getById`（getBusProduct）【写操作：bus_product.association_id/warse_code、bus_user_click_pro、bus_visitor_record、bus_pro_log、bus_msg、bus_user_msg、bus_red_money_clog、bus_red_envelope_record（条件）】
- body：`{"id":n, "shareUserId":n?, "visitorResource":n?}`（id 必需）。
- jwtInfo：非 null 时触发访客逻辑（见下）；null 时跳过。
- 流程：
  1. `bp=getById(id)`；null → 控制器返回 `fail(602,"产品已下架")`。
  2. association_id 为 null → 生成 32 位 uuid 写回 bus_product（【写】）。
  3. `status==0 && warse_code==null` → warseCode=createNewCode(bp)（2.4，可能 999），**仅内存赋值不落库**。
  4. userMap：feign `/chainAuthUser/getById`（chain_auth_user 行 JSON，camelCase 键）；data null 时 userMap 保持 null（后续 1639 行 NPE→999）。
  5. companyId 非 null → etcInfo=`SELECT * FROM chain_basic_enterprise WHERE id=${companyId}`（**snake_case 键**）；groupId 非 null → groupInfo=selEtcById(groupId)；groupId null 且 companyId 非 null → `SELECT * FROM chain_basic_enterprise WHERE group_id=${companyId}` 有行时 groupInfo=selEtcById(companyId)。
  6. reportListMap = feign `/chainCommonFile/list` {f_type:2, is_compressed_file:0, association_id}；imgListMap = 同接口 {f_type:4, is_compressed_file:1}；videoList = `SELECT * FROM bus_videos WHERE pro_id=?`；videoListMap = feign {f_type:1, is_compressed_file:0}。
  7. warehouse_nb 非 null → areaProList=`select (thName 子查),p.* from bus_product where warehouse_nb like '{首段}%' and id != {id} and status=1 ORDER BY create_time`（首段=`warehouseNb.split(" ")[0]`）。
  8. bus_config where msg_type=3：非空→clientPhone=configs[0].adminPhone、clientName=**常量"李志强"**；⚠️随后 `configs.get(0)` 再取一次客服用户（feign `/chainAuthUser/list` {phone:configs[0].adminPhone, channel:2}）`.get(0)` → **无 msg_type=3 配置或查无此人 → IndexOutOfBounds→999**；businessHead=该用户 wxFaceImg（null→""）。
  9. brokerList：coalRelName 非 null → {username:coalRelName, phone:coalRelPhone, isBorker:false}；warehouse_nb 含"凯嘉" → 追加硬编码 {"username":"张慧庆","phone":"18634880999","isBorker":false}；businessId 非 null → feign getById → {…用户全部字段, borkerYears:Period(brokerStartDate,今天).getYears(), isBorker:true}（**brokerStartDate 缺失→parse 异常→999**）；areaBuserId 非 null 且 businessId 非 null 且二者不等 → 同上追加。**businessId 为 null 而 areaBuserId 非 null → `busProduct.getBusinessId().intValue()` NPE→999**。
  10. tags=selTags（`SELECT t.* FROM bus_tag t, bus_pro_tag pt WHERE t.id=pt.sign_id AND pt.pro_id=${proId}`）；无→[]。
  11. vtPeopCount = count(bus_visitor_record where pro_id=?)。
  12. thName=bus_big_prodect.getById(threeProId).productName（**三级分类缺失→NPE→999**）。
  13. busMemberBenefits = bus_member_benefits where member_type=userMap["memberType"]（getOne）。
  14. prags=[]。
  15. **jwtInfo 非 null（updByVisitor）**【写】：
      - collCount=0；bus_pro_collection where pro_id=? and user_id=jwt.userId 存在→collCount=1。
      - bus_user_click_pro where user_id=jwt.userId：无→INSERT(userId,lastProId=id,lastTypeId/ threeProId=pro.threeProId,clickTime=now,isLastClick=true)；有→UPDATE 同字段+isLastClick=true。
      - isCanCareReward：selTodayCare（`count bus_pro_log where user_id=? and pro_id=${proId} and operate_type=6 and operate_obj=4 and to_days(create_time)=to_days(now())`）≥3→false；该用户对该 pro 的 type=6 记录 count>0→false；否则 true。
      - lastShareTime=selLastShareTime（当日 type 4/5 obj=4 最新 create_time，可能 null）。
      - 分享红包：selTodayUser（`count chain_auth_user where id=${uid} and to_days(create_time)=to_days(now())`）非 null 且 body.shareUserId 非 null 且 jwt.userId != shareUserId → selVisitShare（`select * from bus_pro_log where user_id=? and share_user_id is not null and operate_type=7 and (operate_obj=4 or operate_obj=3)`）为空 → shareRedEnvelopes(shareUserId)：generateAmount(type5)→INSERT bus_red_money_clog(create_user=访客, receive_type=5, amount=区间随机, is_write_off=false)；receiveAmount→校验（config.openRedAmount/redStartTime~redEndTime、每日上限、redis clogId 1h 锁、is_write_off、receive_type 一致）→INSERT bus_red_envelope_record+核销 clog。
      - 访客记录：`chainAuthUserMapper.selectById(jwt.userId)`；`pro.createUserId != jwt.userId 且 visitUser.channel=="2"` → saveVisitInfo：INSERT bus_msg(msg_type=4, content="有1位客户浏览过您发布的煤种，快去查看吧!", remark=proId, jumpId=proId) + INSERT bus_user_msg(create_user_id=pro.createUserId, msg_id)；bus_visitor_record where visitor_phone=? and pro_id=? 无→INSERT(proId,visitorId=jwt.userId,username,nickname,phone,wxHeadImg=**nickname（bug 照抄）**,now,visitorResource) 有→UPDATE visitor_id/update_time；INSERT bus_pro_log(user_id,user_phone,pro_id,target_user_id=createUserId,operate_obj=4,operate_type=7,share_user_id=body.shareUserId,ware_identity=visitUser.wareIdentity)。visitUser 查无→NPE→999。
- 响应：data=完整 BusProduct（含 userMap/etcInfo/groupInfo/reportListMap/imgListMap/videoList/videoListMap/areaProList/clientPhone/clientName/businessHead/brokerList/tags/vtPeopCount/thName/busMemberBenefits/collCount/isCanCareReward/lastShareTime/prags）。

### 3.4 POST `/busProduct/add`（saveBusProduct）【写操作】
- body：BusProduct（表单全字段）+ reportList:List<ChainCommonFile>? + imgList:List<ChainCommonFile>? + videoList:List<BusVideos>? + proTags:List<BusProTag>? + ckhyFileName?。
- jwtInfo：userId → create_user_id、文件 ownerId、日志、积分归属。
- 校验（publishType==null 或 ≠2 时）：
  - wareIdentity、coalType、bigProId、threeProId、(exFactoryPrice 或 platePrice 至少一)、coalRelName、coalRelPhone 任一为 null → `fail(602,"煤种身份、煤种区别、煤炭分类、煤炭品种、出厂现汇价或车板价、联系人姓名、联系人电话、煤棚照片为必填项")`。
  - wareIdentity≠4 且 companyId==null → `fail(602,"请选择企业")`。
- companyId 非 null：etcMap=selEtcById(companyId)；enterpriseName=etcMap["enterprise_name"].toString()（**企业不存在→NPE→999**）；group_id 非 null→groupId；is_authentication 为 "1" 或 "true" → isCertifiedEtc=true。
- associationId=32位uuid；createUserId=jwt.userId。
- reportList 非空：每条 set fType=2,ownerType="1",ownerId=jwt.userId,relationId=busProduct.id(**add 时为 null**),relationType=1,effective=true,associationId → feign updateById 逐条【外部-feign；实际写 chain_common_file，relation_id 为 null】；isHasReport=true。imgList 同（fType=4）逐条 feign updateById。
- preExFactoryPrice=exFactoryPrice；prePlatePrice=platePrice。
- secName=bus_big_prodect[bigProId].productName、thName=bus_big_prodect[threeProId].productName（**分类 id 无效→getOne null→NPE→999**）。
- INSERT bus_product（status 等未赋值字段由 DB 默认；create_time/update_time fill）。
- videoList：每条 video.proId=新id、video.userId=createUserId → saveBusVideos(video,uuid,1)：videoId 为 null→`fail(602,"视频上传失败，请重试")`（**但外层不检查返回值，继续执行**）；feign 查用户 nickname 写入；INSERT bus_videos；feign `/chainCommonFile/updateById`（videoId 对应文件 fType=1? 此处 fType 参数=1, ownerId, relationId=proId, effective=true, associationId）。
- ckhyFileName 非 null：feign updateById（ChainCommonFile 只带 fileName+effective=true，**id 为 null → MySQL `WHERE id=NULL` 匹配 0 行，实际无操作**，照抄为 no-op）。
- proTags 非 null：saveBatch INSERT bus_pro_tag。
- 积分：publishType==null 或 ==1 → score=updUserScore(bp,2,2,null)（冻结"发布煤种"）；publishType==2 → score=0。
- sendAdminMsg（短信 2068872，2.5）。
- 日志：feign 无（本地）`SELECT * FROM chain_auth_user WHERE id=jwt.userId`（chainAuthUserMapper.selectById）→ INSERT bus_common_log(user_id,username,phone,pro_id=新id, content="新增了煤种，新增信息为：\n"+getCustAddLog(实体)（逐 @ExcelColumn 非null字段 `标签：值\n`，标签取注解值"（"前段，cust.* 环境变量映射可忽略）, operation_type=1)。
- 响应：data=score（Integer 标量）。
- 【外部依赖】短信、feign saas-common 文件更新、（视频时）feign saas-auth。

### 3.5 POST `/busProduct/deleteById`（delete）【写操作】
- body：`{"id":n}`；jwtInfo：userId 用于日志。
- 顺序（无事务）：① `INSERT INTO bus_del_product SELECT * FROM bus_product WHERE id=${id}`（bus_del_product 备份）；② `getById(id)`（备份后查询，仍能查到）；③ `DELETE FROM bus_product WHERE id=?`；④ saveDelLog→INSERT bus_pro_log(type=3)。pro 为 null（id 不存在）→ ①空插 ②null→saveDelLog NPE→999（前两步已生效，无回滚）。
- 响应：`success("删除成功，id为："+id)` → data=字符串。

### 3.6 POST `/busProduct/updateById`（updateBusProduct）【写操作；外部：MQ/RocketMQ、feign、短信(无)、积分】
- body：BusProduct（部分字段）+ reportList/imgList/proTags/isForceUpdPrice/effective 等；jwtInfo：userId。
- 逻辑顺序：
  1. `bp=getById(id)`；**null → bp.getAssociationId() NPE→999**。uuid=bp.associationId；null→生成并放进 busProduct（随 update 落库）。
  2. reportList 非空：每条强制 fType=2/ownerType="1"/ownerId=jwt.userId/relationId=id/relationType=1/effective=true/associationId → feign `/chainCommonFile/updBthById` 批量；isHasReport=**true**。reportList 为空/null → isHasReport=**false**（显式落库）。
  3. updateTime=now。
  4. createPhone 非 null 且 `!bp.createPhone.equals(新)`（**bp.createPhone 为 null→NPE→999**）：selectWxUser（`chain_auth_user where (phone=? or account=?) and channel=2 order by level desc`）有→createUserId=users[0].id；无→INSERT chain_auth_user(phone=account=新手机,channel="2") → createUserId=新id（【写】auth 库用户表）。
  5. status：bp.status==0 && 新 status==1 → lastCheckTime=now；bp.status==0 && 新 status==3 → downTime=now（`bp.getStatus()==0` Integer 拆箱，bp.status null→NPE→999）。
  6. `updateById`（非null字段+update_time）。
  7. isForceUpdPrice==true：exFactoryPrice==null → `update bus_product set ex_factory_price=null where id=${id}`（forceUpdExPrice）；platePrice==null → forceUpdPlPrice。
  8. `effective==0`（Integer，null 安全）→ **return 0**（跳过后面全部）。
  9. imgList 非 null：每条 fType=4/ownerType="1"/ownerId/relationId/relationType=1/effective=true/associationId → **逐条** feign updateById。
  10. proTags 非 null：DELETE bus_pro_tag where pro_id=? → saveBatch 插入。
  11. 日志：chain_auth_user selectById(jwt.userId) → INSERT bus_common_log(…, content="修改了煤种："+bp.thName+"\n"+getCustChangeLog(bp,新,"cust."), operation_type=2)。getCustChangeLog：遍历新值非null的 @ExcelColumn 字段，BigDecimal 用 compareTo，其他 equals，输出 `标签：旧值——》新值\n`。
  12. isFirst 判定：`新status==1 && bp.status==0 && bp.lastCheckTime==null` → "true"。
  13. 【外部依赖-MQ】`sendNormalMsg(syncTopicName, "{id},{jwt.userId},{isFirst},{bp.status},{当前毫秒}", "syncOtherInfo", "syncOtherInfo{id}")`；发送后写 bus_mq_consume_log。消费侧行为见 2.6/2.7（Go 单体内同步/异步等价，需登记）。
  14. score=getUpdPriceScore（2.2，可能触发积分写 bus_user_score/bus_score_log/bus_score_rule）。
- 响应：data=score（Integer；effective=0 时 0）。

### 3.7 POST `/busProduct/importExcel`（importExcel+saveIrt）【写操作】
- 请求：multipart `file`。jwtInfo 不用；**用户身份取 Spring Security principal UserInfo**（Go 侧：当前登录用户 id/phone）。
- importExcel：`ExcelUtils.readExcel(BusProduct.class,file,null,0)`：按扩展名 xls/xlsx 解析 sheet0；**表头行=各字段 @ExcelColumn value 文本**，逐行映射字段（数字/文本），必填校验失败进 errInfos；文件空→errInfos=["上传文件为空"]；扩展名不对→["上传文件格式不正确"]；解析异常→["解析 Excel 异常，请检查文件格式和内容"]。返回 ImportResult{dataList,errInfos}。
- saveIrt：batchNumber=UUID（**未使用**）。errInfos 非空 → 逐条前缀 `"产品sheet中，"+s`，errs=["导入失败 0"]，**不导入**。errInfos 空：
  - `SELECT * FROM bus_warehouse` 建 warehouseNb→id 映射；
  - 逐行：warehouseNb 不在映射 → errinfos.add("云仓号不存在："+nb)（**该列表最终未写入返回值——丢失，静默**）；warehouseId=映射.get(nb)（不存在时为 null 照插）；createUserId=user.id；createPhone=user.phone；**status=1**。
  - `saveBatch` 全量 INSERT bus_product；sucInfos=["云仓产品管理数据成功 N"]（N=dataList.size）。
- 响应：data=ImportResult 序列化 `{dataList:[导入的实体], sucInfos:[...], errs:[...], errInfos:[...]}`。
- 怪癖：云仓号校验失败不阻断、不报错、照插 warehouse_id=null；无事务。

### 3.8 GET `/busProduct/exportExcel`【只读】
- 无 body/jwtInfo。`SELECT * FROM bus_product` 全表 → `ExcelUtils.writeExcel`，xlsx 附件下载，文件名"导出数据"，Content 为 Excel 流（**无 JsonResult**）。Go 侧：全表导出 xlsx 即可。

### 3.9 POST `/busProduct/delBatch`【写操作】
- body：`{"ids":[...]}`；jwtInfo：userId 进日志。
- 逐 id：① INSERT INTO bus_del_product SELECT * WHERE id=${id}；② getById(id)；③ saveDelLog（INSERT bus_pro_log type=3）。循环后 `DELETE FROM bus_product WHERE id IN (ids)`。无事务；pro null→999（前序已生效）。
- 响应：`success("删除成功")` data=字符串。

### 3.10 POST `/busProduct/warelist`【只读】
- body：BusProduct；jwtInfo 不用。
- searchSql=getSql(t,"fstProId,bigProId,threeProId")；再拼 ` and fst_pro_id = n / big_pro_id = n / three_pro_id = n`（非null时，数值裸拼）。
- SQL：`select (fstName),(secName),(thName),(warehouseVideo: bus_warehouse by warehouse_id),(warehouseImg),bus_product.* from bus_product ${searchSql} and status=1 GROUP BY warehouse_nb`（**无 ORDER BY；GROUP BY 任取每组一行（MySQL 惯例取组内第一扫描行）**）。无 tags。
- 响应：data=List<BusProduct>。

### 3.11 POST `/busProduct/buslist`【只读】
- 无参数。`SELECT product_name from bus_product GROUP BY product_name`。
- 响应：data=List<BusProduct>，每项仅 productName 有值。

### 3.12 POST `/busProduct/hisPrice`【只读】
- body：`{"id":n, "priceType":1|2|3, "hisPriceType":1|2}`。
- 控制器：`hisPriceType==null` → `fail(802,"参数异常")`。
- 算法：
  1. bp=getById(id)（null→NPE→999）。
  2. cal=Calendar(现在)；ldate=今天。priceType==1→cal-7天、ldate-7天、prices=hisPriceWeek（`TO_DAYS(NOW())-TO_DAYS(create_time)<=7`）；==2→-30天、hisPriceMonth；==3→-5月、hisPriceHYear（`create_time between date_sub(DATE_ADD(curdate(),interval -day(curdate())+1 day),interval 5 month) and now()`）。三条 SQL 均：`select * from (select * from bus_price_log where pro_id=${id} and price_type=${hisPriceType} [窗口] HAVING 1 ORDER BY create_time desc) n GROUP BY date_format(create_time,'%Y-%m-%d') ORDER BY create_time`（周/月外层升序；半年 SQL 无外层分组，直接升序）。**priceType 缺失或 ∉{1,2,3} → prices 保持 null → `prices.size()` NPE→999**。
  3. 若 bp.createTime 在 ldate 之后：cal=createTime、ldate=createTime 日期。
  4. prices 非空：若 prices[0].price≠prices[0].prePrice 且 prePrice.intValue()!=0 且 priceType∈{2,3} → pMap[ldate]=prices[0].prePrice、xlist+=ldate。遍历 prices（升序）：time=createTime 前 10 字符；type∈{2,3} 且 xlist 未含→xlist+=time；pMap[time]=price。type≠3：若 type==1 → xlist=pMap 全部 key（TreeMap 字典序）。
  5. prices 空：price = hisPriceType==1?bp.exFactoryPrice:bp.platePrice；非 null → xlist=[cal 格式化"yyyy-MM-dd"]、pMap[ldate]=price。
  6. 返回 `{"priceData":{TreeMap按key字典序}, "xlist":[...]}`。
- 响应：data=该 Map。怪癖：priceData 为 TreeMap（**key 日期字符串排序输出**）；pMap 的 key 与 xlist 的内容可能不一致（照抄）。

### 3.13 POST `/busProduct/shareTime`【写操作】
- body：`{"id":n, "busProLog":{BusProLog}?}`；jwtInfo：仅判空。
- getById(id) 为 null → **直接 return（无更新）**，仍 success。非 null → sharingTimes=(null?0:+1)；`updateById`（全非null字段+update_time 回写）。
- jwtInfo 非 null 且 body.busProLog 非 null → INSERT bus_pro_log（**body 原样字段直插，包括 user_id/pro_id/operate_obj/type 等**）。
- 响应：success()，data=null。

### 3.14 POST `/busProduct/selCheckPro`【只读】
- body：`Map<String,Object>`（key=BusProduct 驼峰字段名）；jwtInfo 不用。
- searchSql=getSql(maps,BusProduct.class)（1.3：命中 @ExcelColumn 的 key → ` and snake like #{maps.key}`，**值不包裹=精确值 LIKE**）。
- SQL：`select (warehouseVideo),(warehouseImg),(collCount #{maps.userId}),(isShowEtcName),(isShowEtcPhone),p.* from bus_product p ${searchSql} and CONCAT(19槽同 3.2) like '%${maps.searchStr}%' order by status asc, update_time desc` + tags @Many。`maps.userId` 未传→`collCount` 子查询 user_id=NULL→0。
- 响应：data=List<BusProduct>。

### 3.15 POST `/busProduct/addViewTime`【写操作】
- body：`{"id":n}`。控制器内直接：getById(id)（**null→NPE→999**）；viewTimes=(null?0:+1)；updateById（非null字段全量回写+update_time）。
- 响应：success() data=null。无 jwtInfo。

### 3.16 POST `/busProduct/addPhoneTime`【写操作】
- body：`{"id":n, "phoneLog":{BusPhoneLog}?}`；jwtInfo 取出不用。
- getById(id)→null→NPE→999；phoneTimes+1→updateById。
- body.phoneLog 非 null：feign `/chainAuthUser/list` {phone:phoneLog.targetUserPhone, channel:2}；有结果→phoneLog.targetUserId=users[0].id。
- `busPhoneLogService.save(pl)`：**pl 为 null 时 MP insert(null) 抛异常→999**；非 null→INSERT bus_phone_log（body 字段直插）。
- 响应：success()。

### 3.17 POST `/busProduct/getDetailImg`（getDetailImg2）【写操作：bus_product 无；chain_common_file（feign）、redis、本地文件】【外部依赖：feign、redis、文件系统】
- body：`{"id":n}`；jwtInfo 不用。
- 全程 try/catch：**任何异常→返回 null→响应 success(null)**（不报错）。
- 流程：pro=getById(id)；`middPrice=exFactoryPrice - preExFactoryPrice`（任一 null→NPE→被吞→null）；fuhao=▲/▼，priceZf=fuhao+|差|+"元"；以 `/home/ubuntu/uploadfiles/frame73.jpg` 为底合成 JPG：云仓号(微软雅黑粗体40px @210,130 黑)、日期 yyyy/MM/dd(30px @440,50 白)、企业名（**仅当创建者会员权益 isShowEtcName**，25px @90,230）、联系人/电话（isShowEtcPhone？真实值 : "***"/"***********"，@280,751/793 黑）、"出厂现汇价：{价}元"(25px @90,275)、middPrice≠0 时"较上次：{▲/▼x元}"(25px @380,275)、businessName(25px @340,1016 白)、businessPhone(20px @280,1060 白)、指标两列布局（灰分/水分/硫/挥发/粘结/Y/CSR/岩相，非null项按 idx 奇偶左右列 28px，左列 x=90,y=330 起 y+=55，右列 x=370,yr=330 起 yr+=55）。合成结果 BufferedImage→Base64。
- 落盘 `/home/ubuntu/uploadfiles/{毫秒时间戳}.jpg`；feign `/chainCommonFile/add`（ChainCommonFile(fileName,fileName,filePath="/home/ubuntu/uploadfiles",fType="5",ownerType="2")）；`redis SET {fileName} = base64`。
- 响应：data=fileName 字符串（如 `"1731812345678.jpg"`），异常时 null。

### 3.18 POST `/busProduct/updateBatch`（updateBatchById）【写操作】
- body：**JSON 数组** `List<BusProduct>`；jwtInfo 取出不用。
- MP updateBatchById：逐条 `UPDATE bus_product SET <非null字段>…, update_time=now() WHERE id=?`。无任何校验/日志/MQ。
- 响应：success() data=null。body 空数组→同样 success。

### 3.19 POST `/busProduct/updByCompanyId`【写操作：无（恒定怪癖）】
- body：BusProduct（wareIdentity/isSignEtc/companyId 会被取出传参）；jwtInfo 不用。
- mapper 方法标了 `@Update` 但 SQL 是 **`select * from bus_product where available_tonnage is not null ORDER BY update_time desc`**——参数完全未引用，执行为只读 SELECT（写 0 行）。**Go 侧等价实现：什么都不做，直接 success()**（登记差异：Java 侧从不改数据）。
- 响应：success() data=null。

### 3.20 POST `/busProduct/selNewPro`【只读】
- 无 body 消费；jwtInfo 不用。
- newPro=`select (thName 子查 bus_big_prodect by three_pro_id), bus_product.* from bus_product where status=1 and (fprice_upd_time is null and pprice_upd_time is null) and TIMESTAMPDIFF(SECOND, last_check_time, update_time) < 10 ORDER BY update_time desc LIMIT 1`（last_check_time 为 NULL → 条件为 NULL → 不命中）。
- pricePro=同结构，`(fprice_upd_time is not null or pprice_upd_time is not null)`。
- 响应：data=`{"newPro":{...}|null, "pricePro":{...}|null}`。

### 3.21 POST `/busProduct/selCnPro`【只读】
- body：BusProduct（**仅 page:{current,size} 被消费，其余字段全忽略**）；jwtInfo 不用。
- `SELECT * FROM bus_product WHERE available_tonnage IS NOT NULL AND status=1 ORDER BY update_time DESC` 分页（Page 参数）。
- 响应：data=IPage 对象（0.4 形状），records 内实体无 tags。

### 3.22 POST `/busProduct/selNowNotice`【只读；外部依赖：redis】
- body/jwtInfo 均不用。
- count=`SELECT COUNT(*) FROM bus_product WHERE status=1`；broadcastMsg=`"今日"+count+"个《炼焦煤》直供平台"`；redis `GET broadcastMsg` 存在则**整体替换**（**environmentType==3 时写侧用的是 zscbroadcastMsg，本接口只读 broadcastMsg——照抄该不一致**）。
- 响应：data=`broadcastMsg.split(";")` → **String[] JSON 数组**（尾随分号产生空串元素，照抄）。

### 3.23 POST `/busProduct/pushOfficialMsg`【外部依赖：feign、微信 HTTP；不写业务表】
- body 不消费；jwtInfo 不用。
- `pushTemplateMsg(new JSONObject())`：先 feign `/chainAuthChat/getAccessToken` {appIdentification:"cloudAccount"}（**在 try 外，feign 失败→999**）→ POST `https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=...`，body=`{}`（微信必返 errcode，结果仅日志）。发送异常被吞。
- bus_mq_consume_log 不写。响应：恒 success()。

### 3.24 POST `/busProduct/getCareOfficial`（judgeIsFollow）【外部依赖：feign、微信 HTTP；只读】
- body：`{"openid":"..."}`。**openid 缺失 → `maps.get("openid").toString()` NPE→999**；jwtInfo 不用。
- token=feign getAccessToken(cloudAccount)（try 外，失败→999）；GET `https://api.weixin.qq.com/cgi-bin/user/info?access_token=...&openid=...&lang=zh_CN`；解析 `subscribe` 字段（微信异常/网络失败 → catch → `ResultTool.fail("判断是否关注公众号异常：")` → `{success:false,code:999,msg:"失败",data:"判断是否关注公众号异常："}`）。
- subscribe==1 → `success("已关注")`；否则 → `success("未关注")`（data=字符串）。

### 3.25 POST `/busProduct/syncEtcScore`（updEtcScore(null,null)）【写操作：chain_basic_enterprise 三列】
- body 不消费；jwtInfo 不用。proId/companyId 传 null → proWhere=""、etcWhere=""（全量模式）。
- 顺序：
  1. clearScore("")：`update chain_basic_enterprise set upd_share_score=0, report_score=0, upd_price_score=0 where exists (select 1 from bus_product where company_id=chain_basic_enterprise.id)`。
  2. selUpdSharePro("")：`SELECT distinct company_id from bus_product where status=1 and ((fprice_upd_time is not null and TO_DAYS(NOW())-TO_DAYS(fprice_upd_time)<=30) or (pprice_upd_time 同)) and exists (select 1 from bus_pro_log where pro_id=bus_product.id and operate_obj=4 and (operate_type=4 or 5) and user_id=bus_product.create_user_id)`；非空 → getWhere 拼 `where (id = a or id = b ...)` → updPriceScore：`update chain_basic_enterprise set upd_share_score=0.5 <where>`。
  3. selHydPro("")：`... status=1 and exists (select 1 from chain_common_file where association_id=bus_product.association_id and f_type=2 and effective=1)` → updHydScore：`set report_score=0.25`。
  4. 价格分 selPricePro/updPceScore **已注释禁用**（Go 不要实现）。
- 响应：success()。

### 3.26 POST `/busProduct/selBasCount`【只读】
- body/jwtInfo 不用。
- proCount=`SELECT COUNT(*) FROM bus_product WHERE effective=1`；userCount=`SELECT COUNT(*) FROM chain_auth_user WHERE phone IS NOT NULL AND channel=2`；visitCount=`SELECT COUNT(*) FROM bus_pro_log WHERE operate_type=7`。
- 响应：data=`{"proCount":n,"userCount":n,"visitCount":n}`。

---

## 4. NPE / 恒定行为 / 越界 快查表（→999 汇总）
| 端点 | 触发条件 | 结果 |
|---|---|---|
| listByPage | searchStr=null 且 5 个 isDirect count 全 0 | NPE→999 |
| listByPage | provinceCode/cityCode/areaCode 非数字裸拼 | SQL 语法错→999 |
| listByPage | warehouseType=1 且无 jwtInfo | NPE→999 |
| listByPage | 分词分支"硫/挥发/水"数字词为最后一个词 | IndexOOB→999 |
| listByPage | bus_score_rule 空表（经积分链路） | IndexOOB→999 |
| getById | id 不存在 | fail(602,"产品已下架") |
| getById | userMap（feign data）为 null / msg_type=3 配置空 / 客服用户查无 / brokerStartDate 缺失 / businessId=null 但 areaBuserId 非 null / 三级分类 id 无效 / createNewCode 的 HY/KJ 或空结果 | NPE 或 IndexOOB 或 NumberFormat→999 |
| add | 企业 id 无效 / bigProId/threeProId 无效 | NPE→999 |
| updateById | id 不存在 / bp.createPhone=null 且传新 createPhone / bp.status=null | NPE→999 |
| deleteById/delBatch | id 不存在 | bus_pro_log NPE→999（删除已发生，无回滚） |
| hisPrice | hisPriceType=null | fail(802,"参数异常") |
| hisPrice | priceType ∉{1,2,3} | NPE→999 |
| addViewTime/addPhoneTime/shareTime 前者 | id 不存在 | NPE→999（shareTime 例外：静默 success） |
| addPhoneTime | body 无 phoneLog | insert(null)→999 |
| getDetailImg | 任何异常（含价差字段 null、背景图缺失） | 吞掉→success(null) |
| getCareOfficial | openid 缺失 | NPE→999 |
| pushOfficialMsg | feign getAccessToken 失败 | 999 |
| updByCompanyId | 任意 | 恒 success，无写入 |
| selNowNotice | 恒定默认文案或 redis broadcastMsg | 只读 |
| pushOfficialMsg | 恒 success（微信侧失败仅日志） |

## 5. 外部依赖登记表（Go 单体改造点）
| 类别 | 键/地址 | 用途 | 端点 |
|---|---|---|---|
| Redis | `userId{uid}` hash.phone | 登录用户手机号（etc_rel_phone 条件） | listByPage |
| Redis | `broadcastMsg`（env3 写 `zscbroadcastMsg`） | 公告广播追加/读取 | selNowNotice（读）；写入方在 MQ 消费侧 |
| Redis | `{fileName}`=base64 | 详情图缓存 | getDetailImg |
| Redis | `smsPhone{uid}`、`{proId}t8`、`updProNotice*`、`addProNotice*`、`{clogId}`、`{uid}t1` | 短信/红包去重锁 | MQ 消费侧/红包 |
| RocketMQ | `config.syncTopicName` tag syncOtherInfo | 更新后同步（价格历史/审核消息/积分/红包/广播/推送/企业分） | updateById（发出）；消费侧 2.6/2.7 |
| RocketMQ | `config.noticeTopicName` tag updPro/addPro | 公众号通知 | 消费侧 |
| 微信 HTTP | api.weixin.qq.com/cgi-bin/message/template/send | 公众号模板推送 | pushOfficialMsg（空体）、消费侧 |
| 微信 HTTP | api.weixin.qq.com/cgi-bin/user/info | 是否关注公众号 | getCareOfficial |
| unipush HTTP | config.unipushAdr | App 推送 | 消费侧 |
| 短信 | 模板 2068872/2068879/1767390/2119567/2165435/2165442，签名单"山西海云启辰科技" | 平台/用户通知 | add、updateById、消费侧 |
| Feign saas-auth | /chainAuthUser/getById、/selByIds、/list、/chainAuthChat/getAccessToken | 用户信息/客服/token | getById、addPhoneTime、pushOfficialMsg、getCareOfficial、红包链路 |
| Feign saas-common | /chainCommonFile/updateById、/updBthById、/list、/add | 化验单/图片/视频文件 | add、updateById、getById、getDetailImg |
| HanLP | segment + 词典"主焦","1/3焦" | 全字段搜索分词 | listByPage（Go：登记差异，可用简化分词器，决策树照抄） |
| 本地文件 | /home/ubuntu/uploadfiles/frame73.jpg、落盘 *.jpg、background.jpg（旧版 getDetailImg 未暴露） | 图片合成 | getDetailImg |

## 6. 表写入汇总（全部【写操作】端点）
| 表 | 端点 | 列 |
|---|---|---|
| bus_product | add INSERT；updateById/updateBatch/shareTime/addViewTime/addPhoneTime UPDATE（非null字段）；getById UPDATE(association_id)；消费侧 UPDATE(价格时间/pre价/认证态/分类名/last_check_time/down_time)；deleteById/delBatch DELETE | 全列 |
| bus_del_product | deleteById/delBatch INSERT…SELECT 备份 | 全列 |
| bus_pro_log | deleteById/delBatch(type3)、shareTime(body直插)、getById 访客(type7)、消费侧 | 依构造/直插 |
| bus_pro_tag | add saveBatch；updateById DELETE+saveBatch | 全列 |
| bus_common_log | add(type1)/updateById(type2) INSERT | user_id,username,phone,pro_id,operation,operation_type |
| bus_msg / bus_user_msg | getById 访客(type4)、消费侧审核(type3)/价格(type1) INSERT | msg_type,content,remark,jump_id / user_id,msg_id |
| bus_visitor_record | getById 访客 INSERT/UPDATE | pro_id,visitor_id,username,nickname,phone,wx_head_img,create/update_time,visitor_resource |
| bus_user_click_pro | getById 访客 INSERT/UPDATE | user_id,last_pro_id,last_type_id,three_pro_id,click_time,is_last_click |
| bus_phone_log | addPhoneTime INSERT（body 直插） | 全列 |
| bus_price_log | 消费侧改价 INSERT | user_id,fst/big/three_pro_id,pro_id,pre_price,price,price_type |
| bus_videos / chain_common_file(feign) | add 视频链路 / add、updateById 文件状态 / getDetailImg add(fType5) | 见 3.4/3.6/3.17 |
| bus_user_score / bus_score_log / bus_score_rule | updUserScore（add/updateById/消费侧） | 见 2.1 |
| bus_red_money_clog / bus_red_envelope_record | 红包链路 INSERT/核销 | 见 2.1/2.6 |
| chain_basic_enterprise | updCheckScore、clearScore、updPriceScore、updHydScore（updEtcScore/消费侧） | coal_verify_score、upd_share_score、report_score、upd_price_score |
| chain_auth_user | updateById 改 createPhone 时 INSERT 新用户 | phone,account,channel |
| bus_mq_consume_log | 每次 MQ 发送 INSERT | ip,topic,tag,create_time,message_id,body |

> 逐字节一致性要点：IPage 形状（0.4）、LocalDateTime 无 @JsonFormat 的 ISO `T` 格式、TreeMap 键序（hisPrice.priceData）、String[]（selNowNotice）、fail 文案原样、`success`/`code`/`msg` 三键恒在、`data:null` 显式输出、GROUP BY 任取一行、`LIKE '%%v%%'` 等价语义、分词分支全未命中时返回全表。
