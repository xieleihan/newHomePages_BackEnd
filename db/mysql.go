package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"gin/config"
)

var DB *gorm.DB

func InitMysql(){
	var err error
	DB, err = gorm.Open(mysql.Open(config.Mysqlusername + ":" + config.Mysqlpassword + "@tcp(" + config.Mysqlhost + ":" + string(rune(config.Mysqlport)) + ")/" + config.Mysqldb + "?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接Mysql数据库失败: %v", err)
	}
	log.Println("连接Mysql数据库成功")
}