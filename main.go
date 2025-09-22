package main

import (
	"github.com/gin-gonic/gin"
	"time"
	"gin/config"
)

func main()  {
	r := gin.Default() // 创建一个默认的路由引擎

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": config.AppName + " is running!",
			"code":200,
			"time":time.Now().Format("2006-01-02 15:04:05"),
			"version": config.Version,
		})
	})

	// 启动服务
	r.Run(config.Port)
}