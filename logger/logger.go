package logger

import (
	"fmt"

	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	log           *zap.Logger
	sentryEnabled bool
	sentryDSN     string
	logSentry     bool
}

func getOptions() []zap.Option {
	return []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.AddStacktrace(zap.ErrorLevel),
	}
}

func NewLogger() *Logger {
	config := zap.NewProductionConfig()
	l, err := config.Build(getOptions()...)
	if err != nil {
		fmt.Printf("failed to initialize zap logger: %v", err)
	}
	return &Logger{log: l}
}

func NewNamedLogger(name string) *Logger {
	l, err := zap.NewProduction(getOptions()...)
	if err != nil {
		fmt.Printf("failed to initialize zap logger: %v", err)
	}
	return &Logger{log: l.Named(name)}
}

func (l *Logger) WithEnv(env string) *Logger {
	var err error
	switch env {
	case "development", "dev":
		l.log, err = zap.NewDevelopment(getOptions()...)
		if err != nil {
			fmt.Printf("failed to initialize zap logger: %v", err)
		}
	case "local":
		l.log, err = zap.NewDevelopment(getOptions()...)
		if err != nil {
			fmt.Printf("failed to initialize zap logger: %v", err)
		}
	}
	l.log = l.log.With(zap.String("env", env))
	return l
}

func (l *Logger) AddCaller() *Logger {
	l.log = l.log.WithOptions(zap.AddCaller())
	return l
}

func (l *Logger) AddCallerSkip(skip int) *Logger {
	l.log = l.log.WithOptions(zap.AddCallerSkip(skip))
	return l
}

func (l *Logger) WithSentry(client *sentry.Client) *Logger {
	l.sentryEnabled = true
	l.logSentry = true
	l.log = modifyToSentryLogger(l.log, client)
	return l
}

// info log
func (l *Logger) Info(msg string, args ...interface{}) {
	l.log.With(getZapFields(args...)...).Info(msg)
}

// error log
func (l *Logger) Error(msg string, args ...interface{}) {
	l.log.With(getZapFields(args...)...).Error(msg)
}

// warn log
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log.With(getZapFields(args...)...).Warn(msg)
}

// debug log
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log.With(getZapFields(args...)...).Debug(msg)
}

// fatal log
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.log.With(getZapFields(args...)...).Fatal(msg)
}

// panic log
func (l *Logger) Panic(msg string, args ...interface{}) {
	l.log.With(getZapFields(args...)...).Panic(msg)
}

func getZapFields(args ...interface{}) (fields []zapcore.Field) {
	var fError error
	fields = make([]zap.Field, 0)
	errors := make([]error, 0)
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			break
		}
		key := args[i]
		value := args[i+1]

		keyStr, ok := key.(string)
		if !ok {
			keyStr = fmt.Sprintf("%v", key)
		}
		if err, ok := value.(error); ok {
			errors = append(errors, err)
			continue
		}
		fields = append(fields, zap.Any(keyStr, value))
	}

	if (len(fields)*2 + len(errors)) < len(args) {
		// check if value is of type error
		if err, ok := args[len(args)-1].(error); ok {
			fError = err
		}
	}
	if len(errors) > 0 {
		fields = append(fields, zap.Errors("errors", errors))
	}
	if fError != nil {
		fields = append(fields, zap.Error(fError))
	}
	return fields
}

func modifyToSentryLogger(log *zap.Logger, client *sentry.Client) *zap.Logger {
	cfg := zapsentry.Configuration{
		Level:             zapcore.ErrorLevel, // when to send message to sentry
		EnableBreadcrumbs: true,               // enable sending breadcrumbs to Sentry
		BreadcrumbLevel:   zapcore.InfoLevel,  // at what level should we sent breadcrumbs to sentry, this level can't be higher than `Level`
	}
	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromClient(client))
	// don't use value if error was returned. Noop core will be replaced to nil soon.
	if err != nil {
		core = zapcore.NewNopCore()
	}

	log = zapsentry.AttachCoreToLogger(core, log)

	// if you have web service, create a new scope somewhere in middleware to have valid breadcrumbs.
	return log.With(zapsentry.NewScope())
}
