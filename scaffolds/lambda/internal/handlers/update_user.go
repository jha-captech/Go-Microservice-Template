package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/captechconsulting/go-microservice-templates/lambda/internal/log"
	"github.com/captechconsulting/go-microservice-templates/lambda/internal/model"
)

type userUpdater interface {
	UpdateUser(ctx context.Context, ID int, user model.User) (model.User, error)
}

func HandleUpdateUser(logger log.Logger, service userUpdater) HandlerFunc {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// get and validate ID
		idString := request.PathParameters["ID"]
		ID, err := strconv.Atoi(idString)
		if err != nil {
			logger.Error("error getting ID", "error", err)
			return encodeResponse(logger, http.StatusBadRequest, responseErr{
				Error: "Not a valid ID",
			})
		}

		// get and validate body as object
		userIn, problems, err := decodeValidateBody[inputUser, model.User](request.Body)
		if err != nil {
			switch {
			case len(problems) > 0:
				logger.Error("Problems validating input", "error", err, "problems", problems)
				return encodeResponse(logger, http.StatusBadRequest, responseErr{
					ValidationErrors: problems,
				})
			default:
				logger.Error("BodyParser error", "error", err)
				return encodeResponse(logger, http.StatusBadRequest, responseErr{
					Error: "missing values or malformed body",
				})
			}
		}

		// update object in database
		user, err := service.UpdateUser(ctx, ID, userIn)
		if err != nil {
			logger.Error("error updating object in database", "error", err)
			return encodeResponse(logger, http.StatusInternalServerError, responseErr{
				Error: "Error updating object",
			})
		}

		// return response
		userOut := mapOutput(user)
		return encodeResponse(logger, http.StatusOK, responseUser{
			User: userOut,
		})
	}
}
