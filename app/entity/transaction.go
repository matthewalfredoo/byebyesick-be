package entity

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"time"
)

type Transaction struct {
	Id                  int64           `json:"id"`
	Date                time.Time       `json:"date"`
	PaymentProof        string          `json:"payment_proof"`
	TransactionStatusId int64           `json:"transaction_status_id"`
	PaymentMethodId     int64           `json:"payment_method_id"`
	Address             string          `json:"address"`
	UserId              int64           `json:"user_id"`
	TotalPayment        decimal.Decimal `json:"total_payment"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
	DeletedAt           sql.NullTime    `json:"deleted_at"`
	Orders              []*Order        `json:"orders"`
	TransactionStatus   *TransactionStatus
	PaymentMethod       *PaymentMethod
}

func (u *Transaction) ToTransactionJoinResponse() responsedto.TransactionJoinResponse {
	var res []*responsedto.OrderResponse
	statusRes := u.TransactionStatus.ToTransactionStatusResponse()
	for _, order := range u.Orders {
		conv := order.ToOrderResponse()
		res = append(res, &conv)
	}

	return responsedto.TransactionJoinResponse{
		Id:                u.Id,
		Date:              u.Date,
		PaymentProof:      u.PaymentProof,
		TransactionStatus: &statusRes,
		PaymentMethod:     u.PaymentMethod.Name,
		Address:           u.Address,
		TotalPayment:      u.TotalPayment.String(),
		Orders:            res,
	}
}

func (u *Transaction) ToTransactionResponse() responsedto.TransactionResponse {

	return responsedto.TransactionResponse{
		Id:                  u.Id,
		Date:                u.Date,
		PaymentProof:        u.PaymentProof,
		TransactionStatusId: u.TransactionStatusId,
		PaymentMethodId:     u.PaymentMethodId,
		Address:             u.Address,
		TotalPayment:        u.TotalPayment.String(),
	}
}

func (u *Transaction) GetEntityName() string {
	return "transactions"
}

func (u *Transaction) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(u).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (u *Transaction) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", u.GetEntityName(), u.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

type TransactionStatus struct {
	Id        int64              `json:"id"`
	Name      string             `json:"name"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
}

func (e *TransactionStatus) ToTransactionStatusResponse() responsedto.TransactionStatusResponse {
	return responsedto.TransactionStatusResponse{
		Id:   e.Id,
		Name: e.Name,
	}
}

type PaymentMethod struct {
	Id        int64              `json:"id"`
	Name      string             `json:"name"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
}

type TransactionPaymentAndStatus struct {
	TotalPayment        string `json:"total_payment"`
	TransactionStatusId int64  `json:"transaction_status_id"`
}
