package queryparamdto

import (
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
	"strconv"
	"strings"
)

type GetAllProductsQuery struct {
	Search              string `form:"search"`
	SortBy              string `form:"sort_by"`
	Sort                string `form:"sort"`
	DrugClassifications string `form:"drug_class" validate:"omitempty,comma_separated=number"`
	Latitude            string `form:"latitude" validate:"omitempty,latitude"`
	Longitude           string `form:"longitude" validate:"omitempty,longitude"`
	Limit               string `form:"limit"`
	Page                string `form:"page"`
}

func (q *GetAllProductsQuery) ToGetAllParams() (*GetAllParams, string, string, error) {
	const (
		sortByName = "name"
		sortByDate = "date"
	)

	param := NewGetAllParams()
	product := new(entity.Product)

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
		column := fmt.Sprintf("%s.%s", product.GetEntityName(), product.GetFieldStructTag("DrugClassificationId", appconstant.JsonStructTag))
		param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.In, q.DrugClassifications))
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

	return param, q.Latitude, q.Longitude, nil
}

func (q *GetAllProductsQuery) GetCurrentLocation() (string, string) {
	return q.Latitude, q.Longitude
}
