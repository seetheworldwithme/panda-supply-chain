package commonrest

import (
	"bytes"
	"context"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestSmallImageCompressionKeepsJavaBytes(t *testing.T) {
	b := []byte("small-image-under-800KiB")
	got, e := compressUploadImage(b)
	if e != nil || !bytes.Equal(b, got) {
		t.Fatalf("%s %v", got, e)
	}
}
func TestExcelUploadProducesSheetPreview(t *testing.T) {
	if _, e := exec.LookPath("soffice"); e != nil {
		t.Skip("soffice not installed")
	}
	if _, e := exec.LookPath("pdftoppm"); e != nil {
		t.Skip("pdftoppm not installed")
	}
	dir := t.TempDir()
	w := httptest.NewRecorder()
	writeXlsx(w, [][]string{{"phase6"}, {"preview"}})
	p := filepath.Join(dir, "fixture.xlsx")
	os.WriteFile(p, w.Body.Bytes(), 0600)
	files, e := renderExcelUpload(context.Background(), p, filepath.Join(dir, "preview"))
	if e != nil || len(files) != 1 {
		t.Fatalf("%v %v", files, e)
	}
	b, e := os.ReadFile(filepath.Join(dir, "preview", files[0]))
	if e != nil || !bytes.HasPrefix(b, []byte("\x89PNG")) {
		t.Fatalf("%v", e)
	}
}
