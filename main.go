package main

import (
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"simUI/config"
	"simUI/constant"
	"simUI/controller"
	"simUI/db"
	"simUI/modules"
	"simUI/utils"
	"strings"
)

//go:embed all:frontend/dist
var assets embed.FS

type FileLoader struct {
	http.Handler
}

func main() {

	rootpath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	//测试环境配置
	if utils.FileExists(".env") {
		rootpath, _ = utils.ReadFile(".env", true)
		constant.DEV = true
	}

	rootpath = utils.ReplcePathSeparator(rootpath)
	constant.VERSION_NO = "1.0.0"                     //当前版本id
	constant.ROOT_PATH = rootpath + "/"               //当前软件的绝对路径
	constant.LANG_PATH = rootpath + "/language/"      //语言目录
	constant.CACHE_PATH = rootpath + "/cache/"        //缓存路径
	constant.RESOURCE_PATH = rootpath + "/resources/" //软件资源目录

	constant.CACHE_UNZIP_PATH = constant.CACHE_PATH + "unzip/"     //解压缓存路径
	constant.CACHE_THUMB_PATH = constant.CACHE_PATH + "thumb_bak/" //展示图备份目录
	constant.CACHE_UNOWNED_PATH = constant.CACHE_PATH + "unowned/" //无效资源备份目录
	constant.CACHE_REPEAT_PATH = constant.CACHE_PATH + "repeat/"   //重复ROM备份目录

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

	//检查数据库升级
	modules.UpgradeDB()

	//初始化配置
	config.Cfg = &config.ConfStruct{}
	errConf := config.InitConf()
	if errConf != nil {
		utils.DialogError("Error", errConf.Error())
		fmt.Println(errConf)
		os.Exit(1)
		return
	}

	if len(config.Cfg.Lang) == 0 {
		utils.WriteLog("没有找到语言文件或语言文件为空\n language files is not exists")
		utils.DialogError("Error", "没有找到语言文件或语言文件为空\n language files is not exists")
		return
	}

	//读软件配置
	conf := (&db.Config{}).GetConfig(false)

	defer func() {

		//删除解压的缓存
		utils.DeleteDir(constant.CACHE_UNZIP_PATH)

		if r := recover(); r != nil {
			var trace [1024]byte
			n := runtime.Stack(trace[:], true)
			utils.WriteLog("==================")
			utils.WriteLog("recover:" + fmt.Sprintf("%s", r))
			utils.WriteLog("trace:" + string(trace[:n]))
			utils.WriteLog("==================")
			fmt.Println("recover:", fmt.Sprintf("%s", r))
			fmt.Println("trace:" + string(trace[:n]))
		}
	}()

	var stat = options.Normal
	if conf.WindowState == 1 {
		stat = options.Maximised
	} else if conf.WindowState == 3 {
		stat = options.Fullscreen
	}

	if conf.WindowWidth < 800 {
		conf.WindowWidth = 800
	}
	if conf.WindowHeight < 600 {
		conf.WindowHeight = 600
	}

	// https://wails.io/zh-Hans/docs/reference/options
	app := NewApp()
	ctl := controller.NewController()
	err := wails.Run(&options.App{
		Title:                    conf.SoftName.Name,
		Width:                    conf.WindowWidth,
		Height:                   conf.WindowHeight,
		MinWidth:                 800,
		MinHeight:                600,
		Frameless:                true,         //无边框
		WindowStartState:         stat,         //窗口启动状态
		StartHidden:              constant.DEV, //启动时隐藏窗口
		LogLevel:                 logger.WARNING,
		LogLevelProduction:       logger.ERROR,
		EnableDefaultContextMenu: constant.DEV,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: NewFileLoader(),
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
			ctl,
		},
		Windows: &windows.Options{
			IsZoomControlEnabled: true,            //启用窗口缩放
			ZoomFactor:           conf.WindowZoom, //窗口缩放比例
			ResizeDebounceMS:     200,
			//DisableFramelessWindowDecorations: false,
			//Messages: &windows.Messages{
			//	InstallationRequired: "",
			//	UpdateRequired:       "",
			//	MissingRequirements:  "",
			//	Webview2NotInstalled: "",
			//	Error:                "",
			//	FailedToInstall:      "",
			//	DownloadPage:         "",
			//	PressOKToInstall:     "",
			//	ContactAdmin:         "",
			//	InvalidFixedWebview2: "",
			//	WebView2ProcessCrash: "",
			//},
		},
	})

	if err != nil {
		utils.DialogError("Error", err.Error())
	}

}

func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	if strings.Contains(strings.ToLower(req.URL.Path), "http") {
		return
	}

	f := req.URL.Path
	if utils.WailsPathCheck(f) {
		f = utils.WailsPathDecode(f)
	} else if !utils.IsAbsPath(f) {
		f = utils.ToAbsPath(f, "")
	}

	fileData, err := os.ReadFile(f)
	if err != nil {
		fmt.Println("没有找到磁盘文件:", f)
		return
	} else {
		fmt.Println("找到磁盘文件:", f)
		res.Write(fileData)
	}

}
