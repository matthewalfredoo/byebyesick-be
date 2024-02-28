package entity

import (
	"database/sql"
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"time"
)

type ConsultationSession struct {
	Id                          int64        `json:"id"`
	UserId                      int64        `json:"user_id"`
	DoctorId                    int64        `json:"doctor_id"`
	ConsultationSessionStatusId int64        `json:"consultation_session_status_id"`
	CreatedAt                   time.Time    `json:"created_at"`
	UpdatedAt                   time.Time    `json:"updated_at"`
	DeletedAt                   sql.NullTime `json:"deleted_at"`
	ConsultationSessionStatus   *ConsultationSessionStatus
	UserProfile                 *UserProfile
	DoctorProfile               *DoctorProfile
	Prescription                *Prescription
	SickLeaveForm               *SickLeaveForm
	Message                     []*ConsultationMessage
}

func (e *ConsultationSession) GetEntityName() string {
	return "consultation_sessions"
}

func (e *ConsultationSession) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(e).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (e *ConsultationSession) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", e.GetEntityName(), e.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (e *ConsultationSession) ToResponse() *responsedto.ConsultationSessionResponse {
	if e == nil {
		return nil
	}

	messageResp := make([]*responsedto.WsConsultationMessage, 0)
	for _, message := range e.Message {
		messageResp = append(messageResp, message.ToWsMessage())
	}

	return &responsedto.ConsultationSessionResponse{
		Id:                          e.Id,
		UserId:                      e.UserId,
		DoctorId:                    e.DoctorId,
		ConsultationSessionStatusId: e.ConsultationSessionStatusId,
		CreatedAt:                   e.CreatedAt,
		UpdatedAt:                   e.UpdatedAt,
		ConsultationSessionStatus:   e.ConsultationSessionStatus.ToResponse(),
		UserProfile:                 e.UserProfile.GetProfile().ToResponse(),
		DoctorProfile:               e.DoctorProfile.GetProfile().ToResponse(),
		Prescription:                e.Prescription.ToResponse(),
		SickLeaveForm:               e.SickLeaveForm.ToResponse(),
		Message:                     messageResp,
	}
}
