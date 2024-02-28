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

type ProductStockMutationRequestUseCase interface {
	Add(ctx context.Context, mutationRequest entity.ProductStockMutationRequest) (*entity.ProductStockMutationRequest, error)
	GetAllIncoming(ctx context.Context, pharmacyOriginId int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
	GetAllOutgoing(ctx context.Context, pharmacyDestId int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
	EditStatus(ctx context.Context, id int64, mutationRequest entity.ProductStockMutationRequest) (*entity.ProductStockMutationRequest, error)
}

type ProductStockMutationRequestUseCaseImpl struct {
	productStockMutationRequestRepo repository.ProductStockMutationRequestRepository
	pharmacyProductRepo             repository.PharmacyProductRepository
	pharmacyRepo                    repository.PharmacyRepository
}

func NewProductStockMutationRequestUseCaseImpl(productStockMutationRequestRepo repository.ProductStockMutationRequestRepository, pharmacyProductRepo repository.PharmacyProductRepository, pharmacyRepo repository.PharmacyRepository) *ProductStockMutationRequestUseCaseImpl {
	return &ProductStockMutationRequestUseCaseImpl{productStockMutationRequestRepo: productStockMutationRequestRepo, pharmacyProductRepo: pharmacyProductRepo, pharmacyRepo: pharmacyRepo}
}

func (uc *ProductStockMutationRequestUseCaseImpl) findPharmacyProductById(ctx context.Context, id int64) (*entity.PharmacyProduct, error) {
	pharmacyProduct, err := uc.pharmacyProductRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(pharmacyProduct, "Id", id)
		}
		return nil, err
	}
	return pharmacyProduct, nil
}

func (uc *ProductStockMutationRequestUseCaseImpl) Add(ctx context.Context, mutationRequest entity.ProductStockMutationRequest) (*entity.ProductStockMutationRequest, error) {
	destPharmacyProduct, err := uc.findPharmacyProductById(ctx, mutationRequest.PharmacyProductDestId)
	if err != nil {
		return nil, err
	}

	originPharmacyProduct, err := uc.findPharmacyProductById(ctx, mutationRequest.PharmacyProductOriginId)
	if err != nil {
		return nil, err
	}

	destPharmacy, err := uc.pharmacyRepo.FindById(ctx, destPharmacyProduct.PharmacyId)
	if err != nil {
		return nil, err
	}

	if destPharmacy.PharmacyAdminId != ctx.Value(appconstant.ContextKeyUserId) {
		return nil, apperror.ErrForbiddenModifyEntity
	}

	if destPharmacyProduct.ProductId != originPharmacyProduct.ProductId {
		return nil, apperror.ErrRequestStockMutationDifferentProduct
	}

	if destPharmacyProduct.Id == originPharmacyProduct.Id {
		return nil, apperror.ErrRequestStockMutationFromOwnPharmacy
	}

	if originPharmacyProduct.Stock < mutationRequest.Stock {
		return nil, apperror.ErrInsufficientProductStock
	}

	mutationRequest.ProductStockMutationRequestStatusId = appconstant.StockMutationRequestStatusPending
	created, err := uc.productStockMutationRequestRepo.Create(ctx, mutationRequest)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (uc *ProductStockMutationRequestUseCaseImpl) GetAllIncoming(ctx context.Context, pharmacyOriginId int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	pharmacy, err := uc.pharmacyRepo.FindById(ctx, pharmacyOriginId)
	if errors.Is(err, apperror.ErrRecordNotFound) {
		return nil, apperror.NewNotFound(pharmacy, "Id", pharmacyOriginId)
	}
	if err != nil {
		return nil, err
	}
	if pharmacy.PharmacyAdminId != ctx.Value(appconstant.ContextKeyUserId) {
		return nil, apperror.ErrForbiddenViewEntity
	}

	mutationRequest, err := uc.productStockMutationRequestRepo.FindAllJoin(ctx, param)
	if err != nil {
		return nil, err
	}

	totalItems, err := uc.productStockMutationRequestRepo.CountFindAllJoin(ctx, param)
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
		int64(len(mutationRequest)),
		int64(*param.PageId),
		mutationRequest,
	)
	return paginatedItems, nil
}

func (uc *ProductStockMutationRequestUseCaseImpl) GetAllOutgoing(ctx context.Context, pharmacyDestId int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	pharmacy, err := uc.pharmacyRepo.FindById(ctx, pharmacyDestId)
	if errors.Is(err, apperror.ErrRecordNotFound) {
		return nil, apperror.NewNotFound(pharmacy, "Id", pharmacyDestId)
	}
	if err != nil {
		return nil, err
	}
	if pharmacy.PharmacyAdminId != ctx.Value(appconstant.ContextKeyUserId) {
		return nil, apperror.ErrForbiddenViewEntity
	}

	mutationRequest, err := uc.productStockMutationRequestRepo.FindAllJoin(ctx, param)
	if err != nil {
		return nil, err
	}

	totalItems, err := uc.productStockMutationRequestRepo.CountFindAllJoin(ctx, param)
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
		int64(len(mutationRequest)),
		int64(*param.PageId),
		mutationRequest,
	)
	return paginatedItems, nil
}

func (uc *ProductStockMutationRequestUseCaseImpl) EditStatus(ctx context.Context, id int64, mutationRequest entity.ProductStockMutationRequest) (*entity.ProductStockMutationRequest, error) {
	mutationRequestdb, err := uc.productStockMutationRequestRepo.FindByIdJoinPharmacyOrigin(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(mutationRequestdb, "Id", id)
		}
		return nil, err
	}

	if mutationRequestdb.PharmacyProductOrigin.Pharmacy.PharmacyAdminId != ctx.Value(appconstant.ContextKeyUserId) {
		return nil, apperror.ErrForbiddenModifyEntity
	}

	if mutationRequest.ProductStockMutationRequestStatusId == appconstant.StockMutationRequestStatusAccepted &&
		mutationRequestdb.PharmacyProductOrigin.Stock < mutationRequest.Stock {
		return nil, apperror.ErrInsufficientProductStock
	}

	if mutationRequestdb.ProductStockMutationRequestStatusId != appconstant.StockMutationRequestStatusPending {
		return nil, apperror.ErrAlreadyFinishedRequest
	}

	mutationRequestdb.ProductStockMutationRequestStatusId = mutationRequest.ProductStockMutationRequestStatusId
	updated, err := uc.productStockMutationRequestRepo.Update(ctx, *mutationRequestdb)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
