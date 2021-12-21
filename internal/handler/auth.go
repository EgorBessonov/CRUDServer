package handler

import (
	"CRUDServer/internal/repository"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

// CustomClaims struct represent user information in tokens
type CustomClaims struct {
	email    string
	userName string
	jwt.StandardClaims
}

// Registration method is echo authentication method(POST) for creating user
func (h Handler) Registration(c echo.Context) error {
	hashedPassword, err := hashPassword(c.QueryParam("password"))
	if err != nil {
		authOperationError(err, "Registration()")
		return c.String(http.StatusInternalServerError, "registration failed.")
	}

	err = h.rps.CreateAuthUser(repository.RegistrationForm{
		UserName: c.QueryParam("userName"),
		Email:    c.QueryParam("email"),
		Password: hashedPassword,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "error while saving form.")
	}

	return c.String(http.StatusOK, "successfully.")
}

// Authentication method checks user password and if it ok returns access and refresh tokens
func (h Handler) Authentication(c echo.Context) error {
	authUser, err := h.rps.GetAuthUser(c.QueryParam("email"))
	if err != nil {
		return c.String(http.StatusBadRequest, "incorrect email")
	}
	hashedLoginPassword, err := hashPassword(c.QueryParam("password"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "login failed, try again")
	}

	if authUser.Password != hashedLoginPassword {
		return c.String(http.StatusBadRequest, "incorrect password")
	}
	accessTokenString, refreshTokenString, err := h.createTokenPair(authUser)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln(err))
	}
	return c.JSONBlob(
		http.StatusOK,
		[]byte(
			fmt.Sprintf(`{
			"accessToken" : %v,
			"refreshToken" : %v}`, accessTokenString, refreshTokenString),
		),
	)
}

// RefreshToken method checks refresh token for validity and if it ok returns new token pair
func (h Handler) RefreshToken(c echo.Context) error {
	refreshTokenString := c.QueryParam("refreshToken")
	if refreshTokenString == "" {
		return c.String(http.StatusBadRequest, "empty refresh token.")
	}
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRETKEY")), nil
	}
	refreshToken, err := jwt.Parse(refreshTokenString, keyFunc)
	if err != nil {
		authOperationError(err, "RefreshToken()")
		return c.String(http.StatusInternalServerError, "error while parsing token.")
	}
	if !refreshToken.Valid {
		return c.String(http.StatusNonAuthoritativeInfo, "invalid token.")
	}
	claims := refreshToken.Claims.(jwt.MapClaims)
	userUUID := claims["jti"]
	if userUUID == "" {
		return c.String(http.StatusInternalServerError, "error while parsing token.")
	}
	authUser, err := h.rps.GetAuthUserByID(userUUID.(string))
	if err != nil {
		return c.String(http.StatusInternalServerError, "error while parsing token")
	}
	if refreshTokenString != authUser.RefreshToken {
		return c.String(http.StatusBadRequest, "invalid refresh token")
	}
	newAccessTokenString, newRefreshTokenString, err := h.createTokenPair(authUser)
	if err != nil {
		return c.String(http.StatusInternalServerError, "error while creating tokens")
	}
	return c.JSONBlob(
		http.StatusOK,
		[]byte(
			fmt.Sprintf(`{
			"accessToken" : %v,
			"refreshToken" : %v}`, newAccessTokenString, newRefreshTokenString),
		),
	)
}

// Logout method delete user refresh token from database
func (h Handler) Logout(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" {
		return c.String(http.StatusBadRequest, "Empty value")
	}
	err := h.rps.UpdateAuthUser(email, "")
	if err != nil {
		return c.String(http.StatusInternalServerError, "logout error.")
	}
	return c.String(http.StatusOK, "logout successfully")
}

func (h Handler) createTokenPair(authUser repository.RegistrationForm) (string, string, error) {
	expirationTimeAT := time.Now().Add(15 * time.Minute)
	expirationTimeRT := time.Now().Add(time.Hour * 720)

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
		authOperationError(err, "Authentication()")
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
		authOperationError(err, "Authentication()")
		return "", "", fmt.Errorf("error while creating token")
	}

	err = h.rps.UpdateAuthUser(authUser.Email, refreshTokenString)
	if err != nil {
		authOperationError(err, "Authentication()")
		return "", "", fmt.Errorf("error while creating token")
	}
	return accessTokenString, refreshTokenString, nil
}

func hashPassword(password string) (string, error) {
	if password == "" {
		log.WithFields(log.Fields{
			"method": "hashPassword()",
			"err":    errors.New("no input supplied"),
		}).Info("Authentication info.")
		return "", errors.New("no input supplied")
	}
	h := sha256.New()
	h.Write([]byte(password))
	hashedPassword := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return hashedPassword, nil
}

func authOperationError(err error, method string) {
	log.WithFields(log.Fields{
		"method": method,
		"err":    err,
		"status": "Operation failed.",
	}).Info("Authentication info.")
}
