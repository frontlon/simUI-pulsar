package controller

import (
	"simUI/modules"
)

// 读取平台下的所有目录
func (a *Controller) GetMenuList(platform uint32) string {
	return Resp(modules.GetMenuList(platform))
}

// 添加目录
func (a *Controller) AddMenu(platform uint32, path, name string) string {
	return Resp("", modules.AddMenu(platform, path, name))
}

// 编辑目录
func (a *Controller) RenameMenu(platform uint32, path, newName string) string {
	return Resp("", modules.RenameMenu(platform, path, newName))
}

// 删除目录
func (a *Controller) DeleteMenu(platform uint32, path string) string {
	return Resp("", modules.DeleteMenu(platform, path))
}

// 排序目录
func (a *Controller) SortMenu(platform uint32, paths []string) string {
	return Resp("", modules.SortMenu(platform, paths))
}
