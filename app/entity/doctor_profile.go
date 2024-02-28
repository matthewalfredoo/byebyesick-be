package entity

import (
	"database/sql"
	"fmt"
	"github.com/shopspring/decimal"
	"halodeksik-be/app/appconstant"
	"reflect"
	"time"
)

type DoctorProfile struct {
	UserId                 int64           `json:"user_id"`
	Name                   string          `json:"name"`
	ProfilePhoto           string          `json:"profile_photo"`
	StartingYear           int32           `json:"starting_year"`
	DoctorCertificate      string          `json:"doctor_certificate"`
	DoctorSpecializationId int64           `json:"doctor_specialization_id"`
	ConsultationFee        decimal.Decimal `json:"consultation_fee"`
	IsOnline               bool            `json:"is_online"`
	CreatedAt              time.Time       `json:"created_at"`
	UpdatedAt              time.Time       `json:"updated_at"`
	DeletedAt              sql.NullTime    `json:"deleted_at"`
	DoctorSpecialization   *DoctorSpecialization
}

func (u *DoctorProfile) GetEntityName() string {
	return "doctor_profiles"
}

func (u *DoctorProfile) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(u).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (u *DoctorProfile) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", u.GetEntityName(), u.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (u *DoctorProfile) GetProfile() *Profile {
	if u == nil {
		return nil
	}
	return &Profile{
		UserId:       u.UserId,
		Name:         u.Name,
		ProfilePhoto: u.ProfilePhoto,
	}
}
