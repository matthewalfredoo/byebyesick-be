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

type PharmacyUseCase interface {
	Add(ctx context.Context, pharmacy entity.Pharmacy) (*entity.Pharmacy, error)
	GetById(ctx context.Context, id int64) (*entity.Pharmacy, error)
	GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
	Edit(ctx context.Context, id int64, pharmacy entity.Pharmacy) (*entity.Pharmacy, error)
	Remove(ctx context.Context, id int64) error
}

type PharmacyUseCaseImpl struct {
	pharmacyRepository    repository.PharmacyRepository
	addressAreaRepository repository.AddressAreaRepository
}

func NewPharmacyUseCaseImpl(pharmacyRepository repository.PharmacyRepository, addressAreaRepository repository.AddressAreaRepository) *PharmacyUseCaseImpl {
	return &PharmacyUseCaseImpl{pharmacyRepository: pharmacyRepository, addressAreaRepository: addressAreaRepository}
}

func (uc *PharmacyUseCaseImpl) Add(ctx context.Context, pharmacy entity.Pharmacy) (*entity.Pharmacy, error) {
	city, err := uc.addressAreaRepository.FindCityById(ctx, pharmacy.CityId)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.ErrInvalidCityProvinceCombi
		}
		return nil, err
	}

	if city.ProvinceId != pharmacy.ProvinceId {
		return nil, apperror.ErrInvalidCityProvinceCombi
	}

	created, err := uc.pharmacyRepository.Create(ctx, pharmacy)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (uc *PharmacyUseCaseImpl) GetById(ctx context.Context, id int64) (*entity.Pharmacy, error) {
	pharmacy, err := uc.pharmacyRepository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(pharmacy, "Id", id)
		}
		return nil, err
	}

	if pharmacy.PharmacyAdminId != ctx.Value(appconstant.ContextKeyUserId) {
		return nil, apperror.ErrForbiddenViewEntity
	}

	return pharmacy, nil
}

func (uc *PharmacyUseCaseImpl) GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	pharmacies, err := uc.pharmacyRepository.FindAll(ctx, param)
	if err != nil {
		return nil, err
	}

	totalItems, err := uc.pharmacyRepository.CountFindAll(ctx, param)
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
		int64(len(pharmacies)),
		int64(*param.PageId),
		pharmacies,
	)
	return paginatedItems, nil
}

func (uc *PharmacyUseCaseImpl) Edit(ctx context.Context, id int64, pharmacy entity.Pharmacy) (*entity.Pharmacy, error) {
	pharmacydb, err := uc.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if pharmacydb.PharmacyAdminId != ctx.Value(appconstant.ContextKeyUserId) {
		return nil, apperror.ErrForbiddenModifyEntity
	}

	city, err := uc.addressAreaRepository.FindCityById(ctx, pharmacy.CityId)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.ErrInvalidCityProvinceCombi
		}
		return nil, err
	}

	if city.ProvinceId != pharmacy.ProvinceId {
		return nil, apperror.ErrInvalidCityProvinceCombi
	}

	pharmacy.Id = id
	updated, err := uc.pharmacyRepository.Update(ctx, pharmacy)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (uc *PharmacyUseCaseImpl) Remove(ctx context.Context, id int64) error {
	pharmacy, err := uc.GetById(ctx, id)
	if err != nil {
		return err
	}

	if pharmacy.PharmacyAdminId != ctx.Value(appconstant.ContextKeyUserId) {
		return apperror.ErrForbiddenModifyEntity
	}

	if err = uc.pharmacyRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
