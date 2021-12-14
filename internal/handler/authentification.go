package handler

import (
	"CRUDServer/internal/repository"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// Registration is echo authentication method(POST) for creating user
func (h Handler) Registration(c echo.Context) error {
	hashedPassword, err := hashPassword(c.QueryParam("password"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Registration failed.")
	}
	err = h.rps.CreateAuthUser(repository.AuthForm{
		UserName: c.QueryParam("userName"),
		Email:    c.QueryParam("email"),
		Password: hashedPassword,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error while saving form.")
	}
	return c.String(http.StatusOK, "Successfully.")
}

func (h Handler) Login(c echo.Echo) error {
	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.WithFields(log.Fields{
			"method": "hashPassword()",
			"err":    err,
		}).Info("Authentication info.")
		return "", nil
	}
	return string(hashedPassword), nil
}
