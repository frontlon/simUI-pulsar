package modules

import (
	"simUI/config"
	"simUI/db"
)

// 设为隐藏
func SetHide(id uint64, hide uint8) error {

	//数据库中读取rom详情
	rom, err := (&db.Rom{}).GetById(id)
	if err != nil {
		return err
	}

	//更新数据
	if err = config.SetRomSettingOneField(rom.Platform, rom.RomName, "Hide", hide, true); err != nil {
		return err
	}

	err = (&db.Rom{RomName: rom.RomName, Hide: hide}).UpdateHide()
	if err != nil {
		return err
	}

	return nil
}

// 设为喜爱
func SetFavorite(id uint64, fav uint8) error {

	//数据库中读取rom详情
	rom, err := (&db.Rom{}).GetById(id)
	if err != nil {
		return err
	}

	//更新数据
	if err = config.SetRomSettingOneField(rom.Platform, rom.RomName, "Favorite", fav, true); err != nil {
		return err
	}

	err = (&db.Rom{RomName: rom.RomName, Favorite: fav}).UpdateFavorite()
	if err != nil {
		return err
	}

	return nil
}
