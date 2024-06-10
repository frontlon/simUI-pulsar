package modules

import (
	"archive/zip"
	"errors"
	"os"
	"path/filepath"
	"simUI/components"
	"simUI/config"
	"simUI/constant"
	"simUI/db"
	"simUI/utils"
)

type MergeDb struct {
	Platform    *db.Platform
	RomCount    int64
	Simulators  []*db.Simulator
	FolderCheck map[string]MergeDbFolderCheck
}

type MergeDbFolderCheck struct {
	Status string //检测结果 suc正常 warn警告 err错误
	Desc   string
	Path   string //目标路径
}

// 导出一个rom
func OutputOneRom(id uint64, zipFile string, outputSubRom, outputRes, outputConf bool) error {

	info, _ := (&db.Rom{}).GetById(id)
	platform := (&db.Platform{}).GetVOById(info.Platform, false)

	slnk := components.GetLinkFileData(utils.ToAbsPath(info.LinkFile, platform.LinkPath))

	zipFile = utils.ToAbsPath(zipFile, "")

	shareData := map[string]*config.Share{}
	shareData[info.RomName] = &config.Share{
		RomFile: utils.GetFileNameAndExt(slnk.RelRomPath),
	}

	//创建压缩文件
	compreFile, _ := os.Create(zipFile)
	zw := zip.NewWriter(compreFile)
	defer zw.Close()

	//导出rom
	if utils.FileExists(slnk.AbsRomPath) {
		openF, err := os.Open(slnk.AbsRomPath)
		defer openF.Close()
		if err == nil {
			_ = components.CompressZip(openF, constant.RES_DIR["rom"], zw)
		}
	}

	//导出资源
	if outputRes {
		paths := components.GetResPath(info.Platform, 0)
		for typ, pth := range paths {
			//主资源
			f := components.GetMasterRes(typ, info.Platform, info.RomName)
			if f != "" {
				openF, err := os.Open(slnk.AbsRomPath)
				defer openF.Close()
				if err == nil {
					_ = components.CompressZip(openF, constant.RES_DIR[typ], zw)
				}
			}

			//子资源
			fs, _ := components.GetSlaveRes(pth, info.RomName)
			if len(fs) > 0 {
				for _, v := range fs {
					openF, err := os.Open(v)
					defer openF.Close()
					if err == nil {
						subDir := constant.RES_DIR[typ] + "/" + info.RomName
						_ = components.CompressZip(openF, subDir, zw)
					}
				}
			}
		}
	}

	//导出配置
	if outputConf {
		//rombase
		rombase, err := config.GetRomBase(info.Platform, false)
		if err == nil {
			if _, ok := rombase[info.RomName]; ok {
				shareData[info.RomName].Rombase = rombase[info.RomName]
			}
		}
	}

	//写分享配置
	out := utils.ToAbsPath(constant.SHARE_FILE_NAME, constant.CACHE_PATH)
	config.WriteDataToShareFile(out, shareData)
	openF, err := os.Open(out)
	if err == nil {
		defer func() {
			openF.Close()
			utils.FileDelete(out)
		}()
		_ = components.CompressZip(openF, "", zw)
	}

	return nil
}

// OutputRomByPlatform 导出一个平台rom
func OutputRomByPlatform(id uint32, zipFile string, outputSim bool) error {

	PlatformInfo := (&db.Platform{}).GetVOById(id, false)

	/*a := "D:\\work\\go\\src\\sim-ui-pulsar\\build\\bin\\games\\fc"
	b := "D:\\work\\go\\src\\sim-ui-pulsar\\build\\bin\\games\\fc\\roms"
	c := "D:\\work\\go\\src\\sim-ui-pulsar\\build\\bin\\games\\fc"
	d := "D:\\work\\go\\src\\sim-ui-pulsar\\build\\bin\\games\\fc\\roms\\a.zip"
	e := "a\\b\\roms\\a.zip"

	aa := utils.ToRelPath(a, PlatformInfo.RootPath)
	bb := utils.ToRelPath(b, PlatformInfo.RootPath)
	cc := utils.ToRelPath(c, PlatformInfo.RootPath)
	dd := utils.ToRelPath(d, PlatformInfo.RootPath)
	ee := utils.ToRelPath(e, PlatformInfo.RootPath)
	fmt.Println(aa)
	fmt.Println(bb)
	fmt.Println(cc)
	fmt.Println(dd)
	fmt.Println(ee)
	return nil*/

	//创建压缩文件
	compreFile, _ := os.Create(zipFile)
	zw := zip.NewWriter(compreFile)
	defer zw.Close()

	//压缩平台目录
	if err := filepath.Walk(PlatformInfo.RootPath,
		func(p string, f os.FileInfo, err error) error {
			if f.IsDir() {
				return nil
			}

			openF, err := os.Open(p)
			defer openF.Close()
			if err != nil {
				return err
			}

			pth := utils.ToRelPath(utils.GetFilePath(p), PlatformInfo.RootPath)

			_ = components.CompressZip(openF, pth, zw)
			return nil
		}); err != nil {
		return errors.New(config.Cfg.Lang["romPathNotFound"])
	}

	////////
	/*shareData := map[string]*config.Share{}
	shareData[info.RomName] = &config.Share{
		RomFile: utils.GetFileNameAndExt(slnk.RelRomPath),
	}



	//导出rom
	if utils.FileExists(slnk.AbsRomPath) {
		openF, err := os.Open(slnk.AbsRomPath)
		defer openF.Close()
		if err == nil {
			_ = components.CompressZip(openF, constant.RES_DIR["rom"], zw)
		}
	}

	//导出子rom
	if outputSubRom {
		subs, _ := (&db.Rom{}).GetSubRom(info.Platform, info.RomName)
		for _, v := range subs {
			slnk = components.GetLinkFileData(utils.ToAbsPath(v.LinkFile, platform.LinkPath))
			if utils.FileExists(slnk.AbsRomPath) {
				openF, err := os.Open(slnk.AbsRomPath)
				defer openF.Close()
				if err == nil {
					_ = components.CompressZip(openF, constant.RES_DIR["rom"], zw)
				}
				shareData[info.RomName].SubGame = append(shareData[info.RomName].SubGame, utils.GetFileNameAndExt(slnk.RelRomPath))
			}
		}
	}

	//导出资源
	if outputRes {
		paths := components.GetResPath(info.Platform, 0)
		for typ, pth := range paths {
			//主资源
			f := components.GetMasterRes(typ, info.Platform, info.RomName)
			if f != "" {
				openF, err := os.Open(slnk.AbsRomPath)
				defer openF.Close()
				if err == nil {
					_ = components.CompressZip(openF, constant.RES_DIR[typ], zw)
				}
			}

			//子资源
			fs, _ := components.GetSlaveRes(pth, info.RomName)
			if len(fs) > 0 {
				for _, v := range fs {
					openF, err := os.Open(v)
					defer openF.Close()
					if err == nil {
						subDir := constant.RES_DIR[typ] + "/" + info.RomName
						_ = components.CompressZip(openF, subDir, zw)
					}
				}
			}
		}
	}

	//导出配置
	if outputConf {
		//rombase
		rombase, err := config.GetRomBase(info.Platform)
		if err == nil {
			if _, ok := rombase[info.RomName]; ok {
				shareData[info.RomName].Rombase = rombase[info.RomName]
			}
		}
	}

	//写分享配置
	out := utils.ToAbsPath(constant.SHARE_FILE_NAME, constant.CACHE_PATH)
	config.WriteDataToShareFile(out, shareData)
	openF, err := os.Open(out)
	if err == nil {
		defer func() {
			openF.Close()
			utils.FileDelete(out)
		}()
		_ = components.CompressZip(openF, "", zw)
	}*/

	return nil
}

