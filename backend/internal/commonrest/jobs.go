package commonrest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

type tlJob struct {
	Order string
	Async bool
}

func (s *Service) RegisterJobs(jobs *platform.Jobs) {
	s.Jobs = jobs
	jobs.Handle("common.tlzf", func(id string, payload json.RawMessage) {
		var job tlJob
		if err := json.Unmarshal(payload, &job); err != nil {
			log.Printf("通联消息格式错误: %v", err)
			return
		}
		if err := s.consumeTLJob(job); err != nil {
			log.Printf("通联消息处理失败: %v", err)
		}
		// 原 common 消费者捕获交易异常后仍记录成功消费。
		_, _ = platform.Insert(s.DB.Zjjg, "bus_mq_consume_log", map[string]any{"msg_id": id, "msg_content": string(payload)})
	})
}

func (s *Service) consumeTLJob(job tlJob) error {
	d := s.DB.Zjjg // 原 RocketMqBeans 消费者明确 setDataSource("zjjg")。
	rows, err := queryRows(d, "SELECT * FROM chain_tlzf_order WHERE mch_order_nbr=?", job.Order)
	if err != nil {
		return err
	}
	if len(rows) != 1 {
		return fmt.Errorf("Java 通联订单为空或不唯一")
	}
	return s.tlSyncStatus(tlRequest(context.Background(), "zjjg"), d, rows[0], job.Async)
}

func (s *Service) submitTLJob(order string, async bool) error {
	job := tlJob{Order: order, Async: async}
	if s.Jobs != nil {
		return s.Jobs.Submit("common.tlzf", job)
	}
	go func() {
		if err := s.consumeTLJob(job); err != nil {
			log.Printf("通联消息处理失败: %v", err)
		}
	}()
	return nil
}
