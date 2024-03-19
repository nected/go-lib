package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger struct {
	log           *zap.Logger
	sentryEnabled bool
	sentryDSN     string
	logSentry     bool
}

func NewLogger() *Logger {
	l, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("failed to initialize zap logger: %v", err)
	}
	return &Logger{log: l}
}

func NewNamedLogger(name string) *Logger {
	l, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("failed to initialize zap logger: %v", err)
	}
	return &Logger{log: l.Named(name)}
}

// info log
func (l *Logger) Info(code string, msg string, fields ...zap.Field) {
	l.log.Info(msg, fields...)
}

// error log
func (l *Logger) Error(code string, msg string, fields ...zap.Field) {
	l.log.Error(msg, fields...)
}

// warn log
func (l *Logger) Warn(code string, msg string, fields ...zap.Field) {
	l.log.Warn(msg, fields...)
}

// debug log
func (l *Logger) Debug(code string, msg string, fields ...zap.Field) {
	l.log.Debug(msg, fields...)
}

// fatal log
func (l *Logger) Fatal(code string, msg string, fields ...zap.Field) {
	l.log.Fatal(msg, fields...)
}

// panic log
func (l *Logger) Panic(code string, msg string, fields ...zap.Field) {
	l.log.Panic(msg, fields...)
}
