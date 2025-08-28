package slime

import (
	"time"

	"github.com/google/uuid"
	"trpc.group/trpc-go/trpc-go/errs"
)

// clientCfg is the configuration of slime client.
type clientCfg struct {
	Services []serviceCfg `yaml:"service"`
}

// serviceCfg is the configuration of slime service.
type serviceCfg struct {
	// 服务名称
	Name string `yaml:"name"`
	// 被调服务名称
	Callee string `yaml:"callee"`
	// 限流配置指针
	Throttle *throttleCfg `yaml:"retry_hedging_throttle"`
	// 重试和对冲配置
	RetryHedging retryHedgingCfg `yaml:"retry_hedging"`
	// 方法配置切片
	Methods []methodCfg `yaml:"methods"`
}

// repair is used to make configuration file consistent with tRPC-Go.
// Slime only needs naming service, aka Name, and does not care about proto service, aka Callee.
func (cfg *serviceCfg) repair() {
	// 用于修复配置，如果 Name 为空，则使用 Callee 的值填充
	if cfg.Name == "" {
		cfg.Name = cfg.Callee
	}
}

// throttleCfg is the configuration of slime throttle.
type throttleCfg struct {
	MaxTokens  float64 `yaml:"max_tokens"`  // 最大令牌数
	TokenRatio float64 `yaml:"token_ratio"` // 令牌生成比率
}

// methodCfg is the configuration of slime method.
type methodCfg struct {
	Callee       string           `yaml:"callee"`        // 被调方法名称
	RetryHedging *retryHedgingCfg `yaml:"retry_hedging"` // 重试和对冲配置指针
}

// retryHedgingCfg is the configuration of slime retry and hedging.
type retryHedgingCfg struct {
	Retry   *retryCfg   `yaml:"retry"`   // 重试配置指针
	Hedging *hedgingCfg `yaml:"hedging"` // 对冲配置指针
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// hedgingCfg is the configuration of slime hedging.
type hedgingCfg struct {
	Name             string        `yaml:"name"`                  // 对冲策略名称
	MaxAttempts      int           `yaml:"max_attempts"`          // 最大尝试次数
	HedgingDelay     time.Duration `yaml:"hedging_delay"`         // 对冲延迟时间
	NonFatalECs      []int         `yaml:"non_fatal_error_codes"` // 非致命错误码列表
	SkipVisitedNodes *bool         `yaml:"skip_visited_nodes"`    // 是否跳过已访问节点的布尔指针
}

var (
	// 默认最大尝试次数为 2
	defaultHedgingMaxAttempts = 2
	// 默认非致命错误码列表
	defaultNonFatalECs = []int{
		int(errs.RetServerTimeout),
		int(errs.RetClientConnectFail),
		int(errs.RetClientRouteErr),
		int(errs.RetClientNetErr),
	}
)

// repair fix hedgingCfg With default values.
// 修复对冲配置
func (cfg *hedgingCfg) repair() {
	if cfg.MaxAttempts == 0 {
		cfg.MaxAttempts = defaultHedgingMaxAttempts
	}
	if cfg.Name == "" {
		// 生成一个唯一名称
		cfg.Name = "hedging-" + uuid.New().String()
	}
	if len(cfg.NonFatalECs) == 0 {
		// 使用默认错误码列表
		cfg.NonFatalECs = defaultNonFatalECs
	}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// retryCfg is the configuration of slime retry.
type retryCfg struct {
	Name             string     `yaml:"name"`                  // 重试策略名称
	MaxAttempts      int        `yaml:"max_attempts"`          // 最大尝试次数
	Backoff          backoffCfg `yaml:"backoff"`               // 退避策略配置
	RetryableECs     []int      `yaml:"retryable_error_codes"` // 可重试错误码列表
	SkipVisitedNodes *bool      `yaml:"skip_visited_nodes"`    // 是否跳过已访问节点的布尔指针
}

var (
	// 默认最大尝试次数为 2
	defaultRetryMaxAttempts = 2
	// 默认可重试错误码列表
	defaultRetryableECs = []int{
		int(errs.RetServerTimeout),
		int(errs.RetClientConnectFail),
		int(errs.RetClientRouteErr),
		int(errs.RetClientNetErr),
	}
)

// repair fix retryCfg with default values.
func (cfg *retryCfg) repair() {
	if cfg.MaxAttempts == 0 {
		cfg.MaxAttempts = defaultRetryMaxAttempts
	}
	if cfg.Name == "" {
		// 生成一个唯一名称
		cfg.Name = "retry-" + uuid.New().String()
	}
	if len(cfg.RetryableECs) == 0 {
		// 使用默认错误码列表
		cfg.RetryableECs = defaultRetryableECs
	}
}

// backoffCfg is the configuration of slime backoff.
// 退避策略配置结构体
type backoffCfg struct {
	// 指数退避策略配置指针
	Exponential *struct {
		Initial    time.Duration `yaml:"initial"`    // 初始延迟
		Maximum    time.Duration `yaml:"maximum"`    // 最大延迟
		Multiplier int           `yaml:"multiplier"` // 乘数
	} `yaml:"exponential"`
	// 线性退避策略的时间间隔列表
	Linear []time.Duration `yaml:"linear"`
}
