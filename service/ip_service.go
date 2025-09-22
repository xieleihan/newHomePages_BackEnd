package service

import (
	"encoding/json"
	"fmt"
	"gin/config"
	"gin/model"
	"io"
	"net/http"
)

func GetIPInfo(ip string) (*model.IPInfo, error) {
	url := fmt.Sprintf("%s%s", config.IP_API_URL, ip)

	// 构建 request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 打印响应状态
	fmt.Printf("状态码: %s", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("请求结果: %s", string(body))

	var ipInfo model.IPInfo
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		return nil, err
	}

	return &ipInfo, nil
}
