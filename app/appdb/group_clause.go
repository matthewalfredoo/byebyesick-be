package appdb

type GroupClause struct {
	Column string
}

func NewGroupClause(column string) GroupClause {
	return GroupClause{Column: column}
}
