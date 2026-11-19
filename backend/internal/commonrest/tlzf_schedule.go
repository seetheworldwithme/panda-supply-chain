package commonrest

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

const tlTaskName = "通联同步交易状态任务"

func (s *Service) startTLTask() {
	ctx, cancel := context.WithCancel(context.Background())
	s.tlTasksMu.Lock()
	if s.tlTasks == nil {
		s.tlTasks = map[string]context.CancelFunc{}
	}
	s.tlTasks[tlTaskName] = cancel // 保留 Java 重复 start 覆盖缓存而未取消旧任务的行为。
	s.tlTasksMu.Unlock()
	go func() {
		for {
			timer := time.NewTimer(time.Until(time.Now().Truncate(5 * time.Second).Add(5 * time.Second)))
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
			}
			r := tlRequest(ctx, "zjjg")
			if err := s.syncTLOrders(r, s.DB.Zjjg, false); err != nil {
				log.Printf("通联定时同步失败: %v", err)
			}
		}
	}()
}

func (s *Service) stopTLTasks(id string, all bool) {
	s.tlTasksMu.Lock()
	defer s.tlTasksMu.Unlock()
	if all {
		for _, cancel := range s.tlTasks {
			cancel()
		}
		return
	} // Java stopAll 不移除缓存。
	if cancel := s.tlTasks[id]; cancel != nil {
		cancel()
		delete(s.tlTasks, id)
	}
}

func (s *Service) syncTLOrders(r *http.Request, d *sql.DB, force bool) error {
	q := "SELECT * FROM chain_tlzf_order WHERE trans_status=0"
	var args []any
	if force {
		m := bodyMap(r)
		if m["contractId"] == nil {
			q += " AND contract_id IS NULL"
		} else {
			q += " AND contract_id=?"
			args = append(args, m["contractId"])
		}
	}
	rows, err := queryRows(d, q, args...)
	if err != nil {
		return err
	}
	count := 0
	for _, row := range rows {
		if row["sel_retry_times"] == nil {
			return fmt.Errorf("Java selRetryTimes 为 null")
		}
		if num(row["sel_retry_times"]) < 10 || force {
			count++
			if err := s.submitTLJob(str(row, "mch_order_nbr"), false); err != nil {
				return err
			}
		}
	}
	if count == 0 {
		s.stopTLTasks("", true)
	}
	return nil
}

