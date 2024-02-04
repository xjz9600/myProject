package orm

import (
	"context"
)

// Queries 定义查询的终结方法,中间方法，定义在对应实现上
type Queries[T any] interface {
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) ([]*T, error)
}

// Executor 为INSERT UPDATE DELETE的终结方法
type Executor interface {
	Exec(ctx context.Context) Result
}

// QueryBuilder 生成sql语句跟参数
type QueryBuilder interface {
	Build() (*Query, error)
}

type Query struct {
	SQL  string
	Args []any
}
