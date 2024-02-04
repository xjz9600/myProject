package orm

type orderBy struct {
	column any
	order  string
}

func Asc(col any) orderBy {
	return orderBy{
		column: col,
		order:  "ASC",
	}
}

func Desc(col any) orderBy {
	return orderBy{
		column: col,
		order:  "DESC",
	}
}
