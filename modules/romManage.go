package modules

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"simUI/components"
	"simUI/config"
	"simUI/constant"
	"simUI/db"
	"simUI/utils"
	"sort"
	"strings"
	"time"
)

/**
 * 复制链接文件
 **/
func CopyRomLink(id uint64, menu string) (*db.Rom, error) {

	info, err := (&db.Rom{}).GetById(id)
	if err != nil {
		return nil, err
	}
	platform := (&db.Platform{}).GetVOById(info.Platform, false)

	//复制链接文件
	src := utils.ToAbsPath(info.LinkFile, platform.LinkPath)
	dst := platform.LinkPath + menu + info.RomName + components.RomLinkExt

	if !utils.FileExists(src) {
		fmt.Println("源文件不存在", id, src, dst)
		return nil, errors.New("源文件" + src + "不存在，请更新缓存")
	}

	//文件存在，加入时间戳
	if utils.FileExists(dst) {
		name := info.RomName + "_" + utils.ToString(time.Now().Unix())
		dst = platform.LinkPath + menu + name + components.RomLinkExt
	}

	if err := utils.FileCopy(src, dst); err != nil {
		fmt.Println("复制文件出错", id, src, dst, err)
		return nil, errors.New("复制文件出错:" + err.Error())
	}

	//写数据库
	rom := info
	rom.LinkFile = utils.ToRelPath(dst, platform.LinkPath)
	rom.Menu = menu
	rom.InfoMd5 = rom.CreateInfoMd5()
	rom.Id = 0
	newId, err := rom.Add()

	if err != nil {
		return nil, err
	}
	rom.Id = newId

	return rom, nil
}

/**
 * 复制模块
 **/
func MoveRomLink(id uint64, menu string) error {

	info, err := (&db.Rom{}).GetById(id)
	if err != nil {
		return err
	}

	platform := (&db.Platform{}).GetVOById(info.Platform, false)

	//移动链接文件
	src := utils.ToAbsPath(info.LinkFile, platform.LinkPath)
	name := utils.GetFileNameAndExt(src)
	dst := platform.LinkPath + menu + name

	if utils.FileExists(dst) {
		return errors.New("目标文件已存在")
	}

	if err := utils.FileMove(src, dst); err != nil {
		return errors.New("移动文件出错:" + err.Error())
	}

	//更新数据库
	err = (&db.Rom{
		Id:       id,
		LinkFile: utils.ToRelPath(dst, platform.LinkPath),
		Menu:     menu,
	}).UpdateRoomLinkById()

	return err

}

// 删除rom链接
func DeleteRomLink(id uint64) error {

	info, err := (&db.Rom{}).GetById(id)

	if err != nil {
		return err
	}

	platform := (&db.Platform{}).GetVOById(info.Platform, false)

	//删除链接文件
	src := utils.ToAbsPath(info.LinkFile, platform.LinkPath)

	if err = utils.FileDelete(src); err != nil {
		return errors.New("文件删除失败:" + err.Error())
	}

	//删除数据
	_ = (&db.Rom{}).DeleteById(id)

	return nil
}

// 删除rom以及相关资源
func DeleteRomAndRes(id uint64, delSubRom, delRes uint8) error {

	//游戏游戏详细数据
	master, err := (&db.Rom{}).GetById(id)
	if err != nil {
		return err
	}

	//读取子游戏数据
	subs, _ := (&db.Rom{}).GetSubRom(master.Platform, master.RomName)

	//读取平台数据
	platform := (&db.Platform{}).GetVOById(master.Platform, false)

	master.LinkFile = utils.ToAbsPath(master.LinkFile, platform.LinkPath)

	//删除主数据
	_ = (&db.Rom{}).DeleteByName(master.RomName)

	//删除子游戏数据
	if delSubRom == 1 {
		_ = (&db.Rom{}).DeleteSubRom(master.Platform, master.RomName)
	} else {
		_ = (&db.Rom{}).UpdatePidByPid(master.Platform, master.RomName, "")
	}

	//删除子游戏(配置文件)
	subConfig, _ := config.GetSubGame(master.Platform, false)
	slaves, _ := config.GetSubGameByParent(master.Platform, master.RomName)
	if len(slaves) > 0 {
		for _, s := range slaves {
			delete(subConfig, s)
		}
		_ = config.WriteDataToSubGameFile(platform.SubGameFile, subConfig)
	}

	go func() {

		//删除主rom文件
		slnk := components.GetLinkFileData(master.LinkFile)
		_ = utils.FileDelete(slnk.AbsRomPath)

		//删除主游戏链接文件
		masterLinks, _ := components.ReadRomLinksByRom(platform.LinkPath, master.RomName)
		if len(masterLinks) > 0 {
			for _, v := range masterLinks {
				_ = utils.FileDelete(v)
			}
		}

		//删除资源文件
		if delRes == 1 {
			resPaths := components.GetResPath(platform.Id, 0)
			for typ, path := range resPaths {

				//删除主资源
				masterFile := components.GetMasterRes(typ, platform.Id, master.RomName)
				if masterFile != "" {
					_ = utils.FileDelete(masterFile)
				}

				//删除子资源
				dir := utils.ToAbsPath(master.RomName, path)
				_ = utils.DeleteDir(dir)
			}
		}

		//删除子游戏rom文件
		if delSubRom == 1 {
			for _, v := range subs {
				slnk = components.GetLinkFileData(utils.ToAbsPath(v.LinkFile, platform.LinkPath))
				_ = utils.FileDelete(slnk.AbsRomPath)

				//删除子游戏链接文件
				slaveLinks, _ := components.ReadRomLinksByRom(platform.LinkPath, v.RomName)
				if len(slaveLinks) > 0 {
					for _, r := range slaveLinks {
						_ = utils.FileDelete(r)
					}
				}
			}
		}
	}()
	return nil
}

