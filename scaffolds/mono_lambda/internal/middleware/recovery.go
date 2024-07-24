package middleware

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/captechconsulting/go-microservice-templates/lambda/internal/log"
)

func Recovery(logger log.Logger) LambdaMiddleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, request events.APIGatewayProxyRequest) (event events.APIGatewayProxyResponse, err error) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Recovered from panic", "err", err)
					event = events.APIGatewayProxyResponse{
						Headers:    map[string]string{"Content-Type": "application/json"},
						StatusCode: http.StatusInternalServerError,
						Body:       `{"error": "Internal server error"}`,
					}
				}
			}()

			return next(ctx, request)
		}
	}
}
