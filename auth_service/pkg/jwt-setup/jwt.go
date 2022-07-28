package jwt_setup

import (
	"auth_service/internal/config"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func CreateToken(cfg *config.Config, userId string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userId,
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(5 * time.Hour)},
	})
	return token.SignedString([]byte(cfg.Keys.JWTSignKey))
}
