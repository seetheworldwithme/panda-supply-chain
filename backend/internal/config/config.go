// Package config 云管理 Go 单体配置。环境变量优先,缺省值即本地开发形态
// (brew MySQL root 无密码,chain/zjjg 两库与 Java 版 Nacos common-dev.yaml 对齐)
package config

import (
	"os"
	"strings"
)

type Config struct {
	Port     string // HTTP 监听端口,与 Java 网关一致(前端 baseUrl 写死 8765)
	ChainDSN string // 主业务库(默认数据源)
	ZjjgDSN  string // 监管库(请求头 applytype=zjjg 时切换)
	KeysDir  string // RSA 密钥对目录(首次启动自动生成)
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func Load() *Config {
	return &Config{
		Port:     getEnv("PORT", "8765"),
		ChainDSN: getEnv("CHAIN_DSN", "root@tcp(127.0.0.1:3306)/chain"),
		ZjjgDSN:  getEnv("ZJJG_DSN", "root@tcp(127.0.0.1:3306)/zjjg"),
		KeysDir:  getEnv("KEYS_DIR", "keys"),
	}
}

// LoadDotEnv 读取运行目录下 .env(KEY=value 逐行),无文件静默跳过;已存在环境变量优先
func LoadDotEnv() {
	data, err := os.ReadFile(".env")
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "export "))
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.Trim(strings.TrimSpace(v), `"'`)
		if os.Getenv(k) == "" {
			os.Setenv(k, v)
		}
	}
}
