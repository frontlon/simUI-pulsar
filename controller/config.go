package controller

import (
	"encoding/json"
	"simUI/db"
	"simUI/modules"
)

// 读取配置
func (a *Controller) GetConfig() string {
	return Resp(modules.GetConfig())
}

// 读取全部平台 - 平台原信息
func (a *Controller) GetBaseConfig() string {
	return Resp(modules.GetBaseConfig())
}

// 更新一条配置
func (a *Controller) UpdateOneConfig(key, data string) string {
	return Resp("", modules.UpdateOneConfig(key, data))
}

// 更新基本配置
func (a *Controller) UpdateBaseConfig(data string) string {
	d := db.ConfigVO{}
	json.Unmarshal([]byte(data), &d)
	return Resp("", modules.UpdateBaseConfig(d))
}

// 更新颜色配置
func (a *Controller) UpdateColorsConfig(data string) string {
	return Resp("", modules.UpdateOneConfig("Colors", data))
}

// 读取当前主题
func (a *Controller) GetTheme() string {
	return Resp(modules.GetTheme())
}

// 设置当前主题
func (a *Controller) SetTheme(theme string) string {
	return Resp("", modules.SetTheme(theme))
}

// 更新展示图排序
func (a *Controller) UpdateThumbsOrders(orders []string) string {
	return Resp("", modules.UpdateThumbsOrders(orders))
}

// 读取字体列表
func (a *Controller) GetFontList() string {
	return Resp(modules.GetFontList())
}
