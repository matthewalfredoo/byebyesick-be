package entity

import (
	"database/sql"
	"fmt"
	"github.com/shopspring/decimal"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"time"
)

type Product struct {
	Id                   int64        `json:"id"`
	Name                 string       `json:"name"`
	GenericName          string       `json:"generic_name"`
	Content              string       `json:"content"`
	ManufacturerId       int64        `json:"manufacturer_id"`
	Description          string       `json:"description"`
	DrugClassificationId int64        `json:"drug_classification_id"`
	ProductCategoryId    int64        `json:"product_category_id"`
	DrugForm             string       `json:"drug_form"`
	UnitInPack           string       `json:"unit_in_pack"`
	SellingUnit          string       `json:"selling_unit"`
	Weight               float64      `json:"weight"`
	Length               float64      `json:"length"`
	Width                float64      `json:"width"`
	Height               float64      `json:"height"`
	Image                string       `json:"image"`
	CreatedAt            time.Time    `json:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at"`
	DeletedAt            sql.NullTime `json:"-"`
	Manufacturer         *Manufacturer
	DrugClassification   *DrugClassification
	ProductCategory      *ProductCategory
	MinimumPrice         decimal.Decimal
	MaximumPrice         decimal.Decimal
}

func (p *Product) GetEntityName() string {
	return "products"
}

func (p *Product) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(p).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (p *Product) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", p.GetEntityName(), p.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (p *Product) ToProductResponse() *responsedto.ProductResponse {
	if p == nil {
		return nil
	}
	minimumPrice := p.MinimumPrice.String()
	maximumPrice := p.MaximumPrice.String()

	if p.MinimumPrice.IsZero() {
		minimumPrice = ""
	}
	if p.MaximumPrice.IsZero() {
		maximumPrice = ""
	}

	return &responsedto.ProductResponse{
		Id:                         p.Id,
		Name:                       p.Name,
		GenericName:                p.GenericName,
		Content:                    p.Content,
		ManufacturerId:             p.ManufacturerId,
		Description:                p.Description,
		DrugClassificationId:       p.DrugClassificationId,
		ProductCategoryId:          p.ProductCategoryId,
		DrugForm:                   p.DrugForm,
		UnitInPack:                 p.UnitInPack,
		SellingUnit:                p.SellingUnit,
		Weight:                     p.Weight,
		Length:                     p.Length,
		Width:                      p.Width,
		Height:                     p.Height,
		Image:                      p.Image,
		ManufacturerResponse:       p.Manufacturer.ToResponse(),
		DrugClassificationResponse: p.DrugClassification.ToResponse(),
		ProductCategoryResponse:    p.ProductCategory.ToResponse(),
		MinimumPrice:               minimumPrice,
		MaximumPrice:               maximumPrice,
	}
}
