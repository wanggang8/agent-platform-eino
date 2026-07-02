package fobrain

import (
	"fmt"
	"strings"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
)

// BuildCurrentUserStructuredResult 将 Fobrain 当前用户业务结果收敛为 StructuredResult candidate。
// 返回的 business schema 只用于展示 payload 校验，不作为 Product Facts 的第二套事实 schema。
func BuildCurrentUserStructuredResult(result CurrentUserContextResult) (product.StructuredResultCandidate, string) {
	return product.StructuredResultCandidate{
		SchemaVersion: facts.StructuredResultSchemaVersion,
		ResultRef:     "result:fobrain:current-user-context",
		SafeSummary:   currentUserSafeSummary(result),
	}, BusinessResultSchemaVersion
}

// BuildMyPermissionsStructuredResult 将权限范围结果收敛为 StructuredResult candidate。
func BuildMyPermissionsStructuredResult(result MyPermissionsResult) (product.StructuredResultCandidate, string) {
	return product.StructuredResultCandidate{
		SchemaVersion: facts.StructuredResultSchemaVersion,
		ResultRef:     "result:fobrain:my-permissions",
		SafeSummary:   permissionsSafeSummary(result),
	}, BusinessResultSchemaVersion
}

// BuildConnectorSecurityStructuredResult 将 connector 状态和凭据摘要收敛为 StructuredResult candidate。
func BuildConnectorSecurityStructuredResult(policyContext capabilities.PolicyContext, connectorStatus capabilities.ConnectorStatus) product.StructuredResultCandidate {
	binding := capabilities.SafeCredentialBinding(policyContext.CredentialBinding)
	status := "available"
	if connectorStatus == capabilities.ConnectorStatusUnavailable {
		status = "unavailable"
	}
	summary := fmt.Sprintf("Fobrain connector：%s，workspace：%s，绑定状态：%s", status, policyContext.WorkspaceID, binding.Status)
	return product.StructuredResultCandidate{
		SchemaVersion: facts.StructuredResultSchemaVersion,
		ResultRef:     "result:fobrain:connector-security",
		SafeSummary:   summary,
	}
}

// currentUserSafeSummary 生成可进入 StructuredResult 的中文安全摘要。
func currentUserSafeSummary(result CurrentUserContextResult) string {
	name := strings.TrimSpace(result.DisplayName)
	department := strings.TrimSpace(result.Department)
	role := strings.TrimSpace(result.Role)
	if name == "" {
		name = "当前用户"
	}
	parts := []string{"Fobrain 当前用户：" + name}
	if department != "" {
		parts = append(parts, "部门："+department)
	}
	if role != "" {
		parts = append(parts, "角色："+role)
	}
	return strings.Join(parts, "，")
}

// permissionsSafeSummary 生成权限范围的中文安全摘要，不包含 raw policy payload。
func permissionsSafeSummary(result MyPermissionsResult) string {
	name := strings.TrimSpace(result.DisplayName)
	if name == "" {
		name = "当前用户"
	}
	parts := []string{"Fobrain 权限范围：" + name}
	if len(result.PermissionNames) > 0 {
		parts = append(parts, "权限："+strings.Join(safeSummaryList(result.PermissionNames), "、"))
	}
	if len(result.DataPermissionNames) > 0 {
		parts = append(parts, "数据范围："+strings.Join(safeSummaryList(result.DataPermissionNames), "、"))
	}
	if len(result.PermissionNames) == 0 && len(result.DataPermissionNames) == 0 {
		parts = append(parts, "未返回权限字段")
	}
	return strings.Join(parts, "，")
}

func safeSummaryList(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		if text := strings.TrimSpace(value); text != "" {
			out = append(out, text)
		}
	}
	return out
}
