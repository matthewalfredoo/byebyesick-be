package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
)

type ProductStockMutationRequestRepository interface {
	Create(ctx context.Context, mutationRequest entity.ProductStockMutationRequest) (*entity.ProductStockMutationRequest, error)
	FindByIdJoinPharmacyOrigin(ctx context.Context, id int64) (*entity.ProductStockMutationRequest, error)
	FindAllJoin(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.ProductStockMutationRequest, error)
	CountFindAllJoin(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error)
	Update(ctx context.Context, mutationRequest entity.ProductStockMutationRequest) (*entity.ProductStockMutationRequest, error)
}

type ProductStockMutationRequestRepositoryImpl struct {
	db *sql.DB
}

func NewProductStockMutationRequestRepositoryImpl(db *sql.DB) *ProductStockMutationRequestRepositoryImpl {
	return &ProductStockMutationRequestRepositoryImpl{db: db}
}

func (repo *ProductStockMutationRequestRepositoryImpl) Create(ctx context.Context, mutationRequest entity.ProductStockMutationRequest) (*entity.ProductStockMutationRequest, error) {
	create := `INSERT INTO product_stock_mutation_requests(pharmacy_product_origin_id, pharmacy_product_dest_id, stock, product_stock_mutation_request_status_id)
VALUES ($1, $2, $3, $4)
RETURNING id, pharmacy_product_origin_id, pharmacy_product_dest_id, stock, product_stock_mutation_request_status_id, created_at`

	row := repo.db.QueryRowContext(ctx, create,
		mutationRequest.PharmacyProductOriginId,
		mutationRequest.PharmacyProductDestId,
		mutationRequest.Stock,
		mutationRequest.ProductStockMutationRequestStatusId,
	)

	var created entity.ProductStockMutationRequest
	if err := row.Scan(
		&created.Id,
		&created.PharmacyProductOriginId,
		&created.PharmacyProductDestId,
		&created.Stock,
		&created.ProductStockMutationRequestStatusId,
		&created.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &created, nil
}

func (repo *ProductStockMutationRequestRepositoryImpl) FindByIdJoinPharmacyOrigin(ctx context.Context, id int64) (*entity.ProductStockMutationRequest, error) {
	getById := `SELECT psmr.id, psmr.pharmacy_product_origin_id, psmr.pharmacy_product_dest_id, psmr.stock, psmr.product_stock_mutation_request_status_id,
	ppo.stock,
    p.pharmacy_admin_id
FROM product_stock_mutation_requests psmr
INNER JOIN pharmacy_products ppo on psmr.pharmacy_product_origin_id = ppo.id
INNER JOIN pharmacies p on ppo.pharmacy_id = p.id
WHERE psmr.id = $1 AND psmr.deleted_at IS NULL`

	row := repo.db.QueryRowContext(ctx, getById, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var (
		mutationRequest entity.ProductStockMutationRequest
		pharmacyProduct entity.PharmacyProduct
		pharmacy        entity.Pharmacy
	)
	err := row.Scan(
		&mutationRequest.Id, &mutationRequest.PharmacyProductOriginId, &mutationRequest.PharmacyProductDestId, &mutationRequest.Stock, &mutationRequest.ProductStockMutationRequestStatusId,
		&pharmacyProduct.Stock,
		&pharmacy.PharmacyAdminId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}
	pharmacyProduct.Pharmacy = &pharmacy
	mutationRequest.PharmacyProductOrigin = &pharmacyProduct
	return &mutationRequest, err
}

func (repo *ProductStockMutationRequestRepositoryImpl) FindAllJoin(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.ProductStockMutationRequest, error) {
	initQuery := `
SELECT product_stock_mutation_requests.id,
       product_stock_mutation_requests.pharmacy_product_origin_id,
       product_stock_mutation_requests.pharmacy_product_dest_id,
       product_stock_mutation_requests.stock,
       product_stock_mutation_requests.product_stock_mutation_request_status_id,
       product_stock_mutation_requests.created_at AS request_date,
       po.name AS pharmacy_origin_name, pd.name AS pharmacy_destination_name,
       pr.name, pr.generic_name, pr.content, pr.manufacturer_id,
       m.name,
       psmrs.name AS status
FROM product_stock_mutation_requests
         INNER JOIN pharmacy_products ppo ON ppo.id = product_stock_mutation_requests.pharmacy_product_origin_id
         INNER JOIN pharmacy_products ppd ON ppd.id = product_stock_mutation_requests.pharmacy_product_dest_id
         INNER JOIN pharmacies po ON ppo.pharmacy_id = po.id
    	 INNER JOIN pharmacies pd ON ppd.pharmacy_id = pd.id
         INNER JOIN products pr ON ppd.product_id = pr.id
    	 INNER JOIN manufacturers m ON pr.manufacturer_id = m.id
         INNER JOIN product_stock_mutation_request_statuses psmrs
                    ON product_stock_mutation_requests.product_stock_mutation_request_status_id = psmrs.id
WHERE product_stock_mutation_requests.deleted_at IS NULL `

	query, values := buildQuery(initQuery, &entity.ProductStockMutationRequest{}, param, true, true)

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.ProductStockMutationRequest, 0)
	for rows.Next() {
		var (
			mutationRequest       entity.ProductStockMutationRequest
			mutationRequestStatus entity.ProductStockMutationRequestStatus
			pharmacyProductOrigin entity.PharmacyProduct
			pharmacyOrigin        entity.Pharmacy
			pharmacyProductDest   entity.PharmacyProduct
			pharmacyDest          entity.Pharmacy
			product               entity.Product
			manufacturer          entity.Manufacturer
		)
		if err := rows.Scan(
			&mutationRequest.Id, &mutationRequest.PharmacyProductOriginId, &mutationRequest.PharmacyProductDestId, &mutationRequest.Stock, &mutationRequest.ProductStockMutationRequestStatusId, &mutationRequest.CreatedAt,
			&pharmacyOrigin.Name,
			&pharmacyDest.Name,
			&product.Name, &product.GenericName, &product.Content, &product.ManufacturerId,
			&manufacturer.Name,
			&mutationRequestStatus.Name,
		); err != nil {
			return nil, err
		}
		product.Manufacturer = &manufacturer
		pharmacyProductOrigin.Pharmacy = &pharmacyOrigin
		pharmacyProductOrigin.Product = &product
		pharmacyProductDest.Pharmacy = &pharmacyDest
		pharmacyProductDest.Product = &product
		mutationRequest.PharmacyProductOrigin = &pharmacyProductOrigin
		mutationRequest.PharmacyProductDest = &pharmacyProductDest
		mutationRequest.ProductStockMutationRequestStatus = &mutationRequestStatus
		items = append(items, &mutationRequest)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *ProductStockMutationRequestRepositoryImpl) CountFindAllJoin(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error) {
	initQuery := `
	SELECT count(product_stock_mutation_requests.id)
FROM product_stock_mutation_requests
         INNER JOIN pharmacy_products ppo ON ppo.id = product_stock_mutation_requests.pharmacy_product_origin_id
         INNER JOIN pharmacy_products ppd ON ppd.id = product_stock_mutation_requests.pharmacy_product_dest_id
         INNER JOIN pharmacies p ON ppo.pharmacy_id = p.id
         INNER JOIN products pr ON ppd.product_id = pr.id
         INNER JOIN product_stock_mutation_request_statuses psmrs
                    ON product_stock_mutation_requests.product_stock_mutation_request_status_id = psmrs.id
WHERE product_stock_mutation_requests.deleted_at IS NULL `

	query, values := buildQuery(initQuery, &entity.ProductStockMutationRequest{}, param, false, false)

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

func (repo *ProductStockMutationRequestRepositoryImpl) Update(ctx context.Context, mutationRequest entity.ProductStockMutationRequest) (*entity.ProductStockMutationRequest, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	update := `UPDATE product_stock_mutation_requests
SET product_stock_mutation_request_status_id = $1
WHERE id = $2
RETURNING id, pharmacy_product_origin_id, pharmacy_product_dest_id, stock, product_stock_mutation_request_status_id, created_at
`

	row1 := tx.QueryRowContext(ctx, update,
		mutationRequest.ProductStockMutationRequestStatusId,
		mutationRequest.Id,
	)

	var updated entity.ProductStockMutationRequest
	if err := row1.Scan(
		&updated.Id, &updated.PharmacyProductOriginId, &updated.PharmacyProductDestId, &updated.Stock, &updated.ProductStockMutationRequestStatusId, &updated.CreatedAt,
	); err != nil {
		return nil, err
	}

	if mutationRequest.ProductStockMutationRequestStatusId != appconstant.StockMutationRequestStatusAccepted {
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return &updated, nil
	}

	getbyId := `SELECT stock FROM pharmacy_products WHERE id = $1 AND deleted_at IS NULL`
	row2 := tx.QueryRowContext(ctx, getbyId, mutationRequest.PharmacyProductOriginId)
	var pharmacyProductOrigin entity.PharmacyProduct
	if err := row2.Scan(
		&pharmacyProductOrigin.Stock,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}

	if pharmacyProductOrigin.Stock < mutationRequest.Stock {
		return nil, apperror.ErrInsufficientProductStock
	}

	createMutation := `INSERT INTO product_stock_mutations(pharmacy_product_id, product_stock_mutation_type_id, stock) VALUES ($1, $2, $3)`
	updateStock := `UPDATE pharmacy_products SET stock = stock + $1 WHERE id = $2`

	if _, err := tx.ExecContext(ctx, createMutation,
		mutationRequest.PharmacyProductOriginId, appconstant.StockMutationTypeReduction, mutationRequest.Stock,
	); err != nil {
		return nil, err
	}
	if _, err := tx.ExecContext(ctx, updateStock,
		0-mutationRequest.Stock, mutationRequest.PharmacyProductOriginId,
	); err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, createMutation,
		mutationRequest.PharmacyProductDestId, appconstant.StockMutationTypeAddition, mutationRequest.Stock,
	); err != nil {
		return nil, err
	}
	if _, err := tx.ExecContext(ctx, updateStock,
		mutationRequest.Stock, mutationRequest.PharmacyProductDestId,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &updated, err
}
