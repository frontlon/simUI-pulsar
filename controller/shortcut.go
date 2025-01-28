package controller

import (
	"encoding/json"
	"simUI/db"
	"simUI/modules"
)

// 读取快捷工具
func (a *Controller) GetShortcuts(isAbs bool) string {
	return Resp(modules.GetShortcuts(isAbs))
}

// 更新快捷工具
func (a *Controller) UpdateShortcut(data string) string {
	d := []*db.Shortcut{}
	json.Unmarshal([]byte(data), &d)
	return Resp(modules.UpdateShortcut(d))
}

// 更新快捷方式排序
func (a *Controller) UpdateShortcutSort(data string) string {
	d := []uint32{}
	json.Unmarshal([]byte(data), &d)
	return Resp("", modules.UpdateShortcutSort(d))
}
