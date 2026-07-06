package fobrain

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"agent-platform-eino/internal/einoapp/capabilities"
)

const (
	currentBusinessListPath                = "/api/business"
	legacyBusinessListPath                 = "/api/v1/business"
	currentExternalHighRiskAssetListPath   = "/api/external_ip_asset"
	legacyExternalHighRiskAssetListPath    = "/api/v1/external_ip_asset"
	currentTicketPendingPath               = "/api/ticket/pending"
	legacyTicketPendingPath                = "/api/v1/ticket/pending"
	currentThreatRelevanceIPStatsPath      = "/api/threat_center/relevance/ip_stats"
	legacyThreatRelevanceIPStatsPath       = "/api/v1/threat_center/relevance/ip_stats"
	currentThreatRelevanceVulStatsPath     = "/api/threat_center/relevance/vul_stats"
	legacyThreatRelevanceVulStatsPath      = "/api/v1/threat_center/relevance/vul_stats"
	batchCVulnerabilityStatusCountName     = "vulnerability_status"
	batchCVulnerabilityStatusAggregateName = "status"
)

// DirectRead 调用 Batch C 直接列表与统计接口，并把真实响应压缩为安全列表或指标。
func (client *HTTPClient) DirectRead(ctx context.Context, credential ResolvedCredential, toolID string, query DirectReadQuery) (DirectReadResult, error) {
	metadata := directReadToolMetadata(toolID)
	request, err := batchCLiveRequest(toolID, query)
	if err != nil {
		return DirectReadResult{}, err
	}
	payload, err := client.firstSuccessfulBatchCPath(ctx, credential, request)
	if err != nil {
		return DirectReadResult{}, err
	}
	result := DirectReadResult{
		ToolID:     toolID,
		Title:      metadata.DisplayName,
		EntityType: metadata.EntityType,
		Query:      query,
	}
	if metadata.Metrics {
		result.Metrics = liveDirectReadMetrics(payload, toolID)
		return result, nil
	}
	for _, item := range livePayloadItems(payload) {
		result.Items = append(result.Items, liveDirectReadItem(metadata.EntityType, item))
	}
	return result, nil
}

type batchCLiveHTTPRequest struct {
	method string
	paths  []string
	query  url.Values
	body   any
}

