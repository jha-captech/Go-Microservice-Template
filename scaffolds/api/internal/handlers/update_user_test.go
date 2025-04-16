package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/captechconsulting/go-microservice-templates/api/internal/models"
	"github.com/captechconsulting/go-microservice-templates/api/internal/testutil"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"
)

func TestHandleUpdateUser(t *testing.T) {
	user := models.User{FirstName: "John", LastName: "Doe", Role: "Customer", UserID: 1001}
	userIn := inputUser{FirstName: "John", LastName: "Doe", Role: "Customer", UserID: 1001}
	userOut := mapOutput(user)

	type mockArgs struct {
		mockCalledCount int
		mockOutput1     models.User
		mockOutput2     error
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
				mockCalledCount: 1,
				mockOutput1:     user,
				mockOutput2:     nil,
			},
			requestIDParam: "1",
			requestBody:    testutil.ToJSONString(userIn),
			expectedCode:   http.StatusOK,
			expectedBody:   testutil.ToJSONString(responseUser{User: userOut}),
		},
		"invalid request body": {
			mockArgs: mockArgs{
				mockCalledCount: 0,
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
				mockCalledCount: 1,
				mockOutput1:     models.User{},
				mockOutput2:     errors.New("some error"),
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
			mockService := new(moqUserUpdater)
			if tc.mockArgs.mockCalledCount > 0 {
				mockService.UpdateUserFunc = func(ctx context.Context, ID int, user models.User) (models.User, error) {
					return tc.mockArgs.mockOutput1, tc.mockArgs.mockOutput2
				}
			}

			logger := httplog.NewLogger("test", httplog.Options{Writer: io.Discard})
			handler := HandleUpdateUser(logger, mockService)

			req, err := http.NewRequest(http.MethodPut, "/lambda/user/"+tc.requestIDParam, strings.NewReader(tc.requestBody))
			assert.NoError(t, err)

			// Add chi URLParam
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("ID", tc.requestIDParam)
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code, "Wrong code received")
			assert.JSONEq(t, tc.expectedBody, rr.Body.String(), "Wrong response body")
			assert.Equal(t, tc.mockArgs.mockCalledCount, len(mockService.UpdateUserCalls()))
		})
	}
}
