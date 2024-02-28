package queryparamdto

import (
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
	"strconv"
	"strings"
)

type GetAllPharmacySellForAdminQuery struct {
	Search string `form:"search"`
	Sort   string `form:"sort"`
	Limit  string `form:"limit"`
	Page   string `form:"page"`
	Year   int64  `form:"year"`
}

func (q GetAllPharmacySellForAdminQuery) ToGetAllParams() (*GetAllParams, error) {

	param := NewGetAllParams()
	pharmacy := new(entity.Pharmacy)
	user := new(entity.User)

	if q.Search != "" {
		words := strings.Split(q.Search, " ")
		wordToSearch := ""
		for _, word := range words {
			wordToSearch += "%" + word + "%"
		}
		param.WhereClauses = append(
			param.WhereClauses,
			appdb.NewWhere(pharmacy.GetSqlColumnFromField("Name"), appdb.ILike, wordToSearch),
		)
	}

	sortClause := appdb.NewSort("total_sells")
	switch q.Sort {
	case strings.ToLower(string(appdb.OrderAsc)):
		sortClause.Order = appdb.OrderAsc
	default:
		sortClause.Order = appdb.OrderDesc
	}
	param.SortClauses = append(param.SortClauses, sortClause)

	param.GroupClauses = append(
		param.GroupClauses,
		appdb.NewGroupClause(user.GetSqlColumnFromField("Email")),
		appdb.NewGroupClause(pharmacy.GetSqlColumnFromField("Id")),
	)

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
