package repository

import (
	"context"

	"github.com/novychok/flagroll/platform/internal/entity"
)

type APIKey interface {
	Create(ctx context.Context, apiKey *entity.APIKey) (*entity.APIKey, error)
	GetByTokenID(ctx context.Context, tokenID string) (*entity.APIKey, error)
	Get(ctx context.Context, id entity.APIKeyID) ([]*entity.APIKey, error)
	Delete(ctx context.Context, id entity.APIKeyID) error
}
