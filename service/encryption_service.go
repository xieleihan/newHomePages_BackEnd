package service

import (
	"fmt"
	"gin/db"
	"gin/model"
)

func EncryptMessage(req model.EncryptionMessage) error {
	if err := db.DB.Create(&req).Error; err != nil {
		return fmt.Errorf("保存加密消息失败: %v", err)
	}
	return nil
}

func DecryptMessage(req model.DecryptionMessage) (model.DecryptionResponse, error) {
	var encryptMsg model.EncryptionMessage
	if err := db.DB.Where("uuid = ?", req.UUID).First(&encryptMsg).Error; err != nil {
		return model.DecryptionResponse{}, fmt.Errorf("查找加密消息失败: %v", err)
	}
	return model.DecryptionResponse{
		EncryptedAESKey: encryptMsg.EncryptedAESKey,
		Iv:              encryptMsg.Iv,
		CipherText:      encryptMsg.CipherText,
	}, nil
}
