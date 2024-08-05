package middleware

import (
	"context"
	"log/slog"
)

func Recovery[E any, R any](logger *slog.Logger) LambdaMiddlewareT[E, R] {
	return func(next HandlerFuncT[E, R]) HandlerFuncT[E, R] {
		return func(ctx context.Context, event E) (response R, err error) {
			defer func() {
				if errAny := recover(); errAny != nil {
					logger.Error("Recovered from panic", "err", errAny)
				}
			}()

			return next(ctx, event)
		}
	}
}

func RecoveryReturn[E any, R any](logger *slog.Logger, returnFn func() R) LambdaMiddlewareT[E, R] {
	return func(next HandlerFuncT[E, R]) HandlerFuncT[E, R] {
		return func(ctx context.Context, event E) (response R, err error) {
			defer func() {
				if errAny := recover(); errAny != nil {
					logger.Error("Recovered from panic", "err", errAny)

					response = returnFn()
				}
			}()

			return next(ctx, event)
		}
	}
}
