package observability_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/observability"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components"
	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

func TestEinoCallbackHandlerRecordsModelTokenTelemetry(t *testing.T) {
	// Eino callback 只能把 token/latency 等内部诊断写入 telemetry，不能保存 prompt 或 completion 内容。
	sink := observability.NewMemorySink()
	handler := observability.NewEinoCallbackHandler(observability.EinoCallbackConfig{
		Sink:        sink,
		RunID:       "run-callback",
		WorkspaceID: "ws-callback",
		Provider:    "mock",
		Model:       "mock-chat",
		Now:         fixedCallbackClock(),
	})
	ctx := callbacks.InitCallbacks(context.Background(), &callbacks.RunInfo{
		Name:      "chat-model",
		Type:      "mock",
		Component: components.ComponentOfChatModel,
	}, handler)

	ctx = callbacks.OnStart(ctx, &einomodel.CallbackInput{
		Messages: []*schema.Message{schema.UserMessage("raw prompt must not appear")},
	})
	callbacks.OnEnd(ctx, &einomodel.CallbackOutput{
		Message: schema.AssistantMessage("raw completion must not appear", nil),
		TokenUsage: &einomodel.TokenUsage{
			PromptTokens:     7,
			CompletionTokens: 11,
			TotalTokens:      18,
		},
	})

	events := sink.Events()
	if len(events) != 1 {
		t.Fatalf("events = %+v", events)
	}
	event := events[0]
	if event.OperationName != "chat.model.generate" ||
		event.RunID != "run-callback" ||
		event.WorkspaceID != "ws-callback" ||
		event.Provider != "mock" ||
		event.Model != "mock-chat" ||
		event.InputTokens != 7 ||
		event.OutputTokens != 11 ||
		event.LatencyMS <= 0 {
		t.Fatalf("callback telemetry = %+v", event)
	}
	payload, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}
	lower := strings.ToLower(string(payload))
	for _, forbidden := range []string{"raw prompt", "raw completion", "must not appear"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("callback telemetry leaked %q: %s", forbidden, payload)
		}
	}
}

func TestEinoCallbackHandlerRecordsSafeErrorTelemetry(t *testing.T) {
	// callback error 只能落安全 failure category，不能把原始错误文本写入 telemetry。
	sink := observability.NewMemorySink()
	handler := observability.NewEinoCallbackHandler(observability.EinoCallbackConfig{
		Sink:        sink,
		RunID:       "run-callback-error",
		WorkspaceID: "ws-callback",
		Provider:    "mock",
		Model:       "mock-chat",
		Now:         fixedCallbackClock(),
	})
	ctx := callbacks.InitCallbacks(context.Background(), &callbacks.RunInfo{
		Name:      "chat-model",
		Type:      "mock",
		Component: components.ComponentOfChatModel,
	}, handler)

	ctx = callbacks.OnStart(ctx, &einomodel.CallbackInput{Messages: []*schema.Message{schema.UserMessage("raw prompt")}})
	callbacks.OnError(ctx, assertSafeError("provider raw body with token"))

	events := sink.Events()
	if len(events) != 1 {
		t.Fatalf("events = %+v", events)
	}
	if events[0].FailureCategory != "model_error" || events[0].LatencyMS <= 0 {
		t.Fatalf("error callback telemetry = %+v", events[0])
	}
	payload, err := json.Marshal(events[0])
	if err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"provider raw body", "with token", "raw prompt"} {
		if strings.Contains(strings.ToLower(string(payload)), forbidden) {
			t.Fatalf("error telemetry leaked %q: %s", forbidden, payload)
		}
	}
}

func TestEinoCallbackHandlerIgnoresNonModelCallbacks(t *testing.T) {
	// observability handler 当前只接 ChatModel callback，避免其它组件的 raw input/output 混入模型 telemetry。
	sink := observability.NewMemorySink()
	handler := observability.NewEinoCallbackHandler(observability.EinoCallbackConfig{
		Sink:        sink,
		RunID:       "run-non-model",
		WorkspaceID: "ws-callback",
		Provider:    "mock",
		Model:       "mock-chat",
	})
	ctx := callbacks.InitCallbacks(context.Background(), &callbacks.RunInfo{
		Name: "unknown",
		Type: "unknown",
	}, handler)

	ctx = callbacks.OnStart(ctx, "raw prompt")
	callbacks.OnEnd(ctx, "raw output")
	if events := sink.Events(); len(events) != 0 {
		t.Fatalf("non-model callback should be ignored: %+v", events)
	}
}

type assertSafeError string

func (err assertSafeError) Error() string {
	return string(err)
}

func fixedCallbackClock() func() time.Time {
	times := []time.Time{
		time.Unix(100, 0).UTC(),
		time.Unix(100, int64(25*time.Millisecond)).UTC(),
	}
	index := 0
	return func() time.Time {
		if index >= len(times) {
			return times[len(times)-1]
		}
		current := times[index]
		index++
		return current
	}
}
