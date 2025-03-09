package utils

import (
	"log"
	"os"
)

// ANSI color codes.
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorGreen  = "\033[32m"
)

// Level defines the log level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Logger defines a logging interface with context support.
type Logger interface {
	// WithContext returns a new logger instance with the provided context.
	WithContext(ctx string) Logger
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg string)
}

// DefaultLogger is a basic logger using the standard log package.
type DefaultLogger struct {
	level   Level
	context string
}

// New returns a new instance of the default logger.
// It sets the logging level based on the LOG_LEVEL environment variable.
func New() Logger {
	level := LevelInfo
	if lvl, ok := os.LookupEnv("LOG_LEVEL"); ok {
		switch lvl {
		case "DEBUG":
			level = LevelDebug
		case "WARN":
			level = LevelWarn
		case "ERROR":
			level = LevelError
		default:
			level = LevelInfo
		}
	}
	return &DefaultLogger{
		level:   level,
		context: "",
	}
}

// WithContext returns a new logger instance with the provided context.
func (l *DefaultLogger) WithContext(ctx string) Logger {
	return &DefaultLogger{
		level:   l.level,
		context: ctx,
	}
}

// formatMsg formats the log message with the context and appropriate color.
func (l *DefaultLogger) formatMsg(prefix string, msg string) string {
	var color string
	switch prefix {
	case "INFO":
		color = ColorBlue
	case "WARN":
		color = ColorYellow
	case "ERROR":
		color = ColorRed
	case "DEBUG":
		color = ColorGreen
	default:
		color = ""
	}

	if l.context != "" {
		return color + prefix + " [" + l.context + "]: " + msg + ColorReset
	}
	return color + prefix + ": " + msg + ColorReset
}

// Info logs informational messages.
func (l *DefaultLogger) Info(msg string) {
	if l.level <= LevelInfo {
		log.Println(l.formatMsg("INFO", msg))
	}
}

// Warn logs warning messages.
func (l *DefaultLogger) Warn(msg string) {
	if l.level <= LevelWarn {
		log.Println(l.formatMsg("WARN", msg))
	}
}

// Error logs error messages.
func (l *DefaultLogger) Error(msg string) {
	if l.level <= LevelError {
		log.Println(l.formatMsg("ERROR", msg))
	}
}

// Debug logs debug messages.
func (l *DefaultLogger) Debug(msg string) {
	if l.level <= LevelDebug {
		log.Println(l.formatMsg("DEBUG", msg))
	}
}
