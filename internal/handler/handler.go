// Package handler replies handler for echo server
package handler

import (
	"fmt"
	configs "github.com/EgorBessonov/CRUDServer/internal/config"
	"github.com/EgorBessonov/CRUDServer/internal/model"
	"github.com/EgorBessonov/CRUDServer/internal/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Handler type replies for handling echo server requests
type Handler struct {
	s   *service.Service
	cfg *configs.Config
}

// NewHandler function create handler for working with
// postgres or mongo database and initialize connection with this db
func NewHandler(s *service.Service, cfg *configs.Config) *Handler {
	return &Handler{s: s, cfg: cfg}
}

// SaveOrder godoc
// @Summary SaveOrder
// @Description SaveOrder is echo handler(POST) which return orderID
// @Tags orders
// @Accept json
// @Produce json
// @Param order body model.Order true "order instance"
// @Success 200 {string} string
// @Failure 500 {object} echo.HTTPError
// @Router /saveOrder [post]
// @Security ApiKeyAuth
func (h *Handler) SaveOrder(c echo.Context) error {
	order := model.Order{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &order); err != nil {
		log.Error(fmt.Errorf("handler: can't save order - %w", err))
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error while saving"))
	}
	orderID, err := h.s.Save(c.Request().Context(), &order)
	if err != nil {
		log.Error(fmt.Errorf("handler: can't save order - %w", err))
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error while saving"))
	}
	return c.JSONBlob(
		http.StatusOK,
		[]byte(
			fmt.Sprintf(`{"orderID" : %v}`, orderID)),
	)
}

// GetOrderByID godoc
// @Summary GetOrderByID
// @Description GetOrderByID is echo handler(GET) which returns json structure of User object
// @Tags orders
// @Accept json
// @Produce json
// @Param orderID query string true "orderID"
// @Success 200 {string} string
// @Failure 500 {object} echo.HTTPError
// @Router /getOrder{orderID} [get]
// @Security ApiKeyAuth
func (h *Handler) GetOrderByID(c echo.Context) error {
	orderID := c.QueryParam("orderID")
	order, err := h.s.Get(c.Request().Context(), orderID)
	if err != nil {
		log.Error(fmt.Errorf("handler: can't get order - %w", err))
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintln("get operation failed"))
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

// DeleteOrderByID godoc
// @Summary DeleteOrderByID
// @Description DeleteOrderByID is echo handler(DELETE)
// @Tags orders
// @Accept json
// @Produce json
// @Param orderID query string true "orderID"
// @Success 200 {string} string
// @Failure 500 {object} echo.HTTPError
// @Router /deleteOrder{orderID} [delete]
// @Security ApiKeyAuth
func (h *Handler) DeleteOrderByID(c echo.Context) error {
	orderID := c.QueryParam("orderID")
	err := h.s.Delete(c.Request().Context(), orderID)
	if err != nil {
		log.Error(fmt.Errorf("handler: can't delete order - %w", err))
		return echo.NewHTTPError(http.StatusInternalServerError, "error while deleting")
	}
	return c.String(http.StatusOK, fmt.Sprintln("successfully deleted."))
}

// UpdateOrderByID godoc
// @Summary UpdateOrderByID
// @Description UpdateOrderByID is echo handler(PUT) which return updating status
// @Tags orders
// @Accept json
// @Produce json
// @Param order body model.Order true "order instance"
// @Success 200 {string} string
// @Failure 500 {object} echo.HTTPError
// @Router /updateOrder [put]
// @Security ApiKeyAuth
func (h *Handler) UpdateOrderByID(c echo.Context) error {
	order := model.Order{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &order); err != nil {
		log.Error("handler: can't update order - error while parsing")
		return echo.NewHTTPError(http.StatusInternalServerError, "error while parsing json")
	}
	err := h.s.Update(c.Request().Context(), &order)
	if err != nil {
		log.Errorf("handler: can't update order - %e", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error while updating user")
	}
	return c.String(http.StatusOK, fmt.Sprintln("successfully updated."))
}

// UploadImage godoc
// @Summary UploadImage
// @Description UploadImage is echo handler(POST) for uploading user images from server
// @Tags images
// @Accept json
// @Produce json
// @Param image body string true "image"
// @Success 200 {string} string
// @Failure 500 {object} echo.HTTPError
// @Router /uploadImage [post]
func (h *Handler) UploadImage(c echo.Context) error {
	imageFile, err := c.FormFile("image")
	if err != nil {
		log.Errorf("handler: can't upload image - %e", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "image uploading failed.")
	}
	err = h.s.UploadImage(imageFile)
	if err != nil {
		log.Errorf("handler: can't upload image - %e", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "image uploading failed.")
	}
	return c.String(http.StatusOK, fmt.Sprintln("successfully uploaded."))
}

// DownloadImage godoc
// @Summary DownloadImage
// @Description DownloadImage is echo handler(GET) for downloading user images
// @Tags images
// @Accept json
// @Produce json
// @Param imageName query string true "image name"
// @Success 200 {string} string
// @Failure 500 {object} echo.HTTPError
// @Router /downloadImage{imageName} [get]
func (h *Handler) DownloadImage(c echo.Context) error {
	imageName := c.QueryParam("imageName")
	if imageName == "" {
		log.Errorf("handler: can't download image - empty value")
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintln("invalid image name."))
	}
	return c.File("images/" + imageName)
}
