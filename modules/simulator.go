package modules

import (
	"simUI/config"
	"simUI/db"
	"simUI/utils"
)

// 读取所有模拟器
func GetAllSimulator() (map[uint32][]*db.Simulator, error) {
	dolist, _ := (&db.Simulator{}).GetAll()
	data := map[uint32][]*db.Simulator{}

	for _, v := range dolist {
		if _, ok := data[v.Platform]; ok {
			data[v.Platform] = append(data[v.Platform], v)
		} else {
			data[v.Platform] = []*db.Simulator{v}
		}
	}
	return data, nil
}

// 读取一个平台下的所有模拟器
func GetSimulatorByPlatform(id uint32) ([]*db.Simulator, error) {
	return (&db.Simulator{}).GetByPlatform(id)
}

// 读取一个模拟器
func GetSimulatorById(id uint32) (*db.Simulator, error) {
	return (&db.Simulator{}).GetById(id)
}

// 添加模拟器
func AddSimulator(platformId uint32, name string) (*db.Simulator, error) {

	sim := &db.Simulator{
		Name:     name,
		Platform: platformId,
	}
	id, err := sim.Add()
	if err != nil {
		return nil, err
	}
	sim.Id = id
	return sim, nil
}

// 更新模拟器
func UpdateSimulator(data db.Simulator) (*db.Simulator, error) {
	sim := &db.Simulator{
		Id:        data.Id,
		Name:      data.Name,
		Platform:  data.Platform,
		Path:      data.Path,
		Cmd:       data.Cmd,
		RunBefore: data.RunBefore,
		RunAfter:  data.RunAfter,
	}

	//更新模拟器
	if err := sim.UpdateById(); err != nil {
		return sim, err
	}

	return sim, nil
}

// 更新模拟器排序
func UpdateSimulatorSort(lists []uint32) error {
	if len(lists) == 0 {
		return nil
	}
	for k, pfId := range lists {
		simulator := &db.Simulator{
			Id:   pfId,
			Sort: uint32(k + 1),
		}
		err := simulator.UpdateSortById()
		if err != nil {
			utils.WriteLog(err.Error())
			return err
		}
	}
	return nil
}

// 删除模拟器
func DelSimulator(id uint32) error {
	if err := (&db.Simulator{Id: id}).DeleteById(); err != nil {
		utils.WriteLog(err.Error())
		return err
	}
	return nil
}

// 设置一个rom的模拟器id
func SetRomSimId(romId uint64, simId uint32) error {

	if err := (&db.Rom{}).UpdateOneField(romId, "sim_id", simId); err != nil {
		return err
	}

	rom, _ := (&db.Rom{}).GetById(romId)

	//更新csv
	go func() {
		config.SetRomSettingOneField(rom.Platform, rom.RomName, "SimId", simId, true)
	}()

	return nil
}

// 批量设置rom的模拟器id
func BatchSetRomSimId(romIds []uint64, simId uint32) error {

	err := (&db.Rom{}).UpdateSimIdByIds(romIds, simId)
	if err != nil {
		return err
	}

	roms, _ := (&db.Rom{}).GetByIds(romIds)
	if len(roms) == 0 {
		return nil
	}
	platform := roms[0].Platform

	//更新csv
	go func() {
		for _, v := range roms {
			config.SetRomSettingOneField(platform, v.RomName, "SimId", simId, true)
		}
		config.FlushRomSetting(platform)
	}()

	return nil
}
