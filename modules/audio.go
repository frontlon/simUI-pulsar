package modules

// 读取rom展示图列表
/*func GetGameAudios(id uint64) ([]GameAudio, error) {

	//游戏游戏详细数据
	info, err := (&db.Rom{}).GetById(id)
	lists := []GameAudio{}
	if err != nil {
		return lists, err
	}

	//资源路径
	resPath := components.GetResPathByType("audio", info.Platform)

	//音频列表
	files, _ := components.GetAudioRes(resPath, info.RomName)

	//子资源文件名排序
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i]) < strings.ToLower(files[j])
	})

	if len(files) > 0 {
		for _, v := range files {
			s := GameAudio{
				Name: utils.GetFileName(v),
				Path: utils.WailsPathEncode(v),
			}
			lists = append(lists, s)
		}
	}

	return lists, nil
}*/

/**
 * 编辑音频
 **/
/*func UpdateGameAudio(paths []string) ([]string, error) {
	return nil, nil
}
*/
/**
 * 添加音频
 **/
/*func AddGameAudio(id uint64, paths []string) ([]string, error) {
	vo, _ := (&db.Rom{}).GetById(id)
	platformInfo := (&db.Platform{}).GetVOById(vo.Platform, false)

	if platformInfo.AudioPath == "" {
		return nil, errors.New(config.Cfg.Lang["AudioMenuCanNotBeEmpty"])
	}

	//检查创建目录
	_ = utils.CreateDir(platformInfo.AudioPath)

	result := []string{}
	for _, p := range paths {
		p = utils.ToAbsPath(p, "")

		if utils.FileExists(p) {
			continue
		}

		name := utils.GetFileNameAndExt(p)
		dst := platformInfo.AudioPath + "/" + vo.RomName + "/" + name
		if err := utils.FileCopy(p, dst); err != nil {
			return nil, err
		}
		result = append(result, utils.WailsPathEncode(dst))
	}

	return result, nil
}*/

/**
 * 音频改名
 **/
/*func RenameGameAudio(pth string, newName string) (string, error) {
	src := utils.WailsPathDecode(pth)
	dst := utils.ReplaceFileNameByPath(src, newName)
	if utils.FileExists(dst) {
		return "", errors.New("文件已存在")
	}
	err := utils.FileMove(src, dst)
	if err != nil {
		return "", err
	}
	return dst, nil
}
*/
/**
 * 删除音频
 **/
/*func DelGameAudio(pth string) error {
	f := utils.WailsPathDecode(pth)
	return utils.FileDelete(f)

}*/
