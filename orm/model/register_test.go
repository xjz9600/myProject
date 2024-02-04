package orm

import (
	"github.com/stretchr/testify/assert"
	"myProject/orm/internal/errs"
	"myProject/orm/model"
	"reflect"
	"testing"
)

func TestRegister(t *testing.T) {
	r := NewSyncMapRegister()
	m, err := r.Register(&model.TableNameTest{}, WithTableName("test"), WithColumnName("Name", "test_name_tt"))
	assert.NoError(t, err)
	wantModel := &model.Model{
		TableName: "test",
		FieldMap: map[string]*model.Field{
			"Name": &model.Field{
				ColName: "test_name_tt",
				GoName:  "Name",
				Typ:     reflect.TypeOf(""),
			},
		},
		ColumnMap: map[string]*model.Field{
			"test_name_tt": &model.Field{
				ColName: "test_name_tt",
				GoName:  "Name",
				Typ:     reflect.TypeOf(""),
			},
		},
	}
	assert.Equal(t, m, wantModel)

	re := NewSyncRegister()
	_, err = re.Register(&model.TableNameTest{}, WithTableName("test"), WithColumnName("XXX", "test_name"))
	assert.Equal(t, err, errs.NewErrUnknownField("XXX"))
}
