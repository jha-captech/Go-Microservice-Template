package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/captechconsulting/go-microservice-templates/sqs-lambda/internal/models"
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

// CreateUser creates am User objects in the database.
func (s UserService) CreateUser(ctx context.Context, user models.User) (int, error) {
	var ID int
	err := s.database.QueryRowContext(
		ctx,
		`
		INSERT INTO "users" ("first_name", "last_name", "role", "user_id")
			VALUES ($1, $2, $3, $4)
		RETURNING "id"
		`,
		user.FirstName,
		user.LastName,
		user.Role,
		user.UserID,
	).Scan(&ID)
	if err != nil {
		return 0, fmt.Errorf("[in services.CreateUser]: %w", err)
	}

	return ID, nil
}
