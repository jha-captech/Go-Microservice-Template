package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/captechconsulting/go-microservice-templates/api/internal/config"
	"github.com/captechconsulting/go-microservice-templates/api/internal/database"
	"github.com/captechconsulting/go-microservice-templates/api/internal/routes"
	"github.com/captechconsulting/go-microservice-templates/api/internal/services"
	"github.com/captechconsulting/go-microservice-templates/api/internal/swagger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalf("Startup failed. err: %v", err)
	}
}

// run is the main function that initializes the configuration, sets up logging, connects to the
// database, initializes the user service, and starts the API with the necessary middleware. It
// also has graceful stop logic implemented.It returns an error if any step in this initialization
// process fails.
func run(ctx context.Context) error {
	// Setup
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("[in run]: %w", err)
	}

	logger := httplog.NewLogger("user-microservice", httplog.Options{
		LogLevel:        cfg.LogLevel,
		JSON:            false,
		Concise:         true,
		ResponseHeaders: false,
	})

	db, err := database.New(
		ctx,
		fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.DBHost,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
			cfg.DBPort,
		),
		logger,
		time.Duration(cfg.DBRetryDuration)*time.Second,
	)
	if err != nil {
		return fmt.Errorf("[in run]: %w", err)
	}

	defer func() {
		if err = db.Close(); err != nil {
			logger.Error("Error closing db connection", "err", err)
		}
	}()

	router := chi.NewRouter()

	router.Use(httplog.RequestLogger(logger))
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "PUT", "POST", "DELETE"},
		MaxAge:         300,
	}))

	svs := services.NewUserService(db)
	routes.RegisterRoutes(router, logger, svs, routes.WithRegisterHealthRoute(true))

	if cfg.HTTPUseSwagger {
		swagger.RunSwagger(router, logger, cfg.HTTPDomain+cfg.HTTPPort)
	}

	serverInstance := &http.Server{
		Addr:              cfg.HTTPDomain + cfg.HTTPPort,
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: 500 * time.Millisecond,
		ReadTimeout:       500 * time.Millisecond,
		WriteTimeout:      500 * time.Millisecond,
		Handler:           router,
	}

	// Graceful shutdown
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		fmt.Println()
		logger.Info("Shutdown signal received")

		shutdownCtx, err := context.WithTimeout(
			serverCtx, time.Duration(cfg.HTTPShutdownDuration)*time.Second,
		)
		if err != nil {
			log.Fatalf("Error creating context.WithTimeout. err: %v", err)
		}

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		if err := serverInstance.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Error shutting down server. err: %v", err)
		}
		serverStopCtx()
	}()

	// Run
	logger.Info(fmt.Sprintf("Server is listening on %s", serverInstance.Addr))
	err = serverInstance.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	<-serverCtx.Done()
	logger.Info("Shutdown complete")
	return nil
}
