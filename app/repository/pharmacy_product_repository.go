package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
)

type PharmacyProductRepository interface {
	Create(ctx context.Context, pharmacyProduct entity.PharmacyProduct) (*entity.PharmacyProduct, error)
	FindById(ctx context.Context, id int64) (*entity.PharmacyProduct, error)
	FindByIdJoinPharmacy(ctx context.Context, id int64) (*entity.PharmacyProduct, error)
	FindByIdJoinPharmacyAndProduct(ctx context.Context, id int64) (*entity.PharmacyProduct, error)
	FindAllJoinPharmacy(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.PharmacyProduct, error)
	CountFindAllJoinPharmacy(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error)

	FindByProductIdJoinPharmacy(ctx context.Context, productId int64, param *queryparamdto.GetAllParams) (*entity.PharmacyProduct, error)
	MaxTotalStocksByProductsId(ctx context.Context, productId int64, param *queryparamdto.GetAllParams) (int32, error)

	FindAllJoinProducts(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.PharmacyProduct, error)
	FindAllByProductId(ctx context.Context, productId int64) ([]*entity.PharmacyProduct, error)
	CountFindAll(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error)
	Update(ctx context.Context, pharmacyProduct entity.PharmacyProduct) (*entity.PharmacyProduct, error)
}

type PharmacyProductRepositoryImpl struct {
	db *sql.DB
}

func NewPharmacyProductRepository(db *sql.DB) *PharmacyProductRepositoryImpl {
	return &PharmacyProductRepositoryImpl{db: db}
}

func (repo *PharmacyProductRepositoryImpl) Create(ctx context.Context, pharmacyProduct entity.PharmacyProduct) (*entity.PharmacyProduct, error) {
	const create = `INSERT INTO pharmacy_products(pharmacy_id, product_id, is_active, price, stock)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, pharmacy_id, product_id, is_active, price, stock
`

	row := repo.db.QueryRowContext(ctx, create,
		pharmacyProduct.PharmacyId,
		pharmacyProduct.ProductId,
		pharmacyProduct.IsActive,
		pharmacyProduct.Price.String(),
		pharmacyProduct.Stock,
	)

	if row.Err() != nil {
		var errPgConn *pgconn.PgError
		if errors.As(row.Err(), &errPgConn) && errPgConn.Code == apperror.PgconnErrCodeUniqueConstraintViolation {
			return nil, apperror.ErrPharmacyProductUniqueConstraint
		}
		return nil, row.Err()
	}

	var created entity.PharmacyProduct
	err := row.Scan(
		&created.Id, &created.PharmacyId, &created.ProductId,
		&created.IsActive, &created.Price, &created.Stock,
	)
	return &created, err
}

