package middlewares

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

type AccessTokenValidatorMiddleware struct {
	Key []byte
}

func (m *AccessTokenValidatorMiddleware) ExtractToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")

		if len(token) == 0 {
			cookie, err := r.Cookie("Authorization")

			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			token = cookie.Value
		}

		ctx := r.Context()

		ctx = context.WithValue(ctx, "access_token", token)

		_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return m.Key, nil
		})

		ctx = context.WithValue(ctx, "access_token_is_valid", err == nil)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RejectInvalidTokens(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		val := ctx.Value("access_token_is_valid")

		if val == nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		flag, ok := val.(bool)

		if !ok {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		if !flag {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
