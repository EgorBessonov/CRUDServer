package main

import (
	"CRUDServer/internal/handler"
	"CRUDServer/internal/repository"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading .env file.")
	}
	cfg := repository.Config{
		CurrentDB:     os.Getenv("CURRENTDB"),
		PostgresdbUrl: os.Getenv("POSTGRESDB_URL"),
		MongodbUrl:    os.Getenv("MONGODB_URL")}

	if err := env.Parse(&cfg); err != nil {
		fmt.Println(err)
	}
	e := echo.New()

	h := handler.NewHandler(cfg)

	e.POST("users/", h.SaveUser)
	e.PUT("users/", h.UpdateUserByID)
	e.DELETE("users/", h.DeleteUserByID)
	e.GET("users/", h.GetUserByID)

	e.Logger.Fatal(e.Start(":8081"))
}
