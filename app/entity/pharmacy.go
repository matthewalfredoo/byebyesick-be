package entity

import (
	"database/sql"
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Pharmacy struct {
	Id                  int64        `json:"id"`
	Name                string       `json:"name"`
	Address             string       `json:"address"`
	SubDistrict         string       `json:"sub_district"`
	District            string       `json:"district"`
	CityId              int64        `json:"city_id"`
	ProvinceId          int64        `json:"province_id"`
	PostalCode          string       `json:"postal_code"`
	Latitude            string       `json:"latitude"`
	Longitude           string       `json:"longitude"`
	PharmacistName      string       `json:"pharmacist_name"`
	PharmacistLicenseNo string       `json:"pharmacist_license_no"`
	PharmacistPhoneNo   string       `json:"pharmacist_phone_no"`
	OperationalHours    string       `json:"operational_hours"`
	OperationalDays     string       `json:"operational_days"`
	PharmacyAdminId     int64        `json:"pharmacy_admin_id"`
	CreatedAt           time.Time    `json:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at"`
	DeletedAt           sql.NullTime `json:"-"`
}

func (p *Pharmacy) GetEntityName() string {
	return "pharmacies"
}

func (p *Pharmacy) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(p).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (p *Pharmacy) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", p.GetEntityName(), p.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (p *Pharmacy) ToPharmacyResponse() *responsedto.PharmacyResponse {
	if p == nil {
		return nil
	}

	var (
		operationalHoursOpen  int
		operationalHoursClose int
	)

	if p.OperationalHours != "" {
		operationalHours := strings.Split(p.OperationalHours, "-")
		operationalHoursOpen, _ = strconv.Atoi(operationalHours[0])
		operationalHoursClose, _ = strconv.Atoi(operationalHours[1])
	}

	return &responsedto.PharmacyResponse{
		Id:                    p.Id,
		Name:                  p.Name,
		Address:               p.Address,
		SubDistrict:           p.SubDistrict,
		District:              p.District,
		CityId:                p.CityId,
		ProvinceId:            p.ProvinceId,
		PostalCode:            p.PostalCode,
		Latitude:              p.Latitude,
		Longitude:             p.Longitude,
		PharmacistName:        p.PharmacistName,
		PharmacistLicenseNo:   p.PharmacistLicenseNo,
		PharmacistPhoneNo:     p.PharmacistPhoneNo,
		OperationalHoursOpen:  operationalHoursOpen,
		OperationalHoursClose: operationalHoursClose,
		OperationalDays:       strings.Split(p.OperationalDays, ","),
		PharmacyAdminId:       p.PharmacyAdminId,
	}
}

func (p *Pharmacy) ToPharmacyIdNameResponse() responsedto.PharmacyIdNameResponse {
	return responsedto.PharmacyIdNameResponse{
		Id:   p.Id,
		Name: p.Name,
	}
}
