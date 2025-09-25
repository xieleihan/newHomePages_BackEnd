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

func SendEmailHandler(c *gin.Context) {
	var req EmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code":500, "timestamp": time.Now().Format("2006-01-02 15:04:05")})
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
		"message": "验证码已发送",
		"code":200,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
}


func VerifyCodeHandler(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),"code":500,"timestamp": time.Now().Format("2006-01-02 15:04:05")})
		return
	}
	fmt.Println("验证的邮箱和验证码:", req.Email, req.Code)
	
	if service.VerifyCode(req.Email, req.Code) {
		c.JSON(http.StatusOK, gin.H{"message": "验证码验证成功","code":200,"timestamp": time.Now().Format("2006-01-02 15:04:05")})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码无效或已过期","code": 400,"timestamp": time.Now().Format("2006-01-02 15:04:05")})
	}
}