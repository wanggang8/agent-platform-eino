package fobrain_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/providers/fobrain"
)

func TestHTTPClientCurrentUserContextUsesConfiguredAuthorizationHeader(t *testing.T) {
	// live client 必须只在 provider 边界内持有 token，并按配置 header 名访问标准当前用户接口。
	var gotPath string
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("authorization")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 0,
			"data": map[string]any{
				"display_name":    "王五",
				"department_name": "安全部",
				"role":            map[string]any{"name": "安全运营"},
			},
		})
	}))
	defer server.Close()

	client, err := fobrain.NewHTTPClient(fobrain.HTTPClientConfig{
		BaseURL:    server.URL,
		Timeout:    time.Second,
		HTTPClient: server.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}

	result, err := client.CurrentUserContext(context.Background(), fobrain.ResolvedCredential{
		WorkspaceID: "ws_fobrain",
		AuthParam:   "authorization",
		APIToken:    "workspace-token",
	})
	if err != nil {
		t.Fatal(err)
	}

	if gotPath != "/api/v1/user" {
		t.Fatalf("path = %q, want /api/v1/user", gotPath)
	}
	if gotAuth != "workspace-token" {
		t.Fatalf("authorization header = %q", gotAuth)
	}
	if result.DisplayName != "王五" || result.Department != "安全部" || result.Role != "安全运营" {
		t.Fatalf("current user result mismatch: %+v", result)
	}
}

func TestHTTPClientCurrentUserContextFoldsErrorsWithoutLeakingToken(t *testing.T) {
	// provider 错误只能暴露稳定原因和安全摘要，不能把 header/token/raw body 拼出去。
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"error":"workspace-token Authorization raw body"}`, http.StatusUnauthorized)
	}))
	defer server.Close()

	client, err := fobrain.NewHTTPClient(fobrain.HTTPClientConfig{
		BaseURL:    server.URL,
		Timeout:    time.Second,
		HTTPClient: server.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.CurrentUserContext(context.Background(), fobrain.ResolvedCredential{
		WorkspaceID: "ws_fobrain",
		AuthParam:   "authorization",
		APIToken:    "workspace-token",
	})
	if !fobrain.HasReason(err, capabilities.PolicyReasonConnectorAuthFailure) {
		t.Fatalf("err = %v, want auth failure", err)
	}
	for _, forbidden := range []string{"workspace-token", "Authorization", "raw body"} {
		if strings.Contains(err.Error(), forbidden) {
			t.Fatalf("safe error leaked %q: %s", forbidden, err.Error())
		}
	}
}
