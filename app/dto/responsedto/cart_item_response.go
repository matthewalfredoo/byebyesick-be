package responsedto

type CartItemResponse struct {
	Id                              int64                            `json:"id"`
	UserId                          int64                            `json:"user_id"`
	ProductId                       int64                            `json:"product_id"`
	Quantity                        int32                            `json:"quantity"`
	ProductResponse                 *ProductResponse                 `json:"product,omitempty"`
	PharmacyProductCheckoutResponse *PharmacyProductCheckoutResponse `json:"pharmacy_product,omitempty"`
}
