package commonrest

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func (s *Service) tlSynchronous(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entChainTlzfOrder.validateBody(m) {
		platform.Unknown(w)
		return
	}
	body, _ := json.Marshal(m)
	r.Body = io.NopCloser(bytes.NewReader(body))
	if err := s.syncTLOrders(r, s.pick(r), true); err != nil {
		platform.Fail(w, 602, "请求失败")
		return
	}
	platform.Success(w, nil)
}

func (s *Service) tlAsync(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		platform.Unknown(w)
		return
	}
	values, exists := r.Form["repMsg"]
	if !exists || len(values) == 0 {
		platform.Unknown(w)
		return
	}
	// Java 异常（含坏签）由 Service 吞掉，Controller 始终返回 success()。
	if err := s.tlAsyncProcess(r, values[0]); err != nil {
		log.Printf("通联回调处理失败: %v", err)
	}
	platform.Success(w, nil)
}

func (s *Service) tlAsyncProcess(r *http.Request, raw string) error {
	if err := platform.VerifyAllinpay(r.Context(), raw); err != nil {
		return err
	}
	response, err := platform.AllinpayResponse(raw)
	if err != nil {
		return err
	}
	if response["resp_code"] != "0000" && response["resp_code"] != "1122" {
		return nil
	}
	order := str(response, "req_trace_num")

	d := s.pick(r)
	var count int
	if err = d.QueryRow("SELECT count(*) FROM chain_tlzf_log WHERE req_type=2 AND res_content LIKE ?", "%"+order+"%").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	if !s.tlCacheAdd(order) {
		return nil
	}
	if _, err = platform.Insert(d, "chain_tlzf_log", map[string]any{"res_content": raw, "req_type": 2}); err != nil {
		return err
	}
	return s.submitTLJob(order, true)
}
