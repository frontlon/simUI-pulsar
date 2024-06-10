package modules

import (
	"simUI/db"
	"strings"
)

// 读取rombase枚举
func UpdateRomBaseEnum(t string, data []string) error {

	//先删除记录
	(&db.RombaseEnum{Type: t}).DeleteByType()

	if len(data) == 0 {
		return nil
	}

	create := []*db.RombaseEnum{}
	for k, v := range data {
		if strings.Trim(v, " ") == "" {
			continue
		}
		c := &db.RombaseEnum{}
		c.Type = t
		c.Name = v
		c.Sort = uint32(k) + 1
		create = append(create, c)
	}
	if err := (&db.RombaseEnum{}).BatchAdd(create); err != nil {
		return err
	}
	return nil
}

// 读取全部枚举数据
func GetRomBaseEnum() (map[string][]string, error) {
	lists, err := (&db.RombaseEnum{}).GetAll()
	if err != nil {
		return nil, err
	}
	result := map[string][]string{}

	result["type"] = []string{}
	result["year"] = []string{}
	result["producer"] = []string{}
	result["publisher"] = []string{}
	result["country"] = []string{}
	result["translate"] = []string{}
	result["version"] = []string{}

	for _, v := range lists {
		result[v.Type] = append(result[v.Type], v.Name)
	}
	return result, nil
}

// 根据类型读取枚举数据
func GetRomBaseEnumByType(t string) ([]string, error) {
	lists, err := (&db.RombaseEnum{}).GetByType(t)
	if err != nil {
		return nil, err
	}
	result := []string{}
	for _, v := range lists {
		result = append(result, v.Name)
	}
	return result, nil
}
