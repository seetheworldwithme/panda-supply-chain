// common.go 二期组织基础收尾:saas-common 服务承载的企业档案/全局配置/e签宝配置/文件上传。
// Java 源:ChainBasicEnterpriseController/ChainGlobalConfigurationController/ChainEqbConfigController/
// UploadController(网关前缀 /saas-common,服务名 saas-common)。
// 已知 Java 缺陷的行为保持(勿"修复"):updateEnterprise 无事务,后段异常时前段写库已生效;
// filelist 分支先删后更(updateBatchById 落在已删行上,恒 0 行);tlzfConfig 为 null 时
// saveOrUpdate 直接返回 false 不报错。
package authadmin

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// ── 企业档案 chain_basic_enterprise ──

// entCols 企业实体列白名单(实体字段驼峰转下划线;载体字段除外)
var entCols = []string{"id", "enterprise_code", "admin_name", "enterprise_name", "simple_name",
	"account", "credit_code", "address", "legal_person", "controller_person", "contacts_person",
	"enterprise_type", "province_code", "city_code", "area_code", "county_code", "type",
	"reg_time", "rate_number", "sh_time", "phone", "leg_identity", "reg_money", "reg_adress",
	"businesser", "county_name", "boss", "status_mark", "design_pro", "rel_output", "pro_name",
	"city_name", "area_name", "country_name", "is_network", "rotation_flag", "pub_flow_id",
	"per_flow_id", "eqb_org_id", "merchant_number", "merchant_account", "sign_free_count",
	"sign_surplus_count", "sign_charge_count", "service_expiration_time", "advertising_video_id",
	"vr_link", "website_img_id", "website_ced_img_id", "is_chain_etc", "is_business_etc",
	"is_sign_etc", "is_member_etc", "ware_identity", "logo", "image_url", "rec_cap", "econ_kind",
	"updated_date", "scope", "is_authentication", "group_id", "tl_sign_num", "tl_acct_sub_no",
	"tl_mcht_no", "tl_open_time", "bus_license_id", "acc_opening_permit_id", "legal_ft_card_id",
	"legal_bk_card_id", "qua_cert_score", "bas_info_score", "coal_verify_score", "upd_share_score",
	"report_score", "qualified_score", "upd_price_score", "ycoal_silo", "jcoal_silo",
	"ycoal_silo_capacity", "jcoal_silo_capacity", "test_equipment_type", "is_displayed",
	"create_time", "update_time", "remark1", "remark2", "remark3", "remark4", "remark5"}

// entFill 企业实体全键(列键 + 载体键,对照 ChainBasicEnterprise Jackson 输出)
var entFill = []string{"id", "enterpriseCode", "adminName", "enterpriseName", "simpleName",
	"account", "creditCode", "address", "legalPerson", "controllerPerson", "contactsPerson",
	"enterpriseType", "provinceCode", "cityCode", "areaCode", "countyCode", "type",
	"regTime", "rateNumber", "shTime", "phone", "legIdentity", "regMoney", "regAdress",
	"businesser", "countyName", "boss", "statusMark", "designPro", "relOutput", "proName",
	"cityName", "areaName", "countryName", "isNetwork", "rotationFlag", "pubFlowId",
	"perFlowId", "eqbOrgId", "merchantNumber", "merchantAccount", "signFreeCount",
	"signSurplusCount", "signChargeCount", "serviceExpirationTime", "advertisingVideoId",
	"vrLink", "websiteImgId", "websiteCedImgId", "isChainEtc", "isBusinessEtc",
	"isSignEtc", "isMemberEtc", "wareIdentity", "logo", "imageUrl", "recCap", "econKind",
	"updatedDate", "scope", "isAuthentication", "groupId", "tlSignNum", "tlAcctSubNo",
	"tlMchtNo", "tlOpenTime", "busLicenseId", "accOpeningPermitId", "legalFtCardId",
	"legalBkCardId", "quaCertScore", "basInfoScore", "coalVerifyScore", "updShareScore",
	"reportScore", "qualifiedScore", "updPriceScore", "ycoalSilo", "jcoalSilo",
	"ycoalSiloCapacity", "jcoalSiloCapacity", "testEquipmentType", "isDisplayed",
	"createTime", "updateTime", "remark1", "remark2", "remark3", "remark4", "remark5",
	"searchStr", "page", "ids", "userId", "filelist", "imagelist", "password", "yyzz", "lbts",
	"pubFlows", "perFlows", "isAutoAdd", "startTime", "endTime", "isTlEtc", "tlzfConfig",
	"saleEtcs", "etcs", "isOtherEtc", "sumScore", "urcs", "selEtcType", "proCount"}

