package service

import (
	"gin/utils"
	"github.com/mojocn/base64Captcha"
	"fmt"
)

/*
生成验证码
*/
func GenerateCaptcha(uuid string) (string, string, error) {
	driver := base64Captcha.NewDriverDigit(80, 240, 6, 0.7, 80) // 高度,宽度,长度,最大扭曲度,背景噪点数
	captcha := base64Captcha.NewCaptcha(driver, utils.RedisStore{})

	id, b64s, _, err := captcha.Generate()

	fmt.Printf("生成的验证码ID: %s, UUID: %s\n", id, uuid)

	// 验证码(不是图片) 6个字符 打印出日志
	fmt.Printf("生成的验证码是: %s\n", utils.RedisStore{}.Get(id, false))

	if err != nil {
		return "", "", err
	}
	return id, b64s, nil
}

/*
验证验证码
*/
func VerifyCaptcha(uuid, code string) bool {
	return utils.RedisStore{}.Verify(uuid, code, true)
}
