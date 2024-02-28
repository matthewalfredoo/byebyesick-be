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

type OrderDetail struct {
	Id          int64           `json:"id"`
	OrderId     int64           `json:"order_id"`
	ProductId   int64           `json:"product_id"`
	Quantity    int32           `json:"quantity"`
	Name        string          `json:"name"`
	GenericName string          `json:"generic_name"`
	Content     string          `json:"content"`
	Description string          `json:"description"`
	Image       string          `json:"image"`
	Price       decimal.Decimal `json:"price"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   sql.NullTime    `json:"deleted_at"`
}

func (o *OrderDetail) GetEntityName() string {
	return "order_details"
}

func (o *OrderDetail) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(o).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (o *OrderDetail) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", o.GetEntityName(), o.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (o *OrderDetail) ToOrderDetailResponse() responsedto.OrderDetailResponse {
	return responsedto.OrderDetailResponse{
		Name:        o.Name,
		GenericName: o.GenericName,
		Content:     o.Content,
		Description: o.Description,
		Image:       o.Image,
		Price:       o.Price.String(),
		Quantity:    o.Quantity,
	}
}
