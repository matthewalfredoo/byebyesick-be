package usecase

import (
	"context"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type DrugClassificationUseCase interface {
	GetAllDrugsWithoutParams(ctx context.Context) (*entity.PaginatedItems, error)
}

type DrugClassificationUseCaseImpl struct {
	repo repository.DrugClassificationRepository
}

func NewDrugClassificationUseCaseImpl(repo repository.DrugClassificationRepository) *DrugClassificationUseCaseImpl {
	return &DrugClassificationUseCaseImpl{repo: repo}
}

func (uc *DrugClassificationUseCaseImpl) GetAllDrugsWithoutParams(ctx context.Context) (*entity.PaginatedItems, error) {
	drugClassifications, err := uc.repo.FindAllWithoutParams(ctx)
	if err != nil {
		return nil, err
	}

	paginatedItems := new(entity.PaginatedItems)
	paginatedItems.Items = drugClassifications
	paginatedItems.TotalItems = int64(len(drugClassifications))
	paginatedItems.TotalPages = 1
	paginatedItems.CurrentPageTotalItems = int64(len(drugClassifications))
	paginatedItems.CurrentPage = 1

	return paginatedItems, nil
}
