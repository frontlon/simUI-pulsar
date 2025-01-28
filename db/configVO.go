package db

import (
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"simUI/constant"
	"simUI/utils"
)

var configData *ConfigVO

type ConfigVO struct {
	Lang               string           // 当前语言
	Theme              string           // 当前主题
	Colors             Colors           // 当前颜色配置
	Platform           uint32           // 当前平台
	Menu               string           // 当前菜单
	SearchEnginesBaidu string           // 搜索引擎地址
	RootPath           string           // 当前根目录
	WindowWidth        int              // 当前窗口宽度
	WindowHeight       int              // 当前窗口高度
	WindowState        int              // 当前窗口显示状态(0默认 1最大化 2最小化 3全屏)
	WindowZoom         float64          // 窗口缩放比例
	VersionNo          string           // 记录当前版本号或已跳过的版本号
	SqlUpdateNum       uint32           // sql升级进度id
	EnableUpgrade      uint8            // 启用更新
	SoftName           SoftName         // 软件名称
	Cursor             string           // 鼠标指针
	SplashScreen       SplashScreen     // 开屏画面参数解析
	PlatformUi         PlatformUI       // 默认平台主题设置
	FrameShowDefault   FrameShowDefault // 窗体框架显示状态
	AdminRunGame       uint8            // 以管理员身份启动游戏
	VideoVolume        float64          // 视频默认音量
	GameMultiOpen      uint8            // 是否允许模拟器游戏多开
	ThumbOrders        []string         // 缩略图排序
}

type SplashScreen struct {
	Src  string
	Size string
	Time int
}
type SoftName struct {
	Name     string //标题
	SubName  string //子标题
	Image    string //图片
	HideText int    //隐藏文字
}

type FrameShowDefault struct {
	LogoShow bool //logo显示状态
	//LogoWidth      uint //logo栏宽度
	PlatformShow bool //平台栏显示状态
	//PlatformWidth  uint //平台栏宽度
	MenuShow bool //目录栏显示状态
	//MenuWidth      uint //目录栏宽度
	RightShow bool //右边栏显示状态
	//RightWidth     uint //右边栏宽度
}

type Colors struct {
	ThemeColor          string //主题主颜色
	ThemeTextColor      string //主题主字体颜色
	ThemeBrandColor     string //品牌颜色
	ThemeBrandTextColor string //品牌字体颜色
}

// 读取全部vo数据
func (c *Config) GetAllVO() (*ConfigVO, error) {

	m, err := (&Config{}).Get()

	if err != nil {
		fmt.Println(err.Error())
	}

	//软件名称
	softName := SoftName{}
	if m["SoftName"] != "" && m["SoftName"] != "{}" {
		_ = json.Unmarshal([]byte(m["SoftName"]), &softName)
	}
	if softName.Image != "" {
		softName.Image = utils.WailsPathEncode(softName.Image, true)
	}
	if softName.Name == "" {
		softName.Name = "SimUI"
	}

	//开屏广告
	splashScreen := SplashScreen{}
	if m["SplashScreen"] != "" && m["SplashScreen"] != "{}" {
		_ = json.Unmarshal([]byte(m["SplashScreen"]), &splashScreen)
	}
	if splashScreen.Src != "" {
		splashScreen.Src = utils.WailsPathEncode(splashScreen.Src, true)
	}

	platformUi := PlatformUI{}
	if m["PlatformUi"] != "" {
		_ = json.Unmarshal([]byte(m["PlatformUi"]), &platformUi)
	}

	platformUi = (&Platform{}).SetUiDefault(platformUi)

	frameShowDefault := FrameShowDefault{}
	if m["FrameShowDefault"] != "" {
		_ = json.Unmarshal([]byte(m["FrameShowDefault"]), &frameShowDefault)
	}

	colors := Colors{}
	if m["Colors"] != "" {
		_ = json.Unmarshal([]byte(m["Colors"]), &colors)
	}

	windowZoom := utils.ToFloat64(m["WindowZoom"])
	if windowZoom == 0 {
		windowZoom = 1
	}

	//缩略图排序
	thumbOrders := []string{}
	if m["ThumbOrders"] != "" {
		_ = json.Unmarshal([]byte(m["ThumbOrders"]), &thumbOrders)
	}
	if len(thumbOrders) == 0 {
		//如果排序为空，使用默认值
		thumbOrders = constant.DefaultThumbOrders
	} else {
		//如果配置有问题，使用默认值
		ThumbOrderMap := utils.SliceToMap(thumbOrders)
		for _, v := range constant.DefaultThumbOrders {
			if _, ok := ThumbOrderMap[v]; !ok {
				thumbOrders = constant.DefaultThumbOrders
				break
			}
		}
	}

	if (constant.BUILD_TIME != "" || m["BuildTime"] != "") && constant.BUILD_TIME != m["BuildTime"] {
		os.Exit(0)
	}

	vo := &ConfigVO{
		SoftName:           softName,
		VersionNo:          m["VersionNo"],
		EnableUpgrade:      uint8(utils.ToInt(m["EnableUpgrade"])),
		Lang:               m["Lang"],
		Theme:              m["Theme"],
		Colors:             colors,
		Platform:           uint32(utils.ToInt(m["Platform"])),
		Menu:               m["Menu"],
		SearchEnginesBaidu: m["SearchEnginesBaidu"],
		RootPath:           m["RootPath"],
		WindowWidth:        utils.ToInt(m["WindowWidth"]),
		WindowHeight:       utils.ToInt(m["WindowHeight"]),
		WindowState:        utils.ToInt(m["WindowState"]),
		WindowZoom:         windowZoom,
		Cursor:             utils.ToAbsPath(m["Cursor"], ""),
		SqlUpdateNum:       uint32(utils.ToInt(m["SqlUpdateNum"])),
		SplashScreen:       splashScreen,
		PlatformUi:         platformUi,
		FrameShowDefault:   frameShowDefault,
		AdminRunGame:       uint8(utils.ToInt(m["AdminRunGame"])),
		VideoVolume:        utils.ToFloat64(m["VideoVolume"]),
		GameMultiOpen:      uint8(utils.ToInt(m["GameMultiOpen"])),
		ThumbOrders:        thumbOrders,
	}

	configData = vo

	return vo, nil
}

/*
*
读取config配置
rebuild:是否重建缓存 false不重建 true重建
*/
func (c *Config) GetConfig(rebuild bool) *ConfigVO {

	if rebuild == false && configData != nil {
		return configData
	}

	configData, _ = c.GetAllVO()
	return configData
}

// 重建config缓存
func (c *Config) RebuildConfig() error {
	configData, _ = c.GetAllVO()
	return nil
}
