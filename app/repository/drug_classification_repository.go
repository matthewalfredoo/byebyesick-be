package repository

import (
	"context"
	"database/sql"
	"halodeksik-be/app/entity"
)

type DrugClassificationRepository interface {
	FindAllWithoutParams(ctx context.Context) ([]*entity.DrugClassification, error)
}

type DrugClassificationRepositoryImpl struct {
	db *sql.DB
}

func NewDrugClassificationRepositoryImpl(db *sql.DB) *DrugClassificationRepositoryImpl {
	return &DrugClassificationRepositoryImpl{db: db}
}

func (repo *DrugClassificationRepositoryImpl) FindAllWithoutParams(ctx context.Context) ([]*entity.DrugClassification, error) {
	const findAll = `
		SELECT id, name, created_at, updated_at, deleted_at FROM drug_classifications
		`

	rows, err := repo.db.QueryContext(ctx, findAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.DrugClassification, 0)
	for rows.Next() {
		var dc entity.DrugClassification
		if err := rows.Scan(
			&dc.Id, &dc.Name, &dc.CreatedAt, &dc.UpdatedAt, &dc.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &dc)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
