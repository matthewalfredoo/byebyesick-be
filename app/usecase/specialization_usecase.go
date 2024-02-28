package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/appcloud"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
	"halodeksik-be/app/util"
)

type DoctorSpecializationUseCase interface {
	Add(ctx context.Context, specialization entity.DoctorSpecialization) (*entity.DoctorSpecialization, error)
	GetById(ctx context.Context, id int64) (*entity.DoctorSpecialization, error)
	GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
	GetAllSpecsWithoutParams(ctx context.Context) (*entity.PaginatedItems, error)
	Edit(ctx context.Context, id int64, specialization entity.DoctorSpecialization) (*entity.DoctorSpecialization, error)
	Remove(ctx context.Context, id int64) error
}

type DoctorSpecializationUseCaseImpl struct {
	repo     repository.DoctorSpecializationRepository
	uploader appcloud.FileUploader
}

func NewDoctorSpecializationUseCaseImpl(repo repository.DoctorSpecializationRepository, uploader appcloud.FileUploader) *DoctorSpecializationUseCaseImpl {
	return &DoctorSpecializationUseCaseImpl{repo: repo, uploader: uploader}
}

func (uc *DoctorSpecializationUseCaseImpl) Add(ctx context.Context, specialization entity.DoctorSpecialization) (*entity.DoctorSpecialization, error) {
	var (
		err      error
		fileName string
	)
	fileHeader := ctx.Value(appconstant.FormImage)

	if fileHeader != nil {
		fileName, err = uc.uploader.UploadFromFileHeader(ctx, fileHeader, specialization.GetEntityName())
		if err != nil {
			return nil, err
		}
	}

	if !util.IsEmptyString(fileName) {
		specialization.Image = fileName
	}

	created, err := uc.repo.Create(ctx, specialization)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (uc *DoctorSpecializationUseCaseImpl) GetById(ctx context.Context, id int64) (*entity.DoctorSpecialization, error) {
	specialization, err := uc.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(specialization, "Id", id)
		}
		return nil, err
	}
	return specialization, nil
}

func (uc *DoctorSpecializationUseCaseImpl) GetAll(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	specializations, err := uc.repo.FindAll(ctx, param)
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

	paginatedItems := entity.NewPaginationInfo(
		totalItems, totalPages, int64(len(specializations)), int64(*param.PageId), specializations,
	)
	return paginatedItems, nil
}

func (uc *DoctorSpecializationUseCaseImpl) GetAllSpecsWithoutParams(ctx context.Context) (*entity.PaginatedItems, error) {
	doctorSpecs, err := uc.repo.FindAllWithoutParams(ctx)
	if err != nil {
		return nil, err
	}

	paginatedItems := new(entity.PaginatedItems)
	paginatedItems.Items = doctorSpecs
	paginatedItems.TotalItems = int64(len(doctorSpecs))
	paginatedItems.TotalPages = 1
	paginatedItems.CurrentPageTotalItems = int64(len(doctorSpecs))
	paginatedItems.CurrentPage = 1

	return paginatedItems, nil
}

func (uc *DoctorSpecializationUseCaseImpl) Edit(ctx context.Context, id int64, specialization entity.DoctorSpecialization) (*entity.DoctorSpecialization, error) {
	var (
		err      error
		fileName string
	)

	specializationDb, err := uc.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	specialization.Id = id
	specialization.Image = specializationDb.Image

	fileHeader := ctx.Value(appconstant.FormImage)

	if fileHeader != nil {
		fileName, err = uc.uploader.UploadFromFileHeader(ctx, fileHeader, specialization.GetEntityName())
		if err != nil {
			return nil, err
		}
	}

	if !util.IsEmptyString(fileName) {
		specialization.Image = fileName
	}

	updated, err := uc.repo.Update(ctx, specialization)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (uc *DoctorSpecializationUseCaseImpl) Remove(ctx context.Context, id int64) error {
	_, err := uc.GetById(ctx, id)
	if err != nil {
		return err
	}
	err = uc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
