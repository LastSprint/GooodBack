package common

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type CustomJWTClaims struct {
	jwt.RegisteredClaims
}

func (c *CustomJWTClaims) Valid() error {

	now := time.Now()

	if !c.VerifyExpiresAt(now, true) {
		delta := now.Sub(c.ExpiresAt.Time)
		return fmt.Errorf("token is expired by %v", delta)
	}

	if !c.VerifyIssuedAt(now, true) {
		return fmt.Errorf("token used before issued")
	}

	return nil
}
