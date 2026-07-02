package main

import (
	"context"
	"testing"

	"agent-platform-eino/internal/einoapp/bootstrap"
	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/llm"
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
			ID:                      "cap.smoke.read",
			ProviderID:              "phase3-smoke",
			ToolName:                "phase3_smoke_read",
			DisplayName:             "Phase 3 只读验证",
			Description:             "验证配置驱动能力注册",
			RiskLevel:               string(capabilities.RiskReadOnly),
			SideEffect:              string(capabilities.SideEffectReadExternal),
			PolicyRef:               "policy:smoke:read:v1",
			PermissionScope:         string(capabilities.PermissionScopeWorkspace),
			CredentialBindingPolicy: string(capabilities.CredentialBindingNone),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	capability, ok := registry.Get("cap.smoke.read")
	if !ok {
		t.Fatal("configured capability was not registered")
	}
	if capability.ProviderID != "phase3-smoke" ||
		capability.RiskLevel != capabilities.RiskReadOnly ||
		capability.SideEffect != capabilities.SideEffectReadExternal ||
		capability.PolicyRef != "policy:smoke:read:v1" ||
		capability.PermissionScope != capabilities.PermissionScopeWorkspace ||
		capability.CredentialBindingPolicy != capabilities.CredentialBindingNone {
		t.Fatalf("configured capability mismatch: %+v", capability)
	}
}

func TestLLMProviderFromConfigSupportsMockAndOpenAICompatible(t *testing.T) {
	mockProvider, err := llmProviderFromConfig(bootstrap.LLMConfig{Provider: "mock"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := mockProvider.(llm.MockProvider); !ok {
		t.Fatalf("mock provider type = %T", mockProvider)
	}

	realProvider, err := llmProviderFromConfig(bootstrap.LLMConfig{Provider: "openai_compatible", APIKey: "sk-local"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := realProvider.(llm.OpenAICompatibleProvider); !ok {
		t.Fatalf("openai compatible provider type = %T", realProvider)
	}
}

func TestLLMProviderFromConfigRejectsUnsupportedProvider(t *testing.T) {
	_, err := llmProviderFromConfig(bootstrap.LLMConfig{Provider: "unsupported"})
	if err == nil {
		t.Fatal("unsupported provider err = nil")
	}
}

func TestCapabilityRuntimeFromConfigRegistersAndInvokesMockMCP(t *testing.T) {
	registry, invoker, err := capabilityRuntimeFromConfig(bootstrap.Config{
		MCPMockServers: []bootstrap.MCPMockServerConfig{
			{
				ServerID: "mock",
				Tools: []bootstrap.MCPMockToolConfig{
					{
						Name:        "asset_lookup",
						Title:       "MCP 资产查询",
						Description: "通过 mock MCP 查询资产",
						InputSchema: bootstrap.MCPJSONSchemaConfig{
							Type:       "object",
							Properties: map[string]string{"query": "string"},
							Required:   []string{"query"},
						},
						OutputSchema: bootstrap.MCPJSONSchemaConfig{Type: "object"},
						Annotations:  bootstrap.MCPToolAnnotationsConfig{ReadOnlyHint: true},
					},
				},
				Results: map[string]bootstrap.MCPMockToolResultConfig{
					"asset_lookup": {
						StructuredContent: map[string]string{
							"result_ref":   "result:mcp:asset_lookup",
							"safe_summary": "MCP 资产查询完成",
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	capability, ok := registry.Get("mcp.mock.tool.asset_lookup")
	if !ok {
		t.Fatalf("mcp capability was not registered: %+v", registry.List())
	}
	if capability.ProviderID != "mcp:mock" || capability.ToolName != "mcp_asset_lookup" {
		t.Fatalf("mcp capability metadata mismatch: %+v", capability)
	}

	candidate, err := invoker.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID: "mcp.mock.tool.asset_lookup",
		Arguments:    map[string]any{"query": "asset"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if candidate.SchemaVersion != facts.StructuredResultSchemaVersion || candidate.SafeSummary != "MCP 资产查询完成" {
		t.Fatalf("mcp candidate mismatch: %+v", candidate)
	}
}
