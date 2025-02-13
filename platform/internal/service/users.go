package service

import (
	"context"

	"github.com/novychok/flagroll/platform/internal/entity"
)

type Users interface {
	Get(ctx context.Context, userID entity.UserID) (*entity.User, error)
}
