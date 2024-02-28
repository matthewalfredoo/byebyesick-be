package responsedto

type AddressResponse struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	SubDistrict string `json:"sub_district"`
	District    string `json:"district"`
	CityId      int64  `json:"city_id"`
	ProvinceId  int64  `json:"province_id"`
	PostalCode  string `json:"postal_code"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	Status      int32  `json:"status"`
	ProfileId   int64  `json:"profile_id"`
}
