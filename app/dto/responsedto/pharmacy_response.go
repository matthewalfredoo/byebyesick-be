package responsedto

type PharmacyResponse struct {
	Id                    int64    `json:"id,omitempty"`
	Name                  string   `json:"name,omitempty"`
	Address               string   `json:"address,omitempty"`
	SubDistrict           string   `json:"sub_district,omitempty"`
	District              string   `json:"district,omitempty"`
	CityId                int64    `json:"city_id,omitempty"`
	ProvinceId            int64    `json:"province_id,omitempty"`
	PostalCode            string   `json:"postal_code,omitempty"`
	Latitude              string   `json:"latitude,omitempty"`
	Longitude             string   `json:"longitude,omitempty"`
	PharmacistName        string   `json:"pharmacist_name,omitempty"`
	PharmacistLicenseNo   string   `json:"pharmacist_license_no,omitempty"`
	PharmacistPhoneNo     string   `json:"pharmacist_phone_no,omitempty"`
	OperationalHoursOpen  int      `json:"operational_hours_open,omitempty"`
	OperationalHoursClose int      `json:"operational_hours_close,omitempty"`
	OperationalDays       []string `json:"operational_days,omitempty"`
	PharmacyAdminId       int64    `json:"pharmacy_admin_id,omitempty"`
}

type PharmacyIdNameResponse struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