// shapeEnterprise 企业行输出:tinyint(1) 布尔化、DECIMAL→json.Number(保原样小数位)、
// 无 @JsonFormat 的 datetime(updateTime/tlOpenTime)→ ISO T 格式
func shapeEnterprise(rows []map[string]any) []map[string]any {
	out := platform.CamelizeRows(rows)
	for _, r := range out {
		if value, ok := r["regTime"].(string); ok && len(value) >= 10 {
			r["regTime"] = value[:10] // Java LocalDate, not LocalDateTime.
		}
		for _, c := range []string{"isNetwork", "rotationFlag", "isChainEtc", "isBusinessEtc",
			"isSignEtc", "isMemberEtc", "isAuthentication", "isDisplayed"} {
			if v, ok := r[c]; ok && v != nil {
				r[c] = num(v) != 0
			}
		}
		for _, c := range []string{"quaCertScore", "basInfoScore", "coalVerifyScore",
			"updShareScore", "reportScore", "qualifiedScore", "updPriceScore",
			"ycoalSiloCapacity", "jcoalSiloCapacity"} {
			if s, ok := r[c].(string); ok && s != "" {
				r[c] = json.Number(s)
			}
		}
		for _, c := range []string{"updateTime", "tlOpenTime"} {
			if s, ok := r[c].(string); ok && len(s) == 19 {
				r[c] = s[:10] + "T" + s[11:]
			}
		}
		fillNulls(r, entFill)
	}
	return out
}

// entGetById = ChainBasicEnterpriseServiceImpl.getEtc:企业行 + e签宝流程 + 通联配置 +
// 供应商 + 子公司 + 关联用户。企业行不存在时 Java 在 etc 上设值 NPE → 999
func (s *Service) entGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	id := str(m, "id")
	if id == "" {
		platform.Unknown(w)
		return
	}
	row, err := platform.GetByID(s.DB.Chain, "chain_basic_enterprise", id)
	if err != nil || row == nil {
		platform.Unknown(w)
		return
	}
	pubFlows, err := queryCam(s.DB.Chain, "SELECT * FROM chain_eqb_config WHERE flow_type=1")
	if err != nil {
		platform.Unknown(w)
		return
	}
	perFlows, err := queryCam(s.DB.Chain, "SELECT * FROM chain_eqb_config WHERE flow_type=2")
	if err != nil {
		platform.Unknown(w)
		return
	}
	tlzf, _ := queryCam(s.DB.Chain, "SELECT * FROM chain_tlzf_config WHERE company_id=? LIMIT 1", id)
	saleEtcs, err := queryCam(s.DB.Chain, "SELECT * FROM chain_tlzf_sale_etc WHERE buy_company_id=?", id)
	if err != nil {
		platform.Unknown(w)
		return
	}
	etcs, err := queryCam(s.DB.Chain, "SELECT * FROM chain_basic_enterprise WHERE group_id=?", id)
	if err != nil {
		platform.Unknown(w)
		return
	}
	urcs, err := queryRows(s.DB.Chain, "SELECT * FROM chain_auth_user WHERE company_id=? AND channel=2", id)
	if err != nil {
		platform.Unknown(w)
		return
	}
	out := shapeEnterprise([]map[string]any{row})[0]
	out["pubFlows"] = shapeEqb(pubFlows)
	out["perFlows"] = shapeEqb(perFlows)
	if len(tlzf) > 0 {
		out["tlzfConfig"] = shapeTlzfConfig(tlzf)[0]
	}
	out["saleEtcs"] = shapeTlzfSaleEtc(saleEtcs)
	out["etcs"] = shapeEnterprise(etcs)
	out["urcs"] = shapeUser(urcs)
	platform.Success(w, out)
}

