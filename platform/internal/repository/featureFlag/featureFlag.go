package featureflag

import (
	"context"
	"database/sql"
	"time"

	"github.com/novychok/flagroll/platform/internal/database/dao"
	"github.com/novychok/flagroll/platform/internal/database/pqmodels"
	"github.com/novychok/flagroll/platform/internal/entity"
	"github.com/novychok/flagroll/platform/internal/pkg/postgres"
	"github.com/novychok/flagroll/platform/internal/repository"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type postgresRepository struct {
	db *sql.DB
}

func (r *postgresRepository) GetByUserAndName(ctx context.Context, userID entity.UserID, name string) (*entity.FeatureFlag, error) {
	featureFlagDB, err := pqmodels.FeatureFlags(
		pqmodels.FeatureFlagWhere.OwnerID.EQ(userID.String()),
		pqmodels.FeatureFlagWhere.Name.EQ(name),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}

	ff := &entity.FeatureFlag{}
	dao.FeatureFlagTo(featureFlagDB, ff)

	return ff, nil
}

func (r *postgresRepository) GetAll(ctx context.Context) ([]*entity.FeatureFlag, error) {
	featureFlagsDB, err := pqmodels.FeatureFlags().All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	var featureFlags []*entity.FeatureFlag
	for _, featureFlagDB := range featureFlagsDB {
		ff := &entity.FeatureFlag{}
		dao.FeatureFlagTo(featureFlagDB, ff)
		featureFlags = append(featureFlags, ff)
	}

	return featureFlags, nil
}

func (r *postgresRepository) Create(ctx context.Context,
	featureFlag *entity.FeatureFlagCreate) (*entity.FeatureFlag, error) {
	featureFlagDB := &pqmodels.FeatureFlag{
		OwnerID:     featureFlag.OwnerID.String(),
		Name:        featureFlag.Name,
		Active:      featureFlag.Active,
		Description: null.StringFrom(featureFlag.Description),
	}

	err := featureFlagDB.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	ff := &entity.FeatureFlag{}
	dao.FeatureFlagTo(featureFlagDB, ff)

	return ff, nil
}

func (r *postgresRepository) GetByID(ctx context.Context,
	id entity.FeatureFlagID) (*entity.FeatureFlag, error) {
	featureFlagDB, err := pqmodels.FeatureFlags(
		pqmodels.FeatureFlagWhere.ID.EQ(id.String()),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}

	ff := &entity.FeatureFlag{}
	dao.FeatureFlagTo(featureFlagDB, ff)

	return ff, nil
}

func (r *postgresRepository) Delete(ctx context.Context,
	id entity.FeatureFlagID) error {
	featureFlagDB, err := pqmodels.FeatureFlags(
		pqmodels.FeatureFlagWhere.ID.EQ(id.String()),
	).One(ctx, r.db)
	if err != nil {
		return err
	}

	_, err = featureFlagDB.Delete(ctx, r.db)
	if err != nil {
		return err
	}

	return nil
}

func (r *postgresRepository) Update(ctx context.Context,
	id entity.FeatureFlagID,
	featureFlag *entity.FeatureFlagUpdate) (*entity.FeatureFlag, error) {
	featureFlagDB, err := pqmodels.FeatureFlags(
		pqmodels.FeatureFlagWhere.ID.EQ(id.String()),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}

	featureFlagDB.Name = featureFlag.Name
	featureFlagDB.Active = featureFlag.Active
	featureFlagDB.Description = null.StringFrom(featureFlag.Description)
	featureFlagDB.UpdatedAt = time.Now()

	if _, err := featureFlagDB.Update(ctx, r.db, boil.Whitelist(
		pqmodels.FeatureFlagColumns.Name,
		pqmodels.FeatureFlagColumns.Active,
		pqmodels.FeatureFlagColumns.Description,
		pqmodels.FeatureFlagColumns.UpdatedAt,
	)); err != nil {
		return nil, err
	}

	ff := &entity.FeatureFlag{}
	dao.FeatureFlagTo(featureFlagDB, ff)

	return ff, nil
}

func (r *postgresRepository) UpdateToggle(ctx context.Context,
	id entity.FeatureFlagID, active bool) (*entity.FeatureFlag, error) {
	featureFlagDB, err := pqmodels.FeatureFlags(
		pqmodels.FeatureFlagWhere.ID.EQ(id.String()),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}

	featureFlagDB.Active = active
	featureFlagDB.UpdatedAt = time.Now()

	if _, err := featureFlagDB.Update(ctx, r.db, boil.Whitelist(
		pqmodels.FeatureFlagColumns.Active,
		pqmodels.FeatureFlagColumns.UpdatedAt,
	)); err != nil {
		return nil, err
	}

	ff := &entity.FeatureFlag{}
	dao.FeatureFlagTo(featureFlagDB, ff)

	return ff, nil
}

func NewPostgres(db postgres.Connection) repository.FeatureFlag {
	return &postgresRepository{
		db: db,
	}
}