func (s *Service) tlSyncStatus(r *http.Request, d *sql.DB, row map[string]any, async bool) error {
	// ponytail: 单进程全局锁替代 Redisson 合同锁；多实例部署时需使用同一分布式锁。
	s.tlSyncMu.Lock()
	defer s.tlSyncMu.Unlock()
	select {
	case <-r.Context().Done():
		return r.Context().Err()
	case <-time.After(time.Second):
	}
	rows, err := queryRows(d, "SELECT * FROM chain_tlzf_order WHERE mch_order_nbr <=> ?", row["mch_order_nbr"])
	if err != nil {
		return err
	}
	if len(rows) != 1 {
		return fmt.Errorf("Java 订单 getOne 失败")
	}
	row = rows[0]
	if row["trans_status"] == nil {
		return fmt.Errorf("Java transStatus 为 null")
	}
	if num(row["trans_status"]) == 1 {
		return nil
	}
	var balanceText sql.NullString
	if err = d.QueryRow("SELECT sum(residue_freeze_amt) FROM chain_contract WHERE sales_unit_id=?", row["sale_company_id"]).Scan(&balanceText); err != nil {
		return err
	}
	balance := big.NewRat(0, 1)
	if balanceText.Valid {
		balance, err = tlAmount(balanceText.String)
		if err != nil {
			return err
		}
	}
	kind := num(row["trans_type"])
	result := "0000"
	if async {
		if kind == 1 {
			result = "1122"
		}
	} else {
		response, err := s.tlQueryOrder(r, d, row)
		if err != nil {
			return err
		}
		if response["qur_rst"] == nil {
			return fmt.Errorf("Java qur_rst 为 null")
		}
		result = str(response, "qur_rst")
	}
	contract, err := platform.GetByID(d, "chain_contract", row["contract_id"])
	if err != nil {
		return err
	}
	if contract == nil {
		return fmt.Errorf("Java contract 为 null")
	}
	var feeText string
	if err = d.QueryRow("SELECT IFNULL(sum(act_trace_fee),0) FROM chain_tlzf_order WHERE trans_type=1 AND trans_status=1 AND contract_id=?", row["contract_id"]).Scan(&feeText); err != nil {
		return err
	}
	fee, err := tlAmount(feeText)
	if err != nil {
		return err
	}
	changed := (kind == 1 && result == "1122") || ((kind == 2 || kind == 3) && result == "0000")
	if changed {
		frozen, err := tlAmount(contract["sum_freeze_amt"])
		if err != nil {
			return err
		}
		thawed, err := tlAmount(contract["sum_thaw_amt"])
		if err != nil {
			return err
		}
		refunded, err := tlAmount(contract["sum_refund_amt"])
		if err != nil {
			return err
		}
		amountKey := "amt_jf"
		if kind == 1 {
			amountKey = "amt_tran"
		}
		amount, err := tlAmount(row[amountKey])
		if err != nil {
			return err
		}
		money := new(big.Rat).Set(amount)
		status := contract["contract_status"]
		switch kind {
		case 1:
			frozen.Add(frozen, amount)
			balance.Add(balance, amount)
			status = 2
		case 2:
			operationFee, err := tlAmount(row["act_trace_fee"])
			if err != nil {
				return err
			}
			thawed.Add(thawed, amount)
			balance.Sub(balance, amount)
			balance.Sub(balance, operationFee)
			money.Add(money, operationFee)
		case 3:
			refunded.Add(refunded, amount)
			balance.Sub(balance, amount)
		}
		remaining := new(big.Rat).Sub(frozen, thawed)
		remaining.Sub(remaining, refunded)
		remaining.Sub(remaining, fee)
		if kind != 1 {
			sum, err := tlAmount(contract["sum_price"])
			if err != nil {
				return err
			}
			if sum.Cmp(new(big.Rat).Add(thawed, refunded)) == 0 {
				if kind == 2 && remaining.Sign() == 0 {
					status = 4
				}
				if kind == 3 {
					var count int
					if err = d.QueryRow("SELECT count(*) FROM chain_tlzf_order WHERE trans_type=2 AND trans_status=1").Scan(&count); err != nil {
						return err
					}
					if count == 0 {
						status = 3
					}
				} // Java 查询未按合同过滤，照原样保留。
			}
			if status == nil {
				return fmt.Errorf("Java contractStatus 为 null")
			}
		}
		if _, err = d.Exec("UPDATE chain_contract SET contract_status=?,sum_freeze_amt=?,sum_thaw_amt=?,sum_refund_amt=?,residue_freeze_amt=?,sum_trace_fees=? WHERE id=?", status, tlMoney(frozen), tlMoney(thawed), tlMoney(refunded), tlMoney(remaining), tlMoney(fee), contract["id"]); err != nil {
			return err
		}
		if _, err = platform.Insert(d, "chain_tlzf_contract_log", map[string]any{"contract_id": contract["id"], "contract_code": contract["contract_code"], "org_sum_freeze_amt": contract["sum_freeze_amt"], "org_sum_thaw_amt": contract["sum_thaw_amt"], "org_sum_refund_amt": contract["sum_refund_amt"], "org_residue_freeze_amt": contract["residue_freeze_amt"], "sum_freeze_amt": tlMoney(frozen), "sum_thaw_amt": tlMoney(thawed), "sum_refund_amt": tlMoney(refunded), "residue_freeze_amt": tlMoney(remaining), "order_id": row["id"], "order_money": tlMoney(money)}); err != nil {
			return err
		}
		row["trans_status"], row["amt_ye"], row["update_time"] = 1, tlMoney(balance), nowSpace()
	}
	if row["sel_retry_times"] == nil {
		return fmt.Errorf("Java selRetryTimes 为 null")
	}
	row["sel_retry_times"] = int(num(row["sel_retry_times"])) + 1
	_, err = platform.UpdateByID(d, "chain_tlzf_order", row["id"], insertVals(pick(row, entChainTlzfOrder.dbcols)))
	return err
}
