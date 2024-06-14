package main

import (
	"erp-user-service/handlers"
	"erp-user-service/middlewares"
	"fmt"
	"log"
	"net/http"
)

type Config struct {
	Middleware *middlewares.MiddlewareConfig
	Handlers   *handlers.HandlerConfig
}

const (
	grpcPort = "50001"
	webPort  = "80"
)

func main() {
	app := Config{
		Middleware: middlewares.New(),
		Handlers:   handlers.New(),
	}

	go app.startGRPC()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}
