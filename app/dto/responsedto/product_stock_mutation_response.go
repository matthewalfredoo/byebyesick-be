package responsedto

import "time"

type ProductStockMutationResponse struct {
	Id                               int64                             `json:"id"`
	PharmacyProductId                int64                             `json:"pharmacy_product_id"`
	ProductStockMutationTypeId       int64                             `json:"product_stock_mutation_type_id"`
	Stock                            int32                             `json:"stock"`
	MutationDate                     time.Time                         `json:"mutation_date"`
	PharmacyProductResponse          *PharmacyProductResponse          `json:"pharmacy_product,omitempty"`
	ProductStockMutationTypeResponse *ProductStockMutationTypeResponse `json:"product_stock_mutation_type,omitempty"`
}
