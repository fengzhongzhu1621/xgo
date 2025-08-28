package metrics

import (
	"context"
	"strconv"
	"time"

	"github.com/fengzhongzhu1621/xgo/datetime"
	"github.com/fengzhongzhu1621/xgo/opentelemetry/metrics"
	"github.com/fengzhongzhu1621/xgo/trpc/plugins/slime/view"
	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/errs"
)

// Report report view.Stat.
type Report struct {
	counter   metrics.ICounter
	histogram metrics.IHistogram
}

// NewReport create a new Report from Emitter and tagPairs.
func NewReport(emitter IEmitter, tagPairs ...string) *Report {
	return &Report{
		counter:   metrics.WrapCounter(emitter, tagPairs...),
		histogram: metrics.WrapHistogram(emitter, tagPairs...),
	}
}

// tagsFromCtx 该方法从 context.Context 中提取调用链信息（调用方、被调用方、方法名），封装为固定长度的字符串数组返回，
// 用于打标（Tagging）监控指标。核心逻辑是通过 codec.Message 解析上下文中的 RPC 消息元数据。
func (r *Report) tagsFromCtx(ctx context.Context) []string {
	var tags [6]string
	tags[0] = TagCaller
	tags[2] = TagCallee
	tags[4] = TagMethod

	msg := codec.Message(ctx)

	// 从上下文提取调用方信息
	if caller := msg.CallerServiceName(); caller == "" {
		tags[1] = "unknown"
	} else {
		tags[1] = caller
	}

	// 提取被调用方信息
	if callee := msg.CalleeServiceName(); callee == "" {
		tags[3] = "unknown"
	} else {
		tags[3] = callee
	}

	// 提取方法名
	if method := msg.CalleeMethod(); method == "" {
		tags[5] = "unknown"
	} else {
		tags[5] = method
	}

	return tags[:]
}

// Report reports info of view.Stat. ctx is used to retrieve caller, callee and method.
func (r *Report) Report(ctx context.Context, stat view.IStat) {
	// 初始化与上下文标签提取
	// 从上下文 ctx 中提取基础标签（如 caller、callee、method）
	tags := r.tagsFromCtx(ctx)

	// 遍历每次重试/对冲请求（IAttempt 接口实现），检查是否存在禁止重试的标志（NoMoreAttempt）。
	var noMoreAttempt bool
	for _, a := range stat.Attempts() {
		// 若任意一次尝试被标记为禁止重试，则全局标记 noMoreAttempt 为 true，反映服务端策略限制
		if a.NoMoreAttempt() {
			noMoreAttempt = true
		}

		// The order of tags must match TagNamesReal.
		// 标签生成：在基础标签上追加尝试级标签
		// TagErrCodes：错误码（如 HTTP 500），分类统计失败原因。
		// TagInflight：标识请求是否在处理中（bool），用于实时流量监控。
		// TagNoMoreAttempt：标记服务端是否禁止重试。
		realTags := append(tags,
			TagErrCodes, string(errs.Code(a.Error())),
			TagInflight, strconv.FormatBool(a.Inflight()),
			TagNoMoreAttempt, strconv.FormatBool(a.NoMoreAttempt()))

		// 计数器更新：记录真实请求量（FQNRealRequest），每次尝试计数 +1，反映实际负载。
		r.counter.Inc(FQNRealRequest, 1, realTags...)

		// 若请求已完成（!Inflight），使用实际结束时间 a.End()；否则使用当前时间。
		var endTime = time.Now()
		if !a.Inflight() {
			endTime = a.End()
		}

		// 计算耗时（endTime - startTime）并转换为毫秒，记录到直方图（FQNRealCostMs），支持分位分析（如 P99）
		r.histogram.Observe(FQNRealCostMs, datetime.Milliseconds(endTime.Sub(a.Start())), realTags...)
	}

	// 应用层指标打标与记录
	// The order of tags must match TagNamesApp.
	appTags := append(tags,
		TagAttempts, strconv.Itoa(len(stat.Attempts())),
		TagErrCodes, string(errs.Code(stat.Error())),
		TagThrottled, strconv.FormatBool(stat.Throttled()),
		TagInflight, strconv.Itoa(stat.InflightN()),
		TagNoMoreAttempt, strconv.FormatBool(noMoreAttempt))

	r.counter.Inc(FQNAppRequest, 1, appTags...)

	r.histogram.Observe(FQNAppCostMs, datetime.Milliseconds(stat.Cost()), appTags...)
}
