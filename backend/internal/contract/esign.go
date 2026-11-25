package contract

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func (s *Service) esignGet(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	row, err := platform.GetByID(s.pick(r), "chain_eqb_flow", m["id"])
	if err != nil || row == nil {
		platform.Unknown(w)
		return
	}
	res := platform.EsignJSON(r.Context(), "GET", "/v1/signflows/"+platform.JnullAny(m["eqbFlowId"])+"/signers", nil)
	data, ok := res["data"].(map[string]any)
	if !ok {
		platform.Unknown(w)
		return
	}
	out := shapeEqbFlow([]map[string]any{row})[0]
	out["signers"] = data["signers"]
	platform.Success(w, out)
}

func (s *Service) esignExternal(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/confirmFlow") {
		s.esignFlow(w, r, false)
		return
	}
	if strings.HasSuffix(r.URL.Path, "/uploadFile") {
		s.esignUpload(w, r)
		return
	}
	if strings.HasSuffix(r.URL.Path, "/getSignCallBack") {
		s.esignCallback(w, r)
		return
	}
	d := s.pick(r)
	u := currentUser(r)
	if strings.HasSuffix(r.URL.Path, "/selLanucnInfo") {
		user, err := platform.GetByID(d, "chain_auth_user", u.UserID)
		if err != nil || user == nil {
			platform.Unknown(w)
			return
		}
		etc, err := platform.GetByID(d, "chain_basic_enterprise", u.CompanyID)
		if err != nil {
			platform.Unknown(w)
			return
		}
		name, credit := any(""), any("")
		if etc != nil {
			name, credit = etc["enterprise_name"], etc["rate_number"]
		}
		platform.Success(w, map[string]any{"companyName": name, "creditCode": credit, "username": user["username"], "phone": user["phone"]})
		return
	}
	m := bodyMap(r)
	if strings.HasSuffix(r.URL.Path, "/v3/selSignUrl") && strings.TrimSpace(str(m, "eqbFlowId")) == "" {
		platform.Fail(w, 602, "请等候其他签署人进行确认")
		return
	}
	flow := platform.JnullAny(m["eqbFlowId"])
	if !flowIDRe.MatchString(flow) { // 路径由请求值拼出:不校验就能让服务器替调用方签名请求 e签宝的任意路径
		platform.FailMsg(w, "签署流程ID无效")
		return
	}
	path := "/v1/signflows/" + flow
	key := "signers"
	switch {
	case strings.HasSuffix(r.URL.Path, "/selSignPeople"):
		path += "/signers"
	case strings.HasSuffix(r.URL.Path, "/selSignUrl"):
		user, err := platform.GetByID(d, "chain_auth_user", u.UserID)
		if err != nil || user == nil {
			platform.Unknown(w)
			return
		}
		if strings.Contains(r.URL.Path, "/v3/") {
			path = "/v3/signflows/" + flow + "/executeUrl?signerAccount=" + url.QueryEscape(platform.JnullAny(user["phone"]))
		} else {
			path += "/executeUrl?accountId=" + platform.JnullAny(user["eqb_per_id"])
		}
		key = "url"
	case strings.HasSuffix(r.URL.Path, "/v3DownFile"):
		path += "/documents"
		key = "docs"
	default:
		platform.Unknown(w)
		return
	}
	res := platform.EsignJSON(r.Context(), "GET", path, nil)
	data, ok := res["data"].(map[string]any)
	if !ok {
		platform.Unknown(w)
		return
	}
	value := data[key]
	if key == "docs" {
		docs, ok := value.([]any)
		if !ok || len(docs) == 0 {
			platform.Unknown(w)
			return
		}
		first, ok := docs[0].(map[string]any)
		if !ok {
			platform.Unknown(w)
			return
		}
		value = first["fileUrl"]
	}
	platform.Success(w, value)
}

