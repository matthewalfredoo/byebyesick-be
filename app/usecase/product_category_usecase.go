package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type ProductCategoryUseCase interface {
	Add(ctx context.Context, category entity.ProductCategory) (*entity.ProductCategory, error)
	GetById(ctx context.Context, id int64) (*entity.ProductCategory, error)
	GetAllProductCategoriesWithoutParams(ctx context.Context) (*entity.PaginatedItems, error)
	GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
	Edit(ctx context.Context, id int64, category entity.ProductCategory) (*entity.ProductCategory, error)
	Remove(ctx context.Context, id int64) error
}

type ProductCategoryUseCaseImpl struct {
	repo repository.ProductCategoryRepository
}

func NewProductCategoryUseCaseImpl(repo repository.ProductCategoryRepository) *ProductCategoryUseCaseImpl {
	return &ProductCategoryUseCaseImpl{repo: repo}
}

func (uc *ProductCategoryUseCaseImpl) Add(ctx context.Context, category entity.ProductCategory) (*entity.ProductCategory, error) {
	created, err := uc.repo.Create(ctx, category)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (uc *ProductCategoryUseCaseImpl) GetById(ctx context.Context, id int64) (*entity.ProductCategory, error) {
	category, err := uc.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(category, "Id", id)
		}
		return nil, err
	}
	return category, nil
}

func (uc *ProductCategoryUseCaseImpl) GetAllProductCategoriesWithoutParams(ctx context.Context) (*entity.PaginatedItems, error) {
	categories, err := uc.repo.FindAllWithoutParams(ctx)
	if err != nil {
		return nil, err
	}

	paginatedItems := new(entity.PaginatedItems)
	paginatedItems.Items = categories
	paginatedItems.TotalItems = int64(len(categories))
	paginatedItems.TotalPages = 1
	paginatedItems.CurrentPageTotalItems = int64(len(categories))
	paginatedItems.CurrentPage = 1

	return paginatedItems, nil
}

func (uc *ProductCategoryUseCaseImpl) GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	categories, err := uc.repo.FindAll(ctx, param)
	if err != nil {
		return nil, err
	}

	totalItems, err := uc.repo.CountFindAll(ctx, param)
	if err != nil {
		return nil, err
	}
	totalPages := totalItems / int64(*param.PageSize)
	if totalItems%int64(*param.PageSize) != 0 || totalPages == 0 {
		totalPages += 1
	}

	paginatedItems := new(entity.PaginatedItems)
	paginatedItems.Items = categories
	paginatedItems.TotalItems = totalItems
	paginatedItems.TotalPages = totalPages
	paginatedItems.CurrentPageTotalItems = int64(len(categories))
	paginatedItems.CurrentPage = int64(*param.PageId)
	return paginatedItems, nil
}

func (uc *ProductCategoryUseCaseImpl) Edit(ctx context.Context, id int64, category entity.ProductCategory) (*entity.ProductCategory, error) {
	if _, err := uc.GetById(ctx, id); err != nil {
		return nil, err
	}
	category.Id = id
	updated, err := uc.repo.Update(ctx, category)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (uc *ProductCategoryUseCaseImpl) Remove(ctx context.Context, id int64) error {
	if _, err := uc.GetById(ctx, id); err != nil {
		return err
	}
	err := uc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
