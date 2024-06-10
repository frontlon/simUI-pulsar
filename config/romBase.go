package config

import (
	"encoding/csv"
	"os"
	"simUI/db"
	"simUI/utils"
	"strings"
)

type RomBase struct {
	RomName   string // rom文件名
	Name      string // 中文名
	NameEN    string // 英文名
	NameJP    string // 日文名
	Type      string // 类型
	Year      string // 年份
	Producer  string // 制造商
	Publisher string // 出品公司
	Country   string // 国家
	Translate string // 汉化组
	Version   string // 版本
	OtherA    string // 其他内容a
	OtherB    string // 其他内容b
	OtherC    string // 其他内容c
	OtherD    string // 其他内容d
	Score     string // 评分
}

var RomBaseList map[uint32]map[string]*RomBase

// 读取游戏资料列表
func GetRomBase(platform uint32, rebuild bool) (map[string]*RomBase, error) {

	if rebuild || len(RomBaseList) == 0 {
		RomBaseList = map[uint32]map[string]*RomBase{}
	}

	//如果已经读取，则直接返回
	if _, ok := RomBaseList[platform]; ok {
		return RomBaseList[platform], nil
	}

	//开始整理数据
	RomBaseList[platform] = map[string]*RomBase{}

	platformInfo := (&db.Platform{}).GetVOById(platform, false)

	if !utils.FileExists(platformInfo.RomBaseFile) {
		_ = CreateNewRomBaseFile(platformInfo.RomBaseFile)
	}

	records, err := utils.ReadCsv(platformInfo.RomBaseFile)
	if err != nil {
		return RomBaseList[platform], nil //直接返回空，不返回错误
	}

	isUtf8 := false
	if len(records) > 0 {
		isUtf8 = utils.IsUTF8(records[0][0])
	} else {
		return RomBaseList[platform], nil
	}

	creates := [][]string{}
	for k, r := range records {

		if k == 0 {
			continue
		}

		create := []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""}
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
		lists := &RomBase{}
		if len(create) >= 1 {
			lists.RomName = strings.TrimSpace(create[0])
		}
		if len(create) >= 2 {
			lists.Name = strings.TrimSpace(create[1])
		}
		if len(create) >= 3 {
			lists.Type = strings.TrimSpace(create[2])
		}
		if len(create) >= 4 {
			lists.Year = strings.TrimSpace(create[3])
		}
		if len(create) >= 5 {
			lists.Publisher = strings.TrimSpace(create[4])
		}
		if len(create) >= 6 {
			lists.Country = strings.TrimSpace(create[5])
		}
		if len(create) >= 7 {
			lists.Translate = strings.TrimSpace(create[6])
		}
		if len(create) >= 8 {
			lists.Version = strings.TrimSpace(create[7])
		}
		if len(create) >= 9 {
			lists.Producer = strings.TrimSpace(create[8])
		}
		if len(create) >= 10 {
			lists.NameEN = strings.TrimSpace(create[9])
		}
		if len(create) >= 11 {
			lists.NameJP = strings.TrimSpace(create[10])
		}
		if len(create) >= 12 {
			lists.OtherA = strings.TrimSpace(create[11])
		}
		if len(create) >= 13 {
			lists.OtherB = strings.TrimSpace(create[12])
		}
		if len(create) >= 14 {
			lists.OtherC = strings.TrimSpace(create[13])
		}
		if len(create) >= 15 {
			lists.OtherD = strings.TrimSpace(create[14])
		}
		if len(create) >= 16 {
			lists.Score = strings.TrimSpace(create[15])
		}
		RomBaseList[platform][create[0]] = lists
	}
	return RomBaseList[platform], nil
}

// 读取一个rom的资料信息
func GetRomBaseById(platform uint32, id string) *RomBase {
	romlist, _ := GetRomBase(platform, false)
	if _, ok := romlist[id]; ok {
		return RomBaseList[platform][id]
	}
	return nil
}

// 写csv文件
func AddRomBase(platform uint32, newData *RomBase) error {

	platformInfo := (&db.Platform{}).GetVOById(platform, false)

	info, _ := GetRomBase(platform, false) //读取老数据
	//如果全为空则删除当前记录
	info[newData.RomName] = newData //并入新数据

	//写入csv文件
	if err := WriteDataToRomBaseFile(platformInfo.RomBaseFile, info); err != nil {
		return err
	}

	//更新全局变量
	RomBaseList[platform][newData.RomName] = newData

	return nil
}

// 覆盖csv文件
func CoverRomBaseFile(platform uint32, newData map[string]*RomBase) error {

	platformInfo := (&db.Platform{}).GetVOById(platform, false)

	//写入csv文件
	if err := WriteDataToRomBaseFile(platformInfo.RomBaseFile, newData); err != nil {
		return err
	}

	//更新全局变量
	RomBaseList[platform] = map[string]*RomBase{}
	RomBaseList[platform] = newData

	return nil
}

// 写入csv文件
func WriteDataToRomBaseFile(filePath string, data map[string]*RomBase) error {

	if filePath == "" {
		return nil
	}

	//转换为切片
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，避免使用Microsoft Excel打开乱码

	writer := csv.NewWriter(f)

	//表头
	writer.Write(getRomBaseTitle())

	for _, v := range data {

		str := strings.Join([]string{v.Name, v.Producer, v.Publisher, v.Year, v.Type, v.Country, v.Translate, v.Version, v.NameEN, v.NameJP, v.OtherA, v.OtherB, v.OtherC, v.OtherD, v.Score}, "")

		if str == "" {
			continue
		}

		writer.Write([]string{
			v.RomName,
			strings.TrimSpace(v.Name),
			strings.TrimSpace(v.Type),
			strings.TrimSpace(v.Year),
			strings.TrimSpace(v.Publisher),
			strings.TrimSpace(v.Country),
			strings.TrimSpace(v.Translate),
			strings.TrimSpace(v.Version),
			strings.TrimSpace(v.Producer),
			strings.TrimSpace(v.NameEN),
			strings.TrimSpace(v.NameJP),
			strings.TrimSpace(v.OtherA),
			strings.TrimSpace(v.OtherB),
			strings.TrimSpace(v.OtherC),
			strings.TrimSpace(v.OtherD),
			strings.TrimSpace(v.Score),
		})

	}
	writer.Flush() // 此时才会将缓冲区数据写入

	return nil
}

// 创建一个新的csv文件
func CreateNewRomBaseFile(p string) error {

	//转换为切片
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，避免使用Microsoft Excel打开乱码

	writer := csv.NewWriter(f)

	//表头
	writer.Write(getRomBaseTitle())
	writer.Flush() // 此时才会将缓冲区数据写入
	return nil
}

// 表头
func getRomBaseTitle() []string {
	return []string{
		Cfg.Lang["romName"],
		Cfg.Lang["alias"],
		Cfg.Lang["baseType"],
		Cfg.Lang["baseYear"],
		Cfg.Lang["basePublisher"],
		Cfg.Lang["baseCountry"],
		Cfg.Lang["baseTranslate"],
		Cfg.Lang["baseVersion"],
		Cfg.Lang["baseProducer"],
		Cfg.Lang["baseNameEn"],
		Cfg.Lang["baseNameJp"],
		Cfg.Lang["baseOther"] + "A",
		Cfg.Lang["baseOther"] + "B",
		Cfg.Lang["baseOther"] + "C",
		Cfg.Lang["baseOther"] + "D",
		Cfg.Lang["rating"],
	}
}
