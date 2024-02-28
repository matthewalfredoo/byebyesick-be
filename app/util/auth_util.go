package util

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"halodeksik-be/app/appconfig"
	"halodeksik-be/app/apperror"
)

type AuthUtil interface {
	ComparePassword(hashedPwd, plainPwd string) bool
	SignToken(token *jwt.Token) (string, error)
	HashAndSalt(pwd string) (string, error)
	GenerateSecureToken() (string, error)
}

func NewAuthUtil() AuthUtil {
	return &AuthUtilImpl{}
}

type AuthUtilImpl struct{}

func (u *AuthUtilImpl) ComparePassword(hashedPwd, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	password := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, password)
	if err != nil {
		return false
	}
	return true
}

func (u *AuthUtilImpl) GenerateSecureToken() (string, error) {
	token, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return token.String(), nil
}

func (u *AuthUtilImpl) SignToken(token *jwt.Token) (string, error) {
	tokenString, err := token.SignedString([]byte(appconfig.Config.JwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (u *AuthUtilImpl) HashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return "", apperror.ErrPasswordTooLong
	}
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
