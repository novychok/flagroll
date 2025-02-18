package platformapiv1

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/novychok/flagroll/platform/internal/entity"
	"github.com/novychok/flagroll/platform/internal/service"
	platformapiv1 "github.com/novychok/flagroll/platform/pkg/api/platform/v1"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"github.com/samber/lo"
)

//go:generate oapi-codegen --config=./oapi-codegen.yaml ../../../api/platform/openapi/v1.yaml
type handler struct {
	mu sync.Mutex

	connections map[string]*websocket.Conn

	authorizationService service.Authorization
	featureFlagService   service.FeatureFlag
	apiKeyService        service.APIKeys
	realtimeService      service.Realtime
}

const (
	tokenKey        = "token"
	refreshTokenKey = "refreshToken"
)

func (h *handler) GetUserByAPIKey(w http.ResponseWriter, r *http.Request) {
	response(w, http.StatusOK, UserFromContext(r.Context()))
}

func (h *handler) GetFeatureFlagByUserAndName(w http.ResponseWriter, r *http.Request, userID openapitypes.UUID, name string) {
	featureFlag, err := h.featureFlagService.GetByUserAndName(r.Context(), entity.UserID(userID.String()), name)
	if err != nil {
		errResponse(w, r, http.StatusNotFound, err.Error())

		return
	}

	rsp := platformapiv1.FeatureFlag{
		Id:          uuid.MustParse(featureFlag.ID.String()),
		CreatedAt:   featureFlag.CreatedAt,
		UpdatedAt:   featureFlag.UpdatedAt,
		OwnerId:     uuid.MustParse(featureFlag.OwnerID.String()),
		Name:        featureFlag.Name,
		Description: &featureFlag.Description,
		Active:      featureFlag.Active,
	}

	response(w, http.StatusOK, rsp)
}

func (h handler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	var req platformapiv1.APIKeyCreate

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errResponse(w, r, http.StatusInternalServerError, err.Error())

		return
	}

	token, err := r.Cookie(tokenKey)
	if err != nil {
		if err == http.ErrNoCookie {
			errResponse(w, r, http.StatusUnauthorized, "no token provided")
			return
		}
		errResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authorizationService.GetUserByToken(r.Context(), token.Value)
	if err != nil {
		errResponse(w, r, http.StatusUnauthorized, err.Error())
		return
	}

	apiKeyRequest := &entity.APIKeyCreate{
		ExpiresAt: lo.FromPtr(&req.ExpiresAt),
	}

	apiKey, err := h.apiKeyService.Create(r.Context(), user.ID, apiKeyRequest)
	if err != nil {
		errResponse(w, r, http.StatusInternalServerError, err.Error())

		return
	}

	response(w, http.StatusCreated, apiKey.ApiKeyRaw)
}

func (h handler) DeleteAPIKey(w http.ResponseWriter, r *http.Request,
	id openapitypes.UUID,
) {
	err := h.apiKeyService.Delete(r.Context(),
		entity.APIKeyID(id.String()))
	if err != nil {
		errResponse(w, r, http.StatusInternalServerError, err.Error())

		return
	}

	response(w, http.StatusNoContent, nil)
}

func (h *handler) ListFeatureFlags(w http.ResponseWriter, r *http.Request) {
	featureFlags, err := h.featureFlagService.GetAll(r.Context())
	if err != nil {
		errResponse(w, r, http.StatusInternalServerError, err.Error())

		return
	}

	var rsp []platformapiv1.FeatureFlagResponse
	for _, featureFlag := range featureFlags {
		rsp = append(rsp, platformapiv1.FeatureFlagResponse{
			OwnerId:     uuid.MustParse(featureFlag.OwnerID.String()),
			Name:        featureFlag.Name,
			Description: featureFlag.Description,
			Active:      featureFlag.Active,
		})
	}

	response(w, http.StatusOK, rsp)
}

func (h *handler) CreateFeatureFlag(w http.ResponseWriter, r *http.Request) {
	var req platformapiv1.FeatureFlagCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errResponse(w, r, http.StatusBadRequest, err.Error())

		return
	}

	token, err := r.Cookie(tokenKey)
	if err != nil {
		if err == http.ErrNoCookie {
			errResponse(w, r, http.StatusUnauthorized, "no token provided")
			return
		}
		errResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authorizationService.GetUserByToken(r.Context(), token.Value)
	if err != nil {
		errResponse(w, r, http.StatusUnauthorized, err.Error())
		return
	}

	featureFlagCreate := &entity.FeatureFlagCreate{
		OwnerID:     user.ID,
		Name:        req.Name,
		Description: *req.Description,
		Active:      req.Active,
	}

	featureFlag, err := h.featureFlagService.Create(r.Context(), featureFlagCreate)
	if err != nil {
		errResponse(w, r, http.StatusInternalServerError, err.Error())

		return
	}

	rsp := platformapiv1.FeatureFlagResponse{
		OwnerId:     uuid.MustParse(featureFlag.OwnerID.String()),
		Name:        featureFlag.Name,
		Description: featureFlag.Description,
		Active:      featureFlag.Active,
	}

	response(w, http.StatusCreated, rsp)
}

