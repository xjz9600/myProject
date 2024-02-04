package value

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"myProject/orm/internal/errs"
	"myProject/orm/model"
	"testing"
)

func TestReflectValue_SetColumn(t *testing.T) {
	testSetColumn(t, NewReflectValue)
}

func TestUnsafeValue_SetColumn(t *testing.T) {
	testSetColumn(t, NewUnsafeValue)
}

func testSetColumn(t *testing.T, cr Creator) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = mockDB.Close() }()
	testCases := []struct {
		name       string
		wantEntity any
		wantErr    any
		genSqlRow  *sql.Rows
	}{
		{
			name: "part test Model",
			wantEntity: &TestModel{
				Id: 1,
			},
			genSqlRow: func() *sql.Rows {
				exp := mock.ExpectQuery("SELECT XXX")
				res := sqlmock.NewRows([]string{"id"})
				res.AddRow([]byte("1"))
				exp.WillReturnRows(res)
				rows, err := mockDB.QueryContext(context.Background(), "SELECT XXX")
				rows.Next()
				assert.NoError(t, err)
				return rows
			}(),
		},
		{
			name: "test Model",
			wantEntity: &TestModel{
				Id:        1,
				FirstName: "XieJunZe",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Emerson"},
			},
			genSqlRow: func() *sql.Rows {
				exp := mock.ExpectQuery("SELECT XXX")
				res := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				res.AddRow([]byte("1"), []byte("XieJunZe"), []byte("18"), []byte("Emerson"))
				exp.WillReturnRows(res)
				rows, err := mockDB.QueryContext(context.Background(), "SELECT XXX")
				rows.Next()
				assert.NoError(t, err)
				return rows
			}(),
		},
		{
			name: "invalid column",
			wantEntity: &TestModel{
				Id:        1,
				FirstName: "XieJunZe",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Emerson"},
			},
			genSqlRow: func() *sql.Rows {
				exp := mock.ExpectQuery("SELECT XXX")
				res := sqlmock.NewRows([]string{"id", "invalid"})
				res.AddRow([]byte("1"), []byte("invalid"))
				exp.WillReturnRows(res)
				rows, err := mockDB.QueryContext(context.Background(), "SELECT XXX")
				rows.Next()
				assert.NoError(t, err)
				return rows
			}(),
			wantErr: errs.NewErrUnknownField("invalid"),
		},
	}
	r := model.NewSyncRegister()
	model := &TestModel{}
	testModel, err := r.Get(model)
	assert.NoError(t, err)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val := cr(testModel, model)
			err := val.SetColumn(tc.genSqlRow)
			assert.Equal(t, err, tc.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantEntity, model)
		})
	}
}

func Benchmark_SetColumn(b *testing.B) {
	fn := func(b *testing.B, cr Creator) {
		mockDB, mock, err := sqlmock.New()
		assert.NoError(b, err)
		defer func() { _ = mockDB.Close() }()
		exp := mock.ExpectQuery("SELECT XXX")
		res := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
		row := []driver.Value{[]byte("1"), []byte("XieJunZe"), []byte("18"), []byte("Emerson")}
		for i := 0; i < b.N; i++ {
			res.AddRow(row...)
		}
		exp.WillReturnRows(res)
		rows, err := mockDB.QueryContext(context.Background(), "SELECT XXX")
		assert.NoError(b, err)
		r := model.NewSyncRegister()
		model := &TestModel{}
		testModel, err := r.Get(model)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			rows.Next()
			cr(testModel, model).SetColumn(rows)
		}
	}
	b.Run("reflect", func(b *testing.B) {
		fn(b, NewReflectValue)
	})

	b.Run("unsafe", func(b *testing.B) {
		fn(b, NewUnsafeValue)
	})
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
