package responsedto

type ManufacturerResponse struct {
	Id    int64  `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Image string `json:"image"`
}
