package responsedto

type PrescriptionResponse struct {
	Id                   int64                                     `json:"id,omitempty"`
	SessionId            int64                                     `json:"session_id,omitempty"`
	Symptoms             string                                    `json:"symptoms,omitempty"`
	Diagnosis            string                                    `json:"diagnosis,omitempty"`
	CreatedAt            string                                    `json:"created_at,omitempty"`
	UpdatedAt            string                                    `json:"updated_at,omitempty"`
	PrescriptionProducts []*PrescriptionProductResponse            `json:"prescription_products,omitempty"`
	User                 *PrescriptionSickLeaveUserProfileResponse `json:"user,omitempty"`
	Doctor               *PrescriptionSickLeaveUserProfileResponse `json:"doctor,omitempty"`
}
