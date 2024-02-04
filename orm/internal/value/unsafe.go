package value

import (
	"database/sql"
	"myProject/orm/internal/errs"
	"myProject/orm/model"
	"reflect"
	"unsafe"
)

type unsafeValue struct {
	model  *model.Model
	entity unsafe.Pointer
}

var _ Creator = NewUnsafeValue
var _ Value = unsafeValue{}

func NewUnsafeValue(model *model.Model, entity any) Value {
	return unsafeValue{
		model:  model,
		entity: reflect.ValueOf(entity).UnsafePointer(),
	}
}

func (u unsafeValue) Field(name string) any {
	fd := u.model.FieldMap[name]
	return reflect.NewAt(fd.Typ, unsafe.Pointer(uintptr(u.entity)+fd.Offset)).Elem().Interface()
}

func (u unsafeValue) SetColumn(rows *sql.Rows) error {
	col, err := rows.Columns()
	if err != nil {
		return err
	}
	var data []any
	for _, c := range col {
		if fd, ok := u.model.ColumnMap[c]; ok {
			val := reflect.NewAt(fd.Typ, unsafe.Pointer(uintptr(u.entity)+fd.Offset))
			data = append(data, val.Interface())
			continue
		}
		return errs.NewErrUnknownField(c)
	}
	err = rows.Scan(data...)
	if err != nil {
		return err
	}
	return nil
}
