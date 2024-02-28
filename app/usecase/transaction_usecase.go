package usecase

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"halodeksik-be/app/appcloud"
	"halodeksik-be/app/appconfig"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type TransactionUseCase interface {
	AddTransaction(ctx context.Context, transaction requestdto.AddTransaction) (*entity.Transaction, error)
	GetTransactionById(ctx context.Context, id int64) (*entity.Transaction, error)
	GetAllTransactions(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
	UploadTransactionPayment(ctx context.Context, id int64) (*entity.Transaction, error)
	CancelTransaction(ctx context.Context, id int64) (*entity.Transaction, error)
	UpdateTransactionStatus(ctx context.Context, id int64, isAccepted bool) (*entity.Transaction, error)
	FindTotalPaymentByTransactionId(ctx context.Context, id int64) (*entity.TransactionPaymentAndStatus, error)
}

type TransactionUseCaseImpl struct {
	transactionRepository     repository.TransactionRepository
	addressRepository         repository.UserAddressRepository
	pharmacyProductRepository repository.PharmacyProductRepository
	uploader                  appcloud.FileUploader
	cloudFolderPaymentProof   string
}

func NewTransactionUseCaseImpl(transRepo repository.TransactionRepository, addressRepo repository.UserAddressRepository, pharmacyProdRepo repository.PharmacyProductRepository, uploader appcloud.FileUploader) *TransactionUseCaseImpl {

	return &TransactionUseCaseImpl{
		transactionRepository:     transRepo,
		addressRepository:         addressRepo,
		pharmacyProductRepository: pharmacyProdRepo,
		uploader:                  uploader,
		cloudFolderPaymentProof:   appconfig.Config.GcloudStoragePaymentProofs,
	}
}

func (uc *TransactionUseCaseImpl) FindTotalPaymentByTransactionId(ctx context.Context, id int64) (*entity.TransactionPaymentAndStatus, error) {
	userId := ctx.Value(appconstant.ContextKeyUserId).(int64)

	totalPayment, idDb, err := uc.transactionRepository.FindTotalPaymentAndStatusByTransactionId(ctx, id)
	if err != nil {
		return nil, err
	}
	if userId != *idDb {
		return nil, apperror.ErrForbiddenViewEntity
	}

	return totalPayment, nil
}

func (uc *TransactionUseCaseImpl) UploadTransactionPayment(ctx context.Context, id int64) (*entity.Transaction, error) {
	userId := ctx.Value(appconstant.ContextKeyUserId).(int64)

	transactionDb, err := uc.transactionRepository.FindTransactionById(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(transactionDb, "Id", id)
		}
		return nil, err
	}

	if transactionDb.UserId != userId {
		return nil, apperror.ErrForbiddenModifyEntity
	}

	if transactionDb.TransactionStatus.Id != appconstant.UnpaidTransactionStatusId &&
		transactionDb.TransactionStatus.Id != appconstant.RejectedTransactionStatusId {
		return nil, apperror.ErrPaymentSent
	}

	proof := ctx.Value(appconstant.FormPaymentProof)
	if proof != nil {
		url, err2 := uc.uploader.UploadFromFileHeader(ctx, proof, uc.cloudFolderPaymentProof)
		if err2 != nil {
			return nil, err2
		}
		transactionDb.PaymentProof = url
	}
	transactionDb.TransactionStatus.Id = appconstant.WaitingTransactionStatusId
	updatedTransaction, err := uc.transactionRepository.UpdateTransaction(ctx, *transactionDb)
	if err != nil {
		return nil, err
	}

	return updatedTransaction, nil

}

func (uc *TransactionUseCaseImpl) GetAllTransactions(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	userId := ctx.Value(appconstant.ContextKeyUserId).(int64)
	roleId := ctx.Value(appconstant.ContextKeyRoleId).(int64)

	var addresses []*entity.Transaction
	var totalItems int64
	var err error
	if roleId == appconstant.UserRoleIdAdmin {
		addresses, err = uc.transactionRepository.FindAllTransactions(ctx, param)
		if err != nil {
			return nil, err
		}
		totalItems, err = uc.transactionRepository.CountFindAllTransactions(ctx, param)
		if err != nil {
			return nil, err
		}
	} else {
		addresses, err = uc.transactionRepository.FindAllTransactionsByUserId(ctx, param, userId)
		if err != nil {
			return nil, err
		}
		totalItems, err = uc.transactionRepository.CountFindAllTransactionsByUserId(ctx, userId, param)
		if err != nil {
			return nil, err
		}
	}

	totalPages := totalItems / int64(*param.PageSize)
	if totalItems%int64(*param.PageSize) != 0 || totalPages == 0 {
		totalPages += 1
	}

	paginatedItems := entity.NewPaginationInfo(
		totalItems,
		totalPages,
		int64(len(addresses)),
		int64(*param.PageId),
		addresses,
	)
	return paginatedItems, nil
}

