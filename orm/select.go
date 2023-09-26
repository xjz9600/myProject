package orm

import (
	"context"
)

type Selector[T any] struct {
	where []predicate
	builder
	db *DB
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		db: db,
	}
}

func (s *Selector[T]) Build() (*Query, error) {
	s.sql.WriteString("SELECT * FROM ")
	var err error
	s.model, err = s.db.r.Get(new(T))
	if err != nil {
		return nil, err
	}
	if s.tableName != "" {
		s.sql.WriteString(s.tableName)
	} else {
		s.sql.WriteByte('`')
		s.sql.WriteString(s.model.tableName)
		s.sql.WriteByte('`')
	}
	if len(s.where) > 0 {
		s.sql.WriteString(" WHERE ")
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}
		err := s.buildExpression(p)
		if err != nil {
			return nil, err
		}
	}
	s.sql.WriteByte(';')
	return &Query{
		SQL:  s.sql.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) WithTableName(tableName string) *Selector[T] {
	s.tableName = tableName
	return s
}

func (s *Selector[T]) WithWhere(p ...predicate) *Selector[T] {
	s.where = p
	return s
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}
