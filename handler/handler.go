package handler

import (
	"CRUDServer/repository"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type IHandler interface {
	SaveUser(c echo.Context) error
	GetUserByID(c echo.Context) error
	DeleteUserByID(c echo.Context) error
	UpdateUserById(c echo.Context) error
}

type Handler struct {
	rps repository.IRepository
}

func (h Handler) SaveUser(c echo.Context) error {

	return nil
}

func (h Handler) GetUserByID(c echo.Context) error {
	return nil
}

func (h Handler) DeleteUserByID(c echo.Context) error {
	userID := c.QueryParam("id")
	err := h.rps.DeleteUser(userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln("Error while deleting."))
	}
	return c.String(http.StatusOK, fmt.Sprintln("Successfully deleted."))
}

func (h Handler) UpdateUserById(c echo.Context) error {
	return nil
}

func NewHandler(rps string) *Handler {
	switch rps {
	case "mongo":
		h := Handler{rps: new(repository.MongoRepository)}
		return &h
	case "postgres":
		h := Handler{rps: new(repository.PostgresRepository)}
		return &h
	}
	return nil
}
