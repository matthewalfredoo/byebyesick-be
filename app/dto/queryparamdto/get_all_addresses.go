package queryparamdto

import (
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
	"strconv"
)

type GetAllAddressesQuery struct {
	Limit string `form:"limit"`
	Page  string `form:"page"`
}

func (q *GetAllAddressesQuery) ToGetAllParams() (*GetAllParams, error) {
	param := NewGetAllParams()
	address := new(entity.Address)

	sortClause := appdb.NewSort(address.GetSqlColumnFromField("Status"))
	sortClause.Order = appdb.OrderAsc
	param.SortClauses = append(param.SortClauses, sortClause)

	pageSize := appconstant.DefaultGetAllPageSize
	if !util.IsEmptyString(q.Limit) {
		noPageSize, err := strconv.Atoi(q.Limit)
		if err == nil && noPageSize > 0 {
			pageSize = noPageSize
		}
	}
	param.PageSize = &pageSize

	pageId := 1
	if !util.IsEmptyString(q.Page) {
		noPageId, err := strconv.Atoi(q.Page)
		if err == nil && noPageId > 0 {
			pageId = noPageId
		}
	}
	param.PageId = &pageId

	return param, nil
}
