package orm

type upsertBuilder[T any] struct {
	i *Inserter[T]
}

type upsert struct {
	assign []Assignable
}

func (i *Inserter[T]) OnDuplicateKey() *upsertBuilder[T] {
	return &upsertBuilder[T]{
		i: i,
	}
}

func (on *upsertBuilder[T]) Update(assigns ...Assignable) *Inserter[T] {
	on.i.onDuplicate = &upsert{
		assign: assigns,
	}
	return on.i
}
