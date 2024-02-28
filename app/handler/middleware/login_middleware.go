package middleware

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"halodeksik-be/app/appconfig"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/handler"
	"strings"
)

func doAuth(ctx *gin.Context) (*entity.Claims, error) {
	c := ctx.GetHeader("Authorization")

	if c == "" {
		return nil, &apperror.AuthError{Err: apperror.ErrLoginNoToken}
	}

	c = strings.ReplaceAll(c, "Bearer ", "")
	claims := &entity.Claims{}

	tkn, err := jwt.ParseWithClaims(c, claims, func(token *jwt.Token) (any, error) {
		return []byte(appconfig.Config.JwtSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, &apperror.AuthError{Err: apperror.ErrLoginTokenInvalidSign}
		}
		return nil, &apperror.AuthError{Err: err}
	}
	if !tkn.Valid {
		return nil, &apperror.AuthError{Err: apperror.ErrLoginTokenNotValid}
	}
	return claims, nil
}

func LoginMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		claim, err := doAuth(ctx)
		if err != nil {
			_ = ctx.Error(handler.WrapError(err))
			ctx.Abort()
			return
		}

		reqCtx1 := ctx.Request.Context()
		reqCtx2 := context.WithValue(reqCtx1, appconstant.ContextKeyUserId, claim.UserId)
		reqCtx3 := context.WithValue(reqCtx2, appconstant.ContextKeyEmail, claim.Email)
		reqCtx4 := context.WithValue(reqCtx3, appconstant.ContextKeyRoleId, claim.RoleId)
		ctx.Request = ctx.Request.WithContext(reqCtx4)
		ctx.Next()

	}
}

func AllowRoles(auths ...int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roleId := ctx.Request.Context().Value("role_id")
		if roleId == nil {
			_ = ctx.Error(handler.WrapError(&apperror.AuthError{Err: apperror.ErrUnauthorized}))
			ctx.Abort()
			return
		}
		isAllowed := false
		for _, auth := range auths {
			if auth == roleId.(int64) {
				isAllowed = true
				break
			}
		}
		if isAllowed == false {
			_ = ctx.Error(handler.WrapError(&apperror.AuthError{Err: apperror.ErrUnauthorized}))
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
