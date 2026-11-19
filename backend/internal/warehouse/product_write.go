// 煤种写接口及积分、日志、持久化后续同步，按服务器运行包保留原行为。
package warehouse

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func init() {
	addExtra(func(mux *http.ServeMux, s *Service) {
		reg := func(path string, h http.HandlerFunc) {
			mux.HandleFunc("POST /saas-warehouse"+path, authGuard(h))
		}
		reg("/busProduct/add", s.bpAdd)
		reg("/busProduct/updateById", s.bpUpdateById)
		reg("/busProduct/deleteById", s.bpDeleteById)
		reg("/busProduct/delBatch", s.bpDelBatch)
		reg("/busProduct/importExcel", s.bpImportExcel)
		reg("/busProduct/shareTime", s.bpShareTime)
		reg("/busProduct/addViewTime", s.bpAddViewTime)
		reg("/busProduct/addPhoneTime", s.bpAddPhoneTime)
		reg("/busProduct/updateBatch", func(w http.ResponseWriter, r *http.Request) { s.tplUpdBatch(w, r, entBusProduct) })
		reg("/busProduct/updByCompanyId", func(w http.ResponseWriter, r *http.Request) {
			platform.Success(w, nil) // Java mapper @Update 里是 SELECT,恒无写入
		})
		reg("/busProduct/syncEtcScore", s.bpSyncEtcScore)
		reg("/busProduct/getDetailImg", s.bpGetDetailImgActual)
		reg("/busProduct/pushOfficialMsg", s.bpPushOfficialActual)
		reg("/busProduct/getCareOfficial", s.bpGetCareActual)
	})
}

// uuid32 32 位 hex(对照 Java EntityUtils.getUuid)
func uuid32() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// bpSaveBusProduct = saveBusProduct(规格 3.4):校验 602 → selEtcById → INSERT →
// 视频/标签/日志 → updUserScore(2,2)
func (s *Service) bpAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusProduct.validateBody(m) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	publishType := m["publishType"]
	if publishType == nil || str(m, "publishType") != "2" {
		for _, k := range []string{"wareIdentity", "coalType", "bigProId", "threeProId", "coalRelName", "coalRelPhone"} {
			if m[k] == nil {
				platform.Fail(w, 602, "煤种身份、煤种区别、煤炭分类、煤炭品种、出厂现汇价或车板价、联系人姓名、联系人电话、煤棚照片为必填项")
				return
			}
		}
		if m["exFactoryPrice"] == nil && m["platePrice"] == nil {
			platform.Fail(w, 602, "煤种身份、煤种区别、煤炭分类、煤炭品种、出厂现汇价或车板价、联系人姓名、联系人电话、煤棚照片为必填项")
			return
		}
	}
	companyId := str(m, "companyId")
	if str(m, "wareIdentity") != "4" && companyId == "" {
		platform.Fail(w, 602, "请选择企业")
		return
	}
	if companyId != "" {
		etc, err := queryRows(d, "SELECT * from chain_basic_enterprise where id="+companyId)
		if err != nil || len(etc) == 0 {
			platform.Unknown(w) // selEtcById null → NPE
			return
		}
		m["enterpriseName"] = etc[0]["enterprise_name"]
		m["groupId"] = etc[0]["group_id"]
	}
	m["createUserId"] = json.Number(strconvI64(currentUser(r).UserID))
	m["associationId"] = uuid32()
	// reportList/imgList:feign 更新 chain_common_file(effective 等)
	for _, key := range []string{"reportList", "imgList"} {
		if arr, ok := m[key].([]any); ok {
			for _, it := range arr {
				f, _ := it.(map[string]any)
				if f == nil {
					continue
				}
				vals := map[string]any{"f_type": map[string]string{"reportList": "2", "imgList": "4"}[key],
					"owner_type": "1", "owner_id": str(m, "createUserId"), "relation_type": "1", "effective": "1",
					"association_id": str(m, "associationId")}
				_, _ = d.Exec("UPDATE chain_common_file SET f_type=?, owner_type=?, owner_id=?, relation_type=?, effective=?, association_id=? WHERE id=?",
					vals["f_type"], vals["owner_type"], vals["owner_id"], vals["relation_type"], vals["effective"], vals["association_id"], str(f, "id"))
			}
		}
	}
	m["isHasReport"] = true
	if v := str(m, "exFactoryPrice"); v != "" {
		m["preExFactoryPrice"] = json.Number(v)
	}
	if v := str(m, "platePrice"); v != "" {
		m["prePlatePrice"] = json.Number(v)
	}
	// secName/thName=bus_big_prodect(无效 id → NPE 999)
	if v := str(m, "bigProId"); v != "" {
		nm := s.bpNameOf(d, v)
		if nm == nil {
			platform.Unknown(w)
			return
		}
		m["secName"] = nm
	}
	if v := str(m, "threeProId"); v != "" {
		nm := s.bpNameOf(d, v)
		if nm == nil {
			platform.Unknown(w)
			return
		}
		m["thName"] = nm
	}
	if v := str(m, "fstProId"); v != "" {
		if nm := s.bpNameOf(d, v); nm != nil {
			m["fstName"] = nm
		}
	}
	vals := pick(m, entBusProduct.dbcols)
	delete(vals, "id")
	nid, err := platform.Insert(d, "bus_product", vals)
	if err != nil {
		platform.Unknown(w)
		return
	}
	proID := strconv.FormatInt(nid, 10)
	// videoList→saveBusVideos:videoId null→fail 602 但被忽略继续;INSERT bus_videos
	if arr, ok := m["videoList"].([]any); ok {
		for _, it := range arr {
			f, _ := it.(map[string]any)
			if f == nil {
				continue
			}
			if f["videoId"] == nil {
				continue // fail(602,"视频上传失败，请重试") 被吞,继续
			}
			_, _ = d.Exec("INSERT INTO bus_videos (pro_id, video_id, video_name, video_title, status) VALUES (?,?,?,?,?)",
				proID, str(f, "videoId"), str(f, "videoName"), str(f, "videoTitle"), "0")
			_, _ = d.Exec("UPDATE chain_common_file SET effective=1 WHERE id=?", str(f, "videoId"))
		}
	}
	// proTags→saveBatch bus_pro_tag
	if arr, ok := m["proTags"].([]any); ok {
		for _, it := range arr {
			if tagId := str(map[string]any{"v": it}, "v"); tagId != "" {
				_, _ = d.Exec("INSERT INTO bus_pro_tag (pro_id, sign_id) VALUES (?,?)", proID, tagId)
			}
		}
	}
	score := 0
	if publishType == nil || str(m, "publishType") == "1" {
		score = s.bpUpdUserScore(w, d, m, 2, 2)
		if score < 0 {
			return
		}
	}
	s.bpLog(w, r, d, "新增了煤种，新增信息为：\n"+s.custLog(m, true), "1")
	if score < 0 {
		return
	}
	platform.Success(w, jsonNumberInt64(int64(score)))
}

