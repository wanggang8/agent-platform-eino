package llm

import (
	"context"
	"io"
)

type MockProvider struct {
	response string
}

type mockChatModel struct {
	response string
}

type mockStream struct {
	response string
	sent     bool
}

func NewMockProvider(response string) MockProvider {
	return MockProvider{response: response}
}

func (provider MockProvider) NewChatModel(_ context.Context, _ Config) (ChatModel, error) {
	return mockChatModel{response: provider.response}, nil
}

func (model mockChatModel) Generate(_ context.Context, _ ChatRequest) (ChatResponse, error) {
	return ChatResponse{Content: model.response}, nil
}

func (model mockChatModel) Stream(_ context.Context, _ ChatRequest) (ChatStream, error) {
	return &mockStream{response: model.response}, nil
}

func (stream *mockStream) Recv() (ChatResponse, error) {
	if stream.sent {
		return ChatResponse{}, io.EOF
	}
	stream.sent = true
	return ChatResponse{Content: stream.response}, nil
}

func (stream *mockStream) Close() error {
	return nil
}
