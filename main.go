package main

import (
	"CRUDServer/internal/configs"
	"CRUDServer/internal/handler"
	"CRUDServer/internal/repository"
	"CRUDServer/internal/service"
	"context"
	"fmt"
	"os"

	"github.com/caarlos0/env"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := configs.Config{
		CurrentDB:     os.Getenv("CURRENTDB"),
		PostgresdbURL: os.Getenv("POSTGRESDB_URL"),
		MongodbURL:    os.Getenv("MONGODB_URL"),
		SecretKey:     os.Getenv("SECRETKEY"),
	}

	if err := env.Parse(&cfg); err != nil {
		fmt.Println(err)
	}
	e := echo.New()

	_repository := dbConnection(cfg)
	h := handler.NewHandler(_repository)

	g := e.Group("/users")
	config := middleware.JWTConfig{
		Claims:     &service.CustomClaims{},
		SigningKey: []byte(cfg.SecretKey),
	}
	g.Use(middleware.JWTWithConfig(config))

	g.POST("/saveOrder/", h.SaveOrder)
	g.PUT("/updateOrder/", h.UpdateOrderByID)
	g.DELETE("/deleteOrder/", h.DeleteOrderByID)
	g.GET("/getOrder", h.GetOrderByID)

	e.POST("registration/", h.Registration)
	e.POST("authentication/", h.Authentication)
	e.GET("refreshToken/", h.RefreshToken)
	e.POST("logout/", h.Logout)

	e.GET("images/downloadImage", h.DownloadImage)
	e.POST("images/uploadImage", h.UploadImage)

	e.Logger.Fatal(e.Start(":8081"))
}

func dbConnection(cfg configs.Config) repository.Repository {
	switch cfg.CurrentDB {
	case "mongo":
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongodbURL))
		if err != nil {
			log.WithFields(log.Fields{
				"status": "connection to mongo database failed.",
				"err":    err,
			}).Info("mongo repository info")
		} else {
			log.WithFields(log.Fields{
				"status": "successfully connected to mongo database.",
			}).Info("mongo repository info.")
		}
		return repository.MongoRepository{DBconn: client}
	case "postgres":
		conn, err := pgxpool.Connect(context.Background(), cfg.PostgresdbURL)
		if err != nil {
			log.WithFields(log.Fields{
				"status": "connection to postgres database failed.",
				"err":    err,
			}).Info("postgres repository info.")
		} else {
			log.WithFields(log.Fields{
				"status": "successfully connected to postgres database.",
			}).Info("postgres repository info.")
		}
		return repository.PostgresRepository{DBconn: conn}
	}
	log.WithFields(log.Fields{
		"status": "database connection failed.",
		"err":    "invalid config",
	}).Info("repository info")
	return nil
}
