package entity

import (
	"database/sql"
	"fmt"
	"github.com/shopspring/decimal"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"time"
)

type Order struct {
	Id               int64           `json:"id"`
	Date             time.Time       `json:"date"`
	PharmacyId       int64           `json:"pharmacy_id"`
	NoOfItems        int32           `json:"no_of_items"`
	PharmacyAddress  string          `json:"pharmacy_address"`
	ShippingMethodId int64           `json:"shipping_method_id"`
	ShippingCost     decimal.Decimal `json:"shipping_cost"`
	TotalPayment     decimal.Decimal `json:"total_payment"`
	TransactionId    int64           `json:"transaction_id"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        sql.NullTime    `json:"deleted_at"`
	OrderDetails     []*OrderDetail  `json:"order_details"`
	Pharmacy         *Pharmacy       `json:"pharmacy"`
	ShippingMethod   *ShippingMethod `json:"shipping_method"`
	LatestStatus     *OrderStatus    `json:"latest_status"`
	UserAddress      string          `json:"user_address"`
}

type OrderIds struct {
	UserId          int64
	PharmacyAdminId int64
}

func (u *Order) GetEntityName() string {
	return "orders"
}

func (u *Order) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(u).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (u *Order) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", u.GetEntityName(), u.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (u *Order) ToOrderResponse() responsedto.OrderResponse {
	var res []*responsedto.OrderDetailResponse
	for _, detail := range u.OrderDetails {
		conv := detail.ToOrderDetailResponse()
		res = append(res, &conv)
	}
	return responsedto.OrderResponse{
		PharmacyName:   u.Pharmacy.Name,
		ShippingMethod: u.ShippingMethod.Name,
		ShippingCost:   u.ShippingCost.String(),
		TotalPayment:   u.TotalPayment.String(),
		OrderDetails:   res,
	}
}

func (u *Order) ToOrderListResponse() responsedto.OrderListResponse {
	pharRes := u.Pharmacy.ToPharmacyIdNameResponse()
	statusRes := u.LatestStatus.ToOrderStatusResponse()
	return responsedto.OrderListResponse{
		Id:            u.Id,
		Pharmacy:      &pharRes,
		Date:          u.Date,
		NoOfItems:     u.NoOfItems,
		TotalPayment:  u.TotalPayment.String(),
		TransactionId: u.TransactionId,
		Status:        &statusRes,
	}
}

func (u *Order) ToOrderDetailResponse() responsedto.OrderDetailFullResponse {
	status := u.LatestStatus.ToOrderStatusResponse()
	shipping := u.ShippingMethod.ToIdNameResponse()
	pharmacy := u.Pharmacy.ToPharmacyIdNameResponse()

	var res []*responsedto.OrderDetailResponse
	for _, detail := range u.OrderDetails {
		conv := detail.ToOrderDetailResponse()
		res = append(res, &conv)
	}

	return responsedto.OrderDetailFullResponse{
		Id:                  u.Id,
		OrderStatusResponse: &status,
		Date:                u.Date,
		ShippingMethod:      shipping,
		ShippingCost:        u.ShippingCost.String(),
		Pharmacy:            &pharmacy,
		UserAddress:         u.UserAddress,
		OrderDetails:        res,
	}
}
