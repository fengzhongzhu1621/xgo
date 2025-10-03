package zaplogger

import (
	"fmt"
	"os"
	"time"

	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/logging/level"
	"github.com/fengzhongzhu1621/xgo/logging/output"
	"github.com/fengzhongzhu1621/xgo/logging/output/rollwriter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewFormatEncoder is the function type for creating a format encoder out of an encoder config.
type NewFormatEncoder func(zapcore.EncoderConfig) zapcore.Encoder

var formatEncoders = map[string]NewFormatEncoder{
	"console": zapcore.NewConsoleEncoder,
	"json":    zapcore.NewJSONEncoder,
}

// RegisterFormatEncoder registers a NewFormatEncoder with the specified formatName key.
// The existing formats include "console" and "json", but you can override these format encoders
// or provide a new custom one.
func RegisterFormatEncoder(formatName string, newFormatEncoder NewFormatEncoder) {
	formatEncoders[formatName] = newFormatEncoder
}

// CustomTimeFormat customize time format.
// Deprecated: Use https://pkg.go.dev/time#Time.Format instead.
func CustomTimeFormat(t time.Time, format string) string {
	return t.Format(format)
}

// DefaultTimeFormat returns the default time format "2006-01-02 15:04:05.000".
// Deprecated: Use https://pkg.go.dev/time#Time.AppendFormat instead.
func DefaultTimeFormat(t time.Time) []byte {
	return defaultTimeFormat(t)
}

func NewConsoleCore(c *config.LogOutputConfig) (zapcore.Core, zap.AtomicLevel) {
	lvl := zap.NewAtomicLevelAt(level.Levels[c.Level])
	return zapcore.NewCore(
		newEncoder(c),
		zapcore.Lock(os.Stdout),
		lvl), lvl
}

func NewFileCore(c *config.LogOutputConfig) (zapcore.Core, zap.AtomicLevel, error) {
	opts := []rollwriter.Option{
		rollwriter.WithMaxAge(c.WriteConfig.MaxAge),
		rollwriter.WithMaxBackups(c.WriteConfig.MaxBackups),
		rollwriter.WithCompress(c.WriteConfig.Compress),
		rollwriter.WithMaxSize(c.WriteConfig.MaxSize),
	}
	// roll by time.
	if c.WriteConfig.RollType != output.RollBySize {
		opts = append(opts, rollwriter.WithRotationTime(c.WriteConfig.TimeUnit.Format()))
	}
	writer, err := rollwriter.NewRollWriter(c.WriteConfig.Filename, opts...)
	if err != nil {
		return nil, zap.AtomicLevel{}, err
	}

	// write mode.
	var ws zapcore.WriteSyncer
	switch m := c.WriteConfig.WriteMode; m {
	case 0, output.WriteFast:
		// Use WriteFast as default mode.
		// It has better performance, discards logs on full and avoid blocking service.
		ws = rollwriter.NewAsyncRollWriter(writer, rollwriter.WithDropLog(true))
	case output.WriteSync:
		ws = zapcore.AddSync(writer)
	case output.WriteAsync:
		ws = rollwriter.NewAsyncRollWriter(writer, rollwriter.WithDropLog(false))
	default:
		return nil, zap.AtomicLevel{}, fmt.Errorf("validating WriteMode parameter: got %d, "+
			"but expect one of WriteFast(%d), WriteAsync(%d), or WriteSync(%d)", m, output.WriteFast, output.WriteAsync, output.WriteSync)
	}

	// log level.
	lvl := zap.NewAtomicLevelAt(level.Levels[c.Level])
	return zapcore.NewCore(
		newEncoder(c),
		ws, lvl,
	), lvl, nil
}

func newEncoder(c *config.LogOutputConfig) zapcore.Encoder {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        GetLogEncoderKey("T", c.FormatConfig.TimeKey),
		LevelKey:       GetLogEncoderKey("L", c.FormatConfig.LevelKey),
		NameKey:        GetLogEncoderKey("N", c.FormatConfig.NameKey),
		CallerKey:      GetLogEncoderKey("C", c.FormatConfig.CallerKey),
		FunctionKey:    GetLogEncoderKey(zapcore.OmitKey, c.FormatConfig.FunctionKey),
		MessageKey:     GetLogEncoderKey("M", c.FormatConfig.MessageKey),
		StacktraceKey:  GetLogEncoderKey("S", c.FormatConfig.StacktraceKey),
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     NewTimeEncoder(c.FormatConfig.TimeFmt),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	if c.EnableColor {
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	if newFormatEncoder, ok := formatEncoders[c.Formatter]; ok {
		return newFormatEncoder(encoderCfg)
	}
	// Defaults to console encoder.
	return zapcore.NewConsoleEncoder(encoderCfg)
}

// NewTimeEncoder creates a time format encoder.
func NewTimeEncoder(format string) zapcore.TimeEncoder {
	switch format {
	case "":
		return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendByteString(defaultTimeFormat(t))
		}
	case "seconds":
		return zapcore.EpochTimeEncoder
	case "milliseconds":
		return zapcore.EpochMillisTimeEncoder
	case "nanoseconds":
		return zapcore.EpochNanosTimeEncoder
	default:
		return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(format))
		}
	}
}

// defaultTimeFormat returns the default time format "2006-01-02 15:04:05.000",
// which performs better than https://pkg.go.dev/time#Time.AppendFormat.
func defaultTimeFormat(t time.Time) []byte {
	t = t.Local()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	micros := t.Nanosecond() / 1000

	buf := make([]byte, 23)
	buf[0] = byte((year/1000)%10) + '0'
	buf[1] = byte((year/100)%10) + '0'
	buf[2] = byte((year/10)%10) + '0'
	buf[3] = byte(year%10) + '0'
	buf[4] = '-'
	buf[5] = byte((month)/10) + '0'
	buf[6] = byte((month)%10) + '0'
	buf[7] = '-'
	buf[8] = byte((day)/10) + '0'
	buf[9] = byte((day)%10) + '0'
	buf[10] = ' '
	buf[11] = byte((hour)/10) + '0'
	buf[12] = byte((hour)%10) + '0'
	buf[13] = ':'
	buf[14] = byte((minute)/10) + '0'
	buf[15] = byte((minute)%10) + '0'
	buf[16] = ':'
	buf[17] = byte((second)/10) + '0'
	buf[18] = byte((second)%10) + '0'
	buf[19] = '.'
	buf[20] = byte((micros/100000)%10) + '0'
	buf[21] = byte((micros/10000)%10) + '0'
	buf[22] = byte((micros/1000)%10) + '0'
	return buf
}

// GetLogEncoderKey gets user defined log output name, uses defKey if empty.
func GetLogEncoderKey(defKey, key string) string {
	if key == "" {
		return defKey
	}
	return key
}
