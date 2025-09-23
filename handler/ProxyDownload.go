package handler
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func ProxyDownloadHandler(c *gin.Context) {
	var request struct{
		Url string `json:"url"`
		FileType string `json:"fileType"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	if request.Url == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少url参数",
		})
		return
	}

	if request.FileType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少fileType参数",
		})
		return
	}

	// 解析文件名
	parsedURL, err := url.Parse(request.Url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的URL",
			"error":   err.Error(),
		})
		return
	}
	fileName := filepath.Base(parsedURL.Path)
	if fileName == "" || strings.Contains(fileName, "?") {
		fileName = "downloaded_file"
	}

	saveDir := filepath.Join("public", "static", request.FileType)
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建目录失败",
			"error":   err.Error(),
		})
		return
	}

	filePath := filepath.Join(saveDir, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		resp, err := http.Get(request.Url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "请求远程资源失败",
				"error":   err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": fmt.Sprintf("下载失败，状态码: %d", resp.StatusCode),
			})
			return
		}
		out, err := os.Create(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "保存文件失败",
				"error":   err.Error(),
			})
			return
		}
		defer out.Close()

		if _, err := io.Copy(out, resp.Body); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "写入文件失败",
				"error":   err.Error(),
			})
			return
		}
	}

	fileURL := fmt.Sprintf("/static/%s/%s", request.FileType, fileName)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "下载成功",
		"url":     fileURL,
	})
}