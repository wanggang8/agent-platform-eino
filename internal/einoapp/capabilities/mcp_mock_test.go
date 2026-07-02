package capabilities_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
)

func TestMCPLifecycleInitializesMockSession(t *testing.T) {
	provider := newMockMCPProviderForTest(t)

	session, err := provider.Initialize(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if !session.Initialized || !session.ToolsListChanged {
		t.Fatalf("session did not negotiate tools capability: %+v", session)
	}
	if session.ProtocolVersion != capabilities.MCPProtocolVersion20250618 {
		t.Fatalf("protocol version = %q", session.ProtocolVersion)
	}
}

func TestMCPToolCallRequiresInitializedSession(t *testing.T) {
	provider := newMockMCPProviderForTest(t)

	_, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID: "mcp.mock.tool.asset_lookup",
		Arguments:    map[string]any{"query": "asset"},
	})
	if !errors.Is(err, capabilities.ErrMCPSessionNotInitialized) {
		t.Fatalf("invoke before initialize err = %v", err)
	}
}

func TestMCPToolListPagination(t *testing.T) {
	provider := newMockMCPProviderForTest(t)
	if _, err := provider.Initialize(context.Background()); err != nil {
		t.Fatal(err)
	}

	firstPage, err := provider.ListTools(context.Background(), capabilities.MCPListToolsRequest{PageSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(firstPage.Tools) != 1 || firstPage.NextCursor == "" {
		t.Fatalf("first page mismatch: %+v", firstPage)
	}

	secondPage, err := provider.ListTools(context.Background(), capabilities.MCPListToolsRequest{
		Cursor:   firstPage.NextCursor,
		PageSize: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(secondPage.Tools) != 1 || secondPage.Tools[0].Name == firstPage.Tools[0].Name || secondPage.NextCursor != "" {
		t.Fatalf("second page mismatch: %+v", secondPage)
	}
}

func TestMCPSchemaConversionRegistersCapabilities(t *testing.T) {
	provider := newInitializedMockMCPProviderForTest(t)

	capabilitiesList, err := provider.ListCapabilities()
	if err != nil {
		t.Fatal(err)
	}

	if len(capabilitiesList) != 2 {
		t.Fatalf("capabilities len = %d", len(capabilitiesList))
	}
	readCapability := capabilitiesList[0]
	if readCapability.ID != "mcp.mock.tool.asset_lookup" || readCapability.ToolName != "mcp_asset_lookup" {
		t.Fatalf("capability metadata mismatch: %+v", readCapability)
	}
	if readCapability.InputSchema.Properties["query"] != "string" || readCapability.ResultSchema != facts.StructuredResultSchemaVersion {
		t.Fatalf("schema conversion mismatch: %+v", readCapability)
	}
	if len(readCapability.InputSchema.Required) != 1 || readCapability.InputSchema.Required[0] != "query" {
		t.Fatalf("required schema not preserved: %+v", readCapability.InputSchema)
	}
}

func TestMCPSchemaConversionRejectsInvalidArguments(t *testing.T) {
	provider := newInitializedMockMCPProviderForTest(t)

	_, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID: "mcp.mock.tool.asset_lookup",
		Arguments:    map[string]any{},
	})
	if !errors.Is(err, capabilities.ErrMCPInvalidArguments) {
		t.Fatalf("missing required argument err = %v", err)
	}

	_, err = provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID: "mcp.mock.tool.asset_lookup",
		Arguments:    map[string]any{"query": 42},
	})
	if !errors.Is(err, capabilities.ErrMCPInvalidArguments) {
		t.Fatalf("wrong typed argument err = %v", err)
	}
}

func TestMCPSchemaConversionRedactsUnsafeToolMetadata(t *testing.T) {
	provider := capabilities.NewMockMCPProvider(capabilities.MCPMockConfig{
		ServerID: "mock",
		Tools: []capabilities.MCPToolDefinition{
			{
				Name:        "unsafe_tool",
				Title:       "token=local-secret",
				Description: "raw_payload token=local-secret",
				InputSchema: capabilities.MCPJSONSchema{
					Type:       "object",
					Properties: map[string]string{"query": "string"},
				},
				OutputSchema: capabilities.MCPJSONSchema{Type: "object"},
			},
		},
	})
	if _, err := provider.Initialize(context.Background()); err != nil {
		t.Fatal(err)
	}

	capabilitiesList, err := provider.ListCapabilities()
	if err != nil {
		t.Fatal(err)
	}
	capability := capabilitiesList[0]
	for _, forbidden := range []string{"token", "local-secret", "raw_payload"} {
		if strings.Contains(capability.DisplayName, forbidden) || strings.Contains(capability.Description, forbidden) {
			t.Fatalf("unsafe mcp metadata leaked %q: %+v", forbidden, capability)
		}
	}
}

func TestMCPSafetyConvertsStructuredContentThroughSafetyGate(t *testing.T) {
	provider := newInitializedMockMCPProviderForTest(t)

	candidate, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID: "mcp.mock.tool.asset_lookup",
		Arguments:    map[string]any{"query": "asset"},
	})
	if err != nil {
		t.Fatal(err)
	}

	result, err := product.NewStructuredResultSafetyGate().Approve(candidate)
	if err != nil {
		t.Fatal(err)
	}
	if result.SafeSummary != "MCP 资产查询完成" || result.ResultRef == "" {
		t.Fatalf("structured result mismatch: %+v", result)
	}
}

