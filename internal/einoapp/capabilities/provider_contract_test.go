package capabilities_test

import (
	"context"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/facts"
)

func TestEinoToolAdapterBuildsToolInfoFromCapabilityMetadata(t *testing.T) {
	// ToolInfo 必须由注册元数据生成，避免按旧工具名或自然语言关键词硬编码。
	capability := capabilities.Capability{
		ID:          "cap.mock.asset.read",
		ProviderID:  "mock",
		ToolName:    "mock_asset_read",
		DisplayName: "资产查询",
		Description: "按授权范围读取资产",
		InputSchema: capabilities.JSONSchema{
			SchemaVersion: "json_schema.v1",
			Properties: map[string]string{
				"query": "string",
			},
		},
		ResultSchema: facts.StructuredResultSchemaVersion,
		RiskLevel:    capabilities.RiskReadOnly,
		Timeout:      5 * time.Second,
	}

	info, err := capabilities.ToolInfoFromCapability(context.Background(), capability)
	if err != nil {
		t.Fatal(err)
	}

	if info.Name != "mock_asset_read" || info.Desc != "按授权范围读取资产" {
		t.Fatalf("tool info not generated from metadata: %+v", info)
	}
	if info.ParamsOneOf == nil {
		t.Fatal("tool info params must be generated from capability input schema")
	}
}

func TestProviderContractConvertsReadResultToStructuredResultCandidate(t *testing.T) {
	provider := capabilities.NewMockProvider("mock", []capabilities.Capability{
		{
			ID:           "cap.mock.asset.read",
			ProviderID:   "mock",
			ToolName:     "mock_asset_read",
			DisplayName:  "资产查询",
			Description:  "按授权范围读取资产",
			ResultSchema: facts.StructuredResultSchemaVersion,
			RiskLevel:    capabilities.RiskReadOnly,
		},
	})

	candidate, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID: "cap.mock.asset.read",
		Arguments:    map[string]any{"query": "asset"},
	})
	if err != nil {
		t.Fatal(err)
	}

	if candidate.SchemaVersion != facts.StructuredResultSchemaVersion {
		t.Fatalf("schema version = %q", candidate.SchemaVersion)
	}
	if candidate.ResultRef == "" || candidate.SafeSummary == "" {
		t.Fatalf("candidate must contain safe ref and summary: %+v", candidate)
	}
}

func TestProviderContractRejectsUnknownCapability(t *testing.T) {
	provider := capabilities.NewMockProvider("mock", nil)

	_, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID: "cap.missing",
	})
	if err == nil {
		t.Fatal("missing capability err = nil")
	}
}
