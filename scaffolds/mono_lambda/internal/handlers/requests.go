package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/captechconsulting/go-microservice-templates/lambda/internal/models"
)

type inputUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	UserID    int    `json:"user_id"`
}

func (user inputUser) MapTo() (models.User, error) {
	return models.User{
		ID:        0,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		UserID:    uint(user.UserID),
	}, nil
}

func (user inputUser) Valid() []problem {
	var problems []problem

	// validate FirstName is not blank
	if user.FirstName == "" {
		problems = append(problems, problem{
			Name:        "first_name",
			Description: "must not be blank",
		})
	}

	// validate LastName is not blank
	if user.LastName == "" {
		problems = append(problems, problem{
			Name:        "last_name",
			Description: "must not be blank",
		})
	}

	// validate role is not blank and is `Customer` or `Employee`
	if user.Role == "" {
		problems = append(problems, problem{
			Name:        "role",
			Description: "must not be blank",
		})
	} else if user.Role != "Customer" && user.Role != "Employee" {
		problems = append(problems, problem{
			Name:        "role",
			Description: `must be "Customer" or "Employee"`,
		})
	}

	// validate UserID greater than 0
	if user.UserID < 1 {
		problems = append(problems, problem{
			Name:        "user_id",
			Description: "must be must be greater than zero",
		})
	}

	return problems
}

type problem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Validator interface {
	Valid() (problems []problem)
}

type Mapper[T any] interface {
	MapTo() (T, error)
}

type ValidatorMapper[T any] interface {
	Validator
	Mapper[T]
}

func decodeValidateBody[I ValidatorMapper[O], O any](requestBody string) (O, []problem, error) {
	var inputModel I

	// decode to JSON
	if err := json.Unmarshal([]byte(requestBody), &inputModel); err != nil {
		return *new(O), nil, fmt.Errorf("[in decodeValidateBody] decode json: %w", err)
	}

	// validate
	if problems := inputModel.Valid(); len(problems) > 0 {
		return *new(O), problems, fmt.Errorf(
			"[in decodeValidateBody] invalid %T: %d problems", inputModel, len(problems),
		)
	}

	// map to return type
	data, err := inputModel.MapTo()
	if err != nil {
		return *new(O), nil, fmt.Errorf(
			"[in decodeValidateBody] error mapping input %T to %T: %w",
			*new(I),
			*new(O),
			err,
		)
	}

	return data, nil, nil
}
