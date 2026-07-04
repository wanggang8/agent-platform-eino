package observability

import (
	"context"
	"strings"
	"sync"
	"time"
	"unicode"
)

const maxSafeLabelLength = 128

// Event 是内部 telemetry 事件，只允许保存安全标签和计数，不承载 prompt、raw provider body 或密钥。
type Event struct {
	TraceID       string
	RunID         string
	WorkspaceID   string
	OperationName string
	Provider      string
	Model         string
	LatencyMS     int64
	InputTokens   int
	OutputTokens  int
	TotalTokens   int
	// EstimatedCostMicrounits 是内部估算成本，单位由配置约定；默认 0 表示只统计 token。
	EstimatedCostMicrounits int64
	ToolCount               int
	FailureCategory         string
	CreatedAt               time.Time
}

// Sink 是内部观测写入接口；实现不得反向驱动 Workbench、Action API 或 Product Facts。
type Sink interface {
	Record(ctx context.Context, event Event) error
}

// NoopSink 是默认 telemetry sink，保证未配置观测后端时执行路径仍稳定。
type NoopSink struct{}

// Record 丢弃 telemetry 事件。
func (NoopSink) Record(context.Context, Event) error {
	return nil
}

// MemorySink 保存脱敏后的 telemetry 事件，主要用于测试和本地诊断。
type MemorySink struct {
	mu     sync.Mutex
	events []Event
}

// NewMemorySink 创建内存 telemetry sink。
func NewMemorySink() *MemorySink {
	return &MemorySink{}
}

// Record 写入一条已脱敏的内部 telemetry 事件。
func (sink *MemorySink) Record(_ context.Context, event Event) error {
	sink.mu.Lock()
	defer sink.mu.Unlock()
	sink.events = append(sink.events, sanitizeEvent(event))
	return nil
}

// Events 返回 telemetry 快照，避免调用方持有内部切片。
func (sink *MemorySink) Events() []Event {
	sink.mu.Lock()
	defer sink.mu.Unlock()
	events := make([]Event, len(sink.events))
	copy(events, sink.events)
	return events
}

func sanitizeEvent(event Event) Event {
	event.TraceID = safeLabel(event.TraceID)
	event.RunID = safeLabel(event.RunID)
	event.WorkspaceID = safeLabel(event.WorkspaceID)
	event.OperationName = safeLabel(event.OperationName)
	event.Provider = safeLabel(event.Provider)
	event.Model = safeLabel(event.Model)
	event.FailureCategory = safeLabel(event.FailureCategory)
	if event.LatencyMS < 0 {
		event.LatencyMS = 0
	}
	if event.InputTokens < 0 {
		event.InputTokens = 0
	}
	if event.OutputTokens < 0 {
		event.OutputTokens = 0
	}
	if event.TotalTokens < 0 {
		event.TotalTokens = 0
	}
	if event.EstimatedCostMicrounits < 0 {
		event.EstimatedCostMicrounits = 0
	}
	if event.ToolCount < 0 {
		event.ToolCount = 0
	}
	return event
}

func safeLabel(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if len([]rune(value)) > maxSafeLabelLength {
		return "redacted"
	}
	lower := strings.ToLower(value)
	for _, forbidden := range []string{
		"authorization", "api_key", "api-key", "api token", "api-token", "api_token", "apikey",
		"x-api-key", "access_token", "refresh_token", "token", "bearer ",
		"secret", "client_secret", "password", "private_key", "credential",
		"provider_payload", "provider payload", "providerpayload", "raw_payload", "raw payload", "rawpayload",
		"raw prompt", "raw provider", "raw body", "set-cookie", "cookie",
		"http://", "https://",
	} {
		if strings.Contains(lower, forbidden) {
			return "redacted"
		}
	}
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			continue
		}
		switch r {
		case '.', '_', ':', '-':
			continue
		default:
			return "redacted"
		}
	}
	return value
}
