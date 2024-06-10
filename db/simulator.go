package db

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Simulator struct {
	Id        uint32
	Name      string
	Platform  uint32
	Path      string // 模拟器文件路径
	Cmd       string // 启动命令参数
	RunBefore string // 游戏运行前执行脚本
	RunAfter  string // 游戏运行后执行脚本
	Unzip     string
	Sort      uint32
}

func (*Simulator) TableName() string {
	return "simulator"
}

// 写入数据
func (m *Simulator) Add() (uint32, error) {
	result := getDb().Create(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return m.Id, result.Error
}

// 根据ID查询一个模拟器参数
func (*Simulator) GetById(id uint32) (*Simulator, error) {
	vo := &Simulator{}
	result := getDb().Where("id=?", id).First(&vo)
	return vo, result.Error
}

// 查询出平台下的所有模拟器
func (*Simulator) GetByPlatform(platform uint32) ([]*Simulator, error) {

	volist := []*Simulator{}
	where := map[string]interface{}{}

	if platform != 0 {
		where["platform"] = platform
	}

	result := getDb().Where(where).Order("sort,id ASC").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, nil
}

// 读取默认模拟器
func (m *Simulator) GetDefaultByPlatform(platform uint32, romSimId uint32) (*Simulator, error) {
	if romSimId == 0 {
		//未设置默认模拟器，读取模拟器第一项
		vo := &Simulator{}
		where := map[string]interface{}{}
		where["platform"] = platform
		result := getDb().Where(where).Order("sort,id ASC").First(&vo)
		if result.Error != nil {
			return nil, result.Error
		}
		return vo, nil
	} else {
		return m.GetById(romSimId)
	}
}

// 读取所有模拟器
func (*Simulator) GetAll() ([]*Simulator, error) {
	volist := []*Simulator{}
	result := getDb().Order("sort ,id ASC").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, nil
}

// 更新
func (m *Simulator) UpdateById() error {

	create := map[string]interface{}{
		"name":       m.Name,
		"path":       m.Path,
		"cmd":        m.Cmd,
		"unzip":      m.Unzip,
		"run_before": m.RunBefore,
		"run_after":  m.RunAfter,
	}
	result := getDb().Table(m.TableName()).Where("id=(?)", m.Id).Updates(create)

	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 删除一个模拟器
func (m *Simulator) DeleteById() error {
	result := getDb().Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 删除一个平台下的所有模拟器
func (m *Simulator) DeleteByPlatform() error {
	result := getDb().Where("platform=?", m.Platform).Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 更新排序
func (m *Simulator) UpdateSortById() error {
	result := getDb().Table(m.TableName()).Where("id=?", m.Id).Update("sort", m.Sort)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}
	return result.Error
}
