package value

import (
	"database/sql"
	"myProject/orm/internal/errs"
	"myProject/orm/model"
	"reflect"
)

type reflectValue struct {
	model  *model.Model
	entity reflect.Value
}

var _ Creator = NewReflectValue
var _ Value = reflectValue{}

func NewReflectValue(model *model.Model, entity any) Value {
	return reflectValue{
		model:  model,
		entity: reflect.ValueOf(entity).Elem(),
	}
}

func (r reflectValue) Field(name string) any {
	return r.entity.FieldByName(name).Interface()
}

func (r reflectValue) SetColumn(rows *sql.Rows) error {
	col, err := rows.Columns()
	if err != nil {
		return err
	}
	var data []any
	var dataVal []reflect.Value
	for _, c := range col {
		if fd, ok := r.model.ColumnMap[c]; ok {
			val := reflect.New(fd.Typ)
			data = append(data, val.Interface())
			dataVal = append(dataVal, val.Elem())
			continue
		}
		return errs.NewErrUnknownField(c)
	}
	err = rows.Scan(data...)
	if err != nil {
		return err
	}
	for i, c := range col {
		if fd, ok := r.model.ColumnMap[c]; ok {
			r.entity.FieldByName(fd.GoName).Set(dataVal[i])
			continue
		}
		return errs.NewErrUnknownField(c)
	}
	return nil
}
