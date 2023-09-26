package orm

import (
	"github.com/stretchr/testify/assert"
	"myProject/orm/internal/errs"
	"testing"
)

func TestRegister(t *testing.T) {
	r := NewSyncMapRegister()
	m, err := r.Register(&TableNameTest{}, WithTableName("test"), WithColumnName("Name", "test_name"))
	assert.NoError(t, err)
	wantModel := &model{
		tableName: "test",
		fieldMap: map[string]*field{
			"Name": &field{
				colName: "test_name",
			},
		},
	}
	assert.Equal(t, m, wantModel)

	re := NewSyncRegister()
	_, err = re.Register(&TableNameTest{}, WithTableName("test"), WithColumnName("XXX", "test_name"))
	assert.Equal(t, err, errs.NewErrUnknownField("XXX"))
}
