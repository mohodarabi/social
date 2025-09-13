package main

import (
	"fmt"
	"net/http"
	"social/internal/db"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	db     db.DbRepo
}

type config struct {
	address  string
	dbConfig dbConfig
}

type dbConfig struct {
	url         string
	maxOpenConnections int
	maxIdleConnections int
	maxIdleTime        string
}

func (app *application) mount() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(middleware.Timeout(time.Second * 60))

	router.Route("/v1", func(route chi.Router) {
		route.Get("/health", app.healthCheckHandler)
	})

	return router
}

func (app *application) run(mux http.Handler) error {

	application := &http.Server{
		Addr:         app.config.address,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	fmt.Printf("App running on %s\n", app.config.address)

	return application.ListenAndServe()
}
