package requestdto

type RequestValidateLatLong struct {
	CityId     int64  `json:"city_id" validate:"required"`
	ProvinceId int64  `json:"province_id" validate:"required"`
	Latitude   string `json:"latitude" validate:"required,latitude"`
	Longitude  string `json:"longitude" validate:"required,longitude"`
}
