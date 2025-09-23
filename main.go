package main

import (
	"github.com/gin-gonic/gin"
	"time"
	"gin/config"
	"gin/handler"
	"github.com/gin-contrib/cors"
	"golang.org/x/time/rate"
)

func RateLimitMiddleware(rps rate.Limit, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rps, burst)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(429, gin.H{
				"error": "太频繁的请求了",
				"code":  429,
				"time":  time.Now().Format("2006-01-02 15:04:05"),
			})
			return
		}
		c.Next()
	}
}

func main()  {
	r := gin.Default() // 创建一个默认的路由引擎

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 建议生产配置具体域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/static", "./public/static")
	api := r.Group("/")
	api.Use(RateLimitMiddleware(5, 10))

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": config.AppName + " is running!",
			"code":200,
			"time":time.Now().Format("2006-01-02 15:04:05"),
			"version": config.Version,
		})
	})

	// 注册IP信息查询路由
	r.GET("/ip",handler.GetIPInfoHandler)// 获取IP的信息路由
	r.GET("/api/static-files", handler.StaticFilesHandler)// 获取静态资源文件列表路由
	r.POST("/api/proxy",handler.ProxyDownloadHandler) // 代理下载路由

	// 启动服务
	r.Run(config.Port)
}