// 绑定子游戏
func BindSubGame(pid uint64, sid uint64) error {
	master, err := (&db.Rom{}).GetById(pid)
	if err != nil {
		return errors.New(config.Cfg.Lang["masterGameNotFound"])
	}

	if master.Pid != "" {
		return errors.New(config.Cfg.Lang["subGameMasterExists"])
	}

	slave, err := (&db.Rom{}).GetById(sid)
	if err != nil {
		return errors.New(config.Cfg.Lang["subGameNotFound"])
	}

	//检查这个子游戏是否被绑定其他子游戏，已有子游戏的rom不能绑定
	subRoms, _ := (&db.Rom{}).GetSubRom(slave.Platform, slave.RomName)
	if len(subRoms) > 0 {
		return errors.New(config.Cfg.Lang["subGameExists"])
	}

	//更新数据库
	if (&db.Rom{}).UpdatePidByRomName(master.Platform, slave.RomName, master.RomName) != nil {
		return err
	}

	//更新子游戏配置文件
	_ = config.SetSubGame(slave.Platform, slave.RomName, master.RomName)

	return nil
}

// 解绑子游戏
func UnBindSubGame(id uint64) error {
	//读取子游戏数据
	vo, err := (&db.Rom{}).GetById(id)

	//更新数据库
	if (&db.Rom{}).UpdatePidByRomName(vo.Platform, vo.RomName, "") != nil {
		return err
	}

	//更新子游戏配置文件
	_ = config.DelSubGame(vo.Platform, vo.RomName)

	return nil
}

// 批量复制链接文件
func BatchCopyRomLink(ids []uint64, menu string) error {
	if len(ids) == 0 {
		return errors.New("请选择需要删除的项目")
	}
	for _, id := range ids {
		if _, err := CopyRomLink(id, menu); err != nil {
			return err
		}
	}
	return nil
}

// 批量移动链接文件
func BatchMoveRomLink(ids []uint64, menu string) error {
	if len(ids) == 0 {
		return errors.New("请选择需要删除的项目")
	}
	for _, id := range ids {
		if err := MoveRomLink(id, menu); err != nil {
			return err
		}
	}
	return nil
}

// 批量删除rom链接
func BatchDeleteRomLink(ids []uint64) error {
	if len(ids) == 0 {
		return errors.New("请选择需要删除的项目")
	}
	for _, id := range ids {
		if err := DeleteRomLink(id); err != nil {
			return err
		}
	}
	return nil
}

// 批量删除rom和资源
func BatchDeleteRomAndRes(ids []uint64, delSubRom, delRes uint8) error {
	if len(ids) == 0 {
		return errors.New("请选择需要删除的项目")
	}
	for _, id := range ids {
		if err := DeleteRomAndRes(id, delSubRom, delRes); err != nil {
			return err
		}
	}
	return nil
}

