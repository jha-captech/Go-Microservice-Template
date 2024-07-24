package handlers

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

// HandlerFunc is an alias for a lambda function that is used with API Gateway.
type HandlerFunc = func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
