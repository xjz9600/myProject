package orm

type Assignable interface {
	assign()
}

type Assignment struct {
	column string
	val    any
}

func (a Assignment) assign() {
}

func Assign(columnName string, val any) Assignment {
	return Assignment{
		column: columnName,
		val:    val,
	}
}
