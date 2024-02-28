package queryparamdto

import (
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
	"strconv"
	"strings"
)

type GetAllPharmacyProductsQuery struct {
	Search              string `form:"search"`
	SortBy              string `form:"sort_by"`
	Sort                string `form:"sort"`
	DrugClassifications string `form:"drug_class"`
	PharmacyId          string `form:"pharmacy_id" validate:"number"`
	Limit               string `form:"limit"`
	Page                string `form:"page"`
}

func (q *GetAllPharmacyProductsQuery) ToGetAllParams() (*GetAllParams, error) {
	const (
		sortByName  = "name"
		sortByPrice = "price"
		sortByStock = "stock"
	)

	param := NewGetAllParams()
	product := new(entity.Product)
	pharmacyProduct := new(entity.PharmacyProduct)
	pharmacyId, _ := strconv.Atoi(q.PharmacyId)

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
	case sortByPrice:
		q.SortBy = pharmacyProduct.GetSqlColumnFromField("Price")
	case sortByStock:
		q.SortBy = pharmacyProduct.GetSqlColumnFromField("Stock")
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

	if pharmacyId != 0 {
		column := pharmacyProduct.GetSqlColumnFromField("PharmacyId")
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.EqualTo, q.PharmacyId))
	}

	if !util.IsEmptyString(q.DrugClassifications) {
		column := product.GetSqlColumnFromField("DrugClassificationId")
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.In, q.DrugClassifications))
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

func (q *GetAllPharmacyProductsQuery) GetPharmacyId() int64 {
	pharmacyId, _ := util.ParseInt64(q.PharmacyId)
	return pharmacyId
}
