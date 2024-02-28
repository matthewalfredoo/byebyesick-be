package requestdto

import "halodeksik-be/app/entity"

type AddEditCartItem struct {
	ProductId int64 `json:"product_id" validate:"required"`
	Quantity  *int32 `json:"quantity" validate:"required"`
}

func (r AddEditCartItem) ToCartItem() entity.CartItem {
	return entity.CartItem{
		ProductId: r.ProductId,
		Quantity:  *r.Quantity,
	}
}
