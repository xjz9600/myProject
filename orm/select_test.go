package orm

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"myProject/orm/internal/errs"
	"testing"
)

func TestSelector_join(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = mockDB.Close() }()
	db := OpenDB(mockDB)
	type Order struct {
		Id        int
		UsingCol1 string
		UsingCol2 string
	}

	type OrderDetail struct {
		OrderId int
		ItemId  int

		UsingCol1 string
		UsingCol2 string
	}

	type Item struct {
		Id int
	}
	testCases := []struct {
		name      string
		builder   QueryBuilder
		wantErr   error
		wantQuery *Query
	}{
		{
			name: "join",
			builder: func() QueryBuilder {
				t1 := TableOf(&OrderDetail{})
				return NewSelector[Order](db).WithTableName(t1)
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `order_detail`;",
			},
		},
		{
			name: "join-using",
			builder: func() QueryBuilder {
				t1 := TableOf(&Order{})
				t2 := TableOf(&OrderDetail{})
				t3 := t1.Join(t2).Using("UsingCol1", "UsingCol2")
				return NewSelector[Order](db).WithTableName(t3)
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM (`order` JOIN `order_detail` USING (`using_col1`,`using_col2`));",
			},
		},
		{
			name: "join-on",
			builder: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{}).As("t2")
				t3 := t1.Join(t2).On(t1.C("Id").EQ(t2.C("OrderId")))
				return NewSelector[Order](db).WithTableName(t3)
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM (`order` AS `t1` JOIN `order_detail` AS `t2` ON `t1`.`id` = `t2`.`order_id`);",
			},
		},
		{
			name: "left-join",
			builder: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{}).As("t2")
				t3 := t1.JoinLeft(t2).On(t1.C("Id").EQ(t2.C("OrderId")))
				return NewSelector[Order](db).WithTableName(t3)
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM (`order` AS `t1` LEFT JOIN `order_detail` AS `t2` ON `t1`.`id` = `t2`.`order_id`);",
			},
		},
		{
			name: "right-join",
			builder: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{}).As("t2")
				t3 := t1.JoinRight(t2).On(t1.C("Id").EQ(t2.C("OrderId")))
				return NewSelector[Order](db).WithTableName(t3)
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM (`order` AS `t1` RIGHT JOIN `order_detail` AS `t2` ON `t1`.`id` = `t2`.`order_id`);",
			},
		},
		{
			name: "join-join",
			builder: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{}).As("t2")
				t3 := t1.Join(t2).On(t1.C("Id").EQ(t2.C("OrderId")))
				t4 := TableOf(&Item{}).As("t4")
				t5 := t3.Join(t4).On(t2.C("ItemId").EQ(t4.C("Id")))
				return NewSelector[Order](db).WithTableName(t5)
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM ((`order` AS `t1` JOIN `order_detail` AS `t2` ON `t1`.`id` = `t2`.`order_id`) JOIN `item` AS `t4` ON `t2`.`item_id` = `t4`.`id`);",
			},
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

