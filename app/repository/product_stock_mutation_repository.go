package repository

import (
	"context"
	"database/sql"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
)

type ProductStockMutationRepository interface {
	Create(ctx context.Context, stockMutation entity.ProductStockMutation) (*entity.ProductStockMutation, error)
	FindAllJoin(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.ProductStockMutation, error)
	CountFindAllJoin(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error)
}

type ProductStockMutationRepositoryImpl struct {
	db *sql.DB
}

func NewProductStockMutationRepositoryImpl(db *sql.DB) *ProductStockMutationRepositoryImpl {
	return &ProductStockMutationRepositoryImpl{db: db}
}

func (repo *ProductStockMutationRepositoryImpl) Create(ctx context.Context, stockMutation entity.ProductStockMutation) (*entity.ProductStockMutation, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	create := `INSERT INTO product_stock_mutations(pharmacy_product_id, product_stock_mutation_type_id, stock)
VALUES ($1, $2, $3)
RETURNING id, pharmacy_product_id, product_stock_mutation_type_id, stock, created_at
`
	row1 := tx.QueryRowContext(ctx, create,
		stockMutation.PharmacyProductId, stockMutation.ProductStockMutationTypeId, stockMutation.Stock,
	)

	var created entity.ProductStockMutation
	if err := row1.Scan(
		&created.Id, &created.PharmacyProductId, &created.ProductStockMutationTypeId, &created.Stock, &created.CreatedAt,
	); err != nil {
		return nil, err
	}

	stock := stockMutation.Stock
	if stockMutation.ProductStockMutationTypeId == appconstant.StockMutationTypeReduction {
		stock = 0 - stock
	}

	updateStock := `UPDATE pharmacy_products
SET stock = stock + $1
WHERE id = $2
	`

	if _, err := tx.ExecContext(ctx, updateStock,
		stock,
		stockMutation.PharmacyProductId,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &created, err
}

func (repo *ProductStockMutationRepositoryImpl) FindAllJoin(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.ProductStockMutation, error) {
	initQuery := `
	SELECT product_stock_mutations.id,
       pharmacy_product_id,
       product_stock_mutation_type_id,
       product_stock_mutations.stock AS stock,
       product_stock_mutations.created_at,
       product_stock_mutation_types.name AS mutation_type,
       pharmacies.name                   AS pharmacy_name,
       products.name                     AS product_name,
       products.generic_name             AS product_generic_name,
       products.content                  AS product_content,
       manufacturers.name                AS manufacturer
FROM product_stock_mutations
         INNER JOIN product_stock_mutation_types
                    ON product_stock_mutation_types.id = product_stock_mutations.product_stock_mutation_type_id
         INNER JOIN pharmacy_products ON pharmacy_products.id = product_stock_mutations.pharmacy_product_id
         INNER JOIN pharmacies ON pharmacies.id = pharmacy_products.pharmacy_id
         INNER JOIN products ON products.id = pharmacy_products.product_id
         INNER JOIN manufacturers ON products.manufacturer_id = manufacturers.id
	WHERE product_stock_mutations.deleted_at IS NULL 
`
	query, values := buildQuery(initQuery, &entity.ProductStockMutation{}, param, true, true)

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.ProductStockMutation, 0)
	for rows.Next() {
		var (
			stockMutation   entity.ProductStockMutation
			mutationType    entity.ProductStockMutationType
			pharmacyProduct entity.PharmacyProduct
			pharmacy        entity.Pharmacy
			product         entity.Product
			manufacturer    entity.Manufacturer
		)
		if err := rows.Scan(
			&stockMutation.Id, &stockMutation.PharmacyProductId, &stockMutation.ProductStockMutationTypeId, &stockMutation.Stock, &stockMutation.CreatedAt,
			&mutationType.Name,
			&pharmacy.Name,
			&product.Name,
			&product.GenericName,
			&product.Content,
			&manufacturer.Name,
		); err != nil {
			return nil, err
		}
		product.Manufacturer = &manufacturer
		pharmacyProduct.Pharmacy = &pharmacy
		pharmacyProduct.Product = &product
		stockMutation.ProductStockMutationType = &mutationType
		stockMutation.PharmacyProduct = &pharmacyProduct
		items = append(items, &stockMutation)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *ProductStockMutationRepositoryImpl) CountFindAllJoin(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error) {
	initQuery := `
	SELECT count(product_stock_mutations.id)
FROM product_stock_mutations
    	 INNER JOIN product_stock_mutation_types
                    ON product_stock_mutation_types.id = product_stock_mutations.product_stock_mutation_type_id
         INNER JOIN pharmacy_products ON pharmacy_products.id = product_stock_mutations.pharmacy_product_id
         INNER JOIN pharmacies ON pharmacies.id = pharmacy_products.pharmacy_id
         INNER JOIN products ON products.id = pharmacy_products.product_id
         INNER JOIN manufacturers ON products.manufacturer_id = manufacturers.id
	WHERE product_stock_mutations.deleted_at IS NULL 
`
	query, values := buildQuery(initQuery, &entity.ProductStockMutation{}, param, false, false)

	var (
		totalItems int64
		temp       int64
	)

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return totalItems, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&temp,
		); err != nil {
			return totalItems, err
		}
		totalItems++
	}

	if err := rows.Err(); err != nil {
		return totalItems, err
	}
	return totalItems, nil
}
