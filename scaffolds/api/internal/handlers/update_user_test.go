package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	serviceMock "github.com/captechconsulting/go-microservice-templates/api/internal/handlers/mock"
	"github.com/captechconsulting/go-microservice-templates/api/internal/models"
	"github.com/captechconsulting/go-microservice-templates/api/internal/testutil"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"
)

func TestHandleUpdateUser(t *testing.T) {
	mockService := new(serviceMock.MockUserUpdater)
	logger := httplog.NewLogger("test")
	handler := HandleUpdateUser(logger, mockService)

	user := models.User{FirstName: "John", LastName: "Doe", Role: "Customer", UserID: 1001}
	userIn := inputUser{FirstName: "John", LastName: "Doe", Role: "Customer", UserID: 1001}
	userOut := mapOutput(user)

	tests := map[string]struct {
		mockCalled     bool
		mockInput      []any
		mockOutput     []any
		requestIDParam string
		requestBody    string
		expectedCode   int
		expectedBody   string
	}{
		"valid request, user updated": {
			mockCalled:     true,
			mockInput:      []any{1, user},
			mockOutput:     []any{user, nil},
			requestIDParam: "1",
			requestBody:    testutil.ToJSONString(userIn),
			expectedCode:   http.StatusOK,
			expectedBody:   testutil.ToJSONString(responseUser{User: userOut}),
		},
		"invalid request body": {
			mockCalled:     false,
			mockInput:      nil,
			mockOutput:     nil,
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
			mockCalled:     true,
			mockInput:      []any{1, user},
			mockOutput:     []any{models.User{}, errors.New("creation error")},
			requestIDParam: "1",
			requestBody:    testutil.ToJSONString(userIn),
			expectedCode:   http.StatusInternalServerError,
			expectedBody:   testutil.ToJSONString(responseErr{Error: "Error updating object"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPut, "/api/user/"+tc.requestIDParam, strings.NewReader(tc.requestBody))
			assert.NoError(t, err)

			// Add chi URLParam
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("ID", tc.requestIDParam)
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			if tc.mockCalled {
				mockService.
					On("UpdateUser", append([]any{ctx}, tc.mockInput...)...).
					Return(tc.mockOutput...).
					Once()
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code, "Wrong code received")
			assert.JSONEq(t, tc.expectedBody, rr.Body.String(), "Wrong response body")

			if tc.mockCalled {
				mockService.AssertExpectations(t)
			} else {
				mockService.AssertNotCalled(t, "UpdateUser")
			}
		})
	}
}
