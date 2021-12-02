package main

import (
	/*"CRUDServer/handler"
	"CRUDServer/repository"
	"net/http"
	"os"
	"github.com/labstack/echo/v4"*/
	"CRUDServer/repository"
	"fmt"
)

func main() {
	u := repository.User{
		UserId:   "123e4567-e89b-12d3-a456-42661417412b",
		UserName: "Andersen",
		UserAge:  18,
		IsAdult:  true,
	}
	h := repository.MongoRepository{}
	newu, err := h.ReadUser(u.UserId)
	if err != nil {
		fmt.Println(newu)
	}
	/*e := echo.New()
	h := handler.NewHandler("mongo")
	e.POST("users/", h.SaveUser)

	s := http.Server{
	  Addr:        ":8080",
	  Handler:     e,

	}

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
	  log.Fatal(err)
	}*/
}
