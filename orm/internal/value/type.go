package value

import (
	"database/sql"
	"myProject/orm/model"
)

type Value interface {
	SetColumn(rows *sql.Rows) error
	Field(name string) any
}

type Creator func(model *model.Model, entity any) Value
