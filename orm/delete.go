package orm

import (
	"context"
	"database/sql"
)

type Deleter[T any] struct {
	where []predicate
	builder
}

func (d *Deleter[T]) Build() (*Query, error) {
	d.sql.WriteString("DELETE FROM ")
	if d.tableName != "" {
		d.sql.WriteString(d.tableName)
	} else {
		model, err := parseModel(new(T))
		if err != nil {
			return nil, err
		}
		d.model = model
		d.sql.WriteByte('`')
		d.sql.WriteString(model.tableName)
		d.sql.WriteByte('`')
	}
	if len(d.where) > 0 {
		d.sql.WriteString(" WHERE ")
		p := d.where[0]
		for i := 1; i < len(d.where); i++ {
			p = p.And(d.where[i])
		}
		if err := d.buildExpression(p); err != nil {
			return nil, err
		}
	}
	d.sql.WriteByte(';')
	return &Query{
		SQL:  d.sql.String(),
		Args: d.args,
	}, nil
}

func (d *Deleter[T]) From(tableName string) *Deleter[T] {
	d.tableName = tableName
	return d
}

func (d *Deleter[T]) WHERE(p ...predicate) *Deleter[T] {
	d.where = p
	return d
}

func (d *Deleter[T]) Exec(ctx context.Context) (sql.Result, error) {
	//TODO implement me
	panic("implement me")
}
