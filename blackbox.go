package blackbox

import (
	"fmt"
	"os"
)

// New creates a new blackbox logger
func New() *Logger {
	return &Logger{
		level:     Trace,
		targetSet: &targetSet{},
		context:   make(context, 0),
	}
}

//WithCtx creates a new blackbox logger with a given context
func WithCtx(contextData Ctx) *Logger {
	logger := New()
	logger.context = logger.context.extend(contextData)
	return logger
}

// Logger will take log messages and write them to the targets provided
type Logger struct {
	level     Level
	targetSet *targetSet
	context   context
}

// Log logs values to the loggers targets at the given log level. Any values
// that have a String method, it's return value will be used instead.
func (l *Logger) Log(level Level, values ...interface{}) *Logger {
	if level < l.level {
		return l
	}
	l.targetSet.log(level, values, l.context)
	return l
}

// Logf works the same way as fmt.Printf. Provide a format string and any
// values you wish.
func (l *Logger) Logf(level Level, format string, values ...interface{}) *Logger {
	l.Log(level, fmt.Sprintf(format, values...))
	return l
}

// Trace is a convenience method for logging values at the trace log level. It
// behaves the same as Log.
func (l *Logger) Trace(values ...interface{}) *Logger {
	l.Log(Trace, values...)
	return l
}

// Tracef is a convenience method for logging values at the trace log level. It
// behaves the same as Logf.
func (l *Logger) Tracef(format string, values ...interface{}) *Logger {
	l.Logf(Trace, format, values...)
	return l
}

// Debug is a convenience method for logging values at the debug log level. It
// behaves the same as Log.
func (l *Logger) Debug(values ...interface{}) *Logger {
	l.Log(Debug, values...)
	return l
}

// Debugf is a convenience method for logging values at the debug log level. It
// behaves the same as Logf.
func (l *Logger) Debugf(format string, values ...interface{}) *Logger {
	l.Logf(Debug, format, values...)
	return l
}

// Verbose is a convenience method for logging values at the verbose log level. It
// behaves the same as Log.
func (l *Logger) Verbose(values ...interface{}) *Logger {
	l.Log(Verbose, values...)
	return l
}

// Verbosef is a convenience method for logging values at the verbose log level. It
// behaves the same as Logf.
func (l *Logger) Verbosef(format string, values ...interface{}) *Logger {
	l.Logf(Verbose, format, values...)
	return l
}

// Info is a convenience method for logging values at the info log level. It
// behaves the same as Log.
func (l *Logger) Info(values ...interface{}) *Logger {
	l.Log(Info, values...)
	return l
}

// Infof is a convenience method for logging values at the info log level. It
// behaves the same as Logf.
func (l *Logger) Infof(format string, values ...interface{}) *Logger {
	l.Logf(Info, format, values...)
	return l
}

// Warn is a convenience method for logging values at the warn log level. It
// behaves the same as Log.
func (l *Logger) Warn(values ...interface{}) *Logger {
	l.Log(Warn, values...)
	return l
}

// Warnf is a convenience method for logging values at the warn log level. It
// behaves the same as Logf.
func (l *Logger) Warnf(format string, values ...interface{}) *Logger {
	l.Logf(Warn, format, values...)
	return l
}

// Error is a convenience method for logging values at the error log level. It
// behaves the same as Log.
func (l *Logger) Error(values ...interface{}) *Logger {
	l.Log(Error, values...)
	return l
}

// Errorf is a convenience method for logging values at the error log level. It
// behaves the same as Logf.
func (l *Logger) Errorf(format string, values ...interface{}) *Logger {
	l.Logf(Error, format, values...)
	return l
}

// Fatal is a convenience method for logging values at the fatal log level. It
// behaves the same as Log with the exception that it exits the program with
// code 1.
func (l *Logger) Fatal(values ...interface{}) {
	l.Log(Fatal, values...)
	os.Exit(1)
}

// Fatalf is a convenience method for logging values at the fatal log level. It
// behaves the same as Logf with the exception that it exits the program with
// code 1.
func (l *Logger) Fatalf(format string, values ...interface{}) {
	l.Logf(Fatal, format, values...)
	os.Exit(1)
}

// Panic is a convenience method for logging values at the panic log level. It
// behaves the same as Log.
func (l *Logger) Panic(values ...interface{}) {
	l.Log(Panic, values...)
	panic(fmt.Sprint(values...))
}

// Panicf is a convenience method for logging values at the panic log level. It
// behaves the same as Logf.
func (l *Logger) Panicf(format string, values ...interface{}) {
	l.Logf(Panic, format, values...)
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

// Ctx takes a context, merging it with the current one, and creates a new
// sub logger from the merged context. This new logger will have the same
// target set as the one Ctx is called upon.
func (l *Logger) Ctx(context Ctx) *Logger {
	return &Logger{
		context:   l.context.extend(context),
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

// Ctx is an alias for map[string]interface{}. This is the format for
// data to me used for extending contexts.
type Ctx map[string]interface{}

// Target is an interface ment to be implemented by types that collect log
// data. blackbox ships with two of these: PrettyTarget and JSONTarget
type Target interface {
	Log(level Level, values []interface{}, context context)
}

// LevelFromString returns a log level matching the given string
func LevelFromString(levelStr string) Level {
	var level Level
	switch levelStr {
	case "trace":
		level = Trace
	case "debug":
		level = Debug
	case "verbose":
		level = Verbose
	case "info":
		level = Info
	case "warn":
		level = Warn
	case "error":
		level = Error
	case "fatal":
		level = Fatal
	case "panic":
		level = Panic
	}
	return level
}

// Level indicates the logging level to be used when logging messages.
type Level int

const (
	// Trace log level
	Trace Level = iota
	// Debug log level
	Debug
	// Verbose log level
	Verbose
	// Info log level
	Info
	// Warn log level
	Warn
	// Error log level
	Error
	// Fatal log level
	Fatal
	// Panic log level
	Panic
)

// String returns the string representation of each Level constant
func (l Level) String() string {
	var levelStr string
	switch l {
	case Trace:
		levelStr = "trace"
	case Debug:
		levelStr = "debug"
	case Verbose:
		levelStr = "verbose"
	case Info:
		levelStr = "info"
	case Warn:
		levelStr = "warn"
	case Error:
		levelStr = "error"
	case Fatal:
		levelStr = "fatal"
	case Panic:
		levelStr = "panic"
	}
	return levelStr
}
