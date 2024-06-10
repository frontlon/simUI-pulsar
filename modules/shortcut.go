package modules

import (
	"simUI/db"
	"simUI/utils"
	"strings"
)

// 读取快捷工具
func GetShortcuts() ([]*db.Shortcut, error) {
	volist, err := (&db.Shortcut{}).GetAll()
	if err != nil {
		utils.WriteLog(err.Error())
		return nil, err
	}

	for k, v := range volist {
		volist[k].Path = utils.ToAbsPath(v.Path, "")
	}
	return volist, nil
}

// 更新快捷工具
func UpdateShortcut(data []*db.Shortcut) ([]*db.Shortcut, error) {
	m := &db.Shortcut{}
	volist, err := m.GetAll()
	if err != nil {
		utils.WriteLog(err.Error())
		return nil, err
	}
	voMap := map[uint32]*db.Shortcut{}
	for _, v := range volist {
		voMap[v.Id] = v
	}

	for _, v := range data {

		if v.Name == "" && v.Path == "" {
			continue
		}

		v.Path = utils.ToRelPath(v.Path, "")
		v.Name = strings.TrimSpace(v.Name)
		//没有找到则添加
		if _, ok := voMap[v.Id]; !ok {
			v.Add()
			continue
		}

		exist := voMap[v.Id]
		//有差异则更新数据
		if v.Name != exist.Name || v.Path != exist.Path || v.Sort != exist.Sort {
			v.UpdateById()
		}
	}

	//检查删除
	dataMap := map[uint32]*db.Shortcut{}
	for _, v := range data {
		dataMap[v.Id] = v
	}
	for _, v := range volist {
		if _, ok := dataMap[v.Id]; !ok {
			v.DeleteById()
		}
	}

	return volist, nil
}

// 更新快捷软件排序
func UpdateShortcutSort(lists []uint32) error {
	if len(lists) == 0 {
		return nil
	}
	for k, pfId := range lists {
		shortcut := &db.Shortcut{
			Id:   pfId,
			Sort: uint32(k + 1),
		}
		err := shortcut.UpdateSortById()
		if err != nil {
			utils.WriteLog(err.Error())
			return err
		}
	}
	return nil
}
