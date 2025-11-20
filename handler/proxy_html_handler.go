package handler

import (
	"fmt"
	"gin/model"
	"gin/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProxyHTMLHandler 代理HTML内容处理器
// @Summary      代理访问网页HTML
// @Description  通过后端代理访问指定URL的网页，获取完整的HTML内容并发送给前端
// @Tags         Proxy
// @Accept       json
// @Produce      json
// @Param        request body model.ProxyHTMLRequest true "代理请求参数"
// @Success      200 {object} model.ProxyHTMLResponse "成功返回HTML内容"
// @Failure      400 {object} model.ProxyHTMLErrorResponse "请求参数错误"
// @Failure      408 {object} model.ProxyHTMLErrorResponse "请求超时"
// @Failure      500 {object} model.ProxyHTMLErrorResponse "服务器错误"
// @Router       /api/proxy-html [post]
func ProxyHTMLHandler(c *gin.Context) {
	var request model.ProxyHTMLRequest

	// 绑定并验证请求参数
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ProxyHTMLErrorResponse{
			Code:    400,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	fmt.Println(" 收到代理HTML请求，URL:", request.URL, "超时时间:", request.Timeout)

	// 创建服务实例
	proxyService := service.NewProxyHTMLService(request.Timeout)

	// 验证URL有效性
	if err := proxyService.ValidateURL(request.URL); err != nil {
		c.JSON(http.StatusBadRequest, model.ProxyHTMLErrorResponse{
			Code:    400,
			Message: "URL验证失败",
			Error:   err.Error(),
		})
		return
	}

	// 获取HTML内容
	response, err := proxyService.FetchHTML(request.URL)
	if err != nil {
		fmt.Println(" 获取HTML失败:", err)
		c.JSON(http.StatusInternalServerError, model.ProxyHTMLErrorResponse{
			Code:    500,
			Message: "获取HTML内容失败",
			Error:   err.Error(),
		})
		return
	}

	fmt.Println(" 返回代理HTML响应，标题:", response.Title)

	c.JSON(http.StatusOK, response)
}
