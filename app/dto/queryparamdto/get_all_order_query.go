package queryparamdto

import (
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
	"strconv"
	"strings"
)

type GetAllOrderAdminRequestQuery struct {
	Limit         string `form:"limit"`
	Page          string `form:"page"`
	Search        string `form:"search"`
	OrderStatusId string `form:"order_status_id" validate:"omitempty,number"`
}

func (q *GetAllOrderAdminRequestQuery) ToGetAllParams() *GetAllParams {
	param := NewGetAllParams()

	pharmacy := new(entity.Pharmacy)
	order := new(entity.Order)
	status := new(entity.OrderStatus)

	statusId, _ := strconv.Atoi(q.OrderStatusId)

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

	if statusId != 0 {
		column := status.GetSqlColumnFromField("Id")
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.EqualTo, q.OrderStatusId))
	}

	sortClause := appdb.NewSort(order.GetSqlColumnFromField("Date"))
	sortClause.Order = appdb.OrderDesc
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

	return param
}

type GetAllOrderUserRequestQuery struct {
	Search        string `form:"search"`
	Limit         string `form:"limit"`
	Page          string `form:"page"`
	OrderStatusId string `form:"order_status_id" validate:"omitempty,number"`
}

func (q GetAllOrderUserRequestQuery) ToGetAllParams() *GetAllParams {
	param := NewGetAllParams()
	pharmacy := new(entity.Pharmacy)
	status := new(entity.OrderStatus)

	statusId, _ := strconv.Atoi(q.OrderStatusId)

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

	if statusId != 0 {
		column := status.GetSqlColumnFromField("Id")
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.EqualTo, q.OrderStatusId))
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
