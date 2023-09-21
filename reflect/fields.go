package reflect

import (
	"errors"
	"reflect"
)

func IterateFields(entity any) (map[string]any, error) {
	if entity == nil {
		return nil, errors.New("不支持 nil")
	}
	typ := reflect.TypeOf(entity)
	val := reflect.ValueOf(entity)
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, errors.New("不支持的类型")
	}
	num := typ.NumField()
	res := make(map[string]any, num)
	for i := 0; i < num; i++ {
		if val.Field(i).CanInterface() {
			res[typ.Field(i).Name] = val.Field(i).Interface()
		} else {
			res[typ.Field(i).Name] = reflect.Zero(val.Field(i).Type()).Interface()
		}
	}
	return res, nil
}

func SetField(entity any, field string, newValue any) error {
	val := reflect.ValueOf(entity)
	for val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return errors.New("不可修改字段")
	}
	fieldVal := val.FieldByName(field)
	if !fieldVal.CanSet() {
		return errors.New("不可修改字段")
	}
	fieldVal.Set(reflect.ValueOf(newValue))
	return nil
}