func batchCLiveRequest(toolID string, query DirectReadQuery) (batchCLiveHTTPRequest, error) {
	values := url.Values{}
	switch toolID {
	case CapabilityBusinessList:
		values.Set("page", fmt.Sprint(positiveInt(query.Page, 1)))
		values.Set("per_page", fmt.Sprint(boundedPageSize(query.PageSize, 20)))
		if keyword := firstNonEmptyLiveText(query.Keyword, query.BusinessName, query.Owner); keyword != "" {
			values.Set("keyword", keyword)
		}
		return batchCLiveHTTPRequest{method: http.MethodGet, paths: []string{currentBusinessListPath, legacyBusinessListPath}, query: values}, nil
	case CapabilityExternalHighRiskAssets:
		values.Set("page", fmt.Sprint(positiveInt(query.Page, 1)))
		values.Set("per_page", fmt.Sprint(boundedPageSize(query.PageSize, 20)))
		if keyword := firstNonEmptyLiveText(query.Keyword, query.Field); keyword != "" {
			values.Set("keyword", keyword)
		}
		if severityCodes := liveThreatLevelCodes(query.Severity); len(severityCodes) > 0 {
			addLiveSearchCondition(values, "risk", severityCodes, "==")
		}
		return batchCLiveHTTPRequest{method: http.MethodGet, paths: []string{currentExternalHighRiskAssetListPath, legacyExternalHighRiskAssetListPath}, query: values}, nil
	case CapabilityVulnerabilityStatusSummary:
		body := []map[string]any{{
			"count_name":        batchCVulnerabilityStatusCountName,
			"aggregation_field": batchCVulnerabilityStatusAggregateName,
			"data_range":        4,
			"search_condition":  batchCCountSearchConditions(query),
		}}
		return batchCLiveHTTPRequest{method: http.MethodPost, paths: []string{currentThreatCountPath, legacyThreatCountPath}, body: body}, nil
	case CapabilityPendingTickets:
		values.Set("page", fmt.Sprint(positiveInt(query.Page, 1)))
		values.Set("page_size", fmt.Sprint(boundedPageSize(query.PageSize, 20)))
		for _, status := range liveThreatStatusCodes(query.Status) {
			values.Add("status", fmt.Sprint(status))
		}
		if keyword := firstNonEmptyLiveText(query.Keyword); keyword != "" {
			values.Set("keyword", keyword)
		}
		// pending_tickets 是工单增量接口，旧项目和真实接口证据都不支持按人员服务端过滤。
		return batchCLiveHTTPRequest{method: http.MethodGet, paths: []string{legacyTicketPendingPath, currentTicketPendingPath}, query: values}, nil
	case CapabilityIPStats:
		if keyword := firstNonEmptyLiveText(query.Keyword, query.Field); keyword != "" {
			values.Set("keyword", keyword)
		}
		return batchCLiveHTTPRequest{method: http.MethodGet, paths: []string{currentThreatRelevanceIPStatsPath, legacyThreatRelevanceIPStatsPath}, query: values}, nil
	case CapabilityVulStats:
		if keyword := firstNonEmptyLiveText(query.Keyword, query.Field); keyword != "" {
			values.Set("keyword", keyword)
		}
		return batchCLiveHTTPRequest{method: http.MethodGet, paths: []string{currentThreatRelevanceVulStatsPath, legacyThreatRelevanceVulStatsPath}, query: values}, nil
	default:
		return batchCLiveHTTPRequest{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 直接读取能力未注册")
	}
}

func batchCCountSearchConditions(query DirectReadQuery) []string {
	conditions := make([]string, 0, 3)
	if keyword := firstNonEmptyLiveText(query.Keyword, query.Field); keyword != "" {
		conditions = append(conditions, encodedLiveSearchCondition("keyword", keyword))
	}
	for _, severity := range liveThreatLevelCodes(query.Severity) {
		conditions = append(conditions, encodedLiveSearchCondition("level", fmt.Sprint(severity)))
	}
	for _, status := range liveThreatStatusCodes(query.Status) {
		conditions = append(conditions, encodedLiveSearchCondition("status", fmt.Sprint(status)))
	}
	return conditions
}

func (client *HTTPClient) firstSuccessfulBatchCPath(ctx context.Context, credential ResolvedCredential, request batchCLiveHTTPRequest) (any, error) {
	if client == nil || client.httpClient == nil || client.baseURL == nil {
		return nil, NewSafeError(capabilities.PolicyReasonConnectorTransportUnavailable, "Fobrain 直接读取 client 未配置")
	}
	if strings.TrimSpace(credential.APIToken) == "" {
		return nil, NewSafeError(capabilities.PolicyReasonCredentialMissing, "Fobrain 凭据未配置")
	}
	authParam, err := explicitAuthParam(credential)
	if err != nil {
		return nil, err
	}
	for _, path := range request.paths {
		payload, ok, err := client.batchCRequestAtPath(ctx, credential, authParam, request.method, path, request.query, request.body)
		if err != nil {
			return nil, err
		}
		if ok {
			return payload, nil
		}
	}
	return nil, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 直接读取接口不存在")
}

func (client *HTTPClient) batchCRequestAtPath(ctx context.Context, credential ResolvedCredential, authParam string, method string, path string, query url.Values, body any) (any, bool, error) {
	requestCtx := ctx
	cancel := func() {}
	if _, ok := ctx.Deadline(); !ok && client.timeout > 0 {
		requestCtx, cancel = context.WithTimeout(ctx, client.timeout)
	}
	defer cancel()

	var bodyReader *bytes.Reader
	if body == nil {
		bodyReader = bytes.NewReader(nil)
	} else {
		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, false, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 直接读取请求编码失败")
		}
		bodyReader = bytes.NewReader(encoded)
	}
	endpoint := client.endpoint(path)
	if len(query) > 0 {
		endpoint = client.endpointWithQuery(path, query)
	}
	req, err := http.NewRequestWithContext(requestCtx, method, endpoint, bodyReader)
	if err != nil {
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 直接读取请求创建失败")
	}
	req.Header.Set(authParam, credential.APIToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		if errors.Is(requestCtx.Err(), context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, false, NewSafeError(capabilities.PolicyReasonConnectorTimeout, "Fobrain 直接读取请求超时")
		}
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorTransportUnavailable, "Fobrain 直接读取网络不可用")
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorAuthFailure, "Fobrain 直接读取认证失败")
	}
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusMethodNotAllowed {
		return nil, false, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 直接读取请求失败")
	}
	var payload any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 直接读取响应不是合法 JSON")
	}
	if !livePayloadCodeOK(payload) {
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 直接读取返回业务错误")
	}
	return payload, true, nil
}

