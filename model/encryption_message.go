package model

type EncryptionMessage struct {
	UUID            string `json:"uuid" gorm:"column:uuid"`              // 消息唯一ID
	EncryptedAESKey string `json:"encryptedAESKey" gorm:"column:encryptedAESKey"` // 加密的AES密钥
	Iv              string `json:"iv" gorm:"column:iv"`                             // IV值
	CipherText      string `json:"cipherText" gorm:"column:cipherText"`            // 密文
}

// TableName 指定表名
func (EncryptionMessage) TableName() string {
	return "encryptionmessage"
}

type DecryptionMessage struct {
	UUID string `json:"uuid"` // 消息唯一ID
}

type DecryptionResponse struct {
	EncryptedAESKey string `json:"encryptedAESKey"` // 加密的AES密钥
	Iv              string `json:"iv"`              // IV值
	CipherText      string `json:"cipherText"`      // 密文
}
