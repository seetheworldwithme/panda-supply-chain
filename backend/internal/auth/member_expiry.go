package auth

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// ScheduleMemberExpiry 对应原 Redis SET member{id} 后 expireAt；过去日期不补发事件。
func ScheduleMemberExpiry(ctx context.Context, userID string, expires time.Time) error {
	key := "member" + userID
	if _, err := platform.RedisDo(ctx, "SET", key, "11"); err != nil {
		return err
	}
	_, err := platform.RedisDo(ctx, "PEXPIREAT", key, strconv.FormatInt(expires.UnixMilli(), 10))
	return err
}

func ScheduleMemberExpiryText(ctx context.Context, userID, expiration string) error {
	expires, err := time.ParseInLocation("2006-01-02 15:04:05.999999999", strings.Replace(expiration, "T", " ", 1), time.Local)
	if err != nil {
		return err
	}
	return ScheduleMemberExpiry(ctx, userID, expires)
}

// StartMemberExpiry 只消费原 Redis 过期事件，不扫描并改写既有会员数据。
// 与原 Spring 订阅一样需要 Redis 已开启 expired 通知；不在此更改服务器配置。
func (s *Service) StartMemberExpiry(ctx context.Context) error {
	if os.Getenv("PANDA_REDIS_ADDR") == "" {
		return nil
	}
	for {
		err := platform.RedisSubscribeExpired(ctx, func(key string) error {
			if err := s.handleMemberExpired(ctx, key, time.Now()); err != nil {
				log.Printf("会员到期事件处理失败: %v", err)
			}
			return nil // 原 Spring 监听器记录单条异常后继续接收后续消息。
		})
		if ctx.Err() != nil {
			return ctx.Err()
		}
		log.Printf("会员到期监听连接中断: %v", err)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
}

func (s *Service) handleMemberExpired(ctx context.Context, key string, now time.Time) error {
	switch {
	case strings.HasPrefix(key, "member"):
		userID := strings.TrimPrefix(key, "member")
		user, err := platform.GetByID(s.DB, "chain_auth_user", userID)
		if err != nil {
			return err
		}
		// 原 cancelMember mapper 使用 ${userId};这里改参数化,Redis 键名里的内容不进 SQL 文本。
		if _, err = s.DB.ExecContext(ctx, "UPDATE chain_auth_user SET member_type=0 WHERE id = ?", userID); err != nil {
			return err
		}
		if user == nil {
			return fmt.Errorf("Java 到期会员为空")
		}
		phone := platform.JnullAny(user["phone"])
		if memberTimeWithin(now, 8*time.Hour, 22*time.Hour) {
			platform.SendLegacySMS(ctx, phone, nil, "2165510", "")
			return nil
		}
		if _, err = platform.RedisDo(ctx, "SET", "delayed-"+phone, "2165510"); err != nil {
			return err
		}
		_, err = platform.RedisDo(ctx, "PEXPIREAT", "delayed-"+phone, strconv.FormatInt(memberReminderTime(now).UnixMilli(), 10))
		return err
	case strings.HasPrefix(key, "delayed"):
		parts := strings.Split(key, "-")
		if len(parts) < 2 {
			return fmt.Errorf("Java 延时短信 key 缺少电话号码")
		}
		phone := parts[1]
		value, err := platform.RedisDo(ctx, "GET", "sendMsg"+phone)
		if err != nil {
			return err
		}
		// 原 Java 的条件确为“标记已存在才发送”，保留冷缓存漏发问题。
		if value != nil {
			platform.SendLegacySMS(ctx, phone, nil, "2165510", "")
		}
		_, err = platform.RedisDo(ctx, "SETEX", "sendMsg"+phone, "3600", "1")
		return err
	}
	return nil
}

func memberTimeWithin(now time.Time, start, end time.Duration) bool {
	clock := time.Duration(now.Hour())*time.Hour + time.Duration(now.Minute())*time.Minute + time.Duration(now.Second())*time.Second + time.Duration(now.Nanosecond())
	return clock >= start && clock <= end
}

func memberReminderTime(now time.Time) time.Time {
	day := now
	if !memberTimeWithin(now, 0, 7*time.Hour+59*time.Minute) {
		day = now.AddDate(0, 0, 1)
	}
	return time.Date(day.Year(), day.Month(), day.Day(), 9, 30, 0, 0, day.Location())
}
