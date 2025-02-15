package apikey

import (
	"context"
	"database/sql"
	"errors"

	"github.com/novychok/flagroll/platform/internal/entity"
	"github.com/novychok/flagroll/platform/internal/pkg/postgres"
	"github.com/novychok/flagroll/platform/internal/repository"
	"github.com/samber/lo"

	"github.com/novychok/flagroll/platform/internal/database/dao"
	"github.com/novychok/flagroll/platform/internal/database/pqmodels"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type postgresRepository struct {
	db *sql.DB
}

func (r *postgresRepository) Create(ctx context.Context,
	apiKey *entity.APIKey) (*entity.APIKey, error) {
	apiKeyDB := &pqmodels.APIKey{
		TokenID:   apiKey.TokenID,
		OwnerID:   apiKey.OwnerID.String(),
		TokenHash: apiKey.TokenHash,
		ExpiresAt: lo.FromPtr(apiKey.ExpiresAt),
	}

	err := apiKeyDB.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	dao.APIKeyTo(apiKeyDB, apiKey)

	return apiKey, nil
}

func (r *postgresRepository) GetByTokenID(ctx context.Context,
	tokenID string) (*entity.APIKey, error) {
	apiKeyDB, err := pqmodels.APIKeys(
		pqmodels.APIKeyWhere.TokenID.EQ(tokenID),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrAPIKeyNotFound
		}

		return nil, err
	}

	apiKey := &entity.APIKey{}
	dao.APIKeyTo(apiKeyDB, apiKey)

	return apiKey, nil
}

func (r *postgresRepository) Get(ctx context.Context, id entity.APIKeyID) ([]*entity.APIKey, error) {

	return nil, nil
}

func (r *postgresRepository) Delete(ctx context.Context, id entity.APIKeyID) error {

	return nil
}

func NewPostgres(db postgres.Connection) repository.APIKey {
	return &postgresRepository{
		db: db,
	}
}
