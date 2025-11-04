package main

import (
	"gin/config"
	"gin/db"
	_ "gin/docs"
	"gin/handler"
	"gin/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

// @title           Home Pages Backend API
// @version         1.0
// @description     这是一个个人主页后端 API 服务器
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8082
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	public := gin.Default()
	private := gin.Default()
	db.InitRedis()

	public.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 建议生产配置具体域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	public.Static("/static", "./public/static")
	auth := public.Group("/")
	public.Use(RateLimitMiddleware(5, 10))

	// 公共接口测试
	public.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": config.AppName + " is running!",
			"code":    200,
			"time":    time.Now().Format("2006-01-02 15:04:05"),
			"version": config.Version,
		})
	})
	// 私有接口测试
	// public.GET("/ptest", func(c *gin.Context) {
	// 	resp, err := http.Get("http://127.0.0.1:8083/private/test")
	// 	if err != nil {
	// 		c.JSON(500, gin.H{
	// 			"message": "请求私有接口失败",
	// 			"error":   err.Error(),
	// 		})
	// 		return
	// 	}
	// 	defer resp.Body.Close()

	// 	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	// })

	public.GET("/ip", handler.GetIPInfoHandler)                            // 获取IP的信息路由
	public.GET("/api/static-files", handler.StaticFilesHandler)            // 获取静态资源文件列表路由
	public.POST("/api/send-email", handler.SendEmailHandler)               // 发送邮箱验证码路由
	public.POST("/api/verify-code", handler.VerifyCodeHandler)             // 验证邮箱验证码路由
	public.POST("/api/captcha", handler.GetCaptchaHandler)                 // 获取图形验证码路由
	public.POST("/api/verify-captcha", handler.VerifyCaptchaHandler)       // 验证图形验证码路由
	public.GET("/api/bili-follow-anime", handler.BilibiliAnimeHandler)     // 获取B站追番列表路由
	public.GET("/api/bili-follow-movie", handler.BilibiliMovieHandler)     // 获取B站追剧列表路由
	public.GET("/api/server-status", handler.GetServerStatusHandler)       // 获取服务器运行状态路由
	public.POST("/api/server-status", handler.GetServerStatusHandler)      // 获取服务器运行状态路由(POST)
	public.POST("/api/download-pictures", handler.DownloadPicturesHandler) // 通用图片下载路由

	// Swagger 文档路由
	public.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth.Use(middleware.JWTAuthMiddleware())
	// {
	// 	auth.POST("/api/proxy",handler.ProxyDownloadHandler) // 代理下载路由
	// }

	private.POST("/api/proxy", handler.ProxyDownloadHandler)
	// private.GET("/private/test", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "Private API is running!",
	// 		"code":    200,
	// 		"time":    time.Now().Format("2006-01-02 15:04:05"),
	// 		"version": config.Version,
	// 	})
	// })

	// 启动公共接口服务
	go func() {
		public.Run(config.Port)
	}()
	private.Run("127.0.0.1" + config.PrivatePort)
}
