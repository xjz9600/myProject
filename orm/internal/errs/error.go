package errs

import (
	"errors"
	"fmt"
)

var (
	ErrPointerOnly = errors.New("orm: 只支持一级指针作为输入，例如 *User")
)

func NewErrUnknownField(fd string) error {
	return fmt.Errorf("orm：未知字段 %s", fd)
}
