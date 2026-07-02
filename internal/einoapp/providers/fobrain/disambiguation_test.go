package fobrain_test

import (
	"strings"
	"testing"

	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
	"agent-platform-eino/internal/einoapp/providers/fobrain"
)

func TestDisambiguationClassifiesAmbiguousPersonCandidates(t *testing.T) {
	// 多候选实体只能输出安全候选引用，不能把手机号、邮箱、内部人员 id 或 raw provider id 写入候选事实。
	result := fobrain.ClassifyEntityCandidates("张三", fobrain.EntityPerson, []fobrain.EntityCandidate{
		{DisplayName: "张三", DepartmentName: "安全部", RoleName: "安全运营", Alias: "zhangsan", RawRef: "person-raw-001", UnsafeNotes: "phone=13800138000"},
		{DisplayName: "张三", DepartmentName: "研发部", RoleName: "后端工程师", Alias: "zhangsan2", RawRef: "employee-002", UnsafeNotes: "email=zhangsan@example.com"},
		{DisplayName: "张三", DepartmentName: "运维部", RoleName: "值班工程师"},
	}, 2)

	if result.Status != fobrain.EntityResolutionWaiting || result.AttentionReason != "person_disambiguation_required" {
		t.Fatalf("result status mismatch: %+v", result)
	}
	if len(result.Candidates) != 2 {
		t.Fatalf("candidates length = %d, want 2: %+v", len(result.Candidates), result.Candidates)
	}
	if result.Candidates[0].CandidateRef != "candidate:fobrain:person:1" || result.Candidates[0].Label != "张三" {
		t.Fatalf("candidate safe ref/label mismatch: %+v", result.Candidates[0])
	}
	joined := strings.ToLower(candidateText(result.Candidates))
	for _, forbidden := range []string{"13800138000", "example.com", "person-raw", "employee", "raw", "phone", "email", "token", "credential"} {
		if strings.Contains(joined, forbidden) {
			t.Fatalf("candidate leaked %q: %s", forbidden, joined)
		}
	}
}

func TestDisambiguationResolvesExactSingleCandidate(t *testing.T) {
	result := fobrain.ClassifyEntityCandidates("支付系统", fobrain.EntityBusiness, []fobrain.EntityCandidate{
		{DisplayName: "支付系统", DepartmentName: "金融事业部", Relation: "owner"},
	}, 5)

	if result.Status != fobrain.EntityResolutionResolved || result.Resolved == nil {
		t.Fatalf("exact single candidate should resolve: %+v", result)
	}
	if result.Resolved.EntityType != fobrain.EntityBusiness || result.Resolved.CandidateRef == "" {
		t.Fatalf("resolved candidate mismatch: %+v", result.Resolved)
	}
}

func TestBuildEntityResolutionStructuredResultUsesSafeCandidateFacts(t *testing.T) {
	resolution := fobrain.ClassifyEntityCandidates("张三", fobrain.EntityPerson, []fobrain.EntityCandidate{
		{DisplayName: "张三", DepartmentName: "安全部"},
		{DisplayName: "张三", DepartmentName: "研发部"},
	}, 5)

	candidate, businessSchema := fobrain.BuildEntityResolutionStructuredResult(resolution)
	if businessSchema != fobrain.BusinessResultSchemaVersion {
		t.Fatalf("business schema = %q", businessSchema)
	}
	if candidate.SchemaVersion != facts.StructuredResultSchemaVersion || candidate.ResultRef != "result:fobrain:entity-resolution:person" {
		t.Fatalf("structured candidate mismatch: %+v", candidate)
	}
	if !strings.Contains(candidate.SafeSummary, "候选 2 个") || strings.Contains(strings.ToLower(candidate.SafeSummary), "raw") {
		t.Fatalf("safe summary mismatch: %s", candidate.SafeSummary)
	}
	if _, err := product.NewStructuredResultSafetyGate().Approve(candidate); err != nil {
		t.Fatalf("entity resolution candidate rejected: %v", err)
	}
}

func candidateText(candidates []facts.PendingCandidate) string {
	var builder strings.Builder
	for _, candidate := range candidates {
		builder.WriteString(candidate.CandidateRef)
		builder.WriteString(candidate.Label)
		builder.WriteString(candidate.Description)
		builder.WriteString(candidate.EntityType)
		for _, field := range candidate.SafeFields {
			builder.WriteString(field.Label)
			builder.WriteString(field.Value)
		}
	}
	return builder.String()
}
