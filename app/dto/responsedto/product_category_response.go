package responsedto

type ProductCategoryResponse struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
