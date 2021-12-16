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
		PostgresdbURL: os.Getenv("POSTGRESDB_URL"),
		MongodbURL:    os.Getenv("MONGODB_URL")}

	if err := env.Parse(&cfg); err != nil {
		fmt.Println(err)
	}
	e := echo.New()

	h := handler.NewHandler(cfg)
	
	e.POST("users/saveUser/", h.SaveUser)
	e.PUT("users/updateUser/", h.UpdateUserByID)
	e.DELETE("users/deleteUser/", h.DeleteUserByID)
	e.GET("users/getUser", h.GetUserByID)

	e.POST("auth/registration/", h.Registration)
	e.POST("auth/authentication/", h.Authentication)
	e.POST("auth/authorization/", h.Authorization)
	e.POST("auth/refreshToken/", h.RefreshToken)
	e.POST("auth/logout/", h.Logout)

	e.Logger.Fatal(e.Start(":8081"))
}
