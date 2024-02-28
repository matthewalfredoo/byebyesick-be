package requestdto

import (
	"halodeksik-be/app/entity"
	"strings"
)

type AddEditProductCategory struct {
	Name string `json:"name" validate:"required"`
}

func (r AddEditProductCategory) ToProductCategory() entity.ProductCategory {
	r.Name = strings.TrimSpace(r.Name)
	return entity.ProductCategory{
		Name: r.Name,
	}
}