func (h *handler) GetFeatureFlag(w http.ResponseWriter,
	r *http.Request, id openapitypes.UUID) {
	featureFlag, err := h.featureFlagService.GetByID(r.Context(),
		entity.FeatureFlagID(id.String()))
	if err != nil {
		errResponse(w, r, http.StatusNotFound, err.Error())

		return
	}

	rsp := platformapiv1.FeatureFlagResponse{
		OwnerId:     uuid.MustParse(featureFlag.OwnerID.String()),
		Name:        featureFlag.Name,
		Description: featureFlag.Description,
		Active:      featureFlag.Active,
	}

	response(w, http.StatusOK, rsp)
}

func (h *handler) DeleteFeatureFlag(w http.ResponseWriter,
	r *http.Request, id openapitypes.UUID) {
	err := h.featureFlagService.Delete(r.Context(), entity.FeatureFlagID(id.String()))
	if err != nil {
		errResponse(w, r, http.StatusNotFound, err.Error())

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) UpdateFeatureFlag(w http.ResponseWriter,
	r *http.Request, id openapitypes.UUID) {
	var req platformapiv1.FeatureFlagUpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errResponse(w, r, http.StatusBadRequest, err.Error())

		return
	}

	featureFlagUpdate := &entity.FeatureFlagUpdate{
		Name:        req.Name,
		Description: *req.Description,
		Active:      req.Active,
	}

	featureFlag, err := h.featureFlagService.Update(r.Context(),
		entity.FeatureFlagID(id.String()), featureFlagUpdate)
	if err != nil {
		errResponse(w, r, http.StatusInternalServerError, err.Error())

		return
	}

	rsp := platformapiv1.FeatureFlagResponse{
		OwnerId:     uuid.MustParse(featureFlag.OwnerID.String()),
		Name:        featureFlag.Name,
		Description: featureFlag.Description,
		Active:      featureFlag.Active,
	}

	response(w, http.StatusOK, rsp)
}

func (h *handler) UpdateFeatureFlagToggle(w http.ResponseWriter,
	r *http.Request, id openapitypes.UUID) {
	var req platformapiv1.FeatureFlagToggleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errResponse(w, r, http.StatusBadRequest, err.Error())

		return
	}

	featureFlag, err := h.featureFlagService.UpdateToggle(r.Context(),
		entity.FeatureFlagID(id.String()), req.Active)
	if err != nil {
		errResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	rsp := platformapiv1.FeatureFlagResponse{
		OwnerId:     uuid.MustParse(featureFlag.OwnerID.String()),
		Name:        featureFlag.Name,
		Description: featureFlag.Description,
		Active:      featureFlag.Active,
	}

	response(w, http.StatusOK, rsp)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var request platformapiv1.LoginRequest

	err := render.DefaultDecoder(r, &request)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.DefaultResponder(w, r, platformapiv1.Error{
			Message: "invalid credentials",
		})

		return
	}

	token, err := h.authorizationService.Login(r.Context(), &entity.Login{
		Email:    string(request.Email),
		Password: request.Password,
	})
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.DefaultResponder(w, r, platformapiv1.Error{
			Message: "invalid credentials",
		})

		return
	}

	tokenCookie := &http.Cookie{
		Name:     tokenKey,
		Value:    token.Token,
		Path:     "/",
		Expires:  token.TokenExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
	}

	refTokenCookie := &http.Cookie{
		Name:     refreshTokenKey,
		Value:    token.RefreshToken,
		Path:     "/",
		Expires:  token.RefreshTokenExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
	}

	http.SetCookie(w, tokenCookie)
	http.SetCookie(w, refTokenCookie)

	render.DefaultResponder(w, r, platformapiv1.Message{
		Message: "ok",
	})
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var request platformapiv1.RegistrationRequest

	err := render.DefaultDecoder(r, &request)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.DefaultResponder(w, r, platformapiv1.Error{
			Message: "invalid data",
		})

		return
	}

	token, err := h.authorizationService.Register(r.Context(), &entity.UserCreate{
		Name:            request.Name,
		Email:           string(request.Email),
		Password:        request.Password,
		PasswordConfirm: request.PasswordConfirm,
	})
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.DefaultResponder(w, r, platformapiv1.Error{
			Message: "invalid data",
		})

		return
	}

	render.DefaultResponder(w, r, platformapiv1.Token{
		RefreshToken:          token.RefreshToken,
		RefreshTokenExpiresAt: token.RefreshTokenExpiresAt,
		Token:                 token.Token,
		TokenExpiresAt:        token.TokenExpiresAt,
	})
}

func NewHandler(
	authorizationService service.Authorization,
	featureFlagService service.FeatureFlag,
	apiKeyService service.APIKeys,
	realtimeService service.Realtime,
) platformapiv1.ServerInterface {
	handler := &handler{
		connections: make(map[string]*websocket.Conn),

		authorizationService: authorizationService,
		featureFlagService:   featureFlagService,
		apiKeyService:        apiKeyService,
		realtimeService:      realtimeService,
	}

	go func() {
		if err := handler.broadcastMessages(); err != nil {
			log.Println(err)
		}
	}() // todo handle error and call it as handler listener

	return handler
}
