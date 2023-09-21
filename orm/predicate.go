package orm

type predicate struct {
	left  Expression
	op    opType
	right Expression
}

// Expression 标记接口，为了可以使用 predicate 跟 column
type Expression interface {
	expr()
}

func (p predicate) expr() {

}

func (left predicate) And(right predicate) predicate {
	return predicate{
		left:  left,
		op:    opAND,
		right: right,
	}
}

func (left predicate) Or(right predicate) predicate {
	return predicate{
		left:  left,
		op:    opOR,
		right: right,
	}
}

func Not(right predicate) predicate {
	return predicate{
		op:    opNOT,
		right: right,
	}
}
