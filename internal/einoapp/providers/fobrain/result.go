package fobrain

import (
	"strings"

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
