package appdb

type SortClause struct {
	Column string
	Order  DBCondition
}

func NewSort(column string, order ...DBCondition) SortClause {
	if len(order) > 0 {
		return SortClause{Column: column, Order: order[0]}
	}
	return SortClause{Column: column, Order: OrderAsc}
}
