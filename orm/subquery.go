package orm

type SubQuery struct {
	s     QueryBuilder
	alias string
	tbl   TableReference
	cols  []Selectable
}

// 使用where
func (s SubQuery) expr() {
	//TODO implement me
	panic("implement me")
}

// 使用from
func (s SubQuery) table() {
	//TODO implement me
	panic("implement me")
}

func (s SubQuery) C(name string) column {
	return column{
		name:  name,
		table: s,
	}
}

func (s SubQuery) Join(right TableReference) JoinBuilder {
	return JoinBuilder{
		left:  s,
		right: right,
		typ:   "JOIN",
	}
}

func (s SubQuery) JoinLeft(right TableReference) JoinBuilder {
	return JoinBuilder{
		left:  s,
		right: right,
		typ:   "LEFT JOIN",
	}
}

func (s SubQuery) JoinRight(right TableReference) JoinBuilder {
	return JoinBuilder{
		left:  s,
		right: right,
		typ:   "RIGHT JOIN",
	}
}

func Exist(sub SubQuery) Predicate {
	return Predicate{
		op:    opExist,
		right: sub,
	}
}

func ALL(sub SubQuery) SubQueryExpr {
	return SubQueryExpr{
		s:    sub,
		pred: "ALL",
	}
}

func SOME(sub SubQuery) SubQueryExpr {
	return SubQueryExpr{
		s:    sub,
		pred: "SOME",
	}
}

func ANY(sub SubQuery) SubQueryExpr {
	return SubQueryExpr{
		s:    sub,
		pred: "ANY",
	}
}

type SubQueryExpr struct {
	s    SubQuery
	pred string
}

func (s SubQueryExpr) expr() {
	//TODO implement me
	panic("implement me")
}