// bpNameOf bus_big_prodect 名查,无行返回 nil(模拟 Java NPE 面)
func (s *Service) bpNameOf(d *sql.DB, id string) *string {
	var nm sql.NullString
	if err := d.QueryRow("SELECT product_name FROM bus_big_prodect WHERE id=?", id).Scan(&nm); err != nil {
		return nil
	}
	return &nm.String
}

// bpLog 写 bus_common_log(无 fill 实体:create_time 走库默认)
func (s *Service) bpLog(w http.ResponseWriter, r *http.Request, d *sql.DB, operation, opType string) {
	u := currentUser(r)
	_, err := d.Exec("INSERT INTO chain_common_log (user_id, user_name, operation, operation_type) VALUES (?,?,?,?)",
		strconvI64(u.UserID), u.Username, operation, opType)
	if err != nil {
		platform.Unknown(w)
	}
}

// bpUpdUserScore 对照 updUserScore(规格 2 节);返回 score,负数表示已写 999 响应
func (s *Service) bpUpdUserScore(w http.ResponseWriter, d *sql.DB, bp map[string]any, getType, scoreType int) int {
	userId := str(bp, "createUserId")
	userPhone := str(bp, "createPhone")
	var ruleMoGetSumScore, publishCoalScore, updPriceScore int
	var ruleID string
	rows, err := queryCam(d, "SELECT * FROM bus_score_rule")
	if err != nil || len(rows) == 0 {
		platform.Unknown(w) // list().get(0) 越界
		return -1
	}
	rule := rows[0]
	ruleID = str(rule, "id")
	ruleMoGetSumScore = int(num(rule["moGetSumScore"]))
	publishCoalScore = int(num(rule["publishCoalScore"]))
	updPriceScore = int(num(rule["updPriceScore"]))
	// userScore 行
	var urow map[string]any
	urows, err := queryCam(d, "SELECT * FROM bus_user_score WHERE user_id="+userId)
	if err == nil && len(urows) > 0 {
		urow = urows[0]
	}
	isFst := urow == nil
	hjjNum, djHjjNum := 0, 0
	if urow != nil {
		hjjNum, djHjjNum = int(num(urow["hjjNum"])), int(num(urow["djHjjNum"]))
	}
	// status==2 → 返回当月冻结积分(仅 updateById 审核 paths 调用,status 均来自 bp)
	if st := str(bp, "status"); st == "2" {
		var fs int
		_ = d.QueryRow("select ifnull(score,0) from bus_score_log where user_id = " + userId + " and pro_id = " + str(bp, "id") + " and score_type=2 and score_get_type=3 limit 1").Scan(&fs)
		newMo := ruleMoGetSumScore + fs
		_, _ = d.Exec("UPDATE bus_score_rule SET mo_get_sum_score=? WHERE id=?", newMo, ruleID)
		if urow != nil {
			_, _ = d.Exec("UPDATE bus_user_score SET dj_hjj_num=? WHERE id=?", djHjjNum-fs, str(urow, "id"))
			s.bpScoreLog(d, userId, userPhone, "审核未通过，返回冻结积分", 0, hjjNum, hjjNum, djHjjNum, djHjjNum-fs, str(bp, "id"), 4, 1, getType)
		}
		return 0
	}
	updScore := 0
	content := ""
	ruleType := 1
	orgA, orgF := hjjNum, djHjjNum
	finA, finF := 0, 0
	switch getType {
	case 1:
		content = "首次登录"
		updScore = minInt(ruleMoGetSumScore, int(num(rule["registScore"])))
		hjjNum += updScore
		finA, finF = orgA+updScore, orgF
	case 2:
		content = "发布煤种"
		updScore = minInt(ruleMoGetSumScore, publishCoalScore)
		djHjjNum += updScore
		finA, finF = orgA, orgF+updScore
	case 3:
		content = "更新煤种价格"
		updScore = minInt(ruleMoGetSumScore, updPriceScore)
		djHjjNum += updScore
		finA, finF = orgA, orgF+updScore
	case 5:
		content = "审核通过发布煤种"
		var fs int
		_ = d.QueryRow("select ifnull(score,0) from bus_score_log where user_id = " + userId + " and pro_id = " + str(bp, "id") + " and score_type=2 and score_get_type=2 limit 1").Scan(&fs)
		if fs == 0 {
			fs = publishCoalScore
			s.bpScoreLog(d, userId, userPhone, "发布煤种", fs, orgA, orgA, orgF, orgF+fs, str(bp, "id"), 1, 2, 2)
		}
		hjjNum += fs
		if djHjjNum > 0 {
			djHjjNum -= fs
		} else {
			newMo := ruleMoGetSumScore - publishCoalScore
			_, _ = d.Exec("UPDATE bus_score_rule SET mo_get_sum_score=? WHERE id=?", newMo, ruleID)
			ruleMoGetSumScore = newMo
		}
		ruleType = 3
		finA = orgA + fs
		finF = orgF - fs
		if finF < 0 {
			finF = 0
		}
		updScore = fs
	case 6:
		content = "审核通过更新价格"
		var fs int
		_ = d.QueryRow("select ifnull(score,0) from bus_score_log where user_id = " + userId + " and pro_id = " + str(bp, "id") + " and score_type=2 and score_get_type=3 limit 1").Scan(&fs)
		if fs == 0 {
			fs = updPriceScore
			s.bpScoreLog(d, userId, userPhone, "更新煤种价格", fs, orgA, orgA, orgF, orgF+fs, str(bp, "id"), 1, 2, 3)
		}
		hjjNum += fs
		if djHjjNum > 0 {
			djHjjNum -= fs
		} else {
			newMo := ruleMoGetSumScore - updPriceScore
			_, _ = d.Exec("UPDATE bus_score_rule SET mo_get_sum_score=? WHERE id=?", newMo, ruleID)
			ruleMoGetSumScore = newMo
		}
		ruleType = 3
		finA = orgA + fs
		finF = orgF - fs
		if finF < 0 {
			finF = 0
		}
		updScore = fs
	}
	// 当月上限:moGetSumScore<=0 → return 0(不落库)
	if ruleMoGetSumScore <= 0 {
		return 0
	}
	if ruleMoGetSumScore >= updScore {
		newMo := ruleMoGetSumScore - updScore
		_, _ = d.Exec("UPDATE bus_score_rule SET mo_get_sum_score=? WHERE id=?", newMo, ruleID)
	}
	if isFst {
		if _, err := d.Exec("INSERT INTO bus_user_score (user_id, hjj_num, dj_hjj_num) VALUES (?,?,?)", userId, hjjNum, djHjjNum); err != nil {
			platform.Unknown(w)
			return -1
		}
	} else {
		if _, err := d.Exec("UPDATE bus_user_score SET hjj_num=?, dj_hjj_num=? WHERE id=?", hjjNum, djHjjNum, str(urow, "id")); err != nil {
			platform.Unknown(w)
			return -1
		}
	}
	s.bpScoreLog(d, userId, userPhone, content, updScore, orgA, finA, orgF, finF, str(bp, "id"), ruleType, scoreType, getType)
	return updScore
}

