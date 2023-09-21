package orm

type opType string

const (
	opEQ  = "="
	opLT  = "<"
	opGT  = ">"
	opAND = "AND"
	opOR  = "OR"
	opNOT = "NOT"
)

type column struct {
	name string
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

func (c column) LT(val any) predicate {
	return predicate{
		left:  c,
		op:    opLT,
		right: valueOf(val),
	}
}

func (c column) GT(val any) predicate {
	return predicate{
		left:  c,
		op:    opGT,
		right: valueOf(val),
	}
}

func (c column) Eq(val any) predicate {
	return predicate{
		left:  c,
		op:    opEQ,
		right: valueOf(val),
	}
}
