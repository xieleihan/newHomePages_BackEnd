package service

import (
	"errors"
	"fmt"
	"gin/db"
	"gin/model"
	"github.com/google/uuid"
	"time"
)

// RegisterService 用户注册服务
// 验证用户信息并创建新用户账户
// 使用SRP协议，Salt和Verifier由客户端生成
func RegisterService(req model.Register) error {
	// 检查用户名与邮箱是否已存在
	var existingUser model.User
	if err := db.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		// 用户名或邮箱已存在
		fmt.Println("用户名或邮箱已存在 - 用户名:", req.Username, "邮箱:", req.Email)
		return errors.New("用户名或邮箱已存在")
	}

	// 验证邮箱验证码
	if !VerifyCode(req.Email, req.EmailVerificationCode) {
		fmt.Println("邮箱验证码验证失败 - 邮箱:", req.Email)
		return errors.New("邮箱验证码无效或已过期")
	}

	// 验证图片验证码
	if !VerifyCaptcha(req.HumanCheckKey, req.HumanCheckCode) {
		fmt.Println("图片验证码验证失败 - Key:", req.HumanCheckKey)
		return errors.New("图形验证码无效或已过期")
	}

	// 创建新用户 - 只保存必要的字段
	newUser := model.User{
		Username: req.Username,
		Email:    req.Email,
		Salt:     req.Salt,
		Verifier: req.Verifier,
		UserId:   uuid.New().String(),
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	if err := db.DB.Create(&newUser).Error; err != nil {
		fmt.Println("创建用户失败:", err)
		return errors.New("创建用户失败: " + err.Error())
	}

	fmt.Println("用户注册成功 - 用户名:", newUser.Username, "邮箱:", newUser.Email)
	return nil
}
