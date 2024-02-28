package usecase

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
	"halodeksik-be/app/util"
	"math"
)

type ShippingMethodUseCase interface {
	GetAll(ctx context.Context, addressId int64, checkoutItems []entity.CheckoutItem) (*entity.PaginatedItems, error)
}

type ShippingMethodUseCaseImpl struct {
	shippingMethodRepository  repository.ShippingMethodRepository
	addressRepository         repository.UserAddressRepository
	addressAreaRepository     repository.AddressAreaRepository
	pharmacyProductRepository repository.PharmacyProductRepository
	ongkirUtil                util.OngkirUtil
}

func NewShippingMethodUseCaseImpl(shippingMethodRepository repository.ShippingMethodRepository, addressRepository repository.UserAddressRepository, addressAreaRepository repository.AddressAreaRepository, pharmacyProductRepository repository.PharmacyProductRepository, ongkirUtil util.OngkirUtil) *ShippingMethodUseCaseImpl {
	return &ShippingMethodUseCaseImpl{shippingMethodRepository: shippingMethodRepository, addressRepository: addressRepository, addressAreaRepository: addressAreaRepository, pharmacyProductRepository: pharmacyProductRepository, ongkirUtil: ongkirUtil}
}

func (uc *ShippingMethodUseCaseImpl) GetAll(ctx context.Context, addressId int64, checkoutItems []entity.CheckoutItem) (*entity.PaginatedItems, error) {
	address, err := uc.addressRepository.FindById(ctx, addressId)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(address, "Id", addressId)
		}
		return nil, err
	}

	if address.ProfileId != ctx.Value(appconstant.ContextKeyUserId) {
		return nil, apperror.ErrForbiddenViewEntity
	}

	if len(checkoutItems) < 1 {
		return nil, apperror.ErrGetShipmentMethodNoItems
	}

	var totalWeight float64
	pharmacy := new(entity.Pharmacy)
	for _, item := range checkoutItems {
		pharmacyProduct, err := uc.pharmacyProductRepository.FindByIdJoinPharmacyAndProduct(ctx, item.PharmacyProductId)
		if err != nil {
			if errors.Is(err, apperror.ErrRecordNotFound) {
				return nil, apperror.NewNotFound(pharmacyProduct, "Id", item.PharmacyProductId)
			}
			return nil, err
		}
		if pharmacy.Id != pharmacyProduct.PharmacyId && pharmacy.Id != 0 {
			return nil, apperror.ErrGetShipmentMethodDifferentPharmacy
		}
		pharmacy = pharmacyProduct.Pharmacy
		totalWeight += pharmacyProduct.Product.Weight * float64(item.Quantity)
	}

	cost, err := uc.ongkirUtil.GetCost(pharmacy.CityId, address.CityId, int32(math.Ceil(totalWeight)))
	if err != nil {
		return nil, err
	}

	shippingMethods, err := uc.shippingMethodRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	distance, err := uc.shippingMethodRepository.CalculateDistance(ctx, pharmacy.Latitude, pharmacy.Longitude, address.Latitude, address.Longitude)
	if err != nil {
		return nil, err
	}

	for _, method := range shippingMethods {
		switch method.Id {
		case appconstant.ShippingMethodOfficialInstant:
			method.Cost = decimal.NewFromFloat(math.Ceil(distance) * appconstant.ShippingCostPerKMOfficialInstant)
		case appconstant.ShippingMethodOfficialSameDay:
			method.Cost = decimal.NewFromFloat(math.Ceil(distance) * appconstant.ShippingCostPerKMOfficialSameDay)
		default:
			method.Cost = cost
		}
	}

	paginatedItems := entity.NewPaginationInfo(
		int64(len(shippingMethods)),
		1,
		int64(len(shippingMethods)),
		1,
		shippingMethods,
	)
	return paginatedItems, nil
}
