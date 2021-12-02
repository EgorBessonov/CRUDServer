package main

import (
	"CRUDServer/handler"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()
	h := handler.NewHandler("postgres")

	e.POST("users/", h.SaveUser)
	e.PUT("users/", h.UpdateUserById)
	e.DELETE("users/", h.DeleteUserByID)
	e.GET("users/", h.GetUserByID)

	e.Logger.Fatal(e.Start(":8080"))
}