func liveDirectReadItem(entityType string, item map[string]any) QueryResultItem {
	switch entityType {
	case QueryEntityAsset:
		return liveAssetItem(item)
	case QueryEntityVulnerability:
		return liveVulnerabilityItem(item)
	case myScopeEntityBusinessSystem:
		return liveBusinessSystemItem(item)
	case directReadEntityTicket:
		return liveTicketItem(item)
	default:
		return QueryResultItem{EntityRef: "entity:fobrain:" + safeRefPart(firstLiveString(item, "id"), "unknown"), DisplayName: firstLiveString(item, "name", "title")}
	}
}

func liveBusinessSystemItem(item map[string]any) QueryResultItem {
	id := safeRefPart(firstLiveString(item, "id", "business_id"), firstLiveString(item, "business_name", "name"))
	name := firstNonEmptyLiveText(firstLiveString(item, "business_name", "name", "system_name"), "未命名业务系统")
	return QueryResultItem{
		EntityRef:   "business_system:fobrain:" + id,
		DisplayName: name,
		OwnerName:   firstNonEmptyLiveText(liveNestedText(item["person_info"], "name"), liveNestedText(item["person"], "name"), firstLiveString(item, "owner", "principal")),
		Status:      firstLiveString(item, "status", "running_state", "reliability"),
		Affected:    intFromAny(firstPresentValue(item, "risk_num", "asset_count", "vul_count", "ip_count")),
		Summary:     liveJoinSummary(firstLiveString(item, "department_name", "department"), firstLiveString(item, "address")),
	}
}

func liveTicketItem(item map[string]any) QueryResultItem {
	id := safeRefPart(firstLiveString(item, "id", "ticket_id", "external_ticket_id", "task_id"), firstLiveString(item, "poc_id"))
	title := firstNonEmptyLiveText(firstLiveString(item, "title", "name", "vul_name", "poc_name", "threat_name"), "待处理工单")
	return QueryResultItem{
		EntityRef:   "ticket:fobrain:" + id,
		DisplayName: title,
		OwnerName:   firstNonEmptyLiveText(firstLiveString(item, "handler", "owner", "person_name"), liveNestedText(item["person_info"], "name")),
		Status:      firstLiveString(item, "status", "ticket_status", "task_status"),
		Severity:    liveThreatLevelLabel(firstPresentValue(item, "level", "severity", "threat_level")),
		Affected:    intFromAny(firstPresentValue(item, "retry_count", "risk_num", "vul_count")),
		Summary:     liveJoinSummary(firstLiveString(item, "external_ticket_id"), firstLiveString(item, "last_error", "remark")),
	}
}

func liveDirectReadMetrics(payload any, toolID string) []RiskMetric {
	if toolID == CapabilityVulnerabilityStatusSummary {
		if metrics := liveThreatCountAggregationMetrics(batchEPayloadData(payload)); len(metrics) > 0 {
			return metrics
		}
	}
	items := liveRiskMetricItems(batchEPayloadData(payload))
	metrics := make([]RiskMetric, 0, len(items))
	for _, item := range items {
		label := firstNonEmptyLiveText(
			firstLiveString(item, "label", "name", "key", "status_name", "status", "ip", "vul_name", "business_name"),
			"未命名",
		)
		count := intFromAny(firstPresentValue(item, "count", "doc_count", "value", "total", "vul_count", "ip_count"))
		metrics = append(metrics, RiskMetric{Label: label, Count: count})
	}
	if len(metrics) == 0 && batchCDirectReadCapabilityByID(toolID) {
		if total := intFromAny(firstPresentValue(payloadMap(batchEPayloadData(payload)), "total", "count")); total > 0 {
			metrics = append(metrics, RiskMetric{Label: directReadToolMetadata(toolID).DisplayName, Count: total})
		}
	}
	return metrics
}

func liveThreatCountAggregationMetrics(value any) []RiskMetric {
	items := liveRiskMetricItems(value)
	metrics := make([]RiskMetric, 0, len(items))
	for _, item := range items {
		buckets, ok := item["aggregation"].([]any)
		if !ok || len(buckets) == 0 {
			continue
		}
		for _, bucket := range liveMapsFromSlice(buckets) {
			label := firstNonEmptyLiveText(firstLiveString(bucket, "key", "label", "name"), "未命名")
			count := intFromAny(firstPresentValue(bucket, "count", "doc_count", "value", "total"))
			metrics = append(metrics, RiskMetric{Label: label, Count: count})
		}
	}
	return metrics
}

func payloadMap(value any) map[string]any {
	if item, ok := value.(map[string]any); ok {
		return item
	}
	return nil
}
