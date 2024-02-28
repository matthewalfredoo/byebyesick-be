package entity

import (
	"database/sql"
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"time"
)

type Address struct {
	Id          int64        `json:"id"`
	Name        string       `json:"name"`
	Address     string       `json:"address"`
	SubDistrict string       `json:"sub_district"`
	District    string       `json:"district"`
	CityId      int64        `json:"city_id"`
	ProvinceId  int64        `json:"province_id"`
	PostalCode  string       `json:"postal_code"`
	Latitude    string       `json:"latitude"`
	Longitude   string       `json:"longitude"`
	Status      int32        `json:"status"`
	ProfileId   int64        `json:"profile_id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at"`
}

func (p *Address) GetEntityName() string {
	return "addresses"
}

func (p *Address) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(p).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (p *Address) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", p.GetEntityName(), p.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (p *Address) ToAddressResponse() *responsedto.AddressResponse {

	return &responsedto.AddressResponse{
		Id:          p.Id,
		Name:        p.Name,
		Address:     p.Address,
		SubDistrict: p.SubDistrict,
		District:    p.District,
		CityId:      p.CityId,
		ProvinceId:  p.ProvinceId,
		PostalCode:  p.PostalCode,
		Latitude:    p.Latitude,
		Longitude:   p.Longitude,
		Status:      p.Status,
		ProfileId:   p.ProfileId,
	}

}
