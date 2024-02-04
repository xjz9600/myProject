package querylog

import (
	"context"
	"fmt"
	"myProject/orm"
)

type MiddlewareBuilder struct {
	logFunc func(sql string, args ...any)
}

func NewBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(sql string, args ...any) {
			fmt.Println(sql, args)
		},
	}
}

func (m *MiddlewareBuilder) LogFunc(logFunc func(sql string, args ...any)) *MiddlewareBuilder {
	m.logFunc = logFunc
	return m
}

func (m *MiddlewareBuilder) Builder() orm.MiddleWare {
	return func(handlerFunc orm.HandlerFunc) orm.HandlerFunc {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			q, err := qc.Builder.Build()
			if err != nil {
				return &orm.QueryResult{
					Err: err,
				}
			}
			m.logFunc(q.SQL, q.Args...)
			return handlerFunc(ctx, qc)
		}
	}
}
