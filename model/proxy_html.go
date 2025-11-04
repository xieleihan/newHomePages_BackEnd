package model

// ProxyHTMLRequest 代理HTML请求结构体
type ProxyHTMLRequest struct {
	URL     string `json:"url" binding:"required,url"`
	Timeout int    `json:"timeout"` // 超时时间(秒)，默认30秒
}

// ProxyHTMLResponse 代理HTML响应结构体
type ProxyHTMLResponse struct {
	URL        string            `json:"url"`
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	HTML       string            `json:"html"`
	Title      string            `json:"title"`
	ContentLen int64             `json:"contentLen"`
}

// ProxyHTMLErrorResponse 代理HTML错误响应结构体
type ProxyHTMLErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
