package entity

import (
	"database/sql"
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"time"
)

type DrugClassification struct {
	Id        int64        `json:"id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

func (e *DrugClassification) GetEntityName() string {
	return "drug_classifications"
}

func (e *DrugClassification) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(e).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (e *DrugClassification) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", e.GetEntityName(), e.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (e *DrugClassification) ToResponse() *responsedto.DrugClassificationResponse {
	if e == nil {
		return nil
	}
	return &responsedto.DrugClassificationResponse{
		Id:   e.Id,
		Name: e.Name,
	}
}
