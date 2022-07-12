package logger

type Logger interface {
	Info(funcName string, msg string, attempt uint16)
	Warn(funcName string, msg string, attempt uint16)
	Error(funcName string, msg string, attempt uint16)
}
