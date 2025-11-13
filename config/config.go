package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	AppName = "Go Gin Web Framework"
	Version = "v1.0.0"
	BiliFollowURL = "https://api.bilibili.com/x/space/bangumi/follow/list"
	PageSize      = 20
)

var (
	Port                string
	PrivatePort         string
	IP_API_URL          string
	SecretKey           string
	RedisAddr           string
	RedisHost           string
	RedisUsername       string
	RedisPassword       string
	RedisDB             int
	TokenExpireDuration time.Duration
	Mysqlhost           string
	Mysqlport           int
	Mysqldb             string
	Mysqlusername       string
	Mysqlpassword       string
	EmailHost           string
	EmailPort           int
	Email               string
	EmailPassword       string
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("加载.env文件失败，使用系统环境变量")
	}

	Port = getEnv("PORT")
	PrivatePort = getEnv("PRIVATE_PORT")
	IP_API_URL = getEnv("IP_API_URL")
	SecretKey = getEnv("SECRET_KEY")
	RedisAddr = getEnv("REDIS_ADDR")
	RedisHost = getEnv("REDIS_HOST")
	RedisUsername = getEnv("REDIS_USERNAME")
	RedisPassword = getEnv("REDIS_PASSWORD")
	RedisDB = getEnvAsInt("REDIS_DB")
	TokenExpireDuration = time.Duration(getEnvAsInt("TOKEN_EXPIRE_DURATION")) * time.Second
	Mysqlhost = getEnv("MYSQLHOST")
	Mysqlport = getEnvAsInt("MYSQLPORT")
	Mysqldb = getEnv("MYSQLDB")
	Mysqlusername = getEnv("MYSQLUSERNAME")
	Mysqlpassword = getEnv("MYSQLPASSWORD")
	EmailHost = getEnv("EMAIL_HOST")
	EmailPort = getEnvAsInt("EMAIL_PORT")
	Email = getEnv("EMAIL")
	EmailPassword = getEnv("EMAIL_PASSWORD")
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return ""
}

func getEnvAsInt(key string) int {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return 0
}
