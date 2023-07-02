package blackbox

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

// Level indicates the logging level to be used when logging messages.
type Level int

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
