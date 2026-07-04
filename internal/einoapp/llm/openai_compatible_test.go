package llm_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"agent-platform-eino/internal/einoapp/llm"
)

func TestOpenAICompatibleGeneratePostsChatCompletion(t *testing.T) {
	// OpenAI-compatible provider 只替换模型边界，请求体来自安全 messages 投影。
	var capturedAuth string
	var capturedBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuth = r.Header.Get("Authorization")
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&capturedBody); err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"真实模型回答"}}]}`))
	}))
	defer server.Close()

	model := newTestOpenAIModel(t, server.URL, "sk-test")
	resp, err := model.Generate(context.Background(), llm.ChatRequest{
		Messages: []llm.Message{
			{Role: "system", Content: "安全系统提示"},
			{Role: "user", Content: "你好"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if resp.Content != "真实模型回答" {
		t.Fatalf("content = %q", resp.Content)
	}
	if capturedAuth != "Bearer sk-test" {
		t.Fatalf("authorization header = %q", capturedAuth)
	}
	if capturedBody["model"] != "real-chat" {
		t.Fatalf("model = %+v", capturedBody["model"])
	}
	messages, ok := capturedBody["messages"].([]any)
	if !ok || len(messages) != 2 {
		t.Fatalf("messages = %+v", capturedBody["messages"])
	}
}

func TestOpenAICompatibleGenerateMapsUsage(t *testing.T) {
	// provider usage 只能作为安全计数进入 LLM 响应，不能携带 raw response 或凭据。
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"choices":[{"message":{"role":"assistant","content":"带 usage 的回答"}}],
			"usage":{"prompt_tokens":13,"completion_tokens":21,"total_tokens":34}
		}`))
	}))
	defer server.Close()

	model := newTestOpenAIModel(t, server.URL, "sk-test")
	resp, err := model.Generate(context.Background(), llm.ChatRequest{Messages: []llm.Message{{Role: "user", Content: "hello"}}})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Usage.InputTokens != 13 || resp.Usage.OutputTokens != 21 || resp.Usage.TotalTokens != 34 {
		t.Fatalf("usage = %+v", resp.Usage)
	}
}

func TestOpenAICompatibleGenerateNormalizesNegativeUsage(t *testing.T) {
	// usage 属于 provider 外部输入，进入项目安全响应前必须压掉不可信负数。
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"choices":[{"message":{"role":"assistant","content":"带异常 usage 的回答"}}],
			"usage":{"prompt_tokens":-13,"completion_tokens":21,"total_tokens":-34}
		}`))
	}))
	defer server.Close()

	model := newTestOpenAIModel(t, server.URL, "sk-test")
	resp, err := model.Generate(context.Background(), llm.ChatRequest{Messages: []llm.Message{{Role: "user", Content: "hello"}}})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Usage.InputTokens != 0 || resp.Usage.OutputTokens != 21 || resp.Usage.TotalTokens != 0 {
		t.Fatalf("usage = %+v", resp.Usage)
	}
}

func TestOpenAICompatibleNon2xxReturnsRedactedError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "Authorization: Bearer sk-secret raw body", http.StatusUnauthorized)
	}))
	defer server.Close()

	model := newTestOpenAIModel(t, server.URL, "sk-secret")
	_, err := model.Generate(context.Background(), llm.ChatRequest{Messages: []llm.Message{{Role: "user", Content: "hello"}}})
	if err == nil {
		t.Fatal("Generate error = nil")
	}
	var redacted llm.RedactedError
	if !errors.As(err, &redacted) {
		t.Fatalf("err type = %T", err)
	}
	if redacted.ReasonCode != "provider_unavailable" {
		t.Fatalf("reason_code = %q", redacted.ReasonCode)
	}
	for _, forbidden := range []string{"sk-secret", "Authorization", "raw body"} {
		if strings.Contains(err.Error(), forbidden) {
			t.Fatalf("redacted error leaked %q: %v", forbidden, err)
		}
	}
}

func TestOpenAICompatibleMalformedResponseReturnsRedactedError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"choices":[]}`))
	}))
	defer server.Close()

	model := newTestOpenAIModel(t, server.URL, "sk-test")
	_, err := model.Generate(context.Background(), llm.ChatRequest{Messages: []llm.Message{{Role: "user", Content: "hello"}}})
	if err == nil {
		t.Fatal("Generate error = nil")
	}
	var redacted llm.RedactedError
	if !errors.As(err, &redacted) || redacted.ReasonCode != "provider_malformed_response" {
		t.Fatalf("err = %+v", err)
	}
}

