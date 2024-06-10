package config

import (
	"context"
	"errors"
	"github.com/go-ini/ini"
	"os"
	"simUI/constant"
	"simUI/db"
	"simUI/utils"
)

// 配置文件
var (
	Cfg *ConfStruct //公共配置
	Ctx context.Context
)

// 配置文件
type ConfStruct struct {
	LangList []string          //语言列表
	Lang     map[string]string //语言项
}

/*
初始化读取配置
@author frontLon
*/
func InitConf() error {

	err := errors.New("")

	//config配置
	configData := (&db.Config{}).GetConfig(true)

	//语言列表
	Cfg.LangList, err = getLangList()
	if err != nil {
		return err
	}
	//语言配置定义
	Cfg.Lang, err = getLang(configData.Lang)
	if err != nil {
		return err
	}
	return nil
}

// 读取语言参数配置
func getLang(lang string) (map[string]string, error) {
	langpath := constant.LANG_PATH
	fpath := langpath + lang + ".ini"
	section := make(map[string]string)

	//如果默认语言不存在，则读取列表中的其他语言
	if !utils.FileExists(fpath) {
		if len(Cfg.LangList) > 0 {
			for langName, langFile := range Cfg.LangList {
				fpath = langpath + langFile
				//如果找到其他语言，则将第一项更新到数据库配置中
				if err := (&db.Config{}).UpdateOne("lang", langName); err != nil {
					return section, err
				}
				break
			}
		}
	}

	file, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, fpath)

	if err != nil {
		return section, err
	}

	section = file.Section("").KeysHash()
	return section, nil
}

// 读取语言文件列表
func getLangList() ([]string, error) {
	list := []string{}
	lists, _ := os.ReadDir(constant.LANG_PATH)
	for _, fi := range lists {
		if !fi.IsDir() { // 忽略目录
			list = append(list, utils.GetFileName(fi.Name()))
		}
	}
	return list, nil
}
