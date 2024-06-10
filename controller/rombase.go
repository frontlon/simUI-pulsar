package controller

import (
	"encoding/json"
	"simUI/db"
	"simUI/modules"
)

// 读取全部枚举
func (a *Controller) GetRomBaseEnum() string {
	return Resp(modules.GetRomBaseEnum())
}

// 根据类型读取资料枚举
func (a *Controller) GetRomBaseEnumByType(typ string) string {
	return Resp(modules.GetRomBaseEnumByType(typ))
}

// 更新资料枚举
func (a *Controller) UpdateRomBaseEnumByType(typ string, data string) string {
	d := []string{}
	_ = json.Unmarshal([]byte(data), &d)
	return Resp("", modules.UpdateRomBaseEnum(typ, d))
}

// 读取资料项别名
func (a *Controller) GetRomBaseAlias(platform uint32) string {
	return Resp(modules.GetRomBaseAliasByPlatform(platform))
}

// 更新资料项别名
func (a *Controller) UpdateRomBaseAlias(platform uint32, data string) string {
	d := map[string]string{}
	_ = json.Unmarshal([]byte(data), &d)
	return Resp("", modules.UpdateRomBaseAlias(platform, d))
}

// 读取一条rombase信息
func (a *Controller) GetRomBase(id uint64) string {
	return Resp(modules.GetRomBase(id))
}

// 编辑rombase信息
func (a *Controller) SetRomBase(id uint64, data, name string) string {
	d := map[string]string{}
	_ = json.Unmarshal([]byte(data), &d)
	return Resp(modules.SetRomBase(id, d, name))
}

// 批量rombase信息
func (a *Controller) BatchSetRomBase(platform uint32, data string) string {
	d := map[string]*db.RomSimpleVO{}
	_ = json.Unmarshal([]byte(data), &d)
	return Resp("", modules.BatchSetRomBase(platform, d))
}
