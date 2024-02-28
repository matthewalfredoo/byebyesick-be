package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
)

type DoctorSpecializationRepository interface {
	Create(ctx context.Context, specialization entity.DoctorSpecialization) (*entity.DoctorSpecialization, error)
	FindById(ctx context.Context, id int64) (*entity.DoctorSpecialization, error)
	FindAll(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.DoctorSpecialization, error)
	CountFindAll(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error)
	FindAllWithoutParams(ctx context.Context) ([]*entity.DoctorSpecialization, error)
	Update(ctx context.Context, specialization entity.DoctorSpecialization) (*entity.DoctorSpecialization, error)
	Delete(ctx context.Context, id int64) error
}

type DoctorSpecializationRepositoryImpl struct {
	db *sql.DB
}

func NewDoctorSpecializationRepositoryImpl(db *sql.DB) *DoctorSpecializationRepositoryImpl {
	return &DoctorSpecializationRepositoryImpl{db: db}
}

func (repo *DoctorSpecializationRepositoryImpl) Create(ctx context.Context, specialization entity.DoctorSpecialization) (*entity.DoctorSpecialization, error) {
	const create = `INSERT INTO doctor_specializations(name, image)
	VALUES ($1, $2) RETURNING id, name, image, created_at, updated_at, deleted_at`

	row := repo.db.QueryRowContext(ctx, create, specialization.Name, specialization.Image)
	var created entity.DoctorSpecialization
	err := row.Scan(
		&created.Id, &created.Name, &created.Image, &created.CreatedAt, &created.UpdatedAt, &created.DeletedAt,
	)

	return &created, err
}

func (repo *DoctorSpecializationRepositoryImpl) FindById(ctx context.Context, id int64) (*entity.DoctorSpecialization, error) {
	const findById = `SELECT id, name, image, created_at, updated_at, deleted_at FROM doctor_specializations WHERE id = $1 AND deleted_at IS NULL`

	row := repo.db.QueryRowContext(ctx, findById, id)
	var specialization entity.DoctorSpecialization
	err := row.Scan(
		&specialization.Id, &specialization.Name, &specialization.Image, &specialization.CreatedAt, &specialization.UpdatedAt, &specialization.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}

	return &specialization, err
}

func (repo *DoctorSpecializationRepositoryImpl) FindAll(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.DoctorSpecialization, error) {
	initQuery := `SELECT id, name, image FROM doctor_specializations WHERE deleted_at IS NULL `
	query, values := buildQuery(initQuery, &entity.DoctorSpecialization{}, param, true, true)

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.DoctorSpecialization, 0)
	for rows.Next() {
		var specialization entity.DoctorSpecialization
		if err := rows.Scan(
			&specialization.Id, &specialization.Name, &specialization.Image,
		); err != nil {
			return nil, err
		}
		items = append(items, &specialization)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *DoctorSpecializationRepositoryImpl) CountFindAll(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error) {
	initQuery := `SELECT count(id) FROM doctor_specializations WHERE deleted_at IS NULL `
	query, values := buildQuery(initQuery, &entity.DoctorSpecialization{}, param, false, false)

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

func (repo *DoctorSpecializationRepositoryImpl) FindAllWithoutParams(ctx context.Context) ([]*entity.DoctorSpecialization, error) {
	const getDoctorSpecs = `
	SELECT id, name, created_at, updated_at, deleted_at FROM doctor_specializations
	ORDER BY id
	`

	rows, err := repo.db.QueryContext(ctx, getDoctorSpecs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.DoctorSpecialization, 0)
	for rows.Next() {
		var ds entity.DoctorSpecialization
		if err := rows.Scan(
			&ds.Id, &ds.Name, &ds.CreatedAt, &ds.UpdatedAt, &ds.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &ds)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *DoctorSpecializationRepositoryImpl) Update(ctx context.Context, specialization entity.DoctorSpecialization) (*entity.DoctorSpecialization, error) {
	const update = `UPDATE doctor_specializations
	SET name = $1, image = $2 WHERE id = $3 AND deleted_at IS NULL RETURNING id, name, image, created_at, updated_at, deleted_at`

	row := repo.db.QueryRowContext(ctx, update, specialization.Name, specialization.Image, specialization.Id)
	var updated entity.DoctorSpecialization
	err := row.Scan(
		&updated.Id, &updated.Name, &updated.Image, &updated.CreatedAt, &updated.UpdatedAt, &updated.DeletedAt,
	)

	return &updated, err
}

func (repo *DoctorSpecializationRepositoryImpl) Delete(ctx context.Context, id int64) error {
	const deleteQ = `UPDATE doctor_specializations SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL `
	_, err := repo.db.ExecContext(ctx, deleteQ, id)
	return err
}
