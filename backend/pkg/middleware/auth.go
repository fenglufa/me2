package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/me2/pkg/errcode"
	"github.com/me2/pkg/response"
)

// JWTAuth JWT 认证中间件
func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 从 Header 获取 Token
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.HttpError(w, errcode.ErrUnauthorized)
				return
			}

			// 解析 Bearer Token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.HttpError(w, errcode.ErrUnauthorized)
				return
			}

			tokenString := parts[1]

			// 验证 Token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				response.HttpError(w, errcode.ErrInvalidToken)
				return
			}

			// 提取用户信息
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				response.HttpError(w, errcode.ErrInvalidToken)
				return
			}

			userID := int64(claims["user_id"].(float64))

			// 将用户 ID 存入 context
			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID 从 context 获取用户 ID
func GetUserID(ctx context.Context) int64 {
	if userID, ok := ctx.Value("user_id").(int64); ok {
		return userID
	}
	return 0
}
