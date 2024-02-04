package orm

import (
	"context"
	"database/sql"
	"myProject/orm/model"
)

type Deleter[T any] struct {
	where []Predicate
	builder
	tableName string
}

func (d *Deleter[T]) Build() (*Query, error) {
	d.quoter = '`'
	d.sql.WriteString("DELETE FROM ")
	if d.tableName != "" {
		d.sql.WriteString(d.tableName)
	} else {
		model, err := model.ParseModel(new(T))
		if err != nil {
			return nil, err
		}
		d.model = model
		d.sql.WriteByte('`')
		d.sql.WriteString(model.TableName)
		d.sql.WriteByte('`')
	}
	if len(d.where) > 0 {
		d.sql.WriteString(" WHERE ")
		p := d.where[0]
		for i := 1; i < len(d.where); i++ {
			p = p.And(d.where[i])
		}
		if err := d.buildExpression(true, p); err != nil {
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

func (d *Deleter[T]) WHERE(p ...Predicate) *Deleter[T] {
	d.where = p
	return d
}

func (d *Deleter[T]) Exec(ctx context.Context) (sql.Result, error) {
	//TODO implement me
	panic("implement me")
}
