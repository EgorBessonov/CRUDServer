// Package handler replies handlers for echo server
package handler

import (
	"CRUDServer/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handler type replies for handling echo server requests
type Handler struct {
	rps repository.IRepository
}

// SaveUser is echo handler which return creation status and UserId
func (h Handler) SaveUser(c echo.Context) error {
	userAge, err := strconv.Atoi(c.QueryParam("userAge"))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln("Error while converting data."))
	}

	isAdult, err := strconv.ParseBool(c.QueryParam("isAdult"))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln("Error while converting data."))
	}
	user := repository.User{
		UserID:   uuid.New().Version().String(),
		UserName: c.QueryParam("userName"),
		UserAge:  userAge,
		IsAdult:  isAdult,
	}
	err = h.rps.CreateUser(user)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln("Error while adding User to db."))
	}
	return c.String(http.StatusOK, fmt.Sprintln("Successfully added."))
}

// GetUserByID is echo handler which returns json structure of User object
func (h Handler) GetUserByID(c echo.Context) error {
	userID := c.QueryParam("userId")
	user, err := h.rps.ReadUser(userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln("Error while reading."))
	}
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

// DeleteUserByID is echo handler which return deletion status
func (h Handler) DeleteUserByID(c echo.Context) error {
	userID := c.QueryParam("userId")
	err := h.rps.DeleteUser(userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln("Error while deleting."))
	}
	return c.String(http.StatusOK, fmt.Sprintln("Successfully deleted."))
}

// UpdateUserByID is echo handler which return updating status
func (h Handler) UpdateUserByID(c echo.Context) error {
	userAge, err := strconv.Atoi(c.QueryParam("userAge"))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln("Error while converting data."))
	}

	isAdult, err := strconv.ParseBool(c.QueryParam("isAdult"))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln("Error while converting data."))
	}

	user := repository.User{
		UserID:   c.QueryParam("userId"),
		UserName: c.QueryParam("userName"),
		UserAge:  userAge,
		IsAdult:  isAdult,
	}
	err = h.rps.UpdateUser(user)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln("Error while updating user"))
	}
	return c.String(http.StatusOK, fmt.Sprintln("Successfully updated."))
}

// NewHandler function create handler for working with
// postgres or mongo database
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
