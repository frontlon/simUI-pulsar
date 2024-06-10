package controller

import (
	"encoding/json"
	"simUI/modules"
	"simUI/request"
)

// 读取全部平台 - VO数据
func (a *Controller) GetPlatform() string {
	return Resp(modules.GetAllPlatform(), nil)
}

// 读取全部平台 - 平台原信息
func (a *Controller) GetPlatformOriginal() string {
	return Resp(modules.GetPlatformOriginal())
}

// 读取一个平台信息
func (a *Controller) GetPlatformById(id uint32) string {
	return Resp(modules.GetPlatformById(id))
}

// 读取平台标签
func (a *Controller) GetAllPlatformTag() string {
	return Resp(modules.GetAllPlatformTag())
}

// 读取平台UI配置
func (a *Controller) GetPlatformUi(id uint32, theme string) string {
	return Resp(modules.GetPlatformUi(id, theme))
}

// 添加一个平台
func (a *Controller) AddPlatform(name string) string {
	return Resp(modules.AddPlatform(name))
}

// 编辑平台信息
func (a *Controller) UpdatePlatform(data string) string {
	d := request.UpdatePlatform{}
	json.Unmarshal([]byte(data), &d)
	return Resp(modules.UpdatePlatform(d))
}

// 编辑平台简介
func (a *Controller) UpdatePlatformDesc(id uint32, desc string) string {
	return Resp("", modules.UpdatePlatformDesc(id, desc))
}

// 更新平台排序
func (a *Controller) UpdatePlatformSort(data string) string {
	d := []uint32{}
	json.Unmarshal([]byte(data), &d)
	return Resp("", modules.UpdatePlatformSort(d))
}

// 删除一个平台
func (a *Controller) DelPlatform(id uint32) string {
	return Resp("", modules.DelPlatform(id))
}

// 更新平台UI设置
func (a *Controller) UpdatePlatformUi(id uint32, theme string, data string) string {
	return Resp("", modules.UpdatePlatformUi(id, theme, data))
}

// 在资源管理器中打开目录
func (a *Controller) OpenFolder(opt string, id uint64) string {
	return Resp("", modules.OpenFolder(opt, id))
}

// 在资源管理器中打开目录 - 直接访问文件
func (a *Controller) OpenFolderByPath(pth string) string {
	return Resp("", modules.OpenFolderByPath(pth))
}
