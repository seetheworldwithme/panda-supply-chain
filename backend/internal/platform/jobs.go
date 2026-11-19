package platform

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Jobs 保留原 MQ 的异步、持久化投递；消费者与原 Java 一样，业务异常记录后确认消费。
// 队列目录随业务数据备份，运行一个 Go 实例即可，无需继续运行 Java MQ 客户端。
type Jobs struct {
	dir      string
	mu       sync.RWMutex
	handlers map[string]func(string, json.RawMessage)
}

func NewJobs(dir string) *Jobs {
	return &Jobs{dir: dir, handlers: map[string]func(string, json.RawMessage){}}
}
func (q *Jobs) Handle(name string, handler func(string, json.RawMessage)) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.handlers[name] = handler
}

type persistedJob struct {
	Name    string          `json:"name"`
	Payload json.RawMessage `json:"payload"`
}

func (q *Jobs) Submit(name string, payload any) error {
	if err := os.MkdirAll(q.dir, 0700); err != nil {
		return err
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	data, err := json.Marshal(persistedJob{Name: name, Payload: raw})
	if err != nil {
		return err
	}
	file, err := os.CreateTemp(q.dir, ".pending-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	if _, err = file.Write(data); err == nil {
		err = file.Sync()
	}
	closeErr := file.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	var random [16]byte
	if _, err = rand.Read(random[:]); err != nil {
		return err
	}
	target := filepath.Join(q.dir, fmt.Sprintf("%020d-%s.json", time.Now().UnixNano(), hex.EncodeToString(random[:])))
	if err = os.Rename(file.Name(), target); err != nil {
		return err
	}
	dir, err := os.Open(q.dir)
	if err != nil {
		return err
	}
	defer dir.Close()
	return dir.Sync()
}

// 原 warehouse 在业务前写一小时消费标记，业务异常后重投也会跳过。
func (q *Jobs) MarkDelivery(id string, ttl time.Duration) (bool, error) {
	path := filepath.Join(q.dir, id+".seen")
	data, err := os.ReadFile(path)
	if err == nil {
		var expiry int64
		if json.Unmarshal(data, &expiry) != nil {
			return false, fmt.Errorf("消费标记损坏")
		}
		if time.Now().UnixNano() < expiry {
			return false, nil
		}
	} else if !os.IsNotExist(err) {
		return false, err
	}
	if err := os.MkdirAll(q.dir, 0700); err != nil {
		return false, err
	}
	file, err := os.CreateTemp(q.dir, ".seen-")
	if err != nil {
		return false, err
	}
	defer os.Remove(file.Name())
	data, _ = json.Marshal(time.Now().Add(ttl).UnixNano())
	if _, err = file.Write(data); err == nil {
		err = file.Sync()
	}
	closeErr := file.Close()
	if err != nil {
		return false, err
	}
	if closeErr != nil {
		return false, closeErr
	}
	if err = os.Rename(file.Name(), path); err != nil {
		return false, err
	}
	dir, err := os.Open(q.dir)
	if err != nil {
		return false, err
	}
	defer dir.Close()
	return true, dir.Sync()
}

func (q *Jobs) Run(ctx context.Context) error {
	if err := os.MkdirAll(q.dir, 0700); err != nil {
		return err
	}
	lock, err := os.OpenFile(filepath.Join(q.dir, ".worker.lock"), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer lock.Close()
	if err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		return fmt.Errorf("队列目录已有消费者: %w", err)
	}
	defer syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	for {
		paths, err := filepath.Glob(filepath.Join(q.dir, "*.json"))
		if err != nil {
			return err
		}
		for _, path := range paths {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			raw, err := os.ReadFile(path)
			var job persistedJob
			if err != nil || json.Unmarshal(raw, &job) != nil {
				log.Printf("后台消息无法读取: %s", filepath.Base(path))
				continue
			}
			q.mu.RLock()
			handler := q.handlers[job.Name]
			q.mu.RUnlock()
			if handler == nil {
				continue
			}
			id := strings.TrimSuffix(strings.SplitN(filepath.Base(path), "-", 2)[1], ".json")
			handler(id, job.Payload)
			if err := os.Remove(path); err != nil {
				return err
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
