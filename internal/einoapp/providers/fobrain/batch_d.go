package fobrain

import (
	"fmt"
	"strings"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
)

const (
	QueryEntityAsset         = "asset"
	QueryEntityVulnerability = "vulnerability"
)

// ParameterizedQuery 是 Batch D 工具的安全入参模型，只保留业务查询值。
type ParameterizedQuery struct {
	PersonName     string
	PersonStaffID  string
	DepartmentName string
	IP             string
	Severity       string
	Status         string
	Page           int
	PageSize       int
}

// QueryResultItem 是 Batch D mock mapper 的安全行摘要，不承载 raw provider payload。
type QueryResultItem struct {
	EntityRef   string
	DisplayName string
	OwnerName   string
	Status      string
	Severity    string
	Affected    int
	Summary     string
}

// ParameterizedQueryResult 是 Batch D provider 边界内的安全业务结果。
type ParameterizedQueryResult struct {
	ToolID     string
	Title      string
	EntityType string
	Query      ParameterizedQuery
	Items      []QueryResultItem
}

type parameterizedMetadata struct {
	DisplayName string
	EntityType  string
	Required    string
	ResultSlug  string
}

func batchDParameterizedCapabilityByID(capabilityID string) bool {
	_, ok := batchDParameterizedMetadata(capabilityID)
	return ok
}

func parameterizedToolMetadata(toolID string) parameterizedMetadata {
	metadata, ok := batchDParameterizedMetadata(toolID)
	if !ok {
		return parameterizedMetadata{DisplayName: "Fobrain 参数化查询", EntityType: "entity", ResultSlug: "parameterized-query"}
	}
	return metadata
}

func batchDParameterizedMetadata(toolID string) (parameterizedMetadata, bool) {
	switch toolID {
	case CapabilityListAssetsByOwner:
		return parameterizedMetadata{DisplayName: "查询负责人资产", EntityType: QueryEntityAsset, Required: "person_name", ResultSlug: "list-assets-by-owner"}, true
	case CapabilityListVulnerabilitiesByOwner:
		return parameterizedMetadata{DisplayName: "查询负责人漏洞", EntityType: QueryEntityVulnerability, Required: "person_name", ResultSlug: "list-vulnerabilities-by-owner"}, true
	case CapabilityListAssetsByDepartment:
		return parameterizedMetadata{DisplayName: "查询部门资产", EntityType: QueryEntityAsset, Required: "department_name", ResultSlug: "list-assets-by-department"}, true
	case CapabilityListVulnerabilitiesByDepartment:
		return parameterizedMetadata{DisplayName: "查询部门漏洞", EntityType: QueryEntityVulnerability, Required: "department_name", ResultSlug: "list-vulnerabilities-by-department"}, true
	case CapabilityListAssetsByIP:
		return parameterizedMetadata{DisplayName: "查询 IP 资产", EntityType: QueryEntityAsset, Required: "ip", ResultSlug: "list-assets-by-ip"}, true
	case CapabilityListVulnerabilitiesByIP:
		return parameterizedMetadata{DisplayName: "查询 IP 漏洞", EntityType: QueryEntityVulnerability, Required: "ip", ResultSlug: "list-vulnerabilities-by-ip"}, true
	default:
		return parameterizedMetadata{}, false
	}
}

func parameterizedQueryFromArguments(toolID string, arguments map[string]any) (ParameterizedQuery, error) {
	metadata, ok := batchDParameterizedMetadata(toolID)
	if !ok {
		return ParameterizedQuery{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 参数化查询能力未注册")
	}
	query := ParameterizedQuery{
		PersonName:     safeArgumentString(arguments["person_name"]),
		PersonStaffID:  safeArgumentString(arguments["person_staff_id"]),
		DepartmentName: safeArgumentString(arguments["department_name"]),
		IP:             safeArgumentString(arguments["ip"]),
		Severity:       safeArgumentString(arguments["severity"]),
		Status:         safeArgumentString(arguments["status"]),
		Page:           positiveInt(arguments["page"], 1),
		PageSize:       boundedPageSize(arguments["page_size"], 20),
	}
	if !requiredQueryValuePresent(metadata.Required, query) {
		return ParameterizedQuery{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 参数化查询缺少必填参数")
	}
	if unsafeQuery(query) {
		return ParameterizedQuery{}, facts.ErrUnsafeFactMaterial
	}
	return query, nil
}

// BuildParameterizedQueryStructuredResult 将 Batch D 安全结果收敛为 Product Facts StructuredResult candidate。
func BuildParameterizedQueryStructuredResult(result ParameterizedQueryResult) (product.StructuredResultCandidate, string) {
	metadata := parameterizedToolMetadata(result.ToolID)
	// 标题取自本地元数据，避免 provider 回填文案把 raw 信息带入事实摘要。
	title := metadata.DisplayName
	itemCount := len(result.Items)
	summary := fmt.Sprintf("Fobrain %s：%s，返回 %d 条", title, parameterizedQueryTarget(result.Query), itemCount)
	return product.StructuredResultCandidate{
		SchemaVersion: facts.StructuredResultSchemaVersion,
		ResultRef:     "result:fobrain:" + metadata.ResultSlug,
		SafeSummary:   summary,
	}, BusinessResultSchemaVersion
}

func requiredQueryValuePresent(required string, query ParameterizedQuery) bool {
	switch required {
	case "person_name":
		return strings.TrimSpace(query.PersonName) != ""
	case "department_name":
		return strings.TrimSpace(query.DepartmentName) != ""
	case "ip":
		return strings.TrimSpace(query.IP) != ""
	default:
		return false
	}
}

func parameterizedQueryTarget(query ParameterizedQuery) string {
	for _, value := range []string{query.PersonName, query.DepartmentName, query.IP} {
		if text := strings.TrimSpace(value); text != "" {
			return text
		}
	}
	return "默认范围"
}

func unsafeQuery(query ParameterizedQuery) bool {
	for _, value := range []string{query.PersonName, query.PersonStaffID, query.DepartmentName, query.IP, query.Severity, query.Status} {
		if facts.ContainsUnsafeMaterial(value) {
			return true
		}
	}
	return false
}

func safeArgumentString(value any) string {
	return strings.TrimSpace(stringFromAny(value))
}

func positiveInt(value any, fallback int) int {
	number := intFromAny(value)
	if number <= 0 {
		return fallback
	}
	return number
}

func boundedPageSize(value any, fallback int) int {
	number := positiveInt(value, fallback)
	if number > 100 {
		return 100
	}
	return number
}

func intFromAny(value any) int {
	switch typed := value.(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	case jsonNumber:
		number, _ := typed.Int64()
		return int(number)
	default:
		return 0
	}
}

type jsonNumber interface {
	Int64() (int64, error)
}
