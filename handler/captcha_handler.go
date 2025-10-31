package handler

import (
	"gin/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CaptchaRequest struct {
	UUID string `json:"uuid" binding:"required"`
}

type VerifyCaptchaRequest struct {
	UUID string `json:"uuid" binding:"required"`
	Code string `json:"code" binding:"required,len=6"`
}

// GetCaptchaHandler 获取图形验证码
// @Summary 获取图形验证码
// @Description 生成图形验证码，返回base64编码的图片
// @Tags 验证码
// @Accept json
// @Produce json
// @Param request body CaptchaRequest true "UUID（可选，不传则自动生成）"
// @Success 200 {object} map[string]interface{} "验证码生成成功"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /api/captcha [post]
func GetCaptchaHandler(c *gin.Context) {
	var req CaptchaRequest
	_ = c.ShouldBindJSON(&req)
	if req.UUID == "" {
		req.UUID = uuid.NewString()
	}
	id, b64img, err := service.GenerateCaptcha(req.UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成验证码失败", "code": 500, "timestamp": time.Now().Format("2006-01-02 15:04:05")})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"captcha":   b64img,
		"code":      200,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"messages":  "获取captcha成功",
		"UUID":      id,
	})
}

// VerifyCaptchaHandler 验证图形验证码
// @Summary 验证图形验证码
// @Description 验证用户输入的图形验证码
// @Tags 验证码
// @Accept json
// @Produce json
// @Param request body VerifyCaptchaRequest true "UUID和验证码"
// @Success 200 {object} map[string]interface{} "验证成功"
// @Failure 400 {object} map[string]interface{} "验证失败"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /api/verify-captcha [post]
func VerifyCaptchaHandler(c *gin.Context) {
	var req VerifyCaptchaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": 500, "timestamp": time.Now().Format("2006-01-02 15:04:05")})
		return
	}
	if service.VerifyCaptcha(req.UUID, req.Code) {
		c.JSON(http.StatusOK, gin.H{"message": "验证码验证成功", "code": 200, "timestamp": time.Now().Format("2006-01-02 15:04:05")})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码无效或已过期", "code": 400, "timestamp": time.Now().Format("2006-01-02 15:04:05")})
	}

}
