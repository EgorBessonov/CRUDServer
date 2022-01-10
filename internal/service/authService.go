// Package service replies server logic
package service

import (
	"CRUDServer/internal/cache"
	"CRUDServer/internal/model"
	"CRUDServer/internal/repository"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	rps repository.Repository
	orderCache *cache.OrderCache
}

func NewService(_rps repository.Repository, _orderCache *cache.OrderCache) *Service{
	return &Service{ rps:_rps, orderCache: _orderCache}
}
	
const(
	accessTokenExTime = 15
	refreshTokenExTime = 720
)

// CustomClaims struct represent user information in tokens
type CustomClaims struct {
	email    string
	userName string
	jwt.StandardClaims
}

func(s Service) Registration( ctx context.Context, authUser *model.AuthUser) error {
	hPassword, err := hashPassword(authUser.Password, "Registration()")
	if err != nil {
		return err
	}
	authUser.Password = hPassword
	err = s.rps.CreateAuthUser(ctx, authUser)
	if err != nil {
		return err
	}
	return nil
}

func(s Service) RefreshToken(ctx context.Context, refreshTokenString string) (string, string, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRETKEY")), nil
	}
	refreshToken, err := jwt.Parse(refreshTokenString, keyFunc)
	if err != nil {
		servicesOperationError(err, "RefreshToken()")
		return "", "", err
	}
	if !refreshToken.Valid {
		servicesOperationError(errors.New("invalid token"), "RefreshToken()")
		return "", "", errors.New("invalid token")
	}
	claims := refreshToken.Claims.(jwt.MapClaims)
	userUUID := claims["jti"]
	if userUUID == "" {
		servicesOperationError(errors.New("error while parsing token claims"), "refreshToken()")
		return "", "", errors.New("error while parsing token claims")
	}
	authUser, err := s.rps.GetAuthUserByID(ctx, userUUID.(string))
	if err != nil {
		return "", "", errors.New("error while parsing token")
	}
	if refreshTokenString != authUser.RefreshToken {
		servicesOperationError(errors.New("invalid refresh token"), "RefreshToken()")
		return "", "", errors.New("invalid refresh token")
	}
	return createTokenPair(s.rps, ctx, &authUser)
}

func(s Service) Authentication(ctx context.Context, email, password string)(string, string, error){
	hashPassword, err := hashPassword(password, "Authentication()")
	if err != nil{
		return "", "", err
	}
	authForm, err := s.rps.GetAuthUser(ctx, email)
	if err != nil{
		return "", "", err
	}
	if authForm.Password != hashPassword{
		servicesOperationError(errors.New("invalid password"), "Authentication()")
		return "", "", errors.New("invalid password")
	}
	return createTokenPair(s.rps, ctx, &authForm)
}

func(s Service) UpdateAuthUser(ctx context.Context, email string, refreshToken string) error {
	return s.rps.UpdateAuthUser(ctx, email, refreshToken)
}

func createTokenPair(rps repository.Repository, ctx context.Context, authUser *model.AuthUser) (string, string, error) {
	expirationTimeAT := time.Now().Add(accessTokenExTime * time.Minute)
	expirationTimeRT := time.Now().Add(time.Hour * refreshTokenExTime)

	atClaims := &CustomClaims{
		userName: authUser.UserName,
		email:    authUser.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTimeAT.Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("SECRETKEY")))
	if err != nil {
		servicesOperationError(err, "Authentication()")
		return "", "", fmt.Errorf("error while creating token")
	}

	rtClaims := &CustomClaims{
		userName: authUser.UserName,
		email:    authUser.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTimeRT.Unix(),
			Id:        authUser.UserUUID,
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("SECRETKEY")))
	if err != nil {
		servicesOperationError(err, "Authentication()")
		return "", "", fmt.Errorf("error while creating token")
	}

	err = rps.UpdateAuthUser(ctx, authUser.Email, refreshTokenString)
	if err != nil {
		servicesOperationError(err, "Authentication()")
		return "", "", fmt.Errorf("error while creating token")
	}
	return accessTokenString, refreshTokenString, nil
}

func hashPassword(password, method string) (string, error) {
	if password == "" {
		servicesOperationError(errors.New("no input supplied"), method)
		return "", errors.New("no input supplied")
	}
	h := sha256.New()
	h.Write([]byte(password))
	hashedPassword := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return hashedPassword, nil
}

func servicesOperationError(err error, method string) {
	log.WithFields(log.Fields{
		"method": method,
		"err":    err,
		"status": "operation failed.",
	}).Info("Services info.")
}
