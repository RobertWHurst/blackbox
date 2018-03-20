package main

import (
	"os"

	"github.com/RobertWHurst/blackbox"
)

func main() {
	logger := blackbox.New()

	logger.AddTarget(blackbox.NewPrettyTarget(os.Stdout, os.Stderr))

	logger.
		Ctx(blackbox.Ctx{"key1": "value1"}).
		Trace("a trace level message").
		Debug("a debug level message").
		Ctx(blackbox.Ctx{"key2": "value2"}).
		Verbose("a verbose level message").
		Info("a info level message").
		Ctx(blackbox.Ctx{"key3": "value3"}).
		Warn("a warn level message").
		Error("a error level message").
		Ctx(blackbox.Ctx{"key3": nil}).
		Fatal("a fatal level message")
}
