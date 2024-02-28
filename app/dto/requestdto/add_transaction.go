package requestdto

type AddTransaction struct {
	AddressId int64      `json:"address_id" validate:"required"`
	Orders    []AddOrder `json:"orders" validate:"min=1,required,dive"`
}

type AddOrder struct {
	ShippingMethodId int64             `json:"shipping_method_id" validate:"required"`
	ShippingCost     string            `json:"shipping_cost" validate:"required,numeric,numericgt=0"`
	OrderDetails     []AddOrderDetails `json:"order_details" validate:"min=1,required,dive"`
}

type AddOrderDetails struct {
	Quantity          int32 `json:"quantity" validate:"required,min=1"`
	PharmacyProductId int64 `json:"pharmacy_product_id" validate:"required"`
}
