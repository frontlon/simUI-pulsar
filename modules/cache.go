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
	"strconv"
	"strings"
)

/**
 * 创建缓存入口
 **/
func CreateRomCache(getPlatform uint32) (string, error) {

	//重新更新下config
	config.InitConf()

	//重新更新下config
	(&db.Config{}).GetConfig(true)

	//检查更新一个平台还是所有平台
	PlatformList := map[uint32]*db.PlatformVO{}
	if getPlatform == 0 { //所有平台
		PlatformList = (&db.Platform{}).GetPlatformVOMap(true)
	} else { //一个平台
		PlatformList[getPlatform] = (&db.Platform{}).GetVOById(getPlatform, true)
	}

	if len(PlatformList) == 0 {
		return "", nil
	}

	//先检查平台，将不存在的平台数据先干掉
	ClearPlatform()

	//开始重建缓存
	for platform, pfInfo := range PlatformList {

		if pfInfo.RootPath == "" || len(pfInfo.RomPath) == 0 {
			continue
		}

		if !utils.DirExists(pfInfo.RootPath) {
			continue
		}

		//读取rom数据，创建rom链接
		romMap, err := ReadRomFiles(platform)
		if err != nil {
			utils.WriteLog(err.Error())
			continue
		}

		romlist, err := GetRomData(platform, romMap)
		if err != nil {
			utils.WriteLog(err.Error())
			continue
		}

		//第三步：更新rom数据
		updateRomDB(platform, romlist)

		//第四步：读取rom目录
		menu, _ := CreateMenuList(platform)

		//第五步：更新menu数据
		UpdateMenuDB(platform, menu)

		//第六步：更新filter数据
		updateFilterDB(platform, romlist)

		//第六步：清理rom_config和rom_sub
		clearRomConfigFile(platform, romlist)
	}

	//收缩数据库
	db.Vacuum()

	return "", nil

}

func ReadRomFiles(platform uint32) (map[string]*db.Rom, error) {

	romlist := map[string]*db.Rom{}
	PlatformInfo := (&db.Platform{}).GetVOById(platform, false)
	RomPaths := PlatformInfo.RomPath //rom文件路径

	//获取扩展名并转换成map
	RomExtMap := map[string]bool{}
	for _, v := range PlatformInfo.RomExts {
		RomExtMap[v] = true
	}

	linkMap, _ := components.ReadAllRomLinks(PlatformInfo.LinkPath)

	for _, RomPath := range RomPaths {
		relRomPath := utils.ToRelPath(RomPath, "")
		if err := filepath.Walk(RomPath,
			func(p string, f os.FileInfo, err error) error {

				if f == nil {
					return nil
				}

				if f.IsDir() {
					return nil
				}

				romExists := false //rom是否存在

				romExt := strings.ToLower(utils.GetFileExt(f.Name())) //获取文件后缀
				if _, ok := RomExtMap[romExt]; ok {
					romExists = true //rom存在
				} else if romExt == "" {
					//没有扩展名的rom
					if _, ok = RomExtMap[".noext"]; ok {
						romExists = true //rom存在
					}
				}

				//该类型不导入
				if !romExists {
					return nil
				}

				//转换为相对路径
				pRel := utils.ToRelPath(p, "")

				rp := utils.GetFilePath(strings.Replace(pRel, relRomPath+"/", "", 1))
				if rp != "" && strings.HasPrefix(rp, "__") {
					//隐藏目录下的rom不加载
					return nil
				}

				romName := utils.GetFileName(f.Name())
				//rom路径下的相对路径
				romPathRel := utils.ToRelPath(p, RomPath)
				memu := utils.HandleRomMenu(romPathRel)
				size := f.Size()

				//不存在链接文件，创建
				if _, ok := linkMap[romName]; !ok {
					//创建链接文件
					linkPath := PlatformInfo.LinkPath + "/" + memu
					linkFileName := romName + components.RomLinkExt
					linkFile := utils.ToAbsPath(linkFileName, linkPath)
					utils.CreateDir(linkPath)
					utils.CreateFile(linkFile, pRel)
				}

				//如果是链接文件，文件大小找实体文件
				if romExt == components.RomLinkExt {
					size = 0
					link := components.GetLinkFileData(p)
					if link.AbsRomPath != "" {
						st, err := os.Stat(link.AbsRomPath)
						if err == nil {
							size = st.Size()
						}
					}
				}

				rinfo := &db.Rom{
					LinkFile: pRel, //rom路径
					RomName:  utils.GetFileName(pRel),
					Size:     utils.GetFileSizeString(size),
				}
				romlist[romName] = rinfo

				return nil
			}); err != nil {
			return nil, errors.New(config.Cfg.Lang["romPathNotFound"])
		}
	}

	//清理掉不存在的slnk文件
	if len(linkMap) > 0 {
		tmpMap := map[string]int{}
		for _, v := range romlist {
			tmpMap[v.RomName] = 1
		}
		for k, v := range linkMap {
			if _, ok := tmpMap[k]; !ok {
				absPath := utils.ToAbsPath(v, "")
				linkData := components.GetLinkFileData(absPath)
				if !utils.IsExist(linkData.AbsRomPath) {
					utils.FileDelete(absPath)
				}
			}
		}
	}

	return romlist, nil
}

