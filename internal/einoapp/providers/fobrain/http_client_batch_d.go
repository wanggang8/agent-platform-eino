package fobrain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/facts"
)

const (
	currentAssetListPath         = "/api/asset"
	legacyAssetListPath          = "/api/v1/asset"
	currentThreatCenterPath      = "/api/threat_center"
	legacyThreatCenterPath       = "/api/v1/threat_center"
	defaultThreatCenterDataRange = "4"
)

// ParameterizedQuery 调用当前真实 Fobrain 参数化只读接口，并把 raw payload 收敛成安全行摘要。
func (client *HTTPClient) ParameterizedQuery(ctx context.Context, credential ResolvedCredential, toolID string, query ParameterizedQuery) (ParameterizedQueryResult, error) {
	metadata := parameterizedToolMetadata(toolID)
	values, paths, err := batchDLiveRequest(toolID, query)
	if err != nil {
		return ParameterizedQueryResult{}, err
	}
	payload, err := client.firstSuccessfulBatchDPath(ctx, credential, paths, values)
	if err != nil {
		return ParameterizedQueryResult{}, err
	}
	items := livePayloadItems(payload)
	result := ParameterizedQueryResult{
		ToolID:     toolID,
		Title:      metadata.DisplayName,
		EntityType: metadata.EntityType,
		Query:      query,
		Items:      make([]QueryResultItem, 0, len(items)),
	}
	for _, item := range items {
		if metadata.EntityType == QueryEntityVulnerability {
			result.Items = append(result.Items, liveVulnerabilityItem(item))
			continue
		}
		result.Items = append(result.Items, liveAssetItem(item))
	}
	return result, nil
}

func batchDLiveRequest(toolID string, query ParameterizedQuery) (url.Values, []string, error) {
	values := url.Values{}
	values.Set("page", fmt.Sprint(positiveInt(query.Page, 1)))
	values.Set("per_page", fmt.Sprint(boundedPageSize(query.PageSize, 20)))
	switch toolID {
	case CapabilityListAssetsByOwner:
		addLivePersonSearch(values, "oper_info", query)
		return values, []string{currentAssetListPath, legacyAssetListPath}, nil
	case CapabilityListAssetsByDepartment:
		addLiveSearchCondition(values, "business_department.name.keyword", []string{query.DepartmentName}, "==")
		return values, []string{currentAssetListPath, legacyAssetListPath}, nil
	case CapabilityListAssetsByIP:
		values.Set("keyword", strings.TrimSpace(query.IP))
		addLiveSearchCondition(values, "ip", []string{query.IP}, "==")
		return values, []string{currentAssetListPath, legacyAssetListPath}, nil
	case CapabilityListVulnerabilitiesByOwner:
		values.Set("data_range", defaultThreatCenterDataRange)
		addLivePersonSearch(values, "person_info", query)
		addLiveThreatFilters(values, query)
		return values, []string{currentThreatCenterPath, legacyThreatCenterPath}, nil
	case CapabilityListVulnerabilitiesByDepartment:
		values.Set("data_range", defaultThreatCenterDataRange)
		addLiveSearchCondition(values, "person_department.name.keyword", []string{query.DepartmentName}, "==")
		addLiveThreatFilters(values, query)
		return values, []string{currentThreatCenterPath, legacyThreatCenterPath}, nil
	case CapabilityListVulnerabilitiesByIP:
		values.Set("data_range", defaultThreatCenterDataRange)
		// 旧 Fobrain 的 ip query 参数会触发 IP 画像分支并强制 data_range=1；Batch D 需要完整非回收站范围。
		addLiveSearchCondition(values, "ip", []string{query.IP}, "==")
		addLiveThreatFilters(values, query)
		return values, []string{currentThreatCenterPath, legacyThreatCenterPath}, nil
	default:
		return nil, nil, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 参数化查询能力未注册")
	}
}

func addLivePersonSearch(values url.Values, fieldPrefix string, query ParameterizedQuery) {
	if staffID := strings.TrimSpace(query.PersonStaffID); staffID != "" {
		addLiveSearchCondition(values, fieldPrefix+".id", []string{staffID}, "in")
		return
	}
	addLiveSearchCondition(values, fieldPrefix+".name", []string{query.PersonName}, "==")
}

func addLiveThreatFilters(values url.Values, query ParameterizedQuery) {
	if severityCodes := liveThreatLevelCodes(query.Severity); len(severityCodes) > 0 {
		addLiveSearchCondition(values, "level", severityCodes, "==")
	}
	if statusCodes := liveThreatStatusCodes(query.Status); len(statusCodes) > 0 {
		addLiveSearchCondition(values, "status", statusCodes, "==")
	}
}

