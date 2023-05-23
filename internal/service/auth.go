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

type UserInfo struct {
	Id   string
	Type string
}

func Verify(token string, cfg *config.Config) (*UserInfo, error) {
	tokenJwt, err := jwt.Parse(
		token,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.HS256_SECRET), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("token parse failed: %w", err)
	}

	claims, ok := tokenJwt.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrTokenClaims
	}

	if !claims.VerifyExpiresAt(time.Now().UTC().Unix(), true) {
		return nil, ErrTokenExpired
	}

	userType, ok := claims["type"]
	if !ok {
		return nil, ErrTokenType
	}
	typeStr, ok := userType.(string)
	if !ok {
		return nil, ErrTokenType
	}

	id, ok := claims["user_id"]
	if !ok {
		return nil, ErrTokenId
	}

	return &UserInfo{
		fmt.Sprint(id),
		typeStr,
	}, nil
}
