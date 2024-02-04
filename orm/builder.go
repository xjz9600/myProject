package orm

import (
	"myProject/orm/internal/errs"
	model2 "myProject/orm/model"
	"strings"
)

type builder struct {
	sql    strings.Builder
	args   []any
	model  *model2.Model
	quoter byte
	core
}

func (b *builder) buildExpression(isColumns bool, expr Expression) error {
	switch ex := expr.(type) {
	case SubQuery:
		subQuery, err := ex.s.Build()
		if err != nil {
			return err
		}
		b.sql.WriteString("(")
		// 去掉结束符号;
		b.sql.WriteString(subQuery.SQL[:len(subQuery.SQL)-1])
		b.addArgs(subQuery.Args...)
		b.sql.WriteString(")")
	case SubQueryExpr:
		b.sql.WriteString(ex.pred)
		b.sql.WriteByte(' ')
		b.buildExpression(isColumns, ex.s)
	case Predicate:
		_, ok := ex.left.(Predicate)
		if ok {
			b.sql.WriteByte('(')
		}
		err := b.buildExpression(isColumns, ex.left)
		if err != nil {
			return err
		}
		if ok {
			b.sql.WriteByte(')')
		}
		if ex.right == nil {
			return nil
		}
		b.sql.WriteByte(' ')
		b.sql.WriteString(string(ex.op))
		b.sql.WriteByte(' ')
		_, ok = ex.right.(Predicate)
		if ok {
			b.sql.WriteByte('(')
		}
		err = b.buildExpression(isColumns, ex.right)
		if err != nil {
			return err
		}
		if ok {
			b.sql.WriteByte(')')
		}
	case column:
		err := b.buildColumn(isColumns, ex)
		if err != nil {
			return err
		}
	case value:
		b.sql.WriteByte('?')
		b.addArgs(ex.val)
	case RawExpr:
		b.sql.WriteString(ex.raw)
		b.addArgs(ex.args...)
	case Aggregate:
		err := b.buildAgg(ex)
		if err != nil {
			return err
		}
	default:
	}
	return nil
}

func (b *builder) buildColumns(isColumns bool, cols ...Selectable) error {
	if len(cols) == 0 {
		b.sql.WriteString("*")
		return nil
	}
	for i, c := range cols {
		if i > 0 {
			b.sql.WriteByte(',')
		}
		switch exp := c.(type) {
		case column:
			err := b.buildColumn(isColumns, exp)
			if err != nil {
				return err
			}
		case Aggregate:
			err := b.buildColumAgg(exp)
			if err != nil {
				return err
			}
		case RawExpr:
			b.sql.WriteString(exp.raw)
			b.addArgs(exp.args...)
		}
	}
	return nil
}

func (b *builder) buildColumn(isColumns bool, col column) error {
	switch expr := col.table.(type) {
	case nil:
		field, ok := b.model.FieldMap[col.name]
		if !ok {
			return errs.NewErrUnknownField(col.name)
		}
		if !isColumns && len(col.alias) != 0 {
			b.quote(col.alias)
			return nil
		}
		b.quote(field.ColName)
		if len(col.alias) > 0 {
			b.sql.WriteString(" AS `")
			b.sql.WriteString(col.alias)
			b.sql.WriteByte('`')
		}
	case Table:
		table := col.table.(Table)
		m, err := b.r.Get(table.entity)
		if err != nil {
			return err
		}
		field, ok := m.FieldMap[col.name]
		if !ok {
			return errs.NewErrUnknownField(col.name)
		}
		if len(table.alias) > 0 {
			b.quote(table.alias)
			b.sql.WriteByte('.')
		}
		b.quote(field.ColName)
	case SubQuery:
		tableName := expr.alias
		b.quote(tableName)
		b.sql.WriteByte('.')
		if len(expr.cols) == 0 {
			col := column{name: col.name, alias: col.alias, table: expr.tbl}
			err := b.buildColumn(isColumns, col)
			return err
		}
		for _, cl := range expr.cols {
			if cl.aliasName() == col.name {
				b.quote(cl.aliasName())
				return nil
			}
			if cl.fieldName() == col.name {
				col := column{name: col.name, alias: col.alias, table: expr.tbl}
				err := b.buildColumn(isColumns, col)
				return err
			}
		}
		return errs.NewErrUnknownField(col.name)

	}
	return nil
}

func (b *builder) quote(name string) {
	b.sql.WriteByte(b.quoter)
	b.sql.WriteString(name)
	b.sql.WriteByte(b.quoter)
}

func (b *builder) addArgs(args ...any) {
	if len(args) == 0 {
		return
	}
	if b.args == nil {
		// 很少有查询能够超过八个参数
		// INSERT 除外
		b.args = make([]any, 0, 8)
	}
	b.args = append(b.args, args...)
}

func (b *builder) buildOrderAgg(orderBy []orderBy) error {
	b.sql.WriteString(" ORDER BY ")
	for idx, o := range orderBy {
		if idx > 0 {
			b.sql.WriteByte(',')
		}
		switch exp := o.column.(type) {
		case column:
			err := b.buildColumn(false, exp)
			if err != nil {
				return err
			}
		case Aggregate:
			err := b.buildAgg(exp)
			if err != nil {
				return err
			}
		}
		b.sql.WriteByte(' ')
		b.sql.WriteString(o.order)
	}
	return nil
}

func (b *builder) buildColumAgg(exp Aggregate) error {
	b.sql.WriteString(exp.fn)
	b.sql.WriteByte('(')
	err := b.buildColumn(true, C(exp.colName))
	if err != nil {
		return err
	}
	b.sql.WriteByte(')')
	if len(exp.alias) > 0 {
		b.sql.WriteString(" AS `")
		b.sql.WriteString(exp.alias)
		b.sql.WriteByte('`')
	}
	return nil
}

func (b *builder) buildAgg(exp Aggregate) error {
	if len(exp.alias) != 0 {
		_, ok := b.model.FieldMap[C(exp.colName).name]
		if !ok {
			return errs.NewErrUnknownField(C(exp.colName).name)
		}
		b.quote(exp.alias)
	} else {
		b.sql.WriteString(exp.fn)
		b.sql.WriteByte('(')
		err := b.buildColumn(true, C(exp.colName))
		if err != nil {
			return err
		}
		b.sql.WriteByte(')')
	}
	return nil
}
