package responsedto

import (
	"time"
)

type ConsultationSessionResponse struct {
	Id                          int64                              `json:"id"`
	UserId                      int64                              `json:"user_id"`
	DoctorId                    int64                              `json:"doctor_id"`
	ConsultationSessionStatusId int64                              `json:"consultation_session_status_id"`
	CreatedAt                   time.Time                          `json:"created_at"`
	UpdatedAt                   time.Time                          `json:"updated_at"`
	ConsultationSessionStatus   *ConsultationSessionStatusResponse `json:"consultation_session_status,omitempty"`
	UserProfile                 *ProfileResponse                   `json:"user,omitempty"`
	DoctorProfile               *ProfileResponse                   `json:"doctor,omitempty"`
	Prescription                *PrescriptionResponse              `json:"prescription,omitempty"`
	SickLeaveForm               *SickLeaveFormResponse             `json:"sick_leave_form,omitempty"`
	Message                     []*WsConsultationMessage           `json:"messages"`
}
