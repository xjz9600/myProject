package orm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	creator "myProject/orm/internal/value"
	"myProject/orm/model"
	"time"
)

type DB struct {
	db *sql.DB
	core
}

func (d *DB) getCore() core {
	return d.core
}

func (d *DB) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.db.QueryContext(ctx, query, args...)
}

func (d *DB) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}

type DBOptions func(*DB)

func Open(driverName, dataSourceName string, opts ...DBOptions) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...), err
}

func OpenDB(db *sql.DB, opts ...DBOptions) *DB {
	core := core{
		r:       model.NewSyncMapRegister(),
		dialect: MySQL,
		creator: creator.NewUnsafeValue,
	}
	res := &DB{
		db:   db,
		core: core,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func DBWithDialect(dialect Dialect) DBOptions {
	return func(db *DB) {
		db.dialect = dialect
	}
}

func DBWithUnsafeParse() DBOptions {
	return func(db *DB) {
		db.creator = creator.NewUnsafeValue
	}
}

func DBWithSyncMapRegister() DBOptions {
	return func(db *DB) {
		db.r = model.NewSyncMapRegister()
	}
}

func DBWithMiddleWare(middleware ...MiddleWare) DBOptions {
	return func(db *DB) {
		db.mdl = middleware
	}
}

func (d *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := d.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{
		db: d,
		tx: tx,
	}, nil
}

func (d *DB) DoTx(ctx context.Context, fn func(ctx context.Context, tx *Tx) error, opts *sql.TxOptions) error {
	tx, err := d.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	// 防止用户代码崩溃后事务没有提交或回滚
	panicked := true
	defer func() {
		if err != nil || panicked {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	err = fn(ctx, tx)
	if err != nil {
		return err
	}
	panicked = false
	return err
}

func (d *DB) Wait() error {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := d.db.Ping()
			if err == nil {
				return nil
			}
			if err != driver.ErrBadConn {
				return err
			}
		}
	}
}
