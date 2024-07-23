package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/captechconsulting/go-microservice-templates/lambda/internal/model"
)

type UserService struct {
	database *sql.DB
}

// NewUserService returns a new UserService struct.
func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		database: db,
	}
}

// ListUsers returns a list of all UserService objects from the database.
func (s UserService) ListUsers(ctx context.Context) ([]model.User, error) {
	rows, err := s.database.QueryContext(
		ctx,
		`SELECT * FROM "users"`,
	)
	if err != nil {
		return []model.User{}, fmt.Errorf("[in service.ListUsers]: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Role, &user.UserID)
		if err != nil {
			return []model.User{}, fmt.Errorf("[in service.ListUsers]: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return []model.User{}, fmt.Errorf("[in service.ListUsers]: %w", err)
	}

	return users, nil
}

// UpdateUser updates am UserService objects from the database by ID.
func (s UserService) UpdateUser(ctx context.Context, ID int, user model.User) (model.User, error) {
	_, err := s.database.ExecContext(
		ctx,
		`
		UPDATE
			"users"
		SET
			"first_name" = $1,
			"last_name" = $2,
			"role" = $3,
			"user_id" = $4
		WHERE
			"id" = $5
		`,
		user.FirstName,
		user.LastName,
		user.Role,
		user.UserID,
		ID,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("[in service.UpdateUser]: %w", err)
	}

	user.ID = uint(ID)
	return user, nil
}
