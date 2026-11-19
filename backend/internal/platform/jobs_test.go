package platform

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestJobsResumePersistedDelivery(t *testing.T) {
	root := t.TempDir()
	first := NewJobs(root)
	if err := first.Submit("fixture", map[string]any{"id": "fixture"}); err != nil {
		t.Fatal(err)
	}
	restarted := NewJobs(root)
	delivered := make(chan string, 1)
	restarted.Handle("fixture", func(id string, payload json.RawMessage) {
		if len(id) != 32 {
			t.Errorf("消息ID超出原表限制: %s", id)
		}
		delivered <- string(payload)
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go restarted.Run(ctx)
	select {
	case payload := <-delivered:
		if payload != `{"id":"fixture"}` {
			t.Fatal(payload)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("重启后未恢复已持久化消息")
	}
}

func TestJobsKeepOriginalPreConsumptionMarkerAcrossRestart(t *testing.T) {
	root := t.TempDir()
	first := NewJobs(root)
	id := strings.Repeat("a", 32)
	if fresh, err := first.MarkDelivery(id, time.Hour); err != nil || !fresh {
		t.Fatalf("初次标记失败: %v", err)
	}
	restarted := NewJobs(root)
	if fresh, err := restarted.MarkDelivery(id, time.Hour); err != nil || fresh {
		t.Fatalf("重启后不应重复消费: fresh=%v err=%v", fresh, err)
	}
}
