package middleware

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func mockHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "handler executed",
	}, nil
}

func mockMiddlewareOne(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		resp, err := next(ctx, req)
		resp.Body += " - middleware one"
		return resp, err
	}
}

func mockMiddlewareTwo(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		resp, err := next(ctx, req)
		resp.Body += " - middleware two"
		return resp, err
	}
}

func TestAddToHandler(t *testing.T) {
	tests := map[string]struct {
		handler     HandlerFunc
		middlewares []LambdaMiddleware
		wantBody    string
	}{
		"no middleware": {
			handler:     mockHandler,
			middlewares: nil,
			wantBody:    "handler executed",
		},
		"one middleware": {
			handler:     mockHandler,
			middlewares: []LambdaMiddleware{mockMiddlewareOne},
			wantBody:    "handler executed - middleware one",
		},
		"two middlewares": {
			handler:     mockHandler,
			middlewares: []LambdaMiddleware{mockMiddlewareOne, mockMiddlewareTwo},
			wantBody:    "handler executed - middleware two - middleware one",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			combinedHandler := AddToHandler(tt.handler, tt.middlewares...)
			resp, err := combinedHandler(context.Background(), events.APIGatewayProxyRequest{})
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBody, resp.Body)
		})
	}
}
