package handler

import (
	"gin/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

func BilibiliAnimeHandler(c *gin.Context) {
	uid := c.Query("uid")
	if(uid == ""){
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

func BilibiliMovieHandler(c *gin.Context) {
	uid := c.Query("uid")
	if(uid == ""){
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