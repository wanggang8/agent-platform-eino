package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

// OpenAICompatibleProviderConfig 保存 provider 私有密钥和 HTTP 客户端，不进入 Product Facts 或配置摘要。
type OpenAICompatibleProviderConfig struct {
	APIKey     string
	HTTPClient *http.Client
}

// OpenAICompatibleProvider 通过 OpenAI Chat Completions 兼容接口创建 ChatModel。
type OpenAICompatibleProvider struct {
	apiKey     string
	httpClient *http.Client
}

// NewOpenAICompatibleProvider 创建 OpenAI-compatible provider；api key 保留在 provider 私有边界。
func NewOpenAICompatibleProvider(config OpenAICompatibleProviderConfig) OpenAICompatibleProvider {
	client := config.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	return OpenAICompatibleProvider{
		apiKey:     config.APIKey,
		httpClient: client,
	}
}

// NewChatModel 校验配置和网络策略后返回非流式 Chat Completions 模型。
func (provider OpenAICompatibleProvider) NewChatModel(_ context.Context, cfg Config) (ChatModel, error) {
	if strings.TrimSpace(provider.apiKey) == "" {
		return nil, redactedProviderError("auth", "provider_config_missing", "模型凭据未配置", false)
	}
	if strings.TrimSpace(cfg.Model) == "" {
		return nil, redactedProviderError("config", "provider_config_invalid", "模型名称未配置", false)
	}
	if _, err := validateNetworkPolicy(cfg.BaseURL, cfg.NetworkSafety); err != nil {
		return nil, err
	}
	client := policyHTTPClient(provider.httpClient, cfg.NetworkSafety)
	if cfg.TimeoutMillis > 0 {
		client.Timeout = time.Duration(cfg.TimeoutMillis) * time.Millisecond
	}
	return openAICompatibleChatModel{
		baseURL: strings.TrimRight(cfg.BaseURL, "/"),
		model:   cfg.Model,
		apiKey:  provider.apiKey,
		client:  client,
	}, nil
}

type openAICompatibleChatModel struct {
	baseURL string
	model   string
	apiKey  string
	client  *http.Client
}

type chatCompletionRequest struct {
	Model    string                  `json:"model"`
	Messages []chatCompletionMessage `json:"messages"`
	Stream   bool                    `json:"stream"`
}

type chatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatCompletionMessage `json:"message"`
	} `json:"choices"`
}

// Generate 调用 OpenAI-compatible /chat/completions，并只返回 assistant 安全文本候选。
func (model openAICompatibleChatModel) Generate(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	payload := chatCompletionRequest{
		Model:    model.model,
		Messages: toChatCompletionMessages(req.Messages),
		Stream:   false,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return ChatResponse{}, redactedProviderError("config", "provider_config_invalid", "模型请求构造失败", false)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, model.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return ChatResponse{}, redactedProviderError("config", "provider_config_invalid", "模型请求地址无效", false)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+model.apiKey)

	resp, err := model.client.Do(httpReq)
	if err != nil {
		var redacted RedactedError
		if errors.As(err, &redacted) {
			return ChatResponse{}, redacted
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return ChatResponse{}, redactedProviderError("network", "provider_timeout", "模型请求超时", true)
		}
		return ChatResponse{}, redactedProviderError("network", "provider_unavailable", "模型服务暂不可用", true)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
		return ChatResponse{}, redactedProviderError("provider", "provider_unavailable", "模型服务返回错误", true)
	}

	var completion chatCompletionResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&completion); err != nil {
		return ChatResponse{}, redactedProviderError("provider", "provider_malformed_response", "模型响应格式无效", false)
	}
	if len(completion.Choices) == 0 || strings.TrimSpace(completion.Choices[0].Message.Content) == "" {
		return ChatResponse{}, redactedProviderError("provider", "provider_malformed_response", "模型响应缺少可展示内容", false)
	}
	return ChatResponse{Content: completion.Choices[0].Message.Content}, nil
}

// Stream 真实流式 provider 在后续阶段接入；当前避免半成品 streaming 绕过安全门。
func (model openAICompatibleChatModel) Stream(context.Context, ChatRequest) (ChatStream, error) {
	return nil, redactedProviderError("provider", "provider_stream_not_enabled", "模型流式输出尚未启用", false)
}

// toChatCompletionMessages 将项目安全消息映射为 OpenAI-compatible messages。
func toChatCompletionMessages(messages []Message) []chatCompletionMessage {
	result := make([]chatCompletionMessage, 0, len(messages))
	for _, message := range messages {
		result = append(result, chatCompletionMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}
	return result
}
