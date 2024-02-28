package queryparamdto

import (
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
	"strconv"
	"strings"
)

type GetAllManufacturersQuery struct {
	Search string `form:"search"`
	SortBy string `form:"sort_by"`
	Sort   string `form:"sort"`
	Limit  string `form:"limit"`
	Page   string `form:"page"`
}

func (q *GetAllManufacturersQuery) ToGetAllParams() *GetAllParams {
	const (
		sortByName = "name"
	)

	param := NewGetAllParams()
	manufacturer := new(entity.Manufacturer)

	if q.Search != "" {
		words := strings.Split(q.Search, " ")
		wordToSearch := ""
		for _, word := range words {
			wordToSearch += "%" + word + "%"
		}
		param.WhereClauses = append(
			param.WhereClauses,
			appdb.NewWhere(manufacturer.GetSqlColumnFromField("Name"), appdb.ILike, wordToSearch),
		)
	}

	switch q.SortBy {
	case sortByName:
		q.SortBy = manufacturer.GetSqlColumnFromField("Name")
	default:
		q.SortBy = ""
	}
	sortClause := appdb.NewSort(q.SortBy)
	switch q.Sort {
	case strings.ToLower(string(appdb.OrderAsc)):
		sortClause.Order = appdb.OrderAsc
	default:
		sortClause.Order = appdb.OrderDesc
	}
	if !util.IsEmptyString(q.SortBy) {
		param.SortClauses = append(param.SortClauses, sortClause)
	}

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

	return param
}