// entUpdateById = ChainBasicEnterpriseServiceImpl.updateEnterprise。try 段(用户校验→
// 签约次数清空→scope→评分→updateById→认证文件置有效)任何异常 → 602;其后分支异常 → 999。
// 无事务:Java 各语句独立提交,后段失败前段保留
func (s *Service) entUpdateById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if cu := auth.From(r.Context()); !isPlatformAdmin(cu) {
		// 原先按 body 的 id 直写:任何登录用户可改任何企业。非平台管理员只能改本企业档案,
		// 且不能重绑用户(urcs)、调整集团 / 子公司(etcs)和通联对手方(saleEtcs)。
		if cu.UserID == 0 || cu.CompanyID == 0 || num(m["id"]) != float64(cu.CompanyID) ||
			m["urcs"] != nil || m["etcs"] != nil || m["saleEtcs"] != nil {
			deny(w)
			return
		}
	}
	// ── try 段:异常统一 602"操作失败,企业名称或信用代码重复" ──
	if err := s.entUpdateTry(r, m); err != nil {
		platform.Fail(w, 602, "操作失败，企业名称或信用代码重复")
		return
	}
	// ── catch 后段:异常走全局 999 ──
	id := str(m, "id")
	// filelist(file_type=3):先按 relation 删(removeFile 的库内效果),再对每项 update(恒 0 行)
	if v, ok := m["filelist"]; ok && v != nil {
		items, _ := v.([]any)
		s.removeFileByRel(id, "3")
		for _, it := range items {
			if !s.commonFileUpdate(r, it, m, "3") {
				platform.Unknown(w)
				return
			}
		}
	}
	// imagelist(file_type=4):不在 ids 内的行删除,其余项 update
	if v, ok := m["imagelist"]; ok && v != nil {
		items, _ := v.([]any)
		ids, _ := m["ids"].([]any)
		if m["ids"] == nil {
			// Java getIds() 为 null 时 contains 即 NPE → 999
			platform.Unknown(w)
			return
		}
		rows, err := queryRows(s.DB.Chain, "SELECT id FROM chain_common_file WHERE relation_id=? AND file_type='4'", id)
		if err != nil {
			platform.Unknown(w)
			return
		}
		for _, fr := range rows {
			fid := str(fr, "id")
			in := false
			for _, x := range ids {
				if str(map[string]any{"v": x}, "v") == fid {
					in = true
					break
				}
			}
			if !in {
				s.removeFileRow(fid)
			}
		}
		for _, it := range items {
			if !s.commonFileUpdate(r, it, m, "4") {
				platform.Unknown(w)
				return
			}
		}
	}
	// saleEtcs:全删重插
	if v, ok := m["saleEtcs"]; ok && v != nil {
		items, _ := v.([]any)
		if _, err := s.DB.Chain.Exec("DELETE FROM chain_tlzf_sale_etc WHERE buy_company_id=?", id); err != nil {
			platform.Unknown(w)
			return
		}
		for _, it := range items {
			em, _ := it.(map[string]any)
			if em == nil {
				continue
			}
			if _, err := platform.Insert(s.DB.Chain, "chain_tlzf_sale_etc",
				pick(em, "id", "buy_company_id", "sale_company_id", "buy_ect_name", "sale_ect_name",
					"create_time", "update_time", "remark1", "remark2", "remark3", "remark4", "remark5")); err != nil {
				platform.Unknown(w)
				return
			}
		}
	}
	// etcs:子公司——循环内 clearGroupByEtc+updGroupByEtc(Java 原文如此),再清集团侧并批量更新
	if v, ok := m["etcs"]; ok && v != nil {
		items, _ := v.([]any)
		if len(items) > 0 {
			for _, it := range items {
				em, _ := it.(map[string]any)
				if _, err := s.DB.Chain.Exec("UPDATE bus_product SET group_id=NULL WHERE group_id=?", id); err != nil {
					platform.Unknown(w)
					return
				}
				etcID := ""
				if em != nil {
					etcID = str(em, "id")
				}
				q := "UPDATE bus_product SET group_id=? WHERE company_id="
				args := []any{id}
				if etcID == "" {
					q += "NULL" // Java ${etcId} 为 null 时字面 NULL,匹配不到行
				} else {
					q += "?"
					args = append(args, etcID)
				}
				if _, err := s.DB.Chain.Exec(q, args...); err != nil {
					platform.Unknown(w)
					return
				}
			}
			if _, err := s.DB.Chain.Exec("UPDATE chain_basic_enterprise SET group_id=NULL WHERE group_id=?", id); err != nil {
				platform.Unknown(w)
				return
			}
			for _, it := range items {
				em, _ := it.(map[string]any)
				if em == nil || em["id"] == nil {
					platform.Unknown(w) // MP updateById 空 id 抛异常
					return
				}
				if _, err := platform.UpdateByIDNoFill(s.DB.Chain, "chain_basic_enterprise", em["id"], pick(em, entCols...)); err != nil {
					platform.Unknown(w)
					return
				}
			}
		} else {
			if _, err := s.DB.Chain.Exec("UPDATE chain_basic_enterprise SET group_id=NULL WHERE group_id=?", id); err != nil {
				platform.Unknown(w)
				return
			}
			if _, err := s.DB.Chain.Exec("UPDATE bus_product SET group_id=NULL WHERE group_id=?", id); err != nil {
				platform.Unknown(w)
				return
			}
		}
	}
	// urcs:先解绑全部关联用户,再逐个更新
	if v, ok := m["urcs"]; ok && v != nil {
		if _, err := s.DB.Chain.Exec("UPDATE chain_auth_user SET company_id=NULL WHERE company_id=?", id); err != nil {
			platform.Unknown(w)
			return
		}
		items, _ := v.([]any)
		for _, it := range items {
			em, _ := it.(map[string]any)
			if em == nil || em["id"] == nil {
				platform.Unknown(w)
				return
			}
			if _, err := platform.UpdateByIDNoFill(s.DB.Chain, "chain_auth_user", em["id"], pick(em, userWriteCols...)); err != nil {
				platform.Unknown(w)
				return
			}
		}
	}
	// logo:置文件有效
	if v, ok := m["logo"]; ok && v != nil {
		if _, err := s.DB.Chain.Exec("UPDATE chain_common_file SET effective=1 WHERE id=?", str(m, "logo")); err != nil {
			platform.Unknown(w)
			return
		}
	}
	// tlzfConfig:saveOrUpdate(null) 不作为;有 id 且存在→更新,否则插入
	if v, ok := m["tlzfConfig"]; ok && v != nil {
		tcm, _ := v.(map[string]any)
		if tcm != nil {
			vals := pick(tcm, "id", "company_id", "bnk_id", "acct_type", "acct_num", "cnaps_code",
				"create_time", "update_time", "remark1", "remark2", "remark3", "remark4", "remark5")
			tid := str(tcm, "id")
			if tid == "" {
				if _, err := platform.Insert(s.DB.Chain, "chain_tlzf_config", vals); err != nil {
					platform.Unknown(w)
					return
				}
			} else {
				var exists int
				if err := s.DB.Chain.QueryRow("SELECT COUNT(1) FROM chain_tlzf_config WHERE id=?", tid).Scan(&exists); err != nil {
					platform.Unknown(w)
					return
				}
				var err error
				if exists == 0 {
					_, err = platform.Insert(s.DB.Chain, "chain_tlzf_config", vals)
				} else {
					_, err = platform.UpdateByID(s.DB.Chain, "chain_tlzf_config", tid, vals)
				}
				if err != nil {
					platform.Unknown(w)
					return
				}
			}
		}
	}
	platform.Success(w, "操作成功")
}

