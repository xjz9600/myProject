package orm

type TableReference interface {
	table()
}

type Table struct {
	entity any
	alias  string
}

func (t Table) table() {
	//TODO implement me
	panic("implement me")
}

func TableOf(entity any) Table {
	return Table{
		entity: entity,
	}
}

func (t Table) As(alias string) Table {
	return Table{
		entity: t.entity,
		alias:  alias,
	}
}

func (t Table) Join(right TableReference) JoinBuilder {
	return JoinBuilder{
		left:  t,
		right: right,
		typ:   "JOIN",
	}
}

func (t Table) JoinLeft(right TableReference) JoinBuilder {
	return JoinBuilder{
		left:  t,
		right: right,
		typ:   "LEFT JOIN",
	}
}

func (t Table) JoinRight(right TableReference) JoinBuilder {
	return JoinBuilder{
		left:  t,
		right: right,
		typ:   "RIGHT JOIN",
	}
}

type Join struct {
	left  TableReference
	right TableReference
	typ   string
	on    []Predicate
	using []string
}

func (t Table) C(name string) column {
	return column{
		name:  name,
		table: t,
	}
}

func (j Join) Join(right TableReference) JoinBuilder {
	return JoinBuilder{
		left:  j,
		right: right,
		typ:   "JOIN",
	}
}

func (j Join) JoinLeft(right TableReference) JoinBuilder {
	return JoinBuilder{
		left:  j,
		right: right,
		typ:   "LEFT JOIN",
	}
}

func (j Join) JoinRight(right TableReference) JoinBuilder {
	return JoinBuilder{
		left:  j,
		right: right,
		typ:   "RIGHT JOIN",
	}
}

type JoinBuilder struct {
	left  TableReference
	right TableReference
	typ   string
}

func (j JoinBuilder) On(p ...Predicate) Join {
	return Join{
		left:  j.left,
		typ:   j.typ,
		right: j.right,
		on:    p,
	}
}

func (j JoinBuilder) Using(cols ...string) Join {
	return Join{
		left:  j.left,
		typ:   j.typ,
		right: j.right,
		using: cols,
	}
}

func (j Join) table() {
	//TODO implement me
	panic("implement me")
}
