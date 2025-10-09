package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// axios GET 请求封装
func AxiosGet(api string, params map[string]string, result any) error {
	client := &http.Client{}

	// 构造带参数的 URL
	u, _ := url.Parse(api)
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	// 发请求
	resp, err := client.Get(u.String())
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("JSON 解析失败: %w", err)
	}

	return nil
}
