package service

import (
	"context"

	"github.com/novychok/flagroll/platform/internal/entity"
)

type Authorization interface {
	Login(ctx context.Context, login *entity.Login) (*entity.Token, error)
	RefreshToken(ctx context.Context, refreshRequest *entity.RefreshToken) (*entity.Token, error)
	Register(ctx context.Context, user *entity.UserCreate) (*entity.Token, error)
	VerifyToken(ctx context.Context, verifyRequest *entity.VerifyToken) error
	GetUserByToken(ctx context.Context, token string) (*entity.User, error)
	GetUserByApiKey(ctx context.Context, apiKeyRaw string) (*entity.User, error)
}
