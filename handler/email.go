package handler

import (
	"fmt"
	"gin/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type EmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

// SendEmailHandler 发送邮箱验证码
// @Summary 发送邮箱验证码
// @Description 向指定邮箱发送6位数字验证码
// @Tags 邮箱验证
// @Accept json
// @Produce json
// @Param request body EmailRequest true "邮箱地址"
// @Success 200 {object} map[string]interface{} "验证码发送成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 429 {object} map[string]interface{} "发送太频繁"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /api/send-email [post]
func SendEmailHandler(c *gin.Context) {
	var req EmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": 500, "timestamp": time.Now().Format("2006-01-02 15:04:05")})
		return
	}

	if err := service.SendEmailCode(req.Email); err != nil {
		fmt.Println("发送邮件错误:", err)
		if _, ok := err.(*service.TooFrequentError); ok {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error(), "code": 429, "timestamp": time.Now().Format("2006-01-02 15:04:05")})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": 400, "timestamp": time.Now().Format("2006-01-02 15:04:05")})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "验证码已发送",
		"code":      200,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
}

// VerifyCodeHandler 验证邮箱验证码
// @Summary 验证邮箱验证码
// @Description 验证用户输入的6位数字验证码
// @Tags 邮箱验证
// @Accept json
// @Produce json
// @Param request body VerifyRequest true "邮箱和验证码"
// @Success 200 {object} map[string]interface{} "验证成功"
// @Failure 400 {object} map[string]interface{} "验证失败"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /api/verify-code [post]
func VerifyCodeHandler(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": 500, "timestamp": time.Now().Format("2006-01-02 15:04:05")})
		return
	}
	fmt.Println("验证的邮箱和验证码:", req.Email, req.Code)

	if service.VerifyCode(req.Email, req.Code) {
		c.JSON(http.StatusOK, gin.H{"message": "验证码验证成功", "code": 200, "timestamp": time.Now().Format("2006-01-02 15:04:05")})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码无效或已过期", "code": 400, "timestamp": time.Now().Format("2006-01-02 15:04:05")})
	}
}
