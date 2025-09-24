package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"gin/config"
)

var RDB *redis.Client
var Ctx = context.Background()

func InitRedis(){
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr + ":" + string(rune(config.RedisHost)),
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	if err := RDB.Ping(Ctx).Err(); err != nil {
		log.Fatalf("连接Redis数据库失败: %v", err)
	}

	log.Println("连接Redis数据库成功")
}