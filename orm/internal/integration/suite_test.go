package integration

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"myProject/orm"
)

type Suite struct {
	suite.Suite
	db             *orm.DB
	driverName     string
	dataSourceName string
}

func (s *Suite) SetupTest() {
	db, err := orm.Open(s.driverName, s.dataSourceName)
	assert.NoError(s.T(), err)
	err = db.Wait()
	assert.NoError(s.T(), err)
	s.db = db
}
