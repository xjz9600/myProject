package orm

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"myProject/orm/internal/errs"
	"testing"
)

func TestSelector_Build(t *testing.T) {
	db := NewDB()
	testCases := []struct {
		name      string
		builder   QueryBuilder
		wantErr   error
		wantQuery *Query
	}{
		{
			name:    "select",
			builder: NewSelector[TestModel](db),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			name:    "use from",
			builder: NewSelector[TestModel](db).WithTableName("`test_model`"),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			// 调用 FROM，同时出入看了 DB
			name:    "with db",
			builder: NewSelector[TestModel](db).WithTableName("`test_db`.`test_model`"),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_db`.`test_model`;",
			},
		},
		{
			name:    "select eq",
			builder: NewSelector[TestModel](db).WithWhere(C("Id").Eq(18)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `id` = ?;",
				Args: []any{18},
			},
		},
		{
			name:    "select not",
			builder: NewSelector[TestModel](db).WithWhere(Not(C("Id").Eq(18))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE  NOT (`id` = ?);",
				Args: []any{18},
			},
		},
		{
			name:    "select and",
			builder: NewSelector[TestModel](db).WithWhere(C("Id").GT(18).And(C("Age").LT(18))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`id` > ?) AND (`age` < ?);",
				Args: []any{18, 18},
			},
		},
		{
			name:    "select or",
			builder: NewSelector[TestModel](db).WithWhere(C("Id").GT(18).Or(C("Age").LT(18))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`id` > ?) OR (`age` < ?);",
				Args: []any{18, 18},
			},
		},
		{
			name:    "invalid column",
			builder: NewSelector[TestModel](db).WithWhere(C("Invalid").GT(18)),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.builder.Build()
			assert.Equal(t, err, tc.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, query, tc.wantQuery)
		})
	}
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
