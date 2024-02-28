package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction entity.Transaction) (*entity.Transaction, error)
	FindTransactionById(ctx context.Context, id int64) (*entity.Transaction, error)
	FindAllTransactionsByUserId(ctx context.Context, param *queryparamdto.GetAllParams, userId int64) ([]*entity.Transaction, error)
	CountFindAllTransactionsByUserId(ctx context.Context, userId int64, param *queryparamdto.GetAllParams) (int64, error)
	FindAllTransactions(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.Transaction, error)
	CountFindAllTransactions(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error)
	UpdateTransaction(ctx context.Context, transaction entity.Transaction) (*entity.Transaction, error)
	FindTotalPaymentAndStatusByTransactionId(ctx context.Context, id int64) (*entity.TransactionPaymentAndStatus, *int64, error)
}

type TransactionRepositoryImpl struct {
	db *sql.DB
}

func NewTransactionRepositoryImpl(db *sql.DB) *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{db: db}
}

func (repo *TransactionRepositoryImpl) FindTotalPaymentAndStatusByTransactionId(ctx context.Context, id int64) (*entity.TransactionPaymentAndStatus, *int64, error) {
	const getTotalPayment = `SELECT total_payment, transaction_status_id,user_id FROM transactions WHERE id = $1`
	row := repo.db.QueryRowContext(ctx, getTotalPayment,
		id,
	)
	var payment entity.TransactionPaymentAndStatus
	var userId int64
	err := row.Scan(
		&payment.TotalPayment,
		&payment.TransactionStatusId,
		&userId,
	)
	return &payment, &userId, err
}

func (repo *TransactionRepositoryImpl) UpdateTransaction(ctx context.Context, transaction entity.Transaction) (*entity.Transaction, error) {
	const updateTransaction = `UPDATE transactions 
	SET payment_proof = $1, transaction_status_id = $2, updated_at = now() WHERE id = $3 AND deleted_at IS NULL
	RETURNING id, date, payment_proof, transaction_status_id, payment_method_id, address, user_id, total_payment`

	row := repo.db.QueryRowContext(ctx, updateTransaction,
		transaction.PaymentProof,
		transaction.TransactionStatus.Id,
		transaction.Id,
	)
	var updatedTransaction entity.Transaction
	err := row.Scan(
		&updatedTransaction.Id,
		&updatedTransaction.Date,
		&updatedTransaction.PaymentProof,
		&updatedTransaction.TransactionStatusId,
		&updatedTransaction.PaymentMethodId,
		&updatedTransaction.Address,
		&updatedTransaction.UserId,
		&updatedTransaction.TotalPayment,
	)
	return &updatedTransaction, err
}

func (repo *TransactionRepositoryImpl) FindTransactionById(ctx context.Context, id int64) (*entity.Transaction, error) {
	const viewTransaction = `SELECT transactions.id, date,payment_proof, payment_methods.name,transaction_statuses.id, transaction_statuses.name, address, user_id, total_payment, transactions.created_at
	FROM transactions
	INNER JOIN payment_methods ON transactions.payment_method_id = payment_methods.id
	INNER JOIN transaction_statuses ON transactions.transaction_status_id = transaction_statuses.id
	WHERE transactions.id = $1 AND transactions.deleted_at IS NULL `
	row := repo.db.QueryRowContext(ctx, viewTransaction, id)
	var transaction entity.Transaction
	var method entity.PaymentMethod
	var status entity.TransactionStatus
	err := row.Scan(
		&transaction.Id,
		&transaction.Date,
		&transaction.PaymentProof,
		&method.Name,
		&status.Id,
		&status.Name,
		&transaction.Address,
		&transaction.UserId,
		&transaction.TotalPayment,
		&transaction.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}

	orders, err := repo.findAllOrdersByTransactionId(ctx, transaction.Id)
	if err != nil {
		return nil, err
	}

	transaction.Orders = orders
	transaction.PaymentMethod = &method
	transaction.TransactionStatus = &status
	return &transaction, nil
}

