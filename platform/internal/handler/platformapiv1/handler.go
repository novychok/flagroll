package platformapiv1

import (
	"net/http"

	platformapiv1 "github.com/novychok/flagroll/platform/pkg/api/platform/v1"
	openapitypes "github.com/oapi-codegen/runtime/types"
)

//go:generate oapi-codegen --config=./oapi-codegen.yaml ../../../api/platform/openapi/v1.yaml
type handler struct {
}

func (h *handler) ListFeatureFlags(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) CreateFeatureFlag(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) GetFeatureFlag(w http.ResponseWriter, r *http.Request, id openapitypes.UUID) {

}

func (h *handler) DeleteFeatureFlag(w http.ResponseWriter, r *http.Request, id openapitypes.UUID) {

}

func (h *handler) UpdateFeatureFlag(w http.ResponseWriter, r *http.Request, id openapitypes.UUID) {

}

func (h *handler) UpdateFeatureFlagToggle(w http.ResponseWriter, r *http.Request, id openapitypes.UUID) {

}

func NewHandler() platformapiv1.ServerInterface {
	return &handler{}
}
