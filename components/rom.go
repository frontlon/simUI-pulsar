package components

import (
	"os"
	"path/filepath"
	"simUI/constant"
	"simUI/utils"
	"strings"
)

// 读取游戏主资源
func GetMasterRes(typ string, platform uint32, romName string) string {

	fileName := ""
	resName := ""
	pth := GetResPathByType(typ, platform)
	types := GetResExts()
	if pth != "" {
		for _, v := range types[typ] {
			fileName = pth + "/" + romName + v
			if utils.FileExists(fileName) {
				resName = fileName
				break
			}
		}
	}

	return resName
}

/**
 * 读取游戏子资源
 **/
func GetSlaveRes(dir string, masterRomName string) ([]string, error) {
	
	dir = utils.ToAbsPath(masterRomName, dir)

	files := []string{}
	extMap := utils.SliceToMap(constant.MEDIA_EXTS)

	resList, err := utils.ScanCurrentDir(dir)
	if err != nil {
		return nil, nil
	}

	for _, f := range resList {
		if f.IsDir() {
			// 忽略目录
			continue
		}

		ext := utils.GetFileExt(f.Name())
		if _, ok := extMap[ext]; !ok {
			//不是图片，忽略
			continue
		}

		abs := utils.ToAbsPath(f.Name(), dir)
		files = append(files, abs)

	}

	return files, err
}

// rom封装主展示图
func GetMasterRomThumbs(thumbType string, platform uint32, romNameMap map[string]string) map[string]string {
	dir := GetResPathByType(thumbType, platform)
	if dir == "" {
		return map[string]string{}
	}

	thumbsMap := map[string]string{}
	picMap := utils.SliceToMap(constant.MEDIA_EXTS)

	filepath.Walk(dir, func(filename string, f os.FileInfo, err error) error { //遍历目录
		if err != nil {                                                        //忽略错误
			return err
		}

		if f.IsDir() { // 忽略目录
			return nil
		}

		fExt := strings.ToLower(utils.GetFileExt(f.Name()))
		if _, ok := picMap[fExt]; !ok {
			return nil
		}

		fName := utils.GetFileName(f.Name())

		if _, ok := romNameMap[fName]; !ok {
			return nil
		}

		if _, ok := thumbsMap[fName]; ok {
			return nil
		}

		thumbsMap[fName] = filename

		return nil
	})

	return thumbsMap
}

/**
 * 读取游戏文件资源
 **/
func GetFileRes(dir string, romName string) ([]string, error) {

	dir = utils.ToAbsPath(romName, dir)

	files := []string{}

	resList, err := utils.ScanCurrentDir(dir)
	if err != nil {
		return nil, nil
	}

	for _, f := range resList {
		if f.IsDir() {
			// 忽略目录
			continue
		}

		abs := utils.ToAbsPath(f.Name(), dir)
		files = append(files, abs)
	}

	return files, err
}
