package logging

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/logging/output"
	"github.com/fengzhongzhu1621/xgo/logging/zaplogger"
	"github.com/fengzhongzhu1621/xgo/plugin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	pluginType = "log"
)

var (
	// DefaultConsoleWriterFactory is the default console output implementation.
	DefaultConsoleWriterFactory = &ConsoleWriterFactory{}
	// DefaultFileWriterFactory is the default file output implementation.
	DefaultFileWriterFactory = &FileWriterFactory{}
)

var _ plugin.IDecoder = (*Decoder)(nil)

// Decoder decodes the log.
type Decoder struct {
	OutputConfig *config.LogOutputConfig
	Core         zapcore.Core
	ZapLevel     zap.AtomicLevel
}

// Decode decodes writer configuration, copy one.
func (d *Decoder) Decode(cfg interface{}) error {
	output, ok := cfg.(**config.LogOutputConfig)
	if !ok {
		return fmt.Errorf("decoder config type:%T invalid, not **OutputConfig", cfg)
	}
	// 将 d.OutputConfig 的值覆盖输入参数
	*output = d.OutputConfig
	return nil
}

// Factory is the log plugin factory.
// When server start, the configuration is feed to Factory to generate a log instance.
type Factory struct{}

// Type returns the log plugin type.
func (f *Factory) Type() string {
	return pluginType
}

// Setup starts, load and register logs.
func (f *Factory) Setup(name string, dec plugin.IDecoder) error {
	if dec == nil {
		return errors.New("log config decoder empty")
	}
	cfg, callerSkip, err := f.setupConfig(dec)
	if err != nil {
		return err
	}
	logger := NewZapLogWithCallerSkip(cfg, callerSkip)
	if logger == nil {
		return errors.New("new zap logger fail")
	}
	Register(name, logger)
	return nil
}

func (f *Factory) setupConfig(configDec plugin.IDecoder) (config.LogOutputConfigs, int, error) {
	cfg := config.LogOutputConfigs{}
	if err := configDec.Decode(&cfg); err != nil {
		return nil, 0, err
	}
	if len(cfg) == 0 {
		return nil, 0, errors.New("log config output empty")
	}

	// If caller skip is not configured, use 2 as default.
	callerSkip := 2
	for i := 0; i < len(cfg); i++ {
		if cfg[i].CallerSkip != 0 {
			callerSkip = cfg[i].CallerSkip
		}
	}
	return cfg, callerSkip, nil
}

// ConsoleWriterFactory is the console writer instance.
type ConsoleWriterFactory struct {
}

// Type returns the log plugin type.
func (f *ConsoleWriterFactory) Type() string {
	return pluginType
}

// Setup starts, loads and registers console output writer.
func (f *ConsoleWriterFactory) Setup(name string, dec plugin.IDecoder) error {
	if dec == nil {
		return errors.New("console writer decoder empty")
	}
	decoder, ok := dec.(*Decoder)
	if !ok {
		return errors.New("console writer log decoder type invalid")
	}

	// 将 yaml 配置转换为 LogOutputConfig 结构体
	cfg := &config.LogOutputConfig{}
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}

	decoder.Core, decoder.ZapLevel = zaplogger.NewConsoleCore(cfg)
	return nil
}

// FileWriterFactory is the file writer instance Factory.
type FileWriterFactory struct {
}

// Type returns log file type.
func (f *FileWriterFactory) Type() string {
	return pluginType
}

// Setup starts, loads and register file output writer.
func (f *FileWriterFactory) Setup(name string, dec plugin.IDecoder) error {
	if dec == nil {
		return errors.New("file writer decoder empty")
	}
	decoder, ok := dec.(*Decoder)
	if !ok {
		return errors.New("file writer log decoder type invalid")
	}
	if err := f.setupConfig(decoder); err != nil {
		return err
	}
	return nil
}

func (f *FileWriterFactory) setupConfig(decoder *Decoder) error {
	cfg := &config.LogOutputConfig{}
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}
	if cfg.WriteConfig.LogPath != "" {
		cfg.WriteConfig.Filename = filepath.Join(cfg.WriteConfig.LogPath, cfg.WriteConfig.Filename)
	}
	if cfg.WriteConfig.RollType == "" {
		cfg.WriteConfig.RollType = output.RollBySize
	}

	core, level, err := zaplogger.NewFileCore(cfg)
	if err != nil {
		return err
	}
	decoder.Core, decoder.ZapLevel = core, level
	return nil
}
