package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/captechconsulting/go-microservice-templates/lambda/internal/models"
)

type userService interface {
	ListUsers(ctx context.Context) ([]models.User, error)
	UpdateUser(ctx context.Context, ID int, user models.User) (models.User, error)
}

// API returns a HandlerFunc that handles incoming API Gateway proxy requests. It routes the
// requests to the appropriate handler based on the HTTP method.
func API(logger *slog.Logger, service userService) HandlerFunc {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.HTTPMethod {
		case http.MethodGet:
			return HandleListUsers(logger, service)(ctx, request)
		case http.MethodPost:
			return HandleUpdateUser(logger, service)(ctx, request)
		default:
			logger.Warn("Unsupported route", "method", request.HTTPMethod, "path", request.Path)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       http.StatusText(http.StatusNotFound),
			}, nil
		}
	}
}
