package fobrain_test

import (
	"errors"
	"strings"
	"testing"

	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
	"agent-platform-eino/internal/einoapp/providers/fobrain"
)

func TestReadonlyPoCStructuredResultUsesFactsSchemaAndBusinessSchema(t *testing.T) {
	// Product Facts 只能保存 StructuredResult；Fobrain 业务 schema 只作为展示 payload 契约。
	candidate, businessSchema := fobrain.BuildCurrentUserStructuredResult(fobrain.CurrentUserContextResult{
		DisplayName: "张三",
		Department:  "安全部",
		Role:        "安全运营",
	})

	if candidate.SchemaVersion != product.StructuredResultSchemaVersion {
		t.Fatalf("structured schema = %q", candidate.SchemaVersion)
	}
	if businessSchema != fobrain.BusinessResultSchemaVersion {
		t.Fatalf("business schema = %q", businessSchema)
	}
	if candidate.ResultRef != "result:fobrain:current-user-context" {
		t.Fatalf("result ref = %q", candidate.ResultRef)
	}
	for _, expected := range []string{"Fobrain 当前用户：张三", "部门：安全部", "角色：安全运营"} {
		if !strings.Contains(candidate.SafeSummary, expected) {
			t.Fatalf("safe summary missing %q: %s", expected, candidate.SafeSummary)
		}
	}
	if _, err := product.NewStructuredResultSafetyGate().Approve(candidate); err != nil {
		t.Fatalf("safe candidate was rejected: %v", err)
	}
}

func TestStructuredResultRejectsUnsafeCurrentUserMaterial(t *testing.T) {
	candidate, _ := fobrain.BuildCurrentUserStructuredResult(fobrain.CurrentUserContextResult{
		DisplayName: "Authorization: Bearer secret",
		Department:  "安全部",
	})

	_, err := product.NewStructuredResultSafetyGate().Approve(candidate)
	if !errors.Is(err, facts.ErrUnsafeFactMaterial) {
		t.Fatalf("unsafe current user candidate err = %v", err)
	}
}
