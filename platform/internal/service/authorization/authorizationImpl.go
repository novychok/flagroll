package authorization

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"log/slog"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/novychok/flagroll/platform/internal/entity"
	"github.com/novychok/flagroll/platform/internal/pkg/jwts"
	"github.com/novychok/flagroll/platform/internal/repository"
	"github.com/novychok/flagroll/platform/internal/service"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenDuration        = 15 * time.Minute
	refreshTokenDuration = 24 * time.Hour
	issuer               = "nazyk"
)

type srv struct {
	l *slog.Logger
	v *validator.Validate

	userRepository   repository.User
	apikeyRepository repository.APIKey
	jwtSecretManager *jwts.SecretManager
}

const (
	TokenPayloadSize = 60
	ApiKeyRawPrefix  = "ser."
)

func (s *srv) GetUserByApiKey(ctx context.Context, apiKeyRaw string) (*entity.User, error) {
	l := s.l.With(slog.String("method", "GetUserByApiKey"))

	tokenWithoutPrefix := strings.TrimPrefix(apiKeyRaw, ApiKeyRawPrefix)
	tkLengthWithNoPrefix := 128

	if len(tokenWithoutPrefix) != tkLengthWithNoPrefix {
		l.ErrorContext(ctx, "invalid token length", slog.Any("error", entity.ErrInvalidaTokenLength))

		return nil, entity.ErrInvalidaTokenLength
	}

	encodedTokenID := tokenWithoutPrefix[:len(tokenWithoutPrefix)-80]

	tokenID, err := base64.StdEncoding.DecodeString(encodedTokenID)
	if err != nil {
		l.ErrorContext(ctx, "decode toke id failed", slog.Any("error", err))

		return nil, err
	}

	apiKey, err := s.apikeyRepository.GetByTokenID(ctx, string(tokenID))
	if err != nil {
		l.ErrorContext(ctx, "get by token id failed", slog.Any("error", err))

		return nil, err
	}

	sha256Hash := sha256.Sum256([]byte(apiKeyRaw))
	if err := bcrypt.CompareHashAndPassword([]byte(apiKey.TokenHash), sha256Hash[:]); err != nil {
		l.ErrorContext(ctx, "compare hash failed", slog.Any("error", err))

		return nil, entity.ErrInvalidTokenHash
	}

	user, err := s.userRepository.GetByID(ctx, entity.UserID(apiKey.OwnerID))
	if err != nil {
		l.Error("get user by id failed", slog.Any("error", err))

		return nil, err
	}

	return user, nil
}

func (s *srv) GetUserByToken(ctx context.Context, token string) (*entity.User, error) {
	l := s.l.With(slog.String("method", "GetUserByToken"))

	claims, err := s.getClaims(token)
	if err != nil {
		l.Error("get claims failed", slog.Any("error", err))

		return nil, err
	}

	user, err := s.userRepository.GetByID(ctx, entity.UserID(claims.Subject))
	if err != nil {
		l.Error("get user by id failed", slog.Any("error", err))

		return nil, err
	}

	return user, nil
}

func (s *srv) getClaims(token string) (*entity.Claims, error) {
	claims := &entity.Claims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecretManager.PublicKey(), nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (s *srv) VerifyToken(ctx context.Context, verifyRequest *entity.VerifyToken) error {
	l := s.l.With(slog.String("method", "VerifyToken"))

	err := s.v.StructCtx(ctx, verifyRequest)
	if err != nil {
		l.Error("validation failed", slog.Any("error", err))

		return err
	}

	_, err = s.getClaims(verifyRequest.Token)
	if err != nil {
		l.Error("get claims failed", slog.Any("error", err))

		return err
	}

	return nil
}

func (s *srv) generateToken(user *entity.User) (*entity.Token, error) {
	now := time.Now()
	tokenID := uuid.New().String()
	refreshTokenID := uuid.New().String()
	tokenExpiresAt := now.Add(tokenDuration)
	refreshTokenExpiresAt := now.Add(refreshTokenDuration)

	claims := entity.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(tokenExpiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        tokenID,
		},
	}

	refreshClaims := entity.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        refreshTokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)

	tokenString, err := token.SignedString(s.jwtSecretManager.PrivateKey())
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := refreshToken.SignedString(s.jwtSecretManager.PrivateKey())
	if err != nil {
		return nil, err
	}

	return &entity.Token{
		Token:                 tokenString,
		TokenExpiresAt:        tokenExpiresAt,
		RefreshToken:          refreshTokenString,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}, nil
}

func (s *srv) Login(ctx context.Context, login *entity.Login) (*entity.Token, error) {
	l := s.l.With(slog.String("method", "LogIn"))

	err := s.v.StructCtx(ctx, login)
	if err != nil {
		l.Error("validation failed", slog.Any("error", err))

		return nil, err
	}

	user, err := s.userRepository.GetByEmail(ctx, login.Email)
	if err != nil {
		l.Error("get user by email failed", slog.Any("error", err))

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(login.Password))
	if err != nil {
		l.Error("passwords do not match", slog.Any("error", err))

		return nil, err
	}

	token, err := s.generateToken(user)
	if err != nil {
		l.Error("generate token failed", slog.Any("error", err))

		return nil, err
	}

	return token, nil
}

func (s *srv) RefreshToken(ctx context.Context, req *entity.RefreshToken) (*entity.Token, error) {
	l := s.l.With(slog.String("method", "RefreshToken"))

	err := s.v.StructCtx(ctx, req)
	if err != nil {
		l.Error("validation failed", slog.Any("error", err))

		return nil, err
	}

	claims := &entity.Claims{}

	_, err = jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecretManager.PublicKey(), nil
	})
	if err != nil {
		l.Error("parse refresh token failed", slog.Any("error", err))

		return nil, err
	}

	user, err := s.userRepository.GetByID(ctx, entity.UserID(claims.Subject))
	if err != nil {
		l.Error("get user by id failed", slog.Any("error", err))

		return nil, err
	}

	token, err := s.generateToken(user)
	if err != nil {
		l.Error("generate token failed", slog.Any("error", err))

		return nil, err
	}

	return token, nil
}

func (s *srv) Register(ctx context.Context, userCreate *entity.UserCreate) (*entity.Token, error) {
	l := s.l.With(slog.String("method", "Register"))

	err := s.v.StructCtx(ctx, userCreate)
	if err != nil {
		l.Error("validation failed", slog.Any("error", err))

		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(userCreate.Password), bcrypt.DefaultCost)
	if err != nil {
		l.Error("generate password hash failed", slog.Any("error", err))

		return nil, err
	}

	user := &entity.User{
		Name:         userCreate.Name,
		Email:        userCreate.Email,
		PasswordHash: string(hash),
	}

	err = s.userRepository.Create(ctx, user)
	if err != nil {
		l.Error("create user failed", slog.Any("error", err))

		return nil, err
	}

	token, err := s.generateToken(user)
	if err != nil {
		l.Error("generate token failed", slog.Any("error", err))

		return nil, err
	}

	return token, nil
}

func New(
	l *slog.Logger,
	v *validator.Validate,

	userRepository repository.User,
	apikeyRepository repository.APIKey,
	jwtSecretManager *jwts.SecretManager,
) service.Authorization {
	return &srv{
		l: l.With(slog.String("service", "authorization")),
		v: v,

		userRepository:   userRepository,
		apikeyRepository: apikeyRepository,
		jwtSecretManager: jwtSecretManager,
	}
}
