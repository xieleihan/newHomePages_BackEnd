package service

import (
	"fmt"
	"regexp"
	"strings"
)

// PlatformType 平台类型常量
const (
	PlatformXHS      = "xiaohongshu"
	PlatformDouyin   = "douyin"
	PlatformBilibili = "bilibili"
	PlatformUnknown  = "unknown"
)

// DetectPlatform 检测分享文本所属的平台
func DetectPlatform(shareText string) string {
	text := strings.ToLower(shareText)

	// 检查小红书
	if strings.Contains(text, "xhslink.com") || strings.Contains(text, "小红书") {
		return PlatformXHS
	}

	// 检查抖音
	if strings.Contains(text, "v.douyin.com") || strings.Contains(text, "dy.com") ||
		strings.Contains(text, "抖音") || strings.Contains(text, "douyin") {
		return PlatformDouyin
	}

	// 检查 B 站
	if strings.Contains(text, "bilibili.com") || strings.Contains(text, "b23.tv") ||
		strings.Contains(text, "哔哩哔哩") || strings.Contains(text, "bili") {
		return PlatformBilibili
	}

	return PlatformUnknown
}

// DownloadDouyinPictures 抖音图片下载（预留函数）
func DownloadDouyinPictures(shareText string) ([]string, error) {
	douyinURL := extractDouyinLink(shareText)
	if douyinURL == "" {
		return nil, fmt.Errorf("未找到有效的抖音链接")
	}

	// TODO: 实现抖音链接解析和图片提取逻辑
	// 目前返回空列表作为占位符
	return []string{}, fmt.Errorf("抖音图片下载功能尚未实现")
}

// DownloadBilibiliPictures B站图片下载（预留函数）
func DownloadBilibiliPictures(shareText string) ([]string, error) {
	biliURL := extractBilibiliLink(shareText)
	if biliURL == "" {
		return nil, fmt.Errorf("未找到有效的 B 站链接")
	}

	// TODO: 实现 B 站链接解析和图片提取逻辑
	// 目前返回空列表作为占位符
	return []string{}, fmt.Errorf("B 站图片下载功能尚未实现")
}

// extractXHSLink 从分享文本中提取小红书链接（供平台检测调用）
func extractXHSLinkForPlatform(text string) string {
	// 匹配 http://xhslink.com 或 https://xhslink.com
	// 匹配路径中可能的字符：字母、数字、下划线、连字符、斜杠
	re := regexp.MustCompile(`https?://xhslink\.com/[a-zA-Z0-9_\-/]+`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}

// extractDouyinLink 从分享文本中提取抖音链接
func extractDouyinLink(text string) string {
	// 匹配抖音链接
	re := regexp.MustCompile(`https?://(?:v\.douyin\.com|dy\.com)/[a-zA-Z0-9_\-@]+`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 0 {
		return matches[0]
	}

	// 检查短链接
	re = regexp.MustCompile(`https?://(?:dy\.com|douyin\.com)/[a-zA-Z0-9_\-]+`)
	matches = re.FindStringSubmatch(text)
	if len(matches) > 0 {
		return matches[0]
	}

	return ""
}

// extractBilibiliLink 从分享文本中提取 B 站链接
func extractBilibiliLink(text string) string {
	// 匹配 B 站长链接
	re := regexp.MustCompile(`https?://(?:www\.)?bilibili\.com/video/BV[a-zA-Z0-9]+`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 0 {
		return matches[0]
	}

	// 匹配 B 站短链接
	re = regexp.MustCompile(`https?://b23\.tv/[a-zA-Z0-9]+`)
	matches = re.FindStringSubmatch(text)
	if len(matches) > 0 {
		return matches[0]
	}

	return ""
}
