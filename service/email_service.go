package service

import (
    "crypto/rand"
    "fmt"
    "gin/config"
    "gin/db"
    "math/big"
	"gopkg.in/gomail.v2"
    "github.com/redis/go-redis/v9"
	"time"
)

type TooFrequentError struct {
    Message string
}

func (e *TooFrequentError) Error() string {
    return e.Message
}

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

func SendEmailCode(email string) error {
	key := "verify:" + email
    rateLimitKey := "rate_limit:" + email
    count, err := db.RDB.Incr(db.Ctx, rateLimitKey).Result()
    if err != nil {
        return err
    }
    if count == 1 {
        if err := db.RDB.Expire(db.Ctx, rateLimitKey, time.Minute).Err(); err != nil {
            return err
        }
    }
    if count > 3 {
        // return fmt.Errorf("请求过于频繁，请稍后再试")
        return &TooFrequentError{Message: "请求过于频繁，请稍后再试"}
    }
	ttl, err := db.RDB.TTL(db.Ctx, key).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if ttl > 0 {
		return fmt.Errorf("验证码仍在有效期内，请 %.0f 秒后再试", ttl.Seconds())
	}

	code := GenerateCode()
	if err := db.RDB.Set(db.Ctx, key, code, 3*time.Minute).Err(); err != nil {
		return err
	}

	if err := SendEmail(email, code); err != nil {
		return err
	}
	return nil
}


func SendEmail(to, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.Email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "账号注册验证码")
	m.SetBody("text/html", `
<div style="width: 400px;height: 50px;display: flex;flex-direction: row ;align-items: center;">
  <img style="width:50px;height:50px;margin-right: 10px;" src="https://github.com/xieleihan/QingluanSearch-AndroidDev/raw/main/peacock_flat.png" alt="" />
  <span style="font-weight: bold;font-family: kaiti;">
    南秋SouthAki
    <span style="font-family: kaiti;color: #ccc;display: block;margin-left: 10px;font-size: 10px;">
      邮箱验证平台
    </span>
  </span>
</div>
<h1>您好：</h1>
<p style="font-size: 18px;color:#000;">
  您的验证码为：
  <span style="font-size: 16px;color:#f00;"><b>` + code + `</b></span>
</p>
<p>您当前正在使用南秋的邮箱验证服务，验证码告知他人将会导致数据信息被盗，请勿泄露!</p>
<p>他人之招,谨防上当受骗.</p>
<p style="font-size: 1.5rem;color:#999;">3分钟内有效</p>
`)


	d := gomail.NewDialer(config.EmailHost, config.EmailPort, config.Email, config.EmailPassword)

	if config.EmailPort == 465 {
		d.SSL = true
	}

	return d.DialAndSend(m)
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