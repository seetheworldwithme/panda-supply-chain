package commonrest

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func (s *Service) wxJSON(r *http.Request, method, endpoint string, body any) (map[string]any, error) {
	return platform.WechatJSON(r.Context(), s.HTTPClient, method, endpoint, body)
}

func (s *Service) chatAccessToken(r *http.Request, app string) (any, error) {
	return platform.WechatAccessToken(r.Context(), s.pick(r), s.HTTPClient, app)
}

func (s *Service) chatGetAccessToken(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	token, err := s.chatAccessToken(r, str(m, "appIdentification"))
	if err != nil {
		platform.Success(w, "")
		return
	} // Java 捕获异常后返回空字符串。
	platform.Success(w, token)
}

func (s *Service) chatOfficialCare(r *http.Request, openid string) (bool, error) {
	token, err := s.chatAccessToken(r, "cloudAccount")
	if err != nil {
		return false, err
	}
	result, err := s.wxJSON(r, "GET", "https://api.weixin.qq.com/cgi-bin/user/info?"+url.Values{"access_token": {fmt.Sprint(token)}, "openid": {openid}, "lang": {"zh_CN"}}.Encode(), nil)
	if err != nil {
		return false, err
	}
	if result["subscribe"] == nil {
		return false, fmt.Errorf("微信关注状态缺失")
	}
	return num(result["subscribe"]) == 1, nil
}
