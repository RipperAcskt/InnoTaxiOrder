package service

import (
	"fmt"
	"time"

	"github.com/RipperAcskt/innotaxiorder/config"
	"github.com/golang-jwt/jwt"
)

var (
	ErrTokenExpired = fmt.Errorf("token expired")
	ErrTokenClaims  = fmt.Errorf("jwt map claims failed")
	ErrTokenId      = fmt.Errorf("jwt get user id failed")
	ErrTokenType    = fmt.Errorf("jwt get user type failed")
	ErrValidation   = fmt.Errorf("validation failed")
)

func Verify(token string, cfg *config.Config) (string, error) {
	tokenJwt, err := jwt.Parse(
		token,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.HS256_SECRET), nil
		},
	)

	if err != nil {
		return "", fmt.Errorf("token parse failed: %w", err)
	}

	claims, ok := tokenJwt.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrTokenClaims
	}

	if !claims.VerifyExpiresAt(time.Now().UTC().Unix(), true) {
		return "", ErrTokenExpired
	}

	userType, ok := claims["type"]
	if !ok {
		return "", ErrTokenType
	}
	str, ok := userType.(string)
	if !ok {
		return "", ErrTokenType
	}
	if str != "user" {
		return "", ErrValidation
	}

	id, ok := claims["user_id"]
	if !ok {
		return "", ErrTokenId
	}

	return fmt.Sprint(id), nil
}
