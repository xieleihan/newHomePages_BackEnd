package handler

import (
	"fmt"
	"gin/model"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ALLOWED_STATIC_DIR = "./public/static"
	ALLOWED_EXTENSIONS = ".js,.css,.png,.jpg,.jpeg,.gif,.svg,.ico,.woff,.woff2,.ttf,.eot,.otf,.map"
	MAX_FILE_SIZE      = 10 * 1024 * 1024 // 10MB
)

func isAllowedFile(ext string) bool {
	allowedExts := strings.Split(ALLOWED_EXTENSIONS, ",")
	for _, allow := range allowedExts {
		if ext == allow {
			return true
		}
	}
	return false
}

func sanitizeFileName(name string) string {
	// 移除路径遍历字符和特殊字符
	invalidChars := []string{".", "/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, ch := range invalidChars {
		name = strings.ReplaceAll(name, ch, "")
	}
	return name
}

func formatBytes(bytes int64) string {
	if bytes == 0 {
		return "0 Bytes"
	}
	const k = 1024
	sizes := []string{"Bytes", "KB", "MB", "GB"}
	i := 0
	f := float64(bytes)
	for f >= k && i < len(sizes)-1 {
		f /= k
		i++
	}
	return fmt.Sprintf("%.2f %s", f, sizes[i])
}

// StaticFilesHandler 获取静态文件列表
// @Summary 获取静态文件列表
// @Description 获取public/static目录下的所有允许访问的静态文件信息
// @Tags 静态文件
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "静态文件列表"
// @Failure 404 {object} map[string]interface{} "静态资源不存在"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/static-files [get]
func StaticFilesHandler(c *gin.Context) {
	if _, err := os.Stat(ALLOWED_STATIC_DIR); os.IsNotExist(err) {
		c.JSON(404, gin.H{
			"code":    404,
			"message": "静态资源不存在",
			"time":    time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	var files []model.StaticFileInfo

	err := filepath.WalkDir(ALLOWED_STATIC_DIR, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("无效的路径 %s: %v\n", path, err)
			return nil
		}

		// 跳过目录
		if d.IsDir() {
			return nil
		}

		// 检查扩展名
		ext := strings.ToLower(filepath.Ext(d.Name()))
		if !isAllowedFile(ext) {
			return nil
		}

		// 获取文件信息
		info, err := d.Info()
		if err != nil {
			fmt.Printf("获取文件信息出错 %s: %v\n", path, err)
			return nil
		}

		// 检查文件大小
		if info.Size() > MAX_FILE_SIZE {
			fmt.Printf("文件太大了,跳过: %s (%d bytes)\n", d.Name(), info.Size())
			return nil
		}

		// 构建相对路径
		relPath, _ := filepath.Rel(ALLOWED_STATIC_DIR, path)
		relPath = filepath.ToSlash(relPath)

		files = append(files, model.StaticFileInfo{
			Name:          sanitizeFileName(d.Name()),
			Path:          "/" + relPath,
			Ext:           strings.TrimPrefix(ext, "."),
			Size:          info.Size(),
			SizeFormatted: formatBytes(info.Size()),
			LastModified:  info.ModTime().UTC().Format(time.RFC3339),
		})

		return nil
	})

	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "服务器出错",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":      200,
		"message":   "success",
		"data":      files,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
