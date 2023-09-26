package orm

import (
	"myProject/orm/internal/errs"
	"strings"
)

type builder struct {
	tableName string
	sql       strings.Builder
	args      []any
	model     *model
}

func (b *builder) buildExpression(expr Expression) error {
	switch ex := expr.(type) {
	case predicate:
		_, ok := ex.left.(predicate)
		if ok {
			b.sql.WriteByte('(')
		}
		err := b.buildExpression(ex.left)
		if err != nil {
			return err
		}
		if ok {
			b.sql.WriteByte(')')
		}
		b.sql.WriteByte(' ')
		b.sql.WriteString(string(ex.op))
		b.sql.WriteByte(' ')
		_, ok = ex.right.(predicate)
		if ok {
			b.sql.WriteByte('(')
		}
		err = b.buildExpression(ex.right)
		if err != nil {
			return err
		}
		if ok {
			b.sql.WriteByte(')')
		}
	case column:
		b.sql.WriteByte('`')
		field, ok := b.model.fieldMap[ex.name]
		if !ok {
			return errs.NewErrUnknownField(ex.name)
		}
		b.sql.WriteString(field.colName)
		b.sql.WriteByte('`')
	case value:
		b.sql.WriteByte('?')
		b.args = append(b.args, ex.val)
	default:
	}
	return nil
}
