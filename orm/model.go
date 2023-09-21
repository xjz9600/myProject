package orm

import (
	"myProject/orm/internal/errs"
	"reflect"
	"unicode"
)

type model struct {
	tableName string
	fieldMap  map[string]*field
}

type field struct {
	colName string
}

func parseModel(val any) (*model, error) {
	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	typ = typ.Elem()
	numFields := typ.NumField()
	res := make(map[string]*field)
	for i := 0; i < numFields; i++ {
		fd := typ.Field(i)
		res[fd.Name] = &field{
			colName: underscoreName(fd.Name),
		}
	}
	return &model{
		tableName: underscoreName(typ.Name()),
		fieldMap:  res,
	}, nil
}

// underscoreName 驼峰转字符串命名
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}
