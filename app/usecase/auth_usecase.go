package usecase

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"halodeksik-be/app/appcloud"
	"halodeksik-be/app/appconfig"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
	"halodeksik-be/app/util"
	"strconv"
	"time"
)

type AuthUsecase interface {
	Register(ctx context.Context, user entity.User, token string, name string) (*entity.User, error)
	Login(ctx context.Context, req requestdto.LoginRequest) (*entity.User, *responsedto.GenericProfileResponse, error)
	ChangePassword(ctx context.Context, newPassword string, token string) (*entity.User, error)
}

type AuthUseCaseImpl struct {
	userRepository          repository.UserRepository
	profileRepository       repository.ProfileRepository
	forgotTokenRepository   repository.ForgotTokenRepository
	registerTokenRepository repository.RegisterTokenRepository
	authUtil                util.AuthUtil
	uploader                appcloud.FileUploader
	cloudUrl                string
	cloudFolder             string
	loginExpired            int
	forgotTokenUseCase      ForgotTokenUseCase
	registerTokenUseCase    RegisterTokenUseCase
}

type AuthRepos struct {
	UserRepo      repository.UserRepository
	TForgotRepo   repository.ForgotTokenRepository
	TRegisterRepo repository.RegisterTokenRepository
	ProfileRepo   repository.ProfileRepository
}

type AuthUseCases struct {
	TForgotUseCase   ForgotTokenUseCase
	TRegisterUseCase RegisterTokenUseCase
}

func NewAuthUsecase(authRepos AuthRepos, aUtil util.AuthUtil, uploader appcloud.FileUploader, cases AuthUseCases) AuthUsecase {

	expiryLogin, err := strconv.Atoi(appconfig.Config.LoginTokenExpired)
	if err != nil {
		return nil
	}

	return &AuthUseCaseImpl{
		userRepository:          authRepos.UserRepo,
		forgotTokenRepository:   authRepos.TForgotRepo,
		registerTokenRepository: authRepos.TRegisterRepo,
		profileRepository:       authRepos.ProfileRepo,
		authUtil:                aUtil,
		uploader:                uploader,
		cloudUrl:                appconfig.Config.GcloudStorageCdn,
		cloudFolder:             appconfig.Config.GcloudStorageFolderCertificates,
		loginExpired:            expiryLogin,
		forgotTokenUseCase:      cases.TForgotUseCase,
		registerTokenUseCase:    cases.TRegisterUseCase,
	}
}

func (uc *AuthUseCaseImpl) ChangePassword(ctx context.Context, newPassword string, token string) (*entity.User, error) {
	dbToken, err := uc.forgotTokenUseCase.VerifyForgetToken(ctx, token)
	if err != nil {
		return nil, err
	}

	registeredUser, err := uc.userRepository.FindById(ctx, dbToken.UserId)
	if err != nil {
		return nil, err
	}

	newHashedPw, err := uc.authUtil.HashAndSalt(newPassword)
	if err != nil {
		return nil, err
	}

	changedUser, err := uc.userRepository.ChangePassword(ctx, *registeredUser, newHashedPw)
	if err != nil {
		return nil, err
	}

	_, err = uc.forgotTokenRepository.DeactivateForgotToken(ctx, *dbToken)
	if err != nil {
		return nil, err
	}

	return changedUser, nil
}

func (uc *AuthUseCaseImpl) Register(ctx context.Context, user entity.User, token string, name string) (*entity.User, error) {
	verifiedToken, err := uc.verifyToken(ctx, token, user.Email)
	if err != nil {
		return nil, err
	}

	if user.UserRoleId == appconstant.UserRoleIdAdmin || user.UserRoleId == appconstant.UserRoleIdPharmacyAdmin {
		return nil, apperror.ErrInvalidRegisterRole
	}

	hashedPw, err := uc.authUtil.HashAndSalt(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPw
	user.IsVerified = true
	doctorProfile := entity.DoctorProfile{}
	userProfile := entity.UserProfile{}

	fileHeaderAny := ctx.Value(appconstant.FormCertificate)
	if user.UserRoleId == appconstant.UserRoleIdDoctor && fileHeaderAny != nil {
		url, err2 := uc.uploader.UploadFromFileHeader(ctx, fileHeaderAny, uc.cloudFolder)
		if err2 != nil {
			return nil, err2
		}
		doctorProfile.DoctorCertificate = url
	}

	userProfile.Name = name
	doctorProfile.Name = name
	createdUser, err := uc.createUser(ctx, user, &doctorProfile, &userProfile)
	if err != nil {
		return nil, err
	}

	_, err = uc.registerTokenRepository.DeactivateRegisterToken(ctx, *verifiedToken)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (uc *AuthUseCaseImpl) verifyToken(ctx context.Context, token string, email string) (*entity.VerificationToken, error) {
	verifiedToken, err := uc.registerTokenUseCase.VerifyRegisterToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if verifiedToken.Email != email {
		return nil, apperror.ErrRegisterTokenInvalid
	}
	return verifiedToken, nil
}

func (uc *AuthUseCaseImpl) createUser(ctx context.Context, user entity.User, doctorProfile *entity.DoctorProfile, userProfile *entity.UserProfile) (*entity.User, error) {
	createdUser := &entity.User{}
	var err error
	if user.UserRoleId == appconstant.UserRoleIdDoctor {
		doctorProfile.DoctorSpecializationId = appconstant.DoctorSpecializationGeneral
		createdUser, err = uc.userRepository.CreateAndDoctorProfile(ctx, user, *doctorProfile)
		if err != nil {
			return nil, err
		}
	} else {
		createdUser, err = uc.userRepository.CreateAndUserProfile(ctx, user, *userProfile)
		if err != nil {
			return nil, err
		}
	}
	return createdUser, nil
}

func (uc *AuthUseCaseImpl) Login(ctx context.Context, req requestdto.LoginRequest) (*entity.User, *responsedto.GenericProfileResponse, error) {
	user, err := uc.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, err
	}
	var name string
	var image string

	if user.UserRoleId == appconstant.UserRoleIdDoctor {
		userProfile, err := uc.profileRepository.FindDoctorProfileByUserId(ctx, user.Id)
		if err != nil {
			return nil, nil, err
		}
		name = userProfile.DoctorProfile.Name
		image = userProfile.DoctorProfile.ProfilePhoto
	} else if user.UserRoleId == appconstant.UserRoleIdUser {
		userProfile, err := uc.profileRepository.FindUserProfileByUserId(ctx, user.Id)
		if err != nil {
			return nil, nil, err
		}
		name = userProfile.UserProfile.Name
		image = userProfile.UserProfile.ProfilePhoto
	}

	if !uc.authUtil.ComparePassword(user.Password, req.Password) {
		return nil, nil, apperror.ErrWrongCredentials
	}

	expirationTime := time.Now().Add(time.Duration(uc.loginExpired) * time.Minute)
	claims := &entity.Claims{
		UserId: user.Id,
		Email:  user.Email,
		RoleId: user.UserRoleId,
		Name:   name,
		Image:  image,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ByeByeSick Healthcare",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := uc.authUtil.SignToken(token)
	if err != nil {
		return nil, nil, err
	}

	return user, &responsedto.GenericProfileResponse{
		Image: image,
		Name:  name,
		Token: tokenString,
	}, nil
}
