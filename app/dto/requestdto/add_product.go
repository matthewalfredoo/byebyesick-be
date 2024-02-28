package requestdto

import (
	"halodeksik-be/app/entity"
	"mime/multipart"
)

type AddProduct struct {
	Name                 string                `json:"name" form:"name" validate:"required"`
	GenericName          string                `json:"generic_name" form:"generic_name" validate:"required"`
	Content              string                `json:"content" form:"content" validate:"required"`
	ManufacturerId       int64                 `json:"manufacturer_id" form:"manufacturer_id" validate:"required"`
	Description          string                `json:"description" form:"description" validate:"required"`
	DrugClassificationId int64                 `json:"drug_classification_id" form:"drug_classification_id" validate:"required"`
	ProductCategoryId    int64                 `json:"product_category_id" form:"product_category_id" validate:"required"`
	DrugForm             string                `json:"drug_form" form:"drug_form" validate:"required"`
	UnitInPack           string                `json:"unit_in_pack" form:"unit_in_pack" validate:"required"`
	SellingUnit          string                `json:"selling_unit" form:"selling_unit" validate:"required"`
	Weight               float64               `json:"weight" form:"weight" validate:"required"`
	Length               float64               `json:"length" form:"length" validate:"required"`
	Width                float64               `json:"width" form:"width" validate:"required"`
	Height               float64               `json:"height" form:"height" validate:"required"`
	Image                *multipart.FileHeader `json:"image" form:"image" validate:"required,filetype=png jpg jpeg,filesize=500"`
}

func (r AddProduct) ToProduct() entity.Product {
	return entity.Product{
		Name:                 r.Name,
		GenericName:          r.GenericName,
		Content:              r.Content,
		ManufacturerId:       r.ManufacturerId,
		Description:          r.Description,
		DrugClassificationId: r.DrugClassificationId,
		ProductCategoryId:    r.ProductCategoryId,
		DrugForm:             r.DrugForm,
		UnitInPack:           r.UnitInPack,
		SellingUnit:          r.SellingUnit,
		Weight:               r.Weight,
		Length:               r.Length,
		Width:                r.Width,
		Height:               r.Height,
	}
}
