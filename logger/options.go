package logger

type LoggerOptions interface {
	apply(*Logger)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*Logger)

func (f optionFunc) apply(log *Logger) {
	f(log)
}

func WithSentryDSN(dsn string) LoggerOptions {
	return optionFunc(func(log *Logger) {
		log.sentryDSN = dsn
		log.sentryEnabled = true
	})
}

func WithSentryLog(logSentry bool) LoggerOptions {
	return optionFunc(func(log *Logger) {
		log.logSentry = logSentry
	})
}
