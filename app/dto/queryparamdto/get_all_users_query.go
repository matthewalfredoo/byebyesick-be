package queryparamdto

import (
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
	"strconv"
	"strings"
)

type GetAllUsersQuery struct {
	Search     string `form:"search"`
	SortBy     string `form:"sort_by"`
	Sort       string `form:"sort"`
	UserRoleId string `form:"role_id"`
	IsVerified string `form:"is_verified"`
	Limit      string `form:"limit"`
	Page       string `form:"page"`
}

func (q *GetAllUsersQuery) ToGetAllParams() (*GetAllParams, error) {
	param := NewGetAllParams()
	user := new(entity.User)

	if q.Search != "" {
		words := strings.Split(q.Search, " ")
		wordToSearch := ""
		for _, word := range words {
			wordToSearch += "%" + word + "%"
		}
		param.WhereClauses = append(
			param.WhereClauses,
			appdb.NewWhere(user.GetSqlColumnFromField("Email"), appdb.ILike, wordToSearch),
		)
	}

	switch q.SortBy {
	case "email":
		q.SortBy = user.GetSqlColumnFromField("Email")
	default:
		q.SortBy = ""
	}
	sortClause := appdb.NewSort(q.SortBy)
	switch q.Sort {
	case strings.ToLower(string(appdb.OrderDesc)):
		sortClause.Order = appdb.OrderDesc
	default:
		sortClause.Order = appdb.OrderAsc
	}
	if !util.IsEmptyString(q.SortBy) {
		param.SortClauses = append(param.SortClauses, sortClause)
	}

	if !util.IsEmptyString(q.UserRoleId) {
		column := user.GetSqlColumnFromField("UserRoleId")
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.EqualTo, q.UserRoleId))
	}

	if !util.IsEmptyString(q.IsVerified) {
		column := user.GetSqlColumnFromField("IsVerified")
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.EqualTo, q.IsVerified))
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

	return param, nil
}
