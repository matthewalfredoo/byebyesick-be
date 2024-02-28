package appconstant

const (
	BroadcastChannelBufferSize = 16

	ConsultationSessionStatusOngoing int64 = 1
	ConsultationSessionStatusEnded   int64 = 2

	DataEncodingBase64  = "base64"
	DataTypeImage       = "image"
	DataTypeApplication = "application"

	MessageTypeRegular = 1
	MessageTypeAlert   = 2

	MessageDoctorCreateLeaveSick = "Sick leave certificate has been issued"
	MessageDoctorUpdateLeaveSick = "Sick leave certificate has been updated"

	MessageDoctorCreatePrescription = "Prescription has been issued"
	MessageDoctorUpdatePrescription = "Prescription has been updated"
)
