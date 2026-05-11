package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
	logger *zap.Logger
}

// New creates a new logger instance with the specified log level
func New(level string, isDev bool) (*Logger, error) {
	var config zap.Config

	if isDev {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
	}

	// Parse and set log level
	zapLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}
	config.Level = zap.NewAtomicLevelAt(zapLevel)

	// Build logger
	zapLogger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return &Logger{
		SugaredLogger: zapLogger.Sugar(),
		logger:        zapLogger,
	}, nil
}

// NewNoop creates a no-op logger for testing
func NewNoop() *Logger {
	zapLogger := zap.NewNop()
	return &Logger{
		SugaredLogger: zapLogger.Sugar(),
		logger:        zapLogger,
	}
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.logger.Sync()
}

// WithField adds a structured field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		SugaredLogger: l.With(key, value),
		logger:        l.logger,
	}
}

// WithFields adds multiple structured fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	var args []interface{}
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{
		SugaredLogger: l.With(args...),
		logger:        l.logger,
	}
}

// Debugf logs a debug message with format
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.SugaredLogger.Debugf(format, args...)
}

// Infof logs an info message with format
func (l *Logger) Infof(format string, args ...interface{}) {
	l.SugaredLogger.Infof(format, args...)
}

// Warnf logs a warning message with format
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.SugaredLogger.Warnf(format, args...)
}

// Errorf logs an error message with format
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.SugaredLogger.Errorf(format, args...)
}

// Fatalf logs a fatal message and exits
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.SugaredLogger.Fatalf(format, args...)
}

// Debug logs a debug message
func (l *Logger) Debug(args ...interface{}) {
	l.SugaredLogger.Debug(args...)
}

// Info logs an info message
func (l *Logger) Info(args ...interface{}) {
	l.SugaredLogger.Info(args...)
}

// Warn logs a warning message
func (l *Logger) Warn(args ...interface{}) {
	l.SugaredLogger.Warn(args...)
}

// Error logs an error message
func (l *Logger) Error(args ...interface{}) {
	l.SugaredLogger.Error(args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(args ...interface{}) {
	l.SugaredLogger.Fatal(args...)
}

// ContextLogger adds request correlation tracking
type ContextLogger struct {
	*Logger
	correlationID string
}

// WithCorrelationID creates a context logger with a correlation ID
func (l *Logger) WithCorrelationID(correlationID string) *ContextLogger {
	return &ContextLogger{
		Logger:        l.WithField("correlation_id", correlationID),
		correlationID: correlationID,
	}
}

// GetCorrelationID returns the correlation ID
func (cl *ContextLogger) GetCorrelationID() string {
	return cl.correlationID
}