/**
 * 扫描link目录，读取数据
 **/
func GetRomData(platform uint32, romMap map[string]*db.Rom) ([]*db.Rom, error) {

	romlist := []*db.Rom{}

	//读取平台信息
	PlatformInfo := (&db.Platform{}).GetVOById(platform, false)

	//csv游戏信息
	BaseInfo, err := config.GetRomBase(platform, true)

	//rom配置信息
	romSetting, err := config.GetRomSettingByPlatform(platform)

	//读取子游戏关系
	romSubGame, _ := config.GetSubGame(PlatformInfo.Id, false)

	if err != nil {
		return nil, errors.New(config.Cfg.Lang["csvFormatError"] + err.Error())
	}

	//进入循环，遍历文件
	if err := filepath.Walk(PlatformInfo.LinkPath,
		func(p string, f os.FileInfo, err error) error {

			p = utils.ReplcePathSeparator(p)

			if f == nil {
				return nil
			}

			if f.IsDir() {
				return nil
			}

			//前缀为__的目录为隐藏目录
			rel := utils.ToRelPath(p, PlatformInfo.LinkPath)
			if strings.HasPrefix(rel, "__") {
				return nil
			}

			//检查是不是链接文件
			linkExt := strings.ToLower(utils.GetFileExt(f.Name()))
			if linkExt != components.RomLinkExt {
				return nil
			}

			linkData := components.GetLinkFileData(p)

			//读取rom数据，先找文件id，没有文件id找文件路径，都没有则删除链接文件
			romData := &db.Rom{}

			if _, ok := romMap[linkData.RomName]; ok {
				romData = romMap[linkData.RomName]
			} else {
				//rom不存在，可能是pc游戏或ps3游戏
				romData.RomName = linkData.RomName
			}

			//转换为相对路径
			fileName := romData.RomName
			fileSize := romData.Size

			menu := utils.GetFilePath(strings.Replace(p, PlatformInfo.LinkPath, "", 1))
			//先读取基础数据，如果没有基础数据，则读取别名
			title := romData.RomName
			base := &config.RomBase{}
			if _, ok := BaseInfo[fileName]; ok {
				base = BaseInfo[fileName]
				if base.Name != "" && title == linkData.RomName {
					title = base.Name
				}
			}

			score, _ := strconv.ParseFloat(base.Score, 64)
			if score > 0 {
				rs := []rune(utils.ToString(base.Score))
				if len(rs) > 3 {
					round, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", score), 64)
					base.Score = utils.ToString(round)
				}
			} else {
				base.Score = ""
			}

			hide := uint8(0)
			fav := uint8(0)
			runNum := uint64(0)
			runLastTime := int64(0)
			simId := uint32(0)
			complete := uint8(0)
			simSetting := ""
			if romSetting[fileName] != nil {
				complete = uint8(utils.ToInt(romSetting[fileName].Complete))
				hide = uint8(utils.ToInt(romSetting[fileName].Hide))
				fav = uint8(utils.ToInt(romSetting[fileName].Favorite))
				runNum = uint64(utils.ToInt(romSetting[fileName].RunNum))
				runLastTime = int64(utils.ToInt(romSetting[fileName].RunLasttime))
				simId = uint32(utils.ToInt(romSetting[fileName].SimId))
				simSetting = romSetting[fileName].SimSetting
			}

			pid := ""
			if _, ok := romSubGame[romData.RomName]; ok {
				if _, ok := romMap[romSubGame[romData.RomName]]; ok {
					pid = romMap[romSubGame[romData.RomName]].RomName
				}
			}

			rinfo := &db.Rom{
				LinkFile:      utils.ToRelPath(p, PlatformInfo.LinkPath),
				Menu:          menu,
				Name:          title,
				Platform:      platform,
				Pid:           pid,
				RomName:       romData.RomName,
				Hide:          hide,
				Favorite:      fav,
				Size:          fileSize,
				RunNum:        runNum,
				RunLasttime:   runLastTime,
				Score:         utils.ToFloat64(base.Score),
				Complete:      complete,
				SimId:         simId,
				BaseType:      base.Type,
				BaseYear:      base.Year,
				BaseProducer:  base.Producer,
				BasePublisher: base.Publisher,
				BaseCountry:   base.Country,
				BaseTranslate: base.Translate,
				BaseVersion:   base.Version,
				BaseNameEn:    base.NameEN,
				BaseNameJp:    base.NameJP,
				BaseOtherA:    base.OtherA,
				BaseOtherB:    base.OtherB,
				BaseOtherC:    base.OtherC,
				BaseOtherD:    base.OtherD,
				SimSetting:    simSetting,
				Pinyin:        utils.TextToPinyin(title),
			}
			rinfo.InfoMd5 = rinfo.CreateInfoMd5()

			romlist = append(romlist, rinfo)

			return nil
		}); err != nil {
		return nil, errors.New(config.Cfg.Lang["romPathNotFound"])
	}

	return romlist, nil
}

