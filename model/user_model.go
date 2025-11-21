package model

// 使用SRP协议 挑战应答机制

// User 用户数据库模型 - 对应 user 表
type User struct {
	Username  string `gorm:"column:username" json:"username"`           // 用户名
	Email     string `gorm:"column:email" json:"email"`                 // 邮箱地址
	Salt      string `gorm:"column:salt;type:varchar(255)" json:"salt"` // SRP密码盐值
	Verifier  string `gorm:"column:verifier;type:text" json:"verifier"` // SRP密码验证器
	UserId    string `gorm:"column:userId" json:"userId"`               // 用户ID
	CreatedAt string `gorm:"column:createdAt" json:"createdAt"`         // 创建时间
	UpdatedAt string `gorm:"column:updatedAt" json:"updatedAt"`         // 更新时间
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}

type Register struct {
	Username              string `json:"username"`                       // 用户名
	Email                 string `json:"email" binding:"required,email"` // 邮箱地址
	Salt                  string `json:"salt"`                           // 密码的盐值
	Verifier              string `json:"verifier"`                       // 密码的验证器
	EmailVerificationCode string `json:"emailVerificationCode"`          // 邮箱验证码
	HumanCheckKey         string `json:"humanCheckKey"`                  // 人机验证验证码对应的key
	HumanCheckCode        string `json:"humanCheckCode"`                 // 人机验证验证码
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
	M2    string `json:"M2"`    // 服务器证据消息
	Token string `json:"token"` // JWT Token
}

type ChangePassword struct {
	Email            string `json:"email" binding:"required,email"`
	EmailVerificationCode string `json:"emailVerificationCode"`
	HumanCheckKey         string `json:"humanCheckKey"`                  // 人机验证验证码对应的key
	HumanCheckCode        string `json:"humanCheckCode"`
	Salt                  string `json:"salt"`                           // 密码的盐值
	Verifier              string `json:"verifier"` 
}

type ChangeEmail struct {
	OldEmail                 string `json:"oldEmail" binding:"required,email"`
	NewEmail                 string `json:"newEmail" binding:"required,email"`
	EmailVerificationCode    string `json:"emailVerificationCode"`
	NewEmailVerificationCode string `json:"newEmailVerificationCode"`
}