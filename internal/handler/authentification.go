package handler

import (
	"CRUDServer/internal/repository"
	"crypto/sha256"
	"net/http"
	"github.com/labstack/echo/v4"
)

// LoginForm struct represents user login information
type LoginForm struct {
	Email    string `json: "email" bson: "email"`
	Password []byte `json: "password" bson: "password"`
}

// Login is echo authentification method(POST) for creating user
func (h Handler) Registration(c echo.Context) error{
	passwordHash := sha256.Sum256([]byte(c.QueryParam("password")))
	err := h.rps.CreateAuthUser(LoginForm{
		Email: c.QueryParam("email"),
		Password: passwordHash[:],
	})
	if err != nil{
		return c.String(http.StatusInternalServerError, "Registration failed.")
	}
	return c.String(http.StatusOK, "Succesfully registrated.")
}

func (h Handler) Login(c echo.Echo)error{

	return nil
}