package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DNSQuery 代表 Google DNS API 的请求参数
type DNSQuery struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type"`
}

// DNSResponse 代表 Google DNS API 的响应
type DNSResponse struct {
	Status   int           `json:"Status"`
	TC       bool          `json:"TC"`
	RD       bool          `json:"RD"`
	RA       bool          `json:"RA"`
	AD       bool          `json:"AD"`
	CD       bool          `json:"CD"`
	Question []DNSQuestion `json:"Question"`
	Answer   []DNSAnswer   `json:"Answer"`
	Comment  string        `json:"Comment"`
}

type DNSQuestion struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

type DNSAnswer struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	TTL  int    `json:"TTL"`
	Data string `json:"data"`
}

// QueryDNS 查询 DNS 记录
func QueryDNS(domainName string, recordType string) (DNSResponse, error) {
	// 设置默认记录类型
	if recordType == "" {
		recordType = "A"
	}

	// 构建 Google DNS API URL
	url := fmt.Sprintf("https://dns.google/resolve?name=%s&type=%s", domainName, recordType)

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return DNSResponse{}, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头，模拟浏览器请求
	setDNSRequestHeaders(req)

	// 执行请求
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return DNSResponse{}, fmt.Errorf("DNS查询请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return DNSResponse{}, fmt.Errorf("DNS查询返回错误 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var dnsResp DNSResponse
	if err := json.NewDecoder(resp.Body).Decode(&dnsResp); err != nil {
		return DNSResponse{}, fmt.Errorf("解析DNS响应失败: %v", err)
	}

	return dnsResp, nil
}

// setDNSRequestHeaders 设置请求头
func setDNSRequestHeaders(req *http.Request) {
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("referer", "https://dns.google/query")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"142\", \"Google Chrome\";v=\"142\", \"Not_A Brand\";v=\"99\"")
	req.Header.Set("sec-ch-ua-arch", "\"x86\"")
	req.Header.Set("sec-ch-ua-bitness", "\"64\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-model", "\"\"")
	req.Header.Set("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Set("sec-ch-ua-platform-version", "\"19.0.0\"")
	req.Header.Set("sec-ch-ua-wow64", "?0")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
}
