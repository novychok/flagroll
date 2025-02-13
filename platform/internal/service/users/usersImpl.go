package users

import (
	"context"
	"log/slog"

	"github.com/novychok/flagroll/platform/internal/entity"
	"github.com/novychok/flagroll/platform/internal/repository"
	"github.com/novychok/flagroll/platform/internal/service"
)

type srv struct {
	l *slog.Logger

	userRepository repository.User
}

func (s *srv) Get(ctx context.Context, userID entity.UserID) (*entity.User, error) {
	l := s.l.With(slog.String("method", "Get"))

	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		l.Error("get user by id failed", slog.Any("error", err))

		return nil, err
	}

	return user, nil
}

func New(
	l *slog.Logger,

	userRepository repository.User,
) service.Users {
	return &srv{
		l: l.With(slog.String("service", "users")),

		userRepository: userRepository,
	}
}
