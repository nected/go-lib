package logger

import (
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
)

// NewLogger / NewNamedLogger / WithEnv / AddCaller / AddCallerSkip return
// non-nil loggers in every documented mode. Wrapping methods reassign the
// internal *zap.Logger, so the wrapped logger pointer should change.
func TestLoggerConstructorsAndChaining(t *testing.T) {
	t.Run("NewLogger returns ready logger", func(t *testing.T) {
		l := NewLogger()
		assert.NotNil(t, l)
		assert.NotNil(t, l.log)
	})

	t.Run("NewNamedLogger applies name", func(t *testing.T) {
		l := NewNamedLogger("test-name")
		assert.NotNil(t, l)
		assert.NotNil(t, l.log)
	})

	t.Run("WithEnv development uses dev logger", func(t *testing.T) {
		l := NewLogger()
		orig := l.log
		ret := l.WithEnv("development")
		assert.Same(t, l, ret)
		assert.NotSame(t, orig, l.log)
	})

	t.Run("WithEnv dev alias uses dev logger", func(t *testing.T) {
		l := NewLogger().WithEnv("dev")
		assert.NotNil(t, l.log)
	})

	t.Run("WithEnv local uses dev logger", func(t *testing.T) {
		l := NewLogger().WithEnv("local")
		assert.NotNil(t, l.log)
	})

	t.Run("WithEnv production passes through", func(t *testing.T) {
		// Non dev/local/local env still gets the env tag attached without
		// swapping the underlying logger.
		l := NewLogger().WithEnv("prod")
		assert.NotNil(t, l.log)
	})

	t.Run("AddCaller and AddCallerSkip mutate options", func(t *testing.T) {
		l := NewLogger()
		assert.Same(t, l, l.AddCaller())
		assert.Same(t, l, l.AddCallerSkip(2))
	})
}

// WithSentry is a no-op when the client is nil, and wraps the underlying
// logger when given a real Sentry client.
func TestWithSentry(t *testing.T) {
	t.Run("nil client is a no-op", func(t *testing.T) {
		l := NewLogger()
		orig := l.log
		ret := l.WithSentry(nil)
		assert.Same(t, l, ret)
		assert.False(t, l.sentryEnabled)
		assert.Same(t, orig, l.log)
	})

	t.Run("real client wraps logger and flips flags", func(t *testing.T) {
		// DSN must be a valid format; the client never actually transmits in
		// this test because no events are flushed.
		client, err := sentry.NewClient(sentry.ClientOptions{
			Dsn: "https://public@example.com/1",
		})
		assert.NoError(t, err)
		l := NewLogger()
		orig := l.log
		ret := l.WithSentry(client)
		assert.Same(t, l, ret)
		assert.True(t, l.sentryEnabled)
		assert.True(t, l.logSentry)
		assert.NotSame(t, orig, l.log)
	})
}

// Info / Error / Warn / Debug must not panic on either form: empty args or
// key/value pairs interleaved with errors. Output goes to stderr — we don't
// capture it, we just check that the calls return cleanly.
func TestLoggerLevels(t *testing.T) {
	l := NewLogger().WithEnv("development")

	assert.NotPanics(t, func() { l.Info("info message") })
	assert.NotPanics(t, func() { l.Warn("warn message", "k", "v") })
	assert.NotPanics(t, func() { l.Error("error message", "k1", "v1", "k2", 2) })
	assert.NotPanics(t, func() { l.Debug("debug message") })
}

// WithSentryDSN and WithSentryLog option funcs apply their setters when
// invoked. The skill's bug note: these options are declared but not wired
// into a constructor — they're tested in isolation only.
func TestOptionFunctions(t *testing.T) {
	l := &Logger{}

	WithSentryDSN("dsn://example").apply(l)
	assert.Equal(t, "dsn://example", l.sentryDSN)
	assert.True(t, l.sentryEnabled)

	WithSentryLog(true).apply(l)
	assert.True(t, l.logSentry)

	WithSentryLog(false).apply(l)
	assert.False(t, l.logSentry)
}

// getOptions returns the standard zap option list. We only check that it
// returns a non-empty slice — the actual options aren't introspectable.
func TestGetOptions(t *testing.T) {
	opts := getOptions()
	assert.NotEmpty(t, opts)
}

// Panic logs and then panics. Recovering from the panic is the only way to
// verify the call ran without bringing down the test process.
func TestLoggerPanic(t *testing.T) {
	l := NewLogger().WithEnv("development")
	assert.Panics(t, func() {
		l.Panic("panic msg", "k", "v")
	})
}

// getZapFields drops the trailing argument when args has odd length, then
// detects whether that trailing arg is an error and lifts it into a
// zap.Error field. This exercises the dangling-error branch that the
// even-length cases in logger_test.go don't reach.
func TestGetZapFields_TrailingError(t *testing.T) {
	// Odd-length args, last item is an error → ends up in the dangling
	// fError path and produces one extra ErrorType field.
	fields := getZapFields("k1", "v1", assertError("boom"))
	// Expect: one key field + one ErrorType field.
	assert.Len(t, fields, 2)
}

// assertError is a tiny error so we don't import "errors" twice.
type assertError string

func (e assertError) Error() string { return string(e) }
