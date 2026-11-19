package platform

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// 微信鉴权和业务通知共用同一传输，调用端保留各自的 Java catch 行为。
func WechatJSON(ctx context.Context, client *http.Client, method, endpoint string, body any) (map[string]any, error) {
	var data []byte
	var err error
	if body != nil {
		data, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("微信接口连接失败")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("微信接口 HTTP %d", response.StatusCode)
	}
	var result map[string]any
	d := json.NewDecoder(response.Body)
	d.UseNumber()
	err = d.Decode(&result)
	return result, err
}

func WechatAccessToken(ctx context.Context, d *sql.DB, client *http.Client, app string) (any, error) {
	rows, err := d.QueryContext(ctx, "SELECT appid,secret FROM chain_auth_chat_authorization WHERE app_identification=?", app)
	if err != nil {
		return nil, err
	}
	configs, err := ScanMaps(rows)
	if err != nil {
		return nil, err
	}
	if len(configs) != 1 {
		return nil, fmt.Errorf("未找到唯一微信应用配置")
	}
	res, err := WechatJSON(ctx, client, "POST", "https://api.weixin.qq.com/cgi-bin/stable_token", map[string]any{"grant_type": "client_credential", "appid": configs[0]["appid"], "secret": configs[0]["secret"]})
	if err != nil {
		return nil, err
	}
	return res["access_token"], nil
}
