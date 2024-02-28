package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
)

type ProductCategoryRepository interface {
	Create(ctx context.Context, category entity.ProductCategory) (*entity.ProductCategory, error)
	FindById(ctx context.Context, id int64) (*entity.ProductCategory, error)
	FindAllWithoutParams(ctx context.Context) ([]*entity.ProductCategory, error)
	FindAll(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.ProductCategory, error)
	CountFindAll(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error)
	Update(ctx context.Context, category entity.ProductCategory) (*entity.ProductCategory, error)
	Delete(ctx context.Context, id int64) error
}

type ProductCategoryRepositoryImpl struct {
	db *sql.DB
}

func NewProductCategoryRepositoryImpl(db *sql.DB) *ProductCategoryRepositoryImpl {
	return &ProductCategoryRepositoryImpl{db: db}
}

func (repo *ProductCategoryRepositoryImpl) Create(ctx context.Context, category entity.ProductCategory) (*entity.ProductCategory, error) {
	const create = `
		INSERT INTO product_categories(name)
		VALUES ($1)
		RETURNING id, name, created_at, updated_at, deleted_at
		`

	row := repo.db.QueryRowContext(ctx, create, category.Name)
	if row.Err() != nil {
		var errPgConn *pgconn.PgError
		if errors.As(row.Err(), &errPgConn) && errPgConn.Code == apperror.PgconnErrCodeUniqueConstraintViolation {
			return nil, apperror.ErrProductCategoryUniqueConstraint
		}
		return nil, row.Err()
	}

	var created entity.ProductCategory
	err := row.Scan(
		&created.Id, &created.Name, &created.CreatedAt, &created.UpdatedAt, &created.DeletedAt,
	)
	return &created, err
}

func (repo *ProductCategoryRepositoryImpl) FindById(ctx context.Context, id int64) (*entity.ProductCategory, error) {
	const getById = `SELECT id, name, created_at, updated_at, deleted_at FROM product_categories WHERE id = $1 AND deleted_at IS NULL`

	row := repo.db.QueryRowContext(ctx, getById, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var category entity.ProductCategory
	err := row.Scan(
		&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt, &category.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}
	return &category, err
}

func (repo *ProductCategoryRepositoryImpl) FindAllWithoutParams(ctx context.Context) ([]*entity.ProductCategory, error) {
	const findAll = `
		SELECT id, name, created_at, updated_at, deleted_at FROM product_categories WHERE deleted_at IS NULL
		`

	rows, err := repo.db.QueryContext(ctx, findAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*entity.ProductCategory, 0)
	for rows.Next() {
		var category entity.ProductCategory
		if err := rows.Scan(
			&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt, &category.DeletedAt,
		); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (repo *ProductCategoryRepositoryImpl) FindAll(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.ProductCategory, error) {
	initQuery := `SELECT id, name FROM product_categories WHERE deleted_at IS NULL `
	query, values := buildQuery(initQuery, &entity.ProductCategory{}, param, true, true)

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*entity.ProductCategory, 0)
	for rows.Next() {
		var category entity.ProductCategory
		if err := rows.Scan(
			&category.Id, &category.Name,
		); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (repo *ProductCategoryRepositoryImpl) CountFindAll(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error) {
	initQuery := `SELECT count(id) FROM product_categories WHERE deleted_at IS NULL `
	query, values := buildQuery(initQuery, &entity.ProductCategory{}, param, false, false)

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

func (repo *ProductCategoryRepositoryImpl) Update(ctx context.Context, category entity.ProductCategory) (*entity.ProductCategory, error) {
	const update = `UPDATE product_categories SET name = $1, updated_at = now() WHERE id = $2
		RETURNING id, name, created_at, updated_at, deleted_at`

	row := repo.db.QueryRowContext(ctx, update, category.Name, category.Id)
	if row.Err() != nil {
		var errPgConn *pgconn.PgError
		if errors.As(row.Err(), &errPgConn) && errPgConn.Code == apperror.PgconnErrCodeUniqueConstraintViolation {
			return nil, apperror.ErrProductCategoryUniqueConstraint
		}
		return nil, row.Err()
	}

	var updated entity.ProductCategory
	err := row.Scan(
		&updated.Id, &updated.Name, &updated.CreatedAt, &updated.UpdatedAt, &updated.DeletedAt,
	)
	return &updated, err
}

func (repo *ProductCategoryRepositoryImpl) Delete(ctx context.Context, id int64) error {
	const checkProducts = `SELECT count(*) FROM products WHERE product_category_id = $1`
	row := repo.db.QueryRowContext(ctx, checkProducts, id)
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return apperror.ErrProductCategoryStillUsedByProducts
	}

	const deletePC = `UPDATE product_categories SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`
	_, err = repo.db.ExecContext(ctx, deletePC, id)
	return err
}
