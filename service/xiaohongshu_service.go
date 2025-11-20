package service

import (
	"encoding/json"
	"fmt"
	"gin/model"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// ParseXHSLink 解析小红书分享链接
func ParseXHSLink(shareText string) (*model.XHSParseResponse, error) {
	// 从文本中提取 xhslink.com 链接
	shortLink := extractXHSLink(shareText)
	fmt.Println("提取到的小红书短链接:", shortLink)
	if shortLink == "" {
		return nil, fmt.Errorf("未找到有效的小红书链接")
	}

	// 获取重定向后的真实链接
	realLink, err := getRedirectLink(shortLink)
	if err != nil {
		return nil, fmt.Errorf("获取重定向链接失败: %v", err)
	}

	// 从真实链接中提取笔记 ID
	noteID := extractNoteID(realLink)
	if noteID == "" {
		return nil, fmt.Errorf("无法提取笔记 ID")
	}

	// 获取笔记详情
	response, err := fetchXHSNoteDetail(noteID)
	if err != nil {
		return nil, fmt.Errorf("获取笔记详情失败: %v", err)
	}

	return response, nil
}

// extractXHSLink 从分享文本中提取小红书链接
func extractXHSLink(text string) string {
	// 匹配 http://xhslink.com 或 https://xhslink.com
	// 使用更宽松的正则表达式，匹配任何非空格的字符到链接结束
	re := regexp.MustCompile(`https?://xhslink\.com/[^\s]+`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 0 {
		link := matches[0]
		// 移除末尾可能的标点符号
		link = strings.TrimRight(link, "，,。.；;")
		return link
	}
	fmt.Println(" 未匹配到任何链接，文本:", text)
	return ""
}

// getRedirectLink 获取重定向后的真实链接
func getRedirectLink(shortLink string) (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 允许跟踪重定向
			return nil
		},
	}

	req, err := http.NewRequest("GET", shortLink, nil)
	if err != nil {
		return "", err
	}

	// 设置 User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println("重定向后的URL:", resp.Request.URL.String())
	// 获取最终 URL（经过所有重定向后）
	return resp.Request.URL.String(), nil
}

