package modules

import (
	"errors"
	"simUI/config"
	"simUI/db"
	"simUI/utils"
	"strings"
)

// 读取菜单列表
func GetMenuList(platform uint32) ([]*db.MenuVO, error) {
	newMenu := []*db.MenuVO{}

	menu, err := (&db.Menu{}).GetVoByPlatform(platform) //从数据库中读取当前平台的分类目录
	if err != nil {
		return newMenu, err
	}

	if platform == 0 {
		subgameList := map[string][]*db.MenuVO{}
		for _, v := range menu {
			if _, ok := subgameList[v.Name]; !ok {
				newMenu = append(newMenu, v)
				subgameList[v.Name] = []*db.MenuVO{}
			} else {
				subgameList[v.Name] = append(subgameList[v.Name], v.SubMenu...)
			}
		}

		//写入subMenu
		for k, v := range newMenu {
			if len(subgameList[v.Name]) > 0 {
				newMenu[k].SubMenu = append(newMenu[k].SubMenu, subgameList[v.Name]...)
			}
		}

		//subMenu去重
		for k, b := range newMenu {
			isset := map[string]bool{}
			newSubMenu := []*db.MenuVO{}
			for _, v := range b.SubMenu {
				if _, ok := isset[v.Name]; !ok {
					newSubMenu = append(newSubMenu, v)
					isset[v.Name] = true
				}
			}
			newMenu[k].SubMenu = newSubMenu
		}
		return newMenu, nil
	}
	return menu, nil
}

// 添加菜单
func AddMenu(platform uint32, path, name string) error {

	platformInfo := (&db.Platform{}).GetVOById(platform, false)

	if platformInfo.RootPath == "" {
		return errors.New("平台路径未配置")
	}

	folder := platformInfo.LinkPath + path + name + "/"

	//检查数据是否存在
	if _, err := (&db.Menu{}).GetByPath(platform, folder); err == nil {
		return errors.New(config.Cfg.Lang["menuExists"])
	}

	//检查目录是否存在
	exists := utils.DirExists(folder)

	if exists {
		return errors.New(config.Cfg.Lang["menuExists"])
	}

	if err := utils.CreateDir(folder); err != nil {
		return err
	}

	//更新数据库
	if err := (&db.Menu{
		Name:     name,
		Path:     strings.Replace(folder, platformInfo.LinkPath, "", 1),
		Platform: platform,
		Pinyin:   utils.TextToPinyin(name),
		Sort:     0,
	}).Add(); err != nil {
		return err
	}

	return nil
}

// 菜单重命名
func RenameMenu(platform uint32, oldPath string, newName string) error {

	platformInfo := (&db.Platform{}).GetVOById(platform, false)

	oldMenu := platformInfo.LinkPath + oldPath

	oldMenuArr := strings.Split(oldPath, "/")
	oldMenuName := oldMenuArr[len(oldMenuArr)-2]
	oldMenuArr[len(oldMenuArr)-2] = newName
	newMenu := strings.Join(oldMenuArr, "/")
	newMenuAbs := platformInfo.LinkPath + newMenu

	//读取目录信息
	info, _ := (&db.Menu{}).GetByPath(platform, oldPath)
	if info == nil {
		return errors.New(config.Cfg.Lang["menuIsNotExists"])
	}

	if !utils.DirExists(oldMenu) {
		return errors.New(config.Cfg.Lang["menuIsNotExists"])
	}

	//目录重命名
	if err := utils.FolderMove(oldMenu, newMenuAbs); err != nil {
		return err
	}

	//菜单管理
	_ = (&db.Menu{}).UpdateNameByPath(platform, oldPath, newName)
	_ = (&db.Menu{}).ReplacePath(platform, oldMenuName, newName)

	//更新ROM数据
	romModel := &db.Rom{}
	trim1 := strings.TrimPrefix(oldPath, "/")
	trim2 := strings.TrimPrefix(newMenu, "/")
	_ = romModel.ReplaceFieldData(platform, "link_file", trim1, trim2)
	_ = romModel.ReplaceFieldData(platform, "menu", oldPath, newMenu)
	return nil
}

func DeleteMenu(platform uint32, path string) error {

	if path == "" || path == "/" {
		return errors.New("不允许操作根目录")
	}

	//读取目录信息
	info, _ := (&db.Menu{}).GetByPath(platform, path)
	if info == nil {
		return errors.New(config.Cfg.Lang["menuIsNotExists"])
	}

	roms, _ := (&db.Rom{}).GetMenuDeep(platform, path+"%")
	if len(roms) > 0 {
		return errors.New(config.Cfg.Lang["menuIsNotEmpty"])
	}

	platformInfo := (&db.Platform{}).GetVOById(platform, false)

	if platformInfo.RootPath == "" {
		return errors.New(config.Cfg.Lang["platformPathNotFound"])
	}

	//删除目录实体
	folder := platformInfo.LinkPath + path
	_ = utils.DeleteDir(folder)

	//删除目录数据
	_ = (&db.Menu{}).DeleteByPathDeep(platform, path+"%")

	return nil
}

// 更新菜单排序
func SortMenu(platform uint32, paths []string) error {
	if len(paths) == 0 {
		return nil
	}
	for k, path := range paths {
		menu := &db.Menu{
			Platform: platform,
			Path:     path,
			Sort:     uint32(k + 1),
		}
		err := menu.UpdateSortByPath()
		if err != nil {
			utils.WriteLog(err.Error())
			return err
		}
	}
	return nil
}
