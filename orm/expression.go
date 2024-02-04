package orm

// Expression 标记接口，为了可以使用 Predicate 跟 column
type Expression interface {
	expr()
}

// Selectable 标记接口，为了查询列或者语句
type Selectable interface {
	selectTable()
	fieldName() string
	aliasName() string
}

type RawExpr struct {
	raw  string
	args []any
}

func (r RawExpr) fieldName() string {
	//TODO implement me
	panic("implement me")
}

func (r RawExpr) aliasName() string {
	//TODO implement me
	panic("implement me")
}

func Raw(sqlStr string, args ...any) RawExpr {
	return RawExpr{
		sqlStr,
		args,
	}
}

func (r RawExpr) AsPredicate() Predicate {
	return Predicate{
		left: r,
	}
}

func (r RawExpr) selectTable() {
}

func (r RawExpr) expr() {
}
