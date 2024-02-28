package entity

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"time"
)

type OrderStatus struct {
	Id        int64              `json:"id"`
	Name      string             `json:"name"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
}

func (u *OrderStatus) ToOrderStatusResponse() responsedto.OrderStatusResponse {
	return responsedto.OrderStatusResponse{
		Id:   u.Id,
		Name: u.Name,
	}
}

func (u *OrderStatus) GetEntityName() string {
	return "order_statuses"
}

func (u *OrderStatus) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(u).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (u *OrderStatus) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", u.GetEntityName(), u.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

type OrderStatusLog struct {
	Id            int64        `json:"id"`
	OrderId       int64        `json:"order_id"`
	OrderStatusId int64        `json:"order_status_id"`
	IsLatest      bool         `json:"is_latest"`
	Description   string       `json:"description"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	DeletedAt     sql.NullTime `json:"deleted_at"`
	OrderStatus   *OrderStatus
}

func (l *OrderStatusLog) ToOrderStatusLogResponse() responsedto.OrderLogResponse {
	return responsedto.OrderLogResponse{
		Id:            l.Id,
		OrderId:       l.OrderId,
		OrderStatusId: l.OrderStatusId,
		IsLatest:      l.IsLatest,
		Description:   l.Description,
	}
}

func (l *OrderStatusLog) ToOrderStatusHistoryResponse() responsedto.OrderHistoryResponse {
	return responsedto.OrderHistoryResponse{
		Id:              l.Id,
		OrderStatusName: l.OrderStatus.Name,
		Date:            l.CreatedAt,
		IsLatest:        l.IsLatest,
		Description:     l.Description,
	}
}
