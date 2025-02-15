package apikeys

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/novychok/flagroll/platform/internal/entity"
	"github.com/novychok/flagroll/platform/internal/repository"
	"github.com/novychok/flagroll/platform/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type srv struct {
	l *slog.Logger
	v *validator.Validate

	apiKeyRepository repository.APIKey
}

const (
	TokenPayloadSize = 60
	ApiKeyRawPrefix  = "ser."
)

func (s *srv) Create(ctx context.Context, ownerID entity.UserID,
	apiKeyCreate *entity.APIKeyCreate) (*entity.APIKeyResponse, error) {
	l := s.l.With(slog.String("method", "Create"))

	err := s.v.StructCtx(ctx, apiKeyCreate)
	if err != nil {
		l.ErrorContext(ctx, "validation failed", slog.Any("error", err))

		return nil, err
	}

	tokenID := uuid.New().String()

	apiKeyRaw, err := s.createToken(tokenID)
	if err != nil {
		l.ErrorContext(ctx, "create token raw failed", slog.Any("error", err))

		return nil, err
	}

	tokenHash, err := s.createTokenHash(apiKeyRaw)
	if err != nil {
		l.ErrorContext(ctx, "create token hash failed", slog.Any("error", err))

		return nil, err
	}

	newAPIKey := &entity.APIKey{
		OwnerID:   ownerID,
		TokenID:   tokenID,
		TokenHash: tokenHash,
		ExpiresAt: apiKeyCreate.ExpiresAt,
	}

	apiKey, err := s.apiKeyRepository.Create(ctx, newAPIKey)
	if err != nil {
		l.ErrorContext(ctx, "create api key failed", slog.Any("error", err))

		return nil, err
	}

	return &entity.APIKeyResponse{
		ID:        apiKey.ID,
		ApiKeyRaw: apiKeyRaw,
		CreatedAt: apiKey.CreatedAt,
		ExpiresAt: apiKey.ExpiresAt,
	}, nil
}

// func (s *srv) GetByTokenID(ctx context.Context,
// 	apiKeyRaw string) (*entity.APIKey, error) {
// 	l := s.l.With(slog.String("method", "GetByTokenID"))

// 	tokenWithoutPrefix := strings.TrimPrefix(apiKeyRaw, ApiKeyRawPrefix)
// 	tkLengthWithNoPrefix := 128

// 	if len(tokenWithoutPrefix) != tkLengthWithNoPrefix {
// 		l.ErrorContext(ctx, "invalid token length", slog.Any("error", entity.ErrInvalidaTokenLength))

// 		return nil, entity.ErrInvalidaTokenLength
// 	}

// 	encodedTokenID := tokenWithoutPrefix[:len(tokenWithoutPrefix)-80]

// 	tokenID, err := base64.StdEncoding.DecodeString(encodedTokenID)
// 	if err != nil {
// 		l.ErrorContext(ctx, "decode toke id failed", slog.Any("error", err))

// 		return nil, err
// 	}

// 	apiKey, err := s.apiKeyRepository.GetByTokenID(ctx, string(tokenID))
// 	if err != nil {
// 		l.ErrorContext(ctx, "get by token id failed", slog.Any("error", err))

// 		return nil, err
// 	}

// 	sha256Hash := sha256.Sum256([]byte(apiKeyRaw))
// 	if err := bcrypt.CompareHashAndPassword([]byte(apiKey.TokenHash), sha256Hash[:]); err != nil {
// 		l.ErrorContext(ctx, "compare hash failed", slog.Any("error", err))

// 		return nil, entity.ErrInvalidTokenHash
// 	}

// 	return apiKey, nil
// }

func (s *srv) Get(ctx context.Context, id entity.APIKeyID) ([]*entity.APIKey, error) {

	return nil, nil
}

func (s *srv) Delete(ctx context.Context, id entity.APIKeyID) error {
	l := s.l.With(slog.String("method", "Delete"))

	err := s.apiKeyRepository.Delete(ctx, id)
	if err != nil {
		l.ErrorContext(ctx, "delete api key failed", slog.Any("error", err))

		return err
	}

	return nil
}

func (s *srv) createToken(tokenID string) (string, error) {
	tokenHeader := ApiKeyRawPrefix
	encodedAPIKeyID := base64.StdEncoding.EncodeToString([]byte(tokenID))
	tokenPayload := make([]byte, TokenPayloadSize)

	_, err := rand.Read(tokenPayload)
	if err != nil {
		return "", err
	}

	stringTokenPayload := base64.StdEncoding.EncodeToString(tokenPayload)
	token := fmt.Sprintf("%s%s%s", tokenHeader, encodedAPIKeyID, stringTokenPayload)

	fmt.Println("created token:-", token)

	return token, nil
}

func (s *srv) createTokenHash(apiKeyRaw string) (string, error) {
	sha256Hash := sha256.Sum256([]byte(apiKeyRaw))

	tokenHash, err := bcrypt.GenerateFromPassword(sha256Hash[:], bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	fmt.Println("created hash:-", string(tokenHash))

	return string(tokenHash), nil
}

func New(
	l *slog.Logger,
	v *validator.Validate,

	apiKeyRepository repository.APIKey,
) service.APIKeys {
	return &srv{
		l: l.With(slog.String("service", "apiKeys")),
		v: v,

		apiKeyRepository: apiKeyRepository,
	}
}
