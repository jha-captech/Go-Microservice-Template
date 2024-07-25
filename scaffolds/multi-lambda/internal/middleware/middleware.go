package middleware

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type (
	HandlerFunc      func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	LambdaMiddleware func(next HandlerFunc) HandlerFunc
)

func AddToHandler(handler HandlerFunc, middlewares ...LambdaMiddleware) HandlerFunc {
	mds := func(next HandlerFunc) HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}

		return next
	}

	return mds(handler)
}
