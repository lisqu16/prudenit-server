package user

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/lisqu16/prudenit-server/config"
)

func SignToken(claims *jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(config.TokenSecret))
	return tokenString
}