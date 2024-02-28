package queryparamdto

import "halodeksik-be/app/appdb"

type GetAllParams struct {
	WhereClauses []appdb.WhereClause
	SortClauses  []appdb.SortClause
	GroupClauses []appdb.GroupClause
	Search       string
	PageId       *int
	PageSize     *int
}

func NewGetAllParams() *GetAllParams {
	return &GetAllParams{
		WhereClauses: make([]appdb.WhereClause, 0),
		SortClauses:  make([]appdb.SortClause, 0),
		GroupClauses: make([]appdb.GroupClause, 0),
	}
}
