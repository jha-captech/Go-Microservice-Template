package handlers

import (
	"net/http"

	"github.com/go-chi/httplog/v2"
)

// HandleHealth is a health check handler
//
// @Summary		Health check response
// @Description	Health check response
// @Tags		health-check
// @Accept		json
// @Produce		json
// @Success		200				{object}	handlers.responseMsg
// @Router		/health-check	[GET]
func HandleHealth(logger *httplog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Health check called")
		encodeResponse(w, logger, http.StatusOK, responseMsg{
			Message: "hello world",
		})
	}
}
