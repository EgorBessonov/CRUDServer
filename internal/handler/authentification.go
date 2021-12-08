package handler

import (
	"github.com/golang-jwt/jwt"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin string `json:"admin"`
	jwt.StandardClaims
}
