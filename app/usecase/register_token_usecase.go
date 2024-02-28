package usecase

import (
	"context"
	"errors"
	"fmt"
	"halodeksik-be/app/appconfig"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/applogger"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
	"halodeksik-be/app/util"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type RegisterTokenUseCase interface {
	SendRegisterToken(ctx context.Context, email string) (string, error)
	VerifyRegisterToken(ctx context.Context, token string) (*entity.VerificationToken, error)
}

type RegisterTokenUseCaseImpl struct {
	registerTokenRepository repository.RegisterTokenRepository
	userRepository          repository.UserRepository
	authUtil                util.AuthUtil
	mailUtil                util.EmailUtil
	registerTokenExpired    int
	loginTokenExpired       int
	frontEndUrl             string
}

func NewRegisterTokenUseCase(uRepo repository.UserRepository, tRegisterRepo repository.RegisterTokenRepository, aUtil util.AuthUtil, eUtil util.EmailUtil) RegisterTokenUseCase {

	expiryRegister, err := strconv.Atoi(appconfig.Config.RegisterTokenExpired)
	if err != nil {
		return nil
	}

	return &RegisterTokenUseCaseImpl{
		userRepository:          uRepo,
		registerTokenRepository: tRegisterRepo,
		authUtil:                aUtil,
		mailUtil:                eUtil,
		registerTokenExpired:    expiryRegister,
		frontEndUrl:             appconfig.Config.FrontendUrl,
	}
}

func (uc *RegisterTokenUseCaseImpl) VerifyRegisterToken(ctx context.Context, token string) (*entity.VerificationToken, error) {
	existedToken, err := uc.registerTokenRepository.FindRegisterTokenByToken(ctx, token)
	if existedToken == nil {
		return nil, apperror.NewNotFound(existedToken, "Token", token)
	}
	if err != nil {
		return nil, err
	}

	if existedToken.IsValid == false {
		return nil, apperror.ErrRegisterTokenInvalid
	}

	if existedToken.ExpiredAt.Before(time.Now()) {
		return nil, apperror.ErrRegisterTokenExpired
	}
	return existedToken, nil
}

func (uc *RegisterTokenUseCaseImpl) SendRegisterToken(ctx context.Context, email string) (string, error) {
	var userVerify entity.VerificationToken

	existedUser, err := uc.userRepository.FindByEmail(ctx, email)
	if existedUser != nil {
		return "", apperror.NewAlreadyExist(existedUser, "Email", email)
	}
	if err != nil && !errors.Is(err, apperror.ErrRecordNotFound) {
		return "", err
	}

	uid, err := uc.authUtil.GenerateSecureToken()
	if err != nil {
		return "", err
	}

	tokenFound, err := uc.registerTokenRepository.FindRegisterTokenByToken(ctx, uid)
	if tokenFound != nil {
		return "", &apperror.AlreadyExist{
			Resource:        tokenFound,
			FieldInResource: "Token",
			Value:           tokenFound.Token,
		}
	}

	activeToken, err := uc.registerTokenRepository.FindRegisterTokenByEmail(ctx, email)
	if activeToken != nil {
		_, err2 := uc.registerTokenRepository.DeactivateRegisterToken(ctx, *activeToken)
		if err2 != nil {
			return "", err2
		}
	}
	if err != nil && !errors.Is(err, apperror.ErrRecordNotFound) {
		return "", err
	}

	userVerify.Token = uid
	userVerify.Email = email
	userVerify.ExpiredAt = time.Now().Add(time.Duration(uc.registerTokenExpired) * time.Minute)
	userVerify.IsValid = true

	_, err = uc.registerTokenRepository.CreateRegisterToken(ctx, userVerify)
	if err != nil {
		return "", err
	}

	to := []string{email}
	subject := "Email Verification"
	message, err := uc.composeEmail(userVerify)
	if err != nil {
		return "", err
	}

	go func() {
		err := uc.mailUtil.SendEmail(to, []string{}, subject, message)
		if err != nil {
			applogger.Log.Error(err.Error())
		}
	}()

	return uid, nil
}

func (uc *RegisterTokenUseCaseImpl) composeEmail(token entity.VerificationToken) (string, error) {
	htmlFilePath := "app/asset/auth/register_token.html"

	message := fmt.Sprintf("%s/verify-register?token=%s", uc.frontEndUrl, token.Token)

	content, err := ioutil.ReadFile(htmlFilePath)
	if err != nil {
		return "", err
	}

	formattedExpiredAt := token.ExpiredAt.Format(appconstant.TimeHourFormatQueryParam)

	htmlString := string(content)
	htmlString = strings.Replace(htmlString, "{{link}}", message, 1)
	htmlString = strings.Replace(htmlString, "{{tokenExpired}}", formattedExpiredAt, 1)
	htmlString = strings.Replace(htmlString, "{{recipient}}", token.Email, 1)

	return htmlString, nil
}
