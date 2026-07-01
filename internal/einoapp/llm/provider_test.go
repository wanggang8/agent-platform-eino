package llm_test

import (
	"context"
	"errors"
	"testing"

	"agent-platform-eino/internal/einoapp/llm"
)

func TestMockProviderReturnsConfiguredResponse(t *testing.T) {
	provider := llm.NewMockProvider("fixture answer")

	model, err := provider.NewChatModel(context.Background(), llm.Config{
		Provider: "mock",
		Model:    "mock-chat",
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := model.Generate(context.Background(), llm.ChatRequest{
		Messages: []llm.Message{{Role: "user", Content: "hello"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Content != "fixture answer" {
		t.Fatalf("content = %q", resp.Content)
	}
}

func TestRedactedProviderErrorDoesNotExposeSecretsOrRawBody(t *testing.T) {
	// provider 错误进入产品层前必须脱敏，不能泄露 token、鉴权头或 raw body。
	err := llm.RedactProviderError(errors.New("Authorization: Bearer sk-secret raw body: forbidden"), llm.RedactionInput{
		Category:       "auth",
		ReasonCode:     "provider_config_invalid",
		SafeSummary:    "模型配置不可用",
		Retryable:      false,
		CorrelationRef: "provider-error-1",
	})

	encoded := err.Error()
	for _, forbidden := range []string{"sk-secret", "Authorization", "forbidden"} {
		if contains(encoded, forbidden) {
			t.Fatalf("redacted error leaked %q: %s", forbidden, encoded)
		}
	}
	if err.SchemaVersion != "eino.provider_redacted_error.v1" {
		t.Fatalf("schema_version = %q", err.SchemaVersion)
	}
	if !err.Redaction.RemovedSecret || !err.Redaction.RemovedAuthHeader || !err.Redaction.RemovedRawBody {
		t.Fatalf("redaction flags not set: %+v", err.Redaction)
	}
}

func contains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
