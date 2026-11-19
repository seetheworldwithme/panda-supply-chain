package commonrest

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func (s *Service) enterpriseSearch(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	app, secret := os.Getenv("QCC_APP_KEY"), os.Getenv("QCC_SECRET_KEY")
	empty := func() { platform.Success(w, []any{}) }
	if app == "" || secret == "" {
		empty()
		return
	}
	stamp := strconv.FormatInt(time.Now().Unix(), 10)
	sum := md5.Sum([]byte(app + stamp + secret))
	endpoint := os.Getenv("QCC_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://api.qichacha.com/NameSearch/GetList"
	}
	req, err := http.NewRequestWithContext(r.Context(), "GET", endpoint+"?"+url.Values{"key": {app}, "searchName": {platform.JnullAny(m["etcName"])}}.Encode(), nil)
	if err != nil {
		empty()
		return
	}
	req.Header.Set("Token", strings.ToUpper(hex.EncodeToString(sum[:])))
	req.Header.Set("Timespan", stamp)
	response, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		empty()
		return
	}
	defer response.Body.Close()
	var payload map[string]any
	dec := json.NewDecoder(response.Body)
	dec.UseNumber()
	if dec.Decode(&payload) != nil {
		empty()
		return
	}
	if str(payload, "Status") != "200" {
		platform.Fail(w, 603, "")
		return
	}
	result, ok := payload["Result"].(map[string]any)
	if !ok {
		empty()
		return
	}
	platform.Success(w, result["Data"])
}
