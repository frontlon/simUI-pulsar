package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"simUI/config"
	"simUI/db"
	"simUI/utils"
)

// 读取一个ROM模拟器配置
func GetRomSimSettingById(romId uint64, simId uint32) (db.RomSimSetting, error) {

	resp := db.RomSimSetting{}

	romInfo, _ := (&db.Rom{}).GetById(romId)
	if romInfo == nil {
		return resp, errors.New(config.Cfg.Lang["romIsNotExists"])
	}

	simInfo, _ := (&db.Simulator{}).GetById(simId)

	if simInfo == nil {
		return resp, errors.New(config.Cfg.Lang["simIsNotExists"])
	}

	simName := utils.GetFileName(simInfo.Path)

	data := (&db.Rom{}).ConvertRomSimple(romInfo, []*db.Rom{})

	if _, ok := data.SimSetting[simName]; ok {
		resp = db.RomSimSetting{
			Cmd:       data.SimSetting[simName].Cmd,
			RunBefore: data.SimSetting[simName].RunBefore,
			RunAfter:  data.SimSetting[simName].RunAfter,
			Unzip:     data.SimSetting[simName].Unzip,
		}
	}

	return resp, nil
}

// 更新ROM模拟器配置
func UpdateRomSimSetting(romId uint64, simId uint32, data *db.RomSimSetting) error {

	data.RunBefore = utils.ToRelPath(data.RunBefore, "")
	data.RunAfter = utils.ToRelPath(data.RunAfter, "")

	//读rom信息
	romInfo, _ := (&db.Rom{}).GetById(romId)
	if romInfo == nil {
		return errors.New(config.Cfg.Lang["romIsNotExists"])
	}
	romVo := (&db.Rom{}).ConvertRomSimple(romInfo, []*db.Rom{})
	romSetting := map[uint32]db.RomSimSetting{}
	if romInfo.SimSetting != "" {
		json.Unmarshal([]byte(romInfo.SimSetting), &romSetting)
	}

	//读模拟器信息
	simInfo, _ := (&db.Simulator{}).GetById(simId)
	if simInfo == nil {
		return errors.New(config.Cfg.Lang["simIsNotExists"])
	}

	simName := utils.GetFileName(simInfo.Path)
	romVo.SimSetting[simName] = data

	//清除空数据
	for k, v := range romVo.SimSetting {
		if v.RunBefore == "" && v.RunAfter == "" && v.Cmd == "" && v.Unzip == "" {
			delete(romVo.SimSetting, k)
		}
	}

	create, _ := json.Marshal(romVo.SimSetting)
	if err := (&db.Rom{}).UpdateSimSettingById(romId, string(create)); err != nil {
		fmt.Println("UpdateRomSimSetting", err)
		return err
	}

	//更新到配置文件
	config.SetRomSettingOneField(romInfo.Platform, romInfo.RomName, "SimSetting", string(create), true)

	return nil

}
