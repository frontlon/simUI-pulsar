package modules

import (
	"errors"
	"simUI/components"
	"simUI/config"
	"simUI/db"
	"simUI/utils"
)

// rom文件重命名
func RenameRomFile(id uint64, newName string) error {

	if newName == "" {
		return errors.New("文件名不能为空")
	}

	//读取老信息
	rom, err := (&db.Rom{}).GetById(id)

	platform := (&db.Platform{}).GetVOById(rom.Platform, false)

	//读取旧rom真实路径
	absPath := utils.ToAbsPath(rom.LinkFile, platform.LinkPath)
	slnkData := components.GetLinkFileData(absPath)
	if slnkData == nil {
		return nil
	}

	if newName == slnkData.RomName || newName == "" { //如果名称一样则不用修改
		return nil
	}

	//如果新名称带扩展名，则去掉扩展名
	oldExt := utils.GetFileExt(slnkData.RelRomPath)
	newExt := utils.GetFileExt(newName)
	if oldExt == newExt {
		newName = utils.GetFileName(newName)
	}

	newRomPath, _ := utils.FileRename(slnkData.AbsRomPath, newName)
	newRomPath = utils.ToRelPath(newRomPath, "")

	//处理链接文件
	links, _ := components.ReadRomLinksByRom(platform.LinkPath, rom.RomName)
	if len(links) > 0 {
		for _, v := range links {
			//替换链接文件内容
			components.UpdateLinkFileData(v, newRomPath, "")
			//没有别名，slnk也改名
			linkFilename := utils.GetFileName(v)
			if linkFilename == rom.RomName {
				utils.FileRename(v, newName)
			}
		}
	}

	//更新数据库 - 主rom rom_name
	err = rom.UpdateRomName(rom.Platform, rom.RomName, newName)
	if err != nil {
		return err
	}

	//更新数据库 - 子游戏关系
	err = rom.UpdatePidByPid(rom.Platform, rom.RomName, newName)
	if err != nil {
		return err
	}

	go func() {
		//更新资料文件
		rombase, err := config.GetRomBase(rom.Platform, false)
		if err == nil {
			create := &config.RomBase{}
			if _, ok := rombase[rom.RomName]; ok {
				create = rombase[rom.RomName]
				create.RomName = newName
				rombase[newName] = create
				delete(rombase, rom.RomName)
				config.CoverRomBaseFile(rom.Platform, rombase)
			}
		}

		//更新子游戏配置文件
		subConfig, _ := config.GetSubGame(rom.Platform, false)
		slaves, _ := config.GetSubGameByParent(rom.Platform, rom.RomName)
		if len(slaves) > 0 {
			for _, s := range slaves {
				subConfig[s] = newName
			}
			_ = config.WriteDataToSubGameFile(platform.SubGameFile, subConfig)
		}
		//重命名资源文件
		_ = renameResFile(newName, rom)
	}()
	return nil
}

// 修改rom别名
func RenameRomLink(id uint64, newName string) error {

	//读取老信息
	rom, err := (&db.Rom{}).GetById(id)
	if err != nil {
		return err
	}

	name := utils.GetFileName(rom.LinkFile)
	if newName == name || newName == "" { //如果名称一样则不用修改
		return nil
	}

	//读取资料数据
	rombase := config.GetRomBaseById(rom.Platform, rom.RomName)
	create := &config.RomBase{}

	if rombase != nil {
		create = rombase
		create.Name = newName
	} else {
		create.RomName = rom.RomName
		create.Name = newName
	}

	//修改配置文件
	if err = config.AddRomBase(rom.Platform, create); err != nil {
		return err
	}

	//更新数据库
	newRom := rom
	newRom.Id = id
	newRom.Name = newName
	newRom.Pinyin = utils.TextToPinyin(newName)
	newRom.InfoMd5 = newRom.CreateInfoMd5()
	err = newRom.UpdateName()
	if err != nil {
		return err
	}
	return nil
}

// rom批量重命名
func BatchRomRename(data []map[string]string) error {
	/*ids := []uint64{}
	create := map[string]map[string]string{}
	for _, v := range data {
		ids = append(ids, uint64(utils.ToInt(v["id"])))
		c := map[string]string{}
		c["id"] = v["id"]
		c["filename"] = v["filename"]
		create[c["filename"]] = c
	}
	//读取老信息
	volist, _ := (&db.Rom{}).GetByIds(ids)
	romlist := map[uint64]*db.Rom{}
	for _, v := range volist {
		romlist[v.Id] = v
		filename := utils.GetFileName(v.LinkFile)
		//同名等于没改名
		if filename == create[filename]["filename"] {
			delete(create, filename)
		}
	}

	if len(create) == 0 {
		return nil
	}

	//开始遍历修改
	for _, v := range create {
		rom := romlist[uint64(utils.ToInt(v["id"]))]
		filename := v["filename"]

		err := errors.New("")

		if err = renameResFile(filename, rom); err != nil {
			return err
		}

		//更新数据库
		fname := rom.LinkFile
		fpath := utils.GetFileAbsPath(rom.LinkFile)
		fext := utils.GetFileExt(rom.LinkFile)
		fname = filename + fext
		if fpath != "." {
			fname = fpath + "/" + filename + fext
		}

		err = (&db.Rom{Id: uint64(utils.ToInt(v["id"])), Name: filename, LinkFile: fname, Pinyin: utils.TextToPinyin(filename)}).UpdateNameByPath()
		if err != nil {
			return err
		}

	}
	*/
	return nil
}

// 修改文件名
func renameResFile(newName string, rom *db.Rom) error {
	platform := rom.Platform
	resPaths := components.GetResPath(platform, 0)

	//遍历资源目录
	for typ, rpath := range resPaths {
		//改名主资源
		masterFile := components.GetMasterRes(typ, platform, rom.RomName)
		if masterFile != "" {
			if _, err := utils.FileRename(masterFile, newName); err != nil {
				return err
			}
		}

		//改名子资源，子资源在目录中
		old := utils.ToAbsPath(rom.RomName, rpath)
		_ = utils.FolderRename(old, newName)

	}
	return nil
}
