package modules

import (
	"errors"
	"fmt"
	"runtime"
	"simUI/components"
	"simUI/config"
	"simUI/constant"
	"simUI/db"
	"simUI/utils"
	"strings"
	"time"
)

/**
 * 运行游戏
 **/
func RunGame(romId uint64, simId uint32) error {
	//数据库中读取rom详情
	romInfo, err := (&db.Rom{}).GetById(romId)

	if err != nil {
		return err
	}

	platform := (&db.Platform{}).GetVOById(romInfo.Platform, false)

	//如果是相对路径，转换成绝对路径
	slnk := components.GetLinkFileData(utils.ToAbsPath(romInfo.LinkFile, platform.LinkPath))

	romInfo.LinkFile = slnk.AbsRomPath

	//检测rom文件是否存在
	if utils.IsExist(romInfo.LinkFile) == false {
		return errors.New(romInfo.LinkFile + " is not exists")
	}

	ext := strings.ToLower(utils.GetFileExt(romInfo.LinkFile))

	//读取模拟器信息
	if simId == 0 {
		sims, _ := (&db.Simulator{}).GetByPlatform(platform.Id)
		if len(sims) > 0 {
			for _, v := range sims {
				if romInfo.SimId == v.Id {
					simId = romInfo.SimId
					break
				}
			}
			if simId == 0 {
				simId = sims[0].Id
			}
		}
	}

	//运行游戏
	if utils.InSlice(ext, constant.RUN_EXTS) {
		//直接运行exe
		cmd := slnk.Params
		err = RunGameExe(romInfo.LinkFile, cmd)
	} else {
		//依赖模拟器
		sim, err := (&db.Simulator{}).GetDefaultByPlatform(romInfo.Platform, simId)
		if err != nil {
			return errors.New(config.Cfg.Lang["simNotFound"])
		}

		//读取rom模拟器独立配置
		simName := utils.GetFileName(sim.Path)
		romInfoVo := (&db.Rom{}).ConvertRomSimple(romInfo, []*db.Rom{})
		if _, ok := romInfoVo.SimSetting[simName]; ok {
			if romInfoVo.SimSetting[simName].Cmd != "" {
				sim.Cmd = romInfoVo.SimSetting[simName].Cmd
			}
			if romInfoVo.SimSetting[simName].RunBefore != "" {
				sim.RunBefore = romInfoVo.SimSetting[simName].RunBefore
			}
			if romInfoVo.SimSetting[simName].RunAfter != "" {
				sim.RunAfter = romInfoVo.SimSetting[simName].RunAfter
			}
			if romInfoVo.SimSetting[simName].Unzip != "" {
				sim.Unzip = romInfoVo.SimSetting[simName].Unzip
			}
		}

		err = runGameSimulator(romInfo, sim)
	}

	if err != nil {
		return err
	}

	go func() {
		//记录运行次数和时间
		rid := romId
		if romInfo.Pid != "" {
			parent, _ := (&db.Rom{}).GetMasterRom(romInfo.Platform, romInfo.Pid)
			if parent != nil {
				rid = parent.Id
			}
		}
		_ = (&db.Rom{}).UpdateRunNumAndTime(rid)

		//更新csv
		setting := config.GetRomSettingByName(romInfo.Platform, romInfo.RomName)
		runNum := utils.ToInt(setting.RunNum) + 1
		runLasttime := utils.ToString(time.Now().Unix())
		config.SetRomSettingOneField(romInfo.Platform, romInfo.RomName, "RunNum", utils.ToString(runNum), false)
		config.SetRomSettingOneField(romInfo.Platform, romInfo.RomName, "RunLasttime", runLasttime, false)
		config.FlushRomSetting(romInfo.Platform)
	}()
	return nil
}

/**
 * 运行程序
 **/
func RunProgram(filePath string) error {

	if filePath == "" {
		return errors.New("file path is empty")
	}

	ext := strings.ToLower(utils.GetFileExt(filePath))

	filePath = utils.ToAbsPath(filePath, "")

	if utils.InSlice[string](ext, constant.RUN_EXTS) {
		fmt.Println("run exe", filePath)
		return RunGameExe(filePath, []string{})
	} else {
		filePath = strings.ReplaceAll(filePath, "/", "\\")
		fmt.Println("run other", filePath)
		return RunGameExplorer(filePath)
	}
}

/**
 * 直接运行exe
 **/
func RunGameExe(romPath string, cmd []string) error {
	return components.RunGame(romPath, cmd)
}

/**
 * 通过explorer运行
 **/
func RunGameExplorer(romPath string) error {
	cmd := []string{romPath}
	return components.RunGame("explorer", cmd)
}

/**
 * 通过模拟器运行游戏
 **/
func runGameSimulator(romInfo *db.Rom, sim *db.Simulator) error {

	platform := (&db.Platform{}).GetVOById(romInfo.Platform, false)

	//读取父游戏信息
	var parent = &db.Rom{}
	var err error
	if romInfo.Pid != "" {
		parent, err = (&db.Rom{}).GetMasterRom(romInfo.Platform, romInfo.Pid)
		if err != nil {
			return errors.New(config.Cfg.Lang["masterGameNotFound"])
		}
		slnk := components.GetLinkFileData(utils.ToAbsPath(parent.LinkFile, platform.LinkPath))
		parent.LinkFile = slnk.AbsRomPath
	}

	//检测模拟器文件是否存在
	sim.Path = utils.ToAbsPath(sim.Path, "")
	if utils.FileExists(sim.Path) == false {
		return errors.New(config.Cfg.Lang["simNotFound"])
	}

	//解压后运行 - 解压zip包
	if sim.Unzip != "" {
		//解压文件，返回解压目录
		unzipPath, err := components.UnzipRom(romInfo.LinkFile)
		if err != nil {
			return err
		}
		//如果指定了执行文件
		romInfo.LinkFile = unzipPath + sim.Unzip
	}

	if runtime.GOOS == "windows" {
		romInfo.LinkFile = strings.ReplaceAll(romInfo.LinkFile, "/", "\\")
	}

	cmd := []string{romInfo.LinkFile}

	if sim.Cmd != "" {
		//如果rom运行参数存在，则使用rom的参数
		cmd = strings.Split(sim.Cmd, " ")
		for k, _ := range cmd {
			//替换变量
			cmd[k] = strings.ReplaceAll(cmd[k], `{RomName}`, utils.GetFileName(romInfo.LinkFile))
			cmd[k] = strings.ReplaceAll(cmd[k], `{RomAlias}`, romInfo.Name)
			cmd[k] = strings.ReplaceAll(cmd[k], `{RomExt}`, utils.GetFileExt(romInfo.LinkFile))
			cmd[k] = strings.ReplaceAll(cmd[k], `{RomFullPath}`, romInfo.LinkFile)
			cmd[k] = strings.ReplaceAll(cmd[k], `{RomMainFullpath}`, parent.LinkFile)
			cmd[k] = strings.ReplaceAll(cmd[k], `{RomMainName}`, utils.GetFileName(parent.LinkFile))
		}
	}

	//运行游戏前，先杀掉之前运行的程序
	config, _ := (&db.Config{}).GetAllVO()

	if config.GameMultiOpen == 0 {
		if err := components.KillGame(); err != nil {
			return err
		}
	}

	//模拟器运行游戏
	if err := components.RunGame(sim.Path, cmd); err != nil {
		return err
	}

	return nil
}
