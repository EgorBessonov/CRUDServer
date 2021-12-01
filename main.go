package main

import (
	"CRUDServer/handler"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func init(){
	log.SetOutput(os.Stdin)
}

func main() {
	e := echo.New()
	h := handler.NewHandler("mongo")
	e.POST("users/", h.SaveUser)

	s := http.Server{
	  Addr:        ":8080",
	  Handler:     e,
	  
	}

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
	  log.Fatal(err)
	}
  }
  