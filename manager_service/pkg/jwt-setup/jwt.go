package jwt_setup

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"manager_service/internal/config"
)

var signKey *rsa.PrivateKey

type DTO interface {
}

type RegisteredClaims struct {
	jwt.RegisteredClaims
	Id string `json:"id"`
}

func ParseToken(tokenString string) (userId string, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().Keys.JWTSignKey), nil
	})

	if claims, ok := token.Claims.(*RegisteredClaims); ok && token.Valid {
		return claims.Id, nil
	} else {
		return "", fmt.Errorf("wrong token: %v", claims.RegisteredClaims.Issuer)

	}
}
