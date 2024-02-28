package responsedto

import "time"

type ProductStockMutationRequestResponse struct {
	Id                                        int64                                      `json:"id"`
	PharmacyProductOriginId                   int64                                      `json:"pharmacy_product_origin_id"`
	PharmacyProductDestId                     int64                                      `json:"pharmacy_product_dest_id"`
	Stock                                     int32                                      `json:"stock"`
	ProductStockMutationRequestStatusId       int64                                      `json:"product_stock_mutation_request_status_id"`
	RequestDate                               time.Time                                  `json:"request_date"`
	PharmacyProductOriginResponse             *PharmacyProductResponse                   `json:"pharmacy_product_origin,omitempty"`
	PharmacyProductDestResponse               *PharmacyProductResponse                   `json:"pharmacy_product_dest,omitempty"`
	ProductStockMutationRequestStatusResponse *ProductStockMutationRequestStatusResponse `json:"product_stock_mutation_request_status,omitempty"`
}
