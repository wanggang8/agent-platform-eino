package capabilities_test

import (
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/capabilities"
)

func TestRegistryRegistersAndListsCapabilitiesByMetadata(t *testing.T) {
	// capability metadata 是工具选择和审批策略的依据，不能退化为工具名硬编码。
	registry := capabilities.NewRegistry()
	capability := capabilities.Capability{
		ID:          "cap.asset.search",
		ProviderID:  "provider.fixture",
		ToolName:    "asset_search",
		DisplayName: "资产查询",
		Description: "按授权范围查询资产",
		InputSchema: capabilities.JSONSchema{
			SchemaVersion: "json_schema.v1",
			Properties:    map[string]string{"query": "string"},
		},
		ResultSchema:     "tool.structured_result.v1",
		RiskLevel:        capabilities.RiskReadOnly,
		ApprovalRequired: false,
		Timeout:          5 * time.Second,
	}

	if err := registry.Register(capability); err != nil {
		t.Fatal(err)
	}

	list := registry.List()
	if len(list) != 1 {
		t.Fatalf("len(list) = %d, want 1", len(list))
	}
	if list[0].ToolName != "asset_search" || list[0].Description == "" {
		t.Fatalf("capability metadata not preserved: %+v", list[0])
	}
}

func TestRegistryRejectsDuplicateCapabilityID(t *testing.T) {
	registry := capabilities.NewRegistry()
	capability := capabilities.Capability{ID: "cap.duplicate", ToolName: "one", ProviderID: "provider.fixture"}

	if err := registry.Register(capability); err != nil {
		t.Fatal(err)
	}
	if err := registry.Register(capability); err == nil {
		t.Fatal("duplicate Register error = nil")
	}
}

func TestPolicyRequiresApprovalForWriteRisk(t *testing.T) {
	// 写域能力默认需要 approval，provider 或前端不能绕过 policy。
	decision := capabilities.EvaluatePolicy(capabilities.Capability{
		ID:               "cap.ticket.close",
		ToolName:         "ticket_close",
		RiskLevel:        capabilities.RiskWrite,
		ApprovalRequired: false,
	})

	if decision.Allowed {
		t.Fatal("write capability without approval metadata must not be directly allowed")
	}
	if !decision.RequiresApproval {
		t.Fatal("write capability must require approval")
	}
}