func TestMCPSafetyRejectsErrorAndUnsafeStructuredContent(t *testing.T) {
	provider := newInitializedMockMCPProviderForTest(t)

	_, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID: "mcp.mock.tool.error_lookup",
		Arguments:    map[string]any{"query": "asset"},
	})
	if !errors.Is(err, capabilities.ErrMCPToolExecutionFailed) {
		t.Fatalf("isError err = %v", err)
	}

	provider.SetToolResultForTest("asset_lookup", capabilities.MCPToolResult{
		StructuredContent: map[string]any{
			"result_ref":   "result:mcp:asset_lookup",
			"safe_summary": "raw_payload={}",
		},
	})
	_, err = provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID: "mcp.mock.tool.asset_lookup",
		Arguments:    map[string]any{"query": "asset"},
	})
	if !errors.Is(err, facts.ErrUnsafeFactMaterial) {
		t.Fatalf("unsafe structured content err = %v", err)
	}
}

func TestMCPRiskPolicyMergesAnnotationsWithoutTrustingThem(t *testing.T) {
	provider := capabilities.NewMockMCPProvider(capabilities.MCPMockConfig{
		ServerID: "mock",
		Tools: []capabilities.MCPToolDefinition{
			{
				Name:        "write_claims_read_only",
				Title:       "伪装写域",
				Description: "MCP annotation 声明只读，但项目策略必须覆盖",
				InputSchema: capabilities.MCPJSONSchema{
					Type:       "object",
					Properties: map[string]string{"ticket": "string"},
				},
				OutputSchema: capabilities.MCPJSONSchema{Type: "object"},
				Annotations: capabilities.MCPToolAnnotations{
					ReadOnlyHint:    true,
					DestructiveHint: true,
				},
				ProjectPolicy: capabilities.MCPProjectPolicy{
					RiskLevel:           capabilities.RiskHigh,
					SideEffect:          capabilities.SideEffectWriteExternal,
					ApprovalRequired:    true,
					IdempotencyRequired: true,
				},
			},
		},
		Results: map[string]capabilities.MCPToolResult{
			"write_claims_read_only": {
				StructuredContent: map[string]any{
					"result_ref":   "result:mcp:write_claims_read_only",
					"safe_summary": "写域待审批",
				},
			},
		},
	})

	capabilitiesList, err := provider.ListCapabilities()
	if err != nil {
		t.Fatal(err)
	}
	if len(capabilitiesList) != 1 {
		t.Fatalf("capabilities len = %d", len(capabilitiesList))
	}
	capability := capabilitiesList[0]
	if capability.RiskLevel != capabilities.RiskHigh || capability.SideEffect != capabilities.SideEffectWriteExternal || !capability.ApprovalRequired {
		t.Fatalf("project policy must override untrusted MCP annotations: %+v", capability)
	}
}

func TestMCPRiskPolicyRequiresIdempotencyForDestructiveAnnotation(t *testing.T) {
	provider := capabilities.NewMockMCPProvider(capabilities.MCPMockConfig{
		ServerID: "mock",
		Tools: []capabilities.MCPToolDefinition{
			{
				Name:        "destructive_write",
				Title:       "破坏性写域",
				Description: "MCP annotation 标记 destructive 时必须提升幂等要求",
				InputSchema: capabilities.MCPJSONSchema{
					Type:       "object",
					Properties: map[string]string{"ticket": "string"},
				},
				OutputSchema: capabilities.MCPJSONSchema{Type: "object"},
				Annotations:  capabilities.MCPToolAnnotations{DestructiveHint: true},
			},
		},
	})

	capabilitiesList, err := provider.ListCapabilities()
	if err != nil {
		t.Fatal(err)
	}
	capability := capabilitiesList[0]
	if capability.SideEffect != capabilities.SideEffectWriteExternal || !capability.ApprovalRequired || !capability.IdempotencyRequired {
		t.Fatalf("destructive annotation must require approval and idempotency: %+v", capability)
	}
}

func newMockMCPProviderForTest(t *testing.T) *capabilities.MockMCPProvider {
	t.Helper()

	provider := capabilities.NewMockMCPProvider(capabilities.MCPMockConfig{
		ServerID: "mock",
		Tools: []capabilities.MCPToolDefinition{
			{
				Name:        "asset_lookup",
				Title:       "MCP 资产查询",
				Description: "通过 mock MCP 查询资产",
				InputSchema: capabilities.MCPJSONSchema{
					Type:       "object",
					Properties: map[string]string{"query": "string"},
					Required:   []string{"query"},
				},
				OutputSchema: capabilities.MCPJSONSchema{Type: "object"},
				Annotations:  capabilities.MCPToolAnnotations{ReadOnlyHint: true},
			},
			{
				Name:        "error_lookup",
				Title:       "MCP 错误查询",
				Description: "模拟 MCP tool execution error",
				InputSchema: capabilities.MCPJSONSchema{
					Type:       "object",
					Properties: map[string]string{"query": "string"},
					Required:   []string{"query"},
				},
				OutputSchema: capabilities.MCPJSONSchema{Type: "object"},
				Annotations:  capabilities.MCPToolAnnotations{ReadOnlyHint: true},
			},
		},
		Results: map[string]capabilities.MCPToolResult{
			"asset_lookup": {
				StructuredContent: map[string]any{
					"result_ref":   "result:mcp:asset_lookup",
					"safe_summary": "MCP 资产查询完成",
				},
			},
			"error_lookup": {
				IsError: true,
				Content: []capabilities.MCPContent{
					{Type: "text", Text: "safe execution error"},
				},
			},
		},
	})
	return provider
}

func newInitializedMockMCPProviderForTest(t *testing.T) *capabilities.MockMCPProvider {
	t.Helper()
	provider := newMockMCPProviderForTest(t)
	if _, err := provider.Initialize(context.Background()); err != nil {
		t.Fatal(err)
	}
	return provider
}
