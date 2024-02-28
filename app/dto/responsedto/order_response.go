package responsedto

import "time"

type OrderResponse struct {
	PharmacyName   string                 `json:"pharmacy_name"`
	ShippingMethod string                 `json:"shipping_method"`
	ShippingCost   string                 `json:"shipping_cost"`
	TotalPayment   string                 `json:"total_payment"`
	OrderDetails   []*OrderDetailResponse `json:"order_details"`
}

type OrderStatusResponse struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type OrderLogResponse struct {
	Id            int64  `json:"id,omitempty"`
	OrderId       int64  `json:"order_id,omitempty"`
	OrderStatusId int64  `json:"order_status_id,omitempty"`
	IsLatest      bool   `json:"is_latest,omitempty"`
	Description   string `json:"description,omitempty"`
}

type OrderHistoryResponse struct {
	Id              int64     `json:"id,omitempty"`
	OrderStatusName string    `json:"order_status_name"`
	Date            time.Time `json:"date"`
	IsLatest        bool      `json:"is_latest"`
	Description     string    `json:"description"`
}

type OrderListResponse struct {
	Id            int64 `json:"id"`
	Pharmacy      *PharmacyIdNameResponse
	Date          time.Time `json:"date"`
	NoOfItems     int32     `json:"no_of_items"`
	TotalPayment  string    `json:"total_payment"`
	TransactionId int64     `json:"transaction_id"`
	Status        *OrderStatusResponse
}
