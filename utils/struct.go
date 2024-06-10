package utils

import (
	"errors"
	"reflect"
)

// 根据Key，读取struct value
// T是值或指针都可以
func GetStructValue[T, V comparable](stu T, k string) V {
	var value V
	switch reflect.TypeOf(stu).Kind() {
	case reflect.Struct:
		value = reflect.ValueOf(stu).FieldByName(k).Interface().(V)
	case reflect.Ptr:
		isset := reflect.ValueOf(stu).Elem().FieldByName(k)
		if isset.IsValid() {
			value = reflect.ValueOf(stu).Elem().FieldByName(k).Interface().(V)
		}
	}
	return value
}

// 根据Key，设置struct value
// T 必须是指针
func SetStructValue[T comparable](stu T, k string, val any) error {

	structValue := reflect.ValueOf(stu).Elem()
	field := structValue.FieldByName(k)

	if !field.IsValid() {
		return errors.New("no such field")
	}

	if !field.CanSet() {
		return errors.New("cannot set this field")
	}

	value := ToString(val)
	field.Set(reflect.ValueOf(value))
	return nil
}