// bpScoreLog INSERT bus_score_log(fill 时间戳)
func (s *Service) bpScoreLog(d *sql.DB, userId, userPhone, content string, score, orgA, finA, orgF, finF int, proID string, ruleType, scoreType, getType int) {
	_, _ = d.Exec("INSERT INTO bus_score_log (user_id, user_phone, content, score, orgin_available_score, final_available_score, orgin_freeze_score, final_freeze_score, pro_id, rule_type, score_type, score_get_type) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)",
		userId, userPhone, content, strconv.Itoa(score), strconv.Itoa(orgA), strconv.Itoa(finA), strconv.Itoa(orgF), strconv.Itoa(finF), proID, strconv.Itoa(ruleType), strconv.Itoa(scoreType), strconv.Itoa(getType))
}

// custLog 新增/修改日志内容(对照 getCustAddLog/getCustChangeLog 关键字段)
func (s *Service) custLog(m map[string]any, add bool) string {
	parts := []string{}
	for k, label := range map[string]string{"productName": "煤种名称", "enterpriseName": "企业名称", "warehouseNb": "云仓号",
		"coalRelName": "联系人", "coalRelPhone": "联系电话", "exFactoryPrice": "出厂现汇价", "platePrice": "车板价"} {
		if v := str(m, k); v != "" {
			parts = append(parts, label+"："+v)
		}
	}
	return strings.Join(parts, "；")
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// bpUpdateById = updateBusProduct(规格 3.6)
func (s *Service) bpUpdateById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusProduct.validateBody(m) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	id := str(m, "id")
	if id == "" {
		platform.Unknown(w) // getById(null) → NPE
		return
	}
	bpRows, _, err := s.queryEnt(r, "SELECT * FROM bus_product WHERE id=?", id)
	if err != nil || len(bpRows) == 0 || bpRows[0] == nil {
		platform.Unknown(w)
		return
	}
	bp := bpRows[0]
	associationId := str(bp, "associationId")
	if associationId == "" {
		associationId = uuid32()
		_, _ = d.Exec("UPDATE bus_product SET association_id=? WHERE id=?", associationId, id)
	}
	// reportList:非空批量 effective 生效;空 → isHasReport=false
	if arr, ok := m["reportList"].([]any); ok && len(arr) > 0 {
		for _, it := range arr {
			f, _ := it.(map[string]any)
			if f != nil {
				_, _ = d.Exec("UPDATE chain_common_file SET effective=1, association_id=? WHERE id=?", associationId, str(f, "id"))
			}
		}
		m["isHasReport"] = true
	} else {
		m["isHasReport"] = false
	}
	m["updateTime"] = nowSpace()
	// createPhone 变更:createUserId 换人或新建 chain_auth_user(channel 2)
	if cp := str(m, "createPhone"); cp != "" && cp != str(bp, "createPhone") {
		if str(bp, "createPhone") == "" {
			platform.Unknown(w) // bp.createPhone null → NPE
			return
		}
		var uid sql.NullString
		_ = d.QueryRow("SELECT id FROM chain_auth_user WHERE phone=? AND channel=2", cp).Scan(&uid)
		if !uid.Valid {
			res, err2 := d.Exec("INSERT INTO chain_auth_user (phone, channel, username) VALUES (?,?,?)", cp, "2", cp)
			if err2 == nil {
				if nid, e3 := res.LastInsertId(); e3 == nil {
					m["createUserId"] = json.Number(strconv.FormatInt(nid, 10))
				}
			}
		} else {
			m["createUserId"] = json.Number(uid.String)
		}
	}
	// status 迁移时间戳
	if st, has := m["status"]; has && st != nil {
		if str(m, "status") == "1" && str(bp, "status") == "0" {
			m["lastCheckTime"] = nowSpace()
		}
		if str(m, "status") == "3" && str(bp, "status") == "0" {
			m["downTime"] = nowSpace()
		}
	}
	vals := pick(m, entBusProduct.dbcols)
	delete(vals, "id")
	if _, err := platform.UpdateByID(d, "bus_product", id, vals); err != nil {
		platform.Unknown(w)
		return
	}
	if str(m, "isForceUpdPrice") == "true" {
		if m["exFactoryPrice"] == nil {
			_, _ = d.Exec("update bus_product set ex_factory_price=null where id=" + id)
		}
		if m["platePrice"] == nil {
			_, _ = d.Exec("update bus_product set plate_price=null where id=" + id)
		}
	}
	if str(m, "effective") == "0" {
		platform.Success(w, jsonNumberInt64(0))
		return
	}
	// imgList feign 生效
	if arr, ok := m["imgList"].([]any); ok {
		for _, it := range arr {
			f, _ := it.(map[string]any)
			if f != nil {
				_, _ = d.Exec("UPDATE chain_common_file SET effective=1, association_id=? WHERE id=?", associationId, str(f, "id"))
			}
		}
	}
	// proTags:先删后插
	if _, ok := m["proTags"]; ok {
		_, _ = d.Exec("DELETE FROM bus_pro_tag WHERE pro_id=?", id)
		if arr, ok2 := m["proTags"].([]any); ok2 {
			for _, it := range arr {
				if tagId := str(map[string]any{"v": it}, "v"); tagId != "" {
					_, _ = d.Exec("INSERT INTO bus_pro_tag (pro_id, sign_id) VALUES (?,?)", id, tagId)
				}
			}
		}
	}
	s.bpLog(w, r, d, "修改了煤种："+str(bp, "thName")+s.custLog(m, false), "2")
	score := s.bpOriginalUpdPriceScore(w, d, bp, m)
	if score < 0 {
		return
	}
	job := r.Clone(context.WithoutCancel(r.Context()))
	job.Header = r.Header.Clone()
	job.Header.Del("applytype") // 原 MQ 消费者没有请求头，使用默认业务库。
	first := str(m, "status") == "1" && str(bp, "status") == "0" && bp["lastCheckTime"] == nil
	userID, oldStatus := strconv.FormatInt(currentUser(r).UserID, 10), str(bp, "status")
	s.submitProductJob(job, productJob{ProID: id, UserID: userID, First: first, OriginalStatus: oldStatus})
	platform.Success(w, jsonNumberInt64(int64(score)))
}

