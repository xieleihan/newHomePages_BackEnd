package handler

import (
	"fmt"
	"gin/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DownloadPicturesHandler 通用图片下载处理器
// @Summary 解析分享链接并获取图片
// @Description 输入分享文本（支持小红书、抖音等平台），自动判断平台类型，解析并返回图片 URL
// @Tags 图片下载
// @Accept json
// @Produce json
// @Param request body object true "分享文本 {share_text: string}"
// @Success 200 {object} map[string]interface{} "图片 URL 列表"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /api/download-pictures [post]
func DownloadPicturesHandler(c *gin.Context) {
	var req struct {
		ShareText string `json:"share_text" binding:"required"`
		Platform  string `json:"platform"`
	}

	fmt.Println("收到下载图片请求:", req)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"msg":   err.Error(),
			"code":  400,
		})
		return
	}

	// 自动判断平台类型
	platformType := service.DetectPlatform(req.ShareText)

	fmt.Println("检测到的平台类型:", platformType)
	fmt.Println("原始分享文本:", req.ShareText)

	var urls []string
	var err error

	// 根据平台类型调用相应的下载函数
	switch platformType {
	case "xiaohongshu":
		urls, err = service.DownloadXHSPictures(req.ShareText)
	case "douyin":
		urls, err = service.DownloadDouyinPictures(req.ShareText)
	case "bilibili":
		urls, err = service.DownloadBilibiliPictures(req.ShareText)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "不支持的平台",
			"msg":   "无法识别分享链接所属平台",
			"code":  400,
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "图片下载失败",
			"msg":   err.Error(),
			"code":  500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"platform": platformType,
			"pictures": urls,
			"count":    len(urls),
		},
		"code": 200,
		"msg":  "下载成功",
	})
}
