package db

import (
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"simUI/constant"
	"simUI/utils"
)

// 为平台配置填充默认配置
func (m *Platform) SetUiDefault(data PlatformUI) PlatformUI {
	data.Default = m.SetUiDefaultDefault(data.Default)
	data.Playnite = m.SetUiDefaultPlaynite(data.Playnite)
	data.Tiny = m.SetUiDefaultTiny(data.Tiny)

	return data
}

func (m *Platform) SetUiDefaultDefault(data PlatformUIDefault) PlatformUIDefault {

	if data.NameType == 0 {
		data.NameType = 1
	}

	if data.RomListStyle == 0 {
		data.RomListStyle = 1
	}

	if data.BlockThumbType == "" {
		data.BlockThumbType = "thumb"
	}

	if data.RomSort == 0 {
		data.RomSort = 1
	}

	if data.BlockSize == 0 {
		data.BlockSize = 5
	}

	if data.BlockMargin == "" {
		data.BlockMargin = "md"
	}

	if data.BaseFontsize == "" {
		data.BaseFontsize = "14px"
	}

	if len(data.RomListColumn) == 0 {
		data.RomListColumn = constant.DefaultListColumns
	}

	//背景图
	if data.BackgroundImage == "" {
		data.BackgroundImage = constant.RESOURCE_PATH + "background.png"
		if !utils.FileExists(data.BackgroundImage) {
			data.BackgroundImage = ""
		}
	}

	if data.BackgroundImage != "" {
		data.BackgroundImage = utils.WailsPathEncode(data.BackgroundImage, true)
	}

	if data.BackgroundMask != "" {
		data.BackgroundMask = utils.WailsPathEncode(data.BackgroundMask, true)
	}

	return data
}
func (m *Platform) SetUiDefaultPlaynite(data PlatformUIPlaynite) PlatformUIPlaynite {

	if data.NameType == 0 {
		data.NameType = 1
	}

	if data.BlockThumbType == "" {
		data.BlockThumbType = "thumb"
	}

	if data.RomSort == 0 {
		data.RomSort = 1
	}

	if data.BaseFontsize == "" {
		data.BaseFontsize = "14px"
	}

	if data.BlockSize == 0 {
		data.BlockSize = 5
	}

	if data.BlockMargin == "" {
		data.BlockMargin = "xs"
	}

	if data.BaseFontsize == "" {
		data.BaseFontsize = "14px"
	}

	//背景图
	if data.BackgroundImage == "" {
		data.BackgroundImage = constant.RESOURCE_PATH + "background.png"
		if !utils.FileExists(data.BackgroundImage) {
			data.BackgroundImage = ""
		}
	}

	if data.BackgroundImage != "" && !utils.WailsPathCheck(data.BackgroundImage) {
		data.BackgroundImage = utils.WailsPathEncode(data.BackgroundImage, true)
	}

	if data.BackgroundMask != "" {
		data.BackgroundMask = utils.WailsPathEncode(data.BackgroundMask, true)
	}

	return data
}
func (m *Platform) SetUiDefaultTiny(data PlatformUITiny) PlatformUITiny {

	if data.NameType == 0 {
		data.NameType = 1
	}

	if data.RomSort == 0 {
		data.RomSort = 1
	}

	if data.BaseFontsize == "" {
		data.BaseFontsize = "18px"
	}

	//背景图
	if data.BackgroundImage == "" {
		data.BackgroundImage = constant.RESOURCE_PATH + "background.png"
		if !utils.FileExists(data.BackgroundImage) {
			data.BackgroundImage = ""
		}
	}

	if data.BackgroundImage != "" && !utils.WailsPathCheck(data.BackgroundImage) {
		data.BackgroundImage = utils.WailsPathEncode(data.BackgroundImage, true)
	}

	if data.BackgroundMask != "" {
		data.BackgroundMask = utils.WailsPathEncode(data.BackgroundMask, true)
	}

	return data
}

// 更新默认主题UI配置
func (m *Platform) UpdateUiDefaultById(id uint32, data PlatformUIDefault) error {

	platform, _ := m.GetById(id)

	oldUi := &PlatformUI{}
	json.Unmarshal([]byte(platform.Ui), &oldUi)

	oldUi.Default = PlatformUIDefault{
		NameType:            data.NameType,
		RomListStyle:        data.RomListStyle,
		BlockThumbType:      data.BlockThumbType,
		RomSort:             data.RomSort,
		BlockSize:           data.BlockSize,
		BlockMargin:         data.BlockMargin,
		BlockDirection:      data.BlockDirection,
		BaseFontsize:        data.BaseFontsize,
		RomListColumn:       data.RomListColumn,
		BackgroundImage:     data.BackgroundImage,
		BackgroundRepeat:    data.BackgroundRepeat,
		BackgroundFuzzy:     data.BackgroundFuzzy,
		BackgroundMask:      data.BackgroundMask,
		BlockHideTitle:      data.BlockHideTitle,
		BlockHideBackground: data.BlockHideBackground,
	}

	create, _ := json.Marshal(oldUi)

	result := getDb().Table(m.TableName()).Where("id=?", m.Id).Update("ui", string(create))
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}

	m.DelPlatformVOCache() //清空缓存

	return result.Error
}

// 清除默认主题UI
func (m *Platform) ClearAllPlatformUi() error {
	result := getDb().Table(m.TableName()).Update("ui_config", "{}")
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}

	m.DelPlatformVOCache() //清空缓存

	return result.Error
}

// 读取主题配置中的rom排序方式
func (m *Platform) GetPlatformUiRomSort(platform uint32) uint8 {
	conf := (&Config{}).GetConfig(false)
	pfInfo := m.GetVOById(platform, true)
	confSort := uint8(1)
	switch conf.Theme {
	case constant.UI_DEFAULT:
		if platform > 0 && pfInfo != nil {
			confSort = pfInfo.Ui.Default.RomSort
		} else {
			confSort = conf.PlatformUi.Default.RomSort
		}
	}
	return confSort
}
