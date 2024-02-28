package requestdto

import "halodeksik-be/app/entity"

type AddProductStockMutation struct {
	PharmacyProductId          int64 `json:"pharmacy_product_id" validate:"required"`
	ProductStockMutationTypeId int64 `json:"product_stock_mutation_type_id" validate:"required,oneof=1 2"`
	Stock                      int32 `json:"stock" validate:"required,min=1"`
}

func (r *AddProductStockMutation) ToProductStockMutation() entity.ProductStockMutation {
	return entity.ProductStockMutation{
		PharmacyProductId:          r.PharmacyProductId,
		ProductStockMutationTypeId: r.ProductStockMutationTypeId,
		Stock:                      r.Stock,
	}
}
