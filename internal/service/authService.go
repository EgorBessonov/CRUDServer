// Package service replies server logic
package service

import (
	"CRUDServer/internal/models"
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

func Registration( ctx context.Context, rps repository.Repository, authUser *models.AuthUser) error {
	hPassword, err := hashPassword(authUser.Password, "Registration()")
	if err != nil {
		return err
	}
	authUser.Password = hPassword
	err = rps.CreateAuthUser(ctx, authUser)
	if err != nil {
		return err
	}
	return nil
}

func RefreshToken(ctx context.Context, rps repository.Repository, refreshTokenString string) (string, string, error) {
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
	authUser, err := rps.GetAuthUserByID(ctx, userUUID.(string))
	if err != nil {
		return "", "", errors.New("error while parsing token")
	}
	if refreshTokenString != authUser.RefreshToken {
		servicesOperationError(errors.New("invalid refresh token"), "RefreshToken()")
		return "", "", errors.New("invalid refresh token")
	}
	return createTokenPair(rps, ctx, &authUser)
}

func Authentication(ctx context.Context, rps repository.Repository, form *models.AuthUser)(string, string, error){
	hashPassword, err := hashPassword(form.Password, "Authentication()")
	if err != nil{
		return "", "", err
	}
	authForm, err := rps.GetAuthUser(ctx, form.Email)
	if err != nil{
		return "", "", err
	}
	if authForm.Password != hashPassword{
		servicesOperationError(errors.New("invalid password"), "Authentication()")
		return "", "", errors.New("invalid password")
	}
	return  createTokenPair(rps, ctx, &authForm)
}

func createTokenPair(rps repository.Repository,ctx context.Context, authUser *models.AuthUser) (string, string, error) {
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
