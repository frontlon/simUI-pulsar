package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type RombaseAlias struct {
	Platform uint32 // 平台id
	Type     string // 类型
	Alias    string // 别名
}

func (*RombaseAlias) TableName() string {
	return "rombase_alias"
}

// 写入cate数据
func (m *RombaseAlias) Add() error {
	result := getDb().Create(&m)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return nil
}

// 更新喜爱状态
func (m *RombaseAlias) UpdateByType() error {
	//更新数据
	result := getDb().Table(m.TableName()).Where("platform=? AND type = ?", m.Platform, m.Type).Update("alias", m.Alias)
	if result.Error != nil {
		fmt.Println(result.Error)
		return result.Error
	}

	return nil
}

// 根据类型读取数据
func (*RombaseAlias) GetByPlatform(platform uint32) (map[string]string, error) {
	volist := []*RombaseAlias{}
	result := getDb().Where("platform = ?", platform).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}

	data := map[string]string{}
	if len(volist) > 0 {
		for _, v := range volist {
			data[v.Type] = v.Alias
		}
	}
	return data, result.Error
}

// 删除记录
func (m *RombaseAlias) DeleteByType() error {
	result := getDb().Where("platform = ? AND type=? ", m.Platform, m.Type).Delete(&m)
	return result.Error
}
