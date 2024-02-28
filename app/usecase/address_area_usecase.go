package usecase

import (
	"context"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
	"halodeksik-be/app/util"
)

type AddressAreaUseCase interface {
	GetAllProvinces(ctx context.Context) ([]*entity.Province, error)
	GetAllCities(ctx context.Context) ([]*entity.City, error)
	ValidateCityWithLatLong(ctx context.Context, cityId int64, provinceId int64, lat string, long string) error
}

type AddressAreaUseCaseImpl struct {
	repo    repository.AddressAreaRepository
	locUtil util.LocationUtil
}

func NewAddressAreaUseCaseImpl(repo repository.AddressAreaRepository, locationUtil util.LocationUtil) *AddressAreaUseCaseImpl {
	return &AddressAreaUseCaseImpl{repo: repo, locUtil: locationUtil}
}

func (uc *AddressAreaUseCaseImpl) GetAllProvinces(ctx context.Context) ([]*entity.Province, error) {
	provinces, err := uc.repo.FindAllProvince(ctx)
	if err != nil {
		return nil, err
	}
	return provinces, nil
}

func (uc *AddressAreaUseCaseImpl) GetAllCities(ctx context.Context) ([]*entity.City, error) {
	cities, err := uc.repo.FindAllCities(ctx)
	if err != nil {
		return nil, err
	}
	return cities, nil
}

func (uc *AddressAreaUseCaseImpl) ValidateCityWithLatLong(ctx context.Context, cityId int64, provinceId int64, lat string, long string) error {
	city, err := uc.repo.FindCityById(ctx, cityId)
	if err != nil {
		return err
	}

	province, err := uc.repo.FindProvinceById(ctx, provinceId)
	if err != nil {
		return err
	}

	if city.ProvinceId != province.Id {
		return apperror.ErrInvalidCityProvinceCombi
	}

	err = uc.locUtil.ValidateLatLong(city.Name, province.Name, lat, long)
	if err != nil {
		return err
	}

	return nil

}
