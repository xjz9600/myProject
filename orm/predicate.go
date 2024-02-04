package orm

type Predicate struct {
	left  Expression
	op    opType
	right Expression
}

func (p Predicate) expr() {

}

func (left Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opAND,
		right: right,
	}
}

func (left Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opOR,
		right: right,
	}
}

func Not(right Predicate) Predicate {
	return Predicate{
		op:    opNOT,
		right: right,
	}
}
