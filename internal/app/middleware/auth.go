package middleware

import (
	"context"
	"net/http"

	"github.com/axzilla/deeploy/internal/app/cookie"
	"github.com/axzilla/deeploy/internal/app/jwt"
	"github.com/axzilla/deeploy/internal/app/services"
)

type AuthMiddleWare struct {
	userService services.UserServiceInterface
}

func NewAuthMiddleware(userService services.UserServiceInterface) *AuthMiddleWare {
	return &AuthMiddleWare{userService: userService}
}

func (m *AuthMiddleWare) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := cookie.GetTokenFromCookie(r)
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		t, claims, err := jwt.ValidateToken(token)
		if err != nil || !t.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID := claims["user_id"].(string)
		user, err := m.userService.GetUserByID(userID)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func RequireAuth(next http.HandlerFunc, redirectTo ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := cookie.GetTokenFromCookie(r)
		if token == "" {
			path := "/"
			if len(redirectTo) > 0 {
				path = redirectTo[0]
			}

			if r.URL.RawQuery != "" {
				path += "?" + r.URL.RawQuery
			}

			http.Redirect(w, r, path, http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func RequireGuest(next http.HandlerFunc, redirectTo ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := cookie.GetTokenFromCookie(r)
		if token != "" {
			path := "/dashboard"
			if len(redirectTo) > 0 {
				path = redirectTo[0]
			}
			http.Redirect(w, r, path, http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}
