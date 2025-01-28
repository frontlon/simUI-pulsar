package modules

import (
	"errors"
	"fmt"
	"github.com/getlantern/elevate"
	"os"
	"simUI/components"
	"simUI/config"
	"simUI/constant"
	"simUI/db"
	"simUI/utils"
	"strings"
	"time"
)

// 更新url
var upgradeUrl = "https://www.simui.net/checkVersion.html"

type Version struct {
	Upgrade int8   //是否有新版本
	Url     string //下载地址
	Version string //新版本号
	Date    string //更新日期
}

// 检查更新
func CheckUpgrade(newVersion string) (int, error) {
	configData := (&db.Config{}).GetConfig(false)
	upgrdade := 0
	if utils.CompareVersion(configData.VersionNo, newVersion) == -1 {
		upgrdade = 1
	}
	return upgrdade, nil
}

// 跳过本次更新
func JumpUpgrade(version string) error {
	return (&db.Config{}).UpsertOne("VersionNo", version)
}

// 升级数据库
func UpgradeDB() {

	//当去当前sql更新num
	oldNum := 0
	sqlUpdateNum, err := (&db.Config{}).GetByName("SqlUpdateNum")
	if err == nil {
		oldNum = utils.ToInt(sqlUpdateNum)
	}

	//如果num大于40，说明是从测试版升上来的，强制刷新，后期需要把这段代码删掉。
	if oldNum > 40 {
		oldNum = 5
	}

	//读取升级sql列表
	sqls := db.DbUpdateSqlList(constant.VERSION_NO)
	currNum := len(sqls)

	//无需升级
	if oldNum >= currNum {
		return
	}

	//备份数据库
	dbName, _ := db.GetDbFileName()
	newName := utils.GetFileName(dbName)
	newExt := utils.GetFileExt(dbName)
	t := time.Now().Format("20060102_150405")
	newFile := constant.CACHE_PATH + newName + "_bak_" + t + newExt
	utils.FileCopy(dbName, newFile)

	//升级sql
	newSqls := sqls[oldNum:]
	for _, sql := range newSqls {
		if sql == "" {
			continue
		}
		fmt.Println("升级SQL：", sql)
		if err = db.Exec(sql); err != nil {
			fmt.Println(err)
		}
	}

	//更新sql_num
	(&db.Config{}).UpdateOne("SqlUpdateNum", currNum)

	//系统alert提示
	//utils.ShowAlertAndExit(config.Cfg.Lang["DbUpgradeTitle"], config.Cfg.Lang["DbUpgradeContent"])
}

// 下载新版本
func DownloadNewVersion(url string) (string, error) {

	zipNameExt := utils.GetFileNameAndExt(url)
	zipName := utils.GetFileName(url)
	zipDst := constant.CACHE_UNZIP_PATH + zipNameExt

	//下载zip包
	if err := utils.DownloadFile(url, zipDst); err != nil {
		return "", err
	}

	//删掉旧解压目录
	zipDir := constant.CACHE_UNZIP_PATH + zipName + "/"
	if err := utils.DeleteDir(zipDir); err != nil {
		return "", err
	}

	//解压zip
	unzipPath, err := components.UnzipRom(zipDst)
	unzipPath = strings.ReplaceAll(unzipPath, "/", "\\")
	if err != nil {
		return "", errors.New(config.Cfg.Lang["upgradeUnzipFail"] + err.Error())
	}

	//删掉下载包中的data.db，防止覆盖数据库
	dbFile := unzipPath + "data.db"
	utils.FileDelete(dbFile)

	//删除zip包
	utils.FileDelete(zipDst)

	return unzipPath, nil
}

// 安装更新
func InstallUpgrade(version, unzipPath string) error {

	//创建更新bat
	upgradeFile, err := components.UpgradeCreateBat()
	if err != nil {
		return err
	}

	//运行bat
	result := elevate.Command("cmd", "/C", upgradeFile, unzipPath)
	if err := result.Start(); err != nil {
		fmt.Println(err)
		return err
	}

	//更新数据库
	if err = (&db.Config{}).UpsertOne("VersionNo", version); err != nil {
		return err
	}

	//退出当前app
	os.Exit(1)

	return nil
}
