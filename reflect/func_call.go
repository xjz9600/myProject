package reflect

import (
	"reflect"
)

func IterateFunc(entity any) (map[string]FuncInfo, error) {
	tpy := reflect.TypeOf(entity)
	numMethod := tpy.NumMethod()
	var result map[string]FuncInfo
	for i := 0; i < numMethod; i++ {
		fn := tpy.Method(i)
		inputNum := fn.Type.NumIn()
		var inputType []reflect.Type
		var callInput []reflect.Value
		callInput = append(callInput, reflect.ValueOf(entity))
		inputType = append(inputType, tpy)
		for i := 1; i < inputNum; i++ {
			inputType = append(inputType, fn.Type.In(i))
			callInput = append(callInput, reflect.Zero(fn.Type.In(i)))
		}
		outPutNum := fn.Type.NumOut()
		var outputType []reflect.Type
		for i := 0; i < outPutNum; i++ {
			outputType = append(outputType, fn.Type.Out(i))
		}
		resValues := fn.Func.Call(callInput)
		if result == nil {
			result = map[string]FuncInfo{}
		}
		var res []any
		for _, r := range resValues {
			res = append(res, r.Interface())
		}
		result[fn.Name] = FuncInfo{
			Name:   fn.Name,
			Input:  inputType,
			Output: outputType,
			Result: res,
		}
	}
	return result, nil
}

type FuncInfo struct {
	Name   string
	Input  []reflect.Type
	Output []reflect.Type
	Result []any
}
