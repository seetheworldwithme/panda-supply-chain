package platform

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"hash/crc32"
	"image"
	"image/png"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// 几 KB 的图片头声明 60000×60000 即解压炸弹:不得解码,原样返回。
func TestCompressUploadImageSkipsDecompressionBomb(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("\x89PNG\r\n\x1a\n")
	put32 := func(v uint32) { buf.Write([]byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}) }
	chunk := func(kind string, data []byte) {
		put32(uint32(len(data)))
		buf.WriteString(kind)
		buf.Write(data)
		put32(crc32.ChecksumIEEE(append([]byte(kind), data...)))
	}
	ihdr := make([]byte, 13)
	binary.BigEndian.PutUint32(ihdr[0:4], 60000)
	binary.BigEndian.PutUint32(ihdr[4:8], 60000)
	ihdr[8], ihdr[9] = 8, 2
	chunk("IHDR", ihdr)
	chunk("IDAT", []byte{0})
	chunk("IEND", nil)
	bomb := append(buf.Bytes(), make([]byte, 800*1024)...) // 补到 800KiB 以上才会进压缩分支
	if out, err := CompressUploadImage(bomb); err != nil || !bytes.Equal(out, bomb) {
		t.Fatalf("解压炸弹被解码或被改写: err=%v", err)
	}

	// 普通大图(噪声填充,约 1.9MB)仍照常压缩到 800KiB 以下
	img := image.NewNRGBA(image.Rect(0, 0, 700, 700))
	rand.New(rand.NewSource(1)).Read(img.Pix)
	var enc bytes.Buffer
	if err := png.Encode(&enc, img); err != nil {
		t.Fatal(err)
	}
	got, err := CompressUploadImage(enc.Bytes())
	if err != nil || len(got) >= 800*1024 {
		t.Fatalf("普通图片未被压缩: before=%d after=%d err=%v", enc.Len(), len(got), err)
	}
}

// e签宝回调:签名正确但时间戳超出 10 分钟(重放)必须拒绝。
func TestVerifyEsignCallbackRejectsStaleTimestamp(t *testing.T) {
	t.Setenv("appSecret", "fixture-callback-key")
	const body = `{"action":"SIGN_FLOW_FINISH"}`
	signed := func(stamp string) *http.Request {
		r := httptest.NewRequest("POST", "/cb", strings.NewReader(body))
		r.Header.Set("X-Tsign-Open-TIMESTAMP", stamp)
		mac := hmac.New(sha256.New, []byte("fixture-callback-key"))
		mac.Write([]byte(stamp + body))
		r.Header.Set("X-Tsign-Open-SIGNATURE", hex.EncodeToString(mac.Sum(nil)))
		return r
	}
	fresh := strconv.FormatInt(time.Now().UnixMilli(), 10)
	if !VerifyEsignCallback(signed(fresh), []byte(body)) {
		t.Fatal("新鲜回调被误拒")
	}
	stale := strconv.FormatInt(time.Now().Add(-time.Hour).UnixMilli(), 10)
	if VerifyEsignCallback(signed(stale), []byte(body)) {
		t.Fatal("过期回调未被拒绝")
	}
}

// 坏消息(处理器 panic、杂散文件名、无处理器)不能带崩 worker,也不能卡住后面的好消息。
func TestJobsSurvivePoisonMessages(t *testing.T) {
	root := t.TempDir()
	q := NewJobs(root)
	if err := os.WriteFile(filepath.Join(root, "stray.json"), []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}
	expired := filepath.Join(root, strings.Repeat("c", 32)+".seen")
	keep := filepath.Join(root, strings.Repeat("d", 32)+".seen")
	os.WriteFile(expired, []byte("1"), 0600)
	os.WriteFile(keep, []byte(strconv.FormatInt(time.Now().Add(time.Hour).UnixNano(), 10)), 0600)
	for _, name := range []string{"panic", "nobody", "fixture"} {
		if err := q.Submit(name, map[string]any{"id": name}); err != nil {
			t.Fatal(err)
		}
	}
	q.Handle("panic", func(string, json.RawMessage) { panic("fixture") })
	delivered := make(chan struct{}, 1)
	q.Handle("fixture", func(string, json.RawMessage) { delivered <- struct{}{} })
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go q.Run(ctx)
	select {
	case <-delivered:
	case <-time.After(3 * time.Second):
		t.Fatal("坏消息之后的好消息未被消费")
	}
	// worker 空闲后的下一拍隔离坏文件并清理过期的消费标记
	deadline := time.Now().Add(3 * time.Second)
	for {
		if _, err := os.Stat(expired); err != nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("过期消费标记未清理")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if _, err := os.Stat(keep); err != nil {
		t.Fatal("未过期的消费标记被误删")
	}
	if _, err := os.Stat(filepath.Join(root, "stray.json.failed")); err != nil {
		t.Fatal("杂散文件未隔离")
	}
	if left, _ := filepath.Glob(filepath.Join(root, "*.json")); len(left) != 0 {
		t.Fatalf("仍有消息残留: %v", left)
	}
}
