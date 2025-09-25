package service

import (
	"gin/utils"
	"github.com/mojocn/base64Captcha"
)

/*
生成验证码
*/
func GenerateCaptcha (uuid string) (string,string,error) {
	driver := base64Captcha.NewDriverDigit(80, 240, 6, 0.7, 80) // 高度,宽度,长度,最大扭曲度,背景噪点数
	captcha := base64Captcha.NewCaptcha(driver, utils.RedisStore{})

	id, b64s, err := captcha.Generate()

	if err != nil {
		return "","", err
	}
	return id,b64s, nil
}

/*
验证验证码
*/
func VerifyCaptcha (uuid, code string) bool {
	return utils.RedisStore{}.Verify(uuid, code, true)
}