package db

import (
	"fmt"
	"gin/config"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMysql() {
	var err error
	// 构建DSN（数据源名称）
	// 格式: username:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Mysqlusername,
		config.Mysqlpassword,
		config.Mysqlhost,
		config.Mysqlport,
		config.Mysqldb,
	)

	log.Printf("连接MySQL: %s:%d/%s (用户: %s)", config.Mysqlhost, config.Mysqlport, config.Mysqldb, config.Mysqlusername)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接Mysql数据库失败: %v", err)
	}
	log.Println("连接Mysql数据库成功")
}
