package requestdto

import (
	"halodeksik-be/app/entity"
)

type RequestAddress struct {
	Name        string `json:"name" validate:"required"`
	Address     string `json:"address" validate:"required"`
	SubDistrict string `json:"sub_district" validate:"required"`
	District    string `json:"district" validate:"required"`
	CityId      int64  `json:"city_id" validate:"required"`
	ProvinceId  int64  `json:"province_id" validate:"required"`
	PostalCode  string `json:"postal_code" validate:"required"`
	Latitude    string `json:"latitude" validate:"required,latitude"`
	Longitude   string `json:"longitude" validate:"required,longitude"`
}

func (r RequestAddress) ToAddress() entity.Address {

	return entity.Address{
		Name:        r.Name,
		Address:     r.Address,
		SubDistrict: r.SubDistrict,
		District:    r.District,
		CityId:      r.CityId,
		ProvinceId:  r.ProvinceId,
		PostalCode:  r.PostalCode,
		Latitude:    r.Latitude,
		Longitude:   r.Longitude,
	}
}
