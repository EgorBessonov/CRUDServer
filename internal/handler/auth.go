package handler

import (
	"CRUDServer/internal/model"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// Registration method is echo authentication method(POST) for creating user
func (h Handler) Registration(c echo.Context) error {
	authUser := model.AuthUser{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &authUser); err != nil {
		handlerOperationError(errors.New("error while parsing json"), "Registration()")
		return c.String(http.StatusInternalServerError, "error while parsing json")
	}
	err := h.s.Registration(c.Request().Context(), &authUser)
	if err != nil {
		return c.String(http.StatusInternalServerError, "error while saving form.")
	}
	return c.String(http.StatusOK, "successfully.")
}

// Authentication method checks user password and if it ok returns access and refresh tokens
func (h Handler) Authentication(c echo.Context) error {
	authUser := struct {
		Email	 string
		Password string
	}{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &authUser); err != nil {
		handlerOperationError(errors.New("error while parsing json"), "Authentication()")
		return c.String(http.StatusInternalServerError, "error while parsing json")
	}
	accessTokenString, refreshTokenString, err := h.s.Authentication(c.Request().Context(), authUser.Email, authUser.Password)
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
		authOperationError(errors.New("empty refresh token"), "RefreshToken()")
		return c.String(http.StatusBadRequest, "empty refresh token.")
	}
	newAccessTokenString, newRefreshTokenString, err := h.s.RefreshToken(c.Request().Context(), refreshTokenString)
	if err != nil {
		authOperationError(errors.New("error while creating token"), "RefreshToken()")
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
		authOperationError(errors.New("empty email value"), "Logout()")
		return c.String(http.StatusBadRequest, "Empty value")
	}
	err := h.s.UpdateAuthUser(c.Request().Context(), email, "")
	if err != nil {
		authOperationError(errors.New("logout error"), "Logout()")
		return c.String(http.StatusInternalServerError, "logout error.")
	}
	log.WithFields(log.Fields{
		"status": "Successfully",
		"method": "Logout",
	})
	return c.String(http.StatusOK, "logout successfully")
}

func authOperationError(err error, method string) {
	log.WithFields(log.Fields{
		"method": method,
		"err":    err,
		"status": "Operation failed.",
	}).Info("Authentication info.")
}
