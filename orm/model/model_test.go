package orm

import (
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"myProject/orm/internal/errs"
	"reflect"
	"testing"
)

func Test_parseModel(t *testing.T) {
	testCases := []struct {
		name       string
		val        any
		wantModel  *Model
		wantFields []*Field
		wantErr    error
	}{
		{
			name:    "test Model",
			val:     TestModel{},
			wantErr: errors.New("orm: 只支持传入结构体，例如 User"),
		},
		{
			// 指针
			name: "pointer",
			val:  &TestModel{},
			wantFields: []*Field{
				{
					ColName: "id",
					GoName:  "Id",
					Typ:     reflect.TypeOf(int64(0)),
				},
				{
					ColName: "first_name",
					GoName:  "FirstName",
					Typ:     reflect.TypeOf(""),
				},
				{
					ColName: "age",
					GoName:  "Age",
					Typ:     reflect.TypeOf(int8(0)),
				},
				{
					ColName: "last_name",
					GoName:  "LastName",
					Typ:     reflect.TypeOf(&sql.NullString{}),
				},
			},
			wantModel: &Model{
				TableName: "test_model",
			},
		},
		{
			// 多级指针
			name: "multiple pointer",
			// 因为 Go 编译器的原因，所以我们写成这样
			val: func() any {
				val := &TestModel{}
				return &val
			}(),
			wantErr: errors.New("orm: 只支持传入结构体，例如 User"),
		},
		{
			name:    "map",
			val:     map[string]string{},
			wantErr: errors.New("orm: 只支持传入结构体，例如 User"),
		},
		{
			name:    "slice",
			val:     []int{},
			wantErr: errors.New("orm: 只支持传入结构体，例如 User"),
		},
		{
			name:    "basic type",
			val:     0,
			wantErr: errors.New("orm: 只支持传入结构体，例如 User"),
		},
		{
			// 标签测试跟实现TableName接口测试
			name: "TableName and tag",
			val:  &TableNameTest{},
			wantModel: &Model{
				TableName: "test_table_name",
				FieldMap: map[string]*Field{
					"Name": {
						ColName: "test_name",
					},
				},
			},
			wantFields: []*Field{
				{
					GoName:  "Name",
					Typ:     reflect.TypeOf(""),
					ColName: "test_name",
				},
			},
		},
		{
			// 标签格式错误
			name:    "err tag",
			val:     &TagErrTest{},
			wantErr: errs.NewErrInvalidTagContent("columnTest"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := parseModel(tc.val)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			for _, fd := range tc.wantFields {
				if tc.wantModel.FieldMap == nil {
					tc.wantModel.FieldMap = make(map[string]*Field, len(tc.wantFields))
				}
				if tc.wantModel.ColumnMap == nil {
					tc.wantModel.ColumnMap = make(map[string]*Field, len(tc.wantFields))
				}
				tc.wantModel.FieldMap[fd.GoName] = fd
				tc.wantModel.ColumnMap[fd.ColName] = fd
			}
			assert.Equal(t, tc.wantModel, m)
		})
	}
}

type TagErrTest struct {
	Name string `orm:"columnTest"`
}

type TableNameTest struct {
	Name string `orm:"column=test_name"`
}

func (t *TableNameTest) TableName() string {
	return "test_table_name"
}

func Test_underscoreName(t *testing.T) {
	testCases := []struct {
		name    string
		srcStr  string
		wantStr string
	}{
		// 我们这些用例就是为了确保
		// 在忘记 underscoreName 的行为特性之后
		// 可以从这里找回来
		// 比如说过了一段时间之后
		// 忘记了 ID 不能转化为 id
		// 那么这个测试能帮我们确定 ID 只能转化为 i_d
		{
			name:    "upper cases",
			srcStr:  "ID",
			wantStr: "i_d",
		},
		{
			name:    "use number",
			srcStr:  "Table1Name",
			wantStr: "table1_name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := underscoreName(tc.srcStr)
			assert.Equal(t, tc.wantStr, res)
		})
	}
}
