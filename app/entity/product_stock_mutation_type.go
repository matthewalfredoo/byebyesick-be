package entity

import (
	"database/sql"
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"time"
)

type ProductStockMutationType struct {
	Id        int64        `json:"id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

func (e *ProductStockMutationType) GetEntityName() string {
	return "product_stock_mutation_types"
}

func (e *ProductStockMutationType) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(e).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (e *ProductStockMutationType) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", e.GetEntityName(), e.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (e *ProductStockMutationType) ToResponse() *responsedto.ProductStockMutationTypeResponse {
	if e == nil {
		return nil
	}
	return &responsedto.ProductStockMutationTypeResponse{
		Id:   e.Id,
		Name: e.Name,
	}
}
