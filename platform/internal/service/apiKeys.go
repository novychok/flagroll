package service

import (
	"context"

	"github.com/novychok/flagroll/platform/internal/entity"
)

type APIKeys interface {
	Create(ctx context.Context, ownerID entity.UserID, apiKeyCreate *entity.APIKeyCreate) (*entity.APIKeyResponse, error)
	// GetByTokenID(ctx context.Context, apiKeyRaw string) (*entity.APIKey, error)
	Get(ctx context.Context, id entity.APIKeyID) ([]*entity.APIKey, error)
	Delete(ctx context.Context, id entity.APIKeyID) error
}
