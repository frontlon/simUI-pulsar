package config

import (
	"encoding/json"
	"errors"
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
		SubGameList[platform] = section
		return section, nil
	}
	content, err := utils.ReadFile(platformInfo.SubGameFile, false)
	if err != nil {
		return section, err
	}

	if content != "" {
		if err := json.Unmarshal([]byte(content), &section); err != nil {
			return section, errors.New("subgame.json parse error: " + err.Error())
		}
	}
	SubGameList[platform] = section
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
	subgame, err := GetSubGame(platform, true)
	if err != nil {
		return err
	}
	subgame[slave] = master

	platformInfo := (&db.Platform{}).GetVOById(platform, false)

	return WriteDataToSubGameFile(platformInfo.SubGameFile, subgame)
}

// 删除子游戏
func DelSubGame(platform uint32, slave string) error {
	subgame, err := GetSubGame(platform, true)
	if err != nil {
		return err
	}

	if _, ok := subgame[slave]; !ok {
		return nil
	}

	delete(subgame, slave)

	platformInfo := (&db.Platform{}).GetVOById(platform, false)

	return WriteDataToSubGameFile(platformInfo.SubGameFile, subgame)
}

// 覆写子游戏配置
func WriteDataToSubGameFile(filePath string, data map[string]string) error {
	content, _ := json.Marshal(data)
	return utils.OverlayWriteFile(filePath, string(content))
}
