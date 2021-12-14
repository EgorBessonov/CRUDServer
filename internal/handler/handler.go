// Package handler replies handlers for echo server
package handler

import (
	"CRUDServer/internal/repository"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
)

// Handler type replies for handling echo server requests
type Handler struct {
	rps repository.IRepository
}

// NewHandler function create handler for working with
// postgres or mongo database and initialize connection with this db
func NewHandler(cfg repository.Config) *Handler {
	switch cfg.CurrentDB {
	case "mongo":
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongodbURL))
		if err != nil {
			log.WithFields(log.Fields{
				"status": "Connection to mongo database failed.",
				"err":    err,
			}).Info("Mongo repository info")
		} else {
			log.WithFields(log.Fields{
				"status": "Successfully connected to mongo database.",
			}).Info("Mongo repository info.")
		}
		h := Handler{rps: repository.MongoRepository{
			DBconn: client,
		}}
		return &h
	case "postgres":
		conn, err := pgxpool.Connect(context.Background(), cfg.PostgresdbURL)
		if err != nil {
			log.WithFields(log.Fields{
				"status": "Connection to postgres database failed.",
				"err":    err,
			}).Info("Postgres repository info.")
		} else {
			log.WithFields(log.Fields{
				"status": "Successfully connected to postgres database.",
			}).Info("Postgres repository info.")
		}
		h := Handler{repository.PostgresRepository{
			DBconn: conn,
		}}
		return &h
	}
	return nil
}

// SaveUser is echo handler(POST) which return creation status and UserId
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

// GetUserByID is echo handler(GET) which returns json structure of User object
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

// DeleteUserByID is echo handler(DELETE) which return deletion status
func (h Handler) DeleteUserByID(c echo.Context) error {
	userID := c.QueryParam("userId")
	fmt.Println(userID)
	err := h.rps.DeleteUser(userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln("Error while deleting."))
	}
	return c.String(http.StatusOK, fmt.Sprintln("Successfully deleted."))
}

// UpdateUserByID is echo handler(PUT) which return updating status
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

// AddImage is echo handler(POST) for saving user images
func (h Handler) AddImage(c echo.Context) error {
	return nil
}

// GetImageByUserID is echo handler(GET) for getting user images
func (h Handler) GetImageByUserID(c echo.Context) error {
	return nil
}
