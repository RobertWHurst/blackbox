package blackbox_test

import (
	"testing"

	"github.com/RobertWHurst/blackbox"
	"github.com/stretchr/testify/assert"
)

func TestNewToReturnLogger(t *testing.T) {
	var _ *blackbox.Logger = blackbox.New()
}

func TestNewWithCtxToReturnLogger(t *testing.T) {
	var logger *blackbox.Logger = blackbox.NewWithCtx(blackbox.Ctx{"key": "value"})
	testTarget := blackbox.NewTestTarget()
	logger.AddTarget(testTarget)

	logger.Log(blackbox.Info, "Message")

	logged, ok := testTarget.LastLogged()

	assert.Equal(t, true, ok)
	assert.Equal(t, "Message", logged.Values[0].(string))
	assert.Equal(t, blackbox.Ctx{"key": "value"}, logged.Context)
}

func TestLoggerLogToReturnLogger(t *testing.T) {
	logger := blackbox.New()
	testTarget := blackbox.NewTestTarget()
	logger.AddTarget(testTarget)

	logger.Log(blackbox.Info, "Message")

	logged, ok := testTarget.LastLogged()

	assert.Equal(t, true, ok)
	assert.Equal(t, blackbox.Info, logged.Level)
	assert.Equal(t, "Message", logged.Values[0].(string))
}

func TestLoggerLogfToReturnLogger(t *testing.T) {
	logger := blackbox.New()
	testTarget := blackbox.NewTestTarget()
	logger.AddTarget(testTarget)

	value := struct {
		A string
		B string
		C string
	}{A: "1", B: "2", C: "3"}
	logger.Logf(blackbox.Info, "%+v", value)

	logged, ok := testTarget.LastLogged()

	assert.Equal(t, true, ok)
	assert.Equal(t, blackbox.Info, logged.Level)
	assert.Equal(t, "{A:1 B:2 C:3}", logged.Values[0].(string))
}

func TestLoggerCtxToReturnLogger(t *testing.T) {
	logger := blackbox.New()
	testTarget := blackbox.NewTestTarget()
	logger.AddTarget(testTarget)

	var subLogger *blackbox.Logger = logger.WithCtx(blackbox.Ctx{"key": "value"})
	subLogger.Log(blackbox.Info, "Message")

	logged, ok := testTarget.LastLogged()

	assert.Equal(t, true, ok)
	assert.Equal(t, blackbox.Info, logged.Level)
	assert.Equal(t, "Message", logged.Values[0].(string))
	assert.Equal(t, blackbox.Ctx{"key": "value"}, logged.Context)
}

func TestLoggerInlineCtx(t *testing.T) {
	logger := blackbox.New()
	testTarget := blackbox.NewTestTarget()
	logger.AddTarget(testTarget)

	logger.Log(blackbox.Info, "Message", blackbox.Ctx{"key": "value"})

	logged, ok := testTarget.LastLogged()

	assert.Equal(t, true, ok)
	assert.Equal(t, blackbox.Info, logged.Level)
	assert.Equal(t, "Message", logged.Values[0].(string))
	assert.Equal(t, blackbox.Ctx{"key": "value"}, logged.Context)
}
