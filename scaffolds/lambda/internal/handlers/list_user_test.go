package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/captechconsulting/go-microservice-templates/lambda/internal/model"
	"github.com/captechconsulting/go-microservice-templates/lambda/internal/testutil"
	"github.com/stretchr/testify/assert"

	serviceMock "github.com/captechconsulting/go-microservice-templates/lambda/internal/handlers/mock"
)

func TestHandleListUsers(t *testing.T) {
	mockService := new(serviceMock.MockUserLister)
	logger := slog.Default()
	handler := HandleListUsers(logger, mockService)

	users := []model.User{
		{ID: 1, FirstName: "John", LastName: "Doe", Role: "Admin", UserID: 1001},
		{ID: 2, FirstName: "Jane", LastName: "Smith", Role: "User", UserID: 1002},
	}

	usersOut := mapMultipleOutput(users)

	ctx := context.Background()

	tests := map[string]struct {
		mockCalled       bool
		mockOutput       []any
		request          events.APIGatewayProxyRequest
		expectedResponse events.APIGatewayProxyResponse
		expectedError    error
	}{
		"users returned": {
			mockCalled: true,
			mockOutput: []any{users, nil},
			request:    events.APIGatewayProxyRequest{},
			expectedResponse: events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       testutil.ToJSONString(responseUsers{Users: usersOut}),
			},
			expectedError: nil,
		},
		"no users found": {
			mockCalled: true,
			mockOutput: []any{[]model.User{}, nil},
			request:    events.APIGatewayProxyRequest{},
			expectedResponse: events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       testutil.ToJSONString(responseUsers{Users: []outputUser{}}),
			},
			expectedError: nil,
		},
		"internal server error": {
			mockCalled: true,
			mockOutput: []any{[]model.User{}, errors.New("teat error")},
			request:    events.APIGatewayProxyRequest{},
			expectedResponse: events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       testutil.ToJSONString(responseErr{Error: "Error retrieving data"}),
			},
			expectedError: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.mockCalled {
				mockService.
					On("ListUsers", ctx).
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
