package queryparamdto

import (
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
	"strings"
)

type GetCartItemCheckoutQuery struct {
	CartItemIds string `form:"cart_item_ids" validate:"required,comma_separated=number"`
	Latitude    string `form:"latitude" validate:"required,latitude"`
	Longitude   string `form:"longitude" validate:"required,longitude"`
}

func (q *GetCartItemCheckoutQuery) ToGetAllParams() *GetAllParams {
	param := NewGetAllParams()
	pharmacy := new(entity.Pharmacy)
	pharmacyProduct := new(entity.PharmacyProduct)

	latColName := pharmacy.GetSqlColumnFromField("Latitude")
	lonColName := pharmacy.GetSqlColumnFromField("Longitude")

	if !util.IsEmptyString(q.Latitude) && !util.IsEmptyString(q.Longitude) {
		param.WhereClauses = append(
			param.WhereClauses,
			appdb.NewWhere(
				fmt.Sprintf("distance(%s, %s, '%s', '%s')", latColName, lonColName, q.Latitude, q.Longitude),
				appdb.LessOrEqualTo,
				appconstant.ClosestPharmacyRangeRadius,
			),
		)
	}

	param.SortClauses = append(
		param.SortClauses,
		appdb.NewSort(
			fmt.Sprintf("distance(%s, %s, '%s', '%s') %s %v", latColName, lonColName, q.Latitude, q.Longitude, appdb.LessOrEqualTo, appconstant.ClosestPharmacyRangeRadius),
			appdb.OrderAsc,
		),
	)

	param.GroupClauses = append(
		param.GroupClauses,
		appdb.NewGroupClause(pharmacyProduct.GetSqlColumnFromField("Id")),
		appdb.NewGroupClause(pharmacy.GetSqlColumnFromField("Id")),
		appdb.NewGroupClause(pharmacy.GetSqlColumnFromField("Latitude")),
		appdb.NewGroupClause(pharmacy.GetSqlColumnFromField("Longitude")),
	)

	pageSize := 1
	param.PageSize = &pageSize

	return param
}

func (q *GetCartItemCheckoutQuery) GetCartItemIds() ([]int64, error) {
	valuesInStr := strings.TrimSpace(q.CartItemIds)

	ids := make([]int64, 0)
	idsStr := strings.Split(valuesInStr, ",")

	for _, idStr := range idsStr {
		id, err := util.ParseInt64(idStr)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}
