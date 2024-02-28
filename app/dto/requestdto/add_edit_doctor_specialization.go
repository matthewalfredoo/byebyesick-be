package requestdto

import (
	"halodeksik-be/app/entity"
)

type AddEditDoctorSpecialization struct {
	Name string `form:"name" validate:"required"`
}

func (r AddEditDoctorSpecialization) ToDoctorSpecialization() entity.DoctorSpecialization {
	return entity.DoctorSpecialization{Name: r.Name}
}
