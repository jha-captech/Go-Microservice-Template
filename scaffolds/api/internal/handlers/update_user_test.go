package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/captechconsulting/go-microservice-templates/api/internal/handlers/mocks"
	"github.com/captechconsulting/go-microservice-templates/api/internal/models"
	"github.com/captechconsulting/go-microservice-templates/api/internal/testutil"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleUpdateUser(t *testing.T) {
	user := models.User{FirstName: "John", LastName: "Doe", Role: "Customer", UserID: 1001}
	userIn := inputUser{FirstName: "John", LastName: "Doe", Role: "Customer", UserID: 1001}
	userOut := mapOutput(user)

	type mockArgs struct {
		mockCalled bool
		mockInput  []any
		mockOutput []any
	}

	tests := map[string]struct {
		mockArgs       mockArgs
		requestIDParam string
		requestBody    string
		expectedCode   int
		expectedBody   string
	}{
		"valid request, user updated": {
			mockArgs: mockArgs{
				mockCalled: true,
				mockInput:  []any{mock.Anything, 1, user},
				mockOutput: []any{user, nil},
			},
			requestIDParam: "1",
			requestBody:    testutil.ToJSONString(userIn),
			expectedCode:   http.StatusOK,
			expectedBody:   testutil.ToJSONString(responseUser{User: userOut}),
		},
		"invalid request body": {
			mockArgs: mockArgs{
				mockCalled: false,
				mockInput:  nil,
				mockOutput: nil,
			},
			requestIDParam: "1",
			requestBody:    `{"first_name":"John","role":"Admin"}`,
			expectedCode:   http.StatusBadRequest,
			expectedBody: testutil.ToJSONString(responseErr{
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
		"error creating user": {
			mockArgs: mockArgs{
				mockCalled: true,
				mockInput:  []any{mock.Anything, 1, user},
				mockOutput: []any{models.User{}, errors.New("creation error")},
			},
			requestIDParam: "1",
			requestBody:    testutil.ToJSONString(userIn),
			expectedCode:   http.StatusInternalServerError,
			expectedBody:   testutil.ToJSONString(responseErr{Error: "Error updating object"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// setup
			mockService := new(mocks.MockUserUpdater)
			logger := httplog.NewLogger("test", httplog.Options{Writer: io.Discard})
			handler := HandleUpdateUser(logger, mockService)

			req, err := http.NewRequest(http.MethodPut, "/lambda/user/"+tc.requestIDParam, strings.NewReader(tc.requestBody))
			assert.NoError(t, err)

			// Add chi URLParam
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("ID", tc.requestIDParam)
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			if tc.mockArgs.mockCalled {
				mockService.
					On("UpdateUser", tc.mockArgs.mockInput...).
					Return(tc.mockArgs.mockOutput...).
					Once()
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code, "Wrong code received")
			assert.JSONEq(t, tc.expectedBody, rr.Body.String(), "Wrong response body")

			if tc.mockArgs.mockCalled {
				mockService.AssertExpectations(t)
			} else {
				mockService.AssertNotCalled(t, "UpdateUser")
			}
		})
	}
}
