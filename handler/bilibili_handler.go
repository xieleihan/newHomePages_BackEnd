package handler

import (
	"fmt"
	"gin/service"

	"github.com/gin-gonic/gin"
)

// BilibiliAnimeHandler 获取B站追番列表
// @Summary 获取B站用户追番列表
// @Description 根据用户UID获取B站追番列表
// @Tags B站API
// @Accept json
// @Produce json
// @Param uid query string true "B站用户UID"
// @Success 200 {object} map[string]interface{} "追番列表"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /api/bili-follow-anime [get]
func BilibiliAnimeHandler(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		c.JSON(400, gin.H{"error": "uid不能为空"})
		return
	}
	items, err := service.GetFollowAnime(uid)
	if err != nil {
		fmt.Println("获取追番列表失败:", err)
		c.JSON(500, gin.H{"error": "获取追番列表失败"})
		return
	}
	fmt.Println("追番列表:", items)
	c.JSON(200, gin.H{"data": items})
}

// BilibiliMovieHandler 获取B站追剧列表
// @Summary 获取B站用户追剧列表
// @Description 根据用户UID获取B站追剧列表
// @Tags B站API
// @Accept json
// @Produce json
// @Param uid query string true "B站用户UID"
// @Success 200 {object} map[string]interface{} "追剧列表"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /api/bili-follow-movie [get]
func BilibiliMovieHandler(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		c.JSON(400, gin.H{"error": "uid不能为空"})
		return
	}
	items, err := service.GetFollowMovie(uid)
	if err != nil {
		fmt.Println("获取追剧列表失败:", err)
		c.JSON(500, gin.H{"error": "获取追剧列表失败"})
		return
	}
	fmt.Println("追剧列表:", items)
	c.JSON(200, gin.H{"data": items})
}
