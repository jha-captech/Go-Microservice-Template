package handlers

import (
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/captechconsulting/go-microservice-templates/lambda/internal/models"
)

type outputUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	UserID    int    `json:"user_id"`
}

func mapOutput(user models.User) outputUser {
	return outputUser{
		ID:        int(user.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		UserID:    int(user.UserID),
	}
}

func mapMultipleOutput(user []models.User) []outputUser {
	usersOut := make([]outputUser, len(user))
	for i := 0; i < len(user); i++ {
		userOut := mapOutput(user[i])
		usersOut[i] = userOut
	}

	return usersOut
}

type responseUser struct {
	User outputUser `json:"user"`
}

type responseUsers struct {
	Users []outputUser `json:"users"`
}

type responseMsg struct {
	Message string `json:"message"`
}

type responseID struct {
	ObjectID int `json:"object_id"`
}

type responseErr struct {
	Error            string    `json:"error,omitempty"`
	ValidationErrors []problem `json:"validation_errors,omitempty"`
}

// encodeResponse encodes data as a JSON response.
func encodeResponse(logger *slog.Logger, status int, data any) (events.APIGatewayProxyResponse, error) {
	JSONData, err := json.Marshal(data)
	if err != nil {
		logger.Error("Error while marshaling data", "err", err, "data", data)
		return events.APIGatewayProxyResponse{
			StatusCode: status,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"Error": "Internal server error"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(JSONData),
	}, nil
}
