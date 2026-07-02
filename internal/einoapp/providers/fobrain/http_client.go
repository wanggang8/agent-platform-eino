package fobrain

import (
	"context"
	"crypto/tls"
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
// HTTPClientConfig 不持有 token；token 只从 ResolvedCredential 传入单次请求。
type HTTPClientConfig struct {
	BaseURL            string
	Timeout            time.Duration
	HTTPClient         *http.Client
	InsecureSkipVerify bool
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
		httpClient = &http.Client{
			Timeout:   timeout,
			Transport: httpTransport(config.InsecureSkipVerify),
		}
	}
	return &HTTPClient{baseURL: parsed, timeout: timeout, httpClient: httpClient}, nil
}

// httpTransport 仅在 Fobrain 私有证书配置明确开启时跳过证书链校验。
func httpTransport(insecureSkipVerify bool) http.RoundTripper {
	if !insecureSkipVerify {
		return http.DefaultTransport
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return transport
}

// CurrentUserContext 调用 Fobrain 标准当前用户接口，并只返回安全展示字段。
func (client *HTTPClient) CurrentUserContext(ctx context.Context, credential ResolvedCredential) (CurrentUserContextResult, error) {
	item, err := client.currentUserMap(ctx, credential)
	if err != nil {
		return CurrentUserContextResult{}, err
	}
	result := currentUserResultFromMap(item)
	if strings.TrimSpace(result.DisplayName) == "" && strings.TrimSpace(result.Department) == "" && strings.TrimSpace(result.Role) == "" {
		return CurrentUserContextResult{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 当前用户响应缺少安全字段")
	}
	return result, nil
}

// MyPermissions 读取当前用户权限范围；当前真实证据来自 /api/v1/user 的安全字段。
func (client *HTTPClient) MyPermissions(ctx context.Context, credential ResolvedCredential) (MyPermissionsResult, error) {
	item, err := client.currentUserMap(ctx, credential)
	if err != nil {
		return MyPermissionsResult{}, err
	}
	result := myPermissionsResultFromMap(item)
	if len(result.PermissionNames) == 0 && len(result.DataPermissionNames) == 0 {
		return MyPermissionsResult{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 权限响应缺少安全字段")
	}
	return result, nil
}

// currentUserMap 执行当前用户 HTTP 请求，并把 raw payload 限制在 provider 边界内。
func (client *HTTPClient) currentUserMap(ctx context.Context, credential ResolvedCredential) (map[string]any, error) {
	if client == nil || client.httpClient == nil || client.baseURL == nil {
		return nil, NewSafeError(capabilities.PolicyReasonConnectorTransportUnavailable, "Fobrain client 未配置")
	}
	if strings.TrimSpace(credential.APIToken) == "" {
		return nil, NewSafeError(capabilities.PolicyReasonCredentialMissing, "Fobrain 凭据未配置")
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
		return nil, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 当前用户请求创建失败")
	}
	req.Header.Set(authParam, credential.APIToken)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		if errors.Is(requestCtx.Err(), context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, NewSafeError(capabilities.PolicyReasonConnectorTimeout, "Fobrain 当前用户请求超时")
		}
		return nil, NewSafeError(capabilities.PolicyReasonConnectorTransportUnavailable, "Fobrain 当前用户网络不可用")
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, NewSafeError(capabilities.PolicyReasonConnectorAuthFailure, "Fobrain 当前用户认证失败")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 当前用户读取失败")
	}
	var payload any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 当前用户响应不是合法 JSON")
	}
	item, ok := unwrapPayloadMap(payload)
	if !ok {
		return nil, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 当前用户响应结构无效")
	}
	return item, nil
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

// myPermissionsResultFromMap 只读取权限和数据范围安全字段，忽略 raw policy payload。
func myPermissionsResultFromMap(item map[string]any) MyPermissionsResult {
	return MyPermissionsResult{
		DisplayName: firstString(item, "display_name", "name", "username", "staff_name"),
		PermissionNames: stringListFromAny(firstPresentValue(item,
			"permissions", "menu_names", "menus")),
		DataPermissionNames: stringListFromAny(firstPresentValue(item,
			"data_permission_names", "data_permissions", "data_permission", "dataPermission")),
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

func firstPresentValue(item map[string]any, keys ...string) any {
	for _, key := range keys {
		if value, ok := item[key]; ok {
			return value
		}
	}
	return nil
}

func stringListFromAny(value any) []string {
	switch typed := value.(type) {
	case []any:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			out = append(out, stringListFromAny(item)...)
		}
		return out
	case []string:
		return cleanStringList(typed)
	case map[string]any:
		return cleanStringList([]string{firstString(typed, "name", "label", "title")})
	case string:
		if strings.Contains(typed, ",") {
			return cleanStringList(strings.Split(typed, ","))
		}
		return cleanStringList([]string{typed})
	default:
		return nil
	}
}

func cleanStringList(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		text := strings.TrimSpace(value)
		if text == "" || unsafePermissionText(text) {
			continue
		}
		if _, ok := seen[text]; ok {
			continue
		}
		seen[text] = struct{}{}
		out = append(out, text)
	}
	return out
}

func unsafePermissionText(value string) bool {
	lower := strings.ToLower(strings.TrimSpace(value))
	for _, marker := range []string{"/api/", "raw", "credential", "secret", "token", "authorization", "bearer"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
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
