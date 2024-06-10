package controller

import "simUI/modules"

// 更新缓存
func (a *Controller) CreateRomCache(platform uint32) string {
	return Resp(modules.CreateRomCache(platform))
}

// 清理游戏统计信息
func (a *Controller) ClearGameStat() string {
	return Resp("", modules.ClearGameStat())
}

// 清理游戏资料
func (a *Controller) ClearNotExistGameConfig() string {
	return Resp("", modules.ClearNotExistGameConfig())
}
