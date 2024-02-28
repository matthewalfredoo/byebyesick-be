package requestdto

import (
	"halodeksik-be/app/entity"

	"github.com/shopspring/decimal"
)

type EditPharmacyProduct struct {
	IsActive *bool  `json:"is_active" validate:"required"`
	Price    string `json:"price" validate:"required,numeric,numericgt=0"`
}

func (r EditPharmacyProduct) ToPharmacyProduct() entity.PharmacyProduct {
	price, _ := decimal.NewFromString(r.Price)
	return entity.PharmacyProduct{
		IsActive: *r.IsActive,
		Price:    price,
	}
}
