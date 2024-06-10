package controller

import (
	"encoding/json"
	"simUI/db"
	"simUI/modules"
)

// 读取所有模拟器
func (a *Controller) GetAllSimulator() string {
	return Resp(modules.GetAllSimulator())
}

// 读取一个平台下的所有模拟器
func (a *Controller) GetSimulatorByPlatform(id uint32) string {
	return Resp(modules.GetSimulatorByPlatform(id))
}

// 读取一个模拟器
func (a *Controller) GetSimulatorById(id uint32) string {
	return Resp(modules.GetSimulatorById(id))
}

// 添加一个模拟器
func (a *Controller) AddSimulator(platformId uint32, name string) string {
	return Resp(modules.AddSimulator(platformId, name))
}

// 更新模拟器
func (a *Controller) UpdateSimulator(data string) string {
	d := db.Simulator{}
	json.Unmarshal([]byte(data), &d)
	return Resp(modules.UpdateSimulator(d))
}

// 更新模拟器排序
func (a *Controller) UpdateSimulatorSort(data string) string {
	d := []uint32{}
	json.Unmarshal([]byte(data), &d)
	return Resp("", modules.UpdateSimulatorSort(d))
}

// 删除模拟器
func (a *Controller) DelSimulator(id uint32) string {
	return Resp("", modules.DelSimulator(id))
}

// 设置rom模拟器id
func (a *Controller) SetRomSimId(romId uint64, simId uint32) string {
	return Resp("", modules.SetRomSimId(romId, simId))
}

// 批量设置rom模拟器id
func (a *Controller) BatchSetRomSimId(romIds []uint64, simId uint32) string {
	return Resp("", modules.BatchSetRomSimId(romIds, simId))
}
