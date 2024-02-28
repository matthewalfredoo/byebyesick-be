package entity

type CheckoutItem struct {
	PharmacyProductId int64 `json:"pharmacy_product_id"`
	Quantity          int32 `json:"quantity"`
}
