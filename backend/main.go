// 云管理 Go 单体后端(panda-supply-chain)。
// 对照 Java 8 微服务 + 网关的对外形状:同端口 8765、同路径前缀 /saas-xxx、
// 响应体 {success,code,msg,data} 与错误码语义保持一致,前端零改动切换。
package main

import (
	"log"
	"net/http"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/bill"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/authadmin"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/capital"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/config"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/contract"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/dashboard"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/logistics"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/menu"
)

func main() {
	config.LoadDotEnv()
	cfg := config.Load()

	database, err := db.Open(cfg)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer database.Close()

	kp, err := auth.LoadOrGenerateKeyPair(cfg.KeysDir)
	if err != nil {
		log.Fatalf("RSA 密钥初始化失败: %v", err)
	}

	authSvc := auth.NewService(database.Chain, kp)
	menuSvc := menu.NewService(database.Chain)
	dashSvc := dashboard.NewService(database)
	adminSvc := authadmin.NewService(database)
	contractSvc := contract.NewService(database)
	logisticsSvc := logistics.NewService(database)

	mux := http.NewServeMux()

	// ── saas-auth:登录族(permitAll,照 Java SecurityConfig 放行清单) ──
	mux.HandleFunc("POST /saas-auth/login", authSvc.Login)
	mux.HandleFunc("POST /saas-auth/screenLogin", authSvc.ScreenLogin)
	mux.HandleFunc("POST /saas-auth/code/phone", authSvc.PhoneCode)
	mux.HandleFunc("POST /saas-auth/code/image", authSvc.ImageCode)
	mux.HandleFunc("POST /saas-auth/chainAuthUser/updatePwd", authSvc.UpdatePwd)

	// ── saas-auth:菜单/权限 ──
	mux.HandleFunc("POST /saas-auth/chainAuthMenu/selUserMenus", menuSvc.SelUserMenus)
	mux.HandleFunc("POST /saas-auth/chainAuthRole/selMenus", menuSvc.SelMenus)

	// ── saas-auth:组织基础管理(角色/菜单/权限分配,二期) ──
	authadmin.Register(mux, adminSvc)

	// ── saas-common / saas-bill:仪表盘统计(第一期主框架验收面) ──
	dashboard.Register(mux, dashSvc)

	// ── saas-contract 全服务 + saas-common esign DB 层(三期合同) ──
	contract.Register(mux, contractSvc)

	// ── saas-logistics 全服务(四期物流;统计 3 条已由 dashboard 注册) ──
	logistics.Register(mux, logisticsSvc)

	// ── saas-capital 全服务(五期资金;统计 selectQuantity/countSett/countDetail 已由本期精确实现接管) ──
	capitalSvc := capital.NewService(database)
	capital.Register(mux, capitalSvc)

	// ── saas-bill 全服务(五期票据;chainBill 统计 7 条与 chainBillRate 2 条自 dashboard 移交精确实现) ──
	billSvc := bill.NewService(database)
	bill.Register(mux, billSvc)

	// 健康检查
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	handler := cors(auth.Middleware(kp, database.Chain)(mux))
	log.Printf("云管理 Go 后端监听 :%s (chain=%s)", cfg.Port, maskDSN(cfg.ChainDSN))
	log.Fatal(http.ListenAndServe(":"+cfg.Port, handler))
}

// cors 对照 Java 网关的跨域放行:前端各端口直连 8765,允许任意来源与所需请求头
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Access-Control-Allow-Origin", "*")
		h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		h.Set("Access-Control-Allow-Headers", "Authorization, Content-Type, applytype, Signature, AuthorizationId, BatchNumber, routeName, X-Requested-With")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// maskDSN 日志里隐藏 DSN 中的密码段(约定不把凭据写进日志)
func maskDSN(dsn string) string {
	for i := 0; i < len(dsn); i++ {
		if dsn[i] == ':' && i < 12 {
			for j := i; j < len(dsn); j++ {
				if dsn[j] == '@' {
					return dsn[:i+1] + "***" + dsn[j:]
				}
			}
		}
	}
	return dsn
}
