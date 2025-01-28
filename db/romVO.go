package db

import (
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
)

type RomVO struct {
}

type RomSimpleVO struct {
	Id            uint64
	Platform      uint32  // 平台ID
	Menu          string  // 菜单
	Name          string  // 中文名
	RomName       string  // rom名称
	RomPath       string  // rom路径
	Hide          uint8   // 展示图路径
	Score         float64 // 评分
	Favorite      uint8   // 我的最爱
	ThumbPic      string  // 展示图路径
	BaseNameEn    string  // 英文名
	BaseNameJp    string  // 日文名
	BaseType      string  // 游戏类型，如RPG
	BaseYear      string  // 游戏年份
	BaseProducer  string  // 游戏出品公司
	BasePublisher string  // 游戏出品公司
	BaseCountry   string  // 游戏国家
	BaseTranslate string  // 汉化组
	BaseVersion   string  // 版本
	BaseOtherA    string  // 附加信息A
	BaseOtherB    string  // 附加信息B
	BaseOtherC    string  // 附加信息C
	BaseOtherD    string  // 附加信息D
	SimId         uint32  // 附加信息D
	SubGames      []*RomSimpleVO
	SimSetting    map[string]*RomSimSetting
}

type RomSimplestVO struct {
	Id       uint64
	Platform uint32 // 平台ID
	Name     string // 中文名
	RomName  string // rom名称
	ThumbPic string //展示图
}

// rom模拟器设置
type RomSimSetting struct {
	Cmd       string // 启动命令参数
	RunBefore string // 游戏运行前执行脚本
	RunAfter  string // 游戏运行后执行脚本
	Unzip     string // 解压后运行
}

func (m *Rom) ConvertRomListSimple(lists []*Rom, showSubGame uint8) []*RomSimpleVO {
	result := []*RomSimpleVO{}
	if len(lists) == 0 {
		return result
	}

	//加载子游戏
	subRoms := map[string][]*Rom{}
	if showSubGame == 1 {
		names := []string{}
		for _, v := range lists {
			names = append(names, v.RomName)
		}
		subRoms, _ = m.GetSubRomByFileNames(names)
	}

	for _, v := range lists {
		s := []*Rom{}
		if _, ok := subRoms[v.LinkFile]; ok {
			s = subRoms[v.LinkFile]
		}
		result = append(result, m.ConvertRomSimple(v, s))
	}

	return result
}

func (m *Rom) ConvertRomListSimplest(lists []*Rom) []*RomSimplestVO {
	result := []*RomSimplestVO{}
	if len(lists) == 0 {
		return result
	}

	for _, v := range lists {
		s := &RomSimplestVO{
			Id:       v.Id,
			Platform: v.Platform,
			Name:     v.Name,
			RomName:  v.RomName,
			ThumbPic: "",
		}
		result = append(result, s)
	}

	return result
}

func (m *Rom) ConvertRomSimple(rom *Rom, subRoms []*Rom) *RomSimpleVO {

	//rom模拟器配置
	simSetting := map[string]*RomSimSetting{}
	if rom.SimSetting != "" {
		json.Unmarshal([]byte(rom.SimSetting), &simSetting)
	}

	r := &RomSimpleVO{
		Id:            rom.Id,
		Name:          rom.Name,
		Platform:      rom.Platform,
		Menu:          rom.Menu,
		RomName:       rom.RomName,
		Hide:          rom.Hide,
		Score:         rom.Score,
		Favorite:      rom.Favorite,
		BaseNameEn:    rom.BaseNameEn,
		BaseNameJp:    rom.BaseNameJp,
		BaseType:      rom.BaseType,
		BaseYear:      rom.BaseYear,
		BaseProducer:  rom.BaseProducer,
		BasePublisher: rom.BasePublisher,
		BaseCountry:   rom.BaseCountry,
		BaseTranslate: rom.BaseTranslate,
		BaseVersion:   rom.BaseVersion,
		BaseOtherA:    rom.BaseOtherA,
		BaseOtherB:    rom.BaseOtherB,
		BaseOtherC:    rom.BaseOtherC,
		BaseOtherD:    rom.BaseOtherD,
		SimId:         rom.SimId,
		SimSetting:    simSetting,
	}

	//加载子游戏
	if len(subRoms) > 0 {
		subRomVOList := []*RomSimpleVO{}
		for _, b := range subRoms {
			subSimSetting := map[string]*RomSimSetting{}
			if b.SimSetting != "" {
				json.Unmarshal([]byte(b.SimSetting), &subSimSetting)
			}

			subRomVO := &RomSimpleVO{
				Id:            b.Id,
				Platform:      b.Platform,
				Name:          b.Name,
				BaseNameEn:    b.BaseNameEn,
				BaseNameJp:    b.BaseNameJp,
				BaseType:      b.BaseType,
				BaseYear:      b.BaseYear,
				BaseProducer:  b.BaseProducer,
				BasePublisher: b.BasePublisher,
				BaseCountry:   b.BaseCountry,
				BaseTranslate: b.BaseTranslate,
				BaseVersion:   b.BaseVersion,
				SimSetting:    subSimSetting,
			}
			subRomVOList = append(subRomVOList, subRomVO)
		}
		r.SubGames = subRomVOList
	}
	return r

}

/*func (m *Rom) ConvertRomSimplest(lists []*Rom, showSubGame uint8) []*RomSimplestVO {
	result := []*RomSimplestVO{}
	if len(lists) == 0 {
		return result
	}

	for _, v := range lists {
		r := &RomSimplestVO{
			Id:       v.Id,
			Platform: v.Platform,
			Name:     v.Name,
			RomName:  v.RomName,
		}
		result = append(result, r)
	}

	//加载子游戏
	if showSubGame == 1 {
		names := []string{}
		for _, v := range lists {
			names = append(names, v.RomName)
		}
		subRoms, _ := m.GetSubRomByFileNames(names)
		for k, v := range lists {
			if _, ok := subRoms[v.LinkFile]; !ok {
				continue
			}
			subRomVOList := []*RomSimplestVO{}
			for _, b := range subRoms[v.LinkFile] {
				subRomVO := &RomSimplestVO{
					Id:       b.Id,
					Platform: b.Platform,
					Name:     b.Name,
					ThumbPic: "",
					SubGames: nil,
				}
				subRomVOList = append(subRomVOList, subRomVO)
			}
			result[k].SubGames = subRomVOList
		}
	}

	return result
}*/
