package entity

import (
	"database/sql"
	"github.com/shopspring/decimal"
	"halodeksik-be/app/dto/responsedto"
	"time"
)

type ShippingMethod struct {
	Id        int64           `json:"id"`
	Name      string          `json:"name"`
	Cost      decimal.Decimal `json:"cost"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt sql.NullTime    `json:"deleted_at"`
}

func (e *ShippingMethod) ToResponse() *responsedto.ShippingMethodResponse {
	return &responsedto.ShippingMethodResponse{
		Id:   e.Id,
		Name: e.Name,
		Cost: e.Cost,
	}
}

func (e *ShippingMethod) ToIdNameResponse() *responsedto.ShippingMethodIdNameResponse {
	return &responsedto.ShippingMethodIdNameResponse{
		Id:   e.Id,
		Name: e.Name,
	}
}