func (repo *PharmacyProductRepositoryImpl) FindById(ctx context.Context, id int64) (*entity.PharmacyProduct, error) {
	getById := `
	SELECT id, pharmacy_id, product_id, is_active, price, stock
	FROM pharmacy_products
	WHERE id = $1 AND deleted_at IS NULL `

	row := repo.db.QueryRowContext(ctx, getById, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var pharmacyProducts entity.PharmacyProduct
	err := row.Scan(
		&pharmacyProducts.Id, &pharmacyProducts.PharmacyId, &pharmacyProducts.ProductId,
		&pharmacyProducts.IsActive, &pharmacyProducts.Price, &pharmacyProducts.Stock,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}

	return &pharmacyProducts, nil
}

func (repo *PharmacyProductRepositoryImpl) FindByIdJoinPharmacy(ctx context.Context, id int64) (*entity.PharmacyProduct, error) {
	getById := `
	SELECT pp.id, pp.pharmacy_id, pp.product_id, pp.is_active, pp.price, pp.stock,
	       p.id, p.name, p.address, p.sub_district, p.district, p.city, p.province, p.postal_code, p.latitude, p.longitude, p.pharmacist_name, p.pharmacist_license_no, p.pharmacist_phone_no, p.operational_hours, p.operational_days, p.pharmacy_admin_id
	FROM pharmacy_products pp
	INNER JOIN pharmacies p on pp.pharmacy_id = p.id
	WHERE pp.id = $1 AND pp.deleted_at IS NULL `

	row := repo.db.QueryRowContext(ctx, getById, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var (
		pharmacyProducts entity.PharmacyProduct
		pharmacy         entity.Pharmacy
	)
	err := row.Scan(
		&pharmacyProducts.Id, &pharmacyProducts.PharmacyId, &pharmacyProducts.ProductId,
		&pharmacyProducts.IsActive, &pharmacyProducts.Price, &pharmacyProducts.Stock,
		&pharmacy.Id, &pharmacy.Name,
		&pharmacy.Address, &pharmacy.SubDistrict, &pharmacy.District, &pharmacy.CityId, &pharmacy.ProvinceId, &pharmacy.PostalCode, &pharmacy.Latitude, &pharmacy.Longitude,
		&pharmacy.PharmacistName, &pharmacy.PharmacistLicenseNo, &pharmacy.PharmacistPhoneNo,
		&pharmacy.OperationalHours, &pharmacy.OperationalDays, &pharmacy.PharmacyAdminId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}
	pharmacyProducts.Pharmacy = &pharmacy

	return &pharmacyProducts, nil
}

func (repo *PharmacyProductRepositoryImpl) FindByIdJoinPharmacyAndProduct(ctx context.Context, id int64) (*entity.PharmacyProduct, error) {
	getById := `
	SELECT pharmacy_products.id, pharmacy_products.pharmacy_id, pharmacy_products.product_id, pharmacy_products.is_active, pharmacy_products.price, pharmacy_products.stock,
	       pharmacies.id, pharmacies.name, pharmacies.address, pharmacies.sub_district, pharmacies.district, pharmacies.city, pharmacies.province, pharmacies.postal_code, pharmacies.latitude, pharmacies.longitude, pharmacies.pharmacist_name, pharmacies.pharmacist_license_no, pharmacies.pharmacist_phone_no, pharmacies.operational_hours, pharmacies.operational_days, pharmacies.pharmacy_admin_id,
		   products.id, products.name, products.generic_name, products.content, products.manufacturer_id, products.description, products.drug_classification_id, products.product_category_id, products.drug_form, products.unit_in_pack, products.selling_unit, products.weight, products.length, products.width, products.height, products.image
	FROM pharmacy_products
	INNER JOIN pharmacies ON pharmacy_products.pharmacy_id = pharmacies.id
	INNER JOIN products ON pharmacy_products.product_id = products.id
	WHERE pharmacy_products.id = $1 AND pharmacy_products.deleted_at IS NULL `

	row := repo.db.QueryRowContext(ctx, getById, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var (
		pharmacyProducts entity.PharmacyProduct
		pharmacy         entity.Pharmacy
		product          entity.Product
	)
	err := row.Scan(
		&pharmacyProducts.Id, &pharmacyProducts.PharmacyId, &pharmacyProducts.ProductId,
		&pharmacyProducts.IsActive, &pharmacyProducts.Price, &pharmacyProducts.Stock,
		&pharmacy.Id, &pharmacy.Name,
		&pharmacy.Address, &pharmacy.SubDistrict, &pharmacy.District, &pharmacy.CityId, &pharmacy.ProvinceId, &pharmacy.PostalCode, &pharmacy.Latitude, &pharmacy.Longitude,
		&pharmacy.PharmacistName, &pharmacy.PharmacistLicenseNo, &pharmacy.PharmacistPhoneNo,
		&pharmacy.OperationalHours, &pharmacy.OperationalDays, &pharmacy.PharmacyAdminId,
		&product.Id, &product.Name, &product.GenericName, &product.Content, &product.ManufacturerId, &product.Description, &product.DrugClassificationId, &product.ProductCategoryId, &product.DrugForm, &product.UnitInPack, &product.SellingUnit, &product.Weight, &product.Length, &product.Width, &product.Height, &product.Image,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}
	pharmacyProducts.Pharmacy = &pharmacy
	pharmacyProducts.Product = &product

	return &pharmacyProducts, nil
}

func (repo *PharmacyProductRepositoryImpl) FindAllJoinPharmacy(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.PharmacyProduct, error) {
	initQuery := `
	SELECT pharmacy_products.id, pharmacy_products.pharmacy_id, pharmacy_products.product_id, pharmacy_products.is_active, pharmacy_products.price, pharmacy_products.stock,
	       pharmacies.id, pharmacies.name, pharmacies.address, pharmacies.sub_district, pharmacies.district, pharmacies.city, pharmacies.province, pharmacies.postal_code, pharmacies.latitude, pharmacies.longitude, pharmacies.pharmacist_name, pharmacies.pharmacist_license_no, pharmacies.pharmacist_phone_no, pharmacies.operational_hours, pharmacies.operational_days, pharmacies.pharmacy_admin_id
	FROM pharmacy_products
			INNER JOIN pharmacies ON pharmacy_products.pharmacy_id = pharmacies.id
	WHERE pharmacy_products.deleted_at IS NULL AND pharmacies.deleted_at IS NULL `

	query, values := buildQuery(initQuery, &entity.PharmacyProduct{}, param, true, true)

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.PharmacyProduct, 0)
	for rows.Next() {
		var (
			pharmacyProduct entity.PharmacyProduct
			pharmacy        entity.Pharmacy
		)
		if err := rows.Scan(
			&pharmacyProduct.Id, &pharmacyProduct.PharmacyId, &pharmacyProduct.ProductId, &pharmacyProduct.IsActive, &pharmacyProduct.Price, &pharmacyProduct.Stock,
			&pharmacy.Id, &pharmacy.Name,
			&pharmacy.Address, &pharmacy.SubDistrict, &pharmacy.District, &pharmacy.CityId, &pharmacy.ProvinceId, &pharmacy.PostalCode, &pharmacy.Latitude, &pharmacy.Longitude,
			&pharmacy.PharmacistName, &pharmacy.PharmacistLicenseNo, &pharmacy.PharmacistPhoneNo,
			&pharmacy.OperationalHours, &pharmacy.OperationalDays, &pharmacy.PharmacyAdminId,
		); err != nil {
			return nil, err
		}
		pharmacyProduct.Pharmacy = &pharmacy
		items = append(items, &pharmacyProduct)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *PharmacyProductRepositoryImpl) CountFindAllJoinPharmacy(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error) {
	initQuery := `SELECT count(pharmacy_products.id) FROM pharmacy_products
			INNER JOIN pharmacies ON pharmacy_products.pharmacy_id = pharmacies.id
	WHERE pharmacy_products.deleted_at IS NULL AND pharmacies.deleted_at IS NULL `
	query, values := buildQuery(initQuery, &entity.PharmacyProduct{}, param, false, false)

	var totalItems int64

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return totalItems, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&totalItems,
		); err != nil {
			return totalItems, err
		}
	}

	if err := rows.Err(); err != nil {
		return totalItems, err
	}
	return totalItems, nil
}

func (repo *PharmacyProductRepositoryImpl) FindByProductIdJoinPharmacy(ctx context.Context, productId int64, param *queryparamdto.GetAllParams) (*entity.PharmacyProduct, error) {
	initQuery := `SELECT pharmacy_products.id, pharmacy_products.pharmacy_id, pharmacy_products.product_id, 
       pharmacy_products.is_active, pharmacy_products.price, pharmacy_products.stock,
       pharmacies.id, pharmacies.name, pharmacies.address, pharmacies.sub_district, pharmacies.district, pharmacies.city, 
       pharmacies.province, pharmacies.postal_code, pharmacies.latitude, pharmacies.longitude, pharmacies.pharmacist_name, 
       pharmacies.pharmacist_license_no, pharmacies.pharmacist_phone_no, pharmacies.operational_hours, pharmacies.operational_days, pharmacies.pharmacy_admin_id
	FROM pharmacy_products
         INNER JOIN pharmacies ON pharmacy_products.pharmacy_id = pharmacies.id
	WHERE pharmacies.deleted_at IS NULL 
	  AND pharmacy_products.deleted_at IS NULL 
	  AND pharmacy_products.is_active = true
	  AND pharmacy_products.product_id = $1 `
	indexPreparedStatement := 1

	query, values := buildQuery(initQuery, &entity.PharmacyProduct{}, param, true, false, indexPreparedStatement)
	values = util.AppendAtIndex(values, 0, interface{}(productId))

	row := repo.db.QueryRowContext(ctx, query, values...)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var (
		pharmacyProduct entity.PharmacyProduct
		pharmacy        entity.Pharmacy
	)
	err := row.Scan(
		&pharmacyProduct.Id, &pharmacyProduct.PharmacyId, &pharmacyProduct.ProductId,
		&pharmacyProduct.IsActive, &pharmacyProduct.Price, &pharmacyProduct.Stock,
		&pharmacy.Id, &pharmacy.Name, &pharmacy.Address, &pharmacy.SubDistrict, &pharmacy.District, &pharmacy.CityId,
		&pharmacy.ProvinceId, &pharmacy.PostalCode, &pharmacy.Latitude, &pharmacy.Longitude, &pharmacy.PharmacistName,
		&pharmacy.PharmacistLicenseNo, &pharmacy.PharmacistPhoneNo, &pharmacy.OperationalHours, &pharmacy.OperationalDays, &pharmacy.PharmacyAdminId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}
	pharmacyProduct.Pharmacy = &pharmacy
	return &pharmacyProduct, nil
}

func (repo *PharmacyProductRepositoryImpl) MaxTotalStocksByProductsId(ctx context.Context, productId int64, param *queryparamdto.GetAllParams) (int32, error) {
	var (
		stock    int32
		maxStock int32
	)

	initQuery := `SELECT max(pharmacy_products.stock) 
	FROM pharmacy_products
	INNER JOIN pharmacies ON pharmacy_products.pharmacy_id = pharmacies.id
	WHERE pharmacies.deleted_at IS NULL 
	  AND pharmacy_products.deleted_at IS NULL 
	  AND pharmacy_products.is_active = true
	  AND pharmacy_products.product_id = $1 `
	indexPreparedStatement := 1

	query, values := buildQuery(initQuery, &entity.PharmacyProduct{}, param, false, false, indexPreparedStatement)
	values = util.AppendAtIndex(values, 0, interface{}(productId))

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if rows.Err() != nil {
		return maxStock, err
	}

	for rows.Next() {
		if err := rows.Scan(&stock); err != nil {
			return maxStock, err
		}
		if stock > maxStock {
			maxStock = stock
		}
	}
	if err := rows.Err(); err != nil {
		return maxStock, err
	}

	return maxStock, nil
}

func (repo *PharmacyProductRepositoryImpl) FindAllJoinProducts(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.PharmacyProduct, error) {
	initQuery := `
	SELECT pharmacy_products.id, pharmacy_products.pharmacy_id, pharmacy_products.product_id, pharmacy_products.is_active, pharmacy_products.price, pharmacy_products.stock,
		products.id, products.name, products.generic_name, products.content, products.manufacturer_id, products.description, products.drug_classification_id, products.product_category_id, products.drug_form, products.unit_in_pack, products.selling_unit, products.weight, products.length, products.width, products.height, products.image,
		product_categories.name, manufacturers.name, drug_classifications.name
	FROM pharmacy_products
			INNER JOIN products ON pharmacy_products.product_id = products.id
			INNER JOIN product_categories ON products.product_category_id = product_categories.id
	        INNER JOIN manufacturers ON products.manufacturer_id = manufacturers.id
	        INNER JOIN drug_classifications ON products.drug_classification_id = drug_classifications.id
	WHERE pharmacy_products.deleted_at IS NULL AND products.deleted_at IS NULL `

	query, values := buildQuery(initQuery, &entity.PharmacyProduct{}, param, true, true)

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.PharmacyProduct, 0)
	for rows.Next() {
		var (
			pharmacyProduct    entity.PharmacyProduct
			product            entity.Product
			productCategory    entity.ProductCategory
			manufacturer       entity.Manufacturer
			drugClassification entity.DrugClassification
		)
		if err := rows.Scan(
			&pharmacyProduct.Id, &pharmacyProduct.PharmacyId, &pharmacyProduct.ProductId, &pharmacyProduct.IsActive, &pharmacyProduct.Price, &pharmacyProduct.Stock,
			&product.Id, &product.Name, &product.GenericName, &product.Content, &product.ManufacturerId, &product.Description, &product.DrugClassificationId, &product.ProductCategoryId, &product.DrugForm, &product.UnitInPack, &product.SellingUnit, &product.Weight, &product.Length, &product.Width, &product.Height, &product.Image,
			&productCategory.Name, &manufacturer.Name, &drugClassification.Name,
		); err != nil {
			return nil, err
		}
		product.ProductCategory = &productCategory
		product.Manufacturer = &manufacturer
		product.DrugClassification = &drugClassification
		pharmacyProduct.Product = &product
		items = append(items, &pharmacyProduct)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *PharmacyProductRepositoryImpl) FindAllByProductId(ctx context.Context, productId int64) ([]*entity.PharmacyProduct, error) {
	const findAllByProductId = `SELECT pharmacy_products.id, pharmacy_id, product_id, is_active, price, stock
	FROM pharmacy_products INNER JOIN products ON pharmacy_products.product_id = products.id
	WHERE product_id = $1 AND public.pharmacy_products.deleted_at IS NULL AND products.deleted_at IS NULL`

	rows, err := repo.db.QueryContext(ctx, findAllByProductId, productId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pharmacyProducts := make([]*entity.PharmacyProduct, 0)
	for rows.Next() {
		var pharmacyProduct entity.PharmacyProduct
		if err := rows.Scan(
			&pharmacyProduct.Id, &pharmacyProduct.PharmacyId, &pharmacyProduct.ProductId,
			&pharmacyProduct.IsActive, &pharmacyProduct.Price, &pharmacyProduct.Stock,
		); err != nil {
			return nil, err
		}
		pharmacyProducts = append(pharmacyProducts, &pharmacyProduct)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return pharmacyProducts, nil
}

func (repo *PharmacyProductRepositoryImpl) CountFindAll(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error) {
	initQuery := `SELECT count(pharmacy_products.id) FROM pharmacy_products
	INNER JOIN products ON pharmacy_products.product_id = products.id
	INNER JOIN product_categories ON products.product_category_id = product_categories.id
	INNER JOIN manufacturers ON products.manufacturer_id = manufacturers.id
	INNER JOIN drug_classifications ON products.drug_classification_id = drug_classifications.id
	WHERE pharmacy_products.deleted_at IS NULL AND products.deleted_at IS NULL `
	query, values := buildQuery(initQuery, &entity.PharmacyProduct{}, param, false, false)

	var totalItems int64

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return totalItems, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&totalItems,
		); err != nil {
			return totalItems, err
		}
	}

	if err := rows.Err(); err != nil {
		return totalItems, err
	}
	return totalItems, nil
}

func (repo *PharmacyProductRepositoryImpl) Update(ctx context.Context, pharmacyProduct entity.PharmacyProduct) (*entity.PharmacyProduct, error) {
	updateById := `
		UPDATE pharmacy_products
		SET is_active = $1, price = $2
		WHERE id = $3
		RETURNING id, pharmacy_id, product_id, is_active, price, stock
	`

	row := repo.db.QueryRowContext(ctx, updateById,
		pharmacyProduct.IsActive,
		pharmacyProduct.Price,
		pharmacyProduct.Id,
	)
	var updated entity.PharmacyProduct
	err := row.Scan(
		&updated.Id,
		&updated.PharmacyId,
		&updated.ProductId,
		&updated.IsActive,
		&updated.Price,
		&updated.Stock,
	)
	return &updated, err
}
