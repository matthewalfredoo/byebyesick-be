package responsedto

type PrescriptionSickLeaveUserProfileResponse struct {
	Email                string `json:"email,omitempty"`
	Name                 string `json:"name,omitempty"`
	DateOfBirth          string `json:"date_of_birth,omitempty"`
	DoctorSpecialization string `json:"doctor_specialization,omitempty"`
}