func addLiveSearchCondition[T any](values url.Values, field string, entries []T, operation string) {
	field = strings.TrimSpace(field)
	if field == "" || len(entries) == 0 {
		return
	}
	cleaned := make([]T, 0, len(entries))
	for _, entry := range entries {
		if strings.TrimSpace(fmt.Sprint(entry)) != "" {
			cleaned = append(cleaned, entry)
		}
	}
	if len(cleaned) == 0 {
		return
	}
	payload := map[string]any{field: cleaned, "operation_type_string": operation}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return
	}
	values.Add("search_condition", string(encoded))
}

func (client *HTTPClient) firstSuccessfulBatchDPath(ctx context.Context, credential ResolvedCredential, paths []string, query url.Values) (any, error) {
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
	for _, path := range paths {
		payload, ok, err := client.batchDGetAtPath(ctx, credential, authParam, path, query)
		if err != nil {
			return nil, err
		}
		if ok {
			return payload, nil
		}
	}
	return nil, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 参数化查询接口不存在")
}

func (client *HTTPClient) batchDGetAtPath(ctx context.Context, credential ResolvedCredential, authParam string, path string, query url.Values) (any, bool, error) {
	requestCtx := ctx
	cancel := func() {}
	if _, ok := ctx.Deadline(); !ok && client.timeout > 0 {
		requestCtx, cancel = context.WithTimeout(ctx, client.timeout)
	}
	defer cancel()

	req, err := http.NewRequestWithContext(requestCtx, http.MethodGet, client.endpointWithQuery(path, query), nil)
	if err != nil {
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 参数化查询请求创建失败")
	}
	req.Header.Set(authParam, credential.APIToken)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		if errors.Is(requestCtx.Err(), context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, false, NewSafeError(capabilities.PolicyReasonConnectorTimeout, "Fobrain 参数化查询请求超时")
		}
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorTransportUnavailable, "Fobrain 参数化查询网络不可用")
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorAuthFailure, "Fobrain 参数化查询认证失败")
	}
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusMethodNotAllowed {
		return nil, false, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 参数化查询失败")
	}
	var payload any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 参数化查询响应不是合法 JSON")
	}
	if !livePayloadCodeOK(payload) {
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 参数化查询返回业务错误")
	}
	return payload, true, nil
}

func livePayloadCodeOK(payload any) bool {
	item, ok := payload.(map[string]any)
	if !ok {
		return true
	}
	code, ok := item["code"]
	return !ok || intFromAny(code) == 0
}

func livePayloadItems(payload any) []map[string]any {
	switch typed := payload.(type) {
	case []any:
		return liveMapsFromSlice(typed)
	case map[string]any:
		for _, key := range []string{"items", "list", "records"} {
			if value, ok := typed[key].([]any); ok {
				return liveMapsFromSlice(value)
			}
		}
		for _, key := range []string{"data", "result"} {
			if value, ok := typed[key]; ok {
				items := livePayloadItems(value)
				if len(items) > 0 {
					return items
				}
			}
		}
	}
	return nil
}

func liveMapsFromSlice(values []any) []map[string]any {
	items := make([]map[string]any, 0, len(values))
	for _, value := range values {
		if item, ok := value.(map[string]any); ok {
			items = append(items, item)
		}
	}
	return items
}

func liveAssetItem(item map[string]any) QueryResultItem {
	id := safeRefPart(firstLiveString(item, "id", "asset_id", "assetId"), firstLiveString(item, "ip"))
	ip := liveFirstText(firstPresentValue(item, "ip", "ips"))
	name := firstNonEmptyLiveText(
		liveFirstText(firstPresentValue(item, "hostname", "name", "asset_name")),
		ip,
	)
	refPart := safeRefPart(id, name)
	return QueryResultItem{
		EntityRef:   "asset:fobrain:" + refPart,
		DisplayName: name,
		OwnerName:   firstNonEmptyLiveText(liveNestedText(item["oper_info"], "name"), liveNestedText(item["oper"], "name")),
		Status:      firstLiveString(item, "status"),
		Affected:    intFromAny(firstPresentValue(item, "poc_num", "vul_count", "vulnerability_count")),
		Summary:     liveJoinSummary(ip, firstLiveString(item, "network_type"), firstLiveString(item, "ip_type")),
	}
}

func liveVulnerabilityItem(item map[string]any) QueryResultItem {
	id := safeRefPart(firstLiveString(item, "id", "vuln_id", "vulnerability_id"), firstLiveString(item, "name", "title"))
	title := firstNonEmptyLiveText(firstLiveString(item, "name", "title", "threat_name", "vulnerability_name", "vul_name"), firstLiveString(item, "cve"))
	return QueryResultItem{
		EntityRef:   "vuln:fobrain:" + id,
		DisplayName: title,
		OwnerName:   firstNonEmptyLiveText(liveNestedText(item["person_info"], "name"), firstLiveString(item, "person_name", "owner")),
		Status:      liveThreatStatusLabel(firstPresentValue(item, "statusCode", "status_code", "status", "threat_status", "state")),
		Severity:    liveThreatLevelLabel(firstPresentValue(item, "level", "threat_level", "severity")),
		Affected:    intFromAny(firstPresentValue(item, "risk_num", "asset_count", "affected_assets", "ip_count", "vul_count", "relevance_num", "poc_num")),
		Summary:     firstLiveString(item, "describe", "summary", "description"),
	}
}