// 导入一个rom
func InputOneShare(pid uint32, menu, romPath, zipFile string) error {

	zipFile = utils.ToAbsPath(zipFile, "")
	if !utils.FileExists(zipFile) {
		return errors.New(zipFile + "不存在")
	}

	ext := utils.GetFileExt(zipFile)
	if ext != constant.SHARE_ZIP_EXT {
		errors.New("不是SIMUI分享文件")
	}

	exists, err := components.CheckZip(zipFile)
	if err != nil || !exists {
		errors.New("不是SIMUI分享文件")
	}

	platform := (&db.Platform{}).GetVOById(pid, false)

	err = components.DecompressZip(zipFile, platform.RootPath, romPath)
	if err != nil {
		return err
	}

	shareFile := utils.ToAbsPath(constant.SHARE_FILE_NAME, platform.RootPath)
	if utils.FileExists(shareFile) {
		share, _ := config.GetShareData(shareFile)
		if share != nil {

			//读取各种配置
			saveRombase := false
			saveSubGame := false
			rombase, _ := config.GetRomBase(platform.Id, false)
			subgame, _ := config.GetSubGame(platform.Id, false)

			for romName, v := range share {

				//生成主游戏链接文件
				rf := utils.ToAbsPath(v.RomFile, romPath)
				rf = utils.ToRelPath(rf, "")
				sp := platform.LinkPath
				if menu != "" {
					sp = platform.LinkPath + menu
					utils.CreateDir(sp)
				}
				sp = utils.ToAbsPath(romName+components.RomLinkExt, "")
				if err = components.CreateLinkFile(sp, rf, ""); err != nil {
					return err
				}

				//生成子游戏链接文件
				if len(v.SubGame) > 0 {
					for _, s := range v.SubGame {
						subgame[utils.GetFileName(s)] = romName

						//子游戏链接文件
						rf = utils.ToAbsPath(s, romPath)
						rf = utils.ToRelPath(rf, "")
						sp = platform.LinkPath
						if menu != "" {
							sp = platform.LinkPath + menu
							utils.CreateDir(sp)
						}
						if err = components.CreateLinkFile(sp, rf, ""); err != nil {
							return err
						}
					}
					saveSubGame = true
				}

				//保存rombase
				if v.Rombase != nil {
					rombase[romName] = v.Rombase
					saveRombase = true
				}

			}

			//游戏资料
			if saveRombase {
				out := utils.ToAbsPath(constant.ROMBASE_FILE_NAME+".csv", platform.RootPath)
				_ = config.WriteDataToRomBaseFile(out, rombase)
			}

			//子游戏
			if saveSubGame {
				out := utils.ToAbsPath(constant.SUBGAME_FILE_NAME, platform.RootPath)
				_ = config.WriteDataToSubGameFile(out, subgame)
			}
		}

		//删除导入配置
		utils.FileDelete(shareFile)
	}

	return nil
}

// 导出一张图片
func OutputOneImage(src, dst string) error {
	if !utils.FileExists(src) {
		return errors.New(config.Cfg.Lang["fileNotFound"])
	}

	ext := utils.GetFileExt(src)
	newName := dst + ext

	return utils.FileCopy(src, newName)
}
