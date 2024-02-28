package responsedto

type ProductResponse struct {
	Id                         int64                       `json:"id,omitempty"`
	Name                       string                      `json:"name,omitempty"`
	GenericName                string                      `json:"generic_name,omitempty"`
	Content                    string                      `json:"content,omitempty"`
	ManufacturerId             int64                       `json:"manufacturer_id,omitempty"`
	Description                string                      `json:"description,omitempty"`
	DrugClassificationId       int64                       `json:"drug_classification_id,omitempty"`
	ProductCategoryId          int64                       `json:"product_category_id,omitempty"`
	DrugForm                   string                      `json:"drug_form,omitempty"`
	UnitInPack                 string                      `json:"unit_in_pack,omitempty"`
	SellingUnit                string                      `json:"selling_unit,omitempty"`
	Weight                     float64                     `json:"weight,omitempty"`
	Length                     float64                     `json:"length,omitempty"`
	Width                      float64                     `json:"width,omitempty"`
	Height                     float64                     `json:"height,omitempty"`
	Image                      string                      `json:"image,omitempty"`
	ManufacturerResponse       *ManufacturerResponse       `json:"manufacturer,omitempty"`
	DrugClassificationResponse *DrugClassificationResponse `json:"drug_classification,omitempty"`
	ProductCategoryResponse    *ProductCategoryResponse    `json:"product_category,omitempty"`
	MinimumPrice               string                      `json:"minimum_price,omitempty"`
	MaximumPrice               string                      `json:"maximum_price,omitempty"`
}
