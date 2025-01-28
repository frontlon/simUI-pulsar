package db

import (
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"regexp"
	"simUI/constant"
	"simUI/utils"
	"strings"
)

var platformVOList []*PlatformVO

type PlatformVO struct {
	Id             uint32
	Name           string
	Icon           string
	Tag            string
	RomExts        []string
	RootPath       string   //平台根目录
	RomPath        []string //rom目录
	Sort           uint32   //排序号
	Desc           string   //平台简介
	HideName       uint8    //是否隐藏平台名称
	ThumbPath      string   //缩略图
	SnapPath       string   //插画
	PosterPath     string   //海报
	PackingPath    string   //包装盒
	TitlePath      string   //标题图
	CassettePath   string   //卡带图
	IconPath       string   //图标
	GifPath        string   //动画
	BackgroundPath string   //背景图
	OptimizedPath  string   //优化图
	DocPath        string   //游戏简介文档
	StrategyPath   string   //游戏攻略
	VideoPath      string   //游戏视频
	FilesPath      string   //游戏素材文件
	//AudioPath      string     //游戏音频
	UploadPath     string     //平台简介资源目录，放一些编辑器上传的图片
	LinkPath       string     //rom链接目录
	RomBaseFile    string     //资料文件路径
	RomSettingFile string     //rom配置文件路径
	SubGameFile    string     //子游戏配置文件
	Ui             PlatformUI //平台界面配置
}

type PlatformUI struct {
	Default  PlatformUIDefault
	Playnite PlatformUIPlaynite
	Tiny     PlatformUITiny
}

type PlatformUIDefault struct {
	NameType            uint8          // 显示名称类型(1:别名;2:文件名)
	BlockThumbType      string         // 当前缩略图显示哪个模块(Thumb,Snap,Title等)
	RomSort             uint8          // rom列表排序方式
	BaseFontsize        string         // 字体大小
	BlockSize           uint8          // 当前模块缩放等级
	BlockMargin         string         // 是否显示模块间距
	BlockDirection      uint8          // 模块显示方向（0自动、1横向、2竖向）
	BlockClickAnimate   string         // 模块点击动画
	BlockHideTitle      uint8          // 是否隐藏标题
	BlockHideBackground uint8          // 是否隐藏模块背景
	RomListStyle        uint8          // 当前列表样式(1:模块;2:表格)
	RomListColumn       []string       // 列表模式显示哪些列
	BackgroundImage     string         // 背景图片
	BackgroundRepeat    string         // 背景循环方式
	BackgroundFuzzy     uint8          // 背景模糊级别
	BackgroundMask      string         // 背景遮罩图
	Font                PlatformUIFont // 页面字体
}

type PlatformUIPlaynite struct {
	NameType            uint8          // 显示名称类型(1:别名;2:文件名)
	RomSort             uint8          // rom列表排序方式
	BaseFontsize        string         // 字体大小
	BlockThumbType      string         // 当前缩略图显示哪个模块(Thumb,Snap,Title等)
	BlockClickAnimate   string         // 模块点击动画
	BlockSize           uint8          // 当前模块缩放等级
	BlockMargin         string         // 是否显示模块间距
	BlockDirection      uint8          // 模块显示方向（0自动、1横向、2竖向）
	BlockHideTitle      uint8          // 是否隐藏标题
	BlockHideBackground uint8          // 是否隐藏模块背景
	BackgroundImage     string         // 背景图片
	BackgroundRepeat    string         // 背景循环方式
	BackgroundFuzzy     uint8          // 背景模糊级别
	BackgroundMask      string         // 背景遮罩图
	HideCarousel        uint8          // 是否隐藏相册
	Font                PlatformUIFont // 页面字体
}

type PlatformUITiny struct {
	NameType         uint8          // 显示名称类型(1:别名;2:文件名)
	RomSort          uint8          // rom列表排序方式
	BlockThumbType   string         // 当前缩略图显示哪个模块(Thumb,Snap,Title等)
	BackgroundImage  string         // 背景图片
	BackgroundFuzzy  uint8          // 背景模糊级别
	BackgroundMask   string         // 背景遮罩图
	BackgroundRepeat string         // 背景循环方式
	HideCarousel     uint8          // 是否隐藏相册
	BaseFontsize     string         // 字体大小
	Font             PlatformUIFont // 页面字体
}

type PlatformUIFont struct {
	Type   uint8  // 字体类型 1系统字体 2用户字体
	Family string // 字体名称
	Format string // 字体格式
	Src    string // 字体路径
}

/**
 * 读取平台数据(列表)
 * @param rebuild 是否重新构建 false否 true是
 **/
func (m *Platform) GetPlatformVOList(rebuild bool) []*PlatformVO {

	if rebuild == true || platformVOList == nil || len(platformVOList) == 0 {
		m.buildPlatformVO()
	}
	return platformVOList
}

/**
 * 读取平台数据(map)
 * @param rebuild 是否重新构建 false否 true是
 **/
func (m *Platform) GetPlatformVOMap(rebuild bool) map[uint32]*PlatformVO {
	if rebuild == true || platformVOList == nil || len(platformVOList) == 0 {
		m.buildPlatformVO()
	}
	mlist := map[uint32]*PlatformVO{}
	for _, v := range platformVOList {
		mlist[v.Id] = v
	}
	return mlist
}

/**
 * 读取单个平台数据
 * @param rebuild 是否重新构建 false否 true是
 **/
