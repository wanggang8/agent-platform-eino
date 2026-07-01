package execution_test

import (
	"testing"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/execution"
)

func TestCapabilitySelectionDefaultsNaturalLanguageToChat(t *testing.T) {
	registry := capabilities.NewRegistry()
	if err := registry.Register(capabilities.Capability{
		ID:          "fobrain.ip.read",
		ProviderID:  "fobrain",
		ToolName:    "fobrain_ip_read",
		DisplayName: "IP 查询",
		RiskLevel:   capabilities.RiskReadOnly,
	}); err != nil {
		t.Fatal(err)
	}

	// 即使文本看起来像 Fobrain 场景，也不能绕过模型/tool call 直接硬编码选工具。
	selection := execution.SelectCapability(registry, execution.SelectionRequest{
		InputText: "帮我查一下这个 IP",
	})

	if selection.Mode != execution.SelectionModeChat {
		t.Fatalf("natural language must default to chat, got %+v", selection)
	}
	if selection.CapabilityID != "" {
		t.Fatalf("natural language must not hardcode tool selection: %+v", selection)
	}
}

func TestCapabilitySelectionUsesHintThroughRegistryAndPolicy(t *testing.T) {
	registry := capabilities.NewRegistry()
	if err := registry.Register(capabilities.Capability{
		ID:          "safe.read",
		ProviderID:  "demo",
		ToolName:    "safe_read",
		DisplayName: "只读查询",
		RiskLevel:   capabilities.RiskReadOnly,
	}); err != nil {
		t.Fatal(err)
	}
	if err := registry.Register(capabilities.Capability{
		ID:          "danger.write",
		ProviderID:  "demo",
		ToolName:    "danger_write",
		DisplayName: "写入",
		RiskLevel:   capabilities.RiskWrite,
	}); err != nil {
		t.Fatal(err)
	}

	allowed := execution.SelectCapability(registry, execution.SelectionRequest{CapabilityHint: "safe.read"})
	if allowed.Mode != execution.SelectionModeCapability || allowed.CapabilityID != "safe.read" || allowed.RequiresApproval {
		t.Fatalf("read capability selection mismatch: %+v", allowed)
	}

	write := execution.SelectCapability(registry, execution.SelectionRequest{CapabilityHint: "danger.write"})
	if write.Mode != execution.SelectionModeCapability || !write.RequiresApproval || write.PolicyReason != "write_requires_approval" {
		t.Fatalf("write capability policy mismatch: %+v", write)
	}
}

func TestCapabilitySelectionRejectsUnknownHint(t *testing.T) {
	registry := capabilities.NewRegistry()

	selection := execution.SelectCapability(registry, execution.SelectionRequest{CapabilityHint: "missing.tool"})

	if selection.Mode != execution.SelectionModeRejected || selection.PolicyReason != "capability_not_registered" {
		t.Fatalf("unknown hint selection mismatch: %+v", selection)
	}
}
