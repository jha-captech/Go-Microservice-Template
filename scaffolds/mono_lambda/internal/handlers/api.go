package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/captechconsulting/go-microservice-templates/lambda/internal/services"
)

func API(logger *slog.Logger, service *services.UserService) HandlerFunc {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.HTTPMethod {
		case http.MethodGet:
			return HandleListUsers(logger, service)(ctx, request)
		case http.MethodPost:
			return HandleUpdateUser(logger, service)(ctx, request)
		default:
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       http.StatusText(http.StatusNotFound),
			}, nil
		}
	}
}
