package appdb

type DBCondition string

const (
	EqualTo          DBCondition = "="
	NotEqualTo       DBCondition = "!="
	GreaterThan      DBCondition = ">"
	GreaterOrEqualTo DBCondition = ">="
	LessThan         DBCondition = "<"
	LessOrEqualTo    DBCondition = "<="
	Is                           = "IS"
	IsNot                        = "IS NOT"
	In                           = "IN"
	Not                          = "NOT"
	Like             DBCondition = "LIKE"
	NotLike          DBCondition = "NOT LIKE"
	ILike            DBCondition = "ILIKE"
	NotILike         DBCondition = "NOT ILIKE"
	OrderAsc         DBCondition = "ASC"
	OrderDesc        DBCondition = "DESC"
	Null             DBCondition = "NULL"
)
