package queryparamdto

import (
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
	"strconv"
	"time"
)

type GetAllStockMutationsQuery struct {
	StartDate  string `form:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate    string `form:"end_date" validate:"omitempty,datetime=2006-01-02"`
	PharmacyId string `form:"pharmacy_id" validate:"omitempty,number"`
	Limit      string `form:"limit"`
	Page       string `form:"page"`
}

func (q *GetAllStockMutationsQuery) ToGetAllParams(pharmacyAdminId int64) (*GetAllParams, int64, error) {
	param := NewGetAllParams()
	stockMutation := new(entity.ProductStockMutation)
	mutationType := new(entity.ProductStockMutationType)
	pharmacyProduct := new(entity.PharmacyProduct)
	pharmacy := new(entity.Pharmacy)
	product := new(entity.Product)
	manufacturer := new(entity.Manufacturer)
	pharmacyId, _ := util.ParseInt64(q.PharmacyId)

	startDate := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := util.GetCurrentDate()

	if !util.IsEmptyString(q.StartDate) {
		startDate, _ = util.ParseDateTime(q.StartDate)
		column := stockMutation.GetSqlColumnFromField("CreatedAt")
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.GreaterOrEqualTo, startDate))
	}

	if !util.IsEmptyString(q.EndDate) {
		endDate, _ = util.ParseDateTime(q.EndDate)
		column := stockMutation.GetSqlColumnFromField("CreatedAt")
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.LessThan, endDate.AddDate(0, 0, 1)))
	}

	if !util.IsEmptyString(q.StartDate) && !util.IsEmptyString(q.EndDate) {
		if startDate.After(endDate) {
			return nil, 0, apperror.ErrStartDateAfterEndDate
		}
	}

	if pharmacyId != 0 {
		column := pharmacyProduct.GetSqlColumnFromField("PharmacyId")
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.EqualTo, q.PharmacyId))
	}

	param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(pharmacy.GetSqlColumnFromField("PharmacyAdminId"), appdb.EqualTo, pharmacyAdminId))

	param.SortClauses = append(param.SortClauses, appdb.NewSort(stockMutation.GetSqlColumnFromField("CreatedAt"), appdb.OrderDesc))

	param.GroupClauses = append(
		param.GroupClauses,
		appdb.NewGroupClause(stockMutation.GetSqlColumnFromField("Id")),
		appdb.NewGroupClause(mutationType.GetSqlColumnFromField("Name")),
		appdb.NewGroupClause(pharmacy.GetSqlColumnFromField("Id")),
		appdb.NewGroupClause(product.GetSqlColumnFromField("Id")),
		appdb.NewGroupClause(manufacturer.GetSqlColumnFromField("Id")),
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

	return param, pharmacyId, nil
}
