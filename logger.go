package blackbox

import (
	"fmt"
	"os"
	"runtime"
)

// Logger will take log messages and write them to the targets provided
type Logger struct {
	level     Level
	targetSet *targetSet
	context   Ctx
}

// New creates a new blackbox logger
func New() *Logger {
	return &Logger{
		level:     Trace,
		targetSet: &targetSet{},
		context:   make(Ctx, 0),
	}
}

// NewWithCtx creates a new blackbox logger with a given context
func NewWithCtx(contextData Ctx) *Logger {
	logger := New()
	return logger.WithCtx(contextData)
}

// Log logs values to the loggers targets at the given log level. Any values
// that have a String method, it's return value will be used instead.
func (l *Logger) Log(level Level, values ...any) *Logger {
	l.log(level, values...)
	return l
}

// Logf works the same way as fmt.Printf. Provide a format string and any
// values you wish.
func (l *Logger) Logf(level Level, format string, values ...any) *Logger {
	l.log(level, fmt.Sprintf(format, values...))
	return l
}

// Trace is a convenience method for logging values at the trace log level. It
// behaves the same as Log.
func (l *Logger) Trace(values ...any) *Logger {
	l.log(Trace, values...)
	return l
}

// Tracef is a convenience method for logging values at the trace log level. It
// behaves the same as Logf.
func (l *Logger) Tracef(format string, values ...any) *Logger {
	l.log(Trace, fmt.Sprintf(format, values...))
	return l
}

// Debug is a convenience method for logging values at the debug log level. It
// behaves the same as Log.
func (l *Logger) Debug(values ...any) *Logger {
	l.log(Debug, values...)
	return l
}

// Debugf is a convenience method for logging values at the debug log level. It
// behaves the same as Logf.
func (l *Logger) Debugf(format string, values ...any) *Logger {
	l.log(Debug, fmt.Sprintf(format, values...))
	return l
}

// Verbose is a convenience method for logging values at the verbose log level. It
// behaves the same as Log.
func (l *Logger) Verbose(values ...any) *Logger {
	l.log(Verbose, values...)
	return l
}

// Verbosef is a convenience method for logging values at the verbose log level. It
// behaves the same as Logf.
func (l *Logger) Verbosef(format string, values ...any) *Logger {
	l.log(Verbose, fmt.Sprintf(format, values...))
	return l
}

// Info is a convenience method for logging values at the info log level. It
// behaves the same as Log.
func (l *Logger) Info(values ...any) *Logger {
	l.log(Info, values...)
	return l
}

// Infof is a convenience method for logging values at the info log level. It
// behaves the same as Logf.
func (l *Logger) Infof(format string, values ...any) *Logger {
	l.log(Info, fmt.Sprintf(format, values...))
	return l
}

// Warn is a convenience method for logging values at the warn log level. It
// behaves the same as Log.
func (l *Logger) Warn(values ...any) *Logger {
	l.log(Warn, values...)
	return l
}

// Warnf is a convenience method for logging values at the warn log level. It
// behaves the same as Logf.
func (l *Logger) Warnf(format string, values ...any) *Logger {
	l.log(Warn, fmt.Sprintf(format, values...))
	return l
}

// Error is a convenience method for logging values at the error log level. It
// behaves the same as Log.
func (l *Logger) Error(values ...any) *Logger {
	l.log(Error, values...)
	return l
}

// Errorf is a convenience method for logging values at the error log level. It
// behaves the same as Logf.
func (l *Logger) Errorf(format string, values ...any) *Logger {
	l.log(Error, fmt.Errorf(format, values...))
	return l
}

// Fatal is a convenience method for logging values at the fatal log level. It
// behaves the same as Log with the exception that it exits the program with
// code 1.
func (l *Logger) Fatal(values ...any) {
	l.log(Fatal, values...)
	os.Exit(1)
}

// Fatalf is a convenience method for logging values at the fatal log level. It
// behaves the same as Logf with the exception that it exits the program with
// code 1.
func (l *Logger) Fatalf(format string, values ...any) {
	l.log(Fatal, fmt.Sprintf(format, values...))
	os.Exit(1)
}

// Panic is a convenience method for logging values at the panic log level. It
// behaves the same as Log.
func (l *Logger) Panic(values ...any) {
	l.log(Panic, values...)
	panic(fmt.Sprint(values...))
}

// Panicf is a convenience method for logging values at the panic log level. It
// behaves the same as Logf.
func (l *Logger) Panicf(format string, values ...any) {
	l.log(Panic, fmt.Sprintf(format, values...))
	panic(fmt.Sprint(values...))
}

// SetLevel sets the log level across all targets at once.
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// AddTarget adds a io.Writer to be written to
func (l *Logger) AddTarget(target Target) {
	l.targetSet.addTarget(target)
}

// WithCtx takes a context, merging it with the current one, and creates a new
// sub logger from the merged context. This new logger will have the same
// target set as the one WithCtx is called upon.
func (l *Logger) WithCtx(context Ctx) *Logger {
	return &Logger{
		level:     l.level,
		context:   l.context.Extend(context),
		targetSet: l.targetSet,
	}
}

// GetCtx returns a ctx instance containing a copy of the logger's internal
// context data.
func (l *Logger) GetCtx() Ctx {
	ctx := make(Ctx, 0)
	for key, value := range l.context {
		ctx[key] = value
	}
	return ctx
}

func (l *Logger) log(level Level, values ...any) {
	if level < l.level {
		return
	}
	pcs := make([]uintptr, 64)
	n := runtime.Callers(2, pcs)
	l.targetSet.log(level, values, l.context, pcs[:n])
}
