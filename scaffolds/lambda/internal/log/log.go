package log

// Logger is a structured logger interface that is compatible with `log/slog` and
// `github.com/go-chi/httplog`.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}
