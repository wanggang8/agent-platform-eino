package fobrain

import (
	"fmt"
	"strings"

	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
)

const (
	EntityPerson        = "person"
	EntityAsset         = "asset"
	EntityVulnerability = "vulnerability"
	EntityDepartment    = "department"
	EntityBusiness      = "business"

	EntityResolutionResolved = "resolved"
	EntityResolutionWaiting  = "waiting"
	EntityResolutionNotFound = "not_found"
)

// EntityCandidate 是 provider 边界内的待消歧候选，RawRef/UnsafeNotes 仅用于边界内丢弃验证。
type EntityCandidate struct {
	DisplayName    string
	DepartmentName string
	RoleName       string
	Relation       string
	Alias          string
	Score          float64
	RawRef         string
	UnsafeNotes    string
}

// EntityResolutionResult 是 Fobrain 查询目标消歧后的安全结果。
// Resolved/Candidates 都已经收敛为 facts.PendingCandidate，不包含 provider 原始 id。
type EntityResolutionResult struct {
	Status          string
	EntityType      string
	Query           string
	Resolved        *facts.PendingCandidate
	Candidates      []facts.PendingCandidate
	AttentionReason string
}

// ClassifyEntityCandidates 根据查询词和候选列表判断是否可直接解析，或需要 clarification。
func ClassifyEntityCandidates(query string, entityType string, candidates []EntityCandidate, limit int) EntityResolutionResult {
	query = strings.TrimSpace(query)
	entityType = normalizeEntityType(entityType)
	safeCandidates := safeEntityCandidates(entityType, candidates, limit)
	result := EntityResolutionResult{
		Status:     EntityResolutionNotFound,
		EntityType: entityType,
		Query:      query,
	}
	if len(safeCandidates) == 0 {
		result.AttentionReason = entityType + "_not_found"
		return result
	}

	exact := exactCandidates(query, safeCandidates)
	if len(exact) == 1 {
		resolved := exact[0]
		result.Status = EntityResolutionResolved
		result.Resolved = &resolved
		return result
	}
	if len(safeCandidates) == 1 && query != "" {
		resolved := safeCandidates[0]
		result.Status = EntityResolutionResolved
		result.Resolved = &resolved
		return result
	}

	result.Status = EntityResolutionWaiting
	result.Candidates = safeCandidates
	result.AttentionReason = entityType + "_disambiguation_required"
	return result
}

// BuildEntityResolutionStructuredResult 生成实体消歧的 StructuredResult 候选摘要。
func BuildEntityResolutionStructuredResult(result EntityResolutionResult) (product.StructuredResultCandidate, string) {
	statusLabel := "未找到"
	switch result.Status {
	case EntityResolutionResolved:
		if result.Resolved != nil {
			statusLabel = "已解析：" + result.Resolved.Label
		} else {
			statusLabel = "已解析"
		}
	case EntityResolutionWaiting:
		statusLabel = fmt.Sprintf("需要澄清，候选 %d 个", len(result.Candidates))
	}
	entityType := normalizeEntityType(result.EntityType)
	if entityType == "" {
		entityType = "entity"
	}
	return product.StructuredResultCandidate{
		SchemaVersion: facts.StructuredResultSchemaVersion,
		ResultRef:     "result:fobrain:entity-resolution:" + entityType,
		SafeSummary:   fmt.Sprintf("Fobrain 实体消歧：%s，类型：%s", statusLabel, entityType),
	}, BusinessResultSchemaVersion
}

func safeEntityCandidates(entityType string, candidates []EntityCandidate, limit int) []facts.PendingCandidate {
	if limit <= 0 {
		limit = 5
	}
	out := make([]facts.PendingCandidate, 0, minInt(limit, len(candidates)))
	for index, candidate := range candidates {
		if len(out) >= limit {
			break
		}
		safe := facts.PendingCandidate{
			CandidateRef: fmt.Sprintf("candidate:fobrain:%s:%d", entityType, index+1),
			Label:        safeCandidateValue(candidate.DisplayName),
			Description:  safeDescription(candidate),
			EntityType:   entityType,
			SafeFields:   safeCandidateFields(candidate),
		}
		if safe.Label == "" {
			safe.Label = entityType + " 候选"
		}
		if facts.UnsafePendingCandidates([]facts.PendingCandidate{safe}) {
			continue
		}
		out = append(out, safe)
	}
	return out
}

func exactCandidates(query string, candidates []facts.PendingCandidate) []facts.PendingCandidate {
	if query == "" {
		return nil
	}
	normalized := normalizeText(query)
	var exact []facts.PendingCandidate
	for _, candidate := range candidates {
		if normalizeText(candidate.Label) == normalized {
			exact = append(exact, candidate)
			continue
		}
		for _, field := range candidate.SafeFields {
			if field.Label == "别名" && normalizeText(field.Value) == normalized {
				exact = append(exact, candidate)
				break
			}
		}
	}
	return exact
}

func safeCandidateFields(candidate EntityCandidate) []facts.PendingCandidateField {
	fields := []facts.PendingCandidateField{}
	appendField := func(label string, value string) {
		value = safeCandidateValue(value)
		if value == "" {
			return
		}
		field := facts.PendingCandidateField{Label: label, Value: value}
		if !facts.UnsafePendingCandidates([]facts.PendingCandidate{{CandidateRef: "candidate:fobrain:safe:field", Label: "候选", EntityType: EntityPerson, SafeFields: []facts.PendingCandidateField{field}}}) {
			fields = append(fields, field)
		}
	}
	appendField("部门", candidate.DepartmentName)
	appendField("角色", candidate.RoleName)
	appendField("关系", candidate.Relation)
	appendField("别名", candidate.Alias)
	return fields
}

func safeDescription(candidate EntityCandidate) string {
	parts := []string{}
	for _, value := range []string{candidate.DepartmentName, candidate.RoleName, candidate.Relation} {
		if safe := safeCandidateValue(value); safe != "" {
			parts = append(parts, safe)
		}
	}
	return strings.Join(parts, " / ")
}

func safeCandidateValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || facts.UnsafePendingCandidates([]facts.PendingCandidate{{CandidateRef: "candidate:fobrain:safe:value", Label: value, EntityType: EntityPerson}}) {
		return ""
	}
	return value
}

func normalizeEntityType(entityType string) string {
	switch strings.TrimSpace(entityType) {
	case EntityPerson, EntityAsset, EntityVulnerability, EntityDepartment, EntityBusiness:
		return strings.TrimSpace(entityType)
	default:
		return "entity"
	}
}

func normalizeText(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func minInt(left int, right int) int {
	if left < right {
		return left
	}
	return right
}
