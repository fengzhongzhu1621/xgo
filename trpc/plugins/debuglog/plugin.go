package debuglog

import (
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/plugin"
)

// Plugin is the implement of the debuglog trpc plugin.
type Plugin struct{}

// Type is the type of debuglog trpc plugin.
func (p *Plugin) Type() string {
	return pluginType
}

// Config is the congifuration for the debuglog trpc plugin.
type Config struct {
	LogType       string `yaml:"log_type"`
	ErrLogLevel   string `yaml:"err_log_level"`
	NilLogLevel   string `yaml:"nil_log_level"`
	ServerLogType string `yaml:"server_log_type"`
	ClientLogType string `yaml:"client_log_type"`
	EnableColor   *bool  `yaml:"enable_color"`
	Include       []*RuleItem
	Exclude       []*RuleItem
}

// Setup initializes the debuglog instance.
func (p *Plugin) Setup(name string, configDec plugin.Decoder) error {
	var conf Config

	// 解析配置
	err := configDec.Decode(&conf)
	if err != nil {
		return err
	}

	var serverOpt []Option
	var clientOpt []Option

	// 设置额外信息日志打印函数
	serverLogType := conf.LogType
	if conf.ServerLogType != "" {
		serverLogType = conf.ServerLogType
	}
	serverOpt = append(serverOpt, WithLogFunc(getLogFunc(serverLogType)))

	clientLogType := conf.LogType
	if conf.ClientLogType != "" {
		clientLogType = conf.ClientLogType
	}
	clientOpt = append(clientOpt, WithLogFunc(getLogFunc(clientLogType)))

	// 设置匹配规则
	for _, in := range conf.Include {
		serverOpt = append(serverOpt, WithInclude(in))
		clientOpt = append(clientOpt, WithInclude(in))
	}
	for _, ex := range conf.Exclude {
		serverOpt = append(serverOpt, WithExclude(ex))
		clientOpt = append(clientOpt, WithExclude(ex))
	}

	// 设置日志打印的函数（使用特定的日志格式化器）
	clientOpt = append(clientOpt,
		WithNilLogLevelFunc(getLogLevelFunc(conf.NilLogLevel, "debug")),
		WithErrLogLevelFunc(getLogLevelFunc(conf.ErrLogLevel, "error")),
	)
	serverOpt = append(serverOpt,
		WithNilLogLevelFunc(getLogLevelFunc(conf.NilLogLevel, "debug")),
		WithErrLogLevelFunc(getLogLevelFunc(conf.ErrLogLevel, "error")),
	)

	// 设置日志颜色
	if conf.EnableColor != nil {
		serverOpt = append(serverOpt, WithEnableColor(*conf.EnableColor))
		clientOpt = append(clientOpt, WithEnableColor(*conf.EnableColor))
	}

	// register server and client filter
	// 只注册了 debuglog 过滤器
	filter.Register(pluginName, ServerFilter(serverOpt...), ClientFilter(clientOpt...))

	return nil
}
