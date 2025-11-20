package service

import (
	"fmt"
	"gin/model"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// ProxyHTMLService 代理HTML服务
type ProxyHTMLService struct {
	client *http.Client
}

// NewProxyHTMLService 创建代理HTML服务实例
func NewProxyHTMLService(timeoutSeconds int) *ProxyHTMLService {
	if timeoutSeconds <= 0 || timeoutSeconds > 300 {
		timeoutSeconds = 30
	}

	return &ProxyHTMLService{
		client: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// 限制重定向次数
				if len(via) >= 10 {
					return fmt.Errorf("太多的重定向")
				}
				return nil
			},
		},
	}
}

// FetchHTML 获取HTML内容
func (s *ProxyHTMLService) FetchHTML(urlStr string) (*model.ProxyHTMLResponse, error) {
	// 验证URL格式
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}

	fmt.Println(" 开始获取HTML，URL:", urlStr)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置User-Agent和其他必要的请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	// 执行请求
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println(" 请求成功，状态码:", resp.StatusCode)

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("HTTP状态码异常: %d", resp.StatusCode)
	}

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %v", err)
	}

	htmlContent := string(bodyBytes)

	// 提取标题
	title := extractHTMLTitle(htmlContent)

	// 构建响应头映射（只保留重要的响应头）
	headers := make(map[string]string)
	importantHeaders := []string{
		"Content-Type",
		"Content-Length",
		"Cache-Control",
		"Last-Modified",
		"ETag",
		"Server",
		"Set-Cookie",
	}

	for _, key := range importantHeaders {
		if value := resp.Header.Get(key); value != "" {
			headers[key] = value
		}
	}

	// 构建响应
	response := &model.ProxyHTMLResponse{
		URL:        resp.Request.URL.String(),
		StatusCode: resp.StatusCode,
		Headers:    headers,
		HTML:       htmlContent,
		Title:      title,
		ContentLen: int64(len(bodyBytes)),
	}

	fmt.Println(" HTML获取成功，长度:", response.ContentLen, "字节")

	return response, nil
}

// extractHTMLTitle 从HTML中提取标题
func extractHTMLTitle(html string) string {
	// 查找 <title> 标签
	re := regexp.MustCompile(`(?i)<title[^>]*>([^<]+)</title>`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		title := strings.TrimSpace(matches[1])
		// 清除特殊字符
		title = strings.ReplaceAll(title, "\n", "")
		title = strings.ReplaceAll(title, "\r", "")
		title = strings.ReplaceAll(title, "\t", "")
		return title
	}

	// 如果没有找到title标签，尝试查找og:title元标签
	re = regexp.MustCompile(`(?i)<meta[^>]*property=["\']og:title["\'][^>]*content=["\']([^"\']+)["\']`)
	matches = re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	return "Unknown Title"
}

// SanitizeHTML 清理HTML内容（移除脚本和危险标签）
func (s *ProxyHTMLService) SanitizeHTML(html string) string {
	// 移除所有script标签
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>[\s\S]*?</script>`)
	html = scriptRegex.ReplaceAllString(html, "")

	// 移除所有style标签
	styleRegex := regexp.MustCompile(`(?i)<style[^>]*>[\s\S]*?</style>`)
	html = styleRegex.ReplaceAllString(html, "")

	// 移除事件监听器
	eventRegex := regexp.MustCompile(`(?i)\s+on\w+\s*=\s*["']([^"']*)["']`)
	html = eventRegex.ReplaceAllString(html, "")

	// 移除javascript协议的链接
	jsLinkRegex := regexp.MustCompile(`(?i)javascript\s*:`)
	html = jsLinkRegex.ReplaceAllString(html, "")

	return html
}

// ValidateURL 验证URL有效性
func (s *ProxyHTMLService) ValidateURL(urlStr string) error {
	if urlStr == "" {
		return fmt.Errorf("URL不能为空")
	}

	if len(urlStr) > 2048 {
		return fmt.Errorf("URL长度超过限制")
	}

	// 检查是否是本地地址或内网地址（安全考虑）
	blockedDomains := []string{
		"localhost",
		"127.0.0.1",
		"0.0.0.0",
		"192.168.",
		"10.0.",
		"172.16.",
		"169.254.",
	}

	for _, domain := range blockedDomains {
		if strings.Contains(strings.ToLower(urlStr), domain) {
			return fmt.Errorf("不允许访问本地或内网地址")
		}
	}

	return nil
}
