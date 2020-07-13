package logging

type LogLevel int

const (
	_ LogLevel = iota
	LogLevelError
	LogLevelWarn
	LogLevelDebug
	LogLevelTrace
)