func TestOpenAICompatibleRequiresAPIKey(t *testing.T) {
	provider := llm.NewOpenAICompatibleProvider(llm.OpenAICompatibleProviderConfig{})

	_, err := provider.NewChatModel(context.Background(), llm.Config{
		Provider:      "openai_compatible",
		BaseURL:       "https://llm.example.test/v1",
		Model:         "real-chat",
		TimeoutMillis: 1000,
	})
	if err == nil {
		t.Fatal("NewChatModel missing key error = nil")
	}
	var redacted llm.RedactedError
	if !errors.As(err, &redacted) || redacted.ReasonCode != "provider_config_missing" {
		t.Fatalf("err = %+v", err)
	}
}

func TestNetworkPolicyBlocksUnallowedHost(t *testing.T) {
	provider := llm.NewOpenAICompatibleProvider(llm.OpenAICompatibleProviderConfig{APIKey: "sk-test"})

	_, err := provider.NewChatModel(context.Background(), llm.Config{
		Provider:      "openai_compatible",
		BaseURL:       "https://blocked.example.test/v1",
		Model:         "real-chat",
		TimeoutMillis: 1000,
		NetworkSafety: llm.NetworkSafety{
			RequireHTTPS:   true,
			AllowedHosts:   []string{"allowed.example.test"},
			AllowLocalHTTP: false,
		},
	})
	if err == nil {
		t.Fatal("NewChatModel unallowed host error = nil")
	}
	var redacted llm.RedactedError
	if !errors.As(err, &redacted) || redacted.ReasonCode != "provider_host_not_allowed" {
		t.Fatalf("err = %+v", err)
	}
}

func TestNetworkPolicyRejectsURLUserinfo(t *testing.T) {
	provider := llm.NewOpenAICompatibleProvider(llm.OpenAICompatibleProviderConfig{APIKey: "sk-test"})

	_, err := provider.NewChatModel(context.Background(), llm.Config{
		Provider:      "openai_compatible",
		BaseURL:       "https://user:pass@llm.example.test/v1",
		Model:         "real-chat",
		TimeoutMillis: 1000,
		NetworkSafety: llm.NetworkSafety{
			RequireHTTPS: true,
		},
	})
	if err == nil {
		t.Fatal("NewChatModel userinfo error = nil")
	}
	var redacted llm.RedactedError
	if !errors.As(err, &redacted) || redacted.ReasonCode != "provider_config_invalid" {
		t.Fatalf("err = %+v", err)
	}
}

func TestOpenAICompatibleBlocksRedirect(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"redirected"}}]}`))
	}))
	defer target.Close()
	redirector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL+"/v1/chat/completions", http.StatusTemporaryRedirect)
	}))
	defer redirector.Close()

	model := newTestOpenAIModel(t, redirector.URL, "sk-test")
	_, err := model.Generate(context.Background(), llm.ChatRequest{Messages: []llm.Message{{Role: "user", Content: "hello"}}})
	if err == nil {
		t.Fatal("Generate redirect error = nil")
	}
	var redacted llm.RedactedError
	if !errors.As(err, &redacted) || redacted.ReasonCode != "provider_redirect_blocked" {
		t.Fatalf("err = %+v", err)
	}
}

func TestNetworkPolicyDialBlocksPrivateDNSResult(t *testing.T) {
	provider := llm.NewOpenAICompatibleProvider(llm.OpenAICompatibleProviderConfig{APIKey: "sk-test"})
	model, err := provider.NewChatModel(context.Background(), llm.Config{
		Provider:      "openai_compatible",
		BaseURL:       "https://localhost:1/v1",
		Model:         "real-chat",
		TimeoutMillis: 1000,
		NetworkSafety: llm.NetworkSafety{
			RequireHTTPS:         true,
			AllowLocalHTTP:       false,
			BlockPrivateNetworks: true,
			AllowedHosts:         []string{"localhost"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = model.Generate(context.Background(), llm.ChatRequest{Messages: []llm.Message{{Role: "user", Content: "hello"}}})
	if err == nil {
		t.Fatal("Generate private DNS error = nil")
	}
	var redacted llm.RedactedError
	if !errors.As(err, &redacted) || redacted.ReasonCode != "provider_network_blocked" {
		t.Fatalf("err = %+v", err)
	}
}

func newTestOpenAIModel(t *testing.T, baseURL string, apiKey string) llm.ChatModel {
	t.Helper()

	provider := llm.NewOpenAICompatibleProvider(llm.OpenAICompatibleProviderConfig{APIKey: apiKey, HTTPClient: http.DefaultClient})
	model, err := provider.NewChatModel(context.Background(), llm.Config{
		Provider:      "openai_compatible",
		BaseURL:       baseURL + "/v1",
		Model:         "real-chat",
		TimeoutMillis: 1000,
		NetworkSafety: llm.NetworkSafety{
			RequireHTTPS:   false,
			AllowLocalHTTP: true,
			AllowedHosts:   []string{"127.0.0.1"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return model
}