// entUpdateTry 对照 updateEnterprise 的 try 块,返回 error 时 Java 走 602 分支
func (s *Service) entUpdateTry(r *http.Request, m map[string]any) error {
	// user = selectById(jwt.uid):查不到(level 为 null)解包即 NPE
	uid := auth.From(r.Context()).UserID
	user, err := platform.GetByID(s.DB.Chain, "chain_auth_user", uid)
	if err != nil {
		return err
	}
	if user == nil || user["level"] == nil {
		return errors.New("user level null: Java 拆箱 NPE 走 602")
	}
	if num(user["level"]) == 1 {
		// level=1 时三个签约次数置 null,NOT_NULL 策略下不参与 UPDATE
		delete(m, "signFreeCount")
		delete(m, "signSurplusCount")
		delete(m, "signChargeCount")
	}
	if v, ok := m["scope"]; ok && v != nil {
		m["isAuthentication"] = true
		// updPro:update bus_product set enterprise_name=#{...}, is_certified_etc=1 where company_id=${id}
		en := m["enterpriseName"] // 可为 null:Java #{ } 传 null → 置 NULL
		q := "UPDATE bus_product SET enterprise_name=?, is_certified_etc=1 WHERE company_id="
		args := []any{en}
		if str(m, "id") == "" {
			q += "NULL"
		} else {
			q += "?"
			args = append(args, str(m, "id"))
		}
		if _, err := s.DB.Chain.Exec(q, args...); err != nil {
			return err
		}
	}
	// 开户许可/煤仓信息非 null → 评分 1,否则 0(Java 判据是 ==null,空串也算 1)
	qua, bas := 0, 0
	if v, ok := m["accOpeningPermitId"]; ok && v != nil {
		qua = 1
	}
	if v, ok := m["ycoalSilo"]; ok && v != nil {
		bas = 1
	}
	vals := pick(m, entCols...)
	vals["qua_cert_score"] = qua
	vals["bas_info_score"] = bas
	if !isPlatformAdmin(auth.From(r.Context())) {
		for _, c := range entPlatformOnly {
			delete(vals, c) // 认证状态、签章份数、服务期、评分、商户配置由平台维护
		}
	}
	if _, err := platform.UpdateByIDNoFill(s.DB.Chain, "chain_basic_enterprise", m["id"], vals); err != nil {
		return err
	}
	// 四类认证文件置有效(isNotBlank 才处理;parseInt 失败 → 602)
	for _, k := range []string{"busLicenseId", "accOpeningPermitId", "legalFtCardId", "legalBkCardId"} {
		v, ok := m[k]
		if !ok || v == nil {
			continue
		}
		sv := str(m, k)
		if strings.TrimSpace(sv) == "" {
			continue
		}
		if _, err := strconv.Atoi(sv); err != nil {
			return err
		}
		if _, err := s.DB.Chain.Exec(
			"UPDATE chain_common_file SET effective=1 WHERE id=?", sv); err != nil {
			return err
		}
	}
	return nil
}

