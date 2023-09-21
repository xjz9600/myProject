package reflect

import (
	"fmt"
	"reflect"
)

func IterateArrayOrSlice(entity any) ([]any, error) {
	val := reflect.ValueOf(entity)
	var result []any
	for i := 0; i < val.Len(); i++ {
		result = append(result, val.Index(i).Interface())
	}
	return result, nil
}

func IterateMap(entity any) ([]any, []any, error) {
	val := reflect.ValueOf(entity)
	var resultKeys []any
	var resultValues []any
	mapRange := val.MapRange()
	for mapRange.Next() {
		resultKeys = append(resultKeys, mapRange.Key().Interface())
		resultValues = append(resultValues, mapRange.Value().Interface())
	}
	fmt.Println(resultKeys)
	fmt.Println(resultValues)
	return resultKeys, resultValues, nil
}
