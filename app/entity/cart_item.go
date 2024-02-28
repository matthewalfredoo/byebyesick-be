package entity

import (
	"database/sql"
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"time"
)

type CartItem struct {
	Id              int64        `json:"id"`
	UserId          int64        `json:"user_id"`
	ProductId       int64        `json:"product_id"`
	Quantity        int32        `json:"quantity"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	DeletedAt       sql.NullTime `json:"deleted_at"`
	Product         *Product
	PharmacyProduct *PharmacyProduct
}

func (ci *CartItem) GetEntityName() string {
	return "cart_items"
}

func (ci *CartItem) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(ci).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (ci *CartItem) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", ci.GetEntityName(), ci.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (ci *CartItem) ToResponse() *responsedto.CartItemResponse {
	return &responsedto.CartItemResponse{
		Id:                              ci.Id,
		UserId:                          ci.UserId,
		ProductId:                       ci.ProductId,
		Quantity:                        ci.Quantity,
		ProductResponse:                 ci.Product.ToProductResponse(),
		PharmacyProductCheckoutResponse: ci.PharmacyProduct.ToPharmacyProductCheckoutResponse(),
	}
}
