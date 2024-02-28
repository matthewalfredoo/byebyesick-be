package entity

import (
	"database/sql"
	"fmt"
	"halodeksik-be/app/appconstant"
	"reflect"
	"time"
)

type UserProfile struct {
	UserId       int64        `json:"user_id"`
	Name         string       `json:"name"`
	ProfilePhoto string       `json:"profile_photo"`
	DateOfBirth  time.Time    `json:"date_of_birth"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	DeletedAt    sql.NullTime `json:"deleted_at"`
}

func (u *UserProfile) GetEntityName() string {
	return "user_profiles"
}

func (u *UserProfile) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(u).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (u *UserProfile) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", u.GetEntityName(), u.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (u *UserProfile) GetProfile() *Profile {
	if u == nil {
		return nil
	}
	return &Profile{
		UserId:       u.UserId,
		Name:         u.Name,
		ProfilePhoto: u.ProfilePhoto,
	}
}