// extractNoteID 从 URL 中提取笔记 ID
func extractNoteID(url string) string {
	fmt.Println("输入URL:", url)

	// 小红书 URL 格式1: https://www.xiaohongshu.com/explore/{noteID}
	re := regexp.MustCompile(`/explore/([a-zA-Z0-9_]+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		fmt.Println(" 匹配到格式1，noteID:", matches[1])
		return matches[1]
	}

	// 小红书 URL 格式2: https://www.xiaohongshu.com/discovery/item/{noteID}
	re = regexp.MustCompile(`/discovery/item/([a-zA-Z0-9]+)(?:\?|$)`)
	matches = re.FindStringSubmatch(url)
	fmt.Println(" 格式2匹配结果:", matches)
	if len(matches) > 1 {
		fmt.Println(" 匹配到格式2，noteID:", matches[1])
		return matches[1]
	}

	// 备用格式
	re = regexp.MustCompile(`noteId=([a-zA-Z0-9_]+)`)
	matches = re.FindStringSubmatch(url)
	if len(matches) > 1 {
		fmt.Println(" 匹配到备用格式，noteID:", matches[1])
		return matches[1]
	}

	fmt.Println("无法从URL提取noteID")
	return ""
}

// fetchXHSNoteDetail 获取小红书笔记详情（通过页面 HTML 解析）
func fetchXHSNoteDetail(noteID string) (*model.XHSParseResponse, error) {
	// 构建笔记链接（使用 discovery/item 格式，这是实际的重定向URL格式）
	noteURL := fmt.Sprintf("https://www.xiaohongshu.com/discovery/item/%s", noteID)
	fmt.Println("构建的笔记URL:", noteURL)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", noteURL, nil)
	if err != nil {
		return nil, err
	}

	// 设置必要的请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(" 获取笔记详情失败:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	htmlContent := string(body)

	// 从 HTML 中提取初始数据 JSON
	initialStateData := extractInitialState(htmlContent)
	if initialStateData == nil {
		return nil, fmt.Errorf("无法从页面中提取笔记数据")
	}

	// 构建响应
	response := buildXHSResponse(noteID, initialStateData)
	return response, nil
}

// extractInitialState 从 HTML 中提取初始状态数据
func extractInitialState(html string) map[string]interface{} {
	// 查找 __INITIAL_STATE__ 数据
	re := regexp.MustCompile(`<script>\s*window\.__INITIAL_STATE__\s*=\s*({.*?});</script>`)
	matches := re.FindStringSubmatch(html)

	if len(matches) > 1 {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(matches[1]), &data); err == nil {
			return data
		}
	}

	// 备用方法：查找 feed 数据
	re = regexp.MustCompile(`"feed":\s*({[^}]*})`)
	matches = re.FindStringSubmatch(html)
	if len(matches) > 1 {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(matches[1]), &data); err == nil {
			return data
		}
	}

	return nil
}

// buildXHSResponse 构建小红书响应数据
func buildXHSResponse(noteID string, data map[string]interface{}) *model.XHSParseResponse {
	response := &model.XHSParseResponse{
		NoteID:     noteID,
		Title:      extractTitle(data),
		Desc:       extractDescription(data),
		Pictures:   extractPictures(data),
		Author:     extractAuthorInfo(data),
		Interact:   extractInteractInfo(data),
		CreateTime: time.Now().Unix(),
	}

	return response
}

// extractTitle 提取标题
func extractTitle(data map[string]interface{}) string {
	if data == nil {
		return ""
	}

	// 尝试多种路径提取标题
	if title, ok := data["title"].(string); ok && title != "" {
		return title
	}

	if title, ok := data["interact"].(map[string]interface{})["title"].(string); ok {
		return title
	}

	return "小红书笔记"
}

// extractDescription 提取描述
func extractDescription(data map[string]interface{}) string {
	if data == nil {
		return ""
	}

	if desc, ok := data["desc"].(string); ok {
		return desc
	}

	if desc, ok := data["interact"].(map[string]interface{})["desc"].(string); ok {
		return desc
	}

	return ""
}

// extractPictures 提取图片列表
func extractPictures(data map[string]interface{}) []model.XHSPictureInfo {
	var pictures []model.XHSPictureInfo

	if data == nil {
		return pictures
	}

	// 从各种可能的位置提取图片信息
	if picList, ok := data["images"].([]interface{}); ok {
		for _, pic := range picList {
			if picMap, ok := pic.(map[string]interface{}); ok {
				if url, ok := picMap["url"].(string); ok && url != "" {
					pictures = append(pictures, model.XHSPictureInfo{
						URL: url,
					})
				}
			}
		}
	}

	// 从 image_list 中提取
	if picList, ok := data["image_list"].([]interface{}); ok {
		for _, pic := range picList {
			if picMap, ok := pic.(map[string]interface{}); ok {
				if url, ok := picMap["url"].(string); ok && url != "" {
					pictures = append(pictures, model.XHSPictureInfo{
						URL: url,
					})
				}
			}
		}
	}

	return pictures
}

// extractAuthorInfo 提取作者信息
func extractAuthorInfo(data map[string]interface{}) model.XHSAuthorInfo {
	author := model.XHSAuthorInfo{}

	if data == nil {
		return author
	}

	// 从 author 字段提取
	if authorMap, ok := data["author"].(map[string]interface{}); ok {
		if userID, ok := authorMap["user_id"].(string); ok {
			author.UserID = userID
		}
		if nickName, ok := authorMap["nick_name"].(string); ok {
			author.NickName = nickName
		}
		if avatar, ok := authorMap["avatar"].(string); ok {
			author.Avatar = avatar
		}
	}

	return author
}

// extractInteractInfo 提取互动信息
func extractInteractInfo(data map[string]interface{}) model.XHSInteractInfo {
	interact := model.XHSInteractInfo{}

	if data == nil {
		return interact
	}

	// 从 interact 字段提取
	if interactMap, ok := data["interact"].(map[string]interface{}); ok {
		if likeCount, ok := interactMap["liked"].(float64); ok {
			interact.LikeCount = int64(likeCount)
		}
		if commentCount, ok := interactMap["comment_count"].(float64); ok {
			interact.CommentCount = int64(commentCount)
		}
		if shareCount, ok := interactMap["share_count"].(float64); ok {
			interact.ShareCount = int64(shareCount)
		}
		if collectCount, ok := interactMap["collect_count"].(float64); ok {
			interact.CollectCount = int64(collectCount)
		}
	}

	return interact
}

// DownloadXHSPictures 小红书图片下载函数
func DownloadXHSPictures(shareText string) ([]string, error) {
	// 调用解析函数
	response, err := ParseXHSLink(shareText)
	if err != nil {
		return nil, err
	}

	var urls []string
	for _, pic := range response.Pictures {
		if pic.URL != "" {
			// 确保 URL 带有 scheme
			if !strings.HasPrefix(pic.URL, "http") {
				pic.URL = "https://" + pic.URL
			}
			urls = append(urls, pic.URL)
		}
	}

	return urls, nil
}
