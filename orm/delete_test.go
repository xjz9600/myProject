package orm

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"myProject/orm/internal/errs"
	"testing"
)

func TestDeleter_Build(t *testing.T) {
	testCase := []struct {
		name      string
		builder   QueryBuilder
		wantErr   error
		wantQuery *Query
	}{
		{
			name:    "no from ",
			builder: &Deleter[TestModel]{},
			wantQuery: &Query{
				SQL: "DELETE FROM `test_model`;",
			},
		},
		{
			name:    "empty from ",
			builder: &Deleter[TestModel]{},
			wantQuery: &Query{
				SQL: "DELETE FROM `test_model`;",
			},
		},
		{
			name:    "from",
			builder: (&Deleter[TestModel]{}).From("`test_Model`"),
			wantQuery: &Query{
				SQL: "DELETE FROM `test_Model`;",
			},
		},
		{
			name:    "from db",
			builder: (&Deleter[TestModel]{}).From("`test_db`.`test_Model`"),
			wantQuery: &Query{
				SQL: "DELETE FROM `test_db`.`test_Model`;",
			},
		},
		{
			name:    "where",
			builder: (&Deleter[TestModel]{}).WHERE(C("Age").Eq(18)),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE `age` = ?;",
				Args: []any{18},
			},
		},
		{
			name:    "not where",
			builder: (&Deleter[TestModel]{}).WHERE(Not(C("Age").Eq(18))),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE  NOT (`age` = ?);",
				Args: []any{18},
			},
		},
		{
			name:    "and where",
			builder: (&Deleter[TestModel]{}).WHERE(C("Age").Eq(18).And(C("FirstName").Eq("Tom"))),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE (`age` = ?) AND (`first_name` = ?);",
				Args: []any{18, "Tom"},
			},
		},
		{
			name:    "or where",
			builder: (&Deleter[TestModel]{}).WHERE(C("Age").Eq(18).Or(C("FirstName").Eq("Tom"))),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE (`age` = ?) OR (`first_name` = ?);",
				Args: []any{18, "Tom"},
			},
		},
		{
			name:    "empty where",
			builder: (&Deleter[TestModel]{}).WHERE(),
			wantQuery: &Query{
				SQL: "DELETE FROM `test_model`;",
			},
		},
		{
			name:    "err basic type",
			builder: (&Deleter[int]{}).WHERE(),
			wantErr: errors.New("orm: 只支持传入结构体，例如 User"),
		},
		{
			name:    "err field",
			builder: (&Deleter[TestModel]{}).WHERE(Not(C("XXX").Eq(18))),
			wantErr: errs.NewErrUnknownField("XXX"),
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.builder.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, q)
		})
	}
}
