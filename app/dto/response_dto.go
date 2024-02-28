package dto

type ResponseDto struct {
	Data    any      `json:"data,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}
