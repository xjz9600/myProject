package orm

import (
	"context"
	"myProject/orm/internal/errs"
	"strings"
)

type Selector[T any] struct {
	tableName string
	where     []predicate
	sql       strings.Builder
	args      []any
	model     *model
}

func (s *Selector[T]) Build() (*Query, error) {
	s.sql.WriteString("SELECT * FROM ")
	model, err := parseModel(new(T))
	if err != nil {
		return nil, err
	}
	s.model = model
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

func (s *Selector[T]) buildExpression(expr Expression) error {
	switch ex := expr.(type) {
	case predicate:
		_, ok := ex.left.(predicate)
		if ok {
			s.sql.WriteByte('(')
		}
		err := s.buildExpression(ex.left)
		if err != nil {
			return err
		}
		if ok {
			s.sql.WriteByte(')')
		}
		s.sql.WriteByte(' ')
		s.sql.WriteString(string(ex.op))
		s.sql.WriteByte(' ')
		_, ok = ex.right.(predicate)
		if ok {
			s.sql.WriteByte('(')
		}
		err = s.buildExpression(ex.right)
		if err != nil {
			return err
		}
		if ok {
			s.sql.WriteByte(')')
		}
	case column:
		s.sql.WriteByte('`')
		field, ok := s.model.fieldMap[ex.name]
		if !ok {
			return errs.NewErrUnknownField(ex.name)
		}
		s.sql.WriteString(field.colName)
		s.sql.WriteByte('`')
	case value:
		s.sql.WriteByte('?')
		s.args = append(s.args, ex.val)
	default:
	}
	return nil
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