func (s *Service) bpLastPrice(d *sql.DB, proId, priceType string) string {
	var p sql.NullString
	_ = d.QueryRow("SELECT price FROM bus_price_log WHERE pro_id=" + proId + " AND price_type=" + priceType + " ORDER BY create_time DESC LIMIT 1").Scan(&p)
	return p.String
}

// bpDeleteById = delete:备份 bus_del_product → DELETE → saveDelLog
func (s *Service) bpDeleteById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusProduct.validateBody(m) {
		platform.Unknown(w)
		return
	}
	id := str(m, "id")
	d := s.pick(r)
	if _, err := d.Exec("INSERT INTO bus_del_product SELECT * FROM bus_product WHERE id=" + id); err != nil {
		platform.Unknown(w)
		return
	}
	bp, err := platform.GetByID(d, "bus_product", id)
	if err != nil || bp == nil {
		platform.Unknown(w) // getById null → saveDelLog NPE
		return
	}
	if _, err := d.Exec("DELETE FROM bus_product WHERE id=" + id); err != nil {
		platform.Unknown(w)
		return
	}
	s.bpSaveDelLog(w, r, d, bp)
	platform.Success(w, "删除成功，id为："+id)
}

// bpDelBatch = delBatch:逐 id 备份+日志,末尾批量 DELETE
func (s *Service) bpDelBatch(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	ids, _ := m["ids"].([]any)
	d := s.pick(r)
	for _, iv := range ids {
		id := str(map[string]any{"v": iv}, "v")
		if _, err := d.Exec("INSERT INTO bus_del_product SELECT * FROM bus_product WHERE id=" + id); err != nil {
			platform.Unknown(w)
			return
		}
		bp, err := platform.GetByID(d, "bus_product", id)
		if err != nil || bp == nil {
			platform.Unknown(w)
			return
		}
		s.bpSaveDelLog(w, r, d, bp)
	}
	if len(ids) > 0 {
		ph := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
		args := []any{}
		args = append(args, ids...)
		if _, err := d.Exec("DELETE FROM bus_product WHERE id IN ("+ph+")", args...); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, "删除成功")
}

