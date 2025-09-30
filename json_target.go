package blackbox

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// JSONTarget is a Target that produces newline separated json output containing
// log data.
type JSONTarget struct {
	showLoggerID  bool
	showTimestamp bool
	showLevel     bool
	showContext   bool
	useSource     bool
	level         Level
	outTarget     io.Writer
	errTarget     io.Writer
}

var _ Target = &JSONTarget{}

// NewJSONTarget creates a JSONTarget for use with a logger
func NewJSONTarget(outTarget io.Writer, errTarget io.Writer) *JSONTarget {
	return &JSONTarget{
		showTimestamp: true,
		showLevel:     true,
		showContext:   true,
		level:         Trace,
		outTarget:     outTarget,
		errTarget:     errTarget,
	}
}

// SetLevel sets the minimum log level that JSONTarget will output. Note that
// this setting is independent of the log level set on the logger itself.
func (j *JSONTarget) SetLevel(level Level) *JSONTarget {
	j.level = level
	return j
}

// ShowTimestamp will enable or disable timestamps in the output depending on
// the boolean value passed.
func (j *JSONTarget) ShowTimestamp(b bool) *JSONTarget {
	j.showTimestamp = b
	return j
}

// ShowLevel will enable or disable level values in the output depending on
// the boolean value passed.
func (j *JSONTarget) ShowLevel(b bool) *JSONTarget {
	j.showLevel = b
	return j
}

// ShowContext will enable or disable context key value pairs in the output
// depending on the boolean value passed.
func (j *JSONTarget) ShowContext(b bool) *JSONTarget {
	j.showContext = b
	return j
}

// UseSource enables the inclusion of source
func (s *JSONTarget) UseSource(b bool) *JSONTarget {
	s.useSource = b
	return s
}

// Log takes a Level and series of values, then outputs them formatted
// accordingly.
func (j *JSONTarget) Log(loggerID string, level Level, values []any, context Ctx, getSource func() *Source) {
	if level < j.level {
		return
	}

	jsonData := make(map[string]any, 1)
	if j.showTimestamp {
		jsonData["time"] = time.Now().Local().Format(time.RFC3339)
	}
	if j.showLevel {
		jsonData["level"] = level.String()
	}
	strValues := make([]string, 0)
	for _, value := range values {
		strValues = append(strValues, fmt.Sprintf("%+v", value))
	}
	jsonData["message"] = strings.Join(strValues, " ")
	if j.showContext {
		jsonData["context"] = context
	}
	if j.showLoggerID {
		jsonData["loggerID"] = loggerID
	}
	if j.useSource {
		jsonData["source"] = getSource()
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}

	jsonBytes = append(jsonBytes, byte('\n'))

	if level >= Warn {
		_, err = j.errTarget.Write(jsonBytes)
	} else {
		_, err = j.outTarget.Write(jsonBytes)
	}
	if err != nil {
		panic(err)
	}
}
