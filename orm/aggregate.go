package orm

type Aggregate struct {
	fn      string
	colName string
	alias   string
}

func (a Aggregate) fieldName() string {
	return a.colName
}

func (a Aggregate) aliasName() string {
	return a.alias
}

func (a Aggregate) expr() {
}

func (a Aggregate) selectTable() {
}

func (a Aggregate) LT(val any) Predicate {
	return Predicate{
		left:  a,
		op:    opLT,
		right: valueOf(val),
	}
}

func (a Aggregate) GT(val any) Predicate {
	return Predicate{
		left:  a,
		op:    opGT,
		right: valueOf(val),
	}
}

func (a Aggregate) Eq(val any) Predicate {
	return Predicate{
		left:  a,
		op:    opEQ,
		right: valueOf(val),
	}
}

func (a Aggregate) As(alias string) Aggregate {
	return Aggregate{
		fn:      a.fn,
		colName: a.colName,
		alias:   alias,
	}
}

func Avg(c string) Aggregate {
	return Aggregate{
		fn:      "AVG",
		colName: c,
	}
}

func Max(c string) Aggregate {
	return Aggregate{
		fn:      "MAX",
		colName: c,
	}
}

func Min(c string) Aggregate {
	return Aggregate{
		fn:      "MIN",
		colName: c,
	}
}

func Count(c string) Aggregate {
	return Aggregate{
		fn:      "COUNT",
		colName: c,
	}
}

func Sum(c string) Aggregate {
	return Aggregate{
		fn:      "SUM",
		colName: c,
	}
}
