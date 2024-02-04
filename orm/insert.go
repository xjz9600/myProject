package orm

import (
	"context"
	"myProject/orm/internal/errs"
	model2 "myProject/orm/model"
)

type Inserter[T any] struct {
	Session
	core
	values []any
	builder
	columns []string
	upsert  *upsert
}

func (i *Inserter[T]) Exec(ctx context.Context) Result {
	var err error
	i.model, err = i.r.Get(new(T))
	if err != nil {
		return Result{err: err}
	}
	res := exec(i.core, i.Session, ctx, &QueryContext{
		Type:    "INSERT",
		Builder: i,
		Model:   i.model,
	})
	return res.Result.(Result)
}

func NewInserter[T any](sess Session) *Inserter[T] {
	c := sess.getCore()
	return &Inserter[T]{
		Session: sess,
		core:    c,
		builder: builder{
			quoter: c.dialect.quoter(),
		},
	}
}

func (i *Inserter[T]) WithValues(val ...any) *Inserter[T] {
	i.values = val
	return i
}

func (i *Inserter[T]) WithColumns(cols ...string) *Inserter[T] {
	i.columns = cols
	return i
}

func (i *Inserter[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return nil, errs.ErrInsertZeroRow
	}
	var err error
	i.model, err = i.r.Get(new(T))
	if err != nil {
		return nil, err
	}
	i.sql.Reset()
	i.sql.WriteString("INSERT INTO ")
	i.quote(i.model.TableName)
	i.sql.WriteString(" (")
	fields := i.model.Fields
	if len(i.columns) > 0 {
		fields = make([]*model2.Field, 0, len(i.columns))
		for _, c := range i.columns {
			f, ok := i.model.FieldMap[c]
			if !ok {
				return nil, errs.NewErrUnknownField(c)
			}
			fields = append(fields, f)
		}
	}
	for idx, f := range fields {
		if idx > 0 {
			i.sql.WriteByte(',')
		}
		i.quote(f.ColName)
	}
	i.sql.WriteString(") VALUES")
	for idx, v := range i.values {
		if idx > 0 {
			i.sql.WriteByte(',')
		}
		i.sql.WriteByte('(')
		for idx, f := range fields {
			if idx > 0 {
				i.sql.WriteByte(',')
			}
			i.sql.WriteByte('?')
			i.addArgs(i.creator(i.model, v).Field(f.GoName))
		}
		i.sql.WriteByte(')')
	}
	if i.upsert != nil {
		err = i.dialect.OnDuplicateKeyUpdate(&i.builder, i.upsert)
		if err != nil {
			return nil, err
		}
	}
	i.sql.WriteByte(';')
	return &Query{
		SQL:  i.sql.String(),
		Args: i.args,
	}, nil
}
