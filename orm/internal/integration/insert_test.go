//go:build e2e

package integration

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"myProject/orm"
	"myProject/orm/internal/errs"
	"myProject/orm/internal/test"
	"testing"
)

type MysqlSuite struct {
	Suite
}

func (m *MysqlSuite) TearDownSuite() {
	orm.NewRawQuerier[any](m.db, "truncate `integration_test`.`simple_struct`;").Exec(context.Background())
}

func TestMysql(t *testing.T) {
	suite.Run(t, &MysqlSuite{
		Suite{
			driverName:     "mysql",
			dataSourceName: "root:root@tcp(localhost:13306)/integration_test",
		},
	})
}

func (m *MysqlSuite) TestInsert() {
	db := m.db
	t := m.T()
	testCases := []struct {
		name         string
		i            *orm.Inserter[test.SimpleStruct]
		wantErr      error
		rowsAffected int64
	}{
		{
			// 查询返回错误
			name: "zero error",
			i: func() *orm.Inserter[test.SimpleStruct] {
				return orm.NewInserter[test.SimpleStruct](db).WithValues()
			}(),
			wantErr: errs.ErrInsertZeroRow,
		},
		{
			name: "insert one data",
			i: func() *orm.Inserter[test.SimpleStruct] {
				return orm.NewInserter[test.SimpleStruct](db).WithValues(test.NewSimpleStruct(12))
			}(),
			rowsAffected: 1,
		},
		{
			name: "insert multi data",
			i: func() *orm.Inserter[test.SimpleStruct] {
				return orm.NewInserter[test.SimpleStruct](db).WithValues(test.NewSimpleStruct(13), test.NewSimpleStruct(14))
			}(),
			rowsAffected: 2,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.i.Exec(context.Background())
			rows, err := res.RowsAffected()
			assert.Equal(t, err, tc.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, rows, tc.rowsAffected)
		})
	}
}
