package service

import (
	"fmt"
	"gin/config"
	"gin/model"
	"gin/utils"
)

func GetFollowList(uid string,followType int)([]model.BiliFollowItem,error){
	page:=1
	var allItems []model.BiliFollowItem

	for {
		fmt.Printf("开始请求第 %d 页\n", page)

		params := map[string]string{
			"vmid":     uid,
			"pn":       fmt.Sprintf("%d", page),
			"ps":       fmt.Sprintf("%d", config.PageSize),
			"type":     fmt.Sprintf("%d", followType),
		}

		var resp model.BiliFollowResponse
		if err:=utils.AxiosGet(config.BiliFollowURL, params, &resp);err!=nil{
			return nil,err
		}

		if resp.Code != 0 {
			return nil, fmt.Errorf("请求失败，错误码: %d, 错误信息: %s", resp.Code, resp.Message)
		}

		list := resp.Data.List
		allItems = append(allItems, list...)

		if len(list) < config.PageSize {
			break
		}
		page++

	}
	return allItems, nil
}

// GetFollowAnime 获取追番列表
func GetFollowAnime(uid string) ([]model.BiliFollowItem, error) {
	return GetFollowList(uid, 1)
}

// GetFollowMovie 获取追剧列表
func GetFollowMovie(uid string) ([]model.BiliFollowItem, error) {
	return GetFollowList(uid, 2)
}