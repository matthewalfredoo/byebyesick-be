package queryparamdto

import (
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
	"strconv"
)

type GetAllTransactionsQuery struct {
	TransactionStatusId string `form:"transaction_status_id"`
	Limit               string `form:"limit"`
	Page                string `form:"page"`
}

func (q *GetAllTransactionsQuery) ToGetAllParams() *GetAllParams {
	param := NewGetAllParams()
	transaction := new(entity.Transaction)

	if !util.IsEmptyString(q.TransactionStatusId) {
		column := transaction.GetSqlColumnFromField("TransactionStatusId")
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.EqualTo, q.TransactionStatusId))
	}

	pageSize := appconstant.DefaultGetAllPageSize
	if !util.IsEmptyString(q.Limit) {
		noPageSize, err := strconv.Atoi(q.Limit)
		if err == nil && noPageSize > 0 {
			pageSize = noPageSize
		}
	}
	param.PageSize = &pageSize

	sortClause := appdb.NewSort(transaction.GetSqlColumnFromField("Date"))
	sortClause.Order = appdb.OrderDesc
	param.SortClauses = append(param.SortClauses, sortClause)

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
