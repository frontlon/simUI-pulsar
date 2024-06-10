package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"simUI/utils"
)

type Config struct {
	Name string
	Data string
	Desc string
}

func (*Config) TableName() string {
	return "config"
}

// 根据id查询一条数据
func (m *Config) Get() (map[string]string, error) {

	volist := []*Config{}
	result := getDb().Table(m.TableName()).Select("name,data").Find(&volist)

	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}

	data := map[string]string{}
	for _, v := range volist {
		data[v.Name] = v.Data
	}

	return data, result.Error
}

// 根据id查询一条数据
func (m *Config) GetByName(name string) (string, error) {

	vo := &Config{}
	result := getDb().Table(m.TableName()).Select("data").Where("name = ?", name).First(&vo)

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return "", result.Error
	}
	return vo.Data, nil
}

// 更新一个字段
func (m *Config) UpdateOne(name string, data interface{}) error {
	result := getDb().Table(m.TableName()).Where("name = ?", name).Update("data", data)

	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}
	return result.Error
}

// 更新一个字段，如果不存在则新增
func (m *Config) UpsertOne(name string, data interface{}) error {
	count := 0
	result := getDb().Table(m.TableName()).Where("name = ?", name).Count(&count)
	if result.Error != nil {
		return result.Error
	}

	if count > 0 {
		result = getDb().Table(m.TableName()).Where("name = ?", name).Update("data", data)
	} else {
		c := &Config{
			Name: name,
			Data: utils.ToString(data),
		}
		result = getDb().Create(&c)
	}

	return result.Error
}

func (m *Config) Add() {
	getDb().Create(&m)
}
