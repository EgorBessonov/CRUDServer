package main

import (
	"CRUDServer/internal/handler"
	"CRUDServer/internal/repository"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
		MongodbURL:    os.Getenv("MONGODB_URL"),
		SecretKey:     os.Getenv("SECRETKEY"),
	}

	if err := env.Parse(&cfg); err != nil {
		fmt.Println(err)
	}
	e := echo.New()

	h := handler.NewHandler(cfg)
	g := e.Group("/users")
	config := middleware.JWTConfig{
		Claims:     &handler.CustomClaims{},
		SigningKey: []byte(cfg.SecretKey),
	}
	g.Use(middleware.JWTWithConfig(config))

	g.POST("/saveUser/", h.SaveUser)
	g.PUT("/updateUser/", h.UpdateUserByID)
	g.DELETE("/deleteUser/", h.DeleteUserByID)
	g.GET("/getUser", h.GetUserByID)

	e.POST("registration/", h.Registration)
	e.POST("authentication/", h.Authentication)
	e.GET("refreshToken/", h.RefreshToken)
	e.POST("logout/", h.Logout)

	e.GET("images/downloadImage", h.DownloadImage)
	e.POST("images/uploadImage", h.UploadImage)

	e.Logger.Fatal(e.Start(":8081"))
}
