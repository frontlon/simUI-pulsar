package config

import (
	"encoding/csv"
	"fmt"
	"os"
	"simUI/db"
	"simUI/utils"
	"strings"
)

type RomSetting struct {
	RomName     string // rom文件名
	Complete    string // 通关状态(0未通关;1已通关;2完美通关)
	Hide        string // 是否隐藏
	RunNum      string // 运行次数
	RunLasttime string // 最后运行时间
	SimId       string // 正在使用的模拟器id
	SimSetting  string //ROM模拟器设置
	Favorite    string //是否为喜爱
}

var RomSettingList map[uint32]map[string]*RomSetting

// 读取游戏资料列表
func GetRomSettingByPlatform(platform uint32) (map[string]*RomSetting, error) {

	if len(RomSettingList) == 0 {
		RomSettingList = map[uint32]map[string]*RomSetting{}
	}

	//如果已经读取，则直接返回
	if _, ok := RomSettingList[platform]; ok {
		return RomSettingList[platform], nil
	}

	RomSettingList[platform] = map[string]*RomSetting{}
	platformInfo := (&db.Platform{}).GetVOById(platform, false)

	if !utils.FileExists(platformInfo.RomSettingFile) {
		_ = CreateNewRomSettingFile(platformInfo.RomSettingFile)
	}

	records, err := utils.ReadCsv(platformInfo.RomSettingFile)
	if err != nil {
		return RomSettingList[platform], nil //直接返回空，不返回错误
	}

	isUtf8 := false
	if len(records) > 0 {
		isUtf8 = utils.IsUTF8(records[0][0])
	} else {
		return RomSettingList[platform], nil
	}

	creates := [][]string{}
	for k, r := range records {

		if k == 0 {
			continue
		}

		create := []string{"", "", "", "", "", "", "", ""}
		createLen := len(create)
		i := 1
		for a, b := range r {
			if i > createLen {
				continue
			}
			i++
			create[a] = b
		}

		//转换成utf-8编码
		if isUtf8 == false {
			for ck, cv := range create {
				create[ck] = utils.ToUTF8(cv)
			}
		}
		creates = append(creates, create)
	}

	for _, create := range creates {
		if create[0] == "" {
			continue
		}
		lists := &RomSetting{}
		if len(create) >= 1 {
			lists.RomName = strings.TrimSpace(create[0])
		}

		if len(create) >= 2 {
			lists.Complete = strings.TrimSpace(create[1])
		}
		if len(create) >= 3 {
			lists.Hide = strings.TrimSpace(create[2])
		}
		if len(create) >= 4 {
			lists.RunNum = strings.TrimSpace(create[3])
		}
		if len(create) >= 5 {
			lists.RunLasttime = strings.TrimSpace(create[4])
		}
		if len(create) >= 6 {
			lists.SimId = strings.TrimSpace(create[5])
		}
		if len(create) >= 7 {
			lists.SimSetting = strings.TrimSpace(create[6])
		}
		if len(create) >= 8 {
			lists.Favorite = strings.TrimSpace(create[7])
		}

		RomSettingList[platform][create[0]] = lists
	}
	return RomSettingList[platform], nil
}

func GetRomSettingByName(platform uint32, name string) *RomSetting {

	if len(RomSettingList) == 0 || len(RomSettingList[platform]) == 0 {
		_, _ = GetRomSettingByPlatform(platform)
	}

	if _, ok := RomSettingList[platform][name]; ok {
		return RomSettingList[platform][name]
	}
	return &RomSetting{}
}

func SetRomSettingOneField(platform uint32, name string, key string, val any, flush bool) error {
	if len(RomSettingList) == 0 {
		_, _ = GetRomSettingByPlatform(platform)
	}
	if len(RomSettingList[platform]) == 0 {
		RomSettingList[platform] = map[string]*RomSetting{}
	}

	if _, ok := RomSettingList[platform][name]; !ok {
		RomSettingList[platform][name] = &RomSetting{
			RomName: name,
		}
	}

	if err := utils.SetStructValue[*RomSetting](RomSettingList[platform][name], key, val); err != nil {
		return err
	}

	if flush {
		FlushRomSetting(platform)
	}

	return nil
}

// 写csv文件
func FlushRomSetting(platform uint32) {

	platformInfo := (&db.Platform{}).GetVOById(platform, false)

	go func() {
		if err := WriteDataToRomSettingFile(platformInfo.RomSettingFile, RomSettingList[platform]); err != nil {
			fmt.Println("FlushRomSetting Error:", err)
		}
	}()
}

// 数据覆写csv文件
func CoverRomSettingFile(platform uint32, newData map[string]*RomSetting) error {

	platformInfo := (&db.Platform{}).GetVOById(platform, false)

	//写入csv文件
	if err := WriteDataToRomSettingFile(platformInfo.RomSettingFile, newData); err != nil {
		return err
	}

	//更新全局变量
	RomSettingList[platform] = newData

	return nil
}

// 创建一个新的csv文件
func CreateNewRomSettingFile(p string) error {

	//转换为切片
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，避免使用Microsoft Excel打开乱码

	writer := csv.NewWriter(f)

	//表头
	writer.Write(getRomSettingTitle())
	writer.Flush() // 此时才会将缓冲区数据写入
	return nil
}

// 写入csv文件
func WriteDataToRomSettingFile(filePath string, data map[string]*RomSetting) error {

	if filePath == "" {
		return nil
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，避免使用Microsoft Excel打开乱码

	writer := csv.NewWriter(f)

	//表头
	writer.Write(getRomSettingTitle())

	for _, v := range data {

		str := strings.Join([]string{v.Complete, v.Hide, v.RunNum, v.RunLasttime, v.SimId, v.SimSetting, v.Favorite}, "")

		if str == "" {
			continue
		}

		writer.Write([]string{
			v.RomName,
			strings.TrimSpace(v.Complete),
			strings.TrimSpace(v.Hide),
			strings.TrimSpace(v.RunNum),
			strings.TrimSpace(v.RunLasttime),
			strings.TrimSpace(v.SimId),
			strings.TrimSpace(v.SimSetting),
			strings.TrimSpace(v.Favorite),
		})

	}

	writer.Flush() // 此时才会将缓冲区数据写入

	return nil
}

// 表头
func getRomSettingTitle() []string {
	return []string{
		"romName",     // rom文件路径
		"complete",    // 通关状态(0未通关;1已通关;2完美通关)
		"hide",        // 是否隐藏
		"runNum",      // 运行次数
		"runLastTime", // 最后运行时间
		"simId",       // 正在使用的模拟器id
		"simSetting",  // 模拟器配置
		"favorite",    // 是否为喜爱
	}
}