// removeFileByRel 对照 delByRelId({relation_id, file_type})。
func (s *Service) removeFileByRel(relID, fileType string) {
	rows, err := queryRows(s.DB.Chain, "SELECT id FROM chain_common_file WHERE relation_id=? AND file_type=?", relID, fileType)
	if err != nil {
		return
	}
	for _, fr := range rows {
		s.removeFileRow(str(fr, "id"))
	}
}

// removeFileRow removeFile 的库内效果:删文件行、清企业引用、删压缩副本行
func (s *Service) removeFileRow(fid string) {
	if fid == "" {
		return
	}
	_, _ = s.DB.Chain.Exec("DELETE FROM chain_common_file WHERE id=?", fid)
	_, _ = s.DB.Chain.Exec("UPDATE chain_basic_enterprise SET advertising_video_id=NULL WHERE advertising_video_id=?", fid)
	_, _ = s.DB.Chain.Exec("UPDATE chain_basic_enterprise SET website_img_id=NULL, website_ced_img_id=NULL WHERE website_img_id=? OR website_ced_img_id=?", fid, fid)
	// 原件因先删行后查行而遗留，转换图片仍按原 mapper 返回路径逐个删除。
	images, err := queryRows(s.DB.Chain, "SELECT file_path,file_name FROM chain_common_img WHERE file_id=?", fid)
	if err == nil {
		for _, image := range images {
			_ = os.Remove(filepath.Join(str(image, "file_path"), str(image, "file_name")))
		}
	}
	_, _ = s.DB.Chain.Exec("DELETE FROM chain_common_file WHERE source_id=?", fid)
}

// commonFileUpdate updateEnterprise filelist/imagelist 分支的单项 update:
// 强制 7 字段 + body 携带的其他非空列,WHERE id=项.id(行已被删时自然 0 行,与 Java 一致)
func (s *Service) commonFileUpdate(r *http.Request, it any, m map[string]any, fileType string) bool {
	em, _ := it.(map[string]any)
	if em == nil || em["id"] == nil {
		return false
	}
	// file_path/file_name 不接受 body:它们会被预览/下载接口直接当磁盘路径,可改写成任意文件读
	vals := pick(em, "file_real_name", "file_description", "f_type",
		"trans_status", "association_id", "is_compressed_file", "source_id", "compressed_id")
	vals["company_id"] = m["id"]
	vals["file_type"] = fileType
	vals["owner_type"] = "1"
	vals["owner_id"] = auth.From(r.Context()).UserID
	vals["relation_id"] = m["id"]
	vals["relation_type"] = 6
	vals["effective"] = true
	_, err := platform.UpdateByIDNoFill(s.DB.Chain, "chain_common_file", em["id"], vals)
	return err == nil
}

// ── 通联配置/供应商范围实体(getEtc 嵌套输出) ──

var tlzfConfigFill = []string{"id", "companyId", "bnkId", "acctType", "acctNum", "cnapsCode",
	"createTime", "updateTime", "remark1", "remark2", "remark3", "remark4", "remark5",
	"searchStr", "page", "ids", "startTime", "endTime", "orgBnkId", "orgAcctNum"}

var tlzfSaleEtcFill = []string{"id", "buyCompanyId", "saleCompanyId", "buyEctName",
	"saleEctName", "createTime", "updateTime", "remark1", "remark2", "remark3", "remark4",
	"remark5", "searchStr", "page", "ids", "startTime", "endTime"}

func shapeTlzfConfig(rows []map[string]any) []map[string]any {
	out := platform.CamelizeRows(rows)
	for _, r := range out {
		fillNulls(r, tlzfConfigFill)
	}
	return out
}

func shapeTlzfSaleEtc(rows []map[string]any) []map[string]any {
	out := platform.CamelizeRows(rows)
	for _, r := range out {
		fillNulls(r, tlzfSaleEtcFill)
	}
	return out
}

// ── 全局配置 chain_global_configuration ──

var glconfigCols = []string{"id", "is_open_pay", "is_open_etc_regist", "red_uplimit_days",
	"is_sendOfficeNotice", "is_open_member", "create_time", "update_time",
	"remark1", "remark2", "remark3", "remark4", "remark5"}

// glconfigWriteCols pick 用:body 键经 Underscore 后为 is_sendofficenotice(实体全小写),
// 与 allEq 的 DDL 原样键不同
var glconfigWriteCols = []string{"id", "is_open_pay", "is_open_etc_regist", "red_uplimit_days",
	"is_sendofficenotice", "is_open_member", "create_time", "update_time",
	"remark1", "remark2", "remark3", "remark4", "remark5"}

var glconfigFill = []string{"id", "isOpenPay", "isOpenEtcRegist", "redUplimitDays",
	"isSendofficenotice", "isOpenMember", "createTime", "updateTime",
	"remark1", "remark2", "remark3", "remark4", "remark5", "searchStr", "page", "ids"}

