package entity

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserId int64  `json:"user_id"`
	Email  string `json:"email"`
	RoleId int64  `json:"user_role_id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	jwt.RegisteredClaims
}