// bpSaveDelLog 删除日志(remark 指标串,content 固定文案)
func (s *Service) bpSaveDelLog(w http.ResponseWriter, r *http.Request, d *sql.DB, bp map[string]any) {
	var segs []string
	for k, label := range map[string]string{"ashContent": "灰分Ad%≤", "waterContent": "水分Mt%≤", "std": "全硫St.d%≤", "volatiles": "挥发分Vdaf%≤", "bondIndex": "粘结指数G≥"} {
		if v := decStr(bp[k]); v != "" && v != "0" {
			segs = append(segs, label+v)
		}
	}
	remark := strings.Join(segs, "，")
	content := "id为" + str(bp, "createUserId") + "的用户删除了产品名称为" + str(bp, "productName") + "的煤种，相关指标为：" + remark
	_, err := d.Exec("INSERT INTO bus_pro_log (user_id, pro_id, operate_type, operate_content) VALUES (?,?,3,?)",
		strconvI64(currentUser(r).UserID), str(bp, "id"), content)
	if err != nil {
		platform.Unknown(w)
	}
}

// bpImportExcel = importExcel+saveIrt(规格 3.7):errInfos 前缀"产品sheet中，";
// warehouseNb 不存在 errinfos 丢失 bug;status=1;saveBatch;响应 ImportResult
func (s *Service) bpImportExcel(w http.ResponseWriter, r *http.Request) {
	fh, ok := formFile(r)
	if !ok {
		platform.Unknown(w)
		return
	}
	irt := readExcelFile(fh, entBusProduct.excel, "", true, emptyDict())
	if len(irt.ErrInfos) == 0 {
		// warehouseNb 校验(结果从不写入返回值——原版 bug)
		for _, row := range irt.DataList {
			wn := row["warehouse_nb"]
			if wn == nil || wn == "" {
				continue
			}
			var cnt int
			_ = s.pick(r).QueryRow("SELECT count(1) FROM bus_warehouse WHERE warehouse_nb=?", wn).Scan(&cnt)
		}
	}
	irt.Errs = []any{}
	if len(irt.ErrInfos) > 0 {
		wrapped := make([]string, 0, len(irt.ErrInfos))
		for _, e := range irt.ErrInfos {
			wrapped = append(wrapped, "产品sheet中，"+e)
		}
		irt.ErrInfos = wrapped
		irt.Errs = []any{"导入失败 0"}
		irt.SucInfos = []any{}
		for i, row := range irt.DataList {
			irt.DataList[i] = platform.EntityEcho(platform.Camelize(row), entBusProduct.keys, nil)
		}
		platform.Success(w, irt)
		return
	}
	u := currentUser(r)
	for i, row := range irt.DataList {
		vals := insertVals(row)
		delete(vals, "id")
		vals["create_user_id"] = strconvI64(u.UserID)
		vals["create_phone"] = u.Account
		vals["status"] = "1"
		nid, err := platform.Insert(s.pick(r), "bus_product", vals)
		if err != nil {
			platform.Unknown(w)
			return
		}
		echo := platform.EntityEcho(platform.Camelize(row), entBusProduct.keys, map[string]any{"id": json.Number(strconv.FormatInt(nid, 10))})
		echo["createUserId"] = json.Number(strconvI64(u.UserID))
		echo["status"] = json.Number("1")
		irt.DataList[i] = echo
	}
	irt.SucInfos = []any{"云仓产品管理数据成功 " + strconv.Itoa(len(irt.DataList))}
	platform.Success(w, irt)
}

