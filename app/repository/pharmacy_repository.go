package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
)

type PharmacyRepository interface {
	Create(ctx context.Context, pharmacy entity.Pharmacy) (*entity.Pharmacy, error)
	FindById(ctx context.Context, id int64) (*entity.Pharmacy, error)
	FindAll(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.Pharmacy, error)
	CountFindAll(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error)
	Update(ctx context.Context, pharmacy entity.Pharmacy) (*entity.Pharmacy, error)
	Delete(ctx context.Context, id int64) error
}

type PharmacyRepositoryImpl struct {
	db *sql.DB
}

func NewPharmacyRepository(db *sql.DB) *PharmacyRepositoryImpl {
	return &PharmacyRepositoryImpl{db: db}
}

func (repo *PharmacyRepositoryImpl) Create(ctx context.Context, pharmacy entity.Pharmacy) (*entity.Pharmacy, error) {
	const create = `INSERT INTO pharmacies(name, address, sub_district, district, city, province, postal_code, latitude, longitude, pharmacist_name, pharmacist_license_no, pharmacist_phone_no, operational_hours, operational_days, pharmacy_admin_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING id, name, address, sub_district, district, city, province, postal_code, latitude, longitude, pharmacist_name, pharmacist_license_no, pharmacist_phone_no, operational_hours, operational_days, pharmacy_admin_id, created_at, updated_at, deleted_at
`

	row := repo.db.QueryRowContext(ctx, create,
		pharmacy.Name,
		pharmacy.Address, pharmacy.SubDistrict, pharmacy.District, pharmacy.CityId, pharmacy.ProvinceId, pharmacy.PostalCode,
		pharmacy.Latitude, pharmacy.Longitude,
		pharmacy.PharmacistName, pharmacy.PharmacistLicenseNo, pharmacy.PharmacistPhoneNo,
		pharmacy.OperationalHours, pharmacy.OperationalDays, pharmacy.PharmacyAdminId,
	)
	var created entity.Pharmacy
	err := row.Scan(
		&created.Id, &created.Name,
		&created.Address, &created.SubDistrict, &created.District, &created.CityId, &created.ProvinceId, &created.PostalCode,
		&created.Latitude, &created.Longitude,
		&created.PharmacistName, &created.PharmacistLicenseNo, &created.PharmacistPhoneNo,
		&created.OperationalHours, &created.OperationalDays, &created.PharmacyAdminId,
		&created.CreatedAt, &created.UpdatedAt, &created.DeletedAt,
	)
	return &created, err
}

func (repo *PharmacyRepositoryImpl) FindById(ctx context.Context, id int64) (*entity.Pharmacy, error) {
	var getById = `SELECT id, name, address, sub_district, district, city, province, postal_code, latitude, longitude, pharmacist_name, pharmacist_license_no, pharmacist_phone_no, operational_hours, operational_days, pharmacy_admin_id
		FROM pharmacies
		WHERE id = $1 AND deleted_at IS NULL`

	row := repo.db.QueryRowContext(ctx, getById, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var pharmacy entity.Pharmacy
	err := row.Scan(
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

	return &pharmacy, nil
}

func (repo *PharmacyRepositoryImpl) FindAll(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.Pharmacy, error) {
	initQuery := `SELECT id, name, address, sub_district, district, city, province, postal_code, latitude, longitude, pharmacist_name, pharmacist_license_no, pharmacist_phone_no, operational_hours, operational_days, pharmacy_admin_id FROM pharmacies WHERE deleted_at IS NULL `
	query, values := buildQuery(initQuery, &entity.Pharmacy{}, param, true, true)

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Pharmacy, 0)
	for rows.Next() {
		var pharmacy entity.Pharmacy
		if err := rows.Scan(
			&pharmacy.Id, &pharmacy.Name,
			&pharmacy.Address, &pharmacy.SubDistrict, &pharmacy.District, &pharmacy.CityId, &pharmacy.ProvinceId, &pharmacy.PostalCode, &pharmacy.Latitude, &pharmacy.Longitude,
			&pharmacy.PharmacistName, &pharmacy.PharmacistLicenseNo, &pharmacy.PharmacistPhoneNo,
			&pharmacy.OperationalHours, &pharmacy.OperationalDays, &pharmacy.PharmacyAdminId,
		); err != nil {
			return nil, err
		}
		items = append(items, &pharmacy)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *PharmacyRepositoryImpl) CountFindAll(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error) {
	initQuery := `SELECT count(id) FROM pharmacies WHERE deleted_at IS NULL `
	query, values := buildQuery(initQuery, &entity.Pharmacy{}, param, false, false)

	var totalItems int64

	row := repo.db.QueryRowContext(ctx, query, values...)
	if row.Err() != nil {
		return totalItems, row.Err()
	}

	if err := row.Scan(
		&totalItems,
	); err != nil {
		return totalItems, err
	}

	return totalItems, nil
}

func (repo *PharmacyRepositoryImpl) Update(ctx context.Context, pharmacy entity.Pharmacy) (*entity.Pharmacy, error) {
	const updateById = `
		UPDATE pharmacies
		SET name=$1, address=$2, sub_district=$3, district=$4, city=$5, province=$6, postal_code=$7,
		latitude=$8, longitude=$9,
		pharmacist_name=$10, pharmacist_license_no=$11, pharmacist_phone_no=$12,
		operational_hours=$13, operational_days=$14, pharmacy_admin_id=$15, updated_at = now()
		WHERE id = $16
		RETURNING id, name, address, sub_district, district, city, province, postal_code, latitude, longitude, pharmacist_name, pharmacist_license_no, pharmacist_phone_no, operational_hours, operational_days, pharmacy_admin_id
		`

	row := repo.db.QueryRowContext(ctx, updateById,
		pharmacy.Name,
		pharmacy.Address, pharmacy.SubDistrict, pharmacy.District, pharmacy.CityId, pharmacy.ProvinceId, pharmacy.PostalCode,
		pharmacy.Latitude, pharmacy.Longitude,
		pharmacy.PharmacistName, pharmacy.PharmacistLicenseNo, pharmacy.PharmacistPhoneNo,
		pharmacy.OperationalHours, pharmacy.OperationalDays, pharmacy.PharmacyAdminId, pharmacy.Id,
	)

	var updated entity.Pharmacy
	err := row.Scan(
		&updated.Id, &updated.Name,
		&updated.Address, &updated.SubDistrict, &updated.District, &updated.CityId, &updated.ProvinceId, &updated.PostalCode, &updated.Latitude, &updated.Longitude,
		&updated.PharmacistName, &updated.PharmacistLicenseNo, &updated.PharmacistPhoneNo,
		&updated.OperationalHours, &updated.OperationalDays, &updated.PharmacyAdminId,
	)
	return &updated, err
}

func (repo *PharmacyRepositoryImpl) Delete(ctx context.Context, id int64) error {
	const deleteById = `
		UPDATE pharmacies
		SET deleted_at = now()
		WHERE id = $1 AND deleted_at IS NULL
		`

	_, err := repo.db.ExecContext(ctx, deleteById, id)
	return err
}
