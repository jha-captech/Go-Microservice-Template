package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/captechconsulting/go-microservice-templates/lambda/internal/models"
	"github.com/captechconsulting/go-microservice-templates/lambda/internal/testutil"
	"github.com/stretchr/testify/assert"

	serviceMock "github.com/captechconsulting/go-microservice-templates/lambda/internal/handlers/mock"
)

func TestAPI(t *testing.T) {
	mockService := new(serviceMock.MockUserService)
	logger := slog.Default()
	handler := API(logger, mockService)

	users := []models.User{
		{ID: 0, FirstName: "John", LastName: "Doe", Role: "Customer", UserID: 1001},
		{ID: 1, FirstName: "John", LastName: "Doe", Role: "Customer", UserID: 1001},
		{ID: 2, FirstName: "Jane", LastName: "Smith", Role: "Employee", UserID: 1002},
	}
	userIn := inputUser{FirstName: "John", LastName: "Doe", Role: "Customer", UserID: 1001}
	usersOut := mapMultipleOutput(users)

	ctx := context.Background()

	tests := map[string]struct {
		mockCalled       bool
		mockSetup        func()
		request          events.APIGatewayProxyRequest
		expectedResponse events.APIGatewayProxyResponse
		expectedError    error
	}{
		"GET list users": {
			mockCalled: true,
			mockSetup: func() {
				mockService.
					On("ListUsers", ctx).
					Return(users, nil).
					Once()
			},
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodGet,
			},
			expectedResponse: events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       testutil.ToJSONString(responseUsers{Users: usersOut}),
			},
			expectedError: nil,
		},
		"POST update user": {
			mockCalled: true,
			mockSetup: func() {
				mockService.
					On("UpdateUser", ctx, 1, users[0]).
					Return(users[1], nil).
					Once()
			},
			request: events.APIGatewayProxyRequest{
				HTTPMethod:     http.MethodPost,
				PathParameters: map[string]string{"ID": "1"},
				Body:           testutil.ToJSONString(userIn),
			},
			expectedResponse: events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       testutil.ToJSONString(responseUser{User: usersOut[1]}),
			},
			expectedError: nil,
		},
		"PUT method not found": {
			mockCalled: false,
			mockSetup:  nil,
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPut,
			},
			expectedResponse: events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       http.StatusText(http.StatusNotFound),
			},
			expectedError: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.mockCalled {
				tc.mockSetup()
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
