package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/captechconsulting/go-microservice-templates/api/internal/models"
	"github.com/captechconsulting/go-microservice-templates/api/internal/testutil"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"
)

func TestHandleListUsers(t *testing.T) {
	users := []models.User{
		{ID: 1, FirstName: "John", LastName: "Doe", Role: "Admin", UserID: 1001},
		{ID: 2, FirstName: "Jane", LastName: "Smith", Role: "User", UserID: 1002},
	}

	usersOut := mapMultipleOutput(users)

	tests := map[string]struct {
		mockCalled    bool
		mockReturn    *sqlmock.Rows
		mockReturnErr error
		expectedCode  int
		expectedBody  string
	}{
		"users returned": {
			mockCalled:    true,
			mockReturn:    testutil.MustStructsToRows(users),
			mockReturnErr: nil,
			expectedCode:  http.StatusOK,
			expectedBody:  testutil.ToJSONString(responseUsers{Users: usersOut}),
		},
		"no users found": {
			mockCalled:    true,
			mockReturn:    testutil.MustStructToEmptyRow(models.User{}),
			mockReturnErr: nil,
			expectedCode:  http.StatusOK,
			expectedBody:  testutil.ToJSONString(responseUsers{Users: []outputUser{}}),
		},
		"internal server error - query failed": {
			mockCalled:    true,
			mockReturn:    &sqlmock.Rows{},
			mockReturnErr: errors.New("test"),
			expectedCode:  http.StatusInternalServerError,
			expectedBody:  testutil.ToJSONString(responseErr{Error: "Error retrieving data"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// setup
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			logger := httplog.NewLogger("test")
			handler := HandleListUsers(logger, db)

			req, err := http.NewRequest(http.MethodGet, "/api/user", nil)
			assert.NoError(t, err)

			// Add chi URLParam
			rctx := chi.NewRouteContext()
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
				WillReturnRows(tc.mockReturn).
				WillReturnError(tc.mockReturnErr)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code, "Wrong code received")
			assert.JSONEq(t, tc.expectedBody, rr.Body.String(), "Wrong response body")

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
