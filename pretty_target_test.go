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

	values := make([]any, 1)
	values[0] = "Hello Test"

	prettyTarget.Log("AAA-AAA", blackbox.Trace, values, blackbox.Ctx{"key": "value"}, nil)

	assert.Regexp(
		t,
		`^\x1b\[\d+m\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} [A-Z]+ \x1b\[0m\x1b\[\d+mtrace\x1b\[`+
			`0m   \x1b\[\d+;?\d*mHello Test\x1b\[0m \x1b\[\d+;?\d*mkey\x1b\[0m=\x1b\[\d+m?value\x1b\[0m\n$`,
		outBuf.String(),
	)
}

func TestPrettyTargetSetLevel(t *testing.T) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	prettyTarget := blackbox.NewPrettyTarget(outBuf, errBuf)

	prettyTarget.SetLevel(blackbox.Info)

	values := make([]any, 1)
	values[0] = "Filtered Message"

	prettyTarget.Log("AAA-AAA", blackbox.Trace, values, blackbox.Ctx{"x": "y"}, nil)

	values = make([]any, 1)
	values[0] = "Hello Test"

	prettyTarget.Log("AAA-AAA", blackbox.Info, values, blackbox.Ctx{"key": "value"}, nil)

	assert.NotRegexp(t, `Filtered Message`, outBuf.String())
	assert.Regexp(t, `Hello Test`, outBuf.String())
}

func TestPrettyTargetShowTimestamp(t *testing.T) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	prettyTarget := blackbox.NewPrettyTarget(outBuf, errBuf)

	prettyTarget.ShowTimestamp(false)

	values := make([]any, 1)
	values[0] = "Hello Test"

	prettyTarget.Log("AAA-AAA", blackbox.Trace, values, blackbox.Ctx{"key": "value"}, nil)

	assert.Regexp(
		t,
		`^\x1b\[\d+mtrace\x1b\[0m   .+Hello Test.+ .+key.+=.+value`,
		outBuf.String(),
	)
}

func TestPrettyTargetSelectContext(t *testing.T) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	prettyTarget := blackbox.NewPrettyTarget(outBuf, errBuf)

	prettyTarget.SelectContext("x")

	values := make([]any, 1)
	values[0] = "Hello Test"

	prettyTarget.Log("AAA-AAA", blackbox.Trace, values, blackbox.Ctx{"key": "value", "x": "y"}, nil)

	output := outBuf.String()
	assert.Contains(t, output, "x")
	assert.Contains(t, output, "y")
	assert.NotContains(t, output, "value")
}

func TestPrettyTargetShowContext(t *testing.T) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	prettyTarget := blackbox.NewPrettyTarget(outBuf, errBuf)

	prettyTarget.ShowContext(false)

	values := make([]any, 1)
	values[0] = "Hello Test"

	prettyTarget.Log("AAA-AAA", blackbox.Trace, values, blackbox.Ctx{"key": "value"}, nil)

	assert.NotRegexp(t, `key`, outBuf.String())
	assert.Regexp(t, `Hello Test`, outBuf.String())
}

func TestPrettyTargetUseColor(t *testing.T) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	prettyTarget := blackbox.NewPrettyTarget(outBuf, errBuf)

	prettyTarget.UseColor(false)

	values := make([]any, 1)
	values[0] = "Hello Test"

	prettyTarget.Log("AAA-AAA", blackbox.Trace, values, blackbox.Ctx{"key": "value"}, nil)

	assert.Regexp(
		t,
		`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} [A-Z]+ trace   Hello Test key=value\n$`,
		outBuf.String(),
	)
}

func TestPrettyTargetUseSource(t *testing.T) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	prettyTarget := blackbox.NewPrettyTarget(outBuf, errBuf)

	prettyTarget.UseColor(false).ShowSource(true)

	values := make([]any, 1)
	values[0] = "Hello Test"

	prettyTarget.Log("AAA-AAA", blackbox.Trace, values, blackbox.Ctx{"key": "value"}, func() *blackbox.Source {
		return &blackbox.Source{
			File:     "file.go",
			Line:     123,
			Function: "functionName",
		}
	})

	assert.Regexp(
		t,
		`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} [A-Z]+ trace   Hello Test key=value @=> file\.go:123 - functionName\n$`,
		outBuf.String(),
	)
}

func TestPrettyTargetHiddenContextKeys(t *testing.T) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	prettyTarget := blackbox.NewPrettyTarget(outBuf, errBuf)

	prettyTarget.UseColor(false)

	values := make([]any, 1)
	values[0] = "Hello Test"

	prettyTarget.Log("AAA-AAA", blackbox.Trace, values, blackbox.Ctx{
		"key":     "value",
		"-hidden": "secret",
	}, nil)

	assert.Regexp(t, `key=value`, outBuf.String())
	assert.NotRegexp(t, `hidden`, outBuf.String())
	assert.NotRegexp(t, `secret`, outBuf.String())
}
