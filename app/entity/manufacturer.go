package entity

import (
	"database/sql"
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"time"
)

type Manufacturer struct {
	Id        int64        `json:"id"`
	Name      string       `json:"name"`
	Image     string       `json:"image"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

func (e *Manufacturer) GetEntityName() string {
	return "manufacturers"
}

func (e *Manufacturer) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(e).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (e *Manufacturer) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", e.GetEntityName(), e.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (e *Manufacturer) ToResponse() *responsedto.ManufacturerResponse {
	if e == nil {
		return nil
	}
	return &responsedto.ManufacturerResponse{
		Id:    e.Id,
		Name:  e.Name,
		Image: e.Image,
	}
}