// glconfigExcelCols @ExcelColumn 字段(getSql 范围;布尔/长整型键走 %…% 包装)
var glconfigExcelCols = []string{"id", "isOpenPay", "isOpenEtcRegist", "redUplimitDays",
	"isSendofficenotice", "isOpenMember", "createTime", "updateTime",
	"remark1", "remark2", "remark3", "remark4", "remark5"}

var glconfigWrap = map[string]bool{"isOpenPay": true, "isOpenEtcRegist": true,
	"redUplimitDays": true, "isSendofficenotice": true, "isOpenMember": true,
	"createTime": true, "updateTime": true}

// shapeGlconfig:tinyint(1) 布尔化 + 列名大小写修正(DDL is_sendOfficeNotice,实体字段全小写)
func shapeGlconfig(rows []map[string]any) []map[string]any {
	out := platform.CamelizeRows(rows)
	for _, r := range out {
		if v, ok := r["isSendOfficeNotice"]; ok {
			delete(r, "isSendOfficeNotice")
			r["isSendofficenotice"] = v
		}
		for _, c := range []string{"isOpenPay", "isOpenEtcRegist", "isSendofficenotice", "isOpenMember"} {
			if v, ok := r[c]; ok && v != nil {
				r[c] = num(v) != 0
			}
		}
		fillNulls(r, glconfigFill)
	}
	return out
}

// glconfigList = QueryWrapper.allEq(maps,false):键原样作列名,白名单外 999
func (s *Service) glconfigList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cond, ok := allEqCond(m, glconfigCols)
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := platform.ListWhere(s.DB.Chain, "chain_global_configuration", cond)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, shapeGlconfig(rows))
}

// glconfigUpdateById = service.updateById:非空列更新,返回 success()
func (s *Service) glconfigUpdateById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !s.requirePlatform(w, r) { // 全局支付 / 会员开关
		return
	}
	if _, err := platform.UpdateByID(s.DB.Chain, "chain_global_configuration", m["id"], pick(m, glconfigWriteCols...)); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// ── e签宝配置 chain_eqb_config ──

// eqbCols pick 用写白名单:body 键为实体全小写字段(Underscore 后原样),MySQL 列名
// 不区分大小写,filepassword 即 filePassword 列;flowType→flow_type
var eqbCols = []string{"id", "plan_name", "encryption", "filepassword", "autoarchive",
	"autoinitiate", "contractremind", "contractvalidity", "batchdropseal", "billisolationcode",
	"countdown", "mobileshieldway", "noticedeveloperurl", "noticetype", "orgavailableauthtypes",
	"personauthadvancedenabled", "personavailableauthtypes", "redirectdelaytime", "redirecturl",
	"signplatform", "willtypes", "remark", "signvalidity", "platformsign", "actorindentitytype",
	"assignedposbean", "autoexecute", "pospage", "posx", "posy", "sealtype", "fontsize",
	"format", "qpospage", "qposx", "qposy", "signdatebeantype", "signtype", "width",
	"signorder", "thirdorderno", "flow_type", "forcereadtime", "create_time", "update_time",
	"remark1", "remark2", "remark3", "remark4", "remark5"}

// eqbRename DDL 驼峰列名 → 实体全小写字段(ChainEqbConfig 字段命名特例)
var eqbRename = map[string]string{
	"filePassword": "filepassword", "autoArchive": "autoarchive", "autoInitiate": "autoinitiate",
	"contractRemind": "contractremind", "contractValidity": "contractvalidity",
	"batchDropSeal": "batchdropseal", "billIsolationCode": "billisolationcode",
	"mobileShieldWay": "mobileshieldway", "noticeDeveloperUrl": "noticedeveloperurl",
	"noticeType": "noticetype", "orgAvailableAuthTypes": "orgavailableauthtypes",
	"personAuthAdvancedEnabled": "personauthadvancedenabled",
	"personAvailableAuthTypes":  "personavailableauthtypes",
	"redirectDelayTime":         "redirectdelaytime", "redirectUrl": "redirecturl",
	"signPlatform": "signplatform", "willTypes": "willtypes", "signValidity": "signvalidity",
	"platformSign": "platformsign", "actorIndentityType": "actorindentitytype",
	"assignedPosbean": "assignedposbean", "autoExecute": "autoexecute",
	"posPage": "pospage", "posX": "posx", "posY": "posy", "sealType": "sealtype",
	"fontSize": "fontsize", "qposPage": "qpospage", "qposX": "qposx", "qposY": "qposy",
	"signDateBeanType": "signdatebeantype", "signType": "signtype", "signOrder": "signorder",
	"thirdOrderNo": "thirdorderno", "forceReadTime": "forcereadtime",
}

