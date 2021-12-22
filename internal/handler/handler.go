// Package handler replies handlers for echo server
package handler

import (
	"CRUDServer/internal/repository"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strconv"
)

// Handler type replies for handling echo server requests
type Handler struct {
	rps repository.Repository
}

// NewHandler function create handler for working with
// postgres or mongo database and initialize connection with this db
func NewHandler(_rps repository.Repository) *Handler {
	h := Handler{rps: _rps}
	return &h
}

// SaveUser is echo handler(POST) which return creation status and UserId
func (h Handler) SaveUser(c echo.Context) error {
	user := repository.User{}
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	fmt.Println(user)
	err = h.rps.CreateUser(user, c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln("Error while adding User to db."))
	}
	return c.String(http.StatusOK, fmt.Sprintln("successfully added."))
}

// GetUserByID is echo handler(GET) which returns json structure of User object
func (h Handler) GetUserByID(c echo.Context) error {
	userID := c.QueryParam("userId")
	user, err := h.rps.ReadUser(userID, c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln("error while reading."))
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

// DeleteUserByID is echo handler(DELETE) which return deletion status
func (h Handler) DeleteUserByID(c echo.Context) error {
	userID := c.QueryParam("userId")
	fmt.Println(userID)
	err := h.rps.DeleteUser(userID, c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln("error while deleting."))
	}
	return c.String(http.StatusOK, fmt.Sprintln("successfully deleted."))
}

// UpdateUserByID is echo handler(PUT) which return updating status
func (h Handler) UpdateUserByID(c echo.Context) error {
	userAge, err := strconv.Atoi(c.QueryParam("userAge"))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln("error while converting data."))
	}

	isAdult, err := strconv.ParseBool(c.QueryParam("isAdult"))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln("error while converting data."))
	}

	user := repository.User{
		UserID:   c.QueryParam("userId"),
		UserName: c.QueryParam("userName"),
		UserAge:  userAge,
		IsAdult:  isAdult,
	}
	err = h.rps.UpdateUser(user, c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln("error while updating user"))
	}
	return c.String(http.StatusOK, fmt.Sprintln("successfully updated."))
}

// UploadImage is echo handler(POST) for uploading user images from server
func (h Handler) UploadImage(c echo.Context) error {
	imageFile, err := c.FormFile("image")
	if err != nil {
		handlerOperationError(err, "UploadImage()")
		return c.String(http.StatusInternalServerError, "operation failed.")
	}
	imageSrc, err := imageFile.Open()
	if err != nil {
		handlerOperationError(err, "UploadImage()")
		return c.String(http.StatusInternalServerError, "operation failed.")
	}
	defer func() {
		err = imageSrc.Close()
		if err != nil {
			handlerOperationError(err, "UploadImage()")
		}
	}()
	dst, err := os.Create(imageFile.Filename)
	if err != nil {
		handlerOperationError(err, "UploadImage()")
		return c.String(http.StatusInternalServerError, "operation failed.")
	}
	if _, err = io.Copy(dst, imageSrc); err != nil {
		handlerOperationError(err, "UploadImage()")
		return c.String(http.StatusInternalServerError, "operation failed.")
	}
	return c.String(http.StatusOK, fmt.Sprintln("successfully uploaded."))
}

// DownloadImage is echo handler(GET) for downloading user images
func (h Handler) DownloadImage(c echo.Context) error {
	imageName := c.QueryParam("imageName")
	if imageName == "" {
		return c.String(http.StatusBadRequest, fmt.Sprintln("invalid image name."))
	}
	return c.File(imageName)
}

func handlerOperationError(err error, method string) {
	log.WithFields(log.Fields{
		"status": "Operation Failed.",
		"method": method,
		"error":  err,
	}).Info("Handler info.")
}
