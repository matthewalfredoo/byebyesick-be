package requestdto

import "halodeksik-be/app/entity"

type AddProductStockMutationRequest struct {
	PharmacyProductOriginId int64 `json:"pharmacy_product_origin_id" validate:"required"`
	PharmacyProductDestId   int64 `json:"pharmacy_product_dest_id" validate:"required"`
	Stock                   int32 `json:"stock" validate:"required,min=1"`
}

func (r *AddProductStockMutationRequest) ToProductStockMutationRequest() entity.ProductStockMutationRequest {
	return entity.ProductStockMutationRequest{
		PharmacyProductOriginId: r.PharmacyProductOriginId,
		PharmacyProductDestId:   r.PharmacyProductDestId,
		Stock:                   r.Stock,
	}
}
