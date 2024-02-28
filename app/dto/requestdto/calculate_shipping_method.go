package requestdto

import "halodeksik-be/app/entity"

type CheckoutItemRequest struct {
	PharmacyProductId int64 `json:"pharmacy_product_id" validate:"required,min=1"`
	Quantity          int32 `json:"quantity" validate:"required,min=1"`
}

func (r *CheckoutItemRequest) ToCheckoutItem() entity.CheckoutItem {
	return entity.CheckoutItem{
		PharmacyProductId: r.PharmacyProductId,
		Quantity:          r.Quantity,
	}
}

type CalculateShippingMethod struct {
	AddressId     int64                 `json:"address_id" validate:"required"`
	CheckoutItems []CheckoutItemRequest `json:"checkout_items" validate:"required,min=1,dive"`
}
