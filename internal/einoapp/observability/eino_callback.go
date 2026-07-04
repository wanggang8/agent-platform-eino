package observability

import (
	"context"
	"time"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components"
	einomodel "github.com/cloudwego/eino/components/model"
)

type callbackStartKey struct{}

// EinoCallbackConfig 定义 Eino callback 到内部 telemetry 的安全映射参数。
type EinoCallbackConfig struct {
	Sink        Sink
	RunID       string
	WorkspaceID string
	Provider    string
	Model       string
	Now         func() time.Time
}

// NewEinoCallbackHandler 创建 Eino callback handler，只写内部 telemetry，不写 Product Facts 或产品 SSE。
func NewEinoCallbackHandler(config EinoCallbackConfig) callbacks.Handler {
	sink := config.Sink
	if sink == nil {
		sink = NoopSink{}
	}
	now := config.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return callbacks.NewHandlerBuilder().
		OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, _ callbacks.CallbackInput) context.Context {
			if !chatModelCallback(info) {
				return ctx
			}
			return context.WithValue(ctx, callbackStartKey{}, now())
		}).
		OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
			if !chatModelCallback(info) {
				return ctx
			}
			_ = sink.Record(ctx, callbackTelemetryEvent(config, callbackLatency(ctx, now()), tokenUsage(output), "", now()))
			return ctx
		}).
		OnErrorFn(func(ctx context.Context, info *callbacks.RunInfo, _ error) context.Context {
			if !chatModelCallback(info) {
				return ctx
			}
			_ = sink.Record(ctx, callbackTelemetryEvent(config, callbackLatency(ctx, now()), einomodel.TokenUsage{}, "model_error", now()))
			return ctx
		}).
		Build()
}

func chatModelCallback(info *callbacks.RunInfo) bool {
	return info != nil && info.Component == components.ComponentOfChatModel
}

func callbackLatency(ctx context.Context, now time.Time) int64 {
	start, ok := ctx.Value(callbackStartKey{}).(time.Time)
	if !ok {
		return 0
	}
	latency := now.Sub(start).Milliseconds()
	if latency <= 0 {
		return 1
	}
	return latency
}

func tokenUsage(output callbacks.CallbackOutput) einomodel.TokenUsage {
	modelOutput := einomodel.ConvCallbackOutput(output)
	if modelOutput == nil || modelOutput.TokenUsage == nil {
		return einomodel.TokenUsage{}
	}
	return *modelOutput.TokenUsage
}

func callbackTelemetryEvent(config EinoCallbackConfig, latencyMS int64, usage einomodel.TokenUsage, failureCategory string, createdAt time.Time) Event {
	return Event{
		TraceID:         config.RunID,
		RunID:           config.RunID,
		WorkspaceID:     config.WorkspaceID,
		OperationName:   "chat.model.generate",
		Provider:        config.Provider,
		Model:           config.Model,
		LatencyMS:       latencyMS,
		InputTokens:     usage.PromptTokens,
		OutputTokens:    usage.CompletionTokens,
		FailureCategory: failureCategory,
		CreatedAt:       createdAt,
	}
}
