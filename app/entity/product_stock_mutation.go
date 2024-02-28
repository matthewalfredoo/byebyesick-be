package entity

import (
	"database/sql"
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"time"
)

type ProductStockMutation struct {
	Id                         int64        `json:"id"`
	PharmacyProductId          int64        `json:"pharmacy_product_id"`
	ProductStockMutationTypeId int64        `json:"product_stock_mutation_type_id"`
	Stock                      int32        `json:"stock"`
	CreatedAt                  time.Time    `json:"created_at"`
	UpdatedAt                  time.Time    `json:"updated_at"`
	DeletedAt                  sql.NullTime `json:"deleted_at"`
	PharmacyProduct            *PharmacyProduct
	ProductStockMutationType   *ProductStockMutationType
}

func (e *ProductStockMutation) GetEntityName() string {
	return "product_stock_mutations"
}

func (e *ProductStockMutation) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(e).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (e *ProductStockMutation) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", e.GetEntityName(), e.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (e *ProductStockMutation) ToResponse() *responsedto.ProductStockMutationResponse {
	pharmacyProduct := e.PharmacyProduct.ToPharmacyProductResponse()
	if pharmacyProduct != nil {
		pharmacyProduct.Price = ""
		pharmacyProduct.Stock = nil
	}
	return &responsedto.ProductStockMutationResponse{
		Id:                               e.Id,
		PharmacyProductId:                e.PharmacyProductId,
		ProductStockMutationTypeId:       e.ProductStockMutationTypeId,
		Stock:                            e.Stock,
		MutationDate:                     e.CreatedAt.UTC(),
		PharmacyProductResponse:          pharmacyProduct,
		ProductStockMutationTypeResponse: e.ProductStockMutationType.ToResponse(),
	}
}
