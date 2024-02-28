package queryparamdto

import (
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
)

type GetAllPharmacySellsMonthlyForAdminQuery struct {
	ProductId         string `form:"product_id"`
	ProductCategoryId string `form:"product_category_id"`
	Year              int64  `form:"year"`
}

func (q GetAllPharmacySellsMonthlyForAdminQuery) ToGetAllParams() (*GetAllParams, error) {
	param := NewGetAllParams()
	product := new(entity.Product)

	param.GroupClauses = append(
		param.GroupClauses,
		appdb.NewGroupClause("month"),
	)

	sortClause := appdb.NewSort("month")
	sortClause.Order = appdb.OrderAsc

	param.SortClauses = append(param.SortClauses, sortClause)

	if !util.IsEmptyString(q.ProductId) {
		column := product.GetSqlColumnFromField("Id")
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.EqualTo, q.ProductId))
	}

	if !util.IsEmptyString(q.ProductCategoryId) {
		column := product.GetSqlColumnFromField("ProductCategoryId")
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.EqualTo, q.ProductCategoryId))
	}
	monthPageSize := appconstant.MonthInAYearPageSize
	param.PageSize = &monthPageSize

	pageId := 1
	param.PageId = &pageId

	return param, nil
}
