// Package db 数据库连接池。单体直连 MySQL,按请求头 applytype 在 chain/zjjg 间路由
// (对照 Java 版 DynamicDataSource: 默认 chain, applytype=zjjg 切监管库)
package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/config"
)

type DB struct {
	Chain *sql.DB
	Zjjg  *sql.DB
}

// ctxKey 私有 key 类型,避免与其他包的 context 键冲突
type ctxKey int

const dbKey ctxKey = 1

func Open(cfg *config.Config) (*DB, error) {
	chain, err := openPool(cfg.ChainDSN)
	if err != nil {
		return nil, err
	}
	zjjg, err := openPool(cfg.ZjjgDSN)
	if err != nil {
		_ = chain.Close()
		return nil, err
	}
	return &DB{Chain: chain, Zjjg: zjjg}, nil
}

func openPool(dsn string) (*sql.DB, error) {
	d, err := sql.Open("mysql", dsn+"?parseTime=true&loc=Local&charset=utf8mb4")
	if err != nil {
		return nil, err
	}
	d.SetMaxOpenConns(20)
	d.SetMaxIdleConns(5)
	d.SetConnMaxLifetime(30 * time.Minute)
	// 启动即探活,配置错尽快暴露
	if err := d.Ping(); err != nil {
		_ = d.Close()
		return nil, err
	}
	return d, nil
}

func (d *DB) Close() {
	_ = d.Chain.Close()
	_ = d.Zjjg.Close()
}

// Pick 按请求头选库。handler 通过它拿"当前请求的库",
// Java 版 TenantLineInnerInterceptor 的租户条件在这些 SQL 上均为 filter=true(跳过),故无需附加条件
func (d *DB) Pick(ctx context.Context) *sql.DB {
	if v, ok := ctx.Value(dbKey).(*sql.DB); ok && v != nil {
		return v
	}
	return d.Chain
}

// WithSource 把选定的库放入请求 context(middleware 调用)
func WithSource(ctx context.Context, applyType string, d *DB) context.Context {
	src := d.Chain
	if applyType == "zjjg" {
		src = d.Zjjg
	}
	return context.WithValue(ctx, dbKey, src)
}
