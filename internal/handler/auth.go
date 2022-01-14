package handler

import (
	"CRUDServer/internal/model"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// Registration godoc
// @Summary Registration
// @Description Registration method is echo authentication method(POST) for creating user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body model.AuthUser true "auth user instance"
// @Success 200 {object} string
// @Failure 500 {object} echo.HTTPError
// @Router registration/ [post]
func (h *Handler) Registration(c echo.Context) error {
	authUser := model.AuthUser{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &authUser); err != nil {
		log.Errorf("handler: registration failed - %e", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error while parsing json")
	}
	err := h.s.Registration(c.Request().Context(), &authUser)
	if err != nil {
		log.Errorf("handler: registration failed - %e", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error while saving form.")
	}
	return c.String(http.StatusOK, "successfully.")
}

// Authentication godoc
// @Summary Authentication
// @Description Authentication method checks user password and if it ok returns access and refresh tokens
// @Tags orders
// @Accept json
// @Produce json
// @Param input body model.AuthUser true "email & password"
// @Success 200 {object} string
// @Failure 500 {object} HTTPError
// @Router authentication/ [post]
func (h *Handler) Authentication(c echo.Context) error {
	authUser := struct {
		Email    string
		Password string
	}{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &authUser); err != nil {
		log.Errorf("handler: authentication failed - %e", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error while parsing json")
	}
	accessTokenString, refreshTokenString, err := h.s.Authentication(c.Request().Context(), authUser.Email, authUser.Password)
	if err != nil {
		log.Errorf("handler: authentication failed - %e", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintln(err))
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

// RefreshToken godoc
// @Summary Logout
// @Description RefreshToken method checks refresh token for validity and if it ok returns new token pair
// @Tags orders
// @Accept json
// @Produce json
// @Param q query string true "refresh token"
// @Success 200 {object} string
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router logout/ [get]
func (h *Handler) RefreshToken(c echo.Context) error {
	refreshTokenString := c.QueryParam("refreshToken")
	if refreshTokenString == "" {
		log.Errorf("handler: token refresh  failed - empty value")
		return echo.NewHTTPError(http.StatusBadRequest, "empty refresh token.")
	}
	newAccessTokenString, newRefreshTokenString, err := h.s.RefreshToken(c.Request().Context(), refreshTokenString)
	if err != nil {
		log.Errorf("handler: token refresh failed - %e", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error while creating tokens")
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

// Logout godoc
// @Summary Logout
// @Description Logout method delete user refresh token from database
// @Tags orders
// @Accept json
// @Produce json
// @Param q query string true "email"
// @Success 200 {object} string
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router logout/ [post]
func (h *Handler) Logout(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" {
		log.Error("handler: logout failed - empty value")
		return echo.NewHTTPError(http.StatusBadRequest, "empty value")
	}
	err := h.s.UpdateAuthUser(c.Request().Context(), email, "")
	if err != nil {
		log.Errorf("handler: logout failed - %e", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "logout error.")
	}
	return c.String(http.StatusOK, "logout successfully")
}
