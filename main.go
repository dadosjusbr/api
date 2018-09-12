// +build heroku
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
)

type config struct {
	Port int `envconfig:"PORT"`
}

func main() {
	var conf config
	err := envconfig.Process("remuneracao-magistrados", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Going to start listening at port:%d\n", conf.Port)
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.Port),
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}
	e := echo.New()
	e.GET("/", func(c echo.Context) error { return c.HTML(http.StatusOK, "<html><h1>Olá!</h1></html>") })
	e.Logger.Fatal(e.StartServer(s))
}
