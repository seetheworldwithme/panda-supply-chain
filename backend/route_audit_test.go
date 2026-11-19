package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
)

func TestAuditRouteInventory(t *testing.T) {
	mux := newMux(&db.DB{}, &auth.KeyPair{})
	b, e := os.ReadFile("../docs/validation/java-go-20261118/java-source-routes.json")
	if e != nil {
		t.Fatal(e)
	}
	var routes []map[string]any
	if e = json.Unmarshal(b, &routes); e != nil {
		t.Fatal(e)
	}
	for _, v := range routes {
		methods := []string{v["method"].(string)}
		if methods[0] == "ANY" {
			methods = []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "TRACE"}
		}
		path := v["path"].(string)
		parts := strings.Split(path, "/")
		for i, p := range parts {
			if strings.HasPrefix(p, "{") {
				parts[i] = "99977799"
			}
		}
		v["go_pattern"] = ""
		var missing []string
		for _, method := range methods {
			r, _ := http.NewRequest(method, strings.Join(parts, "/"), nil)
			_, pattern := mux.Handler(r)
			if pattern == "" {
				missing = append(missing, method)
			} else {
				v["go_pattern"] = pattern
			}
		}
		if len(missing) > 0 {
			v["go_pattern"] = ""
			v["go_missing_methods"] = missing
		}
	}
	b, _ = json.MarshalIndent(routes, "", "  ")
	if e = os.WriteFile("../docs/validation/java-go-20261118/route-coverage.json", b, 0644); e != nil {
		t.Fatal(e)
	}
}

// 兼容基准必须是服务器运行包；仅生成覆盖报告不能防止遗漏路由。
func TestServerRunningJavaRoutesAreRegistered(t *testing.T) {
	data, err := os.ReadFile("../docs/validation/java-go-20261119/server-java-routes.json")
	if err != nil {
		t.Fatal(err)
	}
	var snapshot struct {
		Routes []struct {
			Method string `json:"method"`
			Path   string `json:"path"`
		} `json:"routes"`
	}
	if err = json.Unmarshal(data, &snapshot); err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Routes) == 0 {
		t.Fatal("服务器路由快照为空")
	}
	mux := newMux(&db.DB{}, &auth.KeyPair{})
	for _, route := range snapshot.Routes {
		methods := []string{route.Method}
		if route.Method == "ANY" {
			methods = []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "TRACE"}
		}
		parts := strings.Split(route.Path, "/")
		for i, part := range parts {
			if strings.HasPrefix(part, "{") {
				parts[i] = "99977799"
			}
		}
		for _, method := range methods {
			r, _ := http.NewRequest(method, strings.Join(parts, "/"), nil)
			if _, pattern := mux.Handler(r); pattern == "" {
				t.Errorf("服务器路由缺失: %s %s", method, route.Path)
			}
		}
	}
}
