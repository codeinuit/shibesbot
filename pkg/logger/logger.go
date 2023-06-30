package logger

// Logger interface implementation
type Logger interface {
	Info(...any)
	Warn(...any)
	Error(...any)
}
