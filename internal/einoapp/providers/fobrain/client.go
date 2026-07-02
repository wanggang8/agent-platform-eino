package fobrain

import "context"

// FobrainClient 是 provider 边界内的业务读取接口。
// 实现可以是 mock 或 HTTP client，但 raw provider payload 不能离开该边界。
type FobrainClient interface {
	CurrentUserContext(context.Context, ResolvedCredential) (CurrentUserContextResult, error)
	MyPermissions(context.Context, ResolvedCredential) (MyPermissionsResult, error)
}

// ParameterizedQueryClient 是 Batch D 参数化只读查询的可选 client 能力。
// HTTP live 接入未完成前，provider 只在 client 明确实现该接口时启用 Batch D 调用。
type ParameterizedQueryClient interface {
	ParameterizedQuery(context.Context, ResolvedCredential, string, ParameterizedQuery) (ParameterizedQueryResult, error)
}

// DetailRiskClient 是 Batch E 详情和风险关联只读能力的可选 client 能力。
// 该接口只返回安全业务结果，不能把 raw provider payload 传出 provider 边界。
type DetailRiskClient interface {
	AssetDetail(context.Context, ResolvedCredential, AssetDetailQuery) (DetailRiskResult, error)
	VulnerabilityDetail(context.Context, ResolvedCredential, VulnerabilityDetailQuery) (DetailRiskResult, error)
	BusinessRiskSummary(context.Context, ResolvedCredential, BusinessRiskQuery) (DetailRiskResult, error)
	ThreatRelevanceList(context.Context, ResolvedCredential, ThreatRelevanceQuery) (DetailRiskResult, error)
}

// CurrentUserContextResult 是当前用户 PoC 的安全业务结果。
type CurrentUserContextResult struct {
	DisplayName string
	Department  string
	Role        string
}

// MyPermissionsResult 是 Batch A 权限范围读取的安全业务结果。
type MyPermissionsResult struct {
	DisplayName         string
	PermissionNames     []string
	DataPermissionNames []string
}

// MockClient 是 Phase 5 smoke 使用的本地 Fobrain client，不触达真实网络。
type MockClient struct {
	Result      CurrentUserContextResult
	Permissions MyPermissionsResult
}

// CurrentUserContext 返回安全 fixture 结果，用于证明 provider/Facts 链路。
func (client MockClient) CurrentUserContext(_ context.Context, _ ResolvedCredential) (CurrentUserContextResult, error) {
	if client.Result.DisplayName == "" {
		return CurrentUserContextResult{
			DisplayName: "Fobrain 本地用户",
			Department:  "安全运营",
			Role:        "只读验证",
		}, nil
	}
	return client.Result, nil
}

// MyPermissions 返回安全权限 fixture，用于 Batch A mock 验收。
func (client MockClient) MyPermissions(_ context.Context, _ ResolvedCredential) (MyPermissionsResult, error) {
	if len(client.Permissions.PermissionNames) == 0 && len(client.Permissions.DataPermissionNames) == 0 {
		return MyPermissionsResult{
			DisplayName:         "Fobrain 本地用户",
			PermissionNames:     []string{"只读验证"},
			DataPermissionNames: []string{"本地 fixture 数据范围"},
		}, nil
	}
	return client.Permissions, nil
}

// ParameterizedQuery 返回 Batch D mock 只读结果，用于固定 catalog、参数和 StructuredResult 边界。
func (client MockClient) ParameterizedQuery(_ context.Context, _ ResolvedCredential, toolID string, query ParameterizedQuery) (ParameterizedQueryResult, error) {
	metadata := parameterizedToolMetadata(toolID)
	return ParameterizedQueryResult{
		ToolID:     toolID,
		Title:      metadata.DisplayName,
		EntityType: metadata.EntityType,
		Query:      query,
		Items: []QueryResultItem{
			{
				EntityRef:   metadata.EntityType + ":fobrain:mock-1",
				DisplayName: metadata.DisplayName + "结果",
				OwnerName:   query.PersonName,
				Status:      "resolved",
				Summary:     parameterizedQueryTarget(query),
				Severity:    query.Severity,
				Affected:    1,
			},
		},
	}, nil
}

// AssetDetail 返回 Batch E 资产详情 mock 结果，用于固定详情 StructuredResult 边界。
func (client MockClient) AssetDetail(_ context.Context, _ ResolvedCredential, query AssetDetailQuery) (DetailRiskResult, error) {
	return DetailRiskResult{
		ToolID: CapabilityGetAssetDetail,
		Target: query.AssetID,
		Items: []QueryResultItem{{
			EntityRef:   "asset:fobrain:" + safeRefPart(query.AssetID, "mock-asset"),
			DisplayName: "Fobrain mock 资产",
			Status:      "resolved",
			Summary:     query.NetworkType,
			Affected:    1,
		}},
	}, nil
}

// VulnerabilityDetail 返回 Batch E 漏洞详情 mock 结果。
func (client MockClient) VulnerabilityDetail(_ context.Context, _ ResolvedCredential, query VulnerabilityDetailQuery) (DetailRiskResult, error) {
	return DetailRiskResult{
		ToolID: CapabilityGetVulnerabilityDetail,
		Target: query.VulnerabilityID,
		Items: []QueryResultItem{{
			EntityRef:   "vuln:fobrain:" + safeRefPart(query.VulnerabilityID, "mock-vuln"),
			DisplayName: "Fobrain mock 漏洞",
			Severity:    "high",
			Status:      "open",
			Affected:    1,
		}},
	}, nil
}

// BusinessRiskSummary 返回 Batch E 业务风险聚合 mock 结果。
func (client MockClient) BusinessRiskSummary(_ context.Context, _ ResolvedCredential, query BusinessRiskQuery) (DetailRiskResult, error) {
	return DetailRiskResult{
		ToolID: CapabilityBusinessRiskSummary,
		Target: query.BusinessName,
		Metrics: []RiskMetric{{
			Label: query.BusinessName,
			Count: 1,
		}},
	}, nil
}

// ThreatRelevanceList 返回 Batch E 威胁关联 mock 结果。
func (client MockClient) ThreatRelevanceList(_ context.Context, _ ResolvedCredential, query ThreatRelevanceQuery) (DetailRiskResult, error) {
	return DetailRiskResult{
		ToolID: CapabilityThreatRelevanceList,
		Target: query.VulnerabilityName,
		Items: []QueryResultItem{{
			EntityRef:   "vuln:fobrain:mock-relevance",
			DisplayName: query.VulnerabilityName,
			Severity:    "high",
			Affected:    1,
			Summary:     firstNonEmptyLiveText(query.IP, query.BusinessName),
		}},
	}, nil
}
