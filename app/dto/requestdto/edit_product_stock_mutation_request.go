package requestdto

import "halodeksik-be/app/entity"

type EditProductStockMutationRequest struct {
	ProductStockMutationRequestStatusId int64 `json:"product_stock_mutation_request_status_id" validate:"required,oneof=2 3"`
}

func (r *EditProductStockMutationRequest) ToProductStockMutationRequest() entity.ProductStockMutationRequest {
	return entity.ProductStockMutationRequest{
		ProductStockMutationRequestStatusId: r.ProductStockMutationRequestStatusId,
	}
}
