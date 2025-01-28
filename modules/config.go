package modules

import (
	"encoding/json"
	"os"
	"simUI/components"
	"simUI/config"
	"simUI/constant"
	"simUI/db"
	"simUI/server"
	"simUI/utils"
)

/**
 * 读取配置数据
 **/
func GetConfig() (map[string]any, error) {
	fileDate := utils.GetFileUpdateDate(os.Args[0])

	configData := (&db.Config{}).GetConfig(true)

	result := map[string]any{
		"Config":       configData,                          //界面配置
		"Lang":         config.Cfg.Lang,                     //语言项
		"LangList":     config.Cfg.LangList,                 //语言列表
		"VersionNo":    constant.VERSION_NO,                 //当前版本号
		"BuildTime":    fileDate.Format("2006-01-02 15:04"), //打包时间
		"RootPath":     constant.ROOT_PATH,                  //当前项目根目录
		"UploadServer": server.Addr,                         //上传服务地址
		//"RomListPageSize": utils.ToString(db.ROM_PAGE_NUM),     //rom列表每页加载数量
		"logo": utils.WailsPathEncode(constant.RESOURCE_PATH+"logo.png", true), //侧边栏默认平台图片
	}

	return result, nil
}

// 读取配置原始信息
func GetBaseConfig() (db.ConfigVO, error) {
	m, _ := (&db.Config{}).Get()

	//软件名称
	softName := db.SoftName{}
	if m["SoftName"] != "" && m["SoftName"] != "{}" {
		_ = json.Unmarshal([]byte(m["SoftName"]), &softName)
	}

	if softName.Name == "" {
		softName.Name = "SimUI"
	}

	//开屏广告
	splashScreen := db.SplashScreen{}
	if m["SplashScreen"] != "" && m["SplashScreen"] != "{}" {
		_ = json.Unmarshal([]byte(m["SplashScreen"]), &splashScreen)
	}
	vo := db.ConfigVO{
		SoftName:           softName,
		EnableUpgrade:      uint8(utils.ToInt(m["EnableUpgrade"])),
		Lang:               m["Lang"],
		SearchEnginesBaidu: m["SearchEnginesBaidu"],
		Cursor:             m["Cursor"],
		SplashScreen:       splashScreen,
		AdminRunGame:       uint8(utils.ToInt(m["AdminRunGame"])),
		WindowZoom:         utils.ToFloat64(m["WindowZoom"]),
		GameMultiOpen:      uint8(utils.ToInt(m["GameMultiOpen"])),
	}

	return vo, nil
}

// 更新一个配置
func UpdateOneConfig(key, data string) error {
	(&db.Config{}).UpsertOne(key, data)
	return nil
}

func UpdateBaseConfig(data db.ConfigVO) error {

	if data.Cursor != "" {
		data.Cursor = utils.ToRelPath(data.Cursor, "")
	}
	if data.SoftName.Image != "" {
		data.SoftName.Image = utils.ToRelPath(data.SoftName.Image, "")
	}
	softName, _ := json.Marshal(data.SoftName)

	splashScreen, _ := json.Marshal(data.SplashScreen)

	m := &db.Config{}
	m.UpdateOne("Cursor", data.Cursor)
	m.UpdateOne("EnableUpgrade", data.EnableUpgrade)
	m.UpdateOne("Lang", data.Lang)
	m.UpdateOne("SearchEnginesBaidu", data.SearchEnginesBaidu)
	m.UpdateOne("SoftName", softName)
	m.UpdateOne("AdminRunGame", data.AdminRunGame)
	m.UpsertOne("WindowZoom", data.WindowZoom)
	m.UpsertOne("GameMultiOpen", data.GameMultiOpen)
	m.UpsertOne("SplashScreen", splashScreen)
	return nil
}

// 读取当前主题
func GetTheme() (string, error) {
	m, _ := (&db.Config{}).GetByName("Theme")
	return m, nil
}

// 更新当前主题
func SetTheme(theme string) error {
	err := (&db.Config{}).UpdateOne("Theme", theme)
	return err
}

// 更新展示图排序
func UpdateThumbsOrders(orders []string) error {
	data, _ := json.Marshal(orders)
	return (&db.Config{}).UpsertOne("ThumbOrders", string(data))
}

// 读取字体列表
func GetFontList() ([]map[string]any, error) {
	resp := []map[string]any{}

	//头部加入系统默认字体
	def := map[string]any{
		"label": config.Cfg.Lang["systemDefault"],
		"value": db.PlatformUIFont{},
	}
	resp = append(resp, def)

	user := components.GetUserFontList()
	system := components.GetSystemFontList()

	if len(user) == 0 && len(system) == 0 {
		return []map[string]any{}, nil
	}

	if len(user) > 0 {
		head := map[string]any{
			"label":   config.Cfg.Lang["userFont"],
			"value":   "",
			"disable": true,
		}
		resp = append(resp, head)

		for _, v := range user {
			f := map[string]any{
				"label": v.Family,
				"value": v,
			}
			resp = append(resp, f)
		}
	}

	if len(system) > 0 {
		head := map[string]any{
			"label":   config.Cfg.Lang["systemFont"],
			"value":   "",
			"disable": true,
		}
		resp = append(resp, head)

		for _, v := range system {
			f := map[string]any{
				"label": v.Family,
				"value": v,
			}
			resp = append(resp, f)
		}
	}

	return resp, nil
}
