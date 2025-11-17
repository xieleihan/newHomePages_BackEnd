package model

// 使用SRP协议 挑战应答机制

// User 用户数据库模型 - 对应 user 表
type User struct {
	Username string `gorm:"column:username;primaryKey" json:"username"` // 用户名（主键）
	Email    string `gorm:"column:email;uniqueIndex" json:"email"`      // 邮箱地址（唯一索引）
	Salt     string `gorm:"column:salt" json:"salt"`                    // SRP密码盐值
	Verifier string `gorm:"column:verifier" json:"verifier"`            // SRP密码验证器
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}

type Register struct {
	Username              string `json:"username" binding:"required"`                    // 用户名
	Email                 string `json:"email" binding:"required,email"`                 // 邮箱地址
	Salt                  string `json:"salt" binding:"required"`                        // 密码的盐值
	Verifier              string `json:"verifier" binding:"required"`                    // 密码的验证器
	EmailVerificationCode string `json:"emailVerificationCode" binding:"required,len=6"` // 邮箱验证码
	HumanCheckKey         string `json:"humanCheckKey" binding:"required"`               // 人机验证验证码对应的key
	HumanCheckCode        string `json:"humanCheckCode" binding:"required"`              // 人机验证验证码
}

type Login struct {
	Username string `json:"username" binding:"required"` // 用户名
	A        string `json:"A" binding:"required"`        // 客户端公钥
}

type LoginStep2 struct {
	Username string `json:"username" binding:"required"` // 用户名
	M1       string `json:"M1" binding:"required"`       // 客户端证据消息
}

type LoginResponse struct {
	Salt string `json:"salt"` // 密码的盐值
	B    string `json:"B"`    // 服务器公钥
}

type LoginStep2Response struct {
	M2 string `json:"M2"` // 服务器证据消息
}
