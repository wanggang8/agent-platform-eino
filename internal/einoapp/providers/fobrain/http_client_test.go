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

func TestHTTPClientCurrentUserContextAllowsConfiguredSelfSignedCertificate(t *testing.T) {
	// 私有证书只允许通过显式配置跳过校验，且该行为限定在 Fobrain HTTP client 内。
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"display_name": "王五"},
		})
	}))
	defer server.Close()

	client, err := fobrain.NewHTTPClient(fobrain.HTTPClientConfig{
		BaseURL:            server.URL,
		Timeout:            time.Second,
		InsecureSkipVerify: true,
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
	if result.DisplayName != "王五" {
		t.Fatalf("current user result mismatch: %+v", result)
	}
}

func TestHTTPClientCurrentUserContextRejectsSelfSignedCertificateByDefault(t *testing.T) {
	// 默认 TLS 行为不能跳过证书校验，防止私有证书兼容开关扩大成全局不安全默认值。
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"display_name": "王五"},
		})
	}))
	defer server.Close()

	client, err := fobrain.NewHTTPClient(fobrain.HTTPClientConfig{
		BaseURL: server.URL,
		Timeout: time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.CurrentUserContext(context.Background(), fobrain.ResolvedCredential{
		WorkspaceID: "ws_fobrain",
		AuthParam:   "authorization",
		APIToken:    "workspace-token",
	})
	if !fobrain.HasReason(err, capabilities.PolicyReasonConnectorTransportUnavailable) {
		t.Fatalf("err = %v, want TLS transport failure", err)
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

func TestHTTPClientMyPermissionsUsesCurrentUserEndpoint(t *testing.T) {
	// my_permissions 从当前用户接口读取权限数组和数据范围摘要，不新增未经确认的外部 endpoint。
	var gotPath string
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("authorization")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"display_name":           "王五",
				"permissions":            []any{"资产只读", "漏洞只读"},
				"data_permission_names":  []any{"本部门数据"},
				"raw_policy_payload":     "must stay inside provider",
				"credential_debug_value": "must not be projected",
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

	result, err := client.MyPermissions(context.Background(), fobrain.ResolvedCredential{
		WorkspaceID: "ws_fobrain",
		AuthParam:   "authorization",
		APIToken:    "workspace-token",
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/api/v1/user" || gotAuth != "workspace-token" {
		t.Fatalf("request path/auth mismatch path=%q auth=%q", gotPath, gotAuth)
	}
	if strings.Join(result.PermissionNames, ",") != "资产只读,漏洞只读" ||
		strings.Join(result.DataPermissionNames, ",") != "本部门数据" {
		t.Fatalf("permissions result mismatch: %+v", result)
	}
}

func TestHTTPClientMyPermissionsIgnoresRawAPIPolicyFields(t *testing.T) {
	// API 路径、策略 code 和 raw payload 不是安全展示字段，不能进入权限 StructuredResult 材料。
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"display_name": "王五",
				"permissions": []any{
					map[string]any{"name": "资产只读", "code": "/api/v1/internal/assets"},
					map[string]any{"code": "/api/v1/internal/code-only"},
				},
				"apis":                   []any{map[string]any{"api": "/api/v1/internal/admin"}},
				"api_names":              []any{"raw:/api/v1/internal/report"},
				"data_permission_names":  []any{"本部门数据"},
				"raw_policy_payload":     "raw policy must stay inside provider",
				"credential_debug_value": "secret-like material",
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
	result, err := client.MyPermissions(context.Background(), fobrain.ResolvedCredential{
		WorkspaceID: "ws_fobrain",
		AuthParam:   "authorization",
		APIToken:    "workspace-token",
	})
	if err != nil {
		t.Fatal(err)
	}
	candidate, _ := fobrain.BuildMyPermissionsStructuredResult(result)
	for _, forbidden := range []string{"/api/v1/internal", "raw policy", "credential_debug_value", "secret-like"} {
		if strings.Contains(candidate.SafeSummary, forbidden) {
			t.Fatalf("permissions summary leaked %q: %s", forbidden, candidate.SafeSummary)
		}
	}
	if !strings.Contains(candidate.SafeSummary, "资产只读") || !strings.Contains(candidate.SafeSummary, "本部门数据") {
		t.Fatalf("permissions summary missing safe fields: %s", candidate.SafeSummary)
	}
}
