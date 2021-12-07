package main

import (
	"CRUDServer/handler"
	"fmt"

	"github.com/caarlos0/env"
	"github.com/labstack/echo/v4"
)

type config struct {
	CurrentDB string `env:"CURRENTDB, /.env" envDefault:"postgres"`
}

func main() {

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Println(err)
	}

	e := echo.New()
	h := handler.NewHandler(cfg.CurrentDB)
	fmt.Println(cfg.CurrentDB)

	e.POST("users/", h.SaveUser)
	e.PUT("users/", h.UpdateUserById)
	e.DELETE("users/", h.DeleteUserByID)
	e.GET("users/", h.GetUserByID)

	e.Logger.Fatal(e.Start(":8081"))
}
