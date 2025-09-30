package blackbox

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// PrettyTarget is a Target that produces newline separated human readable
// output suitable for stdout and stderr. It also supports colorized log levels.
type PrettyTarget struct {
	showLoggerID  bool
	showTimestamp bool
	showLevel     bool
	showContext   bool
	useColor      bool
	useSource     bool
	level         Level
	outTarget     io.Writer
	errTarget     io.Writer
	contextFields []string
}

var _ Target = &PrettyTarget{}

// NewPrettyTarget creates a PrettyTarget for use with a logger
func NewPrettyTarget(outTarget io.Writer, errTarget io.Writer) *PrettyTarget {
	return &PrettyTarget{
		showTimestamp: true,
		showLevel:     true,
		showContext:   true,
		useColor:      true,
		level:         Trace,
		outTarget:     outTarget,
		errTarget:     errTarget,
	}
}

// SetLevel sets the minimum log level that PrettyTarget will output. Note that
// this setting is independent of the log level set on the logger itself.
func (s *PrettyTarget) SetLevel(level Level) *PrettyTarget {
	s.level = level
	return s
}

// ShowLoggerID will enable or disable logger ID values in the output depending
// on the boolean value passed. Logger IDs can be useful when the output of
// multiple loggers are viewed together.
func (s *PrettyTarget) ShowLoggerID(b bool) *PrettyTarget {
	s.showLoggerID = b
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

// SelectContext will limit the context key value pairs in the output to only
// those specified as arguments to SelectContext. If select context is called
// no arguments then all context key value pairs will be output.
func (s *PrettyTarget) SelectContext(fields ...string) *PrettyTarget {
	s.contextFields = fields
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

// ShowSource enables the inclusion of source
func (s *PrettyTarget) ShowSource(b bool) *PrettyTarget {
	s.useSource = b
	return s
}

// Log takes a Level and series of values, then outputs them formatted
// accordingly.
func (s *PrettyTarget) Log(loggerID string, level Level, values []any, context Ctx, getSource func() *Source) {
	if level < s.level {
		return
	}

	str := ""
	if s.showLoggerID {
		loggerIDStr := loggerID + " "
		if s.useColor {
			loggerIDStr = wrapStrInColorCodes("loggerID", loggerIDStr)
		}
		str += loggerIDStr
	}

	if s.showTimestamp {
		timestampStr := time.Now().Local().Format("2006-01-02 15:04:05") + " "
		if s.useColor {
			timestampStr = wrapStrInColorCodes("timestamp", timestampStr)
		}
		str += timestampStr
	}

	if s.showLevel {
		levelStr := level.String()
		var padStr string
		for i := len(levelStr); i < 7; i++ {
			padStr += " "
		}
		if s.useColor {
			levelStr = wrapStrInAnsiLevelColorCodes(level, levelStr)
		}
		str += levelStr + padStr + " "
	}

	valueStrs := make([]string, 0)
	for _, value := range values {
		valueStrs = append(valueStrs, fmt.Sprintf("%+v", value))
	}
	valueStr := strings.Join(valueStrs, " ")
	if s.useColor {
		valueStr = wrapStrInColorCodes("value", valueStr)
	}
	str += valueStr

	if s.showContext {
		contextStrs := make([]string, 0)
		for key, value := range context {
			if len(s.contextFields) != 0 {
				skipField := true
				for _, field := range s.contextFields {
					if key == field {
						skipField = false
						break
					}
				}
				if skipField {
					continue
				}
			}
			if s.useColor {
				key = wrapStrInColorCodes("contextKey", key)
			}
			formattedValue := strings.Replace(fmt.Sprintf("%+v", value), "\n", "\\n", -1)
			if s.useColor {
				formattedValue = wrapStrInColorCodes("contextValue", formattedValue)
			}
			contextStrs = append(contextStrs, key+"="+formattedValue)
		}
		sort.Strings(contextStrs)
		contextStr := strings.Join(contextStrs, " ")
		str += " " + contextStr
	}

	if s.useSource {
		source := getSource()
		if source == nil {
			return
		}
		functionAndPackageName := source.Function
		funcPathChunks := strings.Split(functionAndPackageName, "/")
		if len(funcPathChunks) > 0 {
			functionAndPackageName = funcPathChunks[len(funcPathChunks)-1]
		}
		if s.useColor {
			chunks := strings.Split(functionAndPackageName, ".")
			colorizedChunks := make([]string, len(chunks))
			for i, chunk := range chunks {
				colorizedChunks[i] = wrapStrInColorCodes("packageAndFunctionName", chunk)
			}
			functionAndPackageName = strings.Join(colorizedChunks, ".")
		}
		filePath := source.File
		cwd, getwdErr := os.Getwd()
		if getwdErr == nil {
			relFilePath, relErr := filepath.Rel(cwd, source.File)
			if relErr == nil {
				filePath = relFilePath
			}
		}
		if s.useColor {
			filePath = wrapStrInColorCodes("filePath", filePath)
		}
		lineNumber := fmt.Sprintf("%d", source.Line)
		if s.useColor {
			lineNumber = wrapStrInColorCodes("lineNumber", lineNumber)
		}
		separator := "@=>"
		if s.useColor {
			separator = wrapStrInColorCodes("separator", separator)
		}
		sourceStr := fmt.Sprintf(" %s %s:%s - %s", separator, filePath, lineNumber, functionAndPackageName)
		str += sourceStr
	}

	str += "\n"
	s.writeByLevel(level, str)
}

func (s *PrettyTarget) writeByLevel(level Level, str string) {
	var err error
	if level >= Warn {
		_, err = s.errTarget.Write([]byte(str))
	} else {
		_, err = s.outTarget.Write([]byte(str))
	}
	if err != nil {
		panic(err)
	}
}

func wrapStrInAnsiLevelColorCodes(level Level, str string) string {
	switch level {
	case Trace:
		str = "\u001b[35m" + str + "\u001b[0m"
	case Debug:
		str = "\u001b[34m" + str + "\u001b[0m"
	case Verbose:
		str = "\u001b[36m" + str + "\u001b[0m"
	case Info:
		str = "\u001b[32m" + str + "\u001b[0m"
	case Warn:
		str = "\u001b[33m" + str + "\u001b[0m"
	case Error:
		str = "\u001b[31m" + str + "\u001b[0m"
	case Fatal:
		str = "\u001b[37m\u001b[41;1m" + str + "\u001b[0m"
	case Panic:
		str = "\u001b[37m\u001b[45;1m" + str + "\u001b[0m"
	}
	return str
}

func wrapStrInColorCodes(kind string, str string) string {
	switch kind {
	case "loggerID":
		return "\u001b[90;1m" + str + "\u001b[0m"
	case "timestamp":
		return "\u001b[90m" + str + "\u001b[0m"
	case "value":
		return "\u001b[37;1m" + str + "\u001b[0m"
	case "separator":
		return "\u001b[90m" + str + "\u001b[0m"
	case "packageAndFunctionName":
		return "\u001b[32m" + str + "\u001b[0m"
	case "filePath":
		return "\u001b[33m" + str + "\u001b[0m"
	case "lineNumber":
		return "\u001b[35m" + str + "\u001b[0m"
	case "contextKey":
		return "\u001b[90;1m" + str + "\u001b[0m"
	case "contextValue":
		return "\u001b[90m" + str + "\u001b[0m"
	}
	return str
}
