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

type ForgotTokenUseCase interface {
	VerifyForgetToken(ctx context.Context, token string) (*entity.ForgotPasswordToken, error)
	SendForgotToken(ctx context.Context, email string) (string, error)
}

type ForgotTokenUseCaseImpl struct {
	forgotTokenRepository repository.ForgotTokenRepository
	userRepository        repository.UserRepository
	authUtil              util.AuthUtil
	mailUtil              util.EmailUtil
	forgotTokenExpired    int
	frontEndUrl           string
}

func NewForgotTokenUsecase(uRepo repository.UserRepository, tForgotRepo repository.ForgotTokenRepository, aUtil util.AuthUtil, eUtil util.EmailUtil) ForgotTokenUseCase {

	expiryForgot, err := strconv.Atoi(appconfig.Config.ForgotTokenExpired)
	if err != nil {
		return nil
	}

	return &ForgotTokenUseCaseImpl{
		userRepository:        uRepo,
		forgotTokenRepository: tForgotRepo,
		authUtil:              aUtil,
		mailUtil:              eUtil,
		forgotTokenExpired:    expiryForgot,
		frontEndUrl:           appconfig.Config.FrontendUrl,
	}
}

func (uc *ForgotTokenUseCaseImpl) SendForgotToken(ctx context.Context, email string) (string, error) {
	var token entity.ForgotPasswordToken

	existedUser, err := uc.userRepository.FindByEmail(ctx, email)
	if errors.Is(err, apperror.ErrRecordNotFound) {
		return "", apperror.NewNotFound(existedUser, "Email", email)
	}
	if err != nil {
		return "", err
	}

	uid, err := uc.authUtil.GenerateSecureToken()
	if err != nil {
		return "", err
	}

	tokenFound, err := uc.forgotTokenRepository.FindForgotTokenByToken(ctx, uid)
	if tokenFound != nil {
		return "", &apperror.AlreadyExist{
			Resource:        tokenFound,
			FieldInResource: "Token",
			Value:           tokenFound.Token,
		}
	}

	activeToken, err := uc.forgotTokenRepository.FindForgotTokenByUserId(ctx, existedUser.Id)
	if activeToken != nil {
		_, err2 := uc.forgotTokenRepository.DeactivateForgotToken(ctx, *activeToken)
		if err2 != nil {
			return "", err2
		}
	}
	if err != nil && !errors.Is(err, apperror.ErrRecordNotFound) {
		return "", err
	}

	token.Token = uid
	token.UserId = existedUser.Id
	token.ExpiredAt = time.Now().Add(time.Duration(uc.forgotTokenExpired) * time.Minute)
	token.IsValid = true

	_, err = uc.forgotTokenRepository.CreateForgotToken(ctx, token)
	if err != nil {
		return "", err
	}

	to := []string{email}
	subject := "Password Reset"

	message, err := uc.composeEmail(token)
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

func (uc *ForgotTokenUseCaseImpl) VerifyForgetToken(ctx context.Context, token string) (*entity.ForgotPasswordToken, error) {
	existedToken, err := uc.forgotTokenRepository.FindForgotTokenByToken(ctx, token)
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

func (uc *ForgotTokenUseCaseImpl) composeEmail(token entity.ForgotPasswordToken) (string, error) {
	htmlFilePath := "app/asset/auth/reset_password.html"

	message := fmt.Sprintf("%s/reset-password?token=%s", uc.frontEndUrl, token.Token)

	content, err := ioutil.ReadFile(htmlFilePath)
	if err != nil {
		return "", err
	}

	formattedExpiredAt := token.ExpiredAt.Format(appconstant.TimeHourFormatQueryParam)

	htmlString := string(content)
	htmlString = strings.Replace(htmlString, "{{link}}", message, 1)
	htmlString = strings.Replace(htmlString, "{{tokenExpired}}", formattedExpiredAt, 1)

	return htmlString, nil
}
