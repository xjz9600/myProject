package errs

import (
	"errors"
	"fmt"
)

var (
	ErrPointerOnly   = errors.New("orm: 只支持传入结构体，例如 User")
	ErrRowNotFound   = errors.New("orm: 未找到数据")
	ErrInsertZeroRow = errors.New("orm: 未插入数据")
)

func NewErrUnknownField(fd string) error {
	return fmt.Errorf("orm：未知字段 %s", fd)
}

func NewErrInvalidTagContent(tag string) error {
	return fmt.Errorf("orm: 错误的标签设置: %s", tag)
}