func (s *Service) esignFlow(w http.ResponseWriter, r *http.Request, launch bool) {
	m := bodyMap(r)
	u := currentUser(r)
	code, text := 200, ""
	var result any
	err := platform.WithTx(s.pick(r), func(tx *sql.Tx) error {
		if launch {
			etc, err := platform.TxGetByID(tx, "chain_basic_enterprise", u.CompanyID)
			if err != nil {
				return err
			}
			if etc == nil || etc["sign_free_count"] == nil || etc["sign_charge_count"] == nil {
				return fmt.Errorf("Java 签章份数为 null")
			}
			if num(etc["sign_free_count"])+num(etc["sign_charge_count"]) <= 0 { // 份数被并发扣成负数后,==0 的判断会永久失效
				code, text = 602, "签章剩余份数不足，请联系平台管理员"
				return nil
			}
			if etc["pub_flow_id"] == nil {
				code, text = 602, "未找到相关配置方案，请联系平台管理员"
				return nil
			}
			user, err := platform.TxGetByID(tx, "chain_auth_user", u.UserID)
			if err != nil {
				return err
			}
			if user == nil || user["eqb_per_id"] == nil || etc["eqb_org_id"] == nil {
				return fmt.Errorf("Java eqbConfig 账号为 null")
			}
			m["launchPerId"], m["launchPubId"] = user["eqb_per_id"], etc["eqb_org_id"]
			id, err := platform.TxInsertNoFill(tx, "chain_eqb_flow", bodyVals(m, flowCols, nil))
			if err != nil {
				return err
			}
			m["id"] = id
			singers, ok := m["singers"].([]any)
			if !ok {
				return fmt.Errorf("Java singers 为 null")
			}
			for _, item := range singers {
				signer, ok := item.(map[string]any)
				if !ok {
					return fmt.Errorf("Java signer 为 null")
				}
				signer["sysFlowId"], signer["confirmLaunch"] = id, true
				if _, err = platform.TxInsertNoFill(tx, "chain_eqb_singers", bodyVals(signer, singersCols, nil)); err != nil {
					return err
				}
			}
			m["companyId"] = u.CompanyID
		}
		user, err := platform.TxGetByID(tx, "chain_auth_user", u.UserID)
		if err != nil {
			return err
		}
		if user == nil {
			return fmt.Errorf("Java user 为 null")
		}
		// updateBatchById 只更新非 null 字段；空 eqbPerId 不清除已有签署人账号。
		query := "UPDATE chain_eqb_singers SET confirm_launch=1"
		args := []any{}
		if user["eqb_per_id"] != nil {
			query += ",sign_per_id=?"
			args = append(args, user["eqb_per_id"])
		}
		query += " WHERE sys_flow_id <=> ? AND sign_phone <=> ?"
		args = append(args, m["id"], user["phone"])
		if _, err = tx.Exec(query, args...); err != nil {
			return err
		}
		var waiting int
		if err = tx.QueryRow("SELECT count(*) FROM chain_eqb_singers WHERE sys_flow_id <=> ? AND confirm_launch=0", m["id"]).Scan(&waiting); err != nil {
			return err
		}
		result = ""
		if waiting == 0 {
			rows, err := tx.Query("SELECT * FROM chain_eqb_singers WHERE sys_flow_id <=> ? AND confirm_launch=1", m["id"])
			if err != nil {
				return err
			}
			signers, err := platform.ScanMaps(rows)
			if err != nil {
				return err
			}
			etc, err := platform.TxGetByID(tx, "chain_basic_enterprise", m["companyId"])
			if err != nil {
				return err
			}
			if etc == nil {
				return fmt.Errorf("Java 发起企业为 null")
			}
			cfg, err := platform.TxGetByID(tx, "chain_eqb_config", etc["pub_flow_id"])
			if err != nil {
				return err
			}
			if cfg == nil {
				code, text = 602, "未找到相关配置方案，请联系平台管理员"
				return nil
			}
			v3 := launch || strings.Contains(r.URL.Path, "/v3/")
			payload, err := esignPayload(m, cfg, signers, v3)
			if err != nil {
				return err
			}
			path := "/api/v2/signflows/createFlowOneStep"
			if v3 {
				path = "/api/v3/signflows/createFlowOneStep"
			}
			res := platform.EsignJSON(r.Context(), "POST", path, payload)
			if str(res, "code") == "0" {
				data, ok := res["data"].(map[string]any)
				if !ok {
					return fmt.Errorf("Java 签章 data 为 null")
				}
				result = data["flowId"]
			}
			if _, err = tx.Exec("UPDATE chain_eqb_flow SET eqb_flow_id=?,flow_status=1 WHERE id=?", result, m["id"]); err != nil {
				return err
			}
		}
		etc, err := platform.TxGetByID(tx, "chain_basic_enterprise", u.CompanyID)
		if err != nil {
			return err
		}
		if etc == nil || etc["sign_free_count"] == nil {
			return fmt.Errorf("Java 签章企业或份数为 null")
		}
		if num(etc["sign_free_count"]) > 0 {
			_, err = tx.Exec("UPDATE chain_basic_enterprise SET sign_free_count=sign_free_count-1 WHERE id=? AND sign_free_count>0", u.CompanyID) // 条件自减:并发双击不会扣成负数
		}
		return err // 原 Java 即使尚有人未确认也扣免费份数，不额外修正。
	})
	if err != nil {
		log.Printf("签章流程失败: %v", err)
		platform.Unknown(w)
		return
	}
	if code != 200 {
		platform.Fail(w, code, text)
		return
	}
	platform.Success(w, result)
}

func esignConfig(cfg map[string]any, keys string) map[string]any {
	out := map[string]any{}
	for _, key := range strings.Fields(keys) {
		if value := cfg[strings.ToLower(key)]; value != nil {
			out[key] = value
		}
	}
	return out
}

