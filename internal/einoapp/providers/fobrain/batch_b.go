package fobrain

import (
	"fmt"
	"strings"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
)

const myScopeEntityBusinessSystem = "business_system"

// MyScopeQuery 是 Batch B “我的范围”工具的安全入参，只保留关键字和分页。
type MyScopeQuery struct {
	Keyword       string
	Page          int
	PageSize      int
	Department    bool
	ImportantOnly bool
}

// MyScopeResult 是 Batch B provider 边界内的安全业务结果，不承载 raw user payload。
type MyScopeResult struct {
	ToolID     string
	Title      string
	EntityType string
	Query      MyScopeQuery
	Items      []QueryResultItem
}

type myScopeMetadata struct {
	DisplayName   string
	EntityType    string
	ResultSlug    string
	Department    bool
	ImportantOnly bool
}

func batchBMyScopeCapabilityByID(capabilityID string) bool {
	_, ok := batchBMyScopeMetadata(capabilityID)
	return ok
}

func myScopeToolMetadata(toolID string) myScopeMetadata {
	metadata, ok := batchBMyScopeMetadata(toolID)
	if !ok {
		return myScopeMetadata{DisplayName: "Fobrain 我的范围查询", EntityType: "entity", ResultSlug: "my-scope"}
	}
	return metadata
}

func batchBMyScopeMetadata(toolID string) (myScopeMetadata, bool) {
	switch toolID {
	case CapabilityMyAssets:
		return myScopeMetadata{DisplayName: "查询我的资产", EntityType: QueryEntityAsset, ResultSlug: "my-assets"}, true
	case CapabilityMyDepartmentAssets:
		return myScopeMetadata{DisplayName: "查询本部门资产", EntityType: QueryEntityAsset, ResultSlug: "my-department-assets", Department: true}, true
	case CapabilityMyVulnerabilities:
		return myScopeMetadata{DisplayName: "查询我的漏洞", EntityType: QueryEntityVulnerability, ResultSlug: "my-vulnerabilities"}, true
	case CapabilityMyDepartmentVulnerabilities:
		return myScopeMetadata{DisplayName: "查询本部门漏洞", EntityType: QueryEntityVulnerability, ResultSlug: "my-department-vulnerabilities", Department: true}, true
	case CapabilityMyBusinessSystems:
		return myScopeMetadata{DisplayName: "查询我的业务系统", EntityType: myScopeEntityBusinessSystem, ResultSlug: "my-business-systems"}, true
	case CapabilityMyImportantBusinessSystems:
		return myScopeMetadata{DisplayName: "查询我的重要业务系统", EntityType: myScopeEntityBusinessSystem, ResultSlug: "my-important-business-systems", ImportantOnly: true}, true
	default:
		return myScopeMetadata{}, false
	}
}

func myScopeQueryFromArguments(toolID string, arguments map[string]any) (MyScopeQuery, error) {
	metadata, ok := batchBMyScopeMetadata(toolID)
	if !ok {
		return MyScopeQuery{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 我的范围能力未注册")
	}
	query := MyScopeQuery{
		Keyword:       safeArgumentString(arguments["keyword"]),
		Page:          positiveInt(arguments["page"], 1),
		PageSize:      boundedPageSize(arguments["page_size"], 20),
		Department:    metadata.Department,
		ImportantOnly: metadata.ImportantOnly,
	}
	if facts.ContainsUnsafeMaterial(query.Keyword) {
		return MyScopeQuery{}, facts.ErrUnsafeFactMaterial
	}
	return query, nil
}

// BuildMyScopeStructuredResult 将 Batch B 安全结果收敛为 Product Facts StructuredResult candidate。
func BuildMyScopeStructuredResult(result MyScopeResult) (product.StructuredResultCandidate, string) {
	metadata := myScopeToolMetadata(result.ToolID)
	// 摘要只使用本地元数据和已校验 query，避免 client 回填 raw 用户上下文。
	itemCount := len(result.Items)
	summary := fmt.Sprintf("Fobrain %s：%s，返回 %d 条", metadata.DisplayName, myScopeQueryTarget(metadata, result.Query), itemCount)
	return product.StructuredResultCandidate{
		SchemaVersion: facts.StructuredResultSchemaVersion,
		ResultRef:     "result:fobrain:" + metadata.ResultSlug,
		SafeSummary:   summary,
		ItemCount:     itemCount,
	}, BusinessResultSchemaVersion
}

func myScopeQueryTarget(metadata myScopeMetadata, query MyScopeQuery) string {
	if text := strings.TrimSpace(query.Keyword); text != "" {
		return text
	}
	if query.ImportantOnly || metadata.ImportantOnly {
		return "当前用户重要业务范围"
	}
	if query.Department || metadata.Department {
		return "本部门范围"
	}
	return "当前用户范围"
}
