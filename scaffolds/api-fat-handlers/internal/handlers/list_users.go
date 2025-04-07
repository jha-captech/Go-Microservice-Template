package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/captechconsulting/go-microservice-templates/api/internal/models"
	"github.com/go-chi/httplog/v2"
)

// HandleListUsers is a Handler that returns a list of all users.
//
// @Summary		List all users
// @Description	List all users
// @Tags		users
// @Accept		json
// @Produce		json
// @Success		200		{object}	handlers.responseUsers
// @Failure		500		{object}	handlers.responseErr
// @Router		/user	[GET]
func HandleListUsers(logger *httplog.Logger, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup
		ctx := r.Context()

		// get values from database
		rows, err := db.QueryContext(
			ctx,
			`SELECT * FROM "users"`,
		)
		if err != nil {
			logger.Error("failed to get users", slog.String("error", err.Error()))
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error retrieving data",
			})
			return
		}
		defer func() {
			_ = rows.Close()
		}()

		var users []models.User
		for rows.Next() {
			var user models.User
			err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Role, &user.UserID)
			if err != nil {
				logger.Error("failed to scan user from row", slog.String("error", err.Error()))
				encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
					Error: "Error retrieving data",
				})
				return
			}

			users = append(users, user)
		}

		if err = rows.Err(); err != nil {
			logger.Error("failed to scan users", slog.String("error", err.Error()))
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error retrieving data",
			})
			return
		}

		// return response
		usersOut := mapMultipleOutput(users)
		encodeResponse(w, logger, http.StatusOK, responseUsers{
			Users: usersOut,
		})
	}
}
