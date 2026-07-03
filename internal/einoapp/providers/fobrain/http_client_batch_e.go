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
	currentInternalAssetDetailPath = "/api/internal_asset/%s"
	legacyInternalAssetDetailPath  = "/api/v1/internal_asset/%s"
	currentExternalAssetDetailPath = "/api/external_ip_asset/%s"
	legacyExternalAssetDetailPath  = "/api/v1/external_ip_asset/%s"
	currentDeviceDetailPath        = "/api/device/%s"
	legacyDeviceDetailPath         = "/api/v1/device/%s"
	currentDomainDetailPath        = "/api/domain_asset/%s"
	legacyDomainDetailPath         = "/api/v1/domain_asset/%s"
	currentThreatDetailPath        = "/api/threat_center/%s"
	legacyThreatDetailPath         = "/api/v1/threat_center/%s"
	currentThreatCountPath         = "/api/threat_center/count"
	legacyThreatCountPath          = "/api/v1/threat_center/count"
	currentThreatRelevancePath     = "/api/threat_center/relevance/list"
	legacyThreatRelevancePath      = "/api/v1/threat_center/relevance/list"
)

// AssetDetail 调用 Fobrain 资产详情接口，并只返回安全字段摘要。
func (client *HTTPClient) AssetDetail(ctx context.Context, credential ResolvedCredential, query AssetDetailQuery) (DetailRiskResult, error) {
	query.NetworkType = canonicalAssetNetworkType(query.NetworkType)
	payload, err := client.firstSuccessfulBatchEPath(ctx, credential, http.MethodGet, assetDetailPaths(query.NetworkType, query.AssetID), nil, nil)
	if err != nil {
		return DetailRiskResult{}, err
	}
	item, ok := unwrapPayloadMap(payload)
	if !ok {
		return DetailRiskResult{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 资产详情响应结构无效")
	}
	return DetailRiskResult{
		ToolID: CapabilityGetAssetDetail,
		Target: query.AssetID,
		Items:  []QueryResultItem{liveAssetItem(item)},
	}, nil
}

// VulnerabilityDetail 调用 Fobrain 漏洞详情接口，并复用漏洞安全行摘要 mapper。
func (client *HTTPClient) VulnerabilityDetail(ctx context.Context, credential ResolvedCredential, query VulnerabilityDetailQuery) (DetailRiskResult, error) {
	escaped := url.PathEscape(strings.TrimSpace(query.VulnerabilityID))
	payload, err := client.firstSuccessfulBatchEPath(ctx, credential, http.MethodGet, []string{
		fmt.Sprintf(currentThreatDetailPath, escaped),
		fmt.Sprintf(legacyThreatDetailPath, escaped),
	}, nil, nil)
	if err != nil {
		return DetailRiskResult{}, err
	}
	item, ok := unwrapPayloadMap(payload)
	if !ok {
		return DetailRiskResult{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 漏洞详情响应结构无效")
	}
	return DetailRiskResult{
		ToolID: CapabilityGetVulnerabilityDetail,
		Target: query.VulnerabilityID,
		Items:  []QueryResultItem{liveVulnerabilityItem(item)},
	}, nil
}

// BusinessRiskSummary 调用 threat_center/count 聚合接口，响应只保留安全统计 bucket。
func (client *HTTPClient) BusinessRiskSummary(ctx context.Context, credential ResolvedCredential, query BusinessRiskQuery) (DetailRiskResult, error) {
	body := []map[string]any{{
		"count_name":        "business_risk",
		"aggregation_field": "business.name.keyword",
		"data_range":        4,
		"search_condition":  []string{encodedLiveSearchCondition("business_name", query.BusinessName)},
	}}
	payload, err := client.firstSuccessfulBatchEPath(ctx, credential, http.MethodPost, []string{
		currentThreatCountPath,
		legacyThreatCountPath,
	}, nil, body)
	if err != nil {
		return DetailRiskResult{}, err
	}
	return DetailRiskResult{
		ToolID:  CapabilityBusinessRiskSummary,
		Target:  query.BusinessName,
		Metrics: liveRiskMetrics(payload),
	}, nil
}

func encodedLiveSearchCondition(field string, value string) string {
	// Fobrain count 接口沿用旧项目约定：search_condition 是 JSON 字符串数组，而不是对象数组。
	encoded, _ := json.Marshal(map[string]any{
		field:                   []string{strings.TrimSpace(value)},
		"operation_type_string": "==",
	})
	return string(encoded)
}

// ThreatRelevanceList 调用漏洞关联列表接口，vulnerability_name 同步映射到 keyword 和 vul_name。
func (client *HTTPClient) ThreatRelevanceList(ctx context.Context, credential ResolvedCredential, query ThreatRelevanceQuery) (DetailRiskResult, error) {
	values := url.Values{}
	values.Set("page", fmt.Sprint(positiveInt(query.Page, 1)))
	values.Set("per_page", fmt.Sprint(boundedPageSize(query.PageSize, 20)))
	values.Set("keyword", strings.TrimSpace(query.VulnerabilityName))
	values.Set("vul_name", strings.TrimSpace(query.VulnerabilityName))
	if strings.TrimSpace(query.IP) != "" {
		values.Set("ip", strings.TrimSpace(query.IP))
	}
	if strings.TrimSpace(query.BusinessName) != "" {
		values.Set("business_name", strings.TrimSpace(query.BusinessName))
	}
	payload, err := client.firstSuccessfulBatchEPath(ctx, credential, http.MethodGet, []string{
		currentThreatRelevancePath,
		legacyThreatRelevancePath,
	}, values, nil)
	if err != nil {
		return DetailRiskResult{}, err
	}
	items := livePayloadItems(payload)
	result := DetailRiskResult{
		ToolID: CapabilityThreatRelevanceList,
		Target: query.VulnerabilityName,
		Items:  make([]QueryResultItem, 0, len(items)),
	}
	for _, item := range items {
		result.Items = append(result.Items, liveVulnerabilityItem(item))
	}
	return result, nil
}

func assetDetailPaths(networkType string, assetID string) []string {
	escaped := url.PathEscape(strings.TrimSpace(assetID))
	switch canonicalAssetNetworkType(networkType) {
	case assetNetworkExternal:
		return []string{fmt.Sprintf(currentExternalAssetDetailPath, escaped), fmt.Sprintf(legacyExternalAssetDetailPath, escaped)}
	case assetNetworkDevice:
		return []string{fmt.Sprintf(currentDeviceDetailPath, escaped), fmt.Sprintf(legacyDeviceDetailPath, escaped)}
	case assetNetworkDomain:
		return []string{fmt.Sprintf(currentDomainDetailPath, escaped), fmt.Sprintf(legacyDomainDetailPath, escaped)}
	default:
		return []string{fmt.Sprintf(currentInternalAssetDetailPath, escaped), fmt.Sprintf(legacyInternalAssetDetailPath, escaped)}
	}
}

func (client *HTTPClient) firstSuccessfulBatchEPath(ctx context.Context, credential ResolvedCredential, method string, paths []string, query url.Values, body any) (any, error) {
	if client == nil || client.httpClient == nil || client.baseURL == nil {
		return nil, NewSafeError(capabilities.PolicyReasonConnectorTransportUnavailable, "Fobrain 详情风险 client 未配置")
	}
	if strings.TrimSpace(credential.APIToken) == "" {
		return nil, NewSafeError(capabilities.PolicyReasonCredentialMissing, "Fobrain 凭据未配置")
	}
	authParam, err := explicitAuthParam(credential)
	if err != nil {
		return nil, err
	}
	for _, path := range paths {
		payload, ok, err := client.batchERequestAtPath(ctx, credential, authParam, method, path, query, body)
		if err != nil {
			return nil, err
		}
		if ok {
			return payload, nil
		}
	}
	return nil, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 详情风险接口不存在")
}

func (client *HTTPClient) batchERequestAtPath(ctx context.Context, credential ResolvedCredential, authParam string, method string, path string, query url.Values, body any) (any, bool, error) {
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
			return nil, false, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 详情风险请求编码失败")
		}
		bodyReader = bytes.NewReader(encoded)
	}
	endpoint := client.endpoint(path)
	if len(query) > 0 {
		endpoint = client.endpointWithQuery(path, query)
	}
	req, err := http.NewRequestWithContext(requestCtx, method, endpoint, bodyReader)
	if err != nil {
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 详情风险请求创建失败")
	}
	req.Header.Set(authParam, credential.APIToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		if errors.Is(requestCtx.Err(), context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, false, NewSafeError(capabilities.PolicyReasonConnectorTimeout, "Fobrain 详情风险请求超时")
		}
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorTransportUnavailable, "Fobrain 详情风险网络不可用")
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorAuthFailure, "Fobrain 详情风险认证失败")
	}
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusMethodNotAllowed {
		return nil, false, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 详情风险请求失败")
	}
	var payload any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 详情风险响应不是合法 JSON")
	}
	if !livePayloadCodeOK(payload) {
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 详情风险返回业务错误")
	}
	return payload, true, nil
}

func liveRiskMetrics(payload any) []RiskMetric {
	items := liveRiskMetricItems(batchEPayloadData(payload))
	metrics := make([]RiskMetric, 0, len(items))
	for _, item := range items {
		label := firstNonEmptyLiveText(firstLiveString(item, "name", "label", "key", "business_name"), "未命名")
		count := intFromAny(firstPresentValue(item, "count", "doc_count", "value", "total"))
		metrics = append(metrics, RiskMetric{Label: label, Count: count})
	}
	return metrics
}

func liveRiskMetricItems(value any) []map[string]any {
	switch typed := value.(type) {
	case []any:
		return liveMapsFromSlice(typed)
	case map[string]any:
		for _, key := range []string{"items", "buckets", "list", "records"} {
			if entries, ok := typed[key].([]any); ok {
				return liveMapsFromSlice(entries)
			}
		}
	}
	return livePayloadItems(value)
}

func batchEPayloadData(payload any) any {
	item, ok := payload.(map[string]any)
	if !ok {
		return payload
	}
	if data, ok := item["data"]; ok {
		return data
	}
	return payload
}
