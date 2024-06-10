package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

type MenuVO struct {
	Name     string
	Path     string
	Platform uint32
	SubMenu  []*MenuVO //子目录
}

// 根据条件，查询多条数据
func (m *Menu) GetVoByPlatform(platform uint32) ([]*MenuVO, error) {
	where := map[string]interface{}{}

	if platform > 0 {
		where["platform"] = platform
	}

	volist := []*MenuVO{}
	sublist := map[string][]*MenuVO{}

	menus, err := m.GetByPlatform(platform)
	if err != nil {
		fmt.Println(err)
		return volist, err
	}

	for _, v := range menus {
		item := &MenuVO{
			Name:     v.Name,
			Path:     v.Path,
			Platform: v.Platform,
			SubMenu:  []*MenuVO{},
		}
		patharr := strings.Split(v.Path, "/")
		if len(patharr) <= 3 {
			volist = append(volist, item)
		} else {
			if _, ok := sublist[patharr[1]]; ok {
				sublist[patharr[1]] = append(sublist[patharr[1]], item)
			} else {
				sublist[patharr[1]] = []*MenuVO{item}
			}
		}
	}

	//填充子目录
	for k, v := range volist {
		patharr := strings.Split(v.Path, "/")
		if _, ok := sublist[patharr[1]]; ok {
			volist[k].SubMenu = sublist[patharr[1]]
		}
	}

	return volist, nil
}
