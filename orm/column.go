package orm

type opType string

const (
	opEQ    = "="
	opLT    = "<"
	opGT    = ">"
	opAND   = "AND"
	opOR    = "OR"
	opNOT   = "NOT"
	opIN    = "IN"
	opExist = "EXIST"
)

type column struct {
	name  string
	alias string
	table TableReference
}

func (c column) fieldName() string {
	return c.name
}

func (c column) aliasName() string {
	return c.alias
}

func (c column) selectTable() {
}

func (c column) assign() {
}

func (c column) As(alias string) column {
	return column{
		name:  c.name,
		alias: alias,
		table: c.table,
	}
}

type value struct {
	val any
}

func (v value) expr() {

}

func (c column) expr() {

}

func valueOf(val any) Expression {
	switch expr := val.(type) {
	case Expression:
		return expr
	default:
		return value{val: val}
	}
}

func C(name string) column {
	return column{
		name: name,
	}
}

func (c column) LT(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opLT,
		right: valueOf(val),
	}
}

func (c column) GT(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opGT,
		right: valueOf(val),
	}
}

func (c column) EQ(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opEQ,
		right: valueOf(val),
	}
}

func (c column) InQuery(subQuery SubQuery) Predicate {
	return Predicate{
		left:  c,
		op:    opIN,
		right: valueOf(subQuery),
	}
}
