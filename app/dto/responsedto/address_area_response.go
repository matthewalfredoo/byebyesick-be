package responsedto

type ProvinceResponse struct {
	Id   int64  `json:"province_id"`
	Name string `json:"province"`
}

type CityResponse struct {
	Id         int64  `json:"city_id"`
	Name       string `json:"city_name"`
	ProvinceId int64  `json:"province_id"`
}
