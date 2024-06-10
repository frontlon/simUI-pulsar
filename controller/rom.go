package controller

import (
	"encoding/json"
	"simUI/db"
	"simUI/modules"
	"simUI/request"
)

// 读取游戏列表
func (a *Controller) GetGameList(req string) string {
	data := &request.GetGameList{}
	_ = json.Unmarshal([]byte(req), &data)
	return Resp(modules.GetGameList(data))
}

// 读取子游戏列表
func (a *Controller) GetSubGameList(id uint64) string {
	return Resp(modules.GetSubGameList(id))
}

// 读取游戏数
func (a *Controller) GetGameCount(req string) string {
	data := &request.GetGameList{}
	_ = json.Unmarshal([]byte(req), &data)
	return Resp(modules.GetGameCount(data))
}

// 根据菜单读取游戏数
func (a *Controller) GetGameCountByMenu(platform uint32) string {
	return Resp(modules.GetGameCountByMenu(platform))
}

// 读取游戏详情信息
func (a *Controller) GetGameDetail(id uint64) string {
	return Resp(modules.GetGameDetail(id))
}

// 读取游戏攻略
func (a *Controller) GetGameStrategy(id uint64) string {
	return Resp(modules.GetGameDoc("strategy", id, 1))
}

// 更新游戏攻略/简介
func (a *Controller) UpdateGameStrategy(id uint64, typ, content string) string {
	return Resp("", modules.SetGameDoc(typ, id, content))
}

// 读取过滤器列表
func (a *Controller) GetFilter(id uint32) string {
	return Resp(modules.GetFilter(id))
}

// 更新游戏通关状态
func (a *Controller) UpdateGameComplete(id uint64, status uint8) string {
	return Resp("", modules.UpdateGameComplete(id, status))
}

// 更新游戏星级
func (a *Controller) UpdateGameScore(id uint64, score float64) string {
	return Resp("", modules.UpdateGameScore(id, score))
}

// 运行游戏
func (a *Controller) RunGame(id uint64, simId uint32) string {
	return Resp("", modules.RunGame(id, simId))
}

// 设置隐藏
func (a *Controller) SetHide(id uint64, hide uint8) string {
	return Resp("", modules.SetHide(id, hide))
}

// 读取没有子游戏的主游戏列表
func (a *Controller) GetGameListNotSubGame(req string) string {
	data := &request.GetGameList{}
	_ = json.Unmarshal([]byte(req), &data)
	return Resp(modules.GetGameListNotSubGame(data))
}

// 直接运行程序
func (a *Controller) RunProgram(pth string) string {
	return Resp("", modules.RunProgram(pth))
}

// 读取游戏文件
func (a *Controller) GetGameStrategyFiles(id uint64) string {
	return Resp(modules.GetStrategyFiles(id))
}

// 更新游戏文件
func (a *Controller) UpdateStrategyFiles(id uint64, req string) string {
	data := []modules.GameStrategyFile{}
	_ = json.Unmarshal([]byte(req), &data)
	return Resp(modules.UpdateStrategyFiles(id, data))
}

// 读取一个ROM的模拟器独立配置
func (a *Controller) GetRomSimSettingById(romId uint64, simId uint32) string {
	return Resp(modules.GetRomSimSettingById(romId, simId))
}

// 更新ROM独立配置模拟器
func (a *Controller) UpdateRomSimSetting(romId uint64, simId uint32, data string) string {
	d := &db.RomSimSetting{}
	json.Unmarshal([]byte(data), &d)
	return Resp("", modules.UpdateRomSimSetting(romId, simId, d))
}
