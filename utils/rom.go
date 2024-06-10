package utils

import (
	"strings"
)

var prefix = "t0Ic" //前缀

// 创建rom唯一id
func CreateRomUniqId(t string, s int64) string {
	return prefix + "_" + t + "_" + ToString(s)
}

// 分割rom唯一id
func SplitRomUniqId(uniqId string) (string, string) {
	if uniqId == "" {
		return "", ""
	}
	sli := strings.Split(uniqId, "_")
	return sli[1], sli[2] //返回 时间,大小
}

// 检查是不是唯一ID
func HasRomUniqId(uniqId string) bool {
	if uniqId == "" {
		return false
	}
	sli := strings.Split(uniqId, "_")

	if sli[0] == prefix {
		return true
	}
	return false
}

// 读取rom的目录
func getSplitRomMenu(rel string) string {
	p := GetFilePath(rel)
	parr := strings.Split(p, "/")
	parr = SliceRemoveEmpty(parr)
	menu := "_7b9"
	if len(parr) > 0 {
		menu = menu + "/" + parr[0]
	}
	return menu
}

// 处理rom路径为2级
func HandleRomMenu(p string) string {
	rel := GetFilePath(p)
	arr := strings.Split(rel, "/")
	if len(arr) > 3 {
		arr = arr[0:3]
		rel = strings.Join(arr, "/") + "/"
	}
	return rel
}

// 读取map中第一个key和value
func MapGetFirst[K, V comparable](mp map[K]V) (K, V) {
	var key K
	var val V

	if len(mp) == 0 {
		return key, val
	}

	for k, v := range mp {
		return k, v
	}
	return key, val
}