func (uc *TransactionUseCaseImpl) AddTransaction(ctx context.Context, transaction requestdto.AddTransaction) (*entity.Transaction, error) {
	userId := ctx.Value(appconstant.ContextKeyUserId).(int64)

	addressDb, err := uc.addressRepository.FindById(ctx, transaction.AddressId)
	if err != nil {
		return nil, apperror.NewNotFound(addressDb, "Id", transaction.AddressId)
	}
	if addressDb.ProfileId != userId {
		return nil, apperror.ErrForbiddenModifyEntity
	}

	argTransaction := entity.Transaction{
		TransactionStatusId: appconstant.UnpaidTransactionStatusId,
		PaymentMethodId:     appconstant.BankTransferTransactionMethodId,
		Address:             addressDb.Address,
		UserId:              userId,
	}

	var orders []*entity.Order
	var transactionTotalPayment decimal.Decimal
	for _, order := range transaction.Orders {
		firstDetail := order.OrderDetails[0]

		firstPharProd, err := uc.pharmacyProductRepository.FindByIdJoinPharmacy(ctx, firstDetail.PharmacyProductId)
		if err != nil {
			if errors.Is(err, apperror.ErrRecordNotFound) {
				return nil, apperror.NewNotFound(firstPharProd, "Id", firstDetail.PharmacyProductId)
			}
			return nil, err
		}

		decShippingCost, _ := decimal.NewFromString(order.ShippingCost)
		argOrder := entity.Order{
			PharmacyId:       firstPharProd.PharmacyId,
			NoOfItems:        int32(len(order.OrderDetails)),
			PharmacyAddress:  firstPharProd.Pharmacy.Address,
			ShippingMethodId: order.ShippingMethodId,
			ShippingCost:     decShippingCost,
		}

		var orderDetailsPerOrder []*entity.OrderDetail
		var orderTotalPayment decimal.Decimal
		for _, detail := range order.OrderDetails {
			pharProd, err := uc.pharmacyProductRepository.FindByIdJoinPharmacyAndProduct(ctx, detail.PharmacyProductId)
			if err != nil {
				if errors.Is(err, apperror.ErrRecordNotFound) {
					return nil, apperror.NewNotFound(pharProd, "Id", detail.PharmacyProductId)
				}
				return nil, err
			}
			argDetail := entity.OrderDetail{
				ProductId:   pharProd.ProductId,
				Quantity:    detail.Quantity,
				Name:        pharProd.Product.Name,
				GenericName: pharProd.Product.GenericName,
				Content:     pharProd.Product.Content,
				Description: pharProd.Product.Description,
				Image:       pharProd.Product.Image,
				Price:       pharProd.Price,
			}
			paymentPerItem := decimal.NewFromInt32(detail.Quantity).Mul(pharProd.Price)
			orderTotalPayment = orderTotalPayment.Add(paymentPerItem)
			orderDetailsPerOrder = append(orderDetailsPerOrder, &argDetail)
		}
		argOrder.OrderDetails = orderDetailsPerOrder
		argOrder.TotalPayment = orderTotalPayment.Add(decShippingCost)
		transactionTotalPayment = transactionTotalPayment.Add(argOrder.TotalPayment)
		orders = append(orders, &argOrder)
	}
	argTransaction.Orders = orders
	argTransaction.TotalPayment = transactionTotalPayment

	createdTransaction, err := uc.transactionRepository.Create(ctx, argTransaction)
	if err != nil {
		return nil, err
	}
	return createdTransaction, nil
}

func (uc *TransactionUseCaseImpl) GetTransactionById(ctx context.Context, id int64) (*entity.Transaction, error) {
	userId := ctx.Value(appconstant.ContextKeyUserId).(int64)
	userRole := ctx.Value(appconstant.ContextKeyRoleId).(int64)

	transaction, err := uc.transactionRepository.FindTransactionById(ctx, id)
	if errors.Is(err, apperror.ErrRecordNotFound) {
		return nil, apperror.NewNotFound(transaction, "Id", id)
	}
	if err != nil {
		return nil, err
	}

	if transaction.UserId != userId || (userRole != appconstant.UserRoleIdAdmin && userRole != appconstant.UserRoleIdUser) {
		return nil, apperror.ErrForbiddenViewEntity
	}

	return transaction, nil
}

func (uc *TransactionUseCaseImpl) UpdateTransactionStatus(ctx context.Context, id int64, isAccepted bool) (*entity.Transaction, error) {
	transactionDb, err := uc.transactionRepository.FindTransactionById(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(transactionDb, "Id", id)
		}
		return nil, err
	}

	if transactionDb.TransactionStatus.Id == appconstant.PaidTransactionStatusId {
		return nil, apperror.ErrPaymentConfirmed
	}

	if transactionDb.TransactionStatus.Id == appconstant.UnpaidTransactionStatusId ||
		transactionDb.TransactionStatus.Id == appconstant.RejectedTransactionStatusId {
		return nil, apperror.ErrPaymentNotSent
	}

	if isAccepted == true {
		transactionDb.TransactionStatus.Id = appconstant.PaidTransactionStatusId
	} else {
		transactionDb.TransactionStatus.Id = appconstant.RejectedTransactionStatusId
	}

	updatedTransaction, err := uc.transactionRepository.UpdateTransaction(ctx, *transactionDb)
	if err != nil {
		return nil, err
	}

	return updatedTransaction, nil
}

func (uc *TransactionUseCaseImpl) CancelTransaction(ctx context.Context, id int64) (*entity.Transaction, error) {
	userId := ctx.Value(appconstant.ContextKeyUserId).(int64)

	transactionDb, err := uc.transactionRepository.FindTransactionById(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(transactionDb, "Id", id)
		}
		return nil, err
	}

	if transactionDb.UserId != userId {
		return nil, apperror.ErrForbiddenModifyEntity
	}

	if transactionDb.TransactionStatus.Id != appconstant.UnpaidTransactionStatusId {
		return nil, apperror.ErrBadTransactionCancelStatus
	}

	transactionDb.TransactionStatus.Id = appconstant.CanceledTransactionStatusId
	updatedTransaction, err := uc.transactionRepository.UpdateTransaction(ctx, *transactionDb)
	if err != nil {
		return nil, err
	}

	return updatedTransaction, nil
}
