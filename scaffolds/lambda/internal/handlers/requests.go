package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/captechconsulting/go-microservice-templates/lambda/internal/model"
)

type Validator interface {
	Valid() (problems map[string]string)
}

type Mapper[T any] interface {
	MapTo() (T, error)
}

type ValidatorMapper[T any] interface {
	Validator
	Mapper[T]
}

type inputUser struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Role      string `json:"role,omitempty"`
	UserID    int    `json:"user_id,omitempty"`
}

func (user inputUser) MapTo() (model.User, error) {
	return model.User{
		ID:        0,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		UserID:    uint(user.UserID),
	}, nil
}

func (user inputUser) Valid() map[string]string {
	problems := make(map[string]string)

	// validate UserID greater than 0
	if user.UserID < 1 {
		problems["UserID"] = "UserID must be more than 0"
	}

	// validate role is `Customer` or `Employee`
	if user.Role != "Customer" && user.Role != "Employee" {
		problems["Role"] = "Role must be 'Customer' or 'Employee'"
	}

	return problems
}

func decodeValidateBody[I ValidatorMapper[O], O any](requestBody string) (O, map[string]string, error) {
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
