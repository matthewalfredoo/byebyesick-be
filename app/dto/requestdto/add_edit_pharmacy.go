package requestdto

import (
	"fmt"
	"halodeksik-be/app/entity"
	"strings"
)

type AddEditPharmacy struct {
	Name                  string   `json:"name" validate:"required"`
	Address               string   `json:"address" validate:"required"`
	SubDistrict           string   `json:"sub_district" validate:"required"`
	District              string   `json:"district" validate:"required"`
	CityId                int64    `json:"city_id" validate:"required"`
	ProvinceId            int64    `json:"province_id" validate:"required"`
	PostalCode            string   `json:"postal_code" validate:"required,number"`
	Latitude              string   `json:"latitude" validate:"required,latitude"`
	Longitude             string   `json:"longitude" validate:"required,longitude"`
	PharmacistName        string   `json:"pharmacist_name" validate:"required"`
	PharmacistLicenseNo   string   `json:"pharmacist_license_no" validate:"required"`
	PharmacistPhoneNo     string   `json:"pharmacist_phone_no" validate:"required"`
	OperationalHoursOpen  *int     `json:"operational_hours_open" validate:"required,number,min=0,max=23"`
	OperationalHoursClose *int     `json:"operational_hours_close" validate:"required,number,min=0,max=23,gtfield=OperationalHoursOpen"`
	OperationalDays       []string `json:"operational_days" validate:"required,min=1,max=7,unique,dive,oneof=mon tue wed thu fri sat sun"`
	PharmacyAdminId       int64    `json:"pharmacy_admin_id" validate:"required"`
}

func (r AddEditPharmacy) ToPharmacy() entity.Pharmacy {
	operationalDays := ""
	for _, day := range r.OperationalDays {
		operationalDays += day + ","
	}
	operationalDays = strings.TrimSuffix(operationalDays, ",")
	operationalHours := fmt.Sprintf("%d-%d", *r.OperationalHoursOpen, *r.OperationalHoursClose)

	return entity.Pharmacy{
		Name:                r.Name,
		Address:             r.Address,
		SubDistrict:         r.SubDistrict,
		District:            r.District,
		CityId:              r.CityId,
		ProvinceId:          r.ProvinceId,
		PostalCode:          r.PostalCode,
		Latitude:            r.Latitude,
		Longitude:           r.Longitude,
		PharmacistName:      r.PharmacistName,
		PharmacistLicenseNo: r.PharmacistLicenseNo,
		PharmacistPhoneNo:   r.PharmacistPhoneNo,
		OperationalHours:    operationalHours,
		OperationalDays:     operationalDays,
		PharmacyAdminId:     r.PharmacyAdminId,
	}
}
