package featureflag

import (
	"encoding/json"
	"log/slog"

	"context"

	"github.com/go-playground/validator/v10"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/novychok/flagroll/platform/internal/entity"
	"github.com/novychok/flagroll/platform/internal/repository"
	"github.com/novychok/flagroll/platform/internal/service"
)

type srv struct {
	l *slog.Logger
	v *validator.Validate

	featureFlagRepo repository.FeatureFlag
	natsPublisher   jetstream.JetStream
}

func (s *srv) GetByUserAndName(ctx context.Context, userID entity.UserID, name string) (*entity.FeatureFlag, error) {
	l := s.l.With(slog.String("method", "GetByUserAndName"))

	featureFlag, err := s.featureFlagRepo.GetByUserAndName(ctx, userID, name)
	if err != nil {
		l.Error("failed to get feature flag by user and name", slog.Any("error", err))

		return nil, err
	}

	return featureFlag, nil
}

func (s *srv) GetAll(ctx context.Context) ([]*entity.FeatureFlag, error) {
	l := s.l.With(slog.String("method", "GetAll"))

	featureFlags, err := s.featureFlagRepo.GetAll(ctx)
	if err != nil {
		l.Error("failed to get all feature flags", slog.Any("error", err))

		return nil, err
	}

	return featureFlags, nil
}

func (s *srv) Create(ctx context.Context,
	featureFlagCreate *entity.FeatureFlagCreate) (*entity.FeatureFlag, error) {
	l := s.l.With(slog.String("method", "Create"))

	if err := s.v.Struct(featureFlagCreate); err != nil {
		l.Error("validation failed", slog.Any("error", err))

		return nil, err
	}

	createdFeatureFlag, err := s.featureFlagRepo.Create(ctx, featureFlagCreate)
	if err != nil {
		l.Error("failed to create feature flag", slog.Any("error", err))

		return nil, err
	}

	return createdFeatureFlag, nil
}

func (s *srv) GetByID(ctx context.Context,
	id entity.FeatureFlagID) (*entity.FeatureFlag, error) {
	l := s.l.With(slog.String("method", "GetByID"))

	featureFlag, err := s.featureFlagRepo.GetByID(ctx, id)
	if err != nil {
		l.Error("failed to get feature flag by ID", slog.Any("error", err))

		return nil, err
	}

	return featureFlag, nil
}

func (s *srv) Delete(ctx context.Context,
	id entity.FeatureFlagID) error {
	l := s.l.With(slog.String("method", "Delete"))

	err := s.featureFlagRepo.Delete(ctx, id)
	if err != nil {
		l.Error("failed to delete feature flag", slog.Any("error", err))

		return err
	}

	return nil
}

func (s *srv) Update(ctx context.Context, id entity.FeatureFlagID,
	featureFlag *entity.FeatureFlagUpdate) (*entity.FeatureFlag, error) {
	l := s.l.With(slog.String("method", "Update"))

	if err := s.v.Struct(featureFlag); err != nil {
		l.Error("validation failed", slog.Any("error", err))

		return nil, err
	}

	updatedFeatureFlag, err := s.featureFlagRepo.Update(ctx, id, featureFlag)
	if err != nil {
		l.Error("failed to update feature flag", slog.Any("error", err))

		return nil, err
	}

	return updatedFeatureFlag, nil
}

func (s *srv) UpdateToggle(ctx context.Context,
	id entity.FeatureFlagID, active bool) (*entity.FeatureFlag, error) {
	l := s.l.With(slog.String("method", "UpdateToggle"))

	updatedFeatureFlag, err := s.featureFlagRepo.UpdateToggle(ctx, id, active)
	if err != nil {
		l.Error("failed to update feature flag toggle", slog.Any("error", err))

		return nil, err
	}

	// TODO: change to this payload
	// payload := fmt.Sprintf(`"ID":"%s","Name":"%s","Active":%t`,
	// 	updatedFeatureFlag.ID, updatedFeatureFlag.Name, updatedFeatureFlag.Active)

	payload, err := json.Marshal(updatedFeatureFlag)
	if err != nil {
		l.Error("failed to marshal updated feature flag", slog.Any("error", err))
		return nil, err
	}

	_, err = s.natsPublisher.Publish(ctx, "features.toggle.update", []byte(payload))
	if err != nil {
		s.l.ErrorContext(ctx, "failed to publish a message", "err", err)

		return nil, err
	}

	return updatedFeatureFlag, nil
}

func New(
	l *slog.Logger,
	v *validator.Validate,

	featureFlagRepo repository.FeatureFlag,
	natsPublisher jetstream.JetStream,
) service.FeatureFlag {
	return &srv{
		l: l,
		v: v,

		featureFlagRepo: featureFlagRepo,
		natsPublisher:   natsPublisher,
	}
}
