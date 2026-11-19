package platform

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 兼容旧 Java 的 Redis hash；不查数据库补缺失值，保留冷缓存行为。
func RedisHash(ctx context.Context, key string) (map[string]string, error) {
	value, err := redisCommand(ctx, "HGETALL", key)
	if err != nil {
		return nil, err
	}
	out := map[string]string{}
	if value == nil {
		return out, nil
	}
	items, ok := value.([]any)
	if !ok || len(items)%2 != 0 {
		return nil, fmt.Errorf("Redis hash 响应无效")
	}
	for i := 0; i < len(items); i += 2 {
		k, kok := items[i].(string)
		v, vok := items[i+1].(string)
		if !kok || !vok {
			return nil, fmt.Errorf("Redis hash 字段无效")
		}
		out[k] = v
	}
	return out, nil
}

func RedisSetHash(ctx context.Context, key string, fields map[string]string, ttl time.Duration) error {
	if len(fields) == 0 {
		return nil
	}
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	command := []string{"HSET", key}
	for _, k := range keys {
		command = append(command, k, fields[k])
	}
	if _, err := redisCommand(ctx, command...); err != nil {
		return err
	}
	if ttl > 0 {
		_, err := redisCommand(ctx, "EXPIRE", key, strconv.FormatInt(int64(ttl/time.Second), 10))
		return err
	}
	return nil
}

func redisCommand(ctx context.Context, args ...string) (any, error) {
	return RedisDo(ctx, args...)
}

func RedisDo(ctx context.Context, args ...string) (any, error) {
	conn, reader, err := redisConnect(ctx)
	if err != nil || conn == nil {
		return nil, err
	}
	defer conn.Close()
	return redisRun(conn, reader, args...)
}

func RedisSubscribeExpired(ctx context.Context, callback func(string) error) error {
	conn, reader, err := redisConnect(ctx)
	if err != nil || conn == nil {
		return err
	}
	defer conn.Close()
	stop := context.AfterFunc(ctx, func() { conn.Close() })
	defer stop()
	if _, err = redisRun(conn, reader, "PSUBSCRIBE", "__keyevent@*:expired"); err != nil {
		return err
	}
	if err = conn.SetDeadline(time.Time{}); err != nil {
		return err
	}
	for {
		value, err := readRedis(reader)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return err
		}
		items, ok := value.([]any)
		if !ok || len(items) != 4 || items[0] != "pmessage" {
			continue
		}
		key, ok := items[3].(string)
		if !ok {
			continue
		}
		if err = callback(key); err != nil {
			return err
		}
	}
}

func redisRun(conn net.Conn, reader *bufio.Reader, args ...string) (any, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "*%d\r\n", len(args))
	for _, v := range args {
		fmt.Fprintf(&b, "$%d\r\n%s\r\n", len(v), v)
	}
	if _, err := io.WriteString(conn, b.String()); err != nil {
		return nil, err
	}
	return readRedis(reader)
}

func redisConnect(ctx context.Context) (net.Conn, *bufio.Reader, error) {
	address := os.Getenv("PANDA_REDIS_ADDR")
	if address == "" {
		return nil, nil, nil
	}
	conn, err := (&net.Dialer{Timeout: 5 * time.Second}).DialContext(ctx, "tcp", address)
	if err != nil {
		return nil, nil, err
	}
	ok := false
	defer func() {
		if !ok {
			conn.Close()
		}
	}()
	deadline := time.Now().Add(5 * time.Second)
	if d, ok := ctx.Deadline(); ok && d.Before(deadline) {
		deadline = d
	}
	conn.SetDeadline(deadline)
	reader := bufio.NewReader(conn)
	if password := os.Getenv("PANDA_REDIS_PASSWORD"); password != "" {
		if _, err = redisRun(conn, reader, "AUTH", password); err != nil {
			return nil, nil, err
		}
	}
	if db := os.Getenv("PANDA_REDIS_DB"); db != "" && db != "0" {
		if _, err = strconv.Atoi(db); err != nil {
			return nil, nil, fmt.Errorf("Redis DB 配置无效")
		}
		if _, err = redisRun(conn, reader, "SELECT", db); err != nil {
			return nil, nil, err
		}
	}
	ok = true
	return conn, reader, nil
}

func readRedis(r *bufio.Reader) (any, error) {
	prefix, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(line, "\r\n") {
		return nil, fmt.Errorf("Redis 行结束无效")
	}
	line = strings.TrimSuffix(line, "\r\n")
	switch prefix {
	case '+':
		return line, nil
	case '-':
		return nil, fmt.Errorf("Redis 命令失败") // 不将含认证内容的服务端错误写入日志。
	case ':':
		return strconv.ParseInt(line, 10, 64)
	case '$':
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		if n == -1 {
			return nil, nil
		}
		if n < 0 || n > 16<<20 {
			return nil, fmt.Errorf("Redis bulk 长度无效")
		}
		b := make([]byte, n+2)
		if _, err = io.ReadFull(r, b); err != nil {
			return nil, err
		}
		if string(b[n:]) != "\r\n" {
			return nil, fmt.Errorf("Redis bulk 结束无效")
		}
		return string(b[:n]), nil
	case '*':
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		if n == -1 {
			return nil, nil
		}
		if n < 0 || n > 100000 {
			return nil, fmt.Errorf("Redis array 长度无效")
		}
		out := make([]any, n)
		for i := range out {
			out[i], err = readRedis(r)
			if err != nil {
				return nil, err
			}
		}
		return out, nil
	default:
		return nil, fmt.Errorf("Redis 响应类型无效")
	}
}
