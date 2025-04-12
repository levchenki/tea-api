package logx

type AppLogger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}
