package queryparamdto

import (
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
	"strconv"
	"strings"
	"time"
)

type GetAllMutationRequestQuery struct {
	SortBy           string `form:"sort_by"`
	Sort             string `form:"sort"`
	StartDate        string `form:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate          string `form:"end_date" validate:"omitempty,datetime=2006-01-02"`
	PharmacyOriginId string `form:"pharmacy_origin_id" validate:"omitempty,number"`
	PharmacyDestId   string `form:"pharmacy_dest_id" validate:"omitempty,number"`
	MutationStatusId string `form:"mutation_status_id" validate:"omitempty,oneof=1 2 3"`
	Limit            string `form:"limit"`
	Page             string `form:"page"`
}

type GetAllIncomingMutationRequestQuery struct {
	PharmacyOriginId string `json:"pharmacy_origin_id" validate:"required,number,numbergt=0"`
}

type GetAllOutgoingMutationRequestQuery struct {
	PharmacyDestId string `json:"pharmacy_dest_id" validate:"required,number,numbergt=0"`
}

func (q *GetAllMutationRequestQuery) ToGetAllParams(isIncoming bool) (*GetAllParams, int64, error) {
	param := NewGetAllParams()
	mutationRequest := new(entity.ProductStockMutationRequest)
	mutationRequestStatus := new(entity.ProductStockMutationRequestStatus)
	pharmacyProduct := new(entity.PharmacyProduct)
	pharmacyOriginId, _ := util.ParseInt64(q.PharmacyOriginId)
	pharmacyDestId, _ := util.ParseInt64(q.PharmacyDestId)

	startDate := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := util.GetCurrentDate()

	switch q.SortBy {
	default:
		q.SortBy = mutationRequest.GetSqlColumnFromField("CreatedAt")
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

	if !util.IsEmptyString(q.StartDate) {
		startDate, _ = util.ParseDateTime(q.StartDate)
		column := mutationRequest.GetSqlColumnFromField("CreatedAt")
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.GreaterOrEqualTo, startDate))
	}

	if !util.IsEmptyString(q.EndDate) {
		endDate, _ = util.ParseDateTime(q.EndDate)
		column := mutationRequest.GetSqlColumnFromField("CreatedAt")
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.LessThan, endDate.AddDate(0, 0, 1)))
	}

	if !util.IsEmptyString(q.StartDate) && !util.IsEmptyString(q.EndDate) {
		if startDate.After(endDate) {
			return nil, 0, apperror.ErrStartDateAfterEndDate
		}
	}

	if !util.IsEmptyString(q.MutationStatusId) {
		column := strings.ReplaceAll(mutationRequestStatus.GetSqlColumnFromField("Id"), mutationRequestStatus.GetEntityName(), "psmrs")
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.EqualTo, q.MutationStatusId))
	}

	if pharmacyOriginId != 0 {
		column := strings.ReplaceAll(pharmacyProduct.GetSqlColumnFromField("PharmacyId"), pharmacyProduct.GetEntityName(), "ppo")
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.EqualTo, q.PharmacyOriginId))
	}

	if pharmacyDestId != 0 {
		column := strings.ReplaceAll(pharmacyProduct.GetSqlColumnFromField("PharmacyId"), pharmacyProduct.GetEntityName(), "ppd")
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.EqualTo, q.PharmacyDestId))
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

	if isIncoming {
		return param, pharmacyOriginId, nil
	}
	return param, pharmacyDestId, nil
}
