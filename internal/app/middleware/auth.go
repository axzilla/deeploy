package middleware

import (
	"context"
	"net/http"

	"github.com/axzilla/deeploy/internal/app/cookie"
	"github.com/axzilla/deeploy/internal/app/jwt"
	jwtlib "github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := cookie.GetTokenFromCookie(r)
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims := jwtlib.MapClaims{}
		t, err := jwtlib.ParseWithClaims(token, claims, func(t *jwtlib.Token) (interface{}, error) {
			return jwt.JwtSecret, nil
		})

		if err != nil || !t.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func RequireAuth(next http.HandlerFunc, redirectTo ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := cookie.GetTokenFromCookie(r)
		if token == "" {
			path := "/login"
			if len(redirectTo) > 0 {
				path = redirectTo[0]
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
