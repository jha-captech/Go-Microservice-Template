package middleware

import (
	"context"
)

type (
	// HandlerFuncT is a lambda handler type where E is the event type and R is the response type.
	HandlerFuncT[E any, R any] func(context.Context, E) (R, error)

	// LambdaMiddlewareT is a lambda middleware type where E is the event type and R is the response type.
	LambdaMiddlewareT[E any, R any] func(next HandlerFuncT[E, R]) HandlerFuncT[E, R]
)

func AddToHandler[E any, R any](handler HandlerFuncT[E, R], middlewares ...LambdaMiddlewareT[E, R]) HandlerFuncT[E, R] {
	mds := func(next HandlerFuncT[E, R]) HandlerFuncT[E, R] {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}

		return next
	}

	return mds(handler)
}