// 创建菜单列表
func CreateMenuList(platform uint32) (map[string]*db.Menu, error) {

	menuList := map[string]*db.Menu{}

	platformInfo := (&db.Platform{}).GetVOById(platform, false)

	filepath.WalkDir(platformInfo.LinkPath, func(p string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			return nil
		}

		p = utils.ReplcePathSeparator(p)

		p = strings.Replace(p, platformInfo.LinkPath, "", 1)

		if p == "" {
			return nil
		}

		//前缀为__的目录为隐藏目录
		if strings.HasPrefix(p, "/__") {
			return nil
		}

		pth := p + "/"

		arr := strings.Split(p, "/")
		if len(arr) > 3 {
			return nil
		}

		name := utils.GetFileName(p)
		menuList[pth] = &db.Menu{
			Platform: platform,
			Path:     pth,
			Name:     name,
			Pinyin:   utils.TextToPinyin(name),
		}

		return err
	})
	return menuList, nil
}

/**
 * 删除不存在平台的缓存数据
 **/
func ClearPlatform() {
	pfs := []string{}
	platforms := (&db.Platform{}).GetPlatformVOList(false)

	for _, v := range platforms {
		pfs = append(pfs, utils.ToString(v.Id))
	}

	//清空不存在的平台（rom表）
	if err := (&db.Rom{}).DeleteByNotPlatform(pfs); err != nil {
		fmt.Println(err)
	}

	//清空不存在的平台（menu表）
	if err := (&db.Menu{}).DeleteByNotPlatform(pfs); err != nil {
		fmt.Println(err)
	}

	//清空不存在的平台（filter表）
	if err := (&db.Filter{}).DeleteByNotPlatform(pfs); err != nil {
		fmt.Println(err)
	}

}

/**
 * 更新rom cache
 **/
