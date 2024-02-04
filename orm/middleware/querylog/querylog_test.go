package querylog

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"myProject/orm"
	"testing"
)

func TestNewMiddlewareBuilder(t *testing.T) {
	var query string
	var ag []any
	m := NewBuilder().LogFunc(func(sql string, args ...any) {
		query = sql
		ag = args
	})
	db, err := orm.Open("sqlite3", "file:test.db?cache=shared&mode=memory", orm.DBWithMiddleWare(m.Builder()))
	require.NoError(t, err)
	_, _ = orm.NewSelector[TestModel](db).WithWhere(orm.C("Id").EQ(10), orm.C("Age").EQ(18)).Get(context.Background())
	assert.Equal(t, "SELECT * FROM `test_model` WHERE (`id` = ?) AND (`age` = ?);", query)
	assert.Equal(t, []any{10, 18}, ag)
	//orm.NewInserter[TestModel](db).Values(&TestModel{Id: 18}).Exec(context.Background())
	//assert.Equal(t, "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES (?,?,?,?);", query)
	//assert.Equal(t, []any{int64(18), "", int8(0), sql.NullString{}}, args)
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  sql.NullString
}
