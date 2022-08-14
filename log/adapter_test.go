package log

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogFields_Copy(t *testing.T) {
	fields1 := LogFields{"foo": "bar"}

	fields2 := fields1.Copy()
	fields2["foo"] = "baz"

	assert.Equal(t, fields1["foo"], "bar")
	assert.Equal(t, fields2["foo"], "baz")
}

func TestStdLogger_with(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	cleanLogger := NewStdLoggerWithOut(buf, true, true, "")

	// 添加默认字段
	withLogFieldsLogger := cleanLogger.With(LogFields{"foo": "1"})

	for name, logger := range map[string]LoggerAdapter{"clean": cleanLogger, "with": withLogFieldsLogger} {
		logger.Error(name, nil, LogFields{"bar": "2"})
		logger.Info(name, LogFields{"bar": "2"})
		logger.Debug(name, LogFields{"bar": "2"})
		logger.Trace(name, LogFields{"bar": "2"})
	}

	cleanLoggerOut := buf.String()
	assert.Contains(t, cleanLoggerOut, `level=ERROR msg="clean" bar=2 err=<nil>`)
	assert.Contains(t, cleanLoggerOut, `level=INFO  msg="clean" bar=2`)
	assert.Contains(t, cleanLoggerOut, `level=TRACE msg="clean" bar=2`)

	assert.Contains(t, cleanLoggerOut, `level=ERROR msg="with" bar=2 err=<nil> foo=1`)
	assert.Contains(t, cleanLoggerOut, `level=INFO  msg="with" bar=2 foo=1`)
	assert.Contains(t, cleanLoggerOut, `level=TRACE msg="with" bar=2 foo=1`)
}

type stringer struct{}

func (s stringer) String() string {
	return "stringer"
}

func TestStdLoggerAdapter_stringer_field(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	logger := NewStdLoggerWithOut(buf, true, true, "")

	logger.Info("foo", LogFields{"foo": stringer{}})

	out := buf.String()
	assert.Contains(t, out, `foo=stringer`)
}

func TestStdLoggerAdapter_field_with_space(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	logger := NewStdLoggerWithOut(buf, true, true, "")

	logger.Info("foo", LogFields{"foo": `bar baz`})

	out := buf.String()
	assert.Contains(t, out, `foo="bar baz"`)
}

func TestCaptureLoggerAdapter(t *testing.T) {
	var logger LoggerAdapter = NewCaptureLogger()

	err := errors.New("error")

	logger = logger.With(LogFields{"default": "field"})
	logger.Error("error", err, LogFields{"bar": "2"})
	logger.Info("info", LogFields{"bar": "2"})
	logger.Debug("debug", LogFields{"bar": "2"})
	logger.Trace("trace", LogFields{"bar": "2"})

	expectedLogs := map[LogLevel][]CapturedMessage{
		TraceLogLevel: {
			CapturedMessage{
				Level:  TraceLogLevel,
				Fields: LogFields{"bar": "2", "default": "field"},
				Msg:    "trace",
				Err:    error(nil),
			},
		},
		DebugLogLevel: {
			CapturedMessage{
				Level:  DebugLogLevel,
				Fields: LogFields{"default": "field", "bar": "2"},
				Msg:    "debug",
				Err:    error(nil),
			},
		},
		InfoLogLevel: {
			CapturedMessage{
				Level:  InfoLogLevel,
				Fields: LogFields{"default": "field", "bar": "2"},
				Msg:    "info",
				Err:    error(nil),
			},
		},
		ErrorLogLevel: {
			CapturedMessage{
				Level:  ErrorLogLevel,
				Fields: LogFields{"default": "field", "bar": "2"},
				Msg:    "error",
				Err:    err,
			},
		},
	}

	capturedLogger := logger.(*CaptureLoggerAdapter)
	assert.EqualValues(t, expectedLogs, capturedLogger.Captured())

	for _, logs := range expectedLogs {
		for _, log := range logs {
			assert.True(t, capturedLogger.Has(log))
		}
	}

	assert.False(t, capturedLogger.Has(CapturedMessage{
		Level:  0,
		Fields: nil,
		Msg:    "",
		Err:    nil,
	}))

	assert.True(t, capturedLogger.HasError(err))
	assert.False(t, capturedLogger.HasError(errors.New("foo")))
}
