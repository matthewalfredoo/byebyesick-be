package requestdto

import "halodeksik-be/app/entity"

type AddConsultationSession struct {
	DoctorId int64 `json:"doctor_id" validate:"required"`
}

func (r *AddConsultationSession) ToConsultationSessionUseCase() entity.ConsultationSession {
	return entity.ConsultationSession{
		DoctorId: r.DoctorId,
	}
}
