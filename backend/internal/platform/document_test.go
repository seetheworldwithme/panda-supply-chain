package platform

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestRenderDocumentProducesPDF(t *testing.T) {
	dir := t.TempDir()
	source, target := filepath.Join(dir, "fixture.doc"), filepath.Join(dir, "preview.pdf")
	if err := os.WriteFile(source, []byte(`{\rtf1\ansi Java to Go document fixture}`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := RenderDocumentPDF(context.Background(), source, target); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(target)
	if err != nil || len(data) < 5 || string(data[:5]) != "%PDF-" {
		t.Fatalf("文档未生成 PDF: %v", err)
	}
}
