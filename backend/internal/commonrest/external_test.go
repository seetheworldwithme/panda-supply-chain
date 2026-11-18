package commonrest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDownloadMissingLocalJavaMessage(t *testing.T) {
	w := httptest.NewRecorder()
	serveLocalDownload(w, httptest.NewRequest("GET", "/", nil), "D:/xmmc.apk", false)
	if w.Code != 200 || w.Body.String() != "本地测试文件[D:/xmmc.apk]不存在" {
		t.Fatalf("%d %s", w.Code, w.Body.String())
	}
}

func TestVideoRangeJavaReturnsTail(t *testing.T) {
	p := filepath.Join(t.TempDir(), "video.mp4")
	os.WriteFile(p, []byte("0123456789"), 0600)
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Range", "bytes=3-5")
	w := httptest.NewRecorder()
	serveVideo(w, r, p, false)
	if w.Code != 206 || w.Body.String() != "3456789" || w.Header().Get("Content-Range") != "bytes 3-9/10" {
		t.Fatalf("%d %s %s", w.Code, w.Body.String(), w.Header().Get("Content-Range"))
	}
}

func TestShareTokensAreUUIDAndExpire(t *testing.T) {
	tok := storeShareToken("7")
	if len(tok) != 36 || !shareTokenExists(tok) || shareTokenExists("missing") {
		t.Fatal(tok)
	}
}

func TestJavaFileDownloadPadsLastBuffer(t *testing.T) {
	got := javaDownloadBytes([]byte("0123456789"))
	if len(got) != 1024 || string(got[:10]) != "0123456789" || got[10] != 0 {
		t.Fatalf("length=%d", len(got))
	}
}

func TestExternalDownloadRequiresJavaGatewayAuth(t *testing.T) {
	mux := http.NewServeMux()
	registerExternal(mux, &Service{})
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest("GET", "/saas-common/download", nil))
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"code":403`) {
		t.Fatalf("%d %s", w.Code, w.Body.String())
	}
}

func TestOCRAdapterSendsImageAndPreservesResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var m map[string]any
		if json.NewDecoder(r.Body).Decode(&m) != nil || m["imageBase64"] != "data:image/jpg;base64,eA==" || m["type"] != float64(2) {
			t.Errorf("unexpected payload: %v", m)
		}
		io.WriteString(w, `{"success":false,"code":701,"msg":"fixture-failure","data":"FailedOperation.OcrFailed"}`)
	}))
	defer server.Close()
	result, e := callOCRAdapter(context.Background(), server.URL, "data:image/jpg;base64,eA==", 2)
	if e != nil || result.Code != 701 || result.Msg != "fixture-failure" || result.Data != "FailedOperation.OcrFailed" {
		t.Fatalf("%+v %v", result, e)
	}
}
