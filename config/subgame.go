package config

import (
	"encoding/json"
	"simUI/db"
	"simUI/utils"
)

// 子游戏文件名 => 主游戏文件名
var SubGameList = map[uint32]map[string]string{}

// 读取子游戏配置
func GetSubGame(platform uint32, build bool) (map[string]string, error) {

	if len(SubGameList) == 0 {
		SubGameList = map[uint32]map[string]string{}
	}

	//如果已经读取，则直接返回
	if len(SubGameList[platform]) != 0 && build == false {
		return SubGameList[platform], nil
	}

	platformInfo := (&db.Platform{}).GetVOById(platform, false)
	section := map[string]string{}

	if !utils.FileExists(platformInfo.SubGameFile) {
		return section, nil
	}
	content, _ := utils.ReadFile(platformInfo.SubGameFile, false)
	if content != "" {
		json.Unmarshal([]byte(content), &section)
		SubGameList[platform] = section
	}
	return section, nil
}

// 根据主游戏，读取所属的子游戏
func GetSubGameByParent(platform uint32, romName string) ([]string, error) {

	//如果已经读取，则直接返回
	datas, _ := GetSubGame(platform, false)

	result := []string{}
	for k, v := range datas {
		if v == romName {
			result = append(result, k)
		}
	}
	return result, nil
}

// 设置子游戏
func SetSubGame(platform uint32, slave, master string) error {
	if len(SubGameList) == 0 {
		SubGameList = map[uint32]map[string]string{}
	}
	if len(SubGameList[platform]) == 0 {
		SubGameList[platform] = map[string]string{}
	}
	SubGameList[platform][slave] = master

	platformInfo := (&db.Platform{}).GetVOById(platform, false)

	return WriteDataToSubGameFile(platformInfo.SubGameFile, SubGameList[platform])
}

// 删除子游戏
func DelSubGame(platform uint32, slave string) error {
	if len(SubGameList) == 0 {
		SubGameList = map[uint32]map[string]string{}
	}
	if len(SubGameList[platform]) == 0 {
		SubGameList[platform] = map[string]string{}
	}

	if _, ok := SubGameList[platform][slave]; !ok {
		return nil
	}

	delete(SubGameList[platform], slave)

	platformInfo := (&db.Platform{}).GetVOById(platform, false)

	return WriteDataToSubGameFile(platformInfo.SubGameFile, SubGameList[platform])
}

// 覆写子游戏配置
func WriteDataToSubGameFile(filePath string, data map[string]string) error {
	content, _ := json.Marshal(data)
	utils.OverlayWriteFile(filePath, string(content))
	return nil
}
