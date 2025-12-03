package handler

import (
	"gin/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// QueryDNSHandler DNS查询处理器
// @Summary      查询DNS记录
// @Description  通过Google DNS API查询域名的DNS记录
// @Tags         DNS
// @Accept       json
// @Produce      json
// @Param        name  query string true  "域名"
// @Param        type  query string false "记录类型(A/AAAA/MX/TXT等，默认A)"
// @Success      200   {object} service.DNSResponse "DNS查询结果"
// @Failure      400   {object} map[string]interface{} "请求参数错误"
// @Failure      500   {object} map[string]interface{} "服务器错误"
// @Router       /api/dns/query [get]
func QueryDNSHandler(c *gin.Context) {
	// 获取查询参数
	domainName := c.Query("name")
	recordType := c.Query("type")

	// 验证必要参数
	if domainName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"error":     "缺少必要参数",
			"message":   "请提供域名参数 (name)",
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	// 调用服务查询DNS
	dnsResp, err := service.QueryDNS(domainName, recordType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      500,
			"error":     "DNS查询失败",
			"message":   err.Error(),
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":      200,
		"data":      dnsResp,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
}

// QueryDNSPostHandler DNS查询处理器 (POST)
// @Summary      查询DNS记录 (POST)
// @Description  通过Google DNS API查询域名的DNS记录 (POST请求)
// @Tags         DNS
// @Accept       json
// @Produce      json
// @Param        request body service.DNSQuery true "DNS查询请求参数"
// @Success      200     {object} service.DNSResponse "DNS查询结果"
// @Failure      400     {object} map[string]interface{} "请求参数错误"
// @Failure      500     {object} map[string]interface{} "服务器错误"
// @Router       /api/dns/query [post]
func QueryDNSPostHandler(c *gin.Context) {
	var req service.DNSQuery
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":      400,
			"error":     "请求参数错误",
			"message":   err.Error(),
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	// 调用服务查询DNS
	dnsResp, err := service.QueryDNS(req.Name, req.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      500,
			"error":     "DNS查询失败",
			"message":   err.Error(),
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":      200,
		"data":      dnsResp,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
}