func firstLiveString(item map[string]any, keys ...string) string {
	return liveFirstText(firstPresentValue(item, keys...))
}

func liveFirstText(value any) string {
	switch typed := value.(type) {
	case string:
		return safeLiveText(typed)
	case json.Number:
		return safeLiveText(typed.String())
	case float64:
		if typed == float64(int64(typed)) {
			return safeLiveText(fmt.Sprintf("%.0f", typed))
		}
		return safeLiveText(fmt.Sprint(typed))
	case int:
		return safeLiveText(fmt.Sprint(typed))
	case int64:
		return safeLiveText(fmt.Sprint(typed))
	case []any:
		for _, entry := range typed {
			if text := liveFirstText(entry); text != "" {
				return text
			}
		}
	case map[string]any:
		return firstLiveString(typed, "name", "label", "title", "hostname", "ip")
	}
	return ""
}

func liveNestedText(value any, keys ...string) string {
	switch typed := value.(type) {
	case map[string]any:
		return firstLiveString(typed, keys...)
	case []any:
		for _, entry := range typed {
			if text := liveNestedText(entry, keys...); text != "" {
				return text
			}
		}
	}
	return ""
}

func safeLiveText(value string) string {
	text := strings.TrimSpace(value)
	if text == "" || facts.ContainsUnsafeMaterial(text) {
		return ""
	}
	runes := []rune(text)
	if len(runes) > 160 {
		return string(runes[:160])
	}
	return text
}

func firstNonEmptyLiveText(values ...string) string {
	for _, value := range values {
		if text := safeLiveText(value); text != "" {
			return text
		}
	}
	return ""
}

func liveJoinSummary(values ...string) string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		if text := safeLiveText(value); text != "" {
			out = append(out, text)
		}
	}
	return strings.Join(out, " / ")
}

func safeRefPart(value string, fallback string) string {
	text := firstNonEmptyLiveText(value, fallback, "unknown")
	var builder strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' || r == '-' {
			builder.WriteRune(r)
			continue
		}
		builder.WriteRune('-')
	}
	out := strings.Trim(builder.String(), "-")
	if out == "" {
		return "unknown"
	}
	runes := []rune(out)
	if len(runes) > 80 {
		return string(runes[:80])
	}
	return out
}

func liveThreatLevelCodes(severity string) []int {
	switch strings.ToLower(strings.TrimSpace(severity)) {
	case "low", "低", "低危":
		return []int{1}
	case "medium", "中", "中危":
		return []int{2}
	case "high", "高", "高危":
		return []int{3}
	case "critical", "严重", "严重危":
		return []int{4}
	default:
		if number, err := strconv.Atoi(strings.TrimSpace(severity)); err == nil && number > 0 {
			return []int{number}
		}
		return nil
	}
}

func liveThreatStatusCodes(status string) []int {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "open", "in_progress", "待修复", "复测中":
		return []int{0, 1, 10, 11, 12, 13, 14, 15, 17}
	case "fixed", "已修复", "复测通过":
		return []int{30, 40, 41, 42}
	case "ignored", "误报", "忽略":
		return []int{40, 41}
	default:
		if number, err := strconv.Atoi(strings.TrimSpace(status)); err == nil && number > 0 {
			return []int{number}
		}
		return nil
	}
}

func liveThreatLevelLabel(value any) string {
	switch strings.TrimSpace(strings.ToLower(liveFirstText(value))) {
	case "1", "低", "低危":
		return "low"
	case "2", "中", "中危":
		return "medium"
	case "3", "高", "高危":
		return "high"
	case "4", "严重", "严重危":
		return "critical"
	case "5", "未知", "unknown":
		return "unknown"
	default:
		return liveFirstText(value)
	}
}

func liveThreatStatusLabel(value any) string {
	switch strings.TrimSpace(strings.ToLower(liveFirstText(value))) {
	case "0", "1", "10", "11", "12", "13", "14", "15", "16", "17", "新增", "复现", "待修复", "延时", "超时", "复测中", "复测未通过", "待复测", "催促":
		return "open"
	case "30", "42", "复测通过", "缓解修复":
		return "fixed"
	case "40", "41", "误报", "无法修复":
		return "ignored"
	default:
		return liveFirstText(value)
	}
}
