package model

// XHSParseRequest 小红书链接解析请求
type XHSParseRequest struct {
	ShareText string `json:"share_text" binding:"required"` // 小红书分享文本，包含链接
}

// XHSPictureInfo 小红书图片信息
type XHSPictureInfo struct {
	URL      string `json:"url"`       // 图片 URL
	Width    int    `json:"width"`     // 图片宽度
	Height   int    `json:"height"`    // 图片高度
	Format   string `json:"format"`    // 图片格式
	FileSize int64  `json:"file_size"` // 文件大小
}

// XHSParseResponse 小红书链接解析响应
type XHSParseResponse struct {
	NoteID     string           `json:"note_id"`     // 笔记 ID
	Title      string           `json:"title"`       // 笔记标题
	Desc       string           `json:"desc"`        // 笔记描述
	Pictures   []XHSPictureInfo `json:"pictures"`    // 图片列表
	Author     XHSAuthorInfo    `json:"author"`      // 作者信息
	Interact   XHSInteractInfo  `json:"interact"`    // 互动信息
	CreateTime int64            `json:"create_time"` // 创建时间戳
}

// XHSAuthorInfo 小红书作者信息
type XHSAuthorInfo struct {
	UserID   string `json:"user_id"`   // 用户 ID
	NickName string `json:"nick_name"` // 昵称
	Avatar   string `json:"avatar"`    // 头像 URL
}

// XHSInteractInfo 小红书互动信息
type XHSInteractInfo struct {
	LikeCount    int64 `json:"like_count"`    // 点赞数
	CommentCount int64 `json:"comment_count"` // 评论数
	ShareCount   int64 `json:"share_count"`   // 分享数
	CollectCount int64 `json:"collect_count"` // 收藏数
}
