package handler

import (
	"CRUDServer/repository"

	"github.com/labstack/echo/v4"
)

type IHandler interface {
	SaveUser(c echo.Context)error
	GetUserByID(c echo.Context)error
	DeleteUserByID(c echo.Context)error
	UpdateUserById(c echo.Context)error
}

type Handler struct{
	rps repository.IRepository
}

func(h Handler) SaveUser(c echo.Context)error{


	return nil
}

func(h Handler) GetUserByID(c echo.Context)error{
	return nil
}

func(h Handler) DeleteUserByID(c echo.Context)error{
	return nil
}

func(h Handler) UpdateUserById(c echo.Context)error{
	return nil
}

func NewHandler(rps string) *Handler{
	switch rps{
	case "mongo":
		h := Handler{ rps: new(repository.MongoRepository)}
		return &h
	case "postgres":
		h := Handler{ rps: new(repository.PostgreRepository)}
		return &h
	}
	return nil
}