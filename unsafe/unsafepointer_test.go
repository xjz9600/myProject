package unsafe

import (
	"fmt"
	"testing"
)

func TestField(t *testing.T) {
	user := &User{}
	unsafe := GetField(user)
	unsafe.SetField("Age", int8(18))
	fmt.Println(user)
}

type User struct {
	Name string
	Age  int8
}
