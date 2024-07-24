package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/captechconsulting/go-microservice-templates/lambda/internal/config"
	"github.com/captechconsulting/go-microservice-templates/lambda/internal/database"
	"github.com/captechconsulting/go-microservice-templates/lambda/internal/handlers"
	"github.com/captechconsulting/go-microservice-templates/lambda/internal/middleware"
	"github.com/captechconsulting/go-microservice-templates/lambda/internal/services"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalf("Startup failed. err: %v", err)
	}
}

// run is the main function that initializes the configuration, sets up logging, connects to the
// database, initializes the user service, and starts the AWS Lambda handler with the necessary
// middleware. It returns an error if any step in this initialization process fails.
func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("[in main.run] failed to load config: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))

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
		return fmt.Errorf("[in main.run]: %w", err)
	}

	defer func() {
		if err = db.Close(); err != nil {
			logger.Error("Error closing db connection", "err", err)
		}
	}()

	service := services.NewUserService(db)

	handler := handlers.API(logger, service)

	handler = middleware.AddToHandler(
		handler,
		middleware.Recovery(logger),
		middleware.Recovery(logger),
	)

	lambda.Start(handler)

	return nil
}
