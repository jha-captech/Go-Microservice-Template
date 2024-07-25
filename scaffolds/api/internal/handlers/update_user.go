package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/captechconsulting/go-microservice-templates/api/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

type userUpdater interface {
	UpdateUser(ctx context.Context, ID int, user models.User) (models.User, error)
}

// HandleUpdateUser is a Handler that updates a user based on a user object from the request body.
//
// @Summary		Update a user by ID
// @Description	Update a user by ID
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		id			path		int	true						"User ID"
// @Param		user		body		handlers.inputUser		true	"User Object"
// @Success		200			{object}	handlers.responseUser
// @Failure		500			{object}	handlers.responseErr
// @Failure		422			{object}	handlers.responseErr
// @Router		/user/{ID}	[PUT]
func HandleUpdateUser(logger *httplog.Logger, service userUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup
		ctx := r.Context()

		// get and validate ID
		idString := chi.URLParam(r, "ID")
		ID, err := strconv.Atoi(idString)
		if err != nil {
			logger.Error("error getting ID", "error", err)
			encodeResponse(w, logger, http.StatusBadRequest, responseErr{
				Error: "Not a valid ID",
			})
			return
		}

		// get and validate body as object
		userIn, problems, err := decodeValidateBody[inputUser, models.User](r)
		if err != nil {
			switch {
			case len(problems) > 0:
				logger.Error("Problems validating input", "error", err, "problems", problems)
				encodeResponse(w, logger, http.StatusBadRequest, responseErr{
					ValidationErrors: problems,
				})
			default:
				logger.Error("BodyParser error", "error", err)
				encodeResponse(w, logger, http.StatusBadRequest, responseErr{
					Error: "missing values or malformed body",
				})
			}
			return
		}

		// update object in database
		user, err := service.UpdateUser(ctx, ID, userIn)
		if err != nil {
			logger.Error("error updating object in database", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error updating object",
			})
			return
		}

		// return response
		userOut := mapOutput(user)
		encodeResponse(w, logger, http.StatusOK, responseUser{
			User: userOut,
		})
	}
}
