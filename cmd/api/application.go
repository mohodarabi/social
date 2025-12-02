package main

import (
	"net/http"
	"social/internal/db"
	"time"

	"social/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type application struct {
	config config
	db     db.DbRepo
	logger *zap.SugaredLogger
}

type config struct {
	address  string
	dbConfig dbConfig
}

type dbConfig struct {
	url                string
	apiUrl             string
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

		docsUrl := "/v1/swagger/doc.json"
		route.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(docsUrl),
		))

		route.Route("/post", func(route chi.Router) {
			route.Post("/", app.createPostHandler)

			route.Route("/{postID}", func(route chi.Router) {
				route.Use(app.postContentMiddleware)

				route.Get("/", app.getPostHandler)
				route.Delete("/", app.deletePostHandler)
				route.Patch("/", app.updatePostHandler)
			})
		})

		route.Route("/user", func(route chi.Router) {

			route.Post("/", app.createUserHandler)

			route.Route("/{userID}", func(route chi.Router) {
				route.Use(app.userContentMiddleware)

				route.Get("/", app.getUserHandler)
				route.Put("/follow", app.followUserHandler)
				route.Put("/unfollow", app.unfollowUserHandler)
			})

			route.Group(func(r chi.Router) {
				route.Get("/feed", app.getUserFeedHandler)
			})
		})
	})

	return router
}

func (app *application) run(mux http.Handler) error {

	docs.SwaggerInfo.Version = "0.0.1"
	docs.SwaggerInfo.Host = app.config.dbConfig.apiUrl
	docs.SwaggerInfo.BasePath = "/v1"

	application := &http.Server{
		Addr:         app.config.address,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("App running on", "addr", app.config.address)

	return application.ListenAndServe()
}
