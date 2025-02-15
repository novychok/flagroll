package entity

import (
	"errors"
	"time"
)

var (
	ErrAPIKeyNotFound      = errors.New("api key not found")
	ErrInvalidTokenHash    = errors.New("invalid token hash")
	ErrInvalidaTokenLength = errors.New("invalid token length")
)

type APIKeyID string

func (a APIKeyID) String() string {
	return string(a)
}

type APIKey struct {
	ID        APIKeyID
	OwnerID   UserID
	TokenID   string
	TokenHash string
	CreatedAt time.Time
	ExpiresAt *time.Time
}

type APIKeyCreate struct {
	Name      string
	ExpiresAt *time.Time `validate:"omitempty"`
}

type APIKeyResponse struct {
	ID        APIKeyID
	ApiKeyRaw string
	CreatedAt time.Time
	ExpiresAt *time.Time
}
