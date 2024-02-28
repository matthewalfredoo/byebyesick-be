package queryparamdto

import (
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
	"strconv"
	"strings"
)

type GetAllProductsForAdminQuery struct {
	Search              string `form:"search"`
	SortBy              string `form:"sort_by"`
	Sort                string `form:"sort"`
	DrugClassifications string `form:"drug_class"`
	PharmacyId          string `form:"pharmacy_id" validate:"omitempty,number"`
	NotAdded            string `form:"not_added"`
	Limit               string `form:"limit"`
	Page                string `form:"page"`
}

func (q *GetAllProductsForAdminQuery) ToGetAllParams() (*GetAllParams, error) {
	const (
		sortByName = "name"
		sortByDate = "date"

		notAddedFalse = "false"
		notAddedTrue  = "true"
	)

	param := NewGetAllParams()
	product := new(entity.Product)
	pharmacyProduct := new(entity.PharmacyProduct)

	if q.Search != "" {
		words := strings.Split(q.Search, " ")
		wordToSearch := ""
		for _, word := range words {
			wordToSearch += "%" + word + "%"
		}
		param.WhereClauses = append(
			param.WhereClauses,
			appdb.NewWhereParenthesis(product.GetSqlColumnFromField("Name"), appdb.ILike, wordToSearch, true, false, appdb.OR),
			appdb.NewWhere(product.GetSqlColumnFromField("GenericName"), appdb.ILike, wordToSearch, appdb.OR),
			appdb.NewWhere(product.GetSqlColumnFromField("Description"), appdb.ILike, wordToSearch, appdb.OR),
			appdb.NewWhereParenthesis(product.GetSqlColumnFromField("Content"), appdb.ILike, wordToSearch, false, true),
		)
	}

	switch q.SortBy {
	case sortByName:
		q.SortBy = product.GetSqlColumnFromField("Name")
	case sortByDate:
		q.SortBy = product.GetSqlColumnFromField("CreatedAt")
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

	if !util.IsEmptyString(q.DrugClassifications) {
		column := product.GetSqlColumnFromField("DrugClassificationId")
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.In, q.DrugClassifications))
	}

	switch {
	case q.NotAdded == notAddedTrue && !util.IsEmptyString(q.PharmacyId) && q.GetPharmacyId() > 0:
		param.WhereClauses = append(
			param.WhereClauses,
			appdb.NewWhere(pharmacyProduct.GetSqlColumnFromField("ProductId"), appdb.Is, nil),
		)
	case q.NotAdded == notAddedFalse && !util.IsEmptyString(q.PharmacyId) && q.GetPharmacyId() > 0:
		param.WhereClauses = append(
			param.WhereClauses,
			appdb.NewWhere(pharmacyProduct.GetSqlColumnFromField("ProductId"), appdb.IsNot, nil),
		)
	default:
	}

	param.GroupClauses = append(
		param.GroupClauses,
		appdb.NewGroupClause(product.GetSqlColumnFromField("Id")),
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

func (q *GetAllProductsForAdminQuery) GetPharmacyId() int64 {
	pharmacyId, _ := strconv.ParseInt(q.PharmacyId, 10, 64)
	return pharmacyId
}
