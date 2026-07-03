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
		TraceID:         "trace-1",
		RunID:           "run-1",
		WorkspaceID:     "ws-1",
		OperationName:   "chat",
		Provider:        "mock authorization=secret",
		Model:           "mock-chat",
		LatencyMS:       12,
		InputTokens:     3,
		OutputTokens:    5,
		ToolCount:       0,
		FailureCategory: "none",
		CreatedAt:       time.Unix(100, 0).UTC(),
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
