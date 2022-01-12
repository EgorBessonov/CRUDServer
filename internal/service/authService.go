// Package service replies server logic
package service

import (
	"CRUDServer/internal/cache"
	"CRUDServer/internal/model"
	"CRUDServer/internal/repository"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type Service struct {
	rps        repository.Repository
	orderCache *cache.OrderCache
}

func NewService(_rps repository.Repository, _orderCache *cache.OrderCache) *Service {
	return &Service{rps: _rps, orderCache: _orderCache}
}

const (
	accessTokenExTime  = 15
	refreshTokenExTime = 720
)

// CustomClaims struct represent user information in tokens
type CustomClaims struct {
	email    string
	userName string
	jwt.StandardClaims
}

func (s Service) Registration(ctx context.Context, authUser *model.AuthUser) error {
	hPassword, err := hashPassword(authUser.Password)
	if err != nil {
		return err
	}
	authUser.Password = hPassword
	err = s.rps.SaveAuthUser(ctx, authUser)
	if err != nil {
		return fmt.Errorf("service: registration failed - %w", err)
	}
	return nil
}

func (s Service) RefreshToken(ctx context.Context, refreshTokenString string) (string, string, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRETKEY")), nil
	}
	refreshToken, err := jwt.Parse(refreshTokenString, keyFunc)
	if err != nil {
		return "", "", fmt.Errorf("service: can't parse refresh token - %w", err)
	}
	if !refreshToken.Valid {
		return "", "", fmt.Errorf("service: expired refresh token")
	}
	claims := refreshToken.Claims.(jwt.MapClaims)
	userUUID := claims["jti"]
	if userUUID == "" {
		return "", "", fmt.Errorf("service: error while parsing claims")
	}
	authUser, err := s.rps.GetAuthUserByID(ctx, userUUID.(string))
	if err != nil {
		return "", "", fmt.Errorf("service: token refresh failed - %w", err)
	}
	if refreshTokenString != authUser.RefreshToken {
		return "", "", fmt.Errorf("service: invalid refresh token")
	}
	return createTokenPair(s.rps, ctx, authUser)
}

func (s Service) Authentication(ctx context.Context, email, password string) (string, string, error) {
	hashPassword, err := hashPassword(password)
	if err != nil {
		return "", "", err
	}
	authForm, err := s.rps.GetAuthUser(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("service: authentication failed - %w", err)
	}
	if authForm.Password != hashPassword {
		return "", "", fmt.Errorf("service: invalid password")
	}
	return createTokenPair(s.rps, ctx, authForm)
}

func (s Service) UpdateAuthUser(ctx context.Context, email string, refreshToken string) error {
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
		return "", "", fmt.Errorf("service: can't generate access token - %w", err)
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
		return "", "", fmt.Errorf("service: can't generate refresh token - %w", err)
	}

	err = rps.UpdateAuthUser(ctx, authUser.Email, refreshTokenString)
	if err != nil {
		return "", "", fmt.Errorf("service: can't set refresh token - %w", err)
	}
	return accessTokenString, refreshTokenString, nil
}

func hashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("service: zero password value")
	}
	h := sha256.New()
	h.Write([]byte(password))
	hashedPassword := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return hashedPassword, nil
}
