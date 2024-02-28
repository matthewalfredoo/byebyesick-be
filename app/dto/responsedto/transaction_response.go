package responsedto

import "time"

type TransactionJoinResponse struct {
	Id                int64                      `json:"id"`
	Date              time.Time                  `json:"date"`
	PaymentProof      string                     `json:"payment_proof"`
	PaymentMethod     string                     `json:"payment_method"`
	Address           string                     `json:"address"`
	TotalPayment      string                     `json:"total_payment"`
	TransactionStatus *TransactionStatusResponse `json:"transaction_status"`
	Orders            []*OrderResponse           `json:"orders,omitempty"`
}

type TransactionResponse struct {
	Id                  int64     `json:"id"`
	Date                time.Time `json:"date"`
	PaymentProof        string    `json:"payment_proof"`
	TransactionStatusId int64     `json:"transaction_status_id"`
	PaymentMethodId     int64     `json:"payment_method"`
	Address             string    `json:"address"`
	TotalPayment        string    `json:"total_payment"`
}

type TransactionStatusResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
