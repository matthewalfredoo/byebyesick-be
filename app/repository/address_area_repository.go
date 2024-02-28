package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
)

type AddressAreaRepository interface {
	FindAllProvince(ctx context.Context) ([]*entity.Province, error)
	FindCityById(ctx context.Context, cityId int64) (*entity.City, error)
	FindProvinceById(ctx context.Context, cityId int64) (*entity.Province, error)
	FindAllCities(ctx context.Context) ([]*entity.City, error)
}

type AddressAreaRepositoryImpl struct {
	db *sql.DB
}

func NewAddressAreaRepositoryImpl(db *sql.DB) *AddressAreaRepositoryImpl {
	return &AddressAreaRepositoryImpl{db: db}
}

func (repo *AddressAreaRepositoryImpl) FindAllProvince(ctx context.Context) ([]*entity.Province, error) {
	getAll := `SELECT id, name FROM provinces WHERE deleted_at IS NULL ORDER BY id`
	rows, err := repo.db.QueryContext(ctx, getAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Province, 0)
	for rows.Next() {
		var province entity.Province
		if err := rows.Scan(
			&province.Id, &province.Name,
		); err != nil {
			return nil, err
		}
		items = append(items, &province)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *AddressAreaRepositoryImpl) FindCityById(ctx context.Context, cityId int64) (*entity.City, error) {
	getById := `SELECT id, name, province_id FROM cities WHERE id = $1 AND deleted_at IS NULL `
	row := repo.db.QueryRowContext(ctx, getById, cityId)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var city entity.City
	err := row.Scan(
		&city.Id, &city.Name, &city.ProvinceId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}

	return &city, nil
}

func (repo *AddressAreaRepositoryImpl) FindAllCities(ctx context.Context) ([]*entity.City, error) {
	getAll := `SELECT id, name, province_id FROM cities WHERE deleted_at IS NULL ORDER BY province_id, name`
	rows, err := repo.db.QueryContext(ctx, getAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.City, 0)
	for rows.Next() {
		var city entity.City
		if err := rows.Scan(
			&city.Id, &city.Name, &city.ProvinceId,
		); err != nil {
			return nil, err
		}
		items = append(items, &city)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *AddressAreaRepositoryImpl) FindProvinceById(ctx context.Context, provinceId int64) (*entity.Province, error) {
	getById := `SELECT p.id, p.name FROM provinces p WHERE p.id = $1 AND p.deleted_at IS NULL;`
	row := repo.db.QueryRowContext(ctx, getById, provinceId)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var province entity.Province
	err := row.Scan(
		&province.Id, &province.Name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}

	return &province, nil
}
