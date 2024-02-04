package orm

import (
	"context"
	"myProject/orm/model"
)

type RawQuerier[T any] struct {
	sql  string
	args []any
	Session
	core
	model *model.Model
}

func NewRawQuerier[T any](sess Session, sql string, args ...any) *RawQuerier[T] {
	c := sess.getCore()
	return &RawQuerier[T]{
		Session: sess,
		core:    c,
		sql:     sql,
		args:    args,
	}
}

func (r *RawQuerier[T]) Get(ctx context.Context) (*T, error) {
	var err error
	r.model, err = r.r.Get(new(T))
	if err != nil {
		return nil, err
	}
	res := get[T](r.core, r.Session, ctx, &QueryContext{
		Type:    "RAW",
		Builder: r,
		Model:   r.model,
	})
	if res.Err != nil {
		return nil, res.Err
	}
	return res.Result.(*T), nil
}

func (r *RawQuerier[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RawQuerier[T]) Exec(ctx context.Context) Result {
	res := exec(r.core, r.Session, ctx, &QueryContext{
		Type:    "RAW",
		Builder: r,
	})
	return res.Result.(Result)
}

func (r *RawQuerier[T]) Build() (*Query, error) {
	return &Query{
		SQL:  r.sql,
		Args: r.args,
	}, nil
}
