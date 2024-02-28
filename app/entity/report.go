package entity

type SellReportMonthly struct {
	Month     int32 `json:"month"`
	TotalSell int64 `json:"total_sell"`
}

type SellReport struct {
	PharmacyAdminEmail string `json:"pharmacy_admin_email"`
	PharmacyId         int64  `json:"pharmacy_id"`
	PharmacyName       string `json:"pharmacy_name"`
	TotalSells         int64  `json:"total_sells"`
	Year               int64  `json:"year"`
}
