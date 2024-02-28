package entity

import (
	"database/sql"
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"time"
)

type Prescription struct {
	Id                   int64        `json:"id"`
	SessionId            int64        `json:"session_id"`
	Symptoms             string       `json:"symptoms"`
	Diagnosis            string       `json:"diagnosis"`
	CreatedAt            time.Time    `json:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at"`
	DeletedAt            sql.NullTime `json:"deleted_at"`
	PrescriptionProducts []*PrescriptionProduct

	User   *User
	Doctor *User
}

func (e *Prescription) GetEntityName() string {
	return "prescriptions"
}

func (e *Prescription) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(e).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (e *Prescription) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", e.GetEntityName(), e.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (e *Prescription) ToResponse() *responsedto.PrescriptionResponse {
	if e == nil {
		return nil
	}

	prescriptionProducts := make([]*responsedto.PrescriptionProductResponse, 0)
	for _, prescriptionProduct := range e.PrescriptionProducts {
		prescriptionProducts = append(prescriptionProducts, prescriptionProduct.ToResponse())
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

	return &responsedto.PrescriptionResponse{
		Id:                   e.Id,
		SessionId:            e.SessionId,
		Symptoms:             e.Symptoms,
		Diagnosis:            e.Diagnosis,
		CreatedAt:            createdAtResponse,
		UpdatedAt:            updatedAtResponse,
		PrescriptionProducts: prescriptionProducts,
		User:                 userResponse,
		Doctor:               doctorResponse,
	}
}
