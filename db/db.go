package db

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"path/filepath"
	"simUI/constant"
	"simUI/utils"
	"strings"
	"time"
)

var engine *gorm.DB
var LogMode bool = true
var maxVar = 990 //sqlite最大参数个数

// 连接数据库
func Conn() error {
	//连接数据库
	dbFile, err := GetDbFileName()

	if err != nil {
		panic(err.Error())
	}

	engine, err = gorm.Open("sqlite3", dbFile)
	if err != nil {
		panic("连接数据库失败 " + dbFile + " " + err.Error())
	}
	//调试模式下 打印日志
	engine.LogMode(LogMode)

	//配置参数
	engine.Exec("PRAGMA synchronous = OFF;")
	engine.Exec("PRAGMA journal_mode = OFF;")
	engine.Exec("PRAGMA auto_vacuum = 0;")
	engine.Exec("PRAGMA cache_size = 8000;")
	engine.Exec("PRAGMA temp_store = 2;")
	return nil
}

func getDb() *gorm.DB {
	return engine
}

// 收缩数据库
func Vacuum() {
	engine.Exec("VACUUM;")
}

// 读取数据库文件名称
func GetDbFileName() (string, error) {

	dbPath := constant.ROOT_PATH

	dbFile := dbPath + "data.db"
	files, scanErr := utils.ScanCurrentDir(dbPath)

	if utils.FileExists(dbFile) {
		//如果data.db存在，扫描data_xxx.db文件，全删掉
		clearBackDataFile(files, dbFile)
		return dbFile, nil
	}

	minDate := 0
	//如果存在多个数据库文件，则读取时间最小的那个

	if scanErr != nil {
		return "", errors.New(scanErr.Error())
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "data_") {
			cuDate := strings.ReplaceAll(file.Name(), "data_", "")
			cuDate = strings.ReplaceAll(cuDate, ".db", "")
			intDate := utils.ToInt(cuDate)
			if intDate == 0 {
				continue
			}
			if minDate == 0 {
				minDate = intDate
			} else if minDate > intDate {
				minDate = intDate
			}
		}
	}

	if minDate == 0 {
		return "", errors.New("Database does not exist")
	}

	dbFile = "data_" + utils.ToString(minDate) + ".db"

	//清理多余的data_xxx.db文件
	clearBackDataFile(files, dbFile)

	dbFile = dbPath + "/" + dbFile

	return dbFile, nil
}

// 直接运行sql
func Exec(sql string) error {
	return getDb().Exec(sql).Error
}

// 创建数据库文件
func CreateDbFile() error {

	dbFile, nil := GetDbFileName()
	dbPath := constant.ROOT_PATH
	//如果数据库文件已存在
	if dbFile != "" && utils.FileExists(dbFile) {
		//数据库文件存在，不需要升级，干掉cache下data文件
		cacheDbFile := dbPath + "/cache/data.db"
		if utils.FileExists(cacheDbFile) {
			utils.FileDelete(cacheDbFile)
			return nil
		}
		return nil
	}

	//cache目录下数据库文件
	cacheDbFile := dbPath + "/cache/data.db"
	//cache db文件不存在，则跳过
	if !utils.FileExists(cacheDbFile) {
		return nil
	}

	//移动文件
	newDbFile := dbPath + "/data.db"
	if utils.FileExists(newDbFile) {
		newDbFile = dbPath + "/data_" + utils.ToString(time.Now().Unix()) + ".db"
	}
	return utils.FileMove(cacheDbFile, newDbFile)
}

// 清理data_xxx.db文件
func clearBackDataFile(files []os.DirEntry, ignore string) {

	if files == nil || len(files) == 0 {
		return
	}

	dbPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "data_") {
			if file.Name() == ignore {
				continue
			}
			utils.FileDelete(dbPath + "/" + file.Name())
		}
	}
}
