package handler

import (
	"CRUDServer/repository"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	rps repository.IRepository
}

func (h Handler) SaveUser(c echo.Context) error {

	return nil
}

func (h Handler) GetUserByID(c echo.Context) error {
	userID := c.QueryParam("id")
	user, err := h.rps.ReadUser(userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln("Error while reading."))
	}
	//return c.String(http.StatusOK, fmt.Sprintf("userName: %v\nuserAge: %v\nisAdult: %v\n", user.UserName, user.UserAge, user.IsAdult))
	return c.JSONBlob(
		http.StatusOK,
		[]byte(
			fmt.Sprintf(`{
					"userName" : %v,
					"userAge" : %v,
					"isAdult" : %v}`, user.UserName, user.UserAge, userID),
		),
	)
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
