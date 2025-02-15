package repository

import (
	"context"

	"github.com/novychok/flagroll/platform/internal/entity"
)

type FeatureFlag interface {
	GetByUserAndName(ctx context.Context, userID entity.UserID, name string) (*entity.FeatureFlag, error)
	GetAll(ctx context.Context) ([]*entity.FeatureFlag, error)
	Create(ctx context.Context, featureFlag *entity.FeatureFlagCreate) (*entity.FeatureFlag, error)
	GetByID(ctx context.Context, id entity.FeatureFlagID) (*entity.FeatureFlag, error)
	Delete(ctx context.Context, id entity.FeatureFlagID) error
	Update(ctx context.Context, id entity.FeatureFlagID, featureFlag *entity.FeatureFlagUpdate) (*entity.FeatureFlag, error)
	UpdateToggle(ctx context.Context, id entity.FeatureFlagID, active bool) (*entity.FeatureFlag, error)
}
