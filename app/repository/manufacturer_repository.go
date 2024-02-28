package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
)

type ManufacturerRepository interface {
	Create(ctx context.Context, manufacturer entity.Manufacturer) (*entity.Manufacturer, error)
	FindById(ctx context.Context, id int64) (*entity.Manufacturer, error)
	FindAll(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.Manufacturer, error)
	CountFindAll(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error)
	FindAllWithoutParams(ctx context.Context) ([]*entity.Manufacturer, error)
	Update(ctx context.Context, manufacturer entity.Manufacturer) (*entity.Manufacturer, error)
	Delete(ctx context.Context, id int64) error
}

type ManufacturerRepositoryImpl struct {
	db *sql.DB
}

func NewManufacturerRepositoryImpl(db *sql.DB) *ManufacturerRepositoryImpl {
	return &ManufacturerRepositoryImpl{db: db}
}

func (repo *ManufacturerRepositoryImpl) Create(ctx context.Context, manufacturer entity.Manufacturer) (*entity.Manufacturer, error) {
	const create = `INSERT INTO manufacturers(name, image)
	VALUES ($1, $2) RETURNING id, name, image, created_at, updated_at, deleted_at`

	row := repo.db.QueryRowContext(ctx, create, manufacturer.Name, manufacturer.Image)
	var created entity.Manufacturer
	err := row.Scan(
		&created.Id, &created.Name, &created.Image, &created.CreatedAt, &created.UpdatedAt, &created.DeletedAt,
	)

	return &created, err
}

func (repo *ManufacturerRepositoryImpl) FindAll(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.Manufacturer, error) {
	initQuery := `SELECT id, name, image FROM manufacturers WHERE deleted_at IS NULL `
	query, values := buildQuery(initQuery, &entity.Manufacturer{}, param, true, true)

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Manufacturer, 0)
	for rows.Next() {
		var manufacturer entity.Manufacturer
		if err := rows.Scan(
			&manufacturer.Id, &manufacturer.Name, &manufacturer.Image,
		); err != nil {
			return nil, err
		}
		items = append(items, &manufacturer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *ManufacturerRepositoryImpl) CountFindAll(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error) {
	initQuery := `SELECT count(id) FROM manufacturers WHERE deleted_at IS NULL `
	query, values := buildQuery(initQuery, &entity.Manufacturer{}, param, false, false)

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

func (repo *ManufacturerRepositoryImpl) FindAllWithoutParams(ctx context.Context) ([]*entity.Manufacturer, error) {
	const getAllWithoutParams = `
		SELECT id, name, created_at, updated_at, deleted_at FROM manufacturers WHERE deleted_at IS NULL
		`

	rows, err := repo.db.QueryContext(ctx, getAllWithoutParams)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Manufacturer, 0)
	for rows.Next() {
		var manufacturer entity.Manufacturer
		if err := rows.Scan(
			&manufacturer.Id, &manufacturer.Name, &manufacturer.CreatedAt, &manufacturer.UpdatedAt, &manufacturer.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &manufacturer)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *ManufacturerRepositoryImpl) FindById(ctx context.Context, id int64) (*entity.Manufacturer, error) {
	const findById = `SELECT id, name, image, created_at, updated_at, deleted_at FROM manufacturers WHERE id = $1 AND deleted_at IS NULL`

	row := repo.db.QueryRowContext(ctx, findById, id)
	var manufacturer entity.Manufacturer
	err := row.Scan(
		&manufacturer.Id, &manufacturer.Name, &manufacturer.Image, &manufacturer.CreatedAt, &manufacturer.UpdatedAt, &manufacturer.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}

	return &manufacturer, err
}

func (repo *ManufacturerRepositoryImpl) Update(ctx context.Context, manufacturer entity.Manufacturer) (*entity.Manufacturer, error) {
	const update = `UPDATE manufacturers
	SET name = $1, image = $2 WHERE id = $3 AND deleted_at IS NULL RETURNING id, name, image, created_at, updated_at, deleted_at`

	row := repo.db.QueryRowContext(ctx, update, manufacturer.Name, manufacturer.Image, manufacturer.Id)
	var updated entity.Manufacturer
	err := row.Scan(
		&updated.Id, &updated.Name, &updated.Image, &updated.CreatedAt, &updated.UpdatedAt, &updated.DeletedAt,
	)

	return &updated, err
}

func (repo *ManufacturerRepositoryImpl) Delete(ctx context.Context, id int64) error {
	const deleteQ = `UPDATE manufacturers SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL `
	_, err := repo.db.ExecContext(ctx, deleteQ, id)
	return err
}
