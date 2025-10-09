package handler

import (
	"net/http"
	"gin/service"
	"github.com/gin-gonic/gin"
)

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