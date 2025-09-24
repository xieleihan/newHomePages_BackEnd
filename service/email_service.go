package service

import (
    "crypto/rand"
    "fmt"
    "gin/config"
    "gin/db"
    "math/big"
	"gopkg.in/gomail.v2"
	"time"
)

func GenerateCode() string{
	cryptos := []rune("0123456789")
	code := make([]rune, 6)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(cryptos))))
		code[i] = cryptos[n.Int64()]
	}
	fmt.Printf("生成的验证码是: %s\n", string(code))
	return string(code)
}

func SendEmail(to, code string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", config.Email)
    m.SetHeader("To", to)
    m.SetHeader("Subject", "账号注册验证码")
    m.SetBody("text/plain", "您的验证码是：" + code + "，请在10分钟内使用。")

    d := gomail.NewDialer(config.EmailHost, config.EmailPort, config.Email, config.EmailPassword)
    
    if config.EmailPort == 465 {
        d.SSL = true
    }

    if err := d.DialAndSend(m); err != nil {
        return err
    }
    return nil
}

/*
存储邮箱验证码
*/
func StoreCode(email, code string) error {
	return db.RDB.Set(db.Ctx, "verify:"+email, code, 180*time.Second).Err()
}

/*
验证邮箱验证码
*/
func VerifyCode(email, code string) bool {
	key := "verify:" + email

    val, err := db.RDB.Get(db.Ctx, key).Result()
	fmt.Print(val)
	if err != nil {
		return false
	}
	if val == code {
        db.RDB.Del(db.Ctx, key)
        return true
    }
	return val == code
}