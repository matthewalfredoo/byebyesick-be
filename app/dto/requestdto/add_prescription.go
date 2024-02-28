package requestdto

import "halodeksik-be/app/entity"

type AddPrescription struct {
	SessionId            int64                        `json:"session_id" validate:"required"`
	Symptoms             string                       `json:"symptoms" validate:"required"`
	Diagnosis            string                       `json:"diagnosis" validate:"required"`
	PrescriptionProducts []AddEditPrescriptionProduct `json:"prescription_products" validate:"required,dive"`
}

func (r AddPrescription) ToPrescription() entity.Prescription {
	prescriptionProducts := make([]*entity.PrescriptionProduct, 0)
	for _, prescriptionProduct := range r.PrescriptionProducts {
		prescriptionProducts = append(prescriptionProducts, prescriptionProduct.ToPrescriptionProduct())
	}

	return entity.Prescription{
		SessionId:            r.SessionId,
		Symptoms:             r.Symptoms,
		Diagnosis:            r.Diagnosis,
		PrescriptionProducts: prescriptionProducts,
	}
}
