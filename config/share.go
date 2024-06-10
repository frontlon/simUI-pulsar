package config

import (
	"encoding/json"
	"simUI/utils"
)

type Share struct {
	RomFile string   //master.zip
	SubGame []string //sub.zip
	Rombase *RomBase
}

// 读取分享配置
func GetShareData(filePath string) (map[string]*Share, error) {

	share := map[string]*Share{}

	if !utils.FileExists(filePath) {
		return share, nil
	}
	content, _ := utils.ReadFile(filePath, false)
	if content != "" {
		json.Unmarshal([]byte(content), &share)
	}
	return share, nil
}

// 覆写分享配置
func WriteDataToShareFile(filePath string, data map[string]*Share) error {
	content, _ := json.Marshal(data)
	utils.OverlayWriteFile(filePath, string(content))
	return nil
}
