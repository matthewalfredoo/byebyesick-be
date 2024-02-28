package repository

import (
	"database/sql"
	"fmt"
	"halodeksik-be/app/appconstant"
	"strings"
)

type CronRepository interface {
	ValidateTransactions() error
	ValidateOrders() error
	ValidateOrdersConfirmed() error
}

type CronRepoImpl struct {
	db *sql.DB
}

func NewCronRepoImpl(db *sql.DB) *CronRepoImpl {
	return &CronRepoImpl{db: db}
}

func (repo CronRepoImpl) ValidateOrders() error {

	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	const updateOldStatus = `UPDATE order_status_logs
   	SET is_latest = FALSE FROM orders
    WHERE (orders.date::date + INTERVAL '4 day') <= now()
 	AND is_latest = true AND order_status_logs.order_status_id = 1
	AND orders.id = order_status_logs.order_id RETURNING orders.id `

	rows, err := tx.Query(updateOldStatus)
	if err != nil {
		return err
	}
	var orderIds []int64
	for rows.Next() {
		var orderId int64

		if err := rows.Scan(
			&orderId,
		); err != nil {
			return err
		}

		orderIds = append(orderIds, orderId)
	}

	err = repo.bulkInsertStatus(tx, orderIds, false)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo CronRepoImpl) ValidateOrdersConfirmed() error {

	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	const updateOldStatus = `UPDATE order_status_logs
   	SET is_latest = FALSE FROM orders
    WHERE (orders.date::date + INTERVAL '8 day') <= now()
 	AND is_latest = true AND order_status_logs.order_status_id = 3
	AND orders.id = order_status_logs.order_id RETURNING orders.id `

	rows, err := repo.db.Query(updateOldStatus)
	var orderIds []int64
	if err != nil {
		return err
	}
	for rows.Next() {
		var orderId int64

		if err := rows.Scan(
			&orderId,
		); err != nil {
			return err
		}

		orderIds = append(orderIds, orderId)
	}

	err = repo.bulkInsertStatus(tx, orderIds, true)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo CronRepoImpl) ValidateTransactions() error {
	const setExpiredStatus = `UPDATE transactions SET transaction_status_id = 5
	WHERE (date::date + INTERVAL '4 day') <= now() AND transaction_status_id != 4`

	_, err := repo.db.Exec(setExpiredStatus)
	if err != nil {
		return err
	}

	return nil
}

func (repo CronRepoImpl) bulkInsertStatus(tx *sql.Tx, orderIds []int64, isConfirmed bool) error {
	colSize := 4
	valueStrings := make([]string, 0, len(orderIds))
	valueArgs := make([]interface{}, 0, len(orderIds)*colSize)
	i := 0
	for _, id := range orderIds {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*colSize+1, i*colSize+2, i*colSize+3, i*colSize+4))
		valueArgs = append(valueArgs, id)
		if isConfirmed {
			valueArgs = append(valueArgs, appconstant.ConfirmedUserOrderStatusId)
		} else {
			valueArgs = append(valueArgs, appconstant.CanceledByPharmacyOrderStatusId)
		}
		valueArgs = append(valueArgs, true)

		if isConfirmed {
			valueArgs = append(valueArgs, "status changed automatically")
		} else {
			valueArgs = append(valueArgs, "expired")
		}
		i++
	}
	stmt := fmt.Sprintf("INSERT INTO order_status_logs(order_id, order_status_id, is_latest, description) VALUES %s",
		strings.Join(valueStrings, ","))
	_, err := tx.Exec(stmt, valueArgs...)
	return err
}
