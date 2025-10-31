package handler

import (
	"gin/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetIPInfoHandler 获取IP信息
// @Summary 获取IP地址信息
// @Description 获取指定IP或客户端IP的地理位置信息
// @Tags IP信息
// @Accept json
// @Produce json
// @Param query query string false "要查询的IP地址，不填则查询客户端IP"
// @Success 200 {object} map[string]interface{} "IP信息"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /ip [get]
func GetIPInfoHandler(c *gin.Context) {
	ip := c.Query("query")

	if ip == "" {
		ip = c.ClientIP()
	}

	ipInfo, err := service.GetIPInfo(ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ipInfo)
}
