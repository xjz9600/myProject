package orm

import (
	"context"
)

type Selector[T any] struct {
	where []Predicate
	table TableReference
	builder
	cols    []Selectable
	groupBy []Selectable
	having  []Predicate
	orderBy []orderBy
	Session
}

func NewSelector[T any](sess Session) *Selector[T] {
	c := sess.getCore()
	return &Selector[T]{
		Session: sess,
		builder: builder{
			quoter: c.dialect.quoter(),
			core:   c,
		},
	}
}

func (s *Selector[T]) AsSubQuery(alias string) SubQuery {
	tbl := s.table
	if tbl == nil {
		tbl = TableOf(new(T))
	}
	return SubQuery{
		s:     s,
		alias: alias,
		tbl:   tbl,
		cols:  s.cols,
	}
}

func (s *Selector[T]) Build() (*Query, error) {
	s.sql.Reset()
	s.sql.WriteString("SELECT ")
	var err error
	if s.model == nil {
		s.model, err = s.r.Get(new(T))
		if err != nil {
			return nil, err
		}
	}
	err = s.buildColumns(true, s.cols...)
	if err != nil {
		return nil, err
	}
	s.sql.WriteString(" FROM ")
	err = s.buildTable(s.table)
	if err != nil {
		return nil, err
	}
	if len(s.where) > 0 {
		s.sql.WriteString(" WHERE ")
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}
		err = s.buildExpression(true, p)
		if err != nil {
			return nil, err
		}
	}
	if len(s.groupBy) > 0 {
		s.sql.WriteString(" GROUP BY ")
		err = s.buildColumns(false, s.groupBy...)
		if err != nil {
			return nil, err
		}
	}
	if len(s.having) > 0 {
		s.sql.WriteString(" HAVING ")
		p := s.having[0]
		for i := 1; i < len(s.having); i++ {
			p = p.And(s.having[i])
		}
		err = s.buildExpression(false, p)
		if err != nil {
			return nil, err
		}
	}
	if len(s.orderBy) > 0 {
		err = s.buildOrderAgg(s.orderBy)
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

func (s *Selector[T]) WithTableName(tbl TableReference) *Selector[T] {
	s.table = tbl
	return s
}

func (s *Selector[T]) WithGroupBy(groupBy ...Selectable) *Selector[T] {
	s.groupBy = groupBy
	return s
}

func (s *Selector[T]) WithOrderBy(orderBy ...orderBy) *Selector[T] {
	s.orderBy = orderBy
	return s
}

func (s *Selector[T]) WithHaving(having ...Predicate) *Selector[T] {
	s.having = having
	return s
}

func (s *Selector[T]) WithWhere(p ...Predicate) *Selector[T] {
	s.where = p
	return s
}

func (s *Selector[T]) WithColumns(p ...Selectable) *Selector[T] {
	s.cols = p
	return s
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	var err error
	s.model, err = s.r.Get(new(T))
	if err != nil {
		return nil, err
	}
	res := get[T](s.core, s.Session, ctx, &QueryContext{
		Type:    "SELECT",
		Builder: s,
		Model:   s.model,
	})
	if res.Err != nil {
		return nil, res.Err
	}
	return res.Result.(*T), nil
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) buildTable(table TableReference) error {
	switch exp := table.(type) {
	case nil:
		s.quote(s.model.TableName)
	case Table:
		m, err := s.r.Get(exp.entity)
		if err != nil {
			return err
		}
		s.quote(m.TableName)
		if len(exp.alias) > 0 {
			s.sql.WriteString(" AS ")
			s.quote(exp.alias)
		}
	case Join:
		s.sql.WriteByte('(')
		err := s.buildTable(exp.left)
		if err != nil {
			return err
		}
		s.sql.WriteByte(' ')
		s.sql.WriteString(exp.typ)
		s.sql.WriteByte(' ')
		err = s.buildTable(exp.right)
		if err != nil {
			return err
		}
		if len(exp.on) > 0 {
			s.sql.WriteString(" ON ")
			p := exp.on[0]
			for i := 1; i < len(exp.on); i++ {
				p = p.And(exp.on[i])
			}
			err = s.buildExpression(false, p)
			if err != nil {
				return err
			}
		}
		if len(exp.using) > 0 {
			s.sql.WriteString(" USING (")
			for i, u := range exp.using {
				if i > 0 {
					s.sql.WriteByte(',')
				}
				err := s.buildColumn(false, C(u))
				if err != nil {
					return err
				}
			}
			s.sql.WriteByte(')')
		}
		s.sql.WriteByte(')')
	case SubQuery:
		subQuery, err := exp.s.Build()
		if err != nil {
			return err
		}
		s.sql.WriteString("(")
		// 去掉结束符号;
		s.sql.WriteString(subQuery.SQL[:len(subQuery.SQL)-1])
		s.addArgs(subQuery.Args...)
		s.sql.WriteString(") AS ")
		s.quote(exp.alias)
	}
	return nil
}
