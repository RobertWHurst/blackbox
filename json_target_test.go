package blackbox_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/RobertWHurst/blackbox"
	"github.com/stretchr/testify/assert"
)

type JSONOutput struct {
	Context blackbox.Ctx
	Level   string
	Message string
	Time    string
}

func TestJsonTarget(t *testing.T) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	jsonTarget := blackbox.NewJSONTarget(outBuf, errBuf)

	values := make([]interface{}, 1)
	values[0] = "Hello Test"

	jsonTarget.Log(blackbox.Trace, values, blackbox.Ctx{"key": "value"})

	var output JSONOutput
	json.Unmarshal(outBuf.Bytes(), &output)

	assert.NotEmpty(t, output.Time)
	assert.Equal(t, "Hello Test", output.Message)
	assert.Equal(t, "trace", output.Level)
	assert.Equal(t, "value", output.Context["key"])
}

func TestJsonTargetSetLevel(t *testing.T) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	jsonTarget := blackbox.NewJSONTarget(outBuf, errBuf)

	jsonTarget.SetLevel(blackbox.Info)

	values := make([]interface{}, 1)
	values[0] = "Filtered Message"

	jsonTarget.Log(blackbox.Trace, values, blackbox.Ctx{"x": "y"})

	values = make([]interface{}, 1)
	values[0] = "Hello Test"

	jsonTarget.Log(blackbox.Info, values, blackbox.Ctx{"key": "value"})

	assert.NotRegexp(t, `Filtered Message`, outBuf.String())
	assert.Regexp(t, `Hello Test`, outBuf.String())
}

func TestJsonTargetShowTimestamp(t *testing.T) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	jsonTarget := blackbox.NewJSONTarget(outBuf, errBuf)

	jsonTarget.ShowTimestamp(false)

	values := make([]interface{}, 1)
	values[0] = "Hello Test"

	jsonTarget.Log(blackbox.Trace, values, blackbox.Ctx{"key": "value"})

	var output JSONOutput
	json.Unmarshal(outBuf.Bytes(), &output)

	assert.Empty(t, output.Time)
	assert.Equal(t, "Hello Test", output.Message)
	assert.Equal(t, "trace", output.Level)
	assert.Equal(t, "value", output.Context["key"])
}

func TestJsonTargetShowContext(t *testing.T) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	jsonTarget := blackbox.NewJSONTarget(outBuf, errBuf)

	jsonTarget.ShowContext(false)

	values := make([]interface{}, 1)
	values[0] = "Hello Test"

	jsonTarget.Log(blackbox.Trace, values, blackbox.Ctx{"key": "value"})

	var output JSONOutput
	json.Unmarshal(outBuf.Bytes(), &output)

	assert.NotEmpty(t, output.Time)
	assert.Equal(t, "Hello Test", output.Message)
	assert.Equal(t, "trace", output.Level)
	assert.Empty(t, output.Context)
}
