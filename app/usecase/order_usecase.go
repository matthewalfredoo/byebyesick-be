package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type OrderUseCase interface {
	GetAllOrdersByPharmacyAdminId(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
	GetAllOrdersByUserId(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
	GetAllOrders(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
	GetOrderById(ctx context.Context, id int64) (*entity.Order, error)
	ConfirmOrder(ctx context.Context, id int64) (*entity.OrderStatusLog, error)
	RejectOrder(ctx context.Context, id int64) (*entity.OrderStatusLog, error)
	ShipOrder(ctx context.Context, id int64) (*entity.OrderStatusLog, error)
	ReceiveOrder(ctx context.Context, id int64) (*entity.OrderStatusLog, error)
	CancelOrder(ctx context.Context, id int64) (*entity.OrderStatusLog, error)
	GetAllOrderLogsByOrderId(ctx context.Context, id int64) ([]*entity.OrderStatusLog, error)
}

type OrderUseCaseImpl struct {
	repo repository.OrderRepository
}

func NewOrderUseCaseImpl(repo repository.OrderRepository) *OrderUseCaseImpl {
	return &OrderUseCaseImpl{
		repo: repo,
	}
}

func (uc *OrderUseCaseImpl) GetOrderById(ctx context.Context, id int64) (*entity.Order, error) {
	userId := ctx.Value(appconstant.ContextKeyUserId).(int64)
	roleId := ctx.Value(appconstant.ContextKeyRoleId).(int64)

	order, ids, err := uc.repo.FindOrderById(ctx, id)
	if errors.Is(err, apperror.ErrRecordNotFound) {
		return nil, apperror.NewNotFound(order, "Id", id)
	}
	if err != nil {
		return nil, err
	}

	if roleId == appconstant.UserRoleIdPharmacyAdmin {
		if ids.PharmacyAdminId != userId {
			return nil, apperror.ErrForbiddenViewEntity
		}
	} else if roleId == appconstant.UserRoleIdUser {
		if ids.UserId != userId {
			return nil, apperror.ErrForbiddenViewEntity
		}
	}

	return order, nil
}

func (uc *OrderUseCaseImpl) GetAllOrdersByPharmacyAdminId(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	userId := ctx.Value(appconstant.ContextKeyUserId).(int64)

	orders, err := uc.repo.FindAllOrdersByPharmacyAdminId(ctx, param, userId)
	if err != nil {
		return nil, err
	}
	totalItems, err := uc.repo.CountFindAllOrdersByPharmacyAdminId(ctx, userId, param)
	if err != nil {
		return nil, err
	}
	totalPages := totalItems / int64(*param.PageSize)
	if totalItems%int64(*param.PageSize) != 0 || totalPages == 0 {
		totalPages += 1
	}

	paginatedItems := entity.NewPaginationInfo(
		totalItems,
		totalPages,
		int64(len(orders)),
		int64(*param.PageId),
		orders,
	)
	return paginatedItems, nil
}

func (uc *OrderUseCaseImpl) GetAllOrders(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	orders, err := uc.repo.FindAllOrders(ctx, param)
	if err != nil {
		return nil, err
	}
	totalItems, err := uc.repo.CountFindAllOrders(ctx, param)
	if err != nil {
		return nil, err
	}
	totalPages := totalItems / int64(*param.PageSize)
	if totalItems%int64(*param.PageSize) != 0 || totalPages == 0 {
		totalPages += 1
	}

	paginatedItems := entity.NewPaginationInfo(
		totalItems,
		totalPages,
		int64(len(orders)),
		int64(*param.PageId),
		orders,
	)
	return paginatedItems, nil
}

func (uc *OrderUseCaseImpl) GetAllOrdersByUserId(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	userId := ctx.Value(appconstant.ContextKeyUserId).(int64)

	orders, err := uc.repo.FindAllOrdersByUserId(ctx, param, userId)
	if err != nil {
		return nil, err
	}
	totalItems, err := uc.repo.CountFindAllOrdersByUserId(ctx, userId, param)
	if err != nil {
		return nil, err
	}
	totalPages := totalItems / int64(*param.PageSize)
	if totalItems%int64(*param.PageSize) != 0 || totalPages == 0 {
		totalPages += 1
	}

	paginatedItems := entity.NewPaginationInfo(
		totalItems,
		totalPages,
		int64(len(orders)),
		int64(*param.PageId),
		orders,
	)
	return paginatedItems, nil
}

func (uc *OrderUseCaseImpl) ConfirmOrder(ctx context.Context, id int64) (*entity.OrderStatusLog, error) {
	order, err := uc.GetOrderById(ctx, id)
	if err != nil {
		return nil, err
	}

	latestStatus, err := uc.repo.FindLatestOrderStatusByOrderId(ctx, order.Id)
	if err != nil {
		return nil, err
	}

	if latestStatus.OrderStatusId != appconstant.WaitingPharmacyOrderStatusId {
		return nil, apperror.ErrBadConfirmStatus
	}

	newOrder := entity.OrderStatusLog{
		OrderId:       order.Id,
		OrderStatusId: appconstant.ProcessedPharmacyOrderStatusId,
		IsLatest:      true,
	}

	status, err := uc.repo.AcceptOrder(ctx, order.Id, newOrder)
	if err != nil {
		if errors.Is(err, apperror.ErrNoPharmacyToStockTransfer) {
			newOrder.OrderStatusId = appconstant.CanceledByPharmacyOrderStatusId
			newOrder.Description = "There are no pharmacies available with the requested stock"
			status, err := uc.repo.UpdateOrderStatus(ctx, order.Id, newOrder)
			if err != nil {
				return nil, err
			}
			return status, nil
		}
		return nil, err
	}

	return status, nil
}

func (uc *OrderUseCaseImpl) RejectOrder(ctx context.Context, id int64) (*entity.OrderStatusLog, error) {
	order, err := uc.GetOrderById(ctx, id)
	if err != nil {
		return nil, err
	}

	latestStatus, err := uc.repo.FindLatestOrderStatusByOrderId(ctx, order.Id)
	if err != nil {
		return nil, err
	}

	if latestStatus.OrderStatusId != appconstant.WaitingPharmacyOrderStatusId {
		return nil, apperror.ErrBadRejectStatus
	}

	newOrder := entity.OrderStatusLog{
		OrderId:       order.Id,
		OrderStatusId: appconstant.CanceledByPharmacyOrderStatusId,
		IsLatest:      true,
	}

	status, err := uc.repo.UpdateOrderStatus(ctx, order.Id, newOrder)
	if err != nil {
		return nil, err
	}

	return status, nil
}

func (uc *OrderUseCaseImpl) ShipOrder(ctx context.Context, id int64) (*entity.OrderStatusLog, error) {
	order, err := uc.GetOrderById(ctx, id)
	if err != nil {
		return nil, err
	}

	latestStatus, err := uc.repo.FindLatestOrderStatusByOrderId(ctx, order.Id)
	if err != nil {
		return nil, err
	}

	if latestStatus.OrderStatusId != appconstant.ProcessedPharmacyOrderStatusId {
		return nil, apperror.ErrBadShipStatus
	}

	newOrder := entity.OrderStatusLog{
		OrderId:       order.Id,
		OrderStatusId: appconstant.SentPharmacyOrderStatusId,
		IsLatest:      true,
	}

	status, err := uc.repo.UpdateOrderStatus(ctx, order.Id, newOrder)
	if err != nil {
		return nil, err
	}

	return status, nil
}

func (uc *OrderUseCaseImpl) ReceiveOrder(ctx context.Context, id int64) (*entity.OrderStatusLog, error) {
	order, err := uc.GetOrderById(ctx, id)
	if err != nil {
		return nil, err
	}

	latestStatus, err := uc.repo.FindLatestOrderStatusByOrderId(ctx, order.Id)
	if err != nil {
		return nil, err
	}

	if latestStatus.OrderStatusId != appconstant.SentPharmacyOrderStatusId {
		return nil, apperror.ErrBadReceiveStatus
	}

	newOrder := entity.OrderStatusLog{
		OrderId:       order.Id,
		OrderStatusId: appconstant.ConfirmedUserOrderStatusId,
		IsLatest:      true,
	}

	status, err := uc.repo.UpdateOrderStatus(ctx, order.Id, newOrder)
	if err != nil {
		return nil, err
	}

	return status, nil
}

func (uc *OrderUseCaseImpl) CancelOrder(ctx context.Context, id int64) (*entity.OrderStatusLog, error) {
	order, err := uc.GetOrderById(ctx, id)
	if err != nil {
		return nil, err
	}

	latestStatus, err := uc.repo.FindLatestOrderStatusByOrderId(ctx, order.Id)
	if err != nil {
		return nil, err
	}

	if latestStatus.OrderStatusId != appconstant.ProcessedPharmacyOrderStatusId {
		return nil, apperror.ErrBadConfirmStatus
	}

	newOrder := entity.OrderStatusLog{
		OrderId:       order.Id,
		OrderStatusId: appconstant.CanceledByPharmacyOrderStatusId,
		IsLatest:      true,
	}

	status, err := uc.repo.CancelOrder(ctx, order.Id, newOrder)
	if err != nil {
		return nil, err
	}

	return status, nil
}

func (uc *OrderUseCaseImpl) GetAllOrderLogsByOrderId(ctx context.Context, id int64) ([]*entity.OrderStatusLog, error) {
	order, err := uc.GetOrderById(ctx, id)
	if err != nil {
		return nil, err
	}

	logs, err := uc.repo.FindAllOrderStatusLogsByOrderId(ctx, order.Id)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
