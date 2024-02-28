package entity

import (
	"database/sql"
	"fmt"
	"halodeksik-be/app/appconstant"
	"reflect"
	"time"
)

type VerificationToken struct {
	Id        int64        `json:"id"`
	Token     string       `json:"token"`
	IsValid   bool         `json:"is_valid"`
	ExpiredAt time.Time    `json:"expired_at"`
	Email     string       `json:"email"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

func (v *VerificationToken) GetEntityName() string {
	return "verification_tokens"
}

func (v *VerificationToken) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(v).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (v *VerificationToken) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", v.GetEntityName(), v.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}