// bpShareTime = shareTime:产品 null 静默;sharingTimes+1;body.busProLog 原样入库
func (s *Service) bpShareTime(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusProduct.validateBody(m) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	bp, err := platform.GetByID(d, "bus_product", m["id"])
	if err != nil || bp == nil {
		platform.Success(w, nil) // 静默
		return
	}
	vt := int(num(bp["sharing_times"])) + 1
	_, err = d.Exec("UPDATE bus_product SET sharing_times=? WHERE id=?", strconv.Itoa(vt), str(bp, "id"))
	if err != nil {
		platform.Unknown(w)
		return
	}
	if logv, has := m["busProLog"]; has && logv != nil {
		lm, ok2 := logv.(map[string]any)
		if !ok2 {
			platform.Unknown(w)
			return
		}
		lvals := pick(lm, entBusProLog.dbcols)
		delete(lvals, "id")
		if _, err := platform.Insert(d, "bus_pro_log", lvals); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, nil)
}

// bpAddViewTime = addViewTime:产品 null → NPE 999;viewTimes+1
func (s *Service) bpAddViewTime(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusProduct.validateBody(m) {
		platform.Unknown(w)
		return
	}
	bp, err := platform.GetByID(s.pick(r), "bus_product", m["id"])
	if err != nil || bp == nil {
		platform.Unknown(w)
		return
	}
	vt := int(num(bp["view_times"])) + 1
	if _, err := s.pick(r).Exec("UPDATE bus_product SET view_times=? WHERE id=?", strconv.Itoa(vt), str(bp, "id")); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// bpAddPhoneTime = addPhoneTime:phoneTimes+1;body.phoneLog 原样入库(null → insert(null) 999)
func (s *Service) bpAddPhoneTime(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusProduct.validateBody(m) {
		platform.Unknown(w)
		return
	}
	bp, err := platform.GetByID(s.pick(r), "bus_product", m["id"])
	if err != nil || bp == nil {
		platform.Unknown(w)
		return
	}
	pt := int(num(bp["phone_times"])) + 1
	if _, err := s.pick(r).Exec("UPDATE bus_product SET phone_times=? WHERE id=?", strconv.Itoa(pt), str(bp, "id")); err != nil {
		platform.Unknown(w)
		return
	}
	if plv, has := m["phoneLog"]; has && plv != nil {
		lm, ok2 := plv.(map[string]any)
		if !ok2 {
			platform.Unknown(w)
			return
		}
		lvals := pick(lm, entBusPhoneLog.dbcols)
		delete(lvals, "id")
		if _, err := platform.Insert(s.pick(r), "bus_phone_log", lvals); err != nil {
			platform.Unknown(w)
			return
		}
	} else {
		platform.Unknown(w) // save(null) NPE
		return
	}
	platform.Success(w, nil)
}

// bpSyncEtcScore = syncEtcScore:清分→分享分 0.5→资料分 0.25(价格分已注释禁用)
func (s *Service) bpSyncEtcScore(w http.ResponseWriter, r *http.Request) {
	d := s.pick(r)
	if _, err := d.Exec("update chain_basic_enterprise set upd_share_score=0, report_score=0, upd_price_score=0 where exists (select 1 from bus_product where company_id=chain_basic_enterprise.id)"); err != nil {
		platform.Unknown(w)
		return
	}
	if _, err := d.Exec("update chain_basic_enterprise set upd_share_score=0.5 where id in (SELECT distinct company_id from bus_product where status = 1 and ((fprice_upd_time is not null and TO_DAYS(NOW()) - TO_DAYS(fprice_upd_time) <=30) or (pprice_upd_time is not null and TO_DAYS(NOW()) - TO_DAYS(pprice_upd_time) <=30)) and exists (select 1 from bus_pro_log where pro_id=bus_product.id and operate_obj=4 and (operate_type=4 or operate_type=5) and user_id=bus_product.create_user_id))"); err != nil {
		platform.Unknown(w)
		return
	}
	if _, err := d.Exec("update chain_basic_enterprise set report_score=0.25 where id in (SELECT distinct company_id from bus_product where status = 1 and exists (select 1 from chain_common_file where association_id=bus_product.association_id and f_type=2 and effective=1))"); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}
