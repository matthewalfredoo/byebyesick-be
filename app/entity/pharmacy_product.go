package entity

import (
	"database/sql"
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"time"

	"github.com/shopspring/decimal"
)

type PharmacyProduct struct {
	Id         int64           `json:"id"`
	PharmacyId int64           `json:"pharmacy_id"`
	ProductId  int64           `json:"product_id"`
	IsActive   bool            `json:"is_active"`
	Price      decimal.Decimal `json:"price"`
	Stock      int32           `json:"stock"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	DeletedAt  sql.NullTime    `json:"deleted_at"`
	Pharmacy   *Pharmacy
	Product    *Product
}

func (pp *PharmacyProduct) GetEntityName() string {
	return "pharmacy_products"
}

func (pp *PharmacyProduct) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(pp).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (pp *PharmacyProduct) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", pp.GetEntityName(), pp.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (pp *PharmacyProduct) ToPharmacyProductResponse() *responsedto.PharmacyProductResponse {
	if pp == nil {
		return nil
	}
	return &responsedto.PharmacyProductResponse{
		Id:               pp.Id,
		PharmacyId:       pp.PharmacyId,
		ProductId:        pp.ProductId,
		IsActive:         pp.IsActive,
		Price:            pp.Price.String(),
		Stock:            &pp.Stock,
		PharmacyResponse: pp.Pharmacy.ToPharmacyResponse(),
		ProductResponse:  pp.Product.ToProductResponse(),
	}
}

func (pp *PharmacyProduct) ToPharmacyProductCheckoutResponse() *responsedto.PharmacyProductCheckoutResponse {
	if pp == nil {
		return nil
	}
	return &responsedto.PharmacyProductCheckoutResponse{
		Id:         pp.Id,
		PharmacyId: pp.PharmacyId,
		ProductId:  pp.ProductId,
		IsActive:   pp.IsActive,
		Price:      pp.Price.String(),
		Stock:      &pp.Stock,
	}
}
