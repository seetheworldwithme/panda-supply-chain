package platform

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestRedisUserHashUsesOriginalKeysAndPreservesMiss(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	t.Setenv("PANDA_REDIS_ADDR", listener.Addr().String())
	t.Setenv("PANDA_REDIS_PASSWORD", "fixture")
	t.Setenv("PANDA_REDIS_DB", "2")
	done := make(chan error, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			done <- err
			return
		}
		defer conn.Close()
		reader := bufio.NewReader(conn)
		for _, expected := range []string{"AUTH fixture", "SELECT 2", "HGETALL userId42"} {
			value, err := readRedis(reader)
			if err != nil {
				done <- err
				return
			}
			arr := value.([]any)
			words := []string{}
			for _, v := range arr {
				words = append(words, v.(string))
			}
			if strings.Join(words, " ") != expected {
				done <- fmt.Errorf("unexpected command %v", words)
				return
			}
			if strings.HasPrefix(expected, "HGETALL") {
				fmt.Fprint(conn, "*2\r\n$5\r\nphone\r\n$11\r\n00000000001\r\n")
			} else {
				fmt.Fprint(conn, "+OK\r\n")
			}
		}
		done <- nil
	}()
	hash, err := RedisHash(context.Background(), "userId42")
	if err != nil || hash["phone"] != "00000000001" {
		t.Fatalf("hash=%v error=%v", hash, err)
	}
	if err = <-done; err != nil {
		t.Fatal(err)
	}
	t.Setenv("PANDA_REDIS_ADDR", "")
	hash, err = RedisHash(context.Background(), "userId42")
	if err != nil || len(hash) != 0 {
		t.Fatal("未配置应保留冷缓存")
	}
	if err = RedisSetHash(context.Background(), "userId42", map[string]string{"phone": "fixture"}, time.Minute); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PANDA_REDIS_ADDR", "127.0.0.1:1")
	if _, err = RedisHash(context.Background(), "userId42"); err == nil {
		t.Fatal("配置后连接失败不得伪装成缓存未命中")
	}
}

func TestRedisExpiredSubscriptionKeepsOriginalPattern(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	t.Setenv("PANDA_REDIS_ADDR", listener.Addr().String())
	t.Setenv("PANDA_REDIS_PASSWORD", "")
	t.Setenv("PANDA_REDIS_DB", "")
	done := make(chan error, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			done <- err
			return
		}
		defer conn.Close()
		reader := bufio.NewReader(conn)
		value, err := readRedis(reader)
		if err != nil {
			done <- err
			return
		}
		args := value.([]any)
		if len(args) != 2 || args[0] != "PSUBSCRIBE" || args[1] != "__keyevent@*:expired" {
			done <- fmt.Errorf("错误订阅 %v", args)
			return
		}
		fmt.Fprint(conn, "*3\r\n$10\r\npsubscribe\r\n$20\r\n__keyevent@*:expired\r\n:1\r\n")
		fmt.Fprint(conn, "*4\r\n$8\r\npmessage\r\n$20\r\n__keyevent@*:expired\r\n$22\r\n__keyevent@2__:expired\r\n$8\r\nmember42\r\n")
		done <- nil
		io.Copy(io.Discard, conn)
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = RedisSubscribeExpired(ctx, func(key string) error {
		if key != "member42" {
			t.Errorf("错误过期key %s", key)
		}
		cancel()
		return nil
	})
	if err != context.Canceled {
		t.Fatalf("订阅未随上下文关闭 %v", err)
	}
	if err = <-done; err != nil {
		t.Fatal(err)
	}
}
