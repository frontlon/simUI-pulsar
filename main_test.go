package main

import (
	"fmt"
	"os"
	"path/filepath"
	"simUI/config"
	"simUI/constant"
	"simUI/db"
	"simUI/utils"
	"strings"
	"testing"
)

func Program() {

	aa := "/ASS2ET/ASDFASDF"
	fmt.Println(strings.HasPrefix(aa, "/ASSET"))
}

func TestMain(m *testing.M) {

	rootpath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	//测试环境配置
	if utils.FileExists(".env") {
		rootpath, _ = utils.ReadFile(".env", true)
		constant.DEV = true
	}

	rootpath = strings.ReplaceAll(rootpath, "\\", "/")
	constant.ROOT_PATH = rootpath + "/"                            //当前软件的绝对路径
	constant.CACHE_PATH = rootpath + "/cache/"                     //缓存路径
	constant.CACHE_UNZIP_PATH = rootpath + "/cache/unzip/"         //解压缓存路径
	constant.LANG_PATH = rootpath + "/language/"                   //语言目录
	constant.CACHE_UNOWNED_PATH = constant.CACHE_PATH + "unowned/" //无效资源备份目录

	//创建数据库文件
	if err := db.CreateDbFile(); err != nil {
		//系统alert提示
		utils.WriteLog(err.Error())
		utils.DialogError("Error", err.Error())
		return
	}

	//连接数据库
	if err := db.Conn(); err != nil {
		//系统alert提示
		utils.WriteLog("database connect faild!")
		utils.DialogError("Error", err.Error())
		return
	}

	//初始化配置
	config.Cfg = &config.ConfStruct{}
	errConf := config.InitConf()
	if errConf != nil {
		utils.DialogError("Error", errConf.Error())
		fmt.Println(errConf)
		os.Exit(1)
		return
	}

	//数据库升级
	//modules.UpgradeDB()

	if len(config.Cfg.Lang) == 0 {
		utils.WriteLog("没有找到语言文件或语言文件为空\n language files is not exists")
		utils.DialogError("Error", "没有找到语言文件或语言文件为空\n language files is not exists")
		return
	}

	//游戏手柄
	//modules.CheckJoystick()

	//软件启动时检测升级
	//modules.BootCheckUpgrade()

	//读软件配置

	Program()

}
