package services

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/captechconsulting/go-microservice-templates/lambda/internal/models"
	"github.com/captechconsulting/go-microservice-templates/lambda/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testSuit struct {
	suite.Suite
	service *UserService
	dbMock  sqlmock.Sqlmock
}

func TestTestSuit(t *testing.T) {
	suite.Run(t, new(testSuit))
}

func (s *testSuit) SetupSuite() {
	db, mock, err := sqlmock.New()
	assert.NoError(s.T(), err)

	s.dbMock = mock
	s.service = NewUserService(db)
}

func (s *testSuit) TearDownSuite() {
	_ = s.service.database.Close()
}

func (s *testSuit) TestListUsers() {
	t := s.T()

	users := []models.User{
		{ID: 1, FirstName: "John", LastName: "Doe", Role: "Admin", UserID: 1001},
		{ID: 2, FirstName: "Jane", LastName: "Smith", Role: "User", UserID: 1002},
	}

	testCases := map[string]struct {
		mockReturn     *sqlmock.Rows
		mockReturnErr  error
		expectedReturn []models.User
		expectedError  error
	}{
		"Return slice of users": {
			mockReturn:     testutil.MustStructsToRows(users),
			mockReturnErr:  nil,
			expectedReturn: users,
			expectedError:  nil,
		},
		"Error getting users": {
			mockReturn:     &sqlmock.Rows{},
			mockReturnErr:  errors.New("test"),
			expectedReturn: []models.User{},
			expectedError:  fmt.Errorf("[in services.ListUsers] failed to get users: %w", errors.New("test")),
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			exp := `SELECT * FROM "users"`
			s.dbMock.
				ExpectQuery(regexp.QuoteMeta(exp)).
				WillReturnRows(tc.mockReturn).
				WillReturnError(tc.mockReturnErr)

			actualReturn, err := s.service.ListUsers(context.Background())

			assert.Equal(t, tc.expectedError, err, "errors did not match")
			assert.Equal(t, tc.expectedReturn, actualReturn, "returned data does not match")

			err = s.dbMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func (s *testSuit) TestUpdateUser() {
	t := s.T()

	userIn := models.User{ID: 0, FirstName: "John", LastName: "Doe", Role: "Admin", UserID: 1001}
	userOut := models.User{ID: 1, FirstName: "John", LastName: "Doe", Role: "Admin", UserID: 1001}

	testCases := map[string]struct {
		mockInputArgs  []driver.Value
		mockReturn     driver.Result
		mockReturnErr  error
		inputID        int
		inputUser      models.User
		expectedReturn models.User
		expectedError  error
	}{
		"user updated by ID": {
			mockInputArgs:  []driver.Value{userIn.FirstName, userIn.LastName, userIn.Role, userIn.UserID, int(userOut.ID)},
			mockReturn:     sqlmock.NewResult(1, 1),
			mockReturnErr:  nil,
			inputID:        int(userOut.ID),
			inputUser:      userIn,
			expectedReturn: userOut,
			expectedError:  nil,
		},
		"Error updating user": {
			mockInputArgs:  []driver.Value{userIn.FirstName, userIn.LastName, userIn.Role, userIn.UserID, 0},
			mockReturn:     nil,
			mockReturnErr:  errors.New("test"),
			inputID:        0,
			inputUser:      userIn,
			expectedReturn: models.User{},
			expectedError:  fmt.Errorf("[in services.UpdateUser] failed to update user: %w", errors.New("test")),
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			exp := `
				UPDATE
					"users"
				SET
					"first_name" = $1,
					"last_name" = $2,
					"role" = $3,
					"user_id" = $4
				WHERE
					"id" = $5
			`
			s.dbMock.
				ExpectExec(regexp.QuoteMeta(exp)).
				WithArgs(tc.mockInputArgs...).
				WillReturnResult(tc.mockReturn).
				WillReturnError(tc.mockReturnErr)

			actualReturn, err := s.service.UpdateUser(context.Background(), tc.inputID, tc.inputUser)

			assert.Equal(t, tc.expectedError, err, "errors did not match")
			assert.Equal(t, tc.expectedReturn, actualReturn, "returned data does not match")

			err = s.dbMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
