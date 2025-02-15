package platformapiv1

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/novychok/flagroll/platform/internal/entity"
	platformapiv1 "github.com/novychok/flagroll/platform/pkg/api/platform/v1"
)

func (s *Server) auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customCtx := ContextFromRequest(r)
		_, exists := customCtx.Get("authRequired")
		if !exists {
			h.ServeHTTP(w, r.WithContext(r.Context()))

			return
		}

		var tokenValue string
		// get token from cookie
		tokenCookie, err := r.Cookie(tokenKey)
		if err != nil {
			// try to refresh it if refresh token exists
			refreshTokenCookie, err := r.Cookie(refreshTokenKey)
			if err != nil {
				render.Status(r, http.StatusUnauthorized)
				render.DefaultResponder(w, r, platformapiv1.Error{
					Message: "unauthorized",
				})

				return
			}

			token, err := s.authorizationService.RefreshToken(r.Context(), &entity.RefreshToken{
				Token: refreshTokenCookie.Value,
			})
			if err != nil {
				render.Status(r, http.StatusUnauthorized)
				render.DefaultResponder(w, r, platformapiv1.Error{
					Message: "unauthorized",
				})

				return
			}

			setCookies(token, w)

			tokenValue = token.Token

		} else {
			tokenValue = tokenCookie.Value
		}

		user, err := s.authorizationService.GetUserByToken(r.Context(), tokenValue)
		if err != nil {
			h.ServeHTTP(w, r.WithContext(r.Context()))

			return
		}

		ctx := WithUser(r.Context(), user)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
