package orm

type upsertBuilder[T any] struct {
	i       *Inserter[T]
	columns []string
}

type upsert struct {
	assign  []Assignable
	columns []string
}

func (i *Inserter[T]) OnDuplicateKey() *upsertBuilder[T] {
	return &upsertBuilder[T]{
		i: i,
	}
}

func (i *upsertBuilder[T]) ConflictColumns(cols ...string) *upsertBuilder[T] {
	i.columns = cols
	return i
}

func (on *upsertBuilder[T]) Update(assigns ...Assignable) *Inserter[T] {
	on.i.upsert = &upsert{
		assign:  assigns,
		columns: on.columns,
	}
	return on.i
}
