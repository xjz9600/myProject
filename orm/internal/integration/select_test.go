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

type MysqlSelectSuite struct {
	Suite
}

func (s *MysqlSelectSuite) SetupTest() {
	s.Suite.SetupTest()
	orm.NewInserter[test.SimpleStruct](s.db).WithValues(test.NewSimpleStruct(12)).Exec(context.Background())
}

func (s *MysqlSelectSuite) TearDownTest() {
	orm.NewRawQuerier[any](s.db, "truncate `integration_test`.`simple_struct`;").Exec(context.Background())
}

func TestMysqlSelect(t *testing.T) {
	suite.Run(t, &MysqlSelectSuite{
		Suite{
			driverName:     "mysql",
			dataSourceName: "root:root@tcp(localhost:13306)/integration_test",
		},
	})
}

func (m *MysqlSelectSuite) TestSelect() {
	db := m.db
	t := m.T()
	testCases := []struct {
		name    string
		wantErr error
		s       *orm.Selector[test.SimpleStruct]
		wantVal *test.SimpleStruct
	}{
		{
			name: "get data",
			s: func() *orm.Selector[test.SimpleStruct] {
				return orm.NewSelector[test.SimpleStruct](db).WithWhere(orm.C("Id").EQ(12))
			}(),
			wantVal: test.NewSimpleStruct(12),
		},
		{
			name: "get no error",
			s: func() *orm.Selector[test.SimpleStruct] {
				return orm.NewSelector[test.SimpleStruct](db).WithWhere(orm.C("Id").EQ(15))
			}(),
			wantErr: errs.ErrRowNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.s.Get(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, res)
		})
	}
}
