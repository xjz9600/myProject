package orm

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"myProject/orm/internal/errs"
	"testing"
)

func TestSubQuery(t *testing.T) {
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
			name: "sub from ",
			builder: func() QueryBuilder {
				t1 := NewSelector[Order](db).AsSubQuery("sub")
				return NewSelector[OrderDetail](db).WithTableName(t1)
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM (SELECT * FROM `order`) AS `sub`;",
			},
		},
		{
			//SELECT * FROM `order` WHERE `id` IN (SELECT `order_id` FROM `order_detail`);",
			name: "in",
			builder: func() QueryBuilder {
				sub := NewSelector[OrderDetail](db).WithColumns(C("OrderId")).AsSubQuery("sub")
				return NewSelector[Order](db).WithWhere(C("Id").InQuery(sub))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `order` WHERE `id` IN (SELECT `order_id` FROM `order_detail`);",
			},
		},
		{
			name: "exist",
			builder: func() QueryBuilder {
				sub := NewSelector[OrderDetail](db).WithColumns(C("OrderId")).AsSubQuery("sub")
				return NewSelector[Order](db).WithWhere(Exist(sub))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `order` WHERE  EXIST (SELECT `order_id` FROM `order_detail`);",
			},
		},
		{
			name: "not exist",
			builder: func() QueryBuilder {
				sub := NewSelector[OrderDetail](db).WithColumns(C("OrderId")).AsSubQuery("sub")
				return NewSelector[Order](db).WithWhere(Not(Exist(sub)))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `order` WHERE  NOT ( EXIST (SELECT `order_id` FROM `order_detail`));",
			},
		},
		{
			name: "all",
			builder: func() QueryBuilder {
				sub := NewSelector[OrderDetail](db).WithColumns(C("OrderId")).AsSubQuery("sub")
				return NewSelector[Order](db).WithWhere(C("Id").GT(ALL(sub)))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `order` WHERE `id` > ALL (SELECT `order_id` FROM `order_detail`);",
			},
		},
		{
			name: "some and any",
			builder: func() QueryBuilder {
				sub := NewSelector[OrderDetail](db).WithColumns(C("OrderId")).AsSubQuery("sub")
				return NewSelector[Order](db).WithWhere(C("Id").GT(SOME(sub)), C("Id").LT(ANY(sub)))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `order` WHERE (`id` > SOME (SELECT `order_id` FROM `order_detail`)) AND (`id` < ANY (SELECT `order_id` FROM `order_detail`));",
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
func TestSelector_SubqueryAndJoin(t *testing.T) {
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
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// 虽然泛型是 Order，但是我们传入 OrderDetail
			name: "table and join",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{})
				sub := NewSelector[OrderDetail](db).AsSubQuery("sub")
				return NewSelector[Order](db).WithColumns(sub.C("ItemId")).WithTableName(t1.Join(sub).On(t1.C("Id").EQ(sub.C("OrderId"))))
			}(),
			wantQuery: &Query{
				SQL: "SELECT `sub`.`item_id` FROM (`order` JOIN (SELECT * FROM `order_detail`) AS `sub` ON `id` = `sub`.`order_id`);",
			},
		},
		{
			name: "table and left join",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{})
				sub := NewSelector[OrderDetail](db).AsSubQuery("sub")
				return NewSelector[Order](db).WithTableName(sub.Join(t1).On(t1.C("Id").EQ(sub.C("OrderId"))))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM ((SELECT * FROM `order_detail`) AS `sub` JOIN `order` ON `id` = `sub`.`order_id`);",
			},
		},
		{
			name: "join and join",
			q: func() QueryBuilder {
				sub1 := NewSelector[OrderDetail](db).AsSubQuery("sub1")
				sub2 := NewSelector[OrderDetail](db).AsSubQuery("sub2")
				return NewSelector[Order](db).WithTableName(sub1.JoinRight(sub2).Using("Id"))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM ((SELECT * FROM `order_detail`) AS `sub1` RIGHT JOIN (SELECT * FROM `order_detail`) AS `sub2` USING (`id`));",
			},
		},
		{
			name: "join sub sub",
			q: func() QueryBuilder {
				sub1 := NewSelector[OrderDetail](db).AsSubQuery("sub1")
				sub2 := NewSelector[OrderDetail](db).WithTableName(sub1).AsSubQuery("sub2")
				t1 := TableOf(&Order{}).As("o1")
				return NewSelector[Order](db).WithTableName(sub2.Join(t1).Using("Id"))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM ((SELECT * FROM (SELECT * FROM `order_detail`) AS `sub1`) AS `sub2` JOIN `order` AS `o1` USING (`id`));",
			},
		},
		{
			name: "invalid field",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{})
				sub := NewSelector[OrderDetail](db).AsSubQuery("sub")
				return NewSelector[Order](db).WithColumns(sub.C("Invalid")).WithTableName(t1.Join(sub).On(t1.C("Id").EQ(sub.C("OrderId"))))
			}(),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			name: "invalid field in predicates",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{})
				sub := NewSelector[OrderDetail](db).AsSubQuery("sub")
				return NewSelector[Order](db).WithColumns(sub.C("ItemId")).WithTableName(t1.Join(sub).On(t1.C("Id").EQ(sub.C("Invalid"))))
			}(),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			name: "invalid field in aggregate function",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{})
				sub := NewSelector[OrderDetail](db).AsSubQuery("sub")
				return NewSelector[Order](db).WithColumns(Max("Invalid")).WithTableName(t1.Join(sub).On(t1.C("Id").EQ(sub.C("OrderId"))))
			}(),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			name: "not selected",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{})
				sub := NewSelector[OrderDetail](db).WithColumns(C("OrderId")).AsSubQuery("sub")
				return NewSelector[Order](db).WithColumns(sub.C("ItemId")).WithTableName(t1.Join(sub).On(t1.C("Id").EQ(sub.C("OrderId"))))
			}(),
			wantErr: errs.NewErrUnknownField("ItemId"),
		},
		{
			name: "use alias",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{})
				sub := NewSelector[OrderDetail](db).WithColumns(C("OrderId").As("oid")).AsSubQuery("sub")
				return NewSelector[Order](db).WithColumns(sub.C("oid")).WithTableName(t1.Join(sub).On(t1.C("Id").EQ(sub.C("oid"))))
			}(),
			wantQuery: &Query{
				SQL: "SELECT `sub`.`oid` FROM (`order` JOIN (SELECT `order_id` AS `oid` FROM `order_detail`) AS `sub` ON `id` = `sub`.`oid`);",
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
