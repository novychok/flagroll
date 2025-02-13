package platformapiv1

import (
	"net/http"

	"github.com/novychok/flagroll/platform/internal/entity"
)

func setCookies(token *entity.Token, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     tokenKey,
		Value:    token.Token,
		Expires:  token.TokenExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenKey,
		Value:    token.RefreshToken,
		Expires:  token.RefreshTokenExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}
