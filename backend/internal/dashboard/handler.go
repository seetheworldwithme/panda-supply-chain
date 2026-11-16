// Package dashboard 仪表盘统计接口(企业 Main/平台 PtIndex/政府大屏 Xzmain)。
// 接口对照 Java chainBasicMainPage/chainBill 的 count 系列,SQL 逐条翻译(见各 handler 注释)
package dashboard

import (
	"net/http"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
)

type Service struct {
	DB *db.DB
}

func NewService(database *db.DB) *Service { return &Service{DB: database} }

// Register 挂载仪表盘路由
func Register(mux *http.ServeMux, s *Service) {
	_ = mux // 统计接口在探索报告对齐后逐条挂载
}
