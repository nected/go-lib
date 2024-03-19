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
	l, err := zap.NewProduction(zap.WithCaller(true))
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
	l.log.Info(msg, getZapFields(code, fields...)...)
}

// error log
func (l *Logger) Error(code string, msg string, fields ...zap.Field) {
	l.log.Error(msg, getZapFields(code, fields...)...)
}

// warn log
func (l *Logger) Warn(code string, msg string, fields ...zap.Field) {
	l.log.Warn(msg, getZapFields(code, fields...)...)
}

// debug log
func (l *Logger) Debug(code string, msg string, fields ...zap.Field) {
	l.log.Debug(msg, getZapFields(code, fields...)...)
}

// fatal log
func (l *Logger) Fatal(code string, msg string, fields ...zap.Field) {
	l.log.Fatal(msg, getZapFields(code, fields...)...)
}

// panic log
func (l *Logger) Panic(code string, msg string, fields ...zap.Field) {
	l.log.Panic(msg, getZapFields(code, fields...)...)
}

func getZapFields(code string, fields ...zap.Field) []zap.Field {
	return append(fields, zap.String("code", code))
}
