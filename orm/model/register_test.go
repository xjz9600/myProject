package model

import (
	"github.com/stretchr/testify/assert"
	"myProject/orm/internal/errs"
	"reflect"
	"testing"
)

func TestRegister(t *testing.T) {
	r := NewSyncMapRegister()
	m, err := r.Register(&TableNameTest{}, WithTableName("test"), WithColumnName("Name", "test_name_tt"))
	assert.NoError(t, err)
	wantModel := &Model{
		TableName: "test",
		FieldMap: map[string]*Field{
			"Name": &Field{
				ColName: "test_name_tt",
				GoName:  "Name",
				Typ:     reflect.TypeOf(""),
				Offset:  0,
			},
		},
		ColumnMap: map[string]*Field{
			"test_name_tt": &Field{
				ColName: "test_name_tt",
				GoName:  "Name",
				Typ:     reflect.TypeOf(""),
				Offset:  0,
			},
		},
		Fields: []*Field{
			&Field{
				ColName: "test_name_tt",
				GoName:  "Name",
				Typ:     reflect.TypeOf(""),
				Offset:  0,
			},
		},
	}
	assert.Equal(t, m, wantModel)

	re := NewSyncRegister()
	_, err = re.Register(&TableNameTest{}, WithTableName("test"), WithColumnName("XXX", "test_name"))
	assert.Equal(t, err, errs.NewErrUnknownField("XXX"))
}
