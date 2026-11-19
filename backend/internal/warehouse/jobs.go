package warehouse

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

type productJob struct {
	ProID, UserID, OriginalStatus string
	First                         bool
}

func (s *Service) RegisterJobs(jobs *platform.Jobs) {
	s.Jobs = jobs
	jobs.Handle("warehouse.product", func(id string, payload json.RawMessage) {
		var job productJob
		if err := json.Unmarshal(payload, &job); err != nil {
			log.Printf("煤种消息格式错误: %v", err)
			return
		}
		// 原消费者固定使用 warehouse 的数据库，不继承前端 applytype。
		r, _ := http.NewRequestWithContext(context.Background(), "POST", "http://localhost/internal-job", nil)
		fresh, err := jobs.MarkDelivery(id, time.Hour)
		if err != nil {
			log.Printf("煤种消费标记失败: %v", err)
			return
		}
		if !fresh {
			return
		}
		time.Sleep(2 * time.Second) // 原消费者在同步前固定等待两秒。
		if err := s.bpSyncOtherInfo(r, job.ProID, job.UserID, job.First, job.OriginalStatus); err != nil {
			log.Printf("煤种后续同步失败: %v", err)
		}
	})
}

func (s *Service) submitProductJob(r *http.Request, job productJob) {
	if s.Jobs != nil {
		if err := s.Jobs.Submit("warehouse.product", job); err != nil {
			log.Printf("煤种后台消息投递失败: %v", err)
		}
		return
	}
	go func() {
		if err := s.bpSyncOtherInfo(r, job.ProID, job.UserID, job.First, job.OriginalStatus); err != nil {
			log.Printf("煤种后续同步失败: %v", err)
		}
	}()
}
