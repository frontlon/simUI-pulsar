package modules

import (
	"errors"
	"simUI/components"
	"simUI/config"
	"simUI/db"
	"simUI/utils"
	"time"
)

// 读取网络图片
func LoadWebThumbs(engine, keyword string, page int) (map[string][]components.DownThumbs, error) {
	switch engine {
	case "baidu":
		return components.SearchThumbsForBaidu(keyword, page)
	case "hfsdb":
		return components.SearchThumbsForHfsDB(keyword, page)
	}
	return nil, errors.New("引擎不存在")
}

// 下载网络图片
func DownloadThumb(romId uint64, typ string, master uint8, httpUrl, ext string) (string, error) {

	info, err := (&db.Rom{}).GetById(romId)
	if err != nil {
		return "", err
	}

	resPath := components.GetResPathByType(typ, info.Platform)
	if resPath == "" {
		return "", errors.New("资源目录不存在")
	}
	if !utils.DirExists(resPath) {
		return "", errors.New("资源目录不存在")
	}

	fileName := ""
	if master == 1 {
		fileName = info.RomName + "." + ext
	} else {
		fileName = info.RomName + "/" + info.RomName + "_" + utils.ToString(time.Now().Unix()) + "." + ext
	}
	localUrl := resPath + "/" + fileName

	//如果是网络下载
	if err = utils.DownloadFile(httpUrl, localUrl); err != nil {
		return "", err
	}

	return utils.WailsPathEncode(localUrl, true), nil
}

// 添加缩略图
func AddThumb(romId uint64, typ string, imageList []string) ([]string, error) {

	vo, err := (&db.Rom{}).GetById(romId)
	if err != nil {
		return nil, err
	}

	res := components.GetResPathByType(typ, vo.Platform)
	if res == "" {
		return nil, errors.New(config.Cfg.Lang["noSetThumbDir"])
	}

	//检查是否存在主图
	master := components.GetMasterRes(typ, vo.Platform, vo.RomName)
	result := []string{}
	for k, imgPath := range imageList {

		dst := ""
		ext := utils.GetFileExt(imgPath)
		if k == 0 && master == "" {
			//设为主图
			dst = res + "/" + vo.RomName + ext
		} else {
			//设为子图
			nano := utils.ToString(time.Now().UnixNano() + int64(k))
			dst = res + "/" + vo.RomName + "/" + vo.RomName + "_" + nano + ext
		}

		//拷贝文件
		pth := utils.GetFilePath(dst)
		if !utils.DirExists(pth) {
			utils.CreateDir(pth)
		}
		if err = utils.FileCopy(imgPath, dst); err != nil {
			return nil, err
		}
		result = append(result, utils.WailsPathEncode(dst, true))
	}

	return result, nil
}

// 删除缩略图
func DeleteThumb(romId uint64, typ string, master uint8, imgPath string) error {

	imgPath = utils.WailsPathDecode(imgPath)

	vo, err := (&db.Rom{}).GetById(romId)
	if err != nil {
		return err
	}

	res := components.GetResPathByType(typ, vo.Platform)
	if res == "" {
		return errors.New(config.Cfg.Lang["noSetThumbDir"])
	}

	//备份图片
	_ = components.BackupPic(imgPath)

	//删除图片
	if err = utils.FileDelete(imgPath); err != nil {
		return err
	}

	if master == 1 {
		//第一张子图设为主图
		slaves, err := components.GetSlaveRes(res, vo.RomName)
		if err == nil && len(slaves) > 0 {
			ext := utils.GetFileExt(imgPath)
			dst := res + "/" + vo.RomName + ext
			if err = utils.FileMove(slaves[0], dst); err != nil {
				return err
			}
		}
	}

	return nil
}

// 图集排序
func SortGameThumb(romId uint64, typ string, sortList []string) ([]GameThumb, error) {

	info, err := (&db.Rom{}).GetById(romId)
	if err != nil {
		return nil, err
	}

	res := components.GetResPathByType(typ, info.Platform)
	if res == "" {
		return nil, errors.New(config.Cfg.Lang["noSetThumbDir"])
	}

	for k, v := range sortList {
		sortList[k] = utils.WailsPathDecode(v)
	}

	//读取当前排序数据
	imgList, err := GetGameThumbs(romId, typ)
	if err != nil {
		return nil, err
	}
	if len(imgList) == 0 {
		return nil, errors.New("images is not empty")
	}

	for k, v := range imgList {
		imgList[k].Path = utils.WailsPathDecode(v.Path)
	}

	//创建临时备份文件夹
	bakDir := res + "/sort_bak_" + utils.ToString(time.Now().Unix())
	if !utils.DirExists(bakDir) {
		_ = utils.CreateDir(bakDir)
	}

	//先把所有资源都移动备份目录
	bakMap := map[string]string{}
	for k, src := range imgList {
		ext := utils.GetFileExt(src.Path)
		filename := utils.GetFileName(src.Path)
		dst := bakDir + "/" + filename + "_" + utils.ToString(time.Now().UnixNano()+int64(k)) + ext
		if err = utils.FileMove(src.Path, dst); err != nil {
			return nil, err
		}
		bakMap[filename] = dst
	}

	//重新排序移回
	for k, v := range sortList {
		filename := utils.GetFileName(v)

		if _, ok := bakMap[filename]; !ok {
			continue
		}

		ext := utils.GetFileExt(v)
		src := bakMap[filename]
		dst := ""

		if k == 0 { //主图
			dst = res + "/" + info.RomName + ext
		} else { //子图
			dst = res + "/" + info.RomName + "/" + info.RomName + "_" + utils.ToString(time.Now().UnixNano()) + utils.ToString(k) + ext
		}

		if err = utils.FileMove(src, dst); err != nil {
			return nil, err
		}
	}

	//删除排序备份文件夹
	_ = utils.DeleteDir(bakDir)

	//读取新的排序数据
	return GetGameThumbs(romId, typ)
}
