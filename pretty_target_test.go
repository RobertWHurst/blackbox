package blackbox_test

import (
	"bytes"
	"testing"

	"github.com/RobertWHurst/blackbox"
	"github.com/stretchr/testify/assert"
)

func TestPrettyTarget(t *testing.T) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	prettyTarget := blackbox.NewPrettyTarget(outBuf, errBuf)

	values := make([]interface{}, 1)
	values[0] = "Hello Test"

	prettyTarget.Log(blackbox.Trace, values, blackbox.Ctx{"key": "value"})

	assert.Regexp(
		t,
		`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}-\d{2}:\d{2} \x1b\[\d{2}mtrace\x1b\[`+
			`0m   Hello Test \x1b\[\d{2}mkey\x1b\[0m=value\n$`,
		outBuf.String(),
	)
}

func TestPrettyTargetSetLevel(t *testing.T) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	prettyTarget := blackbox.NewPrettyTarget(outBuf, errBuf)

	prettyTarget.SetLevel(blackbox.Info)

	values := make([]interface{}, 1)
	values[0] = "Filtered Message"

	prettyTarget.Log(blackbox.Trace, values, blackbox.Ctx{"x": "y"})

	values = make([]interface{}, 1)
	values[0] = "Hello Test"

	prettyTarget.Log(blackbox.Info, values, blackbox.Ctx{"key": "value"})

	assert.NotRegexp(t, `Filtered Message`, outBuf.String())
	assert.Regexp(t, `Hello Test`, outBuf.String())
}

func TestPrettyTargetShowTimestamp(t *testing.T) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	prettyTarget := blackbox.NewPrettyTarget(outBuf, errBuf)

	prettyTarget.ShowTimestamp(false)

	values := make([]interface{}, 1)
	values[0] = "Hello Test"

	prettyTarget.Log(blackbox.Trace, values, blackbox.Ctx{"key": "value"})

	assert.Regexp(
		t,
		`^\x1b\[\d{2}mtrace\x1b\[0m   Hello Test \x1b\[\d{2}mkey\x1b\[0m=value\n$`,
		outBuf.String(),
	)
}

func TestPrettyTargetSelectContext(t *testing.T) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	prettyTarget := blackbox.NewPrettyTarget(outBuf, errBuf)

	prettyTarget.SelectContext("x")

	values := make([]interface{}, 1)
	values[0] = "Hello Test"

	prettyTarget.Log(blackbox.Trace, values, blackbox.Ctx{"key": "value", "x": "y"})

	assert.Regexp(t, "x[^ ]*=y", outBuf.String())
}

func TestPrettyTargetShowContext(t *testing.T) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	prettyTarget := blackbox.NewPrettyTarget(outBuf, errBuf)

	prettyTarget.ShowContext(false)

	values := make([]interface{}, 1)
	values[0] = "Hello Test"

	prettyTarget.Log(blackbox.Trace, values, blackbox.Ctx{"key": "value"})

	assert.Regexp(
		t,
		`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}-\d{2}:\d{2} \x1b\[\d{2}mtrace\x1b\[`+
			`0m   Hello Test\n$`,
		outBuf.String(),
	)
}

func TestPrettyTargetUseColor(t *testing.T) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	prettyTarget := blackbox.NewPrettyTarget(outBuf, errBuf)

	prettyTarget.UseColor(false)

	values := make([]interface{}, 1)
	values[0] = "Hello Test"

	prettyTarget.Log(blackbox.Trace, values, blackbox.Ctx{"key": "value"})

	assert.Regexp(
		t,
		`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}-\d{2}:\d{2} trace   Hello Test key=`+
			`value\n$`,
		outBuf.String(),
	)
}
