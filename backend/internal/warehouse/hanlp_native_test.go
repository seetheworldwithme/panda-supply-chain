package warehouse

import (
	"encoding/base64"
	"os"
	"reflect"
	"strings"
	"testing"
)

// 对照结果由 tools/validation/ExportHanLP.java 调用原服务器同 SHA 的 HanLP JAR 生成。
func TestNativeHanLPMatchesOriginalPortable184(t *testing.T) {
	golden, err := os.ReadFile("hanlp_data/golden.tsv")
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(strings.TrimSuffix(string(golden), "\n"), "\n") {
		parts := strings.Split(line, "\t")
		decode := func(s string) string {
			b, e := base64.StdEncoding.DecodeString(s)
			if e != nil {
				t.Fatal(e)
			}
			return string(b)
		}
		input := decode(parts[0])
		want := []string{}
		for _, s := range parts[1:] {
			want = append(want, decode(s))
		}
		got := hanlpNative(input)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("原Jar分词 %q\nwant %#v\ngot  %#v", input, want, got)
		}
	}
}
