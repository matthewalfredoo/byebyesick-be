package usecase

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"halodeksik-be/app/appcloud"
	"halodeksik-be/app/appconfig"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type ProfileUseCase interface {
	GetUserProfileByUserId(ctx context.Context, userId int64) (*entity.User, error)
	GetDoctorProfileByUserId(ctx context.Context, userId int64) (*entity.User, error)
	UpdateUserProfile(ctx context.Context, profile entity.UserProfile) (*entity.User, error)
	UpdateDoctorProfile(ctx context.Context, profile entity.DoctorProfile) (*entity.User, error)
	UpdateDoctorIsOnline(ctx context.Context, isOnline bool) (*entity.User, error)
}

type ProfileUseCaseImpl struct {
	repo                   repository.ProfileRepository
	uploader               appcloud.FileUploader
	cloudFolderProfile     string
	cloudFolderCertificate string
}

func NewProfileUseCaseImpl(repo repository.ProfileRepository, uploader appcloud.FileUploader) *ProfileUseCaseImpl {
	return &ProfileUseCaseImpl{repo: repo, cloudFolderProfile: appconfig.Config.GcloudStorageFolderProfiles, cloudFolderCertificate: appconfig.Config.GcloudStorageFolderCertificates, uploader: uploader}
}

func (uc *ProfileUseCaseImpl) UpdateDoctorIsOnline(ctx context.Context, isOnline bool) (*entity.User, error) {
	userId := ctx.(*gin.Context).Request.Context().Value(appconstant.ContextKeyUserId)
	if userId == nil {
		return nil, apperror.ErrUnauthorized
	}

	profile := new(entity.DoctorProfile)

	user, err := uc.repo.FindDoctorProfileByUserId(ctx, userId.(int64))
	if err != nil {
		return nil, err
	}

	profile = user.DoctorProfile
	profile.IsOnline = isOnline

	updatedProfile, err := uc.repo.UpdateDoctorProfileByUserId(ctx, *profile)
	if err != nil {
		return nil, err
	}

	user.DoctorProfile = updatedProfile
	return user, nil
}

func (uc *ProfileUseCaseImpl) UpdateUserProfile(ctx context.Context, profile entity.UserProfile) (*entity.User, error) {
	userId := ctx.(*gin.Context).Request.Context().Value(appconstant.ContextKeyUserId)
	if userId == nil {
		return nil, apperror.ErrUnauthorized
	}
	profile.UserId = userId.(int64)
	user, err := uc.repo.FindUserProfileByUserId(ctx, userId.(int64))
	if err != nil {
		return nil, err
	}
	profile.ProfilePhoto = user.UserProfile.ProfilePhoto
	photo := ctx.Value(appconstant.FormProfilePhoto)
	if photo != nil {
		url, err2 := uc.uploader.UploadFromFileHeader(ctx, photo, uc.cloudFolderProfile)
		if err2 != nil {
			return nil, err2
		}
		profile.ProfilePhoto = url
	}

	updatedProfile, err := uc.repo.UpdateUserProfileByUserId(ctx, profile)
	if err != nil {
		return nil, err
	}

	user.UserProfile = updatedProfile
	return user, nil

}

func (uc *ProfileUseCaseImpl) UpdateDoctorProfile(ctx context.Context, profile entity.DoctorProfile) (*entity.User, error) {
	userId := ctx.(*gin.Context).Request.Context().Value(appconstant.ContextKeyUserId)
	if userId == nil {
		return nil, apperror.ErrUnauthorized
	}
	profile.UserId = userId.(int64)

	user, err := uc.repo.FindDoctorProfileByUserId(ctx, userId.(int64))
	if err != nil {
		return nil, err
	}

	profile.ProfilePhoto = user.DoctorProfile.ProfilePhoto
	profile.IsOnline = user.DoctorProfile.IsOnline
	photo := ctx.Value(appconstant.FormProfilePhoto)
	if photo != nil {
		url, err2 := uc.uploader.UploadFromFileHeader(ctx, photo, uc.cloudFolderProfile)
		if err2 != nil {
			return nil, err2
		}
		profile.ProfilePhoto = url
	}

	profile.DoctorCertificate = user.DoctorProfile.DoctorCertificate
	cert := ctx.Value(appconstant.FormCertificate)
	if cert != nil {
		url, err2 := uc.uploader.UploadFromFileHeader(ctx, cert, uc.cloudFolderCertificate)
		if err2 != nil {
			return nil, err2
		}
		profile.DoctorCertificate = url
	}

	updatedProfile, err := uc.repo.UpdateDoctorProfileByUserId(ctx, profile)
	if err != nil {
		return nil, err
	}

	user.DoctorProfile = updatedProfile
	return user, nil
}

func (uc *ProfileUseCaseImpl) GetUserProfileByUserId(ctx context.Context, userId int64) (*entity.User, error) {
	profile, err := uc.repo.FindUserProfileByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(profile, "Id", userId)
		}
		return nil, err
	}
	return profile, nil
}

func (uc *ProfileUseCaseImpl) GetDoctorProfileByUserId(ctx context.Context, userId int64) (*entity.User, error) {
	profile, err := uc.repo.FindDoctorProfileByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(profile, "Id", userId)
		}
		return nil, err
	}
	return profile, nil
}
