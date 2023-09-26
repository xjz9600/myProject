package orm

import (
	"myProject/orm/internal/errs"
	"reflect"
	"strings"
	"unicode"
)

type modelOpt func(*model) error

const (
	tagKeyColumn = "column"
)

type TableName interface {
	TableName() string
}

type model struct {
	tableName string
	fieldMap  map[string]*field
}

type field struct {
	colName string
}

func parseTag(tag reflect.StructTag) (map[string]string, error) {
	segStr, ok := tag.Lookup(`orm`)
	if !ok {
		return map[string]string{}, nil
	}
	res := make(map[string]string)
	segs := strings.Split(segStr, ",")
	for _, seg := range segs {
		pairs := strings.Split(seg, "=")
		if len(pairs) != 2 {
			return nil, errs.NewErrInvalidTagContent(seg)
		}
		key := pairs[0]
		val := pairs[1]
		res[key] = val
	}
	return res, nil
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
		tags, err := parseTag(fd.Tag)
		if err != nil {
			return nil, err
		}
		if _, ok := tags[tagKeyColumn]; ok {
			res[fd.Name] = &field{
				colName: tags[tagKeyColumn],
			}
		} else {
			res[fd.Name] = &field{
				colName: underscoreName(fd.Name),
			}
		}
	}
	var tableName string
	if tn, ok := val.(TableName); ok {
		tableName = tn.TableName()
	} else {
		tableName = underscoreName(typ.Name())
	}
	md := &model{tableName: tableName, fieldMap: res}
	return md, nil
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