func updateRomDB(platform uint32, romlist []*db.Rom) {

	romlistMd5s := []string{}         //磁盘文件
	romlistMd5Ids := map[string]int{} //md5关联romId
	for k, v := range romlist {       //从romlist列表中抽出md5
		romlistMd5s = append(romlistMd5s, v.InfoMd5)
		romlistMd5Ids[v.InfoMd5] = k
	}

	//数据库中抽出md5
	DbInfoMd5, _ := (&db.Rom{}).GetIdentByPlatform(platform)
	addMd5s := utils.SliceDiff(romlistMd5s, DbInfoMd5) //新增的
	subMd5s := utils.SliceDiff(DbInfoMd5, romlistMd5s) //删除的

	//分批次批量写入
	if len(addMd5s) > 0 {
		addChunks := utils.SliceChunk[string](addMd5s, constant.DB_ADD_MAX_NUM)
		for _, chunk := range addChunks {
			addData := []*db.Rom{}
			for _, md5 := range chunk {
				addData = append(addData, romlist[romlistMd5Ids[md5]])
			}
			fmt.Println("新增了", len(chunk), "条数据")
			(&db.Rom{}).BatchAdd(addData, 1)
		}
	}

	//分批次批量删除
	if len(subMd5s) > 0 {
		subChunks := utils.SliceChunk[string](subMd5s, constant.DB_ADD_MAX_NUM)
		for _, chunk := range subChunks {
			fmt.Println(platform, "删除了", len(chunk), "条数据")
			(&db.Rom{}).DeleteByInfoMd5s(platform, chunk)
		}
	}
}

/**
 * 更新menu cache
 **/
func UpdateMenuDB(platform uint32, menumap map[string]*db.Menu) {

	//磁盘中目录列表
	diskMenus := []string{}
	if len(menumap) > 0 {
		for k, _ := range menumap {
			diskMenus = append(diskMenus, k)
		}
	}

	//数据库中目录列表
	dbPaths, err := (&db.Menu{}).GetAllPathByPlatform(platform)
	if err != nil {
		return
	}

	add := utils.SliceDiff(diskMenus, dbPaths)
	sub := utils.SliceDiff(dbPaths, diskMenus)

	//删除当前平台下不存在的菜单
	(&db.Menu{}).DeleteNotExists(platform, sub)

	//取出需要写入数据库的rom数据。
	saveMenulist := []*db.Menu{}
	if len(add) > 0 {
		for _, v := range add {
			saveMenulist = append(saveMenulist, menumap[v])
		}
	}

	//保存数据到数据库cate表
	if len(saveMenulist) > 0 {
		for _, v := range saveMenulist {
			if err := v.Add(); err != nil {
			}
		}

	}

	//这些变量较大，写入完成后清理变量
	saveMenulist = []*db.Menu{}
	menumap = map[string]*db.Menu{}

}

/**
 * 更新filter cache
 **/
func updateFilterDB(platform uint32, romlist []*db.Rom) {

	dbf, _ := (&db.Filter{}).GetByPlatform(platform)
	romFilters, dbFilter := components.FilterFactory(romlist, dbf)

	//增加过滤器
	if len(romFilters) > 0 {
		addFilter := []string{}
		for t, v := range romFilters {
			if _, ok := dbFilter[t]; !ok {
				//如果数据库没数据，则添加全部
				addFilter = v
			} else {
				addFilter = utils.SliceDiff(romFilters[t], dbFilter[t])
			}

			//开始写入数据库
			if len(addFilter) > 0 {
				create := []*db.Filter{}
				for _, name := range addFilter {
					f := &db.Filter{
						Platform: platform,
						Type:     t,
						Name:     name,
					}
					create = append(create, f)
				}
				(&db.Filter{}).BatchAdd(create)
			}

		}
	}

	//删除过滤器
	if len(dbFilter) > 0 {
		subFilter := []string{}
		for t, v := range dbFilter {
			if _, ok := romFilters[t]; !ok {
				//如果数据库没数据，则添加全部
				subFilter = v
			} else {
				subFilter = utils.SliceDiff(dbFilter[t], romFilters[t])
			}

			if len(subFilter) > 0 {
				(&db.Filter{}).DeleteByFileNames(platform, t, subFilter)
			}
		}
	}

}

