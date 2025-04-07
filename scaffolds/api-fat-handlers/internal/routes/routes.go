package routes

import (
	"database/sql"

	"github.com/captechconsulting/go-microservice-templates/api/internal/handlers"
	"github.com/captechconsulting/go-microservice-templates/api/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

type Option func(*routerOptions)

type routerOptions struct {
	registerHealthRoute bool
}

// WithRegisterHealthRoute controls whether a healthcheck route will be registered. If `false` is
// passed in or this function is not called, the default is `false`.
func WithRegisterHealthRoute(registerHealthRoute bool) Option {
	return func(options *routerOptions) {
		options.registerHealthRoute = registerHealthRoute
	}
}

func RegisterRoutes(router *chi.Mux, logger *httplog.Logger, svs *services.UserService, db *sql.DB, opts ...Option) {
	options := routerOptions{
		registerHealthRoute: false,
	}
	for _, opt := range opts {
		opt(&options)
	}

	if options.registerHealthRoute {
		router.Get("/lambda/health-check", handlers.HandleHealth(logger))
	}

	router.Get("/lambda/user", handlers.HandleListUsers(logger, db))
	router.Put("/lambda/user/{ID}", handlers.HandleUpdateUser(logger, svs))
}
