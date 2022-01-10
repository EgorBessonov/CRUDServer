// Package handler replies handler for echo server
package handler

import (
	configs "CRUDServer/internal/config"
	"CRUDServer/internal/model"
	"CRUDServer/internal/service"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// Handler type replies for handling echo server requests
type Handler struct {
	s     *service.Service
	cfg   *configs.Config
}

// NewHandler function create handler for working with
// postgres or mongo database and initialize connection with this db
func NewHandler(_s *service.Service, _cfg *configs.Config) *Handler {
	return &Handler{s: _s, cfg: _cfg}
}

// SaveOrder is echo handler(POST) which return creation status
func (h Handler) SaveOrder(c echo.Context) error {
	order := model.Order{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &order); err != nil {
		handlerOperationError(errors.New("error while parsing json"), "Authentication()")
		return c.String(http.StatusInternalServerError, "error while parsing json")
	}
	err := h.s.Save(c.Request().Context(), order)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln("Error while adding User to db."))
	}
	return c.String(http.StatusOK, fmt.Sprintln("successfully added."))
}

// GetOrderByID is echo handler(GET) which returns json structure of User object
func (h Handler) GetOrderByID(c echo.Context) error {
	orderID := c.QueryParam("orderID")
	order, err := h.s.Get(c.Request().Context(), h.cfg, orderID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln("error while reading."))
	}
	return c.JSONBlob(
		http.StatusOK,
		[]byte(
			fmt.Sprintf(`{
					"orderName" : %v,
					"orderCost" : %v,
					"isDelivered" : %v}`, order.OrderName, order.OrderCost, order.IsDelivered),
		),
	)
}

// DeleteOrderByID is echo handler(DELETE) which return deletion status
func (h Handler) DeleteOrderByID(c echo.Context) error {
	orderID := c.QueryParam("orderID")
	fmt.Println(orderID)
	err := h.s.Delete(c.Request().Context(), orderID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln("error while deleting."))
	}
	return c.String(http.StatusOK, fmt.Sprintln("successfully updated."))
}

// UpdateOrderByID is echo handler(PUT) which return updating status
func (h Handler) UpdateOrderByID(c echo.Context) error {
	order := model.Order{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &order); err != nil {
		handlerOperationError(errors.New("error while parsing json"), "Registration()")
		return c.String(http.StatusInternalServerError, "error while parsing json")
	}
	err := h.s.Update(c.Request().Context(), h.cfg,order)
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
