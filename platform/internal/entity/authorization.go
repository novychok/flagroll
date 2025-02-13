package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Login struct {
	Email    string `validate:"required"`
	Password string `validate:"required"`
}

type RefreshToken struct {
	Token string `validate:"required,jwt"`
}

type VerifyToken struct {
	Token string `validate:"required,jwt"`
}

type Token struct {
	Token                 string
	TokenExpiresAt        time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

type Claims struct {
	jwt.RegisteredClaims
}
