package main

import (
	"context"
	"fmt"
	"github.com/EgorBessonov/CRUDServer/internal/cache"
	"github.com/EgorBessonov/CRUDServer/internal/config"
	"github.com/EgorBessonov/CRUDServer/internal/handler"
	"github.com/EgorBessonov/CRUDServer/internal/repository"
	"github.com/EgorBessonov/CRUDServer/internal/service"

	_ "github.com/EgorBessonov/CRUDServer/docs"
	"github.com/caarlos0/env"
	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title CRUDServer
// @version 1.0
//description This is a simple crud server for mongo & postgres databases with jwt authentication

// @host localhost:8081
// @BasePath /

//securityDefinition.apikey ApiKeyAuth
// @in header
// @Authentication
func main() {
	cfg := configs.Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Println(err)
	}
	e := echo.New()

	repo := dbConnection(cfg)
	redisClient := redisConnection(cfg)
	defer func() {
		err := redisClient.Close()
		if err != nil {
			log.Errorf("error while closing redis connection - %e", err)
		}
	}()
	c := cache.NewCache(redisClient.Context(), cfg, redisClient)
	s := service.NewService(repo, c)
	h := handler.NewHandler(s, &cfg)
	g := e.Group("/orders")
	config := middleware.JWTConfig{
		Claims:     &service.CustomClaims{},
		SigningKey: []byte(cfg.SecretKey),
	}
	g.Use(middleware.JWTWithConfig(config))

	g.POST("/saveOrder", h.SaveOrder)
	g.PUT("/updateOrder", h.UpdateOrderByID)
	g.DELETE("/deleteOrder", h.DeleteOrderByID)
	g.GET("/getOrder", h.GetOrderByID)

	e.POST("/registration", h.Registration)
	e.POST("/authentication", h.Authentication)
	e.GET("/refreshToken", h.RefreshToken)
	e.POST("/logout", h.Logout)

	e.GET("/images/downloadImage", h.DownloadImage)
	e.POST("/images/uploadImage", h.UploadImage)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
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

func redisConnection(cfg configs.Config) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: "",
		DB:       0,
	})
	if _, err := redisClient.Ping().Result(); err != nil {
		log.WithFields(log.Fields{
			"status": "error while connection to redisdb",
			"err":    err,
		}).Info("redis repository info.")
		return nil
	}
	log.WithFields(log.Fields{
		"status": "successfully connected to redisdb",
	}).Info("redis repository info.")
	return redisClient
}
