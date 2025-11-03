package handler

import (
	"gin/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetServerStatusHandler 获取服务器运行状态
// @Summary 获取服务器运行状态
// @Description 获取服务器的 CPU、内存、磁盘、网络等实时运行状态，支持 Linux(CentOS/Ubuntu/Debian)、macOS、Windows
// @Tags 服务器状态
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "服务器状态信息"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /api/server-status [get]
// @Router /api/server-status [post]
func GetServerStatusHandler(c *gin.Context) {
	status, err := service.GetServerStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取服务器状态失败",
			"msg":   err.Error(),
			"code":  500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": status,
		"code": 200,
		"msg":  "获取服务器状态成功",
	})
}