// 清理配置文件
func clearRomConfigFile(platform uint32, romlist []*db.Rom) {

	if len(romlist) == 0 {
		return
	}

	romSetting, _ := config.GetRomSettingByPlatform(platform)

	if len(romSetting) == 0 {
		return
	}

	nameMap := map[string]bool{}
	for _, v := range romlist {
		nameMap[v.RomName] = true
	}

	//清理romsetting项
	flush := false
	for _, v := range romSetting {
		if _, ok := nameMap[v.RomName]; !ok {
			flush = true
			delete(romSetting, v.RomName)
		}
	}
	if flush {
		config.CoverRomSettingFile(platform, romSetting)
	}
}

/**
 * 清空缓存
 */
func TruncateRomCache(getPlatform uint32) error {

	//检查更新一个平台还是所有平台
	PlatformList := map[uint32]*db.PlatformVO{}
	if getPlatform == 0 { //所有平台
		platforms := (&db.Platform{}).GetPlatformVOMap(false)
		PlatformList = platforms
	} else { //一个平台
		platformInfo := (&db.Platform{}).GetVOById(getPlatform, false)
		if platformInfo != nil {
			PlatformList[getPlatform] = platformInfo
		}
	}

	//开始重建缓存
	for platform, _ := range PlatformList {

		//清空rom表
		if err := (&db.Rom{Platform: platform}).DeleteByPlatform(); err != nil {
			return err
		}

		//清空menu表
		if err := (&db.Menu{Platform: platform}).DeleteByPlatform(); err != nil {
			return err
		}

	}

	//收缩数据库
	db.Vacuum()

	/*if _, err := utils.Window.Call("CB_clearDB"); err != nil {
	}*/

	return nil

}

/**
 * 清理游戏统计信息
 */
func ClearGameStat() error {
	//清理db
	_ = (&db.Rom{}).TruncateGameStat()

	//清理csv
	platforms := (&db.Platform{}).GetPlatformVOList(false)
	for _, pf := range platforms {
		settingMap, _ := config.GetRomSettingByPlatform(pf.Id)
		for name, _ := range settingMap {
			if _, ok := settingMap[name]; ok {
				settingMap[name].RunLasttime = ""
				settingMap[name].RunNum = ""
			}
		}
		//更新数据
		_ = config.CoverRomSettingFile(pf.Id, settingMap)
	}
	return nil
}

/**
 * 清理不存在的配置项
 */
func ClearNotExistGameConfig() error {

	//读平台列表
	platforms := (&db.Platform{}).GetPlatformVOList(false)
	for _, pf := range platforms {
		//读取rom数据
		roms, err := (&db.Rom{}).GetByPlatform(pf.Id)
		if err != nil {
			return err
		}
		romMap := map[string]int8{}
		for _, v := range roms {
			romMap[v.RomName] = 1
		}

		//清理游戏配置
		settingMap, err := config.GetRomSettingByPlatform(pf.Id)
		if err != nil {
			fmt.Println(err)
			continue
		}
		clearSetting := false
		for name, _ := range settingMap {
			if _, ok := romMap[name]; !ok {
				delete(settingMap, name)
				clearSetting = true
			}
		}
		if clearSetting {
			_ = config.CoverRomSettingFile(pf.Id, settingMap)
		}

		//清理csv
		rombaseMap, err := config.GetRomBase(pf.Id, false)
		if err != nil {
			fmt.Println(err)
			continue
		}
		clearRombase := false
		for name, _ := range rombaseMap {
			if _, ok := romMap[name]; !ok {
				delete(rombaseMap, name)
				clearRombase = true
			}
		}
		if clearRombase {
			if err = config.CoverRomBaseFile(pf.Id, rombaseMap); err != nil {
				return err
			}
		}

		//清理子游戏配置
		subgameMap, err := config.GetSubGame(pf.Id, false)
		if err != nil {
			fmt.Println(err)
			continue
		}
		clearSubgame := false
		for name, _ := range subgameMap {
			if _, ok := romMap[name]; !ok {
				delete(subgameMap, name)
				clearSubgame = true
			}
		}
		if clearSubgame {
			if err = config.WriteDataToSubGameFile(pf.SubGameFile, subgameMap); err != nil {
				return err
			}
		}
	}
	return nil
}
