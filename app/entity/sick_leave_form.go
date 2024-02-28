package entity

import (
	"database/sql"
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"time"
)

type SickLeaveForm struct {
	Id           int64        `json:"id"`
	SessionId    int64        `json:"session_id"`
	StartingDate time.Time    `json:"starting_date"`
	EndingDate   time.Time    `json:"ending_date"`
	Description  string       `json:"description"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	DeletedAt    sql.NullTime `json:"deleted_at"`

	Prescription *Prescription
	User         *User
	Doctor       *User
}

func (e *SickLeaveForm) GetEntityName() string {
	return "sick_leave_forms"
}

func (e *SickLeaveForm) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(e).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (e *SickLeaveForm) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", e.GetEntityName(), e.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (e *SickLeaveForm) ToResponse() *responsedto.SickLeaveFormResponse {
	if e == nil {
		return nil
	}

	var (
		userResponse      *responsedto.PrescriptionSickLeaveUserProfileResponse
		doctorResponse    *responsedto.PrescriptionSickLeaveUserProfileResponse
		createdAtResponse string
		updatedAtResponse string
	)

	if e.User != nil && e.User.UserProfile != nil {
		userResponse = e.User.ToPrescriptionSickLeaveUserProfileResponse()
	}

	if e.Doctor != nil && e.Doctor.DoctorProfile != nil {
		doctorResponse = e.Doctor.ToPrescriptionSickLeaveUserProfileResponse()
	}

	if !e.CreatedAt.IsZero() {
		createdAtResponse = e.CreatedAt.Format(time.RFC3339)
	}

	if !e.UpdatedAt.IsZero() {
		updatedAtResponse = e.UpdatedAt.Format(time.RFC3339)
	}

	return &responsedto.SickLeaveFormResponse{
		Id:           e.Id,
		SessionId:    e.SessionId,
		StartingDate: e.StartingDate,
		EndingDate:   e.EndingDate,
		Description:  e.Description,
		CreatedAt:    createdAtResponse,
		UpdatedAt:    updatedAtResponse,
		Prescription: e.Prescription.ToResponse(),
		User:         userResponse,
		Doctor:       doctorResponse,
	}
}
