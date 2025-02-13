package dao

import (
	"github.com/novychok/flagroll/platform/internal/database/pqmodels"
	"github.com/novychok/flagroll/platform/internal/entity"
)

func UserTo(userDB *pqmodels.User, user *entity.User) {
	user.ID = entity.UserID(userDB.ID)
	user.Name = userDB.Name
	user.Email = userDB.Email
	user.PasswordHash = userDB.PasswordHash
	user.CreatedAt = userDB.CreatedAt
	user.UpdatedAt = userDB.UpdatedAt
}