func (m *Platform) GetVOById(id uint32, rebuild bool) *PlatformVO {
	if rebuild == true || platformVOList == nil || len(platformVOList) == 0 {
		m.buildPlatformVO()
	}
	for _, v := range platformVOList {
		if id == v.Id {
			return v
		}
	}
	return nil
}

/**
 * 构建平台数据
 **/
func (m *Platform) buildPlatformVO() {

	platformList, _ := m.GetAll()
	lists := []*PlatformVO{}
	for _, v := range platformList {
		item := &PlatformVO{}
		absPath := utils.ToAbsPath(v.RootPath, "") + "/"

		//解析扩展名
		extarr := strings.Split(v.RomExts, ",")
		extarr = utils.SliceRemoveEmpty(extarr)
		extarr = utils.SliceRemoveDuplicate(extarr)

		//解析rom路径
		romPath := []string{}
		if v.RomPath != "" {
			romPathArr := strings.Split(v.RomPath, ";")
			for _, p := range romPathArr {
				romPath = append(romPath, utils.ToAbsPath(p, ""))
			}
		}

		item.Id = v.Id
		item.Name = v.Name
		item.Icon = utils.WailsPathEncode(v.Icon, true)
		item.Tag = v.Tag
		item.RomExts = extarr
		item.RootPath = utils.ToAbsPath(v.RootPath, "")
		item.RomPath = romPath
		item.Sort = v.Sort
		item.Desc = v.Desc
		item.HideName = v.HideName
		if item.RootPath != "" {
			item.DocPath = absPath + constant.RES_DIR["doc"]
			item.StrategyPath = absPath + constant.RES_DIR["strategy"]
			item.ThumbPath = absPath + constant.RES_DIR["thumb"]
			item.SnapPath = absPath + constant.RES_DIR["snap"]
			item.PosterPath = absPath + constant.RES_DIR["poster"]
			item.PackingPath = absPath + constant.RES_DIR["packing"]
			item.TitlePath = absPath + constant.RES_DIR["title"]
			item.CassettePath = absPath + constant.RES_DIR["cassette"]
			item.GifPath = absPath + constant.RES_DIR["gif"]
			item.BackgroundPath = absPath + constant.RES_DIR["background"]
			item.VideoPath = absPath + constant.RES_DIR["video"]
			item.FilesPath = absPath + constant.RES_DIR["file"]
			//item.AudioPath = absPath + constant.RES_DIR["audio"]
			item.IconPath = absPath + constant.RES_DIR["icon"]
			item.UploadPath = absPath + constant.RES_DIR["upload"]
			item.LinkPath = absPath + constant.RES_DIR["link"]
		}

		//资料文件
		conf := (&Config{}).GetConfig(false)
		romBaseFile := absPath + constant.ROMBASE_FILE_NAME + ".csv"                   //默认资料文件
		romBaseLang := absPath + constant.ROMBASE_FILE_NAME + "_" + conf.Lang + ".csv" //语言资料文件
		item.RomBaseFile = romBaseFile
		if utils.FileExists(romBaseLang) {
			item.RomBaseFile = romBaseLang
		}

		//rom配置文件
		item.RomSettingFile = absPath + constant.SETTING_FILE_NAME

		//子游戏配置文件
		item.SubGameFile = absPath + constant.SUBGAME_FILE_NAME

		//攻略中相对路径都改成绝对路径
		if item.Desc != "" {
			re := regexp.MustCompile(utils.WailsPathPattern)
			item.Desc = re.ReplaceAllStringFunc(item.Desc, func(match string) string {
				pth := utils.WailsPathDecode(match)
				pth = utils.ToAbsPath(pth, "")
				return utils.WailsPathEncode(pth, true)
			})
		}

		//读取界面配置
		ui := PlatformUI{}
		if v.Ui == "" || v.Ui == "{}" {
			//没有配置就读取默认配置
			conf = (&Config{}).GetConfig(false)
			if conf != nil {
				ui = conf.PlatformUi
			}
		} else {
			_ = json.Unmarshal([]byte(v.Ui), &ui)
		}

		//填充默认值
		item.Ui = m.SetUiDefault(ui)

		//检查创建资源目录
		CreateResDir(item)

		lists = append(lists, item)
	}

	platformVOList = lists
}

// 如果没有目录，则创建
func CreateResDir(item *PlatformVO) {

	//检查创建rom目录
	if len(item.RomPath) > 0 {
		for _, v := range item.RomPath {
			_ = utils.CreateDir(v)
		}
	}

	if item.RootPath == "" {
		return
	}
	//创建资源目录
	_ = utils.CreateDir(item.DocPath)
	_ = utils.CreateDir(item.StrategyPath)
	_ = utils.CreateDir(item.ThumbPath)
	_ = utils.CreateDir(item.SnapPath)
	_ = utils.CreateDir(item.PosterPath)
	_ = utils.CreateDir(item.PackingPath)
	_ = utils.CreateDir(item.TitlePath)
	_ = utils.CreateDir(item.CassettePath)
	_ = utils.CreateDir(item.GifPath)
	_ = utils.CreateDir(item.BackgroundPath)
	_ = utils.CreateDir(item.VideoPath)
	_ = utils.CreateDir(item.FilesPath)
	//_ = utils.CreateDir(item.AudioPath)
	_ = utils.CreateDir(item.IconPath)
	_ = utils.CreateDir(item.UploadPath)
	_ = utils.CreateDir(item.LinkPath)
}

// 清空平台缓存
func (m *Platform) DelPlatformVOCache() {
	platformVOList = nil //清空缓存
}
