package controller

import (
	"simUI/modules"
)

// rom文件重命名
func (a *Controller) RenameRomFile(id uint64, name string) string {
	return Resp("", modules.RenameRomFile(id, name))
}

// 删除rom和资源
func (a *Controller) DeleteRom(id uint64, delSubRom uint8, delRes uint8) string {
	return Resp("", modules.DeleteRomAndRes(id, delSubRom, delRes))
}

// 修改游戏别名
func (a *Controller) RenameRomLink(id uint64, name string) string {
	return Resp("", modules.RenameRomLink(id, name))
}

// 复制链接文件
func (a *Controller) CopyRomLink(id uint64, menu string) string {
	return Resp(modules.CopyRomLink(id, menu))
}

// 移动链接文件
func (a *Controller) MoveRomLink(id uint64, menu string) string {
	return Resp("", modules.MoveRomLink(id, menu))
}

// 删除rom链接
func (a *Controller) DeleteRomLink(id uint64) string {
	return Resp("", modules.DeleteRomLink(id))
}

// 绑定子游戏
func (a *Controller) BindSubGame(pid, sid uint64) string {
	return Resp("", modules.BindSubGame(pid, sid))
}

// 解绑子游戏
func (a *Controller) UnBindSubGame(id uint64) string {
	return Resp("", modules.UnBindSubGame(id))
}

// 批量复制链接文件
func (a *Controller) BatchCopyRomLink(ids []uint64, menu string) string {
	return Resp("", modules.BatchCopyRomLink(ids, menu))
}

// 批量移动链接文件
func (a *Controller) BatchMoveRomLink(ids []uint64, menu string) string {
	return Resp("", modules.BatchMoveRomLink(ids, menu))
}

// 批量删除rom链接
func (a *Controller) BatchDeleteRomLink(ids []uint64) string {
	return Resp("", modules.BatchDeleteRomLink(ids))
}

// 批量删除rom和资源
func (a *Controller) BatchDeleteRom(ids []uint64, delSubRom, delRes uint8) string {
	return Resp("", modules.BatchDeleteRomAndRes(ids, delSubRom, delRes))
}

// 查找无效资源
func (a *Controller) CheckUnownedRes(platform uint32) string {
	return Resp(modules.CheckUnownedRes(platform))
}

// 删除无效资源文件
func (a *Controller) DeleteUnownedFile(platform uint32, ids []string) string {
	return Resp("", modules.DeleteUnownedFile(platform, ids))
}

// 打开备份文件夹
func (a *Controller) OpenCacheFolder(typ string, create int) string {
	return Resp("", modules.OpenCacheFolder(typ, create))
}

// 检查重复rom
func (a *Controller) CheckRomRepeat(platform uint32) string {
	return Resp(modules.CheckRomRepeat(platform))
}

// 删除重复ROM文件
func (a *Controller) DeleteRepeatFile(ids []string) string {
	return Resp("", modules.DeleteRepeatFile(ids))
}
