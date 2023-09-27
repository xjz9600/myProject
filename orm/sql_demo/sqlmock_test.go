package sql_demo

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSqlMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mockRows := sqlmock.NewRows([]string{"name"})
	mockRows.AddRow("XieJunZe")
	mock.ExpectQuery("SELECT " + "\\*" + " FROM `user` .*").WillReturnRows(mockRows)
	rows, err := db.Query("SELECT * FROM `user` WHERE name = ?", "XieJunZe")
	assert.NoError(t, err)
	for rows.Next() {
		userMock := &User{}
		err := rows.Scan(&userMock.Name)
		assert.NoError(t, err)
		assert.Equal(t, userMock.Name, "XieJunZe")
	}
	mock.ExpectQuery("SELECT " + "\\*" + " FROM `user` .*").WillReturnError(errors.New("mock error"))
	_, err = db.Query("SELECT * FROM `user` WHERE name = ?", "XieJunZe")
	assert.Equal(t, err, errors.New("mock error"))
}
