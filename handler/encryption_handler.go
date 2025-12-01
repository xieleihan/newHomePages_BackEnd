package handler

import (
	"fmt"
	"gin/model"
	"gin/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// EncryptMessageHandler 处理加密消息请求
// @Summary 存入消息
// @Description 前端传递加密信息,后端存储就行
// @Tags 加密消息
// @Accept json
// @Produce json
// @Param request body model.EncryptionMessage true "加密消息"
// @Success 200 {object} map[string]interface{} "消息存储成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /api/SendEncryptionMessage [post]
func EncryptMessageHandler(c *gin.Context) {
	var req model.EncryptionMessage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": 400, "timestamp": time.Now().Format("2006-01-02 15:04:05")})
		return
	}
	if err := service.EncryptMessage(req); err != nil {
		fmt.Println("存储加密消息错误:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "存储加密消息失败", "code": 500, "timestamp": time.Now().Format("2006-01-02 15:04:05")})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   "消息存储成功",
		"code":      200,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
}

// DecryptMessageHandler 处理解密消息请求
// @Summary 取出消息
// @Description 后端直接抛出uuid所对应的加密消息
// @Tags 加密消息
// @Accept json
// @Produce json
// @Param request body model.DecryptionMessage true "解密请求"
// @Success 200 {object} model.EncryptionMessage "解密消息"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 404 {object} map[string]interface{} "消息未找到"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /api/GetEncryptionMessage [post]
func DecryptMessageHandler(c *gin.Context) {
	var req model.DecryptionMessage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": 400, "timestamp": time.Now().Format("2006-01-02 15:04:05")})
		return
	}
	message, err := service.DecryptMessage(req)
	if err != nil {
		fmt.Println("获取解密消息错误:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取解密消息失败", "code": 500, "timestamp": time.Now().Format("2006-01-02 15:04:05")})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   message,
		"code":      200,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
}
