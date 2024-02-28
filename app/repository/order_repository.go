package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
	"strings"
)

type OrderRepository interface {
	FindAllOrdersByPharmacyAdminId(ctx context.Context, param *queryparamdto.GetAllParams, adminId int64) ([]*entity.Order, error)
	FindAllOrders(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.Order, error)
	FindAllOrdersByUserId(ctx context.Context, param *queryparamdto.GetAllParams, userId int64) ([]*entity.Order, error)
	CountFindAllOrdersByPharmacyAdminId(ctx context.Context, adminId int64, param *queryparamdto.GetAllParams) (int64, error)
	CountFindAllOrdersByUserId(ctx context.Context, adminId int64, param *queryparamdto.GetAllParams) (int64, error)
	CountFindAllOrders(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error)
	FindOrderById(ctx context.Context, id int64) (*entity.Order, *entity.OrderIds, error)
	UpdateOrderStatus(ctx context.Context, orderId int64, orderLog entity.OrderStatusLog) (*entity.OrderStatusLog, error)
	FindLatestOrderStatusByOrderId(ctx context.Context, id int64) (*entity.OrderStatusLog, error)
	AcceptOrder(ctx context.Context, orderId int64, orderLog entity.OrderStatusLog) (*entity.OrderStatusLog, error)
	CancelOrder(ctx context.Context, orderId int64, orderLog entity.OrderStatusLog) (*entity.OrderStatusLog, error)
	FindAllOrderStatusLogsByOrderId(ctx context.Context, orderId int64) ([]*entity.OrderStatusLog, error)
}

type OrderRepositoryImpl struct {
	db *sql.DB
}

func NewOrderRepositoryImpl(db *sql.DB) OrderRepository {
	return &OrderRepositoryImpl{db: db}
}

