package dao

import (
	"github.com/novychok/flagroll/platform/internal/database/pqmodels"
	"github.com/novychok/flagroll/platform/internal/entity"
)

func FeatureFlagTo(featureFlagDB *pqmodels.FeatureFlag, featureFlag *entity.FeatureFlag) {
	featureFlag.ID = entity.FeatureFlagID(featureFlagDB.ID)
	featureFlag.OwnerID = entity.UserID(featureFlagDB.OwnerID)
	featureFlag.Name = featureFlagDB.Name
	featureFlag.Active = featureFlagDB.Active
	featureFlag.Description = featureFlagDB.Description.String
	featureFlag.CreatedAt = featureFlagDB.CreatedAt
	featureFlag.UpdatedAt = featureFlagDB.UpdatedAt
}
