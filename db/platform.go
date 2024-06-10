package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Platform struct {
	Id       uint32
	Name     string
	Icon     string
	Tag      string
	RomExts  string
	RootPath string
	RomPath  string
	Pinyin   string
	Sort     uint32
	Desc     string
	HideName uint8
	Ui       string //UI配置
}

func (*Platform) TableName() string {
	return "platform"
}

// 添加平台
func (m *Platform) Add() (uint32, error) {
	result := getDb().Create(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	m.DelPlatformVOCache()

	return m.Id, result.Error
}

// 根据条件，查询多条数据
func (*Platform) GetAll() ([]*Platform, error) {
	volist := []*Platform{}
	result := getDb().Order("sort ASC,pinyin ASC").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, nil
}

// 根据ID查询一个平台参数
func (*Platform) GetById(id uint32) (*Platform, error) {
	vo := &Platform{}
	result := getDb().Where("id=?", id).Order("sort ASC,pinyin ASC").First(&vo)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return vo, result.Error
}

// 更新平台信息
func (m *Platform) UpdateById(id uint32, create map[string]any) error {

	result := getDb().Table(m.TableName()).Where("id=?", id).Updates(create)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}

	m.DelPlatformVOCache() //清空缓存
	return result.Error
}

// 更新平台简介
func (m *Platform) UpdateDescById() error {

	create := map[string]interface{}{
		"desc": m.Desc,
	}

	result := getDb().Table(m.TableName()).Where("id=?", m.Id).Updates(create)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}

	m.DelPlatformVOCache() //清空缓存
	return result.Error
}

// 更新平台排序
func (m *Platform) UpdateSortById() error {
	result := getDb().Table(m.TableName()).Where("id=?", m.Id).Update("sort", m.Sort)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}
	m.DelPlatformVOCache() //清空缓存
	return result.Error
}

// 更新一个字段
func (m *Platform) UpdateOneField(id uint32, field string, data any) error {
	result := getDb().Table(m.TableName()).Where("id=?", id).Update(field, data)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}
	m.DelPlatformVOCache() //清空缓存
	return result.Error
}

// 删除一个平台
func (m *Platform) DeleteById() error {
	result := getDb().Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}
	m.DelPlatformVOCache() //清空缓存
	return result.Error
}

// 读取全部平台标签
func (m *Platform) GetAllPlatformTag() []string {
	volist := []*Platform{}
	result := getDb().Select("distinct(tag)").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	//去空
	tagMap := map[string]int{}
	for _, v := range volist {
		if v.Tag == "" {
			continue
		}
		tagMap[v.Tag] = 1
	}
	tags := []string{}
	for k, _ := range tagMap {
		tags = append(tags, k)
	}

	return tags
}

// 清空所有平台的UI配置
func (m *Platform) ClearPlatformUi() error {
	result := getDb().Table(m.TableName()).Update("ui", "{}")
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}
	m.DelPlatformVOCache() //清空缓存
	return result.Error
}
