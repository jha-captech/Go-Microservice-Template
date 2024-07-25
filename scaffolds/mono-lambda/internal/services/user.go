package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/captechconsulting/go-microservice-templates/lambda/internal/models"
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
func (s UserService) ListUsers(ctx context.Context) ([]models.User, error) {
	rows, err := s.database.QueryContext(
		ctx,
		`SELECT * FROM "users"`,
	)
	if err != nil {
		return []models.User{}, fmt.Errorf("[in services.ListUsers] failed to get users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Role, &user.UserID)
		if err != nil {
			return []models.User{}, fmt.Errorf("[in services.ListUsers] failed to scan user from row: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return []models.User{}, fmt.Errorf("[in services.ListUsers] failed to scan users: %w", err)
	}

	return users, nil
}

// UpdateUser updates am UserService objects from the database by ID.
func (s UserService) UpdateUser(ctx context.Context, ID int, user models.User) (models.User, error) {
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
		return models.User{}, fmt.Errorf("[in services.UpdateUser] failed to update user: %w", err)
	}

	user.ID = uint(ID)
	return user, nil
}
