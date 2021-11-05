package middlewares

import (
	"context"
	"fmt"
	"github.com/LastSprint/GooodBack/common"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"time"
)

type TokenRefresher interface {
	RefreshAccessToken(refreshTokenStr string) (string, error)
}

type AccessTokenValidatorMiddleware struct {
	Key []byte
	Refresher TokenRefresher
}

func (m *AccessTokenValidatorMiddleware) ExtractToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")

		if len(token) == 0 {
			cookie, err := r.Cookie("Authorization")
			fmt.Println(r.Cookies())

			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			token = cookie.Value
		}

		ctx := r.Context()

		ctx = context.WithValue(ctx, "access_token", token)

		claims := common.CustomJWTClaims{}

		_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
			return m.Key, nil
		})

		ctx = context.WithValue(ctx, ContextKeyUserId, claims.ID)

		if err == nil {
			ctx = context.WithValue(ctx, "access_token_is_valid", true)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if !claims.VerifyIssuedAt(time.Now(), true) {
			ctx = context.WithValue(ctx, "access_token_is_valid", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		cookie, err := r.Cookie("Refreshing")

		if err != nil || len(cookie.Value) == 0 {
			ctx = context.WithValue(ctx, "access_token_is_valid", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		newAccessToken, err := m.Refresher.RefreshAccessToken(cookie.Value)

		if err != nil {
			ctx = context.WithValue(ctx, "access_token_is_valid", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if err = common.SetCookie(w, r, "Authorization", newAccessToken, 60 * 60 * 24 * 30); err != nil {
			log.Println("[ERR] Couldn't update auth cookie ->", err.Error())
			ctx = context.WithValue(ctx, "access_token_is_valid", false)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		ctx = context.WithValue(ctx, "access_token", token)
		ctx = context.WithValue(ctx, "access_token_is_valid", true)
		next.ServeHTTP(w, r.WithContext(ctx))
		return
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
