package orm

type DB struct {
	r Register
}

type DBOptions func(*DB)

func NewDB(opts ...DBOptions) *DB {
	res := &DB{
		r: NewSyncRegister(),
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}
