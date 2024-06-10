package modules

import (
	"errors"
	"simUI/components"
	"simUI/config"
	"simUI/db"
	"simUI/utils"
	"sort"
	"strings"
)

type GameStrategyFile struct {
	Name string //文件名称
	Path string //文件路径
}

// 读取攻略文件
func GetStrategyFiles(id uint64) ([]GameStrategyFile, error) {

	//游戏游戏详细数据
	info, err := (&db.Rom{}).GetById(id)
	lists := []GameStrategyFile{}
	if err != nil {
		return lists, err
	}

	//资源路径
	resPath := components.GetResPathByType("file", info.Platform)

	//音频列表
	files, _ := components.GetFileRes(resPath, info.RomName)

	//文件名排序
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i]) < strings.ToLower(files[j])
	})

	if len(files) > 0 {
		for _, v := range files {
			s := GameStrategyFile{
				Name: utils.GetFileName(v),
				Path: utils.ToRelPath(v, ""),
			}
			lists = append(lists, s)
		}
	}
	return lists, nil
}

// 更新攻略文件
func UpdateStrategyFiles(id uint64, data []GameStrategyFile) ([]GameStrategyFile, error) {

	//游戏游戏详细数据
	info, err := (&db.Rom{}).GetById(id)
	if err != nil {
		return nil, err
	}

	//游戏游戏详细数据
	platformInfo := (&db.Platform{}).GetVOById(info.Platform, false)

	if platformInfo.FilesPath == "" {
		return nil, errors.New(config.Cfg.Lang["platformPathNotFound"])
	}

	//资源路径
	resPath := components.GetResPathByType("file", info.Platform)

	//攻略列表
	files, _ := components.GetFileRes(resPath, info.RomName)
	existMap := map[string]string{}
	for _, v := range files {
		existMap[utils.GetFileName(v)] = v
	}

	newData := map[string]string{}
	for _, v := range data {

		if v.Path == "" {
			continue
		}

		if v.Name == "" {
			v.Name = utils.GetFileName(v.Path)
		}

		absPath := utils.ToAbsPath(v.Path, "")
		v.Name = strings.TrimSpace(v.Name)
		pathName := utils.GetFileName(v.Path)
		pathExt := utils.GetFileExt(v.Path)
		inPath := platformInfo.FilesPath + "/" + info.RomName + "/" + v.Name + pathExt

		if !utils.FileExists(inPath) {
			//文件不存在，复制进来
			err = utils.FileCopy(absPath, inPath)
			if err != nil {
				return nil, err
			}
		} else if utils.FileExists(absPath) && pathName != v.Name {
			//文件存在，但是需要改名
			_, err = utils.FileRename(absPath, v.Name)
			if err != nil {
				return nil, err
			}
		}
		newData[v.Name] = absPath
	}

	//重新读一下音频列表，检查不存在则删除
	files, _ = components.GetFileRes(resPath, info.RomName)

	for _, p := range files {
		nfile := utils.GetFileName(p)
		if _, ok := newData[nfile]; !ok {
			utils.FileDelete(p)
			continue
		}
	}
	return GetStrategyFiles(id)
}
