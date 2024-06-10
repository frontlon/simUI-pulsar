package modules

import (
	"simUI/config"
	"simUI/db"
)

// 这些是必须返回出一些别名配置
var rombaseAliasList = []string{"OtherA", "OtherB", "OtherC", "OtherD"}

// 根据类型读取数据
func GetRomBaseAliasByPlatform(platform uint32) (map[string]string, error) {
	voMap, err := (&db.RombaseAlias{}).GetByPlatform(platform)
	if err != nil {
		voMap = map[string]string{}
	}

	for _, typ := range rombaseAliasList {
		if _, ok := voMap[typ]; !ok {
			voMap[typ] = config.Cfg.Lang["base"+typ]
		}
	}

	return voMap, nil

}

// 更新数据
func UpdateRomBaseAlias(platform uint32, data map[string]string) error {

	existMap, _ := (&db.RombaseAlias{}).GetByPlatform(platform)

	for typ, alias := range data {

		//删除记录
		if alias == "" {
			(&db.RombaseAlias{
				Platform: platform,
				Type:     typ,
			}).DeleteByType()
			continue
		}

		if _, ok := existMap[typ]; ok {
			//修改记录
			(&db.RombaseAlias{
				Platform: platform,
				Type:     typ,
				Alias:    alias,
			}).UpdateByType()
		} else {
			//新增记录
			(&db.RombaseAlias{
				Platform: platform,
				Type:     typ,
				Alias:    alias,
			}).Add()
		}
	}
	return nil
}
