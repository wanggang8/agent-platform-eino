package fobrain

import (
	"fmt"
	"strings"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
)

const (
	assetNetworkInternal = "internal"
	assetNetworkExternal = "external"
	assetNetworkDevice   = "device"
	assetNetworkDomain   = "domain"
)

// AssetDetailQuery 是资产详情工具的安全入参，network_type 已在 provider 边界归一。
type AssetDetailQuery struct {
	AssetID     string
	NetworkType string
}

// VulnerabilityDetailQuery 是漏洞详情工具的安全入参。
type VulnerabilityDetailQuery struct {
	VulnerabilityID string
}

// BusinessRiskQuery 是业务风险聚合工具的安全入参。
type BusinessRiskQuery struct {
	BusinessName string
}

// ThreatRelevanceQuery 是威胁关联列表工具的安全入参。
type ThreatRelevanceQuery struct {
	VulnerabilityName string
	IP                string
	BusinessName      string
	Page              int
	PageSize          int
}

// RiskMetric 是业务风险聚合的安全统计项，不承载 raw request body 或 provider payload。
type RiskMetric struct {
	Label string
	Count int
}

// DetailRiskResult 是 Batch E provider 边界内的安全业务结果。
type DetailRiskResult struct {
	ToolID  string
	Target  string
	Items   []QueryResultItem
	Metrics []RiskMetric
}

type detailRiskMetadata struct {
	DisplayName string
	ResultSlug  string
}

func batchEDetailRiskCapabilityByID(capabilityID string) bool {
	_, ok := batchEDetailRiskMetadata(capabilityID)
	return ok
}

func detailRiskToolMetadata(toolID string) detailRiskMetadata {
	metadata, ok := batchEDetailRiskMetadata(toolID)
	if !ok {
		return detailRiskMetadata{DisplayName: "Fobrain 详情风险查询", ResultSlug: "detail-risk"}
	}
	return metadata
}

func batchEDetailRiskMetadata(toolID string) (detailRiskMetadata, bool) {
	switch toolID {
	case CapabilityGetAssetDetail:
		return detailRiskMetadata{DisplayName: "资产详情", ResultSlug: "get-asset-detail"}, true
	case CapabilityGetVulnerabilityDetail:
		return detailRiskMetadata{DisplayName: "漏洞详情", ResultSlug: "get-vulnerability-detail"}, true
	case CapabilityBusinessRiskSummary:
		return detailRiskMetadata{DisplayName: "业务风险摘要", ResultSlug: "business-risk-summary"}, true
	case CapabilityThreatRelevanceList:
		return detailRiskMetadata{DisplayName: "威胁关联列表", ResultSlug: "threat-relevance-list"}, true
	default:
		return detailRiskMetadata{}, false
	}
}

func assetDetailQueryFromArguments(arguments map[string]any) (AssetDetailQuery, error) {
	query := AssetDetailQuery{
		AssetID:     safeArgumentString(arguments["asset_id"]),
		NetworkType: canonicalAssetNetworkType(safeArgumentString(arguments["network_type"])),
	}
	if strings.TrimSpace(query.AssetID) == "" {
		return AssetDetailQuery{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 资产详情缺少必填参数")
	}
	if detailRiskUnsafeValues(query.AssetID, query.NetworkType) {
		return AssetDetailQuery{}, facts.ErrUnsafeFactMaterial
	}
	return query, nil
}

func vulnerabilityDetailQueryFromArguments(arguments map[string]any) (VulnerabilityDetailQuery, error) {
	query := VulnerabilityDetailQuery{VulnerabilityID: safeArgumentString(arguments["vulnerability_id"])}
	if strings.TrimSpace(query.VulnerabilityID) == "" {
		return VulnerabilityDetailQuery{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 漏洞详情缺少必填参数")
	}
	if detailRiskUnsafeValues(query.VulnerabilityID) {
		return VulnerabilityDetailQuery{}, facts.ErrUnsafeFactMaterial
	}
	return query, nil
}

func businessRiskQueryFromArguments(arguments map[string]any) (BusinessRiskQuery, error) {
	query := BusinessRiskQuery{BusinessName: safeArgumentString(arguments["business_name"])}
	if strings.TrimSpace(query.BusinessName) == "" {
		return BusinessRiskQuery{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 业务风险摘要缺少必填参数")
	}
	if detailRiskUnsafeValues(query.BusinessName) {
		return BusinessRiskQuery{}, facts.ErrUnsafeFactMaterial
	}
	return query, nil
}

func threatRelevanceQueryFromArguments(arguments map[string]any) (ThreatRelevanceQuery, error) {
	query := ThreatRelevanceQuery{
		VulnerabilityName: safeArgumentString(arguments["vulnerability_name"]),
		IP:                safeArgumentString(arguments["ip"]),
		BusinessName:      safeArgumentString(arguments["business_name"]),
		Page:              positiveInt(arguments["page"], 1),
		PageSize:          boundedPageSize(arguments["page_size"], 20),
	}
	if strings.TrimSpace(query.VulnerabilityName) == "" {
		return ThreatRelevanceQuery{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 威胁关联列表缺少必填参数")
	}
	if detailRiskUnsafeValues(query.VulnerabilityName, query.IP, query.BusinessName) {
		return ThreatRelevanceQuery{}, facts.ErrUnsafeFactMaterial
	}
	return query, nil
}

// canonicalAssetNetworkType 将旧项目数字、中文和英文别名统一为四个新项目契约值。
func canonicalAssetNetworkType(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "1", "internal", "intranet", "内网":
		return assetNetworkInternal
	case "2", "external", "internet", "external_ip", "外网", "互联网":
		return assetNetworkExternal
	case "device", "设备":
		return assetNetworkDevice
	case "domain", "domain_asset", "域名":
		return assetNetworkDomain
	default:
		return assetNetworkInternal
	}
}

// BuildDetailRiskStructuredResult 将 Batch E 安全结果收敛为 Product Facts StructuredResult candidate。
func BuildDetailRiskStructuredResult(result DetailRiskResult) (product.StructuredResultCandidate, string) {
	metadata := detailRiskToolMetadata(result.ToolID)
	itemCount := len(result.Items)
	if len(result.Metrics) > 0 {
		itemCount = len(result.Metrics)
	}
	target := strings.TrimSpace(result.Target)
	if target == "" {
		target = "默认范围"
	}
	summary := fmt.Sprintf("Fobrain %s：%s，返回 %d 条", metadata.DisplayName, target, itemCount)
	return product.StructuredResultCandidate{
		SchemaVersion: facts.StructuredResultSchemaVersion,
		ResultRef:     "result:fobrain:" + metadata.ResultSlug,
		SafeSummary:   summary,
		ItemCount:     itemCount,
	}, BusinessResultSchemaVersion
}

func detailRiskUnsafeValues(values ...string) bool {
	for _, value := range values {
		if facts.ContainsUnsafeMaterial(value) {
			return true
		}
	}
	return false
}
