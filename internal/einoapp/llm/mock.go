package llm

import (
	"context"
	"io"
)

// MockProvider 是测试和早期 Phase 使用的模型 provider，不进行网络访问。
type MockProvider struct {
	response string
	usage    TokenUsage
}

type mockChatModel struct {
	response string
	usage    TokenUsage
}

type mockStream struct {
	response string
	sent     bool
}

// NewMockProvider 创建固定响应的模型 provider。
func NewMockProvider(response string) MockProvider {
	return MockProvider{response: response}
}

// NewMockProviderWithUsage 创建带安全 usage 计数的 mock provider，用于 callback/telemetry 回归测试。
func NewMockProviderWithUsage(response string, usage TokenUsage) MockProvider {
	return MockProvider{response: response, usage: usage}
}

// NewChatModel 返回固定响应模型，保持与真实 provider 相同接口。
func (provider MockProvider) NewChatModel(_ context.Context, _ Config) (ChatModel, error) {
	return mockChatModel{response: provider.response, usage: provider.usage}, nil
}

// Generate 返回预设响应，用于 execution 单元测试。
func (model mockChatModel) Generate(_ context.Context, _ ChatRequest) (ChatResponse, error) {
	return ChatResponse{Content: model.response, Usage: model.usage}, nil
}

// Stream 返回只发送一次的测试流。
func (model mockChatModel) Stream(_ context.Context, _ ChatRequest) (ChatStream, error) {
	return &mockStream{response: model.response}, nil
}

// Recv 第一次返回响应，之后返回 io.EOF。
func (stream *mockStream) Recv() (ChatResponse, error) {
	if stream.sent {
		return ChatResponse{}, io.EOF
	}
	stream.sent = true
	return ChatResponse{Content: stream.response}, nil
}

// Close 满足 ChatStream 接口，mock 不持有外部资源。
func (stream *mockStream) Close() error {
	return nil
}
