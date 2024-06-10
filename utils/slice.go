package utils

import (
	"bytes"
	"math"
	"reflect"
)

/*
*
切片合并为字符串
*/
func SlicetoString(glue string, pieces []string) string {
	var buf bytes.Buffer
	l := len(pieces)
	for _, str := range pieces {
		buf.WriteString(str)
		if l--; l > 0 {
			buf.WriteString(glue)
		}
	}
	return buf.String()
}

/*
取差集
*/
func SliceDiff(slice1 []string, slice2 []string) []string {
	newSlice := []string{}
	isset := false
	for _, v := range slice1 {
		isset = false
		for _, v2 := range slice2 {
			if v2 == v {
				isset = true
				break
			}
		}
		if isset == false {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}

/*
取交集
*/
func SliceIntersect(list1 []string, list2 []string) []string {
	b2 := []string{}
	for _, v1 := range list1 {
		for _, v2 := range list2 {
			if v1 == v2 {
				b2 = append(b2, v1)
			}
		}
	}
	return b2
}

/*
删除最后一个元素
*/
func SliceDeleteLast(s []string) []string {
	if len(s) == 0 {
		return s
	}
	s = append(s[:len(s)-1])
	return s
}

// 去重
func SliceRemoveDuplicate(slc []string) []string {
	result := []string{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

// 去除空数据
func SliceRemoveEmpty(s []string) []string {
	if len(s) == 0 {
		return s
	}
	j := 0
	for _, v := range s {
		if v != "" {
			s[j] = v
			j++
		}
	}
	return s[:j]
}

// alice转map
func SliceToMap(sli []string) map[string]int {
	m := map[string]int{}

	if len(sli) == 0 {
		return m
	}
	for k, v := range sli {
		m[v] = k
	}
	return m
}

/**
 * 从struct/map中取出一列，类似于php的 array_column功能
 * SliceColumn[stuc,string](a,"Name")
 * @param arr 切片数据
 * @param k 要取出的字段名
 */
func SliceColumn[T, V any](arr []T, k string) []V {

	if arr == nil || len(arr) == 0 {
		return []V{}
	}

	values := make([]V, len(arr))
	switch reflect.TypeOf(arr).Elem().Kind() {
	case reflect.Map:
		for i, v := range arr {
			values[i] = reflect.ValueOf(v).MapIndex(reflect.ValueOf(k)).Interface().(V)
		}
		break
	case reflect.Struct:
		for i, v := range arr {
			values[i] = reflect.ValueOf(v).FieldByName(k).Interface().(V)
		}
		break
	}
	return values
}

/**
 * 检查一个数据是否在slice中
 * InSlice[string]("a", slice1)
 */
func InSlice[T comparable](s T, sli []T) bool {
	if len(sli) == 0 {
		return false
	}
	for _, v := range sli {
		if s == v {
			return true
		}
	}
	return false
}

// 根据数量拆分成多个切片
func SliceChunk[T comparable](s []T, size int) [][]T {
	if size < 1 {
		return [][]T{}
	}
	length := len(s)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]T
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		n = append(n, s[i*size:end])
		i++
	}
	return n
}
