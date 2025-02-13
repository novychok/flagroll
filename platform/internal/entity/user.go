package entity

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUserNotFound               = errors.New("user not found")
	ErrAuthorizationFailed        = errors.New("authorization failed")
	ErrMissingAuthorizationHeader = errors.New("missing Authorization header")
	ErrTokenExpired               = errors.New("token expired")
	ErrInvalidToken               = errors.New("invalid token")
	ErrFailedToReadSecret         = errors.New("failed to read secret")
	ErrFailedToGetSecretVersion   = errors.New("failed to get secret version")
	ErrFailedToReadSignature      = errors.New("failed to read signature")
)

type UserID string

func (id UserID) String() string {
	return string(id)
}

type User struct {
	ID           UserID
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserCreate struct {
	Name            string `validate:"required"`
	Email           string `validate:"required,email"`
	Password        string `validate:"required,min=8"`
	PasswordConfirm string `validate:"required,eqfield=Password"`
}

type UserLogin struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type UserToken struct {
	Token            string
	RefreshToken     string
	ExpiresAt        time.Time
	RefreshExpiresAt time.Time
}

type RefreshTokenRequest struct {
	RefreshToken string `validate:"required,jwt"`
}

type UserClaims struct {
	Version string `json:"ver"`
	jwt.RegisteredClaims
}