func esignPayload(flow, cfg map[string]any, singers []map[string]any, v3 bool) (map[string]any, error) {
	// 旧表混用驼峰列名，Java 实体字段全小写；先按其实体映射归一化。
	normalized := map[string]any{}
	for key, value := range cfg {
		normalized[strings.ToLower(key)] = value
	}
	cfg = normalized
	for _, key := range []string{"autoarchive", "autoinitiate", "batchdropseal", "platformsign", "assignedposbean", "autoexecute"} {
		if cfg[key] != nil {
			cfg[key] = num(cfg[key]) != 0
		}
	}
	identify := map[string]any{"faceVideoTemplate": ""}
	for _, key := range []string{"orgAvailableAuthTypes", "personAvailableAuthTypes", "willTypes"} {
		value := cfg[strings.ToLower(key)]
		if value == nil {
			return nil, fmt.Errorf("Java %s.split null", key)
		}
		identify[key] = strings.Split(str(cfg, strings.ToLower(key)), ",")
	}
	if str(cfg, "personauthadvancedenabled") != "" {
		identify["personAuthAdvancedEnabled"] = strings.Split(str(cfg, "personauthadvancedenabled"), ",")
	}
	config := esignConfig(cfg, "batchDropSeal billIsolationCode countdown mobileShieldWay noticeDeveloperUrl noticeType redirectDelayTime redirectUrl signPlatform")
	for k, v := range identify {
		config[k] = v
	}
	if v3 {
		config = map[string]any{"identificationConfig": identify, "notifyConfig": esignConfig(cfg, "noticeDeveloperUrl noticeType"), "signConfig": esignConfig(cfg, "redirectUrl signPlatform batchDropSeal countdown redirectDelayTime mobileShieldWay")}
	}
	info := esignConfig(cfg, "autoArchive autoInitiate contractRemind contractValidity remark signValidity")
	info["businessScene"], info["flowConfigInfo"] = flow["topicName"], config
	if !v3 {
		info["initiatorAccountId"] = flow["launchPerId"]
		if flow["launchType"] == "2" {
			info["initiatorAuthorizedAccountId"] = flow["launchPubId"]
		}
	} else {
		delete(info, "remark")
	}
	name := any("")
	if v3 {
		name = flow["contractName"]
	}
	payload := map[string]any{"docs": []any{map[string]any{"encryption": "", "fileId": flow["field"], "fileName": name, "filePassword": ""}}, "flowInfo": info}
	if v3 {
		payload["attachments"] = []any{map[string]any{"fileId": flow["field"], "attachmentName": ""}}
	}
	list := []any{}
	for _, s := range singers {
		if s["sign_role"] == nil {
			return nil, fmt.Errorf("Java signRole 为 null")
		}
		account := map[string]any{"signerAccountId": s["sign_per_id"]}
		pos := esignConfig(cfg, "posPage posX posY")
		if !v3 {
			if num(s["sign_role"]) == 2 {
				account["authorizedAccountId"] = s["sign_etc_id"]
			}
			for k, v := range pos {
				account[k] = v
			}
			pos = map[string]any{}
		} else {
			account = map[string]any{"signerAccount": s["sign_phone"], "signerAccountInfo": map[string]any{"name": s["sign_name"], "bankCardNo": ""}, "authorizedAccount": s["sign_etc_name"], "orgAccountInfo": map[string]any{"legalRepName": "", "legalRepCertNo": "", "legalRepCertType": ""}, "noticeType": cfg["noticetype"]}
		}
		date := esignConfig(cfg, "fontSize format")
		for k, v := range map[string]string{"posPage": "qpospage", "posX": "qposx", "posY": "qposy"} {
			if cfg[v] != nil {
				date[k] = cfg[v]
			}
		}
		field := esignConfig(cfg, "assignedPosbean autoExecute sealType signDateBeanType signType width")
		field["actorIndentityType"], field["fileId"], field["posBean"], field["sealId"], field["signDateBean"] = s["sign_role"], flow["field"], pos, "", date
		if v3 {
			field["sealIds"] = []any{}
		}
		platformSign := cfg["platformsign"]
		if num(s["sign_role"]) == 0 {
			platformSign = false
		}
		signer := map[string]any{"platformSign": platformSign, "signerAccount": account, "signfields": []any{field}, "signOrder": "", "thirdOrderNo": ""}
		if v3 {
			signer["signOrder"], signer["forceReadTime"] = 1, cfg["forcereadtime"]
		}
		list = append(list, signer)
	}
	payload["signers"] = list
	return payload, nil
}

// e签宝流程/文件 ID 只允许字母数字与 -_:它会拼进出站请求路径,不能出现 / . ? 等字符
var flowIDRe = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
