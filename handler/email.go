package handler

import (
	"fmt"
	"gin/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

func SendEmailHandler(c *gin.Context) {
	var req EmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code := service.GenerateCode()
	if err := service.SendEmail(req.Email, code); err != nil {
		fmt.Println("发送邮件错误:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送邮件失败"})
		return
	}

	if err := service.StoreCode(req.Email, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "存储验证码失败"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "验证码已发送"})
}

func VerifyCodeHandler(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("验证的邮箱和验证码:", req.Email, req.Code)
	
	if service.VerifyCode(req.Email, req.Code) {
		c.JSON(http.StatusOK, gin.H{"message": "验证码验证成功"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码无效或已过期"})
	}
}