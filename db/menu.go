package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"simUI/utils"
)

type Menu struct {
	Name     string // 目录名称
	Path     string // 路径
	Platform uint32 // 平台
	Pinyin   string // 拼音
	Sort     uint32 // 排序
}

func (*Menu) TableName() string {
	return "menu"
}

// 写入cate数据
func (m *Menu) Add() error {
	result := getDb().Create(&m)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return nil
}

// 根据条件，查询多条数据
func (*Menu) GetByPlatform(platform uint32) ([]*Menu, error) {
	where := map[string]interface{}{}

	if platform > 0 {
		where["platform"] = platform
	}

	volist := []*Menu{}

	//pageNum := 200
	//offset := int(pages) * pageNum
	result := getDb().Where(where).Order("sort ASC,pinyin ASC").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return volist, nil
}

// 读取所有菜单数据
func (*Menu) GetAll() ([]*Menu, error) {
	volist := []*Menu{}
	result := getDb().Select("name,platform").Order("sort ASC,pinyin ASC").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return volist, nil
}

// 读取一个平台下的所有menu数据
func (*Menu) GetAllPathByPlatform(platform uint32) ([]string, error) {

	pathList := []string{}

	volist := []*Menu{}
	result := getDb().Select("path").Where("platform = (?)", platform).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	//仅读取name
	if len(volist) > 0 {
		for _, v := range volist {
			pathList = append(pathList, v.Path)
		}
	}
	return pathList, result.Error
}

// 读取一个目录信息
func (m *Menu) GetByPath(platform uint32, path string) (*Menu, error) {
	vo := &Menu{}

	where := map[string]interface{}{
		"platform": platform,
		"path":     path,
	}

	result := getDb().Table(m.TableName()).Where(where).First(&vo)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, result.Error
	}

	return vo, nil
}

// 更新排序
func (m *Menu) UpdateSortByPath() error {
	result := getDb().Table(m.TableName()).Where("platform = ? AND path = ?", m.Platform, m.Path).Update("sort", m.Sort)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 更新目录名称
func (m *Menu) UpdateNameByPath(platform uint32, oldName string, newName string) error {

	data := map[string]interface{}{
		"name":   newName,
		"pinyin": utils.TextToPinyin(newName),
	}

	result := getDb().Table(m.TableName()).Where("platform = ? AND path = ?", platform, oldName).Updates(data)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 替换一个字段内容
func (m *Menu) ReplacePath(platform uint32, from, to string) error {
	if from == "" {
		return nil
	}
	result := getDb().Table(m.TableName()).Where("platform", platform).Update("path", gorm.Expr("REPLACE(path, ?, ?)", from, to))
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 删除一个平台下不存在的目录
func (*Menu) DeleteNotExists(platform uint32, menus []string) error {

	result := &gorm.DB{}
	m := &Menu{}
	if len(menus) == 0 {
		return nil
	}

	//数据量不会很大，慢慢删。
	tx := getDb().Begin()
	for _, v := range menus {
		tx.Where("platform=(?) AND path=(?)", platform, v).Delete(&m)
	}
	result = tx.Commit()

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return result.Error
}

// 删除不存在的平台下的所有menu
func (*Menu) DeleteByNotPlatform(platforms []string) error {
	m := &Menu{}
	result := getDb().Not("platform", platforms).Delete(&m)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return result.Error
}

// 删除平台下的所有menu
func (m *Menu) DeleteByPlatform() error {
	result := getDb().Where("platform = ?", m.Platform).Delete(&m)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return result.Error
}

// 删除一个目录，包含子目录
func (m *Menu) DeleteByPathDeep(platform uint32, path string) error {

	result := getDb().Where("platform = ? AND path LIKE ?", platform, path).Delete(&m)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return result.Error
}

// 清空表数据
func (m *Menu) Truncate() error {
	result := getDb().Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}
