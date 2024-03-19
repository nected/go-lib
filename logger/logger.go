package logger

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type Logger struct {
	log           *zap.Logger
	mu            sync.Mutex
	sentryEnabled bool
	sentryDSN     string
	logSentry     bool
}

var logger *Logger

func NewLogger() *Logger {
	if logger == nil {
		logger.mu.Lock()
		defer logger.mu.Unlock()
		l, err := zap.NewProduction()
		if err != nil {
			fmt.Printf("failed to initialize zap logger: %v", err)
		}
		logger = &Logger{
			log: l,
		}
	}
	return logger
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
