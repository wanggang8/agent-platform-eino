package fobrain

import (
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
	businessImportantLevelVeryHigh = 1
	businessImportantLevelHigh     = 2
)

type batchBLiveHTTPRequest struct {
	paths []string
	query url.Values
}

// MyScopeQuery 调用 Batch B “我的范围”真实只读接口，并把当前用户/部门上下文限制在 provider 边界内。
func (client *HTTPClient) MyScopeQuery(ctx context.Context, credential ResolvedCredential, toolID string, query MyScopeQuery) (MyScopeResult, error) {
	metadata, ok := batchBMyScopeMetadata(toolID)
	if !ok {
		return MyScopeResult{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 我的范围能力未注册")
	}
	userContext, err := client.CurrentUserContext(ctx, credential)
	if err != nil {
		return MyScopeResult{}, err
	}
	request, err := batchBLiveRequest(toolID, query, userContext)
	if err != nil {
		return MyScopeResult{}, err
	}
	payload, err := client.firstSuccessfulBatchBPath(ctx, credential, request)
	if err != nil {
		return MyScopeResult{}, err
	}
	items := livePayloadItems(payload)
	result := MyScopeResult{
		ToolID:     toolID,
		Title:      metadata.DisplayName,
		EntityType: metadata.EntityType,
		Query:      query,
		Items:      make([]QueryResultItem, 0, len(items)),
	}
	for _, item := range items {
		result.Items = append(result.Items, liveMyScopeItem(metadata.EntityType, item))
	}
	return result, nil
}

func batchBLiveRequest(toolID string, query MyScopeQuery, userContext CurrentUserContextResult) (batchBLiveHTTPRequest, error) {
	values := url.Values{}
	values.Set("page", fmt.Sprint(positiveInt(query.Page, 1)))
	values.Set("per_page", fmt.Sprint(boundedPageSize(query.PageSize, 20)))
	switch toolID {
	case CapabilityMyAssets:
		if err := requireBatchBUserScope(userContext.DisplayName, "Fobrain 当前用户缺少姓名，无法查询我的资产"); err != nil {
			return batchBLiveHTTPRequest{}, err
		}
		addBatchBKeyword(values, query.Keyword)
		addLiveSearchCondition(values, "oper_info.name", []string{userContext.DisplayName}, "==")
		return batchBLiveHTTPRequest{paths: []string{currentAssetListPath, legacyAssetListPath}, query: values}, nil
	case CapabilityMyDepartmentAssets:
		if err := requireBatchBUserScope(userContext.Department, "Fobrain 当前用户缺少部门，无法查询本部门资产"); err != nil {
			return batchBLiveHTTPRequest{}, err
		}
		addBatchBKeyword(values, query.Keyword)
		addLiveSearchCondition(values, "business_department.name.keyword", []string{userContext.Department}, "==")
		return batchBLiveHTTPRequest{paths: []string{currentAssetListPath, legacyAssetListPath}, query: values}, nil
	case CapabilityMyVulnerabilities:
		if err := requireBatchBUserScope(userContext.DisplayName, "Fobrain 当前用户缺少姓名，无法查询我的漏洞"); err != nil {
			return batchBLiveHTTPRequest{}, err
		}
		values.Set("data_range", defaultThreatCenterDataRange)
		addBatchBKeyword(values, query.Keyword)
		addLiveSearchCondition(values, "person_info.name", []string{userContext.DisplayName}, "==")
		return batchBLiveHTTPRequest{paths: []string{currentThreatCenterPath, legacyThreatCenterPath}, query: values}, nil
	case CapabilityMyDepartmentVulnerabilities:
		if err := requireBatchBUserScope(userContext.Department, "Fobrain 当前用户缺少部门，无法查询本部门漏洞"); err != nil {
			return batchBLiveHTTPRequest{}, err
		}
		values.Set("data_range", defaultThreatCenterDataRange)
		addBatchBKeyword(values, query.Keyword)
		addLiveSearchCondition(values, "person_department.name.keyword", []string{userContext.Department}, "==")
		return batchBLiveHTTPRequest{paths: []string{currentThreatCenterPath, legacyThreatCenterPath}, query: values}, nil
	case CapabilityMyBusinessSystems, CapabilityMyImportantBusinessSystems:
		if err := requireBatchBUserScope(userContext.DisplayName, "Fobrain 当前用户缺少姓名，无法查询我的业务系统"); err != nil {
			return batchBLiveHTTPRequest{}, err
		}
		// 业务系统 owner scope 使用旧源码确认的 person_base.name 条件；用户 keyword 只能额外收窄业务名。
		addLiveSearchCondition(values, "person_base.name", []string{userContext.DisplayName}, "==")
		addBatchBBusinessNameFilter(values, query.Keyword)
		if toolID == CapabilityMyImportantBusinessSystems {
			addLiveSearchCondition(values, "assets_attribute.important_types", []int{businessImportantLevelVeryHigh, businessImportantLevelHigh}, "in")
		}
		return batchBLiveHTTPRequest{paths: []string{currentBusinessListPath, legacyBusinessListPath}, query: values}, nil
	default:
		return batchBLiveHTTPRequest{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 我的范围能力未注册")
	}
}

func requireBatchBUserScope(value string, message string) error {
	if strings.TrimSpace(value) == "" {
		return NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, message)
	}
	return nil
}

func addBatchBKeyword(values url.Values, keyword string) {
	if text := firstNonEmptyLiveText(keyword); text != "" {
		values.Set("keyword", text)
	}
}

func addBatchBBusinessNameFilter(values url.Values, keyword string) {
	if text := firstNonEmptyLiveText(keyword); text != "" {
		// 用户 keyword 只能作为业务系统名称的额外收窄条件，不能覆盖当前用户 owner scope。
		addLiveSearchCondition(values, "business_name", []string{text}, "==")
	}
}

func liveMyScopeItem(entityType string, item map[string]any) QueryResultItem {
	switch entityType {
	case QueryEntityAsset:
		return liveAssetItem(item)
	case QueryEntityVulnerability:
		return liveVulnerabilityItem(item)
	case myScopeEntityBusinessSystem:
		return liveBusinessSystemItem(item)
	default:
		return QueryResultItem{EntityRef: "entity:fobrain:" + safeRefPart(firstLiveString(item, "id"), "unknown"), DisplayName: firstLiveString(item, "name", "title")}
	}
}

func (client *HTTPClient) firstSuccessfulBatchBPath(ctx context.Context, credential ResolvedCredential, request batchBLiveHTTPRequest) (any, error) {
	if client == nil || client.httpClient == nil || client.baseURL == nil {
		return nil, NewSafeError(capabilities.PolicyReasonConnectorTransportUnavailable, "Fobrain 我的范围 client 未配置")
	}
	if strings.TrimSpace(credential.APIToken) == "" {
		return nil, NewSafeError(capabilities.PolicyReasonCredentialMissing, "Fobrain 凭据未配置")
	}
	authParam, err := explicitAuthParam(credential)
	if err != nil {
		return nil, err
	}
	for _, path := range request.paths {
		payload, ok, err := client.batchBGetAtPath(ctx, credential, authParam, path, request.query)
		if err != nil {
			return nil, err
		}
		if ok {
			return payload, nil
		}
	}
	return nil, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 我的范围接口不存在")
}

func (client *HTTPClient) batchBGetAtPath(ctx context.Context, credential ResolvedCredential, authParam string, path string, query url.Values) (any, bool, error) {
	requestCtx := ctx
	cancel := func() {}
	if _, ok := ctx.Deadline(); !ok && client.timeout > 0 {
		requestCtx, cancel = context.WithTimeout(ctx, client.timeout)
	}
	defer cancel()

	req, err := http.NewRequestWithContext(requestCtx, http.MethodGet, client.endpointWithQuery(path, query), nil)
	if err != nil {
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 我的范围请求创建失败")
	}
	req.Header.Set(authParam, credential.APIToken)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		if errors.Is(requestCtx.Err(), context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, false, NewSafeError(capabilities.PolicyReasonConnectorTimeout, "Fobrain 我的范围请求超时")
		}
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorTransportUnavailable, "Fobrain 我的范围网络不可用")
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorAuthFailure, "Fobrain 我的范围认证失败")
	}
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusMethodNotAllowed {
		return nil, false, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 我的范围请求失败")
	}
	var payload any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 我的范围响应不是合法 JSON")
	}
	if !livePayloadCodeOK(payload) {
		return nil, false, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 我的范围返回业务错误")
	}
	return payload, true, nil
}
