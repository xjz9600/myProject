package sql_demo

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type User struct {
	Name string
}

func TestJsonColumn_Value(t *testing.T) {
	js := JsonColumn[User]{Valid: true, Val: User{Name: "Tom"}}
	value, err := js.Value()
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"Name":"Tom"}`), value)
	js = JsonColumn[User]{}
	value, err = js.Value()
	assert.Nil(t, err)
	assert.Nil(t, value)
}

func TestJsonColumn_Scan(t *testing.T) {
	testCases := []struct {
		name    string
		src     any
		wantErr error
		wantVal User
		valid   bool
	}{
		{
			name: "nil",
		},
		{
			name:    "string",
			src:     `{"Name":"Tom"}`,
			wantVal: User{Name: "Tom"},
			valid:   true,
		},
		{
			name:    "bytes",
			src:     []byte(`{"Name":"Tom"}`),
			wantVal: User{Name: "Tom"},
			valid:   true,
		},
		{
			name:    "error",
			src:     16,
			wantErr: errors.New("无法解析成json"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			js := &JsonColumn[User]{}
			err := js.Scan(tc.src)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, js.Val)
			assert.Equal(t, tc.valid, js.Valid)
		})
	}
}
