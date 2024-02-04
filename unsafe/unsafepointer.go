package unsafe

import (
	"errors"
	"reflect"
	"unsafe"
)

func (u *UnsafeAccessor) SetField(fieldName string, val any) error {
	field, ok := u.fields[fieldName]
	if !ok {
		return errors.New("未知字段")
	}
	ptr := unsafe.Pointer(uintptr(u.address) + field.offset)
	reflect.NewAt(field.typ, ptr).Elem().Set(reflect.ValueOf(val))
	return nil
}

func GetField(entity any) *UnsafeAccessor {
	typ := reflect.TypeOf(entity)
	typ = typ.Elem()
	numField := typ.NumField()
	fields := make(map[string]FieldMeta, numField)
	for i := 0; i < numField; i++ {
		fd := typ.Field(i)
		fields[fd.Name] = FieldMeta{
			offset: fd.Offset,
			typ:    fd.Type,
		}
	}
	val := reflect.ValueOf(entity)
	return &UnsafeAccessor{
		address: val.UnsafePointer(),
		fields:  fields,
	}
}

type UnsafeAccessor struct {
	// 结构体起始地址
	address unsafe.Pointer
	fields  map[string]FieldMeta
}

type FieldMeta struct {
	// 相对偏移量
	offset uintptr
	typ    reflect.Type
}
