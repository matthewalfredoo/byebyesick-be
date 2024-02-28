package repository

import (
	"context"
	"database/sql"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
)

type ReportRepository interface {
	FindSalesAllPharmacy(ctx context.Context, year int64, param *queryparamdto.GetAllParams) ([]*entity.SellReport, error)
	CountSalesAllPharmacy(ctx context.Context, year int64, param *queryparamdto.GetAllParams) (int64, error)
	FindSalesAllPharmacyMonthly(ctx context.Context, year int64, param *queryparamdto.GetAllParams) ([]*entity.SellReportMonthly, error)
	CountSalesAllPharmacyMonthly(ctx context.Context, year int64, param *queryparamdto.GetAllParams) (int64, error)
}

type ReportRepositoryImpl struct {
	db *sql.DB
}

func (repo ReportRepositoryImpl) FindSalesAllPharmacy(ctx context.Context, year int64, param *queryparamdto.GetAllParams) ([]*entity.SellReport, error) {
	const sellReportAll = `SELECT users.email as pharmacy_admin_email, pharmacies.id, pharmacies.name, SUM(orders.total_payment) total_sells
		FROM orders
			INNER JOIN pharmacies ON orders.pharmacy_id = pharmacies.id
			INNER JOIN users ON pharmacies.pharmacy_admin_id = users.id
			INNER JOIN order_status_logs ON orders.id = order_status_logs.order_id
		WHERE order_status_logs.is_latest = true AND order_status_logs.order_status_id = 4
		AND extract(YEAR FROM orders.date) = $1 `

	indexPreparedStatement := 1
	query, values := buildQuery(sellReportAll, &entity.Pharmacy{}, param, true, true, indexPreparedStatement)
	values = util.AppendAtIndex(values, 0, interface{}(year))

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.SellReport, 0)
	for rows.Next() {
		var sellReport entity.SellReport

		if err := rows.Scan(
			&sellReport.PharmacyAdminEmail, &sellReport.PharmacyId, &sellReport.PharmacyName, &sellReport.TotalSells,
		); err != nil {
			return nil, err
		}
		sellReport.Year = year
		items = append(items, &sellReport)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

func (repo ReportRepositoryImpl) CountSalesAllPharmacy(ctx context.Context, year int64, param *queryparamdto.GetAllParams) (int64, error) {
	const sellReportAll = `SELECT users.email as pharmacy_admin_email, pharmacies.id, pharmacies.name, SUM(orders.total_payment) total_sells
		FROM orders
			INNER JOIN pharmacies ON orders.pharmacy_id = pharmacies.id
			INNER JOIN users ON pharmacies.pharmacy_admin_id = users.id
			INNER JOIN order_status_logs ON orders.id = order_status_logs.order_id
		WHERE order_status_logs.is_latest = true AND order_status_logs.order_status_id = 4
		AND extract(YEAR FROM orders.date) = $1 `

	indexPreparedStatement := 1
	query, values := buildQuery(sellReportAll, &entity.Pharmacy{}, param, false, false, indexPreparedStatement)
	values = util.AppendAtIndex(values, 0, interface{}(year))

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	items := make([]*entity.SellReport, 0)
	for rows.Next() {
		var sellReport entity.SellReport

		if err := rows.Scan(
			&sellReport.PharmacyAdminEmail, &sellReport.PharmacyId, &sellReport.PharmacyName, &sellReport.TotalSells,
		); err != nil {
			return 0, err
		}
		items = append(items, &sellReport)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	return int64(len(items)), err
}

func (repo ReportRepositoryImpl) FindSalesAllPharmacyMonthly(ctx context.Context, year int64, param *queryparamdto.GetAllParams) ([]*entity.SellReportMonthly, error) {
	const sellsReportMonthly = `SELECT extract(MONTH FROM order_details.created_at) as MONTH, SUM(order_details.price*order_details.quantity) as total_sells
	FROM order_details
			 INNER JOIN orders ON order_details.order_id = orders.id
			 INNER JOIN pharmacies ON orders.pharmacy_id = pharmacies.id
			 INNER JOIN order_status_logs ON orders.id = order_status_logs.order_id
			 INNER JOIN products ON order_details.product_id = products.id
	WHERE order_status_logs.is_latest = true AND order_status_logs.order_status_id = 4
	AND extract(YEAR FROM order_details.created_at) = $1 `

	indexPreparedStatement := 1
	query, values := buildQuery(sellsReportMonthly, &entity.OrderDetail{}, param, true, false, indexPreparedStatement)
	values = util.AppendAtIndex(values, 0, interface{}(year))

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.SellReportMonthly, 0)

	for i := 1; i <= appconstant.MonthInAYear; i++ {
		items = append(items, &entity.SellReportMonthly{
			Month:     int32(i),
			TotalSell: 0,
		})
	}

	for rows.Next() {
		var sellReport entity.SellReportMonthly

		if err := rows.Scan(
			&sellReport.Month, &sellReport.TotalSell,
		); err != nil {
			return nil, err
		}
		items[sellReport.Month-1] = &sellReport
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

func (repo ReportRepositoryImpl) CountSalesAllPharmacyMonthly(ctx context.Context, year int64, param *queryparamdto.GetAllParams) (int64, error) {
	const sellsReportMonthly = `SELECT extract(MONTH FROM order_details.created_at) as MONTH, SUM(order_details.price*order_details.quantity) as total_sells
	FROM order_details
			 INNER JOIN orders ON order_details.order_id = orders.id
			 INNER JOIN pharmacies ON orders.pharmacy_id = pharmacies.id
			 INNER JOIN order_status_logs ON orders.id = order_status_logs.order_id
			 INNER JOIN products ON order_details.product_id = products.id
	WHERE order_status_logs.is_latest = true AND order_status_logs.order_status_id = 4
	AND extract(YEAR FROM order_details.created_at) = $1 `

	indexPreparedStatement := 1
	query, values := buildQuery(sellsReportMonthly, &entity.OrderDetail{}, param, false, false, indexPreparedStatement)
	values = util.AppendAtIndex(values, 0, interface{}(year))

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	items := make([]*entity.SellReportMonthly, 0)
	for rows.Next() {
		var sellReport entity.SellReportMonthly

		if err := rows.Scan(
			&sellReport.Month, &sellReport.TotalSell,
		); err != nil {
			return 0, err
		}
		items = append(items, &sellReport)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	return int64(len(items)), err
}

func NewReportRepositoryImpl(db *sql.DB) *ReportRepositoryImpl {
	repo := ReportRepositoryImpl{db: db}
	return &repo
}
