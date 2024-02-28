package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
)

type UserAddressRepository interface {
	FindAllByUserId(ctx context.Context, userId int64, param *queryparamdto.GetAllParams) ([]*entity.Address, error)
	CountFindAllUserId(ctx context.Context, userId int64, param *queryparamdto.GetAllParams) (int64, error)
	FindById(ctx context.Context, id int64) (*entity.Address, error)
	FindMainByUserId(ctx context.Context, userId int64) (*entity.Address, error)
	Create(ctx context.Context, address entity.Address) (*entity.Address, error)
	Update(ctx context.Context, address entity.Address) (*entity.Address, error)
	Delete(ctx context.Context, id int64) error
}

type UserAddressRepositoryImpl struct {
	db *sql.DB
}

func NewUserAddressRepositoryImpl(db *sql.DB) *UserAddressRepositoryImpl {
	return &UserAddressRepositoryImpl{db: db}
}

func (repo *UserAddressRepositoryImpl) FindMainByUserId(ctx context.Context, userId int64) (*entity.Address, error) {
	const findMainAddress = `
	SELECT id, name, address, sub_district, district, city, province, postal_code, latitude, longitude, status, profile_id FROM addresses where profile_id = $1 AND status = 1 AND deleted_at IS NULL
	`
	row := repo.db.QueryRowContext(ctx, findMainAddress, userId)
	var address entity.Address
	err := row.Scan(
		&address.Id,
		&address.Name,
		&address.Address,
		&address.SubDistrict,
		&address.District,
		&address.CityId,
		&address.ProvinceId,
		&address.PostalCode,
		&address.Latitude,
		&address.Longitude,
		&address.Status,
		&address.ProfileId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}
	return &address, nil
}

func (repo *UserAddressRepositoryImpl) FindAllByUserId(ctx context.Context, userId int64, param *queryparamdto.GetAllParams) ([]*entity.Address, error) {
	const getAddresses = `
	SELECT id, name, address, sub_district, district, city, province, postal_code, latitude, longitude, status, profile_id FROM addresses
	WHERE profile_id = $1 AND deleted_at IS NULL
	`
	indexPreparedStatement := 1

	query, values := buildQuery(getAddresses, &entity.Address{}, param, true, true, indexPreparedStatement)
	values = util.AppendAtIndex(values, 0, interface{}(userId))

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Address, 0)

	for rows.Next() {
		var address entity.Address
		if err := rows.Scan(
			&address.Id, &address.Name, &address.Address, &address.SubDistrict, &address.District, &address.CityId,
			&address.ProvinceId, &address.PostalCode, &address.Latitude, &address.Longitude, &address.Status, &address.ProfileId,
		); err != nil {
			return nil, err
		}
		items = append(items, &address)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *UserAddressRepositoryImpl) CountFindAllUserId(ctx context.Context, userId int64, param *queryparamdto.GetAllParams) (int64, error) {
	const getAddresses = `
	SELECT count(id) FROM addresses 
	WHERE profile_id = $1 AND deleted_at is null
	`
	query, values := buildQuery(getAddresses, &entity.Address{}, param, false, false)
	values = util.AppendAtIndex(values, 0, interface{}(userId))

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

func (repo *UserAddressRepositoryImpl) Create(ctx context.Context, address entity.Address) (*entity.Address, error) {
	const createAddress = `
	insert into addresses(name, address, sub_district, district, city, province, postal_code, latitude, longitude, status, profile_id, created_at, updated_at)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, now(), now())
	RETURNING id, name, address, sub_district, district, city, province, postal_code, latitude, longitude, status, profile_id
	`
	row := repo.db.QueryRowContext(ctx, createAddress,
		address.Name,
		address.Address,
		address.SubDistrict,
		address.District,
		address.CityId,
		address.ProvinceId,
		address.PostalCode,
		address.Latitude,
		address.Longitude,
		address.Status,
		address.ProfileId,
	)
	var createdAddress entity.Address
	err := row.Scan(
		&createdAddress.Id,
		&createdAddress.Name,
		&createdAddress.Address,
		&createdAddress.SubDistrict,
		&createdAddress.District,
		&createdAddress.CityId,
		&createdAddress.ProvinceId,
		&createdAddress.PostalCode,
		&createdAddress.Latitude,
		&createdAddress.Longitude,
		&createdAddress.Status,
		&createdAddress.ProfileId,
	)
	return &createdAddress, err
}

func (repo *UserAddressRepositoryImpl) FindById(ctx context.Context, id int64) (*entity.Address, error) {
	const findAddressById = `
	SELECT id, name, address, sub_district, district, city, province, postal_code, latitude, longitude, status, profile_id from addresses
	WHERE id = $1 AND deleted_at IS NULL
	`
	row := repo.db.QueryRowContext(ctx, findAddressById, id)
	var address entity.Address
	err := row.Scan(
		&address.Id,
		&address.Name,
		&address.Address,
		&address.SubDistrict,
		&address.District,
		&address.CityId,
		&address.ProvinceId,
		&address.PostalCode,
		&address.Latitude,
		&address.Longitude,
		&address.Status,
		&address.ProfileId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}
	return &address, nil
}

func (repo *UserAddressRepositoryImpl) Update(ctx context.Context, address entity.Address) (*entity.Address, error) {
	const updateAddress = `
	UPDATE addresses SET name = $1, address = $2, sub_district = $3, district = $4,
						 city = $5, province = $6, postal_code = $7, latitude = $8, longitude = $9, status = $10
	WHERE id = $11 AND deleted_at IS NULL 
	RETURNING id, name, address, sub_district, district, city, province, postal_code, latitude, longitude, status, profile_id
	`

	row := repo.db.QueryRowContext(ctx, updateAddress,
		address.Name,
		address.Address,
		address.SubDistrict,
		address.District,
		address.CityId,
		address.ProvinceId,
		address.PostalCode,
		address.Latitude,
		address.Longitude,
		address.Status,
		address.Id,
	)

	var createdAddress entity.Address
	err := row.Scan(
		&createdAddress.Id,
		&createdAddress.Name,
		&createdAddress.Address,
		&createdAddress.SubDistrict,
		&createdAddress.District,
		&createdAddress.CityId,
		&createdAddress.ProvinceId,
		&createdAddress.PostalCode,
		&createdAddress.Latitude,
		&createdAddress.Longitude,
		&createdAddress.Status,
		&createdAddress.ProfileId,
	)
	return &createdAddress, err

}

func (repo *UserAddressRepositoryImpl) Delete(ctx context.Context, id int64) error {
	const deleteAddress = `
	UPDATE addresses SET deleted_at = now(), status = 2 WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := repo.db.ExecContext(ctx, deleteAddress, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.ErrRecordNotFound
		}
		return err
	}

	return err
}
