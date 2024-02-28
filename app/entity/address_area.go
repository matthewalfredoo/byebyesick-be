package entity

import (
	"database/sql"
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"time"
)

type Province struct {
	Id        int64        `json:"id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"-"`
}

func (e *Province) GetEntityName() string {
	return "provinces"
}

func (e *Province) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(e).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (e *Province) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", e.GetEntityName(), e.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (e *Province) ToResponse() *responsedto.ProvinceResponse {
	return &responsedto.ProvinceResponse{
		Id:   e.Id,
		Name: e.Name,
	}
}

type City struct {
	Id         int64        `json:"id"`
	Name       string       `json:"name"`
	ProvinceId int64        `json:"province_id"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	DeletedAt  sql.NullTime `json:"-"`
}

func (e *City) GetEntityName() string {
	return "cities"
}

func (e *City) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(e).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (e *City) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", e.GetEntityName(), e.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (e *City) ToResponse() *responsedto.CityResponse {
	return &responsedto.CityResponse{
		Id:         e.Id,
		Name:       e.Name,
		ProvinceId: e.ProvinceId,
	}
}
