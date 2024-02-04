package orm

import (
	"context"
	"myProject/orm/model"
)

type QueryContext struct {
	Type    string
	Builder QueryBuilder
	Model   *model.Model
}

type HandlerFunc func(ctx context.Context, qc *QueryContext) *QueryResult

type MiddleWare func(handlerFunc HandlerFunc) HandlerFunc

type QueryResult struct {
	// Result 在不同的查询里面，类型是不同的
	// Selector.Get 里面，这会是单个结果
	// Selector.GetMulti，这会是一个切片
	// 其它情况下，它会是 Result 类型
	Result any
	Err    error
}

func chain(mds ...MiddleWare) MiddleWare {
	return func(handlerFunc HandlerFunc) HandlerFunc {
		for i := len(mds) - 1; i >= 0; i-- {
			handlerFunc = mds[i](handlerFunc)
		}
		return handlerFunc
	}
}
