package observability_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/observability"
)

func TestMemorySinkRecordsSanitizedTelemetry(t *testing.T) {
	// telemetry 只能保存内部诊断标签，不能把 raw prompt、Authorization 或 provider payload 写入内存事件。
	sink := observability.NewMemorySink()
	err := sink.Record(context.Background(), observability.Event{
		TraceID:                 "trace-1",
		RunID:                   "run-1",
		WorkspaceID:             "ws-1",
		OperationName:           "chat",
		Provider:                "mock authorization=secret",
		Model:                   "mock-chat",
		LatencyMS:               12,
		InputTokens:             3,
		OutputTokens:            5,
		TotalTokens:             8,
		EstimatedCostMicrounits: 13,
		ToolCount:               0,
		FailureCategory:         "none",
		CreatedAt:               time.Unix(100, 0).UTC(),
	})
	if err != nil {
		t.Fatal(err)
	}

	events := sink.Events()
	if len(events) != 1 {
		t.Fatalf("events = %+v", events)
	}
	payload, err := json.Marshal(events[0])
	if err != nil {
		t.Fatal(err)
	}
	lower := strings.ToLower(string(payload))
	for _, forbidden := range []string{"authorization", "secret", "raw prompt", "provider_payload", "bearer "} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("telemetry leaked %q in %s", forbidden, payload)
		}
	}
	if events[0].Provider != "redacted" {
		t.Fatalf("provider = %q, want redacted", events[0].Provider)
	}
	if events[0].TotalTokens != 8 || events[0].EstimatedCostMicrounits != 13 {
		t.Fatalf("telemetry statistics changed unexpectedly: %+v", events[0])
	}
}

func TestMemorySinkClampsNegativeStatistics(t *testing.T) {
	// telemetry 统计来自外部 callback 链路时仍按不可信输入处理，负数不能进入诊断报表。
	sink := observability.NewMemorySink()
	err := sink.Record(context.Background(), observability.Event{
		LatencyMS:               -1,
		InputTokens:             -2,
		OutputTokens:            -3,
		TotalTokens:             -5,
		EstimatedCostMicrounits: -7,
		ToolCount:               -11,
	})
	if err != nil {
		t.Fatal(err)
	}
	event := sink.Events()[0]
	if event.LatencyMS != 0 ||
		event.InputTokens != 0 ||
		event.OutputTokens != 0 ||
		event.TotalTokens != 0 ||
		event.EstimatedCostMicrounits != 0 ||
		event.ToolCount != 0 {
		t.Fatalf("negative statistics were not clamped: %+v", event)
	}
}

func TestTokenCostRatesEstimateUsesConfiguredRatesOnly(t *testing.T) {
	// 成本估算只使用显式配置的单价；默认 0 表示只统计 token，不内置模型价格。
	rates := observability.TokenCostRates{
		InputMicrounitsPerToken:  3,
		OutputMicrounitsPerToken: 7,
	}
	if got := rates.Estimate(11, 13); got != 124 {
		t.Fatalf("estimated cost = %d", got)
	}
	if got := (observability.TokenCostRates{}).Estimate(11, 13); got != 0 {
		t.Fatalf("default estimated cost = %d, want 0", got)
	}
	if got := (observability.TokenCostRates{InputMicrounitsPerToken: -1, OutputMicrounitsPerToken: -1}).Estimate(-11, -13); got != 0 {
		t.Fatalf("negative estimated cost = %d, want 0", got)
	}
}

func TestMemorySinkRedactsUnsafeLabelVariants(t *testing.T) {
	// label 字段也按不可信输入处理；常见 token/key 变体、URL 和超长字符串都必须脱敏。
	sink := observability.NewMemorySink()
	longLabel := strings.Repeat("a", 129)
	err := sink.Record(context.Background(), observability.Event{
		TraceID:       "trace-token=abc",
		RunID:         longLabel,
		WorkspaceID:   "workspace/apikey/value",
		OperationName: "chat.model.generate",
		Provider:      "https://provider.example.test",
		Model:         "mock-chat",
		CreatedAt:     time.Unix(100, 0).UTC(),
	})
	if err != nil {
		t.Fatal(err)
	}

	event := sink.Events()[0]
	if event.TraceID != "redacted" ||
		event.RunID != "redacted" ||
		event.WorkspaceID != "redacted" ||
		event.Provider != "redacted" ||
		event.Model != "mock-chat" {
		t.Fatalf("unsafe labels were not redacted: %+v", event)
	}
}

func TestMemorySinkRedactsEachSecretLabelVariant(t *testing.T) {
	// 每个敏感变体都要有独立断言，避免后续删改 matcher 时测试仍然误通过。
	variants := []string{
		"api_key=abc",
		"api-key=abc",
		"apikey=abc",
		"x-api-key=abc",
		"access_token=abc",
		"refresh_token=abc",
		"client_secret=abc",
		"private_key=abc",
		"set-cookie=session",
		"cookie=session",
		"raw_payload={}",
		"raw payload",
		"provider_payload={}",
		"provider payload",
		"workspace/value",
	}
	for _, variant := range variants {
		t.Run(variant, func(t *testing.T) {
			sink := observability.NewMemorySink()
			if err := sink.Record(context.Background(), observability.Event{Provider: variant}); err != nil {
				t.Fatal(err)
			}
			if got := sink.Events()[0].Provider; got != "redacted" {
				t.Fatalf("provider label %q sanitized to %q, want redacted", variant, got)
			}
		})
	}
}
