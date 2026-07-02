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
		BaseURL:    server.URL + "/api",
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

func TestHTTPClientCurrentUserContextSupportsAPIBasePath(t *testing.T) {
	// 真实私有环境可能把 base_url 配到 /api；client 必须归一化到 /api/v1/user。
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"display_name": "王五"},
		})
	}))
	defer server.Close()

	client, err := fobrain.NewHTTPClient(fobrain.HTTPClientConfig{
		BaseURL:    server.URL + "/api",
		Timeout:    time.Second,
		HTTPClient: server.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.CurrentUserContext(context.Background(), fobrain.ResolvedCredential{
		WorkspaceID: "ws_fobrain",
		AuthParam:   "authorization",
		APIToken:    "workspace-token",
	}); err != nil {
		t.Fatal(err)
	}
	if gotPath != "/api/v1/user" {
		t.Fatalf("path = %q, want /api/v1/user", gotPath)
	}
}

func TestHTTPClientCurrentUserContextFallsBackToAPIUserPath(t *testing.T) {
	// 部分私有部署只暴露 /api/user；标准路径 404 时允许回退，但仍不能泄漏响应 body。
	var gotPaths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPaths = append(gotPaths, r.URL.Path)
		if r.URL.Path == "/api/v1/user" {
			http.NotFound(w, r)
			return
		}
		if r.URL.Path != "/api/user" {
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"display_name": "王五"},
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
	if result.DisplayName != "王五" {
		t.Fatalf("current user result mismatch: %+v", result)
	}
	if strings.Join(gotPaths, ",") != "/api/v1/user,/api/user" {
		t.Fatalf("paths = %v, want standard then fallback", gotPaths)
	}
}

func TestHTTPClientCurrentUserContextDoesNotFallbackAfterAuthFailure(t *testing.T) {
	// 认证/权限失败是终止错误，不能 fallback 到兼容路径掩盖凭据问题。
	var gotPaths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPaths = append(gotPaths, r.URL.Path)
		http.Error(w, "denied", http.StatusUnauthorized)
	}))
	defer server.Close()

	client, err := fobrain.NewHTTPClient(fobrain.HTTPClientConfig{
		BaseURL:    server.URL + "/api",
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
	if strings.Join(gotPaths, ",") != "/api/v1/user" {
		t.Fatalf("paths = %v, want no fallback after auth failure", gotPaths)
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

func TestHTTPClientMyPermissionsAllowsEmptyPermissionFields(t *testing.T) {
	// Batch A 必须覆盖权限空态：真实当前用户接口可能只返回身份，不返回权限数组。
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"display_name": "王五"},
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
	if !strings.Contains(candidate.SafeSummary, "未返回权限字段") {
		t.Fatalf("permissions empty summary mismatch: %s", candidate.SafeSummary)
	}
}

func TestHTTPClientMyPermissionsRejectsEmptyCurrentUserPayload(t *testing.T) {
	// 权限空态只能建立在已识别当前用户上，不能把空对象或错结构响应误判为通过。
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{}})
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
	_, err = client.MyPermissions(context.Background(), fobrain.ResolvedCredential{
		WorkspaceID: "ws_fobrain",
		AuthParam:   "authorization",
		APIToken:    "workspace-token",
	})
	if !fobrain.HasReason(err, capabilities.PolicyReasonConnectorExecutionFailed) {
		t.Fatalf("err = %v, want invalid current user payload", err)
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
