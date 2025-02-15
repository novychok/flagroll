package dao

import (
	"github.com/novychok/flagroll/platform/internal/database/pqmodels"
	"github.com/novychok/flagroll/platform/internal/entity"
)

func APIKeyTo(apiKeyDB *pqmodels.APIKey, apiKey *entity.APIKey) {
	apiKey.ID = entity.APIKeyID(apiKeyDB.ID)
	apiKey.OwnerID = entity.UserID(apiKeyDB.OwnerID)
	apiKey.TokenID = apiKeyDB.TokenID
	apiKey.TokenHash = apiKeyDB.TokenHash
	apiKey.CreatedAt = apiKeyDB.CreatedAt
	apiKey.ExpiresAt = &apiKeyDB.ExpiresAt
}
