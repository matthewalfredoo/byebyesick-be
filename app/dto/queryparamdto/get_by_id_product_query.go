package queryparamdto

import (
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
)

type GetByIdProductQuery struct {
	Latitude  string `form:"latitude" validate:"omitempty,latitude"`
	Longitude string `form:"longitude" validate:"omitempty,longitude"`
}

func (q *GetByIdProductQuery) ToGetAllParams() *GetAllParams {
	param := NewGetAllParams()
	product := new(entity.Product)
	productCategory := new(entity.ProductCategory)
	manufacturer := new(entity.Manufacturer)
	drugClassification := new(entity.DrugClassification)
	pharmacy := new(entity.Pharmacy)

	if !util.IsEmptyString(q.Latitude) && !util.IsEmptyString(q.Longitude) {
		latColName := pharmacy.GetSqlColumnFromField("Latitude")
		lonColName := pharmacy.GetSqlColumnFromField("Longitude")

		param.WhereClauses = append(
			param.WhereClauses,
			appdb.NewWhere(
				fmt.Sprintf("distance(%s, %s, '%s', '%s')", latColName, lonColName, q.Latitude, q.Longitude),
				appdb.LessOrEqualTo,
				appconstant.ClosestPharmacyRangeRadius,
			),
		)
	}

	param.GroupClauses = append(
		param.GroupClauses,
		appdb.NewGroupClause(product.GetSqlColumnFromField("Id")),
		appdb.NewGroupClause(productCategory.GetSqlColumnFromField("Id")),
		appdb.NewGroupClause(manufacturer.GetSqlColumnFromField("Id")),
		appdb.NewGroupClause(drugClassification.GetSqlColumnFromField("Id")),
	)

	return param
}