var eqbFill = append([]string{"id", "planName", "encryption", "filepassword", "autoarchive",
	"autoinitiate", "contractremind", "contractvalidity", "batchdropseal", "billisolationcode",
	"countdown", "mobileshieldway", "noticedeveloperurl", "noticetype", "orgavailableauthtypes",
	"personauthadvancedenabled", "personavailableauthtypes", "redirectdelaytime", "redirecturl",
	"signplatform", "willtypes", "remark", "signvalidity", "platformsign", "actorindentitytype",
	"assignedposbean", "autoexecute", "pospage", "posx", "posy", "sealtype", "fontsize",
	"format", "qpospage", "qposx", "qposy", "signdatebeantype", "signtype", "width",
	"signorder", "thirdorderno", "flowType", "forcereadtime", "createTime", "updateTime",
	"remark1", "remark2", "remark3", "remark4", "remark5"}, "searchStr", "page", "ids")

// eqbExcelCols getSql 迭代范围(实体字段名,不含载体;camelToUnderline 对全小写键原样保留,
// MySQL 列名不区分大小写,filepassword 即 filePassword 列)
var eqbExcelCols = eqbFill[:len(eqbFill)-3]

var eqbWrap = map[string]bool{"autoarchive": true, "autoinitiate": true, "batchdropseal": true,
	"platformsign": true, "assignedposbean": true, "autoexecute": true,
	"contractvalidity": true, "signvalidity": true, "createTime": true, "updateTime": true}

// shapeEqb:布尔化 + 实体字段名重命名 + 载体补 null
func shapeEqb(rows []map[string]any) []map[string]any {
	out := platform.CamelizeRows(rows)
	for _, r := range out {
		for c, nk := range eqbRename {
			if v, ok := r[c]; ok {
				delete(r, c)
				r[nk] = v
			}
		}
		for _, c := range []string{"autoarchive", "autoinitiate", "batchdropseal",
			"platformsign", "assignedposbean", "autoexecute"} {
			if v, ok := r[c]; ok && v != nil {
				r[c] = num(v) != 0
			}
		}
		fillNulls(r, eqbFill)
	}
	return out
}

// eqbListByPage = mapper.listByPage:getSql(全 ExcelColumn)+ CONCAT(id) like '%searchStr%'
func (s *Service) eqbListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if _, ok := m["page"].(map[string]any); !ok {
		platform.Unknown(w) // Java mapper requires a Page object.
		return
	}
	current, size := pageOf(m)
	searchSql, args := getSqlEq(m, eqbExcelCols, map[string]bool{}, eqbWrap)
	q := "SELECT * FROM chain_eqb_config " + searchSql + " AND CONCAT(id) LIKE ? ORDER BY create_time DESC"
	countQ := "SELECT COUNT(*) FROM chain_eqb_config " + searchSql + " AND CONCAT(id) LIKE ?"
	args = append(args, "%"+searchStrOf(m)+"%")
	page, err := paged(s.DB.Chain, q, countQ, args, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	page["records"] = shapeEqb(page["records"].([]map[string]any))
	platform.Success(w, page)
}

// echoEntity 对照 add/updateById 返回 success(实体):按实体全键回显 body(缺键补 null)。
// add 时 id 用数据库自增回填(对照 MP save 的主键回写)
func echoEntity(m map[string]any, keys []string, id any) map[string]any {
	out := make(map[string]any, len(keys))
	for _, k := range keys {
		if v, ok := m[k]; ok {
			out[k] = v
		} else {
			out[k] = nil
		}
	}
	if id != nil {
		out["id"] = id
	}
	return out
}

func (s *Service) eqbGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	row, err := platform.GetByID(s.DB.Chain, "chain_eqb_config", m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	if row == nil {
		platform.Success(w, nil)
		return
	}
	platform.Success(w, shapeEqb([]map[string]any{row})[0])
}

func (s *Service) eqbAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !s.requirePlatform(w, r) { // e签宝 / 通联商户配置
		return
	}
	id, err := platform.InsertNoFill(s.DB.Chain, "chain_eqb_config", pick(m, eqbCols...))
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, echoEntity(m, eqbFill, id))
}

func (s *Service) eqbUpdateById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !s.requirePlatform(w, r) { // e签宝 / 通联商户配置
		return
	}
	if _, err := platform.UpdateByIDNoFill(s.DB.Chain, "chain_eqb_config", m["id"], pick(m, eqbCols...)); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, echoEntity(m, eqbFill, nil))
}

func (s *Service) eqbDelBatch(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !s.requirePlatform(w, r) { // e签宝 / 通联商户配置
		return
	}
	ids, _ := m["ids"].([]any)
	if err := platform.DeleteByIDs(s.DB.Chain, "chain_eqb_config", ids); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功")
}

// ── 文件上传 /saas-common/upload/batch ──

