package requestdto

import (
	"halodeksik-be/app/entity"

	"github.com/shopspring/decimal"
)

type AddPharmacyProduct struct {
	PharmacyId int64  `json:"pharmacy_id" validate:"required"`
	ProductId  int64  `json:"product_id" validate:"required"`
	IsActive   *bool  `json:"is_active" validate:"required"`
	Price      string `json:"price" validate:"required,numeric,numericgt=0"`
}

func (r AddPharmacyProduct) ToPharmacyProduct() entity.PharmacyProduct {
	price, _ := decimal.NewFromString(r.Price)
	return entity.PharmacyProduct{
		PharmacyId: r.PharmacyId,
		ProductId:  r.ProductId,
		IsActive:   *r.IsActive,
		Price:      price,
	}
}
