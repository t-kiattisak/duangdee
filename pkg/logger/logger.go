package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger is a global structured logger wrapper
type Logger struct {
	logger zerolog.Logger
}

// New creates a new structured JSON logger for a specific microservice
func New(serviceName string) *Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	l := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("service", serviceName).
		Logger()

	return &Logger{logger: l}
}

// Info logs an informational message with key-value fields
func (l *Logger) Info(msg string, fields map[string]interface{}) {
	event := l.logger.Info()
	for k, v := range fields {
		event.Interface(k, v)
	}
	event.Msg(msg)
}

// Error logs an error message with traceback/fields
func (l *Logger) Error(err error, msg string, fields map[string]interface{}) {
	event := l.logger.Error().Err(err)
	for k, v := range fields {
		event.Interface(k, v)
	}
	event.Msg(msg)
}

// Debug logs debug messages
func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	event := l.logger.Debug()
	for k, v := range fields {
		event.Interface(k, v)
	}
	event.Msg(msg)
}

// WithTraceID adds a trace_id to logger context for correlation in Kibana
func (l *Logger) WithTraceID(traceID string) zerolog.Logger {
	return l.logger.With().Str("trace_id", traceID).Logger()
}
