package fobrain

import (
	"fmt"
	"strings"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
)

const directReadEntityTicket = "ticket"

// DirectReadQuery 是 Batch C 直接列表和统计工具的安全入参。
type DirectReadQuery struct {
	BusinessName string
	Owner        string
	Keyword      string
	Field        string
	Severity     string
	Status       string
	TimeRange    string
	Person       string
	Page         int
	PageSize     int
}

// DirectReadResult 是 Batch C provider 边界内的安全业务结果，不承载 raw aggregation body。
type DirectReadResult struct {
	ToolID     string
	Title      string
	EntityType string
	Query      DirectReadQuery
	Items      []QueryResultItem
	Metrics    []RiskMetric
}

type directReadMetadata struct {
	DisplayName string
	EntityType  string
	ResultSlug  string
	Metrics     bool
}

func batchCDirectReadCapabilityByID(capabilityID string) bool {
	_, ok := batchCDirectReadMetadata(capabilityID)
	return ok
}

func directReadToolMetadata(toolID string) directReadMetadata {
	metadata, ok := batchCDirectReadMetadata(toolID)
	if !ok {
		return directReadMetadata{DisplayName: "Fobrain 直接读取", EntityType: "entity", ResultSlug: "direct-read"}
	}
	return metadata
}

func batchCDirectReadMetadata(toolID string) (directReadMetadata, bool) {
	switch toolID {
	case CapabilityBusinessList:
		return directReadMetadata{DisplayName: "查询业务系统", EntityType: myScopeEntityBusinessSystem, ResultSlug: "business-list"}, true
	case CapabilityExternalHighRiskAssets:
		return directReadMetadata{DisplayName: "查询外部高风险资产", EntityType: QueryEntityAsset, ResultSlug: "external-high-risk-assets"}, true
	case CapabilityVulnerabilityStatusSummary:
		return directReadMetadata{DisplayName: "汇总漏洞状态", EntityType: QueryEntityVulnerability, ResultSlug: "vulnerability-status-summary", Metrics: true}, true
	case CapabilityPendingTickets:
		return directReadMetadata{DisplayName: "查询待处理工单", EntityType: directReadEntityTicket, ResultSlug: "pending-tickets"}, true
	case CapabilityIPStats:
		return directReadMetadata{DisplayName: "统计 IP 资产", EntityType: QueryEntityAsset, ResultSlug: "ip-stats", Metrics: true}, true
	case CapabilityVulStats:
		return directReadMetadata{DisplayName: "统计漏洞情况", EntityType: QueryEntityVulnerability, ResultSlug: "vul-stats", Metrics: true}, true
	default:
		return directReadMetadata{}, false
	}
}

func directReadQueryFromArguments(toolID string, arguments map[string]any) (DirectReadQuery, error) {
	if _, ok := batchCDirectReadMetadata(toolID); !ok {
		return DirectReadQuery{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 直接读取能力未注册")
	}
	query := DirectReadQuery{
		BusinessName: safeArgumentString(arguments["business_name"]),
		Owner:        safeArgumentString(arguments["owner"]),
		Keyword:      safeArgumentString(arguments["keyword"]),
		Field:        safeArgumentString(arguments["field"]),
		Severity:     safeArgumentString(arguments["severity"]),
		Status:       safeArgumentString(arguments["status"]),
		TimeRange:    safeArgumentString(arguments["time_range"]),
		Person:       safeArgumentString(arguments["person"]),
		Page:         positiveInt(arguments["page"], 1),
		PageSize:     boundedPageSize(arguments["page_size"], 20),
	}
	if directReadUnsafeValues(query.BusinessName, query.Owner, query.Keyword, query.Field, query.Severity, query.Status, query.TimeRange, query.Person) {
		return DirectReadQuery{}, facts.ErrUnsafeFactMaterial
	}
	return query, nil
}

// BuildDirectReadStructuredResult 将 Batch C 安全结果收敛为 Product Facts StructuredResult candidate。
func BuildDirectReadStructuredResult(result DirectReadResult) (product.StructuredResultCandidate, string) {
	metadata := directReadToolMetadata(result.ToolID)
	// 摘要只使用本地元数据和已校验 query，避免 raw aggregation body 进入事实材料。
	itemCount := len(result.Items)
	if len(result.Metrics) > 0 {
		itemCount = len(result.Metrics)
	}
	summary := fmt.Sprintf("Fobrain %s：%s，返回 %d 条", metadata.DisplayName, directReadQueryTarget(metadata, result.Query), itemCount)
	return product.StructuredResultCandidate{
		SchemaVersion: facts.StructuredResultSchemaVersion,
		ResultRef:     "result:fobrain:" + metadata.ResultSlug,
		SafeSummary:   summary,
		ItemCount:     itemCount,
	}, BusinessResultSchemaVersion
}

func directReadQueryTarget(metadata directReadMetadata, query DirectReadQuery) string {
	for _, value := range []string{query.BusinessName, query.Owner, query.Keyword, query.Field, query.Severity, query.Status, query.Person, query.TimeRange} {
		if text := strings.TrimSpace(value); text != "" {
			return text
		}
	}
	if metadata.Metrics {
		return "默认统计范围"
	}
	return "默认列表范围"
}

func directReadUnsafeValues(values ...string) bool {
	for _, value := range values {
		if facts.ContainsUnsafeMaterial(value) {
			return true
		}
	}
	return false
}
