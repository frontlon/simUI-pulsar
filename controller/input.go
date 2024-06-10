package controller

import (
	"errors"
	"simUI/modules"
)

// 导出一个rom
func (a *Controller) OutputOneRom(id uint64, dir string, outputSubRom, outputRes, outputConf bool) string {
	return Resp("", modules.OutputOneRom(id, dir, outputSubRom, outputRes, outputConf))
}

// 导出一个平台rom
func (a *Controller) OutputRomByPlatform(id uint32, dir string, outputSim bool) string {
	return Resp("", modules.OutputRomByPlatform(id, dir, outputSim))
}

// 添加游戏
func (a *Controller) AddGame(opt string, platform uint32, menu string, romPath string, files []string, bootParam string, isBatFile uint8) string {

	if len(files) == 0 {
		return Resp("", errors.New("files not be empty"))
	}

	var err error

	switch opt {
	case "rom":
		err = modules.AddFileGame(platform, menu, romPath, files)
	case "pc":
		err = modules.AddPcOrFolderGame(platform, menu, romPath, files[0], bootParam, isBatFile)
	case "ps3":
		err = modules.AddPcOrFolderGame(platform, menu, romPath, files[0], "", 0)
	case "share":
		err = modules.InputOneShare(platform, menu, romPath, files[0])
	}

	return Resp("", err)
}

// 导出一张图片
func (a *Controller) OutputOneImage(src, dst string) string {
	return Resp("", modules.OutputOneImage(src, dst))

}
