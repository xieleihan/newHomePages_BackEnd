package main

import (
	"github.com/gin-gonic/gin"
	"time"
	"gin/config"
	"gin/handler"
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

	// 注册IP信息查询路由
	r.GET("/ip",handler.GetIPInfoHandler)

	// 启动服务
	r.Run(config.Port)
}