// uploadBatch 对照 ChainCommonFileServiceImpl.upload:存文件 + 插 chain_common_file,
// 图片压缩副本(compressedId)、excel 转 png 与 Java 同源(共享 platform 实现);
// doc 转 pdf 使用本地文档渲染器，转换异常保留原接口中断行为。
// 整个循环在 Java 同一个 try 内:任何异常吞掉并返回已处理部分(这里用 break 对齐)。
func (s *Service) uploadBatch(w http.ResponseWriter, r *http.Request) {
	cu := auth.From(r.Context())
	user, err := platform.GetByID(s.DB.Chain, "chain_auth_user", cu.UserID)
	if err != nil || user == nil {
		platform.Unknown(w) // Java selectById 后 getPhone() NPE
		return
	}
	apath := ""
	if cu.EnterpriseCode == "" {
		apath = filepath.Base(str(user, "phone")) // phone 用户可自改:带 ../ 会把上传写到根目录之外
	}
	dir := filepath.Join(uploadDir(), filepath.FromSlash(apath))
	cfiles := []map[string]any{}
	r.Body = http.MaxBytesReader(w, r.Body, 100<<20) // 营业执照/附件远小于 100MB;原先无上限,超大"图片"整读再 base64 能占数 GB 内存
	if err := r.ParseMultipartForm(32 << 20); err != nil || r.MultipartForm == nil {
		platform.Success(w, cfiles)
		return
	}
	fhs := r.MultipartForm.File["file"]
	for _, fh := range fhs {
		src, err := fh.Open()
		if err != nil {
			break
		}
		suffix := filepath.Ext(fh.Filename)
		fileName := newUUID() + suffix
		if err := os.MkdirAll(dir, 0o755); err != nil {
			src.Close()
			break
		}
		dst, err := os.Create(filepath.Join(dir, fileName))
		if err != nil {
			src.Close()
			break
		}
		_, cpErr := io.Copy(dst, src)
		dst.Close()
		src.Close()
		if cpErr != nil {
			break
		}
		vals := map[string]any{
			"company_id": cu.CompanyID, "file_real_name": fh.Filename, "file_name": fileName,
			"file_path": dir, "file_type": "", "owner_type": "", "owner_id": cu.UserID,
		}
		id, err := platform.InsertNoFill(s.DB.Chain, "chain_common_file", vals)
		if err != nil {
			break
		}
		cid := int64(0)
		if isImgExt(suffix) {
			cid = platform.CompressFileRow(s.DB.Chain, id)
			if _, err := platform.UpdateByIDNoFill(s.DB.Chain, "chain_common_file", id, map[string]any{"compressed_id": cid}); err != nil {
				break
			}
		}
		if strings.Contains(fileName, ".xls") || strings.Contains(fileName, ".xlsx") {
			imgdir := filepath.Join(dir, strings.SplitN(fileName, ".", 2)[0]+"_img")
			imgs, err := platform.RenderExcelUpload(r.Context(), filepath.Join(dir, fileName), imgdir)
			if err != nil {
				break
			}
			for _, img := range imgs {
				if _, err = platform.InsertNoFill(s.DB.Chain, "chain_common_img", map[string]any{"file_id": id, "file_name": img, "file_path": imgdir}); err != nil {
					break
				}
			}
			if err != nil {
				break
			}
		}
		if strings.Contains(fileName, ".doc") || strings.Contains(fileName, ".docx") {
			root := os.Getenv("PANDA_UPLOAD_DIR")
			if root == "" {
				root = "data/uploads"
			}
			imgdir := filepath.Join(root, "convertFiles")
			preview := platform.PreviewName()
			if err = platform.RenderDocumentPDF(r.Context(), filepath.Join(dir, fileName), filepath.Join(imgdir, preview)); err != nil {
				break
			}
			if _, err = platform.InsertNoFill(s.DB.Chain, "chain_common_img", map[string]any{"file_id": id, "file_name": preview, "file_path": imgdir}); err != nil {
				break
			}
		}
		item := map[string]any{"id": strconv.FormatInt(id, 10), "fileName": fileName, "name": fh.Filename}
		if isImgExt(suffix) {
			item["compressedId"] = strconv.FormatInt(cid, 10)
			if b, err := os.ReadFile(filepath.Join(dir, fileName)); err == nil {
				item["url"] = "data:image/jpg;base64," + base64.StdEncoding.EncodeToString(b)
			}
		}
		cfiles = append(cfiles, item)
	}
	platform.Success(w, cfiles)
}

func uploadDir() string {
	if d := os.Getenv("PANDA_UPLOAD_DIR"); d != "" {
		return d
	}
	return "data/uploads"
}

// newUUID 32 位随机 hex(对照 Java UUID 文件名,仅作唯一名使用)
func newUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// isImgExt 对照 FileHelper.isImg(ImageIO 可解码集的常见后缀)
func isImgExt(suffix string) bool {
	switch strings.ToLower(suffix) {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
		return true
	}
	return false
}
