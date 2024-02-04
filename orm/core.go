package orm

import (
	"context"
	creator "myProject/orm/internal/value"
	"myProject/orm/model"
)

type core struct {
	r       model.Register
	creator creator.Creator
	dialect Dialect
	mdl     []MiddleWare
}

func get[T any](c core, sess Session, ctx context.Context, qc *QueryContext) *QueryResult {
	var root = func(ctx context.Context, qc *QueryContext) *QueryResult {
		return getHandler[T](sess, c.creator, ctx, qc)
	}
	if len(c.mdl) > 0 {
		for i := len(c.mdl) - 1; i >= 0; i-- {
			root = c.mdl[i](root)
		}
	}
	return root(ctx, qc)
}

func getHandler[T any](sess Session, creator creator.Creator, ctx context.Context, qc *QueryContext) *QueryResult {
	query, err := qc.Builder.Build()
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	rows, err := sess.queryContext(ctx, query.SQL, query.Args...)
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	for !rows.Next() {
		return &QueryResult{
			Err: ErrNoRow,
		}
	}
	tp := new(T)
	val := creator(qc.Model, tp)
	err = val.SetColumn(rows)
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	return &QueryResult{
		Result: tp,
	}
}

func exec(c core, sess Session, ctx context.Context, qc *QueryContext) *QueryResult {
	var root = func(ctx context.Context, qc *QueryContext) *QueryResult {
		q, err := qc.Builder.Build()
		if err != nil {
			return &QueryResult{
				Result: Result{
					err: err,
				},
			}
		}
		rows, er := sess.execContext(ctx, q.SQL, q.Args...)
		return &QueryResult{
			Result: Result{
				err: er,
				res: rows,
			},
		}
	}
	if len(c.mdl) > 0 {
		for i := len(c.mdl) - 1; i >= 0; i-- {
			root = c.mdl[i](root)
		}
	}
	return root(ctx, qc)
}
