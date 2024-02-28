package requestdto

import (
	"halodeksik-be/app/entity"
)

type AddEditManufacturer struct {
	Name string `form:"name" validate:"required"`
}

func (r AddEditManufacturer) ToManufacturer() entity.Manufacturer {
	return entity.Manufacturer{Name: r.Name}
}