func (repo *TransactionRepositoryImpl) FindAllTransactionsByUserId(ctx context.Context, param *queryparamdto.GetAllParams, userId int64) ([]*entity.Transaction, error) {
	const viewTransactions = `
	SELECT transactions.id, date, payment_proof,payment_methods.name,transaction_statuses.id, transaction_statuses.name,  address,user_id, total_payment, transactions.created_at
	FROM transactions
		INNER JOIN transaction_statuses ON transactions.transaction_status_id = transaction_statuses.id
		INNER JOIN payment_methods ON transactions.payment_method_id = payment_methods.id
	WHERE user_id = $1 AND transactions.deleted_at IS NULL `

	indexPreparedStatement := 1
	query, values := buildQuery(viewTransactions, &entity.Transaction{}, param, true, true, indexPreparedStatement)
	values = util.AppendAtIndex(values, 0, interface{}(userId))

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Transaction, 0)
	for rows.Next() {
		var transaction entity.Transaction
		var method entity.PaymentMethod
		var status entity.TransactionStatus
		if err := rows.Scan(
			&transaction.Id,
			&transaction.Date,
			&transaction.PaymentProof,
			&method.Name,
			&status.Id,
			&status.Name,
			&transaction.Address,
			&transaction.UserId,
			&transaction.TotalPayment,
			&transaction.CreatedAt,
		); err != nil {
			return nil, err
		}
		transaction.PaymentMethod = &method
		transaction.TransactionStatus = &status
		items = append(items, &transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *TransactionRepositoryImpl) CountFindAllTransactionsByUserId(ctx context.Context, userId int64, param *queryparamdto.GetAllParams) (int64, error) {
	const viewTransactions = `
	SELECT count(transactions.id)
	FROM transactions
		INNER JOIN transaction_statuses ON transactions.transaction_status_id = transaction_statuses.id
		INNER JOIN payment_methods ON transactions.payment_method_id = payment_methods.id
	WHERE user_id = $1 AND transactions.deleted_at IS NULL `

	var totalItems int64
	indexPreparedStatement := 1
	query, values := buildQuery(viewTransactions, &entity.Transaction{}, param, false, false, indexPreparedStatement)
	values = util.AppendAtIndex(values, 0, interface{}(userId))

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

func (repo *TransactionRepositoryImpl) FindAllTransactions(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.Transaction, error) {
	const viewTransactions = `
	SELECT transactions.id, date, payment_proof, payment_methods.name,transaction_statuses.id, transaction_statuses.name,address,user_id, total_payment, transactions.created_at
	FROM transactions
		INNER JOIN transaction_statuses ON transactions.transaction_status_id = transaction_statuses.id
		INNER JOIN payment_methods ON transactions.payment_method_id = payment_methods.id
		WHERE transactions.deleted_at IS NULL `

	query, values := buildQuery(viewTransactions, &entity.Transaction{}, param, true, true)
	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Transaction, 0)
	for rows.Next() {
		var transaction entity.Transaction
		var method entity.PaymentMethod
		var status entity.TransactionStatus
		if err := rows.Scan(
			&transaction.Id,
			&transaction.Date,
			&transaction.PaymentProof,
			&method.Name,
			&status.Id,
			&status.Name,
			&transaction.Address,
			&transaction.UserId,
			&transaction.TotalPayment,
			&transaction.CreatedAt,
		); err != nil {
			return nil, err
		}
		transaction.PaymentMethod = &method
		transaction.TransactionStatus = &status
		items = append(items, &transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *TransactionRepositoryImpl) CountFindAllTransactions(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error) {
	const viewTransactions = `
	SELECT count(transactions.id)
	FROM transactions
		INNER JOIN transaction_statuses ON transactions.transaction_status_id = transaction_statuses.id
		INNER JOIN payment_methods ON transactions.payment_method_id = payment_methods.id
		WHERE transactions.deleted_at IS NULL `

	var totalItems int64
	query, values := buildQuery(viewTransactions, &entity.Transaction{}, param, false, false)
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

func (repo *TransactionRepositoryImpl) Create(ctx context.Context, transaction entity.Transaction) (*entity.Transaction, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	const createTransaction = `
	INSERT INTO transactions(date, payment_proof, transaction_status_id, payment_method_id, address, user_id, total_payment)
	values (now(), $1, $2, $3, $4, $5, $6)
	RETURNING id, date, payment_proof, transaction_status_id, payment_method_id, address, user_id, total_payment
	`
	row := tx.QueryRowContext(ctx, createTransaction,
		transaction.PaymentProof, transaction.TransactionStatusId, transaction.PaymentMethodId, transaction.Address,
		transaction.UserId, transaction.TotalPayment,
	)
	var createdTransaction entity.Transaction
	err = row.Scan(
		&createdTransaction.Id,
		&createdTransaction.Date,
		&createdTransaction.PaymentProof,
		&createdTransaction.TransactionStatusId,
		&createdTransaction.PaymentMethodId,
		&createdTransaction.Address,
		&createdTransaction.UserId,
		&createdTransaction.TotalPayment,
	)
	if err != nil {
		return nil, err
	}

	for _, order := range transaction.Orders {
		const createOrder = `
		INSERT INTO orders(date, pharmacy_id, no_of_items, pharmacy_address, shipping_method_id, shipping_cost, total_payment, transaction_id)
		values (now(), $1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
		row = tx.QueryRowContext(ctx, createOrder,
			order.PharmacyId, order.NoOfItems, order.PharmacyAddress, order.ShippingMethodId, order.ShippingCost,
			order.TotalPayment, createdTransaction.Id,
		)
		var orderId int64
		err = row.Scan(
			&orderId,
		)
		if err != nil {
			return nil, err
		}

		for _, detail := range order.OrderDetails {
			const createDetail = `
			INSERT INTO order_details(order_id, product_id, quantity, name, generic_name, content, description, image, price)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

			_, err = tx.ExecContext(ctx, createDetail,
				orderId, detail.ProductId, detail.Quantity, detail.Name, detail.GenericName, detail.Content,
				detail.Description, detail.Image, detail.Price,
			)
			if err != nil {
				return nil, err
			}

		}

		const addStatus = `INSERT INTO order_status_logs(order_id, order_status_id, is_latest, description)
		values ($1, $2, $3, $4)`
		_, err = tx.ExecContext(ctx, addStatus, orderId, appconstant.WaitingPharmacyOrderStatusId, true, "")
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &createdTransaction, err

}

func (repo *TransactionRepositoryImpl) findAllOrdersByTransactionId(ctx context.Context, transId int64) ([]*entity.Order, error) {
	const getAllOrders = `SELECT orders.id, pharmacies.name, pharmacies.address, shipping_methods.name, shipping_cost, total_payment, transaction_id
	FROM orders
	INNER JOIN pharmacies ON orders.pharmacy_id = pharmacies.id
	INNER JOIN shipping_methods ON orders.shipping_method_id = shipping_methods.id
	WHERE transaction_id = $1`

	rows, err := repo.db.QueryContext(ctx, getAllOrders, transId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Order, 0)
	for rows.Next() {
		var order entity.Order
		var pharmacy entity.Pharmacy
		var shipping entity.ShippingMethod
		if err := rows.Scan(
			&order.Id, &pharmacy.Name, &pharmacy.Address, &shipping.Name,
			&order.ShippingCost, &order.TotalPayment, &order.TransactionId,
		); err != nil {
			return nil, err
		}

		details, err := repo.findAllOrderDetailsByOrderId(ctx, order.Id)
		if err != nil {
			return nil, err
		}

		order.OrderDetails = details
		order.Pharmacy = &pharmacy
		order.ShippingMethod = &shipping
		items = append(items, &order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *TransactionRepositoryImpl) findAllOrderDetailsByOrderId(ctx context.Context, orderId int64) ([]*entity.OrderDetail, error) {
	const getAllOrderDetails = `SELECT order_details.id, quantity, name, generic_name, content, description, image, price FROM order_details
	INNER JOIN orders ON order_details.order_id = orders.id WHERE orders.id = $1`

	rows, err := repo.db.QueryContext(ctx, getAllOrderDetails, orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.OrderDetail, 0)
	for rows.Next() {
		var orderDetail entity.OrderDetail
		if err := rows.Scan(
			&orderDetail.Id, &orderDetail.Quantity, &orderDetail.Name, &orderDetail.GenericName,
			&orderDetail.Content, &orderDetail.Description, &orderDetail.Image, &orderDetail.Price,
		); err != nil {
			return nil, err
		}
		items = append(items, &orderDetail)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
