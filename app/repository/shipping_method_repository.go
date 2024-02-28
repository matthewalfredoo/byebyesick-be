package repository

import (
	"context"
	"database/sql"
	"halodeksik-be/app/entity"
)

type ShippingMethodRepository interface {
	FindAll(ctx context.Context) ([]*entity.ShippingMethod, error)
	CalculateDistance(ctx context.Context, originLat, originLong, destLat, destLong string) (float64, error)
}

type ShippingMethodRepositoryImpl struct {
	db *sql.DB
}

func NewShippingMethodRepositoryImpl(db *sql.DB) *ShippingMethodRepositoryImpl {
	return &ShippingMethodRepositoryImpl{db: db}
}

func (repo *ShippingMethodRepositoryImpl) FindAll(ctx context.Context) ([]*entity.ShippingMethod, error) {
	getAll := `SELECT id, name FROM shipping_methods WHERE deleted_at IS NULL`

	rows, err := repo.db.QueryContext(ctx, getAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.ShippingMethod, 0)
	for rows.Next() {
		var shippingMethod entity.ShippingMethod
		if err := rows.Scan(
			&shippingMethod.Id, &shippingMethod.Name,
		); err != nil {
			return nil, err
		}
		items = append(items, &shippingMethod)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *ShippingMethodRepositoryImpl) CalculateDistance(ctx context.Context, originLat, originLong, destLat, destLong string) (float64, error) {
	query := `SELECT distance($1, $2, $3, $4)`
	row := repo.db.QueryRowContext(ctx, query, originLat, originLong, destLat, destLong)
	if row.Err() != nil {
		return 0, row.Err()
	}
	var distance float64
	err := row.Scan(&distance)
	if err != nil {
		return 0, err
	}

	return distance, nil
}
