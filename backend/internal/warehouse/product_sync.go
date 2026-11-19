package warehouse

import (
	"database/sql"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// 原消息消费者与主记录更新分别执行，消费失败不会回滚前端保存。
func (s *Service) bpSyncOtherInfo(r *http.Request, proID, userID string, isFirst bool, orgStatus string) error {
	d := s.pick(r)
	row, err := platform.GetByID(d, "bus_product", proID)
	if err != nil {
		return err
	}
	if row == nil {
		return fmt.Errorf("Java 同步煤种为空")
	}
	p := platform.Camelize(row)
	if p["status"] == nil {
		return fmt.Errorf("Java 煤种状态为空")
	}
	switch str(p, "status") {
	case "0":
		old := s.bpLastPrice(d, proID, "1")
		if _, err = bpPriceDifference(p["exFactoryPrice"], old); err != nil {
			return err
		}
		if err = s.bpNotifyPending(r, p, old, userID); err != nil {
			return err
		}
	case "1":
		if err = s.bpSyncHistory(r, p, userID); err != nil {
			return err
		}
		var count int
		if err = d.QueryRow("SELECT count(*) FROM bus_product WHERE is_real_check=1 AND company_id=?", p["companyId"]).Scan(&count); err != nil {
			return err
		}
		score := 0
		if count > 0 {
			score = 1
		}
		if _, err = d.Exec("UPDATE chain_basic_enterprise SET coal_verify_score=? WHERE id=?", score, p["companyId"]); err != nil {
			return err
		}
		if _, err = d.Exec("UPDATE bus_videos SET status=1 WHERE pro_id=?", proID); err != nil {
			return err
		}
	}
	if p["companyId"] != nil {
		enterprise, err := platform.GetByID(d, "chain_basic_enterprise", p["companyId"])
		if err != nil {
			return err
		}
		if enterprise == nil {
			return fmt.Errorf("Java 企业为空")
		}
		p["isCertifiedEtc"] = num(enterprise["is_authentication"]) == 1
		if enterprise["group_id"] != nil {
			p["groupId"] = enterprise["group_id"]
		}
	}
	if p["bigProId"] != nil && p["secName"] == nil {
		for _, f := range []struct{ ID, Name string }{{"bigProId", "secName"}, {"threeProId", "thName"}} {
			product, err := platform.GetByID(d, "bus_big_prodect", p[f.ID])
			if err != nil {
				return err
			}
			if product == nil {
				return fmt.Errorf("Java 煤种分类为空")
			}
			p[f.Name] = product["product_name"]
		}
	}
	if _, err = platform.UpdateByID(d, "bus_product", proID, pick(p, entBusProduct.dbcols)); err != nil {
		return err
	}
	if str(p, "status") != "0" {
		if orgStatus == "0" {
			msg := ""
			switch str(p, "status") {
			case "1":
				msg = "煤种信息审核通过,快去查看吧!"
			case "2":
				msg = "煤种信息审核未通过,快去查看吧!"
				if err = s.bpSyncScore(d, p, 5, 1); err != nil {
					return err
				}
			}
			msgID, err := platform.Insert(d, "bus_msg", map[string]any{"msg_type": 3, "content": msg, "remark": proID, "jump_id": proID})
			if err != nil {
				return err
			}
			if _, err = platform.Insert(d, "bus_user_msg", map[string]any{"user_id": p["createUserId"], "msg_id": msgID}); err != nil {
				return err
			}
		}
		if err = s.bpNotifyReview(r, p, orgStatus); err != nil {
			return err
		}
	}
	if err = s.bpSyncEnterpriseScore(d, p); err != nil {
		return err
	}
	if isFirst {
		if err = s.bpSyncScore(d, p, 5, 1); err != nil {
			return err
		}
		if err = s.bpSyncReward(r, p, userID, 2); err != nil {
			return err
		}
		bpBroadcastAppend("新增" + platform.JnullAny(p["warehouseNb"]) + ";")
		if err = s.bpNotifyAdded(r, p); err != nil {
			return err
		}
	}
	return nil
}

func bpPriceDifference(a, b any) (*big.Rat, error) {
	x, ok := new(big.Rat).SetString(platform.JnullAny(a))
	if !ok {
		return nil, fmt.Errorf("Java 价格为空或无效")
	}
	y, ok := new(big.Rat).SetString(platform.JnullAny(b))
	if !ok {
		return nil, fmt.Errorf("Java 旧价格为空或无效")
	}
	return x.Sub(x, y), nil
}

func (s *Service) bpOriginalUpdPriceScore(w http.ResponseWriter, d *sql.DB, bp, m map[string]any) int {
	changed := false
	for _, f := range []struct{ field, kind string }{{"exFactoryPrice", "1"}, {"platePrice", "2"}} {
		old := s.bpLastPrice(d, str(m, "id"), f.kind)
		if old != "" && m[f.field] != nil {
			diff, err := bpPriceDifference(m[f.field], old)
			if err != nil {
				platform.Unknown(w)
				return -1
			}
			if diff.Sign() != 0 {
				changed = true
			}
		}
	}
	if changed {
		if str(bp, "status") == "1" && str(m, "status") == "0" {
			return s.bpUpdUserScore(w, d, m, 3, 2)
		}
		if str(bp, "status") == "0" && str(m, "status") == "1" {
			return s.bpUpdUserScore(w, d, m, 6, 1)
		}
	}
	return 0
}

func (s *Service) bpSyncScore(d *sql.DB, p map[string]any, getType, scoreType int) error {
	w := httptest.NewRecorder()
	if s.bpUpdUserScore(w, d, p, getType, scoreType) < 0 {
		return fmt.Errorf("Java 煤种积分更新异常")
	}
	return nil
}

func (s *Service) bpSyncHistory(r *http.Request, p map[string]any, userID string) error {
	d := s.pick(r)
	factoryChanged := false
	for _, f := range []struct{ field, previous, updated, kind string }{{"exFactoryPrice", "preExFactoryPrice", "fpriceUpdTime", "1"}, {"platePrice", "prePlatePrice", "ppriceUpdTime", "2"}} {
		old := s.bpLastPrice(d, str(p, "id"), f.kind)
		value := p[f.field]
		if value == nil {
			continue
		}
		changed := false
		if old != "" {
			difference, err := bpPriceDifference(value, old)
			if err != nil {
				return err
			}
			changed = difference.Sign() != 0
			if !changed {
				continue
			}
		}
		if changed {
			if f.kind == "1" {
				factoryChanged = true
			}
			p[f.previous], p[f.updated], p["lastCheckFacPrice"] = old, nowSpace(), p["exFactoryPrice"] // 车板价也写出厂价，保留原问题。
			if f.kind == "1" {
				if err := s.bpSyncPriceMessage(d, p, old, value); err != nil {
					return err
				}
				if err := s.bpSyncReward(r, p, userID, 8); err != nil {
					return err
				}
			}
		}
		previous := any(old)
		if old == "" {
			previous = "0"
		}
		if _, err := platform.Insert(d, "bus_price_log", map[string]any{"user_id": userID, "fst_pro_id": p["fstProId"], "big_pro_id": p["bigProId"], "three_pro_id": p["threeProId"], "pro_id": p["id"], "pre_price": previous, "price": value, "price_type": f.kind}); err != nil {
			return err
		}
		if changed {
			diff, _ := bpPriceDifference(value, old)
			verb := "上涨到"
			if diff.Sign() < 0 {
				verb = "下跌到"
			}
			label := "（出厂价）"
			if f.kind == "2" {
				label = "（车板价）"
			}
			bpBroadcastAppend(platform.JnullAny(p["warehouseNb"]) + " " + old + "元/吨 " + verb + " " + platform.JnullAny(value) + "元/吨" + label + ";")
			if f.kind == "1" {
				if err := s.bpNotifyPrice(r, p, old, value); err != nil {
					return err
				}
			} else if !factoryChanged {
				if err := s.bpOfficialNotice(r, p, old, value, false); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *Service) bpSyncPriceMessage(d *sql.DB, p map[string]any, old, value any) error {
	if p["warehouseNb"] == nil {
		return nil
	}
	diff, err := bpPriceDifference(value, old)
	if err != nil {
		return err
	}
	verb := "上升了|"
	if diff.Sign() < 0 {
		verb = "下降了|"
	}
	amount := new(big.Rat).Abs(diff).FloatString(2)
	msg := str(p, "warehouseNb") + "|价格|" + platform.JnullAny(value) + "|" + verb + amount + "|元"
	remark := ""
	for _, f := range []struct{ field, label string }{{"maxAshContent", "灰分Ad%≤"}, {"maxWaterContent", "水分Mt%≤"}, {"maxStd", "全硫St.d%≤"}, {"maxVolatiles", "挥发分Vdaf%≤"}, {"minBondIndex", "粘结指数G≥"}} {
		if p[f.field] != nil {
			if remark != "" {
				remark += "，"
			}
			remark += f.label + platform.JnullAny(p[f.field])
		}
	}
	// 原 Java 首个指标不是灰分时会删掉全部中文逗号。
	if p["maxAshContent"] == nil {
		remark = strings.ReplaceAll(remark, "，", "")
	}
	id, err := platform.Insert(d, "bus_msg", map[string]any{"msg_type": 1, "content": msg, "remark": remark, "jump_id": p["id"], "is_effect": true})
	if err != nil {
		return err
	}
	collections, err := queryRows(d, "SELECT user_id FROM bus_pro_collection WHERE pro_id=?", p["id"])
	if err != nil {
		return err
	}
	for _, coll := range collections {
		if _, err = platform.Insert(d, "bus_user_msg", map[string]any{"user_id": coll["user_id"], "msg_id": id}); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) bpSyncReward(r *http.Request, p map[string]any, userID string, kind int) error {
	d := s.pick(r)
	if kind == 8 {
		if _, ok := s.redCachePut(str(p, "id")+"t8", time.Hour, false); !ok {
			return nil
		}
		var count int
		if err := d.QueryRow("SELECT count(*) FROM bus_price_log WHERE pro_id=? AND user_id=? AND YEARWEEK(create_time,1)=YEARWEEK(CURDATE(),1)", p["id"], p["createUserId"]).Scan(&count); err != nil {
			return err
		}
		if count > 0 {
			return nil
		}
	}
	amount, err := randomRedAmount(d, kind)
	if err != nil {
		return err
	}
	first := map[string]any{"user_id": userID, "receive_type": kind, "create_amount": amount, "is_write_off": false}
	if _, err = platform.Insert(d, "bus_red_money_clog", first); err != nil {
		return err
	}
	if kind == 2 {
		var count int
		if err = d.QueryRow("SELECT count(*) FROM bus_red_envelope_record WHERE user_id=? AND receive_type IN(2,10) AND TO_DAYS(create_time)=TO_DAYS(now())", p["createUserId"]).Scan(&count); err != nil {
			return err
		}
		if count >= 3 {
			return nil
		}
	}
	id, err := platform.Insert(d, "bus_red_money_clog", first)
	if err != nil {
		return err
	}
	_, _, _ = s.redReceive(r, map[string]any{"clogId": id, "receiveType": kind, "receiveAmount": amount, "userId": p["createUserId"], "userPhone": p["createPhone"]})
	return nil
}

func (s *Service) bpSyncEnterpriseScore(d *sql.DB, p map[string]any) error {
	company := p["companyId"]
	if _, err := d.Exec("UPDATE chain_basic_enterprise SET upd_share_score=0,report_score=0,upd_price_score=0 WHERE id=? AND EXISTS(SELECT 1 FROM bus_product WHERE company_id=chain_basic_enterprise.id)", company); err != nil {
		return err
	}
	for _, query := range []string{
		"UPDATE chain_basic_enterprise SET upd_share_score=0.5 WHERE id=? AND EXISTS(SELECT 1 FROM bus_product p WHERE p.id=? AND p.status=1 AND p.company_id=chain_basic_enterprise.id AND ((p.fprice_upd_time IS NOT NULL AND TO_DAYS(now())-TO_DAYS(p.fprice_upd_time)<=30) OR (p.pprice_upd_time IS NOT NULL AND TO_DAYS(now())-TO_DAYS(p.pprice_upd_time)<=30)) AND EXISTS(SELECT 1 FROM bus_pro_log l WHERE l.pro_id=p.id AND l.operate_obj=4 AND l.operate_type IN(4,5) AND l.user_id=p.create_user_id))",
		"UPDATE chain_basic_enterprise SET report_score=0.25 WHERE id=? AND EXISTS(SELECT 1 FROM bus_product p WHERE p.id=? AND p.status=1 AND p.company_id=chain_basic_enterprise.id AND EXISTS(SELECT 1 FROM chain_common_file f WHERE f.association_id=p.association_id AND f.f_type=2 AND f.effective=1))",
	} {
		if _, err := d.Exec(query, company, p["id"]); err != nil {
			return err
		}
	}
	return nil
}

// ponytail: 单进程替代 Redis 广播缓存，多实例部署时需共享 Redis。
var productBroadcast struct {
	sync.Mutex
	day, text string
}

func bpBroadcastAppend(text string) {
	productBroadcast.Lock()
	defer productBroadcast.Unlock()
	day := time.Now().Format("2006-01-02") + os.Getenv("PANDA_WAREHOUSE_ENVIRONMENT_TYPE")
	if productBroadcast.day != day {
		productBroadcast.day, productBroadcast.text = day, ""
	}
	productBroadcast.text += text
}
func bpBroadcastText() string {
	productBroadcast.Lock()
	defer productBroadcast.Unlock()
	day := time.Now().Format("2006-01-02") + os.Getenv("PANDA_WAREHOUSE_ENVIRONMENT_TYPE")
	if productBroadcast.day != day {
		return ""
	}
	return productBroadcast.text
}