func (repo *OrderRepositoryImpl) FindAllOrders(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.Order, error) {
	const findAllOrder = `SELECT orders.id, pharmacy_id, pharmacies.name, orders.date, no_of_items, orders.total_payment, 
       transaction_id, order_statuses.id,order_statuses.name
	FROM orders
			 INNER JOIN pharmacies ON orders.pharmacy_id = pharmacies.id
			 INNER JOIN transactions ON orders.transaction_id = transactions.id
			 INNER JOIN order_status_logs ON orders.id = order_status_logs.order_id
			 INNER JOIN order_statuses ON order_status_logs.order_status_id = order_statuses.id
	WHERE transactions.transaction_status_id = 4 AND order_status_logs.is_latest IS TRUE AND orders.deleted_at IS NULL `

	query, values := buildQuery(findAllOrder, &entity.Transaction{}, param, true, true)
	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Order, 0)
	for rows.Next() {
		var order entity.Order
		var pharmacy entity.Pharmacy
		var status entity.OrderStatus
		if err := rows.Scan(
			&order.Id,
			&pharmacy.Id,
			&pharmacy.Name,
			&order.Date,
			&order.NoOfItems,
			&order.TotalPayment,
			&order.TransactionId,
			&status.Id,
			&status.Name,
		); err != nil {
			return nil, err
		}
		order.Pharmacy = &pharmacy
		order.LatestStatus = &status
		items = append(items, &order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *OrderRepositoryImpl) FindAllOrdersByPharmacyAdminId(ctx context.Context, param *queryparamdto.GetAllParams, adminId int64) ([]*entity.Order, error) {
	const findAllOrderUserId = `SELECT orders.id, pharmacy_id, pharmacies.name, orders.date, no_of_items, orders.total_payment, 
       transaction_id, order_statuses.id,order_statuses.name
	FROM orders
			 INNER JOIN pharmacies ON orders.pharmacy_id = pharmacies.id
			 INNER JOIN transactions ON orders.transaction_id = transactions.id
			 INNER JOIN order_status_logs ON orders.id = order_status_logs.order_id
			 INNER JOIN order_statuses ON order_status_logs.order_status_id = order_statuses.id
	WHERE pharmacies.pharmacy_admin_id = $1 AND transactions.transaction_status_id = 4 AND order_status_logs.is_latest IS TRUE AND orders.deleted_at IS NULL `

	indexPreparedStatement := 1
	query, values := buildQuery(findAllOrderUserId, &entity.Order{}, param, true, true, indexPreparedStatement)
	values = util.AppendAtIndex(values, 0, interface{}(adminId))

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Order, 0)
	for rows.Next() {
		var order entity.Order
		var pharmacy entity.Pharmacy
		var status entity.OrderStatus
		if err := rows.Scan(
			&order.Id,
			&pharmacy.Id,
			&pharmacy.Name,
			&order.Date,
			&order.NoOfItems,
			&order.TotalPayment,
			&order.TransactionId,
			&status.Id,
			&status.Name,
		); err != nil {
			return nil, err
		}
		order.Pharmacy = &pharmacy
		order.LatestStatus = &status
		items = append(items, &order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *OrderRepositoryImpl) FindAllOrdersByUserId(ctx context.Context, param *queryparamdto.GetAllParams, userId int64) ([]*entity.Order, error) {
	const findAllOrderUserId = `SELECT orders.id, pharmacy_id, pharmacies.name, orders.date, no_of_items, orders.total_payment, 
       transaction_id, order_statuses.id,order_statuses.name
	FROM orders
			 INNER JOIN pharmacies ON orders.pharmacy_id = pharmacies.id
			 INNER JOIN transactions ON orders.transaction_id = transactions.id
			 INNER JOIN order_status_logs ON orders.id = order_status_logs.order_id
			 INNER JOIN order_statuses ON order_status_logs.order_status_id = order_statuses.id
	WHERE transactions.user_id = $1 AND transactions.transaction_status_id = 4 AND order_status_logs.is_latest IS TRUE AND orders.deleted_at IS NULL `

	indexPreparedStatement := 1
	query, values := buildQuery(findAllOrderUserId, &entity.Order{}, param, true, true, indexPreparedStatement)
	values = util.AppendAtIndex(values, 0, interface{}(userId))

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Order, 0)
	for rows.Next() {
		var order entity.Order
		var pharmacy entity.Pharmacy
		var status entity.OrderStatus
		if err := rows.Scan(
			&order.Id,
			&pharmacy.Id,
			&pharmacy.Name,
			&order.Date,
			&order.NoOfItems,
			&order.TotalPayment,
			&order.TransactionId,
			&status.Id,
			&status.Name,
		); err != nil {
			return nil, err
		}
		order.Pharmacy = &pharmacy
		order.LatestStatus = &status
		items = append(items, &order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil

}

func (repo *OrderRepositoryImpl) CountFindAllOrders(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error) {

	const findAllOrder = `SELECT count(orders.id)
	FROM orders
			 INNER JOIN pharmacies ON orders.pharmacy_id = pharmacies.id
			 INNER JOIN transactions ON orders.transaction_id = transactions.id
			 INNER JOIN order_status_logs ON orders.id = order_status_logs.order_id
			 INNER JOIN order_statuses ON order_status_logs.order_status_id = order_statuses.id
	WHERE transactions.transaction_status_id = 4 AND order_status_logs.is_latest IS TRUE AND orders.deleted_at IS NULL `

	query, values := buildQuery(findAllOrder, &entity.Order{}, param, false, false)
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

func (repo *OrderRepositoryImpl) CountFindAllOrdersByPharmacyAdminId(ctx context.Context, adminId int64, param *queryparamdto.GetAllParams) (int64, error) {

	const findAllOrderUserId = `SELECT count(orders.id)
	FROM orders
			 INNER JOIN pharmacies ON orders.pharmacy_id = pharmacies.id
			 INNER JOIN transactions ON orders.transaction_id = transactions.id
			 INNER JOIN order_status_logs ON orders.id = order_status_logs.order_id
			 INNER JOIN order_statuses ON order_status_logs.order_status_id = order_statuses.id
	WHERE pharmacies.pharmacy_admin_id = $1 AND transactions.transaction_status_id = 4 AND order_status_logs.is_latest IS TRUE AND orders.deleted_at IS NULL `
	indexPreparedStatement := 1
	query, values := buildQuery(findAllOrderUserId, &entity.Order{}, param, false, false, indexPreparedStatement)
	values = util.AppendAtIndex(values, 0, interface{}(adminId))

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

func (repo *OrderRepositoryImpl) CountFindAllOrdersByUserId(ctx context.Context, adminId int64, param *queryparamdto.GetAllParams) (int64, error) {

	const findAllOrderUserId = `SELECT count(orders.id)
	FROM orders
			 INNER JOIN pharmacies ON orders.pharmacy_id = pharmacies.id
			 INNER JOIN transactions ON orders.transaction_id = transactions.id
			 INNER JOIN order_status_logs ON orders.id = order_status_logs.order_id
			 INNER JOIN order_statuses ON order_status_logs.order_status_id = order_statuses.id
	WHERE transactions.user_id = $1 AND transactions.transaction_status_id = 4 AND order_status_logs.is_latest IS TRUE AND orders.deleted_at IS NULL `

	indexPreparedStatement := 1
	query, values := buildQuery(findAllOrderUserId, &entity.Order{}, param, false, false, indexPreparedStatement)
	values = util.AppendAtIndex(values, 0, interface{}(adminId))
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

func (repo *OrderRepositoryImpl) FindOrderById(ctx context.Context, id int64) (*entity.Order, *entity.OrderIds, error) {
	const findOrderById = `SELECT DISTINCT orders.id, order_statuses.id, order_statuses.name, orders.date, shipping_method_id, 
                shipping_methods.name, shipping_cost, pharmacy_id, 
                pharmacies.name, transactions.address, orders.total_payment, transactions.user_id, pharmacies.pharmacy_admin_id
	FROM orders
		INNER JOIN shipping_methods ON orders.shipping_method_id = shipping_methods.id
		INNER JOIN transactions ON orders.transaction_id = transactions.id
		INNER JOIN order_status_logs ON orders.id = order_status_logs.order_id
		INNER JOIN order_statuses ON order_status_logs.order_status_id = order_statuses.id
		INNER JOIN pharmacies ON orders.pharmacy_id = pharmacies.id
	WHERE transactions.transaction_status_id = 4 AND orders.id = $1 AND order_status_logs.is_latest IS TRUE`

	row := repo.db.QueryRowContext(ctx, findOrderById, id)
	var order entity.Order
	var status entity.OrderStatus
	var shipping entity.ShippingMethod
	var pharmacy entity.Pharmacy
	var ids entity.OrderIds

	err := row.Scan(
		&order.Id,
		&status.Id,
		&status.Name,
		&order.Date,
		&shipping.Id,
		&shipping.Name,
		&order.ShippingCost,
		&pharmacy.Id,
		&pharmacy.Name,
		&order.UserAddress,
		&order.TotalPayment,
		&ids.UserId,
		&ids.PharmacyAdminId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, apperror.ErrRecordNotFound
		}
		return nil, nil, err
	}

	details, err := repo.findAllOrderDetailsByOrderId(ctx, order.Id)
	if err != nil {
		return nil, nil, err
	}

	order.OrderDetails = details
	order.Pharmacy = &pharmacy
	order.ShippingMethod = &shipping
	order.LatestStatus = &status
	return &order, &ids, err

}

func (repo *OrderRepositoryImpl) findAllOrderDetailsByOrderId(ctx context.Context, orderId int64) ([]*entity.OrderDetail, error) {
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

func (repo *OrderRepositoryImpl) UpdateOrderStatus(ctx context.Context, orderId int64, orderLog entity.OrderStatusLog) (*entity.OrderStatusLog, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	const updateOldStatus = `UPDATE order_status_logs SET is_latest = FALSE WHERE order_id = $1 AND is_latest = true`
	row1 := tx.QueryRowContext(ctx, updateOldStatus, orderId)
	if row1.Err() != nil {
		return nil, err
	}

	const addStatus = `INSERT INTO order_status_logs(order_id, order_status_id, is_latest, description)
	values ($1, $2, $3, $4) RETURNING id, order_id, order_status_id, is_latest, description`

	row2 := tx.QueryRowContext(ctx, addStatus, orderId, orderLog.OrderStatusId, orderLog.IsLatest, orderLog.Description)
	var createdStatus entity.OrderStatusLog
	err = row2.Scan(
		&createdStatus.Id,
		&createdStatus.OrderId,
		&createdStatus.OrderStatusId,
		&createdStatus.IsLatest,
		&createdStatus.Description,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &createdStatus, nil

}

func (repo *OrderRepositoryImpl) FindLatestOrderStatusByOrderId(ctx context.Context, id int64) (*entity.OrderStatusLog, error) {
	const findLatest = `SELECT DISTINCT id, order_id, order_status_id, is_latest, description FROM order_status_logs WHERE order_id = $1 AND is_latest = TRUE`
	row := repo.db.QueryRowContext(ctx, findLatest, id)
	var orderStatus entity.OrderStatusLog
	err := row.Scan(
		&orderStatus.Id,
		&orderStatus.OrderId,
		&orderStatus.OrderStatusId,
		&orderStatus.IsLatest,
		&orderStatus.Description,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}

	return &orderStatus, nil
}

func (repo *OrderRepositoryImpl) findAllPharmacyProductsAndOrderDetailsByOrderId(ctx context.Context, orderId int64) ([]*entity.PharmacyProduct, []*entity.OrderDetail, error) {
	getAllPharmacyProducts := `
SELECT pp.id, pp.pharmacy_id, pp.product_id, pp.stock, pp.price, pp.is_active,
       od.id, od.order_id, od.product_id, od.quantity, od.name, od.generic_name, od.content, od.description, od.image, od.price
FROM orders o
         INNER JOIN order_details od ON o.id = od.order_id
         INNER JOIN pharmacy_products pp ON (od.product_id = pp.product_id AND o.pharmacy_id = pp.pharmacy_id)
WHERE o.id = $1
`

	rows, err := repo.db.QueryContext(ctx, getAllPharmacyProducts, orderId)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	pharmacyProducts := make([]*entity.PharmacyProduct, 0)
	orderDetails := make([]*entity.OrderDetail, 0)
	for rows.Next() {
		var (
			pharmacyProduct entity.PharmacyProduct
			orderDetail     entity.OrderDetail
		)
		if err := rows.Scan(
			&pharmacyProduct.Id, &pharmacyProduct.PharmacyId, &pharmacyProduct.ProductId, &pharmacyProduct.Stock, &pharmacyProduct.Price, &pharmacyProduct.IsActive,
			&orderDetail.Id, &orderDetail.OrderId, &orderDetail.ProductId, &orderDetail.Quantity, &orderDetail.Name, &orderDetail.GenericName,
			&orderDetail.Content, &orderDetail.Description, &orderDetail.Image, &orderDetail.Price,
		); err != nil {
			return nil, nil, err
		}
		pharmacyProducts = append(pharmacyProducts, &pharmacyProduct)
		orderDetails = append(orderDetails, &orderDetail)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return pharmacyProducts, orderDetails, nil
}

func (repo *OrderRepositoryImpl) findNearestPharmacyProductByPharmacyAndProductIdWithSufficientStock(ctx context.Context, pharmacy entity.Pharmacy, productId int64, stockNeeded int32) (*entity.PharmacyProduct, error) {
	getNearestPharmacyProduct := `SELECT pp.id, pp.pharmacy_id, pp.product_id, pp.stock, pp.price
FROM pharmacy_products pp
         INNER JOIN pharmacies p ON pp.pharmacy_id = p.id
WHERE pp.pharmacy_id != $3 AND pp.product_id = $4 AND pp.stock >= $5 AND distance($1, $2, p.latitude, p.longitude) <= 25
ORDER BY distance($1, $2, p.latitude, p.longitude)`
	row := repo.db.QueryRowContext(ctx, getNearestPharmacyProduct, pharmacy.Latitude, pharmacy.Longitude, pharmacy.Id, productId, stockNeeded)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var pharmacyProduct entity.PharmacyProduct
	if err := row.Scan(
		&pharmacyProduct.Id, &pharmacyProduct.PharmacyId, &pharmacyProduct.ProductId, &pharmacyProduct.Stock, &pharmacyProduct.Price,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}

	return &pharmacyProduct, nil
}

func (repo *OrderRepositoryImpl) getPharmacyIdAndLocationFromOrderId(ctx context.Context, orderId int64) (*entity.Pharmacy, error) {
	query := `SELECT p.id, p.latitude, p.longitude
			  FROM orders o 
    			INNER JOIN pharmacies p ON p.id = o.pharmacy_id
			  WHERE o.id = $1
`
	row1 := repo.db.QueryRowContext(ctx, query, orderId)
	if row1.Err() != nil {
		return nil, row1.Err()
	}

	var pharmacy entity.Pharmacy
	if err := row1.Scan(
		&pharmacy.Id, &pharmacy.Latitude, &pharmacy.Longitude,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}
	return &pharmacy, nil
}

func (repo *OrderRepositoryImpl) AcceptOrder(ctx context.Context, orderId int64, orderLog entity.OrderStatusLog) (*entity.OrderStatusLog, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	pharmacyProducts, orderDetails, err := repo.findAllPharmacyProductsAndOrderDetailsByOrderId(ctx, orderId)
	if err != nil {
		return nil, err
	}

	needStockTransfer := make(map[*entity.OrderDetail]*entity.PharmacyProduct)
	doableOrders := make(map[*entity.OrderDetail]*entity.PharmacyProduct)
	for i := 0; i < len(pharmacyProducts); i++ {
		if pharmacyProducts[i].Stock < orderDetails[i].Quantity {
			needStockTransfer[orderDetails[i]] = pharmacyProducts[i]
			continue
		}
		doableOrders[orderDetails[i]] = pharmacyProducts[i]
	}

	pharmacyDest, err := repo.getPharmacyIdAndLocationFromOrderId(ctx, orderId)
	if err != nil {
		return nil, err
	}

	stockMutationRequestQuery := `INSERT INTO product_stock_mutation_requests(pharmacy_product_origin_id, pharmacy_product_dest_id, stock, product_stock_mutation_request_status_id, order_detail_id) VALUES `
	stockMutationQuery := `INSERT INTO product_stock_mutations(pharmacy_product_id, product_stock_mutation_type_id, stock) VALUES `
	updateStock := `UPDATE pharmacy_products SET stock = stock + $1 WHERE id = $2`

	for orderDetail, pharmacyProductDest := range needStockTransfer {
		requiredStock := orderDetail.Quantity - pharmacyProductDest.Stock
		pharmacyProductOrigin, err := repo.findNearestPharmacyProductByPharmacyAndProductIdWithSufficientStock(ctx,
			*pharmacyDest, orderDetail.ProductId, requiredStock,
		)
		if err != nil {
			if errors.Is(err, apperror.ErrRecordNotFound) {
				return nil, apperror.ErrNoPharmacyToStockTransfer
			}
			return nil, err
		}
		stockMutationRequestQuery += fmt.Sprintf(
			"(%d, %d, %d, %d, %d),",
			pharmacyProductOrigin.Id, pharmacyProductDest.Id, requiredStock,
			appconstant.StockMutationRequestStatusAccepted, orderDetail.Id,
		)
		stockMutationQuery += fmt.Sprintf(
			"(%d, %d, %d), (%d, %d, %d), (%d, %d, %d),",
			pharmacyProductOrigin.Id, appconstant.StockMutationTypeReduction, requiredStock,
			pharmacyProductDest.Id, appconstant.StockMutationTypeAddition, requiredStock,
			pharmacyProductDest.Id, appconstant.StockMutationTypeReduction, orderDetail.Quantity,
		)

		if _, err := tx.ExecContext(ctx, updateStock,
			0-requiredStock,
			pharmacyProductOrigin.Id,
		); err != nil {
			return nil, err
		}
		if _, err := tx.ExecContext(ctx, updateStock,
			requiredStock,
			pharmacyProductDest.Id,
		); err != nil {
			return nil, err
		}
		if _, err := tx.ExecContext(ctx, updateStock,
			0-orderDetail.Quantity,
			pharmacyProductDest.Id,
		); err != nil {
			return nil, err
		}
	}

	if len(needStockTransfer) > 0 {
		stockMutationRequestQuery = strings.TrimSuffix(stockMutationRequestQuery, ",")
		if _, err := tx.ExecContext(ctx, stockMutationRequestQuery); err != nil {
			return nil, err
		}
	}

	for orderDetail, pharmacyProduct := range doableOrders {
		stockMutationQuery += fmt.Sprintf(
			"(%d, %d, %d),", pharmacyProduct.Id, appconstant.StockMutationTypeReduction, orderDetail.Quantity,
		)
		reducedStock := 0 - orderDetail.Quantity
		if _, err := tx.ExecContext(ctx, updateStock,
			reducedStock,
			pharmacyProduct.Id,
		); err != nil {
			return nil, err
		}
	}

	stockMutationQuery = strings.TrimSuffix(stockMutationQuery, ",")
	if _, err := tx.ExecContext(ctx, stockMutationQuery); err != nil {
		return nil, err
	}

	const updateOldStatus = `UPDATE order_status_logs SET is_latest = FALSE WHERE order_id = $1 AND is_latest = true`
	_, err = tx.ExecContext(ctx, updateOldStatus, orderId)
	if err != nil {
		return nil, err
	}

	const addStatus = `INSERT INTO order_status_logs(order_id, order_status_id, is_latest, description)
	values ($1, $2, $3, $4) RETURNING id, order_id, order_status_id, is_latest, description`

	row2 := tx.QueryRowContext(ctx, addStatus, orderId, orderLog.OrderStatusId, orderLog.IsLatest, orderLog.Description)
	var createdStatus entity.OrderStatusLog
	err = row2.Scan(
		&createdStatus.Id,
		&createdStatus.OrderId,
		&createdStatus.OrderStatusId,
		&createdStatus.IsLatest,
		&createdStatus.Description,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &createdStatus, nil
}

func (repo *OrderRepositoryImpl) findProductStockMutationRequestByOrderDetailId(ctx context.Context, orderDetailId int64) (*entity.ProductStockMutationRequest, error) {
	getProductStockMutationRequestByOrderDetailId := `
SELECT id, pharmacy_product_origin_id, pharmacy_product_dest_id, stock, product_stock_mutation_request_status_id
FROM product_stock_mutation_requests
WHERE order_detail_id = $1`

	row := repo.db.QueryRowContext(ctx, getProductStockMutationRequestByOrderDetailId, orderDetailId)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var mutationRequest entity.ProductStockMutationRequest
	if err := row.Scan(
		&mutationRequest.Id,
		&mutationRequest.PharmacyProductOriginId, &mutationRequest.PharmacyProductDestId, &mutationRequest.Stock,
		&mutationRequest.ProductStockMutationRequestStatusId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}
	return &mutationRequest, nil
}

func (repo *OrderRepositoryImpl) CancelOrder(ctx context.Context, orderId int64, orderLog entity.OrderStatusLog) (*entity.OrderStatusLog, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	pharmacyProducts, orderDetails, err := repo.findAllPharmacyProductsAndOrderDetailsByOrderId(ctx, orderId)
	if err != nil {
		return nil, err
	}

	needReturnStockTransfer := make(map[*entity.OrderDetail]*entity.ProductStockMutationRequest)
	doableRestock := make(map[*entity.OrderDetail]*entity.PharmacyProduct)
	for i := 0; i < len(orderDetails); i++ {
		mutationRequest, err := repo.findProductStockMutationRequestByOrderDetailId(ctx, orderDetails[i].Id)
		if err != nil {
			if errors.Is(err, apperror.ErrRecordNotFound) {
				doableRestock[orderDetails[i]] = pharmacyProducts[i]
				continue
			}
			return nil, err
		}
		needReturnStockTransfer[orderDetails[i]] = mutationRequest
	}

	stockMutationQuery := `INSERT INTO product_stock_mutations(pharmacy_product_id, product_stock_mutation_type_id, stock) VALUES `
	updateStock := `UPDATE pharmacy_products SET stock = stock + $1 WHERE id = $2`

	for orderDetail, mutationRequest := range needReturnStockTransfer {
		returnedStock := mutationRequest.Stock
		stockMutationQuery += fmt.Sprintf(
			"(%d, %d, %d), (%d, %d, %d), (%d, %d, %d),",
			mutationRequest.PharmacyProductDestId, appconstant.StockMutationTypeAddition, orderDetail.Quantity,
			mutationRequest.PharmacyProductDestId, appconstant.StockMutationTypeReduction, returnedStock,
			mutationRequest.PharmacyProductOriginId, appconstant.StockMutationTypeAddition, returnedStock,
		)

		if _, err := tx.ExecContext(ctx, updateStock,
			orderDetail.Quantity,
			mutationRequest.PharmacyProductDestId,
		); err != nil {
			return nil, err
		}
		if _, err := tx.ExecContext(ctx, updateStock,
			0-returnedStock,
			mutationRequest.PharmacyProductDestId,
		); err != nil {
			return nil, err
		}
		if _, err := tx.ExecContext(ctx, updateStock,
			returnedStock,
			mutationRequest.PharmacyProductOriginId,
		); err != nil {
			return nil, err
		}
	}

	for orderDetail, pharmacyProduct := range doableRestock {
		stockMutationQuery += fmt.Sprintf(
			"(%d, %d, %d),", pharmacyProduct.Id, appconstant.StockMutationTypeAddition, orderDetail.Quantity,
		)
		if _, err := tx.ExecContext(ctx, updateStock,
			orderDetail.Quantity,
			pharmacyProduct.Id,
		); err != nil {
			return nil, err
		}
	}

	stockMutationQuery = strings.TrimSuffix(stockMutationQuery, ",")
	if _, err := tx.ExecContext(ctx, stockMutationQuery); err != nil {
		return nil, err
	}

	const updateOldStatus = `UPDATE order_status_logs SET is_latest = FALSE WHERE order_id = $1 AND is_latest = true`
	_, err = tx.ExecContext(ctx, updateOldStatus, orderId)
	if err != nil {
		return nil, err
	}

	const addStatus = `INSERT INTO order_status_logs(order_id, order_status_id, is_latest, description)
	values ($1, $2, $3, $4) RETURNING id, order_id, order_status_id, is_latest, description`

	row2 := tx.QueryRowContext(ctx, addStatus, orderId, orderLog.OrderStatusId, orderLog.IsLatest, orderLog.Description)
	var createdStatus entity.OrderStatusLog
	err = row2.Scan(
		&createdStatus.Id,
		&createdStatus.OrderId,
		&createdStatus.OrderStatusId,
		&createdStatus.IsLatest,
		&createdStatus.Description,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &createdStatus, nil
}

func (repo *OrderRepositoryImpl) FindAllOrderStatusLogsByOrderId(ctx context.Context, orderId int64) ([]*entity.OrderStatusLog, error) {
	const getAllLogs = `SELECT order_status_logs.id, order_statuses.name, order_status_logs.created_at as date, order_status_logs.is_latest, order_status_logs.description FROM order_status_logs
	INNER JOIN orders ON order_status_logs.order_id = orders.id
	INNER JOIN order_statuses ON order_status_logs.order_status_id = order_statuses.id
	WHERE orders.id = $1 ORDER BY order_status_id `

	rows, err := repo.db.QueryContext(ctx, getAllLogs, orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.OrderStatusLog, 0)
	for rows.Next() {
		var log entity.OrderStatusLog
		var status entity.OrderStatus

		if err := rows.Scan(
			&log.Id, &status.Name, &log.CreatedAt, &log.IsLatest, &log.Description,
		); err != nil {
			return nil, err
		}
		log.OrderStatus = &status
		items = append(items, &log)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}
