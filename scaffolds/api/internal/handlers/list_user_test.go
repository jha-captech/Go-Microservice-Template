package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/captechconsulting/go-microservice-templates/api/internal/models"
	"github.com/captechconsulting/go-microservice-templates/api/internal/testutil"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"

	serviceMock "github.com/captechconsulting/go-microservice-templates/api/internal/handlers/mock"
)

func TestHandleListUsers(t *testing.T) {
	mockService := new(serviceMock.MockUserLister)
	logger := httplog.NewLogger("test")
	handler := HandleListUsers(logger, mockService)

	users := []models.User{
		{ID: 1, FirstName: "John", LastName: "Doe", Role: "Admin", UserID: 1001},
		{ID: 2, FirstName: "Jane", LastName: "Smith", Role: "User", UserID: 1002},
	}

	usersOut := mapMultipleOutput(users)

	tests := map[string]struct {
		mockCalled   bool
		mockOutput   []any
		expectedCode int
		expectedBody string
	}{
		"users returned": {
			mockCalled:   true,
			mockOutput:   []any{users, nil},
			expectedCode: http.StatusOK,
			expectedBody: testutil.ToJSONString(responseUsers{Users: usersOut}),
		},
		"no users found": {
			mockCalled:   true,
			mockOutput:   []any{[]models.User{}, nil},
			expectedCode: http.StatusOK,
			expectedBody: testutil.ToJSONString(responseUsers{Users: []outputUser{}}),
		},
		"internal server error": {
			mockCalled:   true,
			mockOutput:   []any{[]models.User{}, errors.New("teat error")},
			expectedCode: http.StatusInternalServerError,
			expectedBody: testutil.ToJSONString(responseErr{Error: "Error retrieving data"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/api/user", nil)
			assert.NoError(t, err)

			// Add chi URLParam
			rctx := chi.NewRouteContext()
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			if tc.mockCalled {
				mockService.
					On("ListUsers", ctx).
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
				mockService.AssertNotCalled(t, "ListUsers")
			}
		})
	}
}
