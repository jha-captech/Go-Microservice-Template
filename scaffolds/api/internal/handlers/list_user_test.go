package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/captechconsulting/go-microservice-templates/api/internal/handlers/mocks"
	"github.com/captechconsulting/go-microservice-templates/api/internal/models"
	"github.com/captechconsulting/go-microservice-templates/api/internal/testutil"
	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandleListUsers(t *testing.T) {
	users := []models.User{
		{ID: 1, FirstName: "John", LastName: "Doe", Role: "Admin", UserID: 1001},
		{ID: 2, FirstName: "Jane", LastName: "Smith", Role: "User", UserID: 1002},
	}

	usersOut := mapMultipleOutput(users)

	type mockArgs struct {
		mockCalled bool
		mockOutput []any
	}

	tests := map[string]struct {
		mockArgs     mockArgs
		expectedCode int
		expectedBody string
	}{
		"users returned": {
			mockArgs: mockArgs{
				mockCalled: true,
				mockOutput: []any{users, nil},
			},
			expectedCode: http.StatusOK,
			expectedBody: testutil.ToJSONString(responseUsers{Users: usersOut}),
		},
		"no users found": {
			mockArgs: mockArgs{
				mockCalled: true,
				mockOutput: []any{[]models.User{}, nil},
			},
			expectedCode: http.StatusOK,
			expectedBody: testutil.ToJSONString(responseUsers{Users: []outputUser{}}),
		},
		"internal server error": {
			mockArgs: mockArgs{
				mockCalled: true,
				mockOutput: []any{[]models.User{}, errors.New("teat error")},
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: testutil.ToJSONString(responseErr{Error: "Error retrieving data"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// setup
			mockService := new(mocks.MockUserLister)
			logger := httplog.NewLogger("test", httplog.Options{Writer: io.Discard})
			handler := HandleListUsers(logger, mockService)

			req, err := http.NewRequest(http.MethodGet, "/api/user", nil)
			require.NoError(t, err)

			if tc.mockArgs.mockCalled {
				mockService.
					On("ListUsers", mock.Anything).
					Return(tc.mockArgs.mockOutput...).
					Once()
			}

			// act
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			// assert
			assert.Equal(t, tc.expectedCode, rr.Code, "Wrong code received")
			assert.JSONEq(t, tc.expectedBody, rr.Body.String(), "Wrong response body")

			if tc.mockArgs.mockCalled {
				mockService.AssertExpectations(t)
			} else {
				mockService.AssertNotCalled(t, "ListUsers")
			}
		})
	}
}
