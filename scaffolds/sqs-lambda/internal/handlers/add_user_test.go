package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/captechconsulting/go-microservice-templates/sqs-lambda/internal/handlers/mock"
	"github.com/captechconsulting/go-microservice-templates/sqs-lambda/internal/models"
	"github.com/captechconsulting/go-microservice-templates/sqs-lambda/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreateUsers(t *testing.T) {
	mockService := new(mock.MockUserCreator)
	logger := slog.Default()
	handler := HandleCreateUsers(logger, mockService)

	users := []models.User{
		{ID: 0, FirstName: "John", LastName: "Doe", Role: "Customer", UserID: 1001},
		{ID: 0, FirstName: "Jane", LastName: "Smith", Role: "Employee", UserID: 1002},
	}

	usersIn := []inputUser{
		{FirstName: "John", LastName: "Doe", Role: "Customer", UserID: 1001},
		{FirstName: "Jane", LastName: "Smith", Role: "Employee", UserID: 1002},
	}

	ctx := context.TODO()

	type mockDetail struct {
		mockCalled bool
		mockInput  []any
		mockOutput []any
	}

	tests := map[string]struct {
		mockCalled       bool
		mockDetails      []mockDetail
		request          events.SQSEvent
		expectedResponse returnFailures
		expectedError    error
	}{
		"no issues - users created": {
			mockDetails: []mockDetail{
				{
					mockCalled: true,
					mockInput:  []any{ctx, users[0]},
					mockOutput: []any{1, nil},
				},
				{
					mockCalled: true,
					mockInput:  []any{ctx, users[1]},
					mockOutput: []any{2, nil},
				},
			},
			request: events.SQSEvent{
				Records: []events.SQSMessage{
					{
						MessageId: "1",
						Body:      testutil.ToJSONString(usersIn[0]),
					},
					{
						MessageId: "2",
						Body:      testutil.ToJSONString(usersIn[1]),
					},
				},
			},
			expectedResponse: returnFailures{BatchItemFailures: []failedItems(nil)},
			expectedError:    nil,
		},
		"one validation issue": {
			mockDetails: []mockDetail{
				{
					mockCalled: false,
					mockInput:  nil,
					mockOutput: nil,
				},
				{
					mockCalled: true,
					mockInput:  []any{ctx, users[1]},
					mockOutput: []any{2, nil},
				},
			},
			request: events.SQSEvent{
				Records: []events.SQSMessage{
					{
						MessageId: "1",
						Body: testutil.ToJSONString(
							inputUser{FirstName: "John", LastName: "Doe", Role: "Person", UserID: 1001},
						),
					},
					{
						MessageId: "2",
						Body:      testutil.ToJSONString(usersIn[1]),
					},
				},
			},
			expectedResponse: returnFailures{BatchItemFailures: []failedItems{{
				ItemIdentifier: "1",
			}}},
			expectedError: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			for _, detail := range tc.mockDetails {
				if detail.mockCalled {
					fmt.Println("mock")
					mockService.
						On("CreateUser", detail.mockInput...).
						Return(detail.mockOutput...).
						Once()
				}
			}

			got, err := handler(ctx, tc.request)

			assert.Equal(t, tc.expectedError, err, "Error expectations not met")
			assert.Equal(t, tc.expectedResponse, got, "Wrong response body")

			mockCalled := false
			for _, detail := range tc.mockDetails {
				if detail.mockCalled {
					mockCalled = detail.mockCalled
					break
				}
			}
			if mockCalled {
				mockService.AssertExpectations(t)
			} else {
				mockService.AssertNotCalled(t, "ListUsers")
			}
		})
	}
}
