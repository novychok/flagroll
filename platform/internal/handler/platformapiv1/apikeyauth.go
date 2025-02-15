package platformapiv1

import (
	"net/http"

	"github.com/go-chi/render"
	platformapiv1 "github.com/novychok/flagroll/platform/pkg/api/platform/v1"
)

func (s *Server) keyAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customCtx := ContextFromRequest(r)
		_, exists := customCtx.Get("apiKeyRequired")
		if !exists {
			h.ServeHTTP(w, r.WithContext(r.Context()))

			return
		}

		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			render.Status(r, http.StatusUnauthorized)
			render.DefaultResponder(w, r, platformapiv1.Error{
				Message: "API key required",
			})

			return
		}

		user, err := s.authorizationService.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.DefaultResponder(w, r, platformapiv1.Error{
				Message: "Invalid API key",
			})

			return
		}

		ctx := WithUser(r.Context(), user)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
