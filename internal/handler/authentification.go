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

type customClaims struct {
	email    string
	userName string
	jwt.StandardClaims
}

// Registration is echo authentication method(POST) for creating user
func (h Handler) Registration(c echo.Context) error {
	hashedPassword, err := hashPassword(c.QueryParam("password"))
	if err != nil {
		authOperationError(err, "Registration()")
		return c.String(http.StatusInternalServerError, "Registration failed.")
	}

	err = h.rps.CreateAuthUser(repository.RegistrationForm{
		UserName: c.QueryParam("userName"),
		Email:    c.QueryParam("email"),
		Password: hashedPassword,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error while saving form.")
	}

	return c.String(http.StatusOK, "Successfully.")
}

// Authentication check user password and if it correct returns access and refresh tokens
func (h Handler) Authentication(c echo.Context) error {
	authUser, err := h.rps.GetAuthUser(c.QueryParam("email"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Incorrect email")
	}
	hashedLoginPassword, err := hashPassword(c.QueryParam("password"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Login failed, try again")
	}

	if authUser.Password != hashedLoginPassword {
		return c.String(http.StatusBadRequest, "Incorrect password")
	}
	expirationTimeAT := time.Now().Add(10 * time.Minute)
	expirationTimeRT := time.Now().Add(time.Hour * 720)

	atClaims := &customClaims{
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
		return c.String(http.StatusInternalServerError, "Error while creating token")
	}

	rtClaims := &customClaims{
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
		return c.String(http.StatusInternalServerError, "Error while creating token")
	}

	err = h.rps.UpdateAuthUser(authUser.Email, refreshTokenString)
	if err != nil {
		authOperationError(err, "Authentication()")
		return c.String(http.StatusInternalServerError, "Error while creating token")
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

func (h Handler) Authorization(c echo.Context) error {
	accessTokenString := c.QueryParam("token")
	accessToken, err := jwt.ParseWithClaims(accessTokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRETKEY")), nil
	})
	if err != nil {
		authOperationError(err, "Authorization()")
		return c.String(http.StatusInternalServerError, "error while parsing token")
	}
	fmt.Println(accessToken)
	return nil
}

func (h Handler) RefreshToken(c echo.Context) error {
	refreshToken := c.QueryParam("refreshToken")
	if refreshToken == "" {
		return c.String(http.StatusBadRequest, "Refresh failed.")
	}

	return nil
}

func (h Handler) Logout(c echo.Context) error {

	return nil
}

func hashPassword(password string) (string, error) {
	if len(password) == 0 {
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
