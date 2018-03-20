package blackbox

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

// NewPrettyTarget creates a PrettyTarget for use with a logger
func NewPrettyTarget(outTarget io.Writer, errTarget io.Writer) *PrettyTarget {
	return &PrettyTarget{
		outTarget:     outTarget,
		errTarget:     errTarget,
		level:         Trace,
		showTimestamp: true,
		showLevel:     true,
		showContext:   true,
		useColor:      true,
	}
}

// PrettyTarget is a Target that produces newline separated human readable
// output suitable for stdout and stderr. It also supports colorized log levels.
type PrettyTarget struct {
	outTarget     io.Writer
	errTarget     io.Writer
	level         Level
	showTimestamp bool
	showLevel     bool
	showContext   bool
	useColor      bool
}

// SetLevel sets the minimum log level that PrettyTarget will output. Note that
// this setting is independent of the log level set on the logger itself.
func (s *PrettyTarget) SetLevel(level Level) *PrettyTarget {
	s.level = level
	return s
}

// ShowTimestamp will enable or disable timestamps in the output depending on
// the boolean value passed.
func (s *PrettyTarget) ShowTimestamp(b bool) *PrettyTarget {
	s.showTimestamp = b
	return s
}

// ShowLevel will enable or disable level labels in the output depending on
// the boolean value passed.
func (s *PrettyTarget) ShowLevel(b bool) *PrettyTarget {
	s.showLevel = b
	return s
}

// ShowContext will enable or disable context key value pairs in the output
// depending on the boolean value passed.
func (s *PrettyTarget) ShowContext(b bool) *PrettyTarget {
	s.showContext = b
	return s
}

// UseColor will enable or disable the use of ansi color codes in the output
// depending on the boolean value passed.
func (s *PrettyTarget) UseColor(b bool) *PrettyTarget {
	s.useColor = b
	return s
}

// Log takes a Level and series of values, then outputs them formatted
// accordingly.
func (s *PrettyTarget) Log(level Level, values []interface{}, context map[string]string) {
	if level < s.level {
		return
	}
	if s.showTimestamp {
		s.writeCurrentTimestamp(level)
	}
	if s.showLevel {
		s.writeLevel(level)
	}
	s.writeValues(level, values)
	if s.showContext {
		s.writeContext(level, context)
	}
	s.writeNewline(level)
}

func (s *PrettyTarget) writeCurrentTimestamp(level Level) {
	timestampBytes := []byte(time.Now().Local().Format(time.RFC3339) + " ")
	if level >= Warn {
		s.errTarget.Write(timestampBytes)
	} else {
		s.outTarget.Write(timestampBytes)
	}
}

func (s *PrettyTarget) writeLevel(level Level) {
	levelStr := level.String()

	var padStr string
	for i := len(levelStr); i < 7; i++ {
		padStr += " "
	}

	if s.useColor {
		switch level {
		case Trace:
			levelStr = "\u001b[35m" + levelStr + "\u001b[0m"
		case Debug:
			levelStr = "\u001b[34m" + levelStr + "\u001b[0m"
		case Verbose:
			levelStr = "\u001b[36m" + levelStr + "\u001b[0m"
		case Info:
			levelStr = "\u001b[32m" + levelStr + "\u001b[0m"
		case Warn:
			levelStr = "\u001b[33m" + levelStr + "\u001b[0m"
		case Error:
			levelStr = "\u001b[31m" + levelStr + "\u001b[0m"
		case Fatal:
			levelStr = "\u001b[37m\u001b[41;1m" + levelStr + "\u001b[0m"
		case Panic:
			levelStr = "\u001b[37m\u001b[45;1m" + levelStr + "\u001b[0m"
		}
	}

	levelBytes := []byte(levelStr + padStr + " ")
	if level >= Warn {
		s.errTarget.Write(levelBytes)
	} else {
		s.outTarget.Write(levelBytes)
	}
}

func (s *PrettyTarget) writeValues(level Level, values []interface{}) {
	valueStrs := make([]string, 0)
	for _, value := range values {
		valueStrs = append(valueStrs, fmt.Sprintf("%+v", value))
	}
	messageBytes := []byte(strings.Join(valueStrs, " "))
	if level >= Warn {
		s.errTarget.Write(messageBytes)
	} else {
		s.outTarget.Write(messageBytes)
	}
}

func (s *PrettyTarget) writeContext(level Level, context map[string]string) {
	contextStrs := make([]string, 0)
	for key, value := range context {
		contextStrs = append(contextStrs, key+"="+value)
	}
	sort.Strings(contextStrs)
	contextStr := strings.Join(contextStrs, " ")

	contextBytes := []byte(" " + contextStr)

	if level >= Warn {
		s.errTarget.Write(contextBytes)
	} else {
		s.outTarget.Write(contextBytes)
	}
}

func (s *PrettyTarget) writeNewline(level Level) {
	newLineBytes := []byte("\n")
	if level >= Warn {
		s.errTarget.Write(newLineBytes)
	} else {
		s.outTarget.Write(newLineBytes)
	}
}
