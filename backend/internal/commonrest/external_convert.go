package commonrest

import (
	"net/http"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// 压缩/渲染实现共享在 platform（authadmin uploadBatch 与这里同源 Java upload）。
func (s *Service) compressFile(r *http.Request, id int64) int64 {
	return platform.CompressFileRow(s.pick(r), id)
}