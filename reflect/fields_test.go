package reflect

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateFields(t *testing.T) {
	testCases := []struct {
		name   string
		entity any

		wantErr error
		wantRes map[string]any
	}{
		{
			name: "testReader",
			entity: testReader{
				Name: "XieJunZe",
				age:  18,
			},
			wantRes: map[string]any{"Name": "XieJunZe", "age": 0},
		},
		{
			name: "testReaderPointer",
			entity: &testReader{
				Name: "XieJunZe",
				age:  18,
			},
			wantRes: map[string]any{"Name": "XieJunZe", "age": 0},
		},
		{
			name:    "testReaderNil",
			entity:  nil,
			wantErr: errors.New("不支持 nil"),
		},
		{
			name:    "testReaderNil",
			entity:  6,
			wantErr: errors.New("不支持的类型"),
		},
		{
			name:    "testReaderZero",
			entity:  (*testReader)(nil),
			wantErr: errors.New("不支持的类型"),
		},
		{
			name:    "testReaderZeroValue",
			entity:  testReader{},
			wantRes: map[string]any{"Name": "", "age": 0},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := IterateFields(tc.entity)
			assert.Equal(t, err, tc.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, res, tc.wantRes)
		})
	}
}

func TestSetField(t *testing.T) {
	testCases := []struct {
		name string

		entity any
		field  string
		newVal any

		wantErr error

		// 修改后的 entity
		wantEntity any
	}{
		{
			name: "struct",
			entity: testReader{
				Name: "Tom",
			},
			field:   "Name",
			newVal:  "Jerry",
			wantErr: errors.New("不可修改字段"),
		},

		{
			name: "pointer",
			entity: &testReader{
				Name: "Tom",
			},
			field:  "Name",
			newVal: "Jerry",
			wantEntity: &testReader{
				Name: "Jerry",
			},
		},

		{
			name: "pointer exported",
			entity: &testReader{
				age: 18,
			},
			field:   "age",
			newVal:  19,
			wantErr: errors.New("不可修改字段"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SetField(tc.entity, tc.field, tc.newVal)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantEntity, tc.entity)
		})
	}
}

type testReader struct {
	Name string
	age  int
}
