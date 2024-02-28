package entity

import (
	"database/sql"
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"time"
)

type DoctorSpecialization struct {
	Id        int64        `json:"id"`
	Name      string       `json:"name"`
	Image     string       `json:"image"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

func (e *DoctorSpecialization) ToDoctorSpecializationResponse() responsedto.DoctorSpecializationResponse {
	return responsedto.DoctorSpecializationResponse{
		Id:   e.Id,
		Name: e.Name,
	}
}

func (e *DoctorSpecialization) GetEntityName() string {
	return "doctor_specializations"
}

func (e *DoctorSpecialization) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(e).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (e *DoctorSpecialization) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", e.GetEntityName(), e.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (e *DoctorSpecialization) ToResponse() *responsedto.SpecializationResponse {
	if e == nil {
		return nil
	}
	return &responsedto.SpecializationResponse{
		Id:    e.Id,
		Name:  e.Name,
		Image: e.Image,
	}
}
