package fobrain

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"agent-platform-eino/internal/einoapp/capabilities"
)

const currentUserPath = "/api/v1/user"

// HTTPClientConfig 保存真实 Fobrain HTTP client 的 provider 边界配置。
// BaseURL 和 token 都不得进入 Product Facts；token 只从 ResolvedCredential 传入单次请求。
type HTTPClientConfig struct {
	BaseURL    string
	Timeout    time.Duration
	HTTPClient *http.Client
}

// HTTPClient 是 Fobrain live read 的最小 HTTP 实现。
type HTTPClient struct {
	baseURL    *url.URL
	timeout    time.Duration
	httpClient *http.Client
}

var _ FobrainClient = (*HTTPClient)(nil)

// NewHTTPClient 创建 live HTTP client；只校验连接入口形态，不持有 token。
func NewHTTPClient(config HTTPClientConfig) (*HTTPClient, error) {
	parsed, err := url.Parse(strings.TrimSpace(config.BaseURL))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain base url 无效")
	}
	timeout := config.Timeout
	if timeout <= 0 {
		timeout = currentUserContextTimeout
	}
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: timeout}
	}
	return &HTTPClient{baseURL: parsed, timeout: timeout, httpClient: httpClient}, nil
}

// CurrentUserContext 调用 Fobrain 标准当前用户接口，并只返回安全展示字段。
func (client *HTTPClient) CurrentUserContext(ctx context.Context, credential ResolvedCredential) (CurrentUserContextResult, error) {
	if client == nil || client.httpClient == nil || client.baseURL == nil {
		return CurrentUserContextResult{}, NewSafeError(capabilities.PolicyReasonConnectorTransportUnavailable, "Fobrain client 未配置")
	}
	if strings.TrimSpace(credential.APIToken) == "" {
		return CurrentUserContextResult{}, NewSafeError(capabilities.PolicyReasonCredentialMissing, "Fobrain 凭据未配置")
	}
	authParam := strings.TrimSpace(credential.AuthParam)
	if authParam == "" {
		authParam = "authorization"
	}
	endpoint := client.endpoint(currentUserPath)
	requestCtx := ctx
	cancel := func() {}
	if _, ok := ctx.Deadline(); !ok && client.timeout > 0 {
		requestCtx, cancel = context.WithTimeout(ctx, client.timeout)
	}
	defer cancel()

	req, err := http.NewRequestWithContext(requestCtx, http.MethodGet, endpoint, nil)
	if err != nil {
		return CurrentUserContextResult{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 当前用户请求创建失败")
	}
	req.Header.Set(authParam, credential.APIToken)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		if errors.Is(requestCtx.Err(), context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return CurrentUserContextResult{}, NewSafeError(capabilities.PolicyReasonConnectorTimeout, "Fobrain 当前用户请求超时")
		}
		return CurrentUserContextResult{}, NewSafeError(capabilities.PolicyReasonConnectorTransportUnavailable, "Fobrain 当前用户网络不可用")
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return CurrentUserContextResult{}, NewSafeError(capabilities.PolicyReasonConnectorAuthFailure, "Fobrain 当前用户认证失败")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return CurrentUserContextResult{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 当前用户读取失败")
	}
	var payload any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return CurrentUserContextResult{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 当前用户响应不是合法 JSON")
	}
	item, ok := unwrapPayloadMap(payload)
	if !ok {
		return CurrentUserContextResult{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 当前用户响应结构无效")
	}
	result := currentUserResultFromMap(item)
	if strings.TrimSpace(result.DisplayName) == "" && strings.TrimSpace(result.Department) == "" && strings.TrimSpace(result.Role) == "" {
		return CurrentUserContextResult{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 当前用户响应缺少安全字段")
	}
	return result, nil
}

// endpoint 兼容 base_url 指向服务根路径或已经包含 /api/v1 的本地配置。
func (client *HTTPClient) endpoint(path string) string {
	endpoint := *client.baseURL
	basePath := strings.TrimRight(endpoint.Path, "/")
	if strings.HasSuffix(basePath, "/api/v1") && path == currentUserPath {
		endpoint.Path = basePath + "/user"
	} else {
		endpoint.Path = basePath + path
	}
	endpoint.RawQuery = ""
	endpoint.Fragment = ""
	return endpoint.String()
}

// unwrapPayloadMap 支持常见 {data:{...}} 包装，也支持顶层对象直接作为业务数据。
func unwrapPayloadMap(payload any) (map[string]any, bool) {
	item, ok := payload.(map[string]any)
	if !ok {
		return nil, false
	}
	if data, ok := item["data"].(map[string]any); ok {
		return data, true
	}
	return item, true
}

// currentUserResultFromMap 把 provider raw 字段收敛为当前用户安全摘要材料。
func currentUserResultFromMap(item map[string]any) CurrentUserContextResult {
	return CurrentUserContextResult{
		DisplayName: firstString(item, "display_name", "name", "username", "staff_name"),
		Department:  firstString(item, "department_name", "department"),
		Role:        roleName(item["role"], firstString(item, "role_name")),
	}
}

func firstString(item map[string]any, keys ...string) string {
	for _, key := range keys {
		value := stringFromAny(item[key])
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func roleName(value any, fallback string) string {
	if text := stringFromAny(value); strings.TrimSpace(text) != "" {
		return strings.TrimSpace(text)
	}
	if item, ok := value.(map[string]any); ok {
		return firstString(item, "name", "label", "title")
	}
	return strings.TrimSpace(fallback)
}

func stringFromAny(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case json.Number:
		return typed.String()
	default:
		return ""
	}
}
