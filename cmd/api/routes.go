package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	mux.Use(middleware.Heartbeat("/ping"))
	mux.Use(app.Middleware.GetDeviceInfo)

	mux.Post("/", app.Handlers.CreateUser)
	mux.Post("/login", app.Handlers.Login)

	protectedRoutes := chi.NewRouter()
	protectedRoutes.Use(app.Middleware.Authenticated)
	protectedRoutes.Get("/refresh", app.Handlers.RefreshUserSession)
	protectedRoutes.Get("/whoami", app.Handlers.WhoAmI)
	protectedRoutes.Get("/devices", app.Handlers.GetDeviceList)
	protectedRoutes.Delete("/devices", app.Handlers.RevokeUserDevices)
	protectedRoutes.Post("/thirdparty", app.Handlers.AddThirdPartyAuthenticator)

	mux.Mount("/authentication", protectedRoutes)

	return mux
}
