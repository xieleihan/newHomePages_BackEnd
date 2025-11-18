package handler

import (
	"fmt"
	"gin/model"
	"gin/service"
	"net/http"
	"time"
	"strings"
	"github.com/gin-gonic/gin"
)

// RegisterHandler 用户注册处理器
// @Summary      用户注册
// @Description  使用SRP协议进行安全的用户注册，需要验证邮箱验证码和图形验证码
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        request body model.Register true "用户注册请求参数"
// @Success      200 {object} map[string]interface{} "注册成功"
// @Failure      400 {object} map[string]interface{} "参数错误或验证失败"
// @Failure      409 {object} map[string]interface{} "用户已存在"
// @Failure      500 {object} map[string]interface{} "服务器错误"
// @Router       /api/register [post]
func RegisterHandler(c *gin.Context) {
	var req model.Register

	fmt.Println("收到注册请求")

	// 绑定和验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("注册请求参数错误:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     err.Error(),
			"code":      400,
			"message":   "请求参数错误",
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	fmt.Println("注册请求 - 用户名:", req.Username, "邮箱:", req.Email)

	params := map[string]string{
        "[用户名]":       req.Username,
        "[邮箱]":        req.Email,
        "[Salt]":       req.Salt,
        "[Verifier]":   req.Verifier,
        "[邮箱验证码]":    req.EmailVerificationCode,
        "[图形验证码Key]": req.HumanCheckKey,
        "[图形验证码]":    req.HumanCheckCode,
    }

    var missingParams []string
    for name, value := range params {
        if value == "" {
            missingParams = append(missingParams, name)
        }
    }

    if len(missingParams) > 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":     "缺失参数: " + strings.Join(missingParams, ", "),
            "code":      400,
            "message":   "请求参数缺失",
            "timestamp": time.Now().Format("2006-01-02 15:04:05"),
        })
        return
    }

	// 检查验证码是否为空
	if req.EmailVerificationCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "邮箱验证码不能为空",
			"code":      400,
			"message":   "请获取邮箱验证码",
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	if req.HumanCheckKey == "" || req.HumanCheckCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "图形验证码不能为空",
			"code":      400,
			"message":   "请获取图形验证码",
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	// 调用注册服务
	err := service.RegisterService(req)
	if err != nil {
		fmt.Println("注册失败:", err)

		// 根据错误类型返回不同的状态码
		if err.Error() == "用户名或邮箱已存在" {
			c.JSON(http.StatusConflict, gin.H{
				"error":     err.Error(),
				"code":      409,
				"message":   "用户名或邮箱已存在",
				"timestamp": time.Now().Format("2006-01-02 15:04:05"),
			})
			return
		}

		if err.Error() == "邮箱验证码无效或已过期" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     err.Error(),
				"code":      400,
				"message":   "邮箱验证码无效或已过期",
				"timestamp": time.Now().Format("2006-01-02 15:04:05"),
			})
			return
		}

		if err.Error() == "图形验证码无效或已过期" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     err.Error(),
				"code":      400,
				"message":   "图形验证码无效或已过期",
				"timestamp": time.Now().Format("2006-01-02 15:04:05"),
			})
			return
		}

		// 其他错误返回500
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     err.Error(),
			"code":      500,
			"message":   "注册失败",
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	fmt.Println("用户注册成功:", req.Username)
	c.JSON(http.StatusOK, gin.H{
		"message":   "注册成功",
		"code":      200,
		"username":  req.Username,
		"email":     req.Email,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
}
