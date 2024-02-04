package orm

import "myProject/orm/internal/errs"

type Dialect interface {
	quoter() byte
	OnDuplicateKeyUpdate(b *builder, odk *upsert) error
}

var MySQL Dialect = mysqlDialect{}
var SQLite3 Dialect = sqlite3Dialect{}

type mysqlDialect struct {
}

func (m mysqlDialect) quoter() byte {
	return '`'
}

func (m mysqlDialect) OnDuplicateKeyUpdate(b *builder, odk *upsert) error {
	b.sql.WriteString(" ON DUPLICATE KEY UPDATE ")
	for idx, as := range odk.assign {
		if idx > 0 {
			b.sql.WriteByte(',')
		}
		switch expr := as.(type) {
		case column:
			col, ok := b.model.FieldMap[expr.name]
			if !ok {
				return errs.NewErrUnknownField(expr.name)
			}
			b.quote(col.ColName)
			b.sql.WriteString("=VALUES(`")
			b.sql.WriteString(col.ColName)
			b.sql.WriteString("`)")
		case Assignment:
			col, ok := b.model.FieldMap[expr.column]
			if !ok {
				return errs.NewErrUnknownField(expr.column)
			}
			b.quote(col.ColName)
			b.sql.WriteString("=?")
			b.addArgs(expr.val)
		}
	}
	return nil
}

type sqlite3Dialect struct {
}

func (m sqlite3Dialect) quoter() byte {
	return '`'
}

//ON CONFLICT(`id`) DO UPDATE SET `first_name`=?;"
//ON CONFLICT(`id`) DO UPDATE SET `first_name`=excluded.`first_name`,`last_name`=excluded.`last_name`

func (m sqlite3Dialect) OnDuplicateKeyUpdate(b *builder, odk *upsert) error {
	b.sql.WriteString(" ON CONFLICT")
	if len(odk.columns) > 0 {
		b.sql.WriteByte('(')
		for i, c := range odk.columns {
			if i > 0 {
				b.sql.WriteByte(',')
			}
			err := b.buildColumn(true, C(c))
			if err != nil {
				return err
			}
		}
	}
	b.sql.WriteString(") DO UPDATE SET ")
	for idx, as := range odk.assign {
		if idx > 0 {
			b.sql.WriteByte(',')
		}
		switch expr := as.(type) {
		case column:
			col, ok := b.model.FieldMap[expr.name]
			if !ok {
				return errs.NewErrUnknownField(expr.name)
			}
			b.quote(col.ColName)
			b.sql.WriteString("=excluded.`")
			b.sql.WriteString(col.ColName)
			b.sql.WriteString("`")
		case Assignment:
			col, ok := b.model.FieldMap[expr.column]
			if !ok {
				return errs.NewErrUnknownField(expr.column)
			}
			b.quote(col.ColName)
			b.sql.WriteString("=?")
			b.addArgs(expr.val)
		}
	}
	return nil
}
