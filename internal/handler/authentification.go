package handler

import (
	"CRUDServer/internal/repository"
	"crypto/sha256"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Registration is echo authentication method(POST) for creating user
func (h Handler) Registration(c echo.Context) error {
	passwordHash := sha256.Sum256([]byte(c.QueryParam("password")))
	err := h.rps.CreateAuthUser(repository.LoginForm{
		Email:    c.QueryParam("email"),
		Password: passwordHash[:],
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Registration failed.")
	}
	return c.String(http.StatusOK, "Successfully.")
}

func (h Handler) Login(c echo.Echo) error {

	return nil
}