func TestSelector_Build(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = mockDB.Close() }()
	db := OpenDB(mockDB)
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
			builder: NewSelector[TestModel](db),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			name:    "select eq",
			builder: NewSelector[TestModel](db).WithWhere(C("Id").EQ(18)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `id` = ?;",
				Args: []any{18},
			},
		},
		{
			name:    "select not",
			builder: NewSelector[TestModel](db).WithWhere(Not(C("Id").EQ(18))),
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
		{
			name:    "selectTable column",
			builder: NewSelector[TestModel](db).WithColumns(C("Id"), C("FirstName")),
			wantQuery: &Query{
				SQL: "SELECT `id`,`first_name` FROM `test_model`;",
			},
		},
		{
			name:    "selectTable column alias",
			builder: NewSelector[TestModel](db).WithColumns(C("Id").As("my_id"), C("FirstName").As("my_name")),
			wantQuery: &Query{
				SQL: "SELECT `id` AS `my_id`,`first_name` AS `my_name` FROM `test_model`;",
			},
		},
		{
			name:    "selectTable aggregate",
			builder: NewSelector[TestModel](db).WithColumns(Avg("Id"), C("FirstName")),
			wantQuery: &Query{
				SQL: "SELECT AVG(`id`),`first_name` FROM `test_model`;",
			},
		},
		{
			name:    "selectTable aggregate alias",
			builder: NewSelector[TestModel](db).WithColumns(Avg("Id").As("avg_id")),
			wantQuery: &Query{
				SQL: "SELECT AVG(`id`) AS `avg_id` FROM `test_model`;",
			},
		},
		{
			name:    "selectTable raw",
			builder: NewSelector[TestModel](db).WithColumns(Raw("COUNT(DISTINCT `first_name`)")),
			wantQuery: &Query{
				SQL: "SELECT COUNT(DISTINCT `first_name`) FROM `test_model`;",
			},
		},
		{
			name:    "agg raw",
			builder: NewSelector[TestModel](db).WithWhere(C("Id").EQ(Raw("age+?", 1).AsPredicate())),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `id` = (age+?);",
				Args: []any{1},
			},
		},
		{
			name:    "expression raw",
			builder: NewSelector[TestModel](db).WithWhere(Raw("`age` < ?", 18).AsPredicate()),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `age` < ?;",
				Args: []any{18},
			},
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

func TestSelector_Get(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = mockDB.Close() }()
	db := OpenDB(mockDB)
	testCases := []struct {
		name     string
		query    string
		mockErr  error
		mockRows *sqlmock.Rows
		wantErr  error
		wantVal  *TestModel
	}{
		{
			// 查询返回错误
			name:    "query error",
			mockErr: errors.New("invalid query"),
			wantErr: errors.New("invalid query"),
			query:   "SELECT .*",
		},
		{
			name:     "no row",
			wantErr:  ErrNoRow,
			query:    "SELECT .*",
			mockRows: sqlmock.NewRows([]string{"id"}),
		},
		{
			name:    "too many column",
			wantErr: errs.NewErrUnknownField("extra_column"),
			query:   "SELECT .*",
			mockRows: func() *sqlmock.Rows {
				res := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name", "extra_column"})
				res.AddRow([]byte("1"), []byte("Da"), []byte("18"), []byte("Ming"), []byte("nothing"))
				return res
			}(),
		},
		{
			name:  "get data",
			query: "SELECT .*",
			mockRows: func() *sqlmock.Rows {
				res := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				res.AddRow([]byte("1"), []byte("Da"), []byte("18"), []byte("Ming"))
				return res
			}(),
			wantVal: &TestModel{
				Id:        1,
				FirstName: "Da",
				Age:       18,
				LastName:  &sql.NullString{String: "Ming", Valid: true},
			},
		},
	}

	for _, tc := range testCases {
		exp := mock.ExpectQuery(tc.query)
		if tc.mockErr != nil {
			exp.WillReturnError(tc.mockErr)
		} else {
			exp.WillReturnRows(tc.mockRows)
		}
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := NewSelector[TestModel](db).Get(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, res)
		})
	}
}

