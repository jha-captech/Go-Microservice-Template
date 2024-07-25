package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	serviceMock "github.com/captechconsulting/go-microservice-templates/lambda/internal/handlers/mock"
	"github.com/captechconsulting/go-microservice-templates/lambda/internal/models"
	"github.com/captechconsulting/go-microservice-templates/lambda/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestHandleUpdateUser(t *testing.T) {
	mockService := new(serviceMock.MockUserUpdater)
	logger := slog.Default()
	handler := HandleUpdateUser(logger, mockService)

	user := models.User{FirstName: "John", LastName: "Doe", Role: "Customer", UserID: 1001}
	userIn := inputUser{FirstName: "John", LastName: "Doe", Role: "Customer", UserID: 1001}
	userOut := mapOutput(user)

	ctx := context.Background()

	tests := map[string]struct {
		mockCalled       bool
		mockInput        []any
		mockOutput       []any
		request          events.APIGatewayProxyRequest
		expectedResponse events.APIGatewayProxyResponse
		expectedError    error
	}{
		"valid request, user updated": {
			mockCalled: true,
			mockInput:  []any{ctx, 1, user},
			mockOutput: []any{user, nil},
			request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"ID": "1"},
				Body:           testutil.ToJSONString(userIn),
			},
			expectedResponse: events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       testutil.ToJSONString(responseUser{User: userOut}),
			},
			expectedError: nil,
		},
		"invalid ID": {
			mockCalled: false,
			mockInput:  nil,
			mockOutput: nil,
			request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"ID": "test"},
				Body:           testutil.ToJSONString(responseUser{User: userOut}),
			},
			expectedResponse: events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       testutil.ToJSONString(responseErr{Error: "Not a valid ID"}),
			},
			expectedError: nil,
		},
		"invalid request body": {
			mockCalled: false,
			mockInput:  nil,
			mockOutput: nil,
			request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"ID": "1"},
				Body:           `{"first_name": "John","role": "Admin", "user_id": -1}`,
			},
			expectedResponse: events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body: testutil.ToJSONString(responseErr{
					ValidationErrors: []problem{
						{
							Name:        "last_name",
							Description: "must not be blank",
						},
						{
							Name:        "role",
							Description: `must be "Customer" or "Employee"`,
						},
						{
							Name:        "user_id",
							Description: "must be must be greater than zero",
						},
					},
				}),
			},
			expectedError: nil,
		},
		"error creating user": {
			mockCalled: true,
			mockInput:  []any{ctx, 1, user},
			mockOutput: []any{models.User{}, errors.New("creation error")},
			request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"ID": "1"},
				Body:           testutil.ToJSONString(userIn),
			},
			expectedResponse: events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       testutil.ToJSONString(responseErr{Error: "Error updating object"}),
			},
			expectedError: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.mockCalled {
				mockService.
					On("UpdateUser", tc.mockInput...).
					Return(tc.mockOutput...).
					Once()
			}

			got, err := handler(ctx, tc.request)

			assert.Equal(t, tc.expectedError, err, "Error expectations not met")
			assert.Equal(t, tc.expectedResponse, got, "Wrong response body")

			if tc.mockCalled {
				mockService.AssertExpectations(t)
			} else {
				mockService.AssertNotCalled(t, "ListUsers")
			}
		})
	}
}
