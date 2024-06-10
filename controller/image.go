package controller

import (
	"simUI/modules"
)

// 读取游戏图集
func (a *Controller) GetGameThumbs(id uint64, typ string) string {
	return Resp(modules.GetGameThumbs(id, typ))
}

// 添加展示图
func (a *Controller) AddGameThumb(romId uint64, typ string, imageList []string) string {
	return Resp(modules.AddThumb(romId, typ, imageList))
}

// 删除展示图
func (a *Controller) DelGameThumb(romId uint64, typ string, master uint8, imgPath string) string {
	return Resp("", modules.DeleteThumb(romId, typ, master, imgPath))
}

// 图集排序
func (a *Controller) SortGameThumb(romId uint64, typ string, imgList []string) string {
	return Resp(modules.SortGameThumb(romId, typ, imgList))
}

// 读取网络图片
func (a *Controller) LoadWebThumbs(engine string, keyword string, page int) string {
	return Resp(modules.LoadWebThumbs(engine, keyword, page))
}

// 下载网络图片
func (a *Controller) DownloadThumb(romId uint64, typ string, master uint8, httpUrl string, ext string) string {
	return Resp(modules.DownloadThumb(romId, typ, master, httpUrl, ext))
}

// base64转图片
func (a *Controller) CreateRomResByBase64(id uint64, resType string, slaveRes uint8, fileType, base64Str string) string {
	return Resp(modules.CreateRomResByBase64(id, resType, slaveRes, fileType, base64Str))
}
