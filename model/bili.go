package model

type BiliFollowResponse struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data BiliFollowData `json:"data"`
}

type BiliFollowData struct {
	Total int `json:"total"`
	Page int `json:"page"`
	PageSize int `json:"pagesize"`
	List []BiliFollowItem `json:"list"`
}

type BiliFollowItem struct {
	Mid int64 `json:"mid"`
	Mtime int64 `json:"mtime"`
	Followed bool `json:"followed"`
	MediaId int64 `json:"media_id"`
	Title string `json:"title"`
	Cover string `json:"cover"`
	Url string `json:"url"`
	EpisodeStatus int `json:"episode_status"`
	Badge string `json:"badge"`
}

