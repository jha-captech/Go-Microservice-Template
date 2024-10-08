package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/captechconsulting/go-microservice-templates/lambda/internal/models"
)

type userLister interface {
	ListUsers(ctx context.Context) ([]models.User, error)
}

// HandleListUsers returns a HandlerFunc that handles GET requests to list users. It retrieves the
// list of users from the provided service and returns them in the response.
func HandleListUsers(logger *slog.Logger, service userLister) HandlerFunc {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// get values from database
		users, err := service.ListUsers(ctx)
		if err != nil {
			logger.Error("error getting all locations", "err", err)
			return encodeResponse(logger, http.StatusInternalServerError, responseErr{
				Error: "Error retrieving data",
			})
		}

		// return response
		usersOut := mapMultipleOutput(users)
		return encodeResponse(logger, http.StatusOK, responseUsers{
			Users: usersOut,
		})
	}
}
