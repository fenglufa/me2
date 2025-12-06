package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/me2/pkg/errcode"
	"github.com/me2/pkg/response"
)

type AuthMiddleware struct {
	secret string
}

func NewAuthMiddleware(secret string) func(http.HandlerFunc) http.HandlerFunc {
	m := &AuthMiddleware{secret: secret}
	return m.Handle
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.HttpError(w, errcode.ErrUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.HttpError(w, errcode.ErrUnauthorized)
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(m.secret), nil
		})

		if err != nil || !token.Valid {
			response.HttpError(w, errcode.ErrInvalidToken)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.HttpError(w, errcode.ErrInvalidToken)
			return
		}

		userID := int64(claims["user_id"].(float64))
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next(w, r.WithContext(ctx))
	}
}
