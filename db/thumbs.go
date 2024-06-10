package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Thumbs struct {
	Platform uint32
	Type     string
	ImageId  string
}

var ThumbMap map[uint32]map[string]uint8

func (*Thumbs) TableName() string {
	return "thumbs"
}

// 写入数据
func (m *Thumbs) Add() error {

	result := getDb().Create(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return result.Error
}

func (m *Thumbs) GetByPlatform(platform uint32) (map[string]uint8, error) {

	if _, ok := ThumbMap[platform]; ok {
		return ThumbMap[platform], nil
	}

	volist := []*Thumbs{}

	result := getDb().Select("image_id").Where("platform = ?", platform).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	platformMap := make(map[string]uint8)
	if len(volist) > 0 {
		for _, v := range volist {
			platformMap[v.ImageId] = 1
		}
	} else {
		//随便塞一个数据，防止穿透
		platformMap["a"] = 1
	}

	ThumbMap[platform] = platformMap
	return platformMap, nil
}

func (m *Thumbs) checkImage(platform uint32, imageId string) bool {

	if _, ok := ThumbMap[platform]; !ok {
		m.GetByPlatform(platform)
	}

	if _, ok := ThumbMap[platform][imageId]; ok {
		return true
	}

	return false
}
