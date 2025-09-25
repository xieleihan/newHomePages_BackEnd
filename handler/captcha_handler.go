package handler

import (
	"gin/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"github.com/google/uuid"
)

type CaptchaRequest struct {
	UUID string `json:"uuid" binding:"required"`
}

type VerifyCaptchaRequest struct {
	UUID string `json:"uuid" binding:"required"`
	Code string `json:"code" binding:"required,len=6"`
}

/*
获取Captcha
*/
func GetCaptchaHandler(c *gin.Context) {
	var req CaptchaRequest
	_ = c.ShouldBindJSON(&req)
	if req.UUID == "" {
        req.UUID = uuid.NewString()
    }
	id,b64img, err := service.GenerateCaptcha(req.UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成验证码失败","code": 500,"timestamp": time.Now().Format("2006-01-02 15:04:05")})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"captcha": b64img,
		"code":200,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"messages": "获取captcha成功",
		"UUID": id,
	})
}

/*
验证Captcha
*/
func VerifyCaptchaHandler(c *gin.Context) {
	var req VerifyCaptchaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),"code":500,"timestamp": time.Now().Format("2006-01-02 15:04:05")})
		return
	}
	if service.VerifyCaptcha(req.UUID, req.Code) {
		c.JSON(http.StatusOK, gin.H{"message": "验证码验证成功","code":200,"timestamp": time.Now().Format("2006-01-02 15:04:05")})
	}else{
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码无效或已过期","code": 400,"timestamp": time.Now().Format("2006-01-02 15:04:05")})
	}

}