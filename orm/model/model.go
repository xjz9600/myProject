package model

import (
	"myProject/orm/internal/errs"
	"reflect"
	"strings"
	"unicode"
)

type modelOpt func(*Model) error

const (
	tagKeyColumn = "column"
)

type TableName interface {
	TableName() string
}

type Model struct {
	TableName string
	FieldMap  map[string]*Field
	ColumnMap map[string]*Field
	Fields    []*Field
}

type Field struct {
	ColName string
	GoName  string
	Typ     reflect.Type
	Offset  uintptr
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

func ParseModel(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	typ = typ.Elem()
	numFields := typ.NumField()
	fieldMap := make(map[string]*Field)
	columnMap := make(map[string]*Field)
	fields := make([]*Field, 0, numFields)
	for i := 0; i < numFields; i++ {
		fd := typ.Field(i)
		tags, err := parseTag(fd.Tag)
		if err != nil {
			return nil, err
		}
		f := &Field{
			Typ:     fd.Type,
			GoName:  fd.Name,
			ColName: underscoreName(fd.Name),
			Offset:  fd.Offset,
		}
		if _, ok := tags[tagKeyColumn]; ok {
			f.ColName = tags[tagKeyColumn]
		}
		fieldMap[f.GoName] = f
		columnMap[f.ColName] = f
		fields = append(fields, f)
	}
	var tableName string
	if tn, ok := val.(TableName); ok {
		tableName = tn.TableName()
	} else {
		tableName = underscoreName(typ.Name())
	}
	md := &Model{TableName: tableName, FieldMap: fieldMap, ColumnMap: columnMap, Fields: fields}
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
