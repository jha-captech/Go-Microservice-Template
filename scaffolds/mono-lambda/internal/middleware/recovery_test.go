package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestRecovery(t *testing.T) {
	tests := map[string]struct {
		handler       HandlerFunc
		expectPanic   bool
		expectedLog   string
		expectedEvent events.APIGatewayProxyResponse
	}{
		"handler does not panic": {
			handler: func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
				return events.APIGatewayProxyResponse{
					StatusCode: 200,
					Body:       "OK",
				}, nil
			},
			expectPanic: false,
			expectedEvent: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       "OK",
			},
		},
		"handler panics": {
			handler: func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
				panic("something went wrong")
			},
			expectPanic: true,
			expectedLog: "Recovered from panic",
			expectedEvent: events.APIGatewayProxyResponse{
				Headers:    map[string]string{"Content-Type": "application/json"},
				StatusCode: http.StatusInternalServerError,
				Body:       `{"error": "Internal server error"}`,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			logger := slog.Default()

			recoveryMiddleware := Recovery(logger)
			handlerWithRecovery := recoveryMiddleware(tt.handler)

			resp, err := handlerWithRecovery(context.Background(), events.APIGatewayProxyRequest{})

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedEvent, resp)
		})
	}
}
