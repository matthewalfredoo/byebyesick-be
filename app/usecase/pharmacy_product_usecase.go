package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type PharmacyProductUseCase interface {
	Add(ctx context.Context, pharmacyProduct entity.PharmacyProduct) (*entity.PharmacyProduct, error)
	GetById(ctx context.Context, id int64) (*entity.PharmacyProduct, error)
	GetAllByProductId(ctx context.Context, id int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
	GetAllByPharmacy(ctx context.Context, pharmacyId int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
	Edit(ctx context.Context, id int64, pharmacyProduct entity.PharmacyProduct) (*entity.PharmacyProduct, error)
}

type PharmacyProductUseCaseImpl struct {
	pharmacyProductRepo repository.PharmacyProductRepository
	pharmacyRepo        repository.PharmacyRepository
	productRepo         repository.ProductRepository
}

func NewPharmacyProductUseCaseImpl(
	pharmacyProductRepo repository.PharmacyProductRepository,
	pharmacyRepo repository.PharmacyRepository,
	productRepo repository.ProductRepository,
) *PharmacyProductUseCaseImpl {
	return &PharmacyProductUseCaseImpl{
		pharmacyProductRepo: pharmacyProductRepo,
		pharmacyRepo:        pharmacyRepo,
		productRepo:         productRepo,
	}
}

func (uc *PharmacyProductUseCaseImpl) Add(ctx context.Context, pharmacyProduct entity.PharmacyProduct) (*entity.PharmacyProduct, error) {
	if pharmacy, err := uc.pharmacyRepo.FindById(ctx, pharmacyProduct.PharmacyId); err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(pharmacy, "Id", pharmacyProduct.PharmacyId)
		}
		return nil, err
	}

	if product, err := uc.productRepo.FindById(ctx, pharmacyProduct.ProductId); err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(product, "Id", pharmacyProduct.ProductId)
		}
		return nil, err
	}

	created, err := uc.pharmacyProductRepo.Create(ctx, pharmacyProduct)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (uc *PharmacyProductUseCaseImpl) GetById(ctx context.Context, id int64) (*entity.PharmacyProduct, error) {
	pharmacyProduct, err := uc.pharmacyProductRepo.FindById(ctx, id)
	if errors.Is(err, apperror.ErrRecordNotFound) {
		return nil, apperror.NewNotFound(pharmacyProduct, "Id", id)
	}
	if err != nil {
		return nil, err
	}

	return pharmacyProduct, nil
}

func (uc *PharmacyProductUseCaseImpl) GetAllByProductId(ctx context.Context, id int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	pharmacyProduct, err := uc.pharmacyProductRepo.FindById(ctx, id)
	if errors.Is(err, apperror.ErrRecordNotFound) {
		return nil, apperror.NewNotFound(pharmacyProduct, "Id", id)
	}
	if err != nil {
		return nil, err
	}
	param.WhereClauses = append(param.WhereClauses,
		appdb.NewWhere(pharmacyProduct.GetSqlColumnFromField("ProductId"),
			appdb.EqualTo,
			pharmacyProduct.ProductId,
		),
		appdb.NewWhere(pharmacyProduct.GetSqlColumnFromField("Id"),
			appdb.NotEqualTo,
			pharmacyProduct.Id,
		),
		appdb.NewWhere(pharmacyProduct.GetSqlColumnFromField("Stock"),
			appdb.GreaterThan,
			0,
		),
	)
	pharmacyProducts, err := uc.pharmacyProductRepo.FindAllJoinPharmacy(ctx, param)
	if err != nil {
		return nil, err
	}

	totalItems, err := uc.pharmacyProductRepo.CountFindAllJoinPharmacy(ctx, param)
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
		int64(len(pharmacyProducts)),
		int64(*param.PageId),
		pharmacyProducts,
	)
	return paginatedItems, nil
}

func (uc *PharmacyProductUseCaseImpl) GetAllByPharmacy(ctx context.Context, pharmacyId int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	if pharmacy, err := uc.pharmacyRepo.FindById(ctx, pharmacyId); err != nil {
		return nil, apperror.NewNotFound(pharmacy, "Id", pharmacyId)
	}

	pharmacyProducts, err := uc.pharmacyProductRepo.FindAllJoinProducts(ctx, param)
	if err != nil {
		return nil, err
	}

	totalItems, err := uc.pharmacyProductRepo.CountFindAll(ctx, param)
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
		int64(len(pharmacyProducts)),
		int64(*param.PageId),
		pharmacyProducts,
	)
	return paginatedItems, nil
}

func (uc *PharmacyProductUseCaseImpl) Edit(ctx context.Context, id int64, pharmacyProduct entity.PharmacyProduct) (*entity.PharmacyProduct, error) {
	if _, err := uc.GetById(ctx, id); err != nil {
		return nil, err
	}
	pharmacyProduct.Id = id
	updated, err := uc.pharmacyProductRepo.Update(ctx, pharmacyProduct)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