func TestSelector_OrderBy(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = mockDB.Close() }()
	db := OpenDB(mockDB)
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "column",
			q:    NewSelector[TestModel](db).WithOrderBy(Asc(C("Age"))),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` ORDER BY `age` ASC;",
			},
		},
		{
			name: "columns",
			q:    NewSelector[TestModel](db).WithOrderBy(Asc(C("Age")), Desc(C("Id"))),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` ORDER BY `age` ASC,`id` DESC;",
			},
		},
		{
			name:    "invalid column",
			q:       NewSelector[TestModel](db).WithOrderBy(Asc(C("Invalid"))),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			name: "alias column",
			q:    NewSelector[TestModel](db).WithColumns(C("Id").As("my_id")).WithOrderBy(Asc(C("Id").As("my_id"))),
			wantQuery: &Query{
				SQL: "SELECT `id` AS `my_id` FROM `test_model` ORDER BY `my_id` ASC;",
			},
		},
		{
			name: "agg",
			q:    NewSelector[TestModel](db).WithColumns(Avg("Id")).WithOrderBy(Asc(Avg("Id"))),
			wantQuery: &Query{
				SQL: "SELECT AVG(`id`) FROM `test_model` ORDER BY AVG(`id`) ASC;",
			},
		},
		{
			name: "agg alias",
			q:    NewSelector[TestModel](db).WithColumns(Avg("Id").As("my_id")).WithOrderBy(Asc(Avg("Id").As("my_id"))),
			wantQuery: &Query{
				SQL: "SELECT AVG(`id`) AS `my_id` FROM `test_model` ORDER BY `my_id` ASC;",
			},
		},
		{
			name: "agg alias columns",
			q:    NewSelector[TestModel](db).WithColumns(Avg("Id").As("my_id"), Max("Age").As("max_age")).WithOrderBy(Asc(Avg("Id").As("my_id")), Desc(Max("Age").As("max_age"))),
			wantQuery: &Query{
				SQL: "SELECT AVG(`id`) AS `my_id`,MAX(`age`) AS `max_age` FROM `test_model` ORDER BY `my_id` ASC,`max_age` DESC;",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}

func TestSelector_GroupBy(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = mockDB.Close() }()
	db := OpenDB(mockDB)
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// 调用了，但是啥也没传
			name: "none",
			q:    NewSelector[TestModel](db).WithGroupBy(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			// 单个
			name: "single",
			q:    NewSelector[TestModel](db).WithGroupBy(C("Age")),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` GROUP BY `age`;",
			},
		},
		{
			// 多个
			name: "multiple",
			q:    NewSelector[TestModel](db).WithGroupBy(C("Age"), C("FirstName")),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` GROUP BY `age`,`first_name`;",
			},
		},
		{
			// 不存在
			name:    "invalid column",
			q:       NewSelector[TestModel](db).WithGroupBy(C("Invalid")),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			// 带别名
			name: "alias",
			q:    NewSelector[TestModel](db).WithColumns(C("Age").As("my_age"), C("FirstName").As("my_name")).WithGroupBy(C("Age").As("my_age"), C("FirstName").As("my_name")),
			wantQuery: &Query{
				SQL: "SELECT `age` AS `my_age`,`first_name` AS `my_name` FROM `test_model` GROUP BY `my_age`,`my_name`;",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}

func TestSelector_Having(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = mockDB.Close() }()
	db := OpenDB(mockDB)
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// 调用了，但是啥也没传
			name: "none",
			q:    NewSelector[TestModel](db).WithGroupBy(C("Age")).WithHaving(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` GROUP BY `age`;",
			},
		},
		{
			// 单个条件
			name: "single",
			q: NewSelector[TestModel](db).WithGroupBy(C("Age")).
				WithHaving(C("FirstName").EQ("Deng")),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` GROUP BY `age` HAVING `first_name` = ?;",
				Args: []any{"Deng"},
			},
		},
		{
			// 多个条件
			name: "multiple",
			q: NewSelector[TestModel](db).WithGroupBy(C("Age")).
				WithHaving(C("FirstName").EQ("Deng"), C("LastName").EQ("Ming")),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` GROUP BY `age` HAVING (`first_name` = ?) AND (`last_name` = ?);",
				Args: []any{"Deng", "Ming"},
			},
		},
		{
			// 聚合函数
			name: "avg",
			q: NewSelector[TestModel](db).WithGroupBy(C("Age")).
				WithHaving(Avg("Age").Eq(18)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` GROUP BY `age` HAVING AVG(`age`) = ?;",
				Args: []any{18},
			},
		},
		{
			// 聚合函数
			name: "alias",
			q: NewSelector[TestModel](db).WithColumns(Avg("Age").As("my_age")).WithGroupBy(C("Age")).
				WithHaving(Avg("Age").As("my_age").Eq(18)),
			wantQuery: &Query{
				SQL:  "SELECT AVG(`age`) AS `my_age` FROM `test_model` GROUP BY `age` HAVING `my_age` = ?;",
				Args: []any{18},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}
