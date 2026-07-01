package main

import (
	"testing"

	"agent-platform-eino/internal/einoapp/bootstrap"
	"agent-platform-eino/internal/einoapp/capabilities"
)

func TestCapabilityRegistryFromConfigDoesNotRegisterImplicitCapabilities(t *testing.T) {
	// 启动路径不能内置 smoke 或业务能力；未配置时 registry 必须为空。
	registry, err := capabilityRegistryFromConfig(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(registry.List()) != 0 {
		t.Fatalf("registry should be empty without config: %+v", registry.List())
	}
}

func TestCapabilityRegistryFromConfigRegistersConfiguredCapabilities(t *testing.T) {
	// capability hint 只能命中配置文件声明的能力元数据。
	registry, err := capabilityRegistryFromConfig([]bootstrap.CapabilityConfig{
		{
			ID:          "cap.smoke.read",
			ProviderID:  "phase3-smoke",
			ToolName:    "phase3_smoke_read",
			DisplayName: "Phase 3 只读验证",
			Description: "验证配置驱动能力注册",
			RiskLevel:   string(capabilities.RiskReadOnly),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	capability, ok := registry.Get("cap.smoke.read")
	if !ok {
		t.Fatal("configured capability was not registered")
	}
	if capability.ProviderID != "phase3-smoke" || capability.RiskLevel != capabilities.RiskReadOnly {
		t.Fatalf("configured capability mismatch: %+v", capability)
	}
}
