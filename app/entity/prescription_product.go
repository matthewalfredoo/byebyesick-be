package entity

import (
	"database/sql"
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"time"
)

type PrescriptionProduct struct {
	Id             int64        `json:"id"`
	PrescriptionId int64        `json:"prescription_id"`
	ProductId      int64        `json:"product_id"`
	Note           string       `json:"note"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      sql.NullTime `json:"deleted_at"`
	Product        *Product
}

func (e *PrescriptionProduct) GetEntityName() string {
	return "prescription_products"
}

func (e *PrescriptionProduct) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(e).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (e *PrescriptionProduct) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", e.GetEntityName(), e.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (e *PrescriptionProduct) ToResponse() *responsedto.PrescriptionProductResponse {
	if e == nil {
		return nil
	}
	return &responsedto.PrescriptionProductResponse{
		Id:             e.Id,
		PrescriptionId: e.PrescriptionId,
		ProductId:      e.ProductId,
		Note:           e.Note,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
		Product:        e.Product.ToProductResponse(),
	}
}
