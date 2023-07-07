package blackbox

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

// NewJSONTarget creates a JSONTarget for use with a logger
func NewJSONTarget(outTarget io.Writer, errTarget io.Writer) *JSONTarget {
	return &JSONTarget{
		outTarget:     outTarget,
		errTarget:     errTarget,
		level:         Trace,
		showTimestamp: true,
		showLevel:     true,
		showContext:   true,
	}
}

// JSONTarget is a Target that produces newline separated json output containing
// log data.
type JSONTarget struct {
	outTarget     io.Writer
	errTarget     io.Writer
	level         Level
	showTimestamp bool
	showLevel     bool
	showContext   bool
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

// Log takes a Level and series of values, then outputs them formatted
// accordingly.
func (j *JSONTarget) Log(level Level, values []interface{}, context Ctx) {
	if level < j.level {
		return
	}

	jsonData := make(map[string]interface{}, 1)
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

	fmt.Printf("%+v\n", values)
	fmt.Printf("%+v\n", strValues)

	jsonData["message"] = strings.Join(strValues, " ")
	if j.showContext {
		jsonData["context"] = context
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		log.Println(err)
		return
	}
	jsonBytes = []byte(string(jsonBytes) + "\n")

	if level >= Warn {
		j.errTarget.Write(jsonBytes)
	} else {
		j.outTarget.Write(jsonBytes)
	}
}
