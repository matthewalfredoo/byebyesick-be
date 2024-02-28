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

type ProductStockMutationUseCase interface {
	Add(ctx context.Context, stockMutation entity.ProductStockMutation) (*entity.ProductStockMutation, error)
	GetAllByPharmacy(ctx context.Context, pharmacyId int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
}

type ProductStockMutationUseCaseImpl struct {
	productStockMutationRepo repository.ProductStockMutationRepository
	pharmacyProductRepo      repository.PharmacyProductRepository
	pharmacyRepository       repository.PharmacyRepository
}

func NewProductStockMutationUseCaseImpl(productStockMutationRepo repository.ProductStockMutationRepository, pharmacyProductRepo repository.PharmacyProductRepository, pharmacyRepository repository.PharmacyRepository) *ProductStockMutationUseCaseImpl {
	return &ProductStockMutationUseCaseImpl{productStockMutationRepo: productStockMutationRepo, pharmacyProductRepo: pharmacyProductRepo, pharmacyRepository: pharmacyRepository}
}

func (uc *ProductStockMutationUseCaseImpl) findPharmacyProductByIdJoinPharmacy(ctx context.Context, id int64) (*entity.PharmacyProduct, error) {
	pharmacyProduct, err := uc.pharmacyProductRepo.FindByIdJoinPharmacy(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(pharmacyProduct, "Id", id)
		}
		return nil, err
	}
	return pharmacyProduct, nil
}

func (uc *ProductStockMutationUseCaseImpl) Add(ctx context.Context, stockMutation entity.ProductStockMutation) (*entity.ProductStockMutation, error) {
	pharmacyProduct, err := uc.findPharmacyProductByIdJoinPharmacy(ctx, stockMutation.PharmacyProductId)
	if err != nil {
		return nil, err
	}

	if pharmacyProduct.Pharmacy.PharmacyAdminId != ctx.Value(appconstant.ContextKeyUserId) {
		return nil, apperror.ErrForbiddenModifyEntity
	}

	if stockMutation.ProductStockMutationTypeId == appconstant.StockMutationTypeReduction &&
		pharmacyProduct.Stock-stockMutation.Stock < 0 {
		return nil, apperror.ErrInsufficientProductStock
	}

	created, err := uc.productStockMutationRepo.Create(ctx, stockMutation)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (uc *ProductStockMutationUseCaseImpl) GetAllByPharmacy(ctx context.Context, pharmacyId int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	if pharmacyId != 0 {
		pharmacy, err := uc.pharmacyRepository.FindById(ctx, pharmacyId)
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(pharmacy, "Id", pharmacyId)
		}
		if err != nil {
			return nil, err
		}
		if pharmacy.PharmacyAdminId != ctx.Value(appconstant.ContextKeyUserId) {
			return nil, apperror.ErrForbiddenViewEntity
		}
	}

	productStockMutation, err := uc.productStockMutationRepo.FindAllJoin(ctx, param)
	if err != nil {
		return nil, err
	}

	totalItems, err := uc.productStockMutationRepo.CountFindAllJoin(ctx, param)
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
		int64(len(productStockMutation)),
		int64(*param.PageId),
		productStockMutation,
	)
	return paginatedItems, nil
}
