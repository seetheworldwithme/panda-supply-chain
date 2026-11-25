package platform

import (
	"context"
	"encoding/xml"
	"errors"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// sofficeGate LibreOffice 子进程并发闸门:每个实例几百 MB 内存,不限并发时几十个并发上传就能把 8G 机器 OOM。
var sofficeGate = make(chan struct{}, 2)

// OpenLegacyExcel normalizes a BIFF workbook for the existing OOXML reader.
// The converter runs in a new profile and sees only the bounded temporary input.
func OpenLegacyExcel(input io.Reader) (*os.File, func(), error) {
	const maxWorkbook = 32 << 20
	work, err := os.MkdirTemp("", "panda-legacy-excel-")
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() { _ = os.RemoveAll(work) }
	src := filepath.Join(work, "workbook.xls")
	file, err := os.OpenFile(src, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	size, err := io.Copy(file, io.LimitReader(input, maxWorkbook+1))
	closeErr := file.Close()
	if err != nil || closeErr != nil || size > maxWorkbook {
		cleanup()
		return nil, nil, errors.New("could not stage legacy workbook within input limit")
	}
	bin := os.Getenv("PANDA_SOFFICE_BIN")
	if bin == "" {
		bin = "soffice"
	}
	profile := (&url.URL{Scheme: "file", Path: filepath.Join(work, "profile")}).String()
	sofficeGate <- struct{}{}
	defer func() { <-sofficeGate }()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err = exec.CommandContext(ctx, bin, "-env:UserInstallation="+profile, "--headless", "--convert-to", "xlsx", "--outdir", work, src).Run(); err != nil {
		cleanup()
		return nil, nil, err
	}
	converted, err := os.Open(filepath.Join(work, "workbook.xlsx"))
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	return converted, func() { _ = converted.Close(); cleanup() }, nil
}

// ExcelSharedStrings preserves one index per shared-string item, including
// empty items and cell text formatted as several runs.
func ExcelSharedStrings(data string) []string {
	dec := xml.NewDecoder(strings.NewReader(data))
	var values []string
	var cell strings.Builder
	inItem, inText, inPhonetic := false, false, false
	for {
		token, err := dec.Token()
		if err != nil {
			return values
		}
		switch tag := token.(type) {
		case xml.StartElement:
			switch tag.Name.Local {
			case "si":
				inItem = true
				cell.Reset()
			case "rPh":
				inPhonetic = true
			case "t":
				inText = inItem && !inPhonetic
			}
		case xml.CharData:
			if inText {
				cell.Write(tag)
			}
		case xml.EndElement:
			switch tag.Name.Local {
			case "si":
				values = append(values, cell.String())
				inItem = false
			case "rPh":
				inPhonetic = false
			case "t":
				inText = false
			}
		}
	}
}