// 查询无效资源
func CheckUnownedRes(platform uint32) ([]map[string]string, error) {

	unownedList := []map[string]string{}
	k := 0
	platformList := []*db.Platform{}
	if platform == 0 {
		//检查全部平台
		volist, err := (&db.Platform{}).GetAll()
		if err != nil {
			return nil, err
		}
		for _, v := range volist {
			platformList = append(platformList, v)
		}
	} else {
		//检查一个平台
		vo, err := (&db.Platform{}).GetById(platform)
		if err != nil {
			return nil, err
		}
		platformList = append(platformList, vo)
	}

	for _, vo := range platformList {
		//读取已存在rom
		romlist, _ := (&db.Rom{}).GetMasterRomByPlatform(vo.Id)
		romMap := map[string]string{}
		for _, v := range romlist {
			romMap[strings.Trim(v.RomName, " ")] = ""
		}

		res := components.GetResPath(vo.Id, 0)

		//不检查这些目录
		delete(res, "upload")

		//检查重复资料
		for resName, path := range res {

			if path == "" {
				continue
			}

			if err := filepath.Walk(path,
				func(p string, f os.FileInfo, err error) error {
					if f == nil {
						return nil
					}

					if f.IsDir() {
						return nil
					}

					rel := utils.ToRelPath(p, path)
					menustr := utils.GetFilePath(rel)
					menus := utils.SplitMenu(menustr)
					name := utils.GetFileName(p)
					name = strings.Trim(name, " ")

					//开始检查
					unowned := false
					if len(menus) == 0 { //根目录，不是同名文件则标记
						if _, ok := romMap[name]; !ok {
							unowned = true
						}
					} else if len(menus) == 1 { //一级子目录，目录名不同则标记
						if _, ok := romMap[menus[0]]; !ok {
							unowned = true
						}
					} else { //深层目录，所有都标记
						unowned = true
					}

					if unowned {
						repeat := map[string]string{
							"Id":       utils.ToString(k),
							"Path":     p,
							"Platform": vo.Name,
							"ResName":  resName,
						}
						unownedList = append(unownedList, repeat)
						k++
					}

					return nil
				}); err != nil {
			}
		}
	}

	//文件名排序
	if len(unownedList) > 0 {
		sort.Slice(unownedList, func(i, j int) bool {
			return strings.ToLower(unownedList[i]["Path"]) < strings.ToLower(unownedList[j]["Path"])
		})
	}

	return unownedList, nil
}

// 删除无效资源文件
func DeleteUnownedFile(platform uint32, files []string) error {

	if len(files) == 0 {
		lists, err := CheckUnownedRes(platform)
		if err != nil {
			return err
		}
		for _, v := range lists {
			files = append(files, v["Path"])
		}
	}

	dstPath := constant.CACHE_UNOWNED_PATH + time.Now().Format("2006-01-02_15")

	for _, src := range files {
		relSrc := utils.ToRelPath(src, "")
		dst := dstPath + "/" + relSrc
		utils.CreateDir(utils.GetFilePath(dst))

		if utils.FileExists(dst) {
			filePath := utils.GetFilePath(relSrc)
			fileName := utils.GetFileName(relSrc)
			ext := utils.GetFileExt(relSrc)
			dst = dstPath + "/" + filePath + fileName + "_" + utils.ToString(time.Now().UnixNano()) + ext
		}

		//移动文件
		if err := utils.FileMove(src, dst); err != nil {
			return err
		}
	}

	return nil
}

// 打开备份文件夹
func OpenCacheFolder(typ string, create int) error {
	dir := constant.CACHE_PATH + typ
	if create == 1 && !utils.DirExists(dir) {
		utils.CreateDir(dir)
	}
	OpenFolderByPath(dir)
	return nil
}

// 检测重复rom
func CheckRomRepeat(platformId uint32) ([]map[string]any, error) {

	romlist, err := (&db.Rom{}).GetByPlatform(platformId)

	if err != nil {
		return nil, err
	}

	if len(romlist) == 0 {
		return []map[string]interface{}{}, nil
	}

	platformInfo := (&db.Platform{}).GetVOById(platformId, false)

	if len(platformInfo.RomPath) == 0 {
		return nil, errors.New("请先设置平台ROM目录")
	}

	repeatMap := map[int64][]string{}

	for _, romPath := range platformInfo.RomPath {
		if err = filepath.WalkDir(romPath, func(p string, f fs.DirEntry, err error) error {
			if f.IsDir() {
				return nil
			}

			o, err := os.Stat(p)
			if err != nil {
				return err
			}
			size := o.Size()
			if _, ok := repeatMap[size]; ok {
				repeatMap[size] = append(repeatMap[size], p)
			} else {
				repeatMap[size] = []string{p}
			}

			return nil
		}); err != nil {
			return nil, err
		}
	}

	result := []map[string]any{}

	for size, v := range repeatMap {
		if len(v) <= 1 {
			continue
		}
		for _, p := range v {
			item := map[string]any{
				"Path": p,
				"Size": size,
			}
			result = append(result, item)
		}
	}
	return result, nil
}

// 删除重复ROM文件
func DeleteRepeatFile(files []string) error {

	if len(files) == 0 {
		return errors.New("没有选择任何项目")
	}

	dstPath := constant.CACHE_REPEAT_PATH + time.Now().Format("2006-01-02_15")

	for _, src := range files {
		relSrc := utils.ToRelPath(src, "")
		dst := dstPath + "/" + relSrc
		utils.CreateDir(utils.GetFilePath(dst))

		if utils.FileExists(dst) {
			filePath := utils.GetFilePath(relSrc)
			fileName := utils.GetFileName(relSrc)
			ext := utils.GetFileExt(relSrc)
			dst = dstPath + "/" + filePath + fileName + "_" + utils.ToString(time.Now().UnixNano()) + ext
		}

		//移动文件
		if err := utils.FileMove(src, dst); err != nil {
			return err
		}
	}

	return nil
}
