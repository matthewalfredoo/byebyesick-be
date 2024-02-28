package responsedto

import "time"

type OrderDetailResponse struct {
	Name        string `json:"name"`
	GenericName string `json:"generic_name"`
	Content     string `json:"content"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Price       string `json:"price"`
	Quantity    int32  `json:"quantity"`
}

type OrderDetailFullResponse struct {
	Id                  int64                         `json:"id"`
	OrderStatusResponse *OrderStatusResponse          `json:"order_status"`
	Date                time.Time                     `json:"date"`
	ShippingMethod      *ShippingMethodIdNameResponse `json:"shippingMethod"`
	ShippingCost        string                        `json:"shipping_cost"`
	Pharmacy            *PharmacyIdNameResponse       `json:"pharmacy"`
	UserAddress         string                        `json:"user_address"`
	OrderDetails        []*OrderDetailResponse        `json:"order_details"`
}
