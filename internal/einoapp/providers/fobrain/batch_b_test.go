package fobrain_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
	"agent-platform-eino/internal/einoapp/providers/fobrain"
)

func TestProviderCatalogRegistersBatchBMyScopeTools(t *testing.T) {
	provider := fobrain.NewProvider(fobrain.ProviderConfig{Client: fobrain.MockClient{}})

	catalog, err := provider.ListCapabilities()
	if err != nil {
		t.Fatal(err)
	}
	byID := map[string]capabilities.Capability{}
	for _, capability := range catalog {
		byID[capability.ID] = capability
	}
	// Batch B 工具必须复用当前用户/本部门上下文，不要求用户输入负责人或部门名。
	expectedSchemas := map[string]struct {
		required   []string
		properties map[string]string
	}{
		fobrain.CapabilityMyAssets: {
			required:   []string{},
			properties: map[string]string{"keyword": "string", "page": "integer", "page_size": "integer"},
		},
		fobrain.CapabilityMyDepartmentAssets: {
			required:   []string{},
			properties: map[string]string{"keyword": "string", "page": "integer", "page_size": "integer"},
		},
		fobrain.CapabilityMyVulnerabilities: {
			required:   []string{},
			properties: map[string]string{"keyword": "string", "page": "integer", "page_size": "integer"},
		},
		fobrain.CapabilityMyDepartmentVulnerabilities: {
			required:   []string{},
			properties: map[string]string{"keyword": "string", "page": "integer", "page_size": "integer"},
		},
		fobrain.CapabilityMyBusinessSystems: {
			required:   []string{},
			properties: map[string]string{"keyword": "string", "page": "integer", "page_size": "integer"},
		},
		fobrain.CapabilityMyImportantBusinessSystems: {
			required:   []string{},
			properties: map[string]string{"keyword": "string", "page": "integer", "page_size": "integer"},
		},
	}
	for id, expected := range expectedSchemas {
		capability, ok := byID[id]
		if !ok {
			t.Fatalf("missing Batch B capability %s in %+v", id, catalog)
		}
		if capability.ResultSchema != fobrain.BusinessResultSchemaVersion ||
			capability.PolicyRef != fobrain.PolicyRead ||
			capability.ConnectorID != fobrain.ConnectorID ||
			capability.RiskLevel != capabilities.RiskLow ||
			capability.SideEffect != capabilities.SideEffectReadExternal ||
			capability.ApprovalRequired {
			t.Fatalf("Batch B capability metadata mismatch: %+v", capability)
		}
		if !reflect.DeepEqual(capability.InputSchema.Required, expected.required) {
			t.Fatalf("Batch B required fields mismatch for %s: got %+v want %+v", id, capability.InputSchema.Required, expected.required)
		}
		if !reflect.DeepEqual(capability.InputSchema.Properties, expected.properties) {
			t.Fatalf("Batch B properties mismatch for %s: got %+v want %+v", id, capability.InputSchema.Properties, expected.properties)
		}
	}
}

func TestProviderInvokesBatchBMyAssetsThroughStructuredResult(t *testing.T) {
	client := &recordingBatchBClient{
		result: fobrain.MyScopeResult{
			ToolID:     fobrain.CapabilityMyAssets,
			Title:      "查询我的资产",
			EntityType: fobrain.QueryEntityAsset,
			Query:      fobrain.MyScopeQuery{Keyword: "生产", Page: 1, PageSize: 20},
			Items:      []fobrain.QueryResultItem{{EntityRef: "asset:fobrain:fixture-1", DisplayName: "prod-web-01", Status: "online"}},
		},
	}
	provider := fobrain.NewProvider(fobrain.ProviderConfig{
		WorkspaceID:        "ws_fobrain",
		CredentialBinding:  boundFobrainCredential("ws_fobrain"),
		ConnectorStatus:    capabilities.ConnectorStatusAvailable,
		CredentialResolver: fobrain.StaticCredentialResolver{WorkspaceID: "ws_fobrain", AuthParam: "authorization", APIToken: "local-secret"},
		Client:             client,
	})

	candidate, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID: fobrain.CapabilityMyAssets,
		Arguments: map[string]any{
			"keyword":   "生产",
			"page":      1,
			"page_size": 20,
		},
		PolicyContext: policyContextForFobrain("ws_fobrain"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if client.calls != 1 || client.lastToolID != fobrain.CapabilityMyAssets || client.lastQuery.Keyword != "生产" {
		t.Fatalf("Batch B client call mismatch: calls=%d tool=%s query=%+v", client.calls, client.lastToolID, client.lastQuery)
	}
	if candidate.SchemaVersion != facts.StructuredResultSchemaVersion ||
		candidate.ResultRef != "result:fobrain:my-assets" ||
		candidate.ItemCount != 1 ||
		!strings.Contains(candidate.SafeSummary, "查询我的资产") ||
		!strings.Contains(candidate.SafeSummary, "生产") {
		t.Fatalf("Batch B candidate mismatch: %+v", candidate)
	}
	if strings.Contains(strings.ToLower(candidate.SafeSummary), "authorization") ||
		strings.Contains(strings.ToLower(candidate.SafeSummary), "raw") {
		t.Fatalf("Batch B summary leaked unsafe material: %s", candidate.SafeSummary)
	}
	if _, err := product.NewStructuredResultSafetyGate().Approve(candidate); err != nil {
		t.Fatalf("Batch B candidate rejected: %v", err)
	}
}

func TestProviderInvokesBatchBDepartmentAndImportantScopes(t *testing.T) {
	tests := []struct {
		name         string
		capabilityID string
		wantRef      string
		wantSummary  string
		assertQuery  func(*testing.T, fobrain.MyScopeQuery)
	}{
		{
			name:         "department assets",
			capabilityID: fobrain.CapabilityMyDepartmentAssets,
			wantRef:      "result:fobrain:my-department-assets",
			wantSummary:  "本部门范围",
			assertQuery: func(t *testing.T, query fobrain.MyScopeQuery) {
				t.Helper()
				if !query.Department {
					t.Fatalf("department query not marked: %+v", query)
				}
			},
		},
		{
			name:         "important business",
			capabilityID: fobrain.CapabilityMyImportantBusinessSystems,
			wantRef:      "result:fobrain:my-important-business-systems",
			wantSummary:  "当前用户重要业务范围",
			assertQuery: func(t *testing.T, query fobrain.MyScopeQuery) {
				t.Helper()
				if !query.ImportantOnly {
					t.Fatalf("important query not marked: %+v", query)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &recordingBatchBClient{result: fobrain.MyScopeResult{ToolID: tt.capabilityID}}
			provider := fobrain.NewProvider(fobrain.ProviderConfig{
				WorkspaceID:        "ws_fobrain",
				CredentialBinding:  boundFobrainCredential("ws_fobrain"),
				ConnectorStatus:    capabilities.ConnectorStatusAvailable,
				CredentialResolver: fobrain.StaticCredentialResolver{WorkspaceID: "ws_fobrain", AuthParam: "authorization", APIToken: "local-secret"},
				Client:             client,
			})

			candidate, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
				CapabilityID:  tt.capabilityID,
				Arguments:     map[string]any{},
				PolicyContext: policyContextForFobrain("ws_fobrain"),
			})
			if err != nil {
				t.Fatal(err)
			}
			tt.assertQuery(t, client.lastQuery)
			if candidate.ResultRef != tt.wantRef || !strings.Contains(candidate.SafeSummary, tt.wantSummary) {
				t.Fatalf("Batch B candidate mismatch: %+v", candidate)
			}
		})
	}
}

func TestProviderUsesValidatedBatchBQueryForStructuredResult(t *testing.T) {
	client := &recordingBatchBClient{
		result: fobrain.MyScopeResult{
			ToolID: fobrain.CapabilityMyAssets,
			Title:  "Authorization: Bearer raw",
			Query:  fobrain.MyScopeQuery{Keyword: "Authorization: Bearer raw"},
			Items:  []fobrain.QueryResultItem{{EntityRef: "asset:fobrain:fixture-1", DisplayName: "raw"}},
		},
	}
	provider := fobrain.NewProvider(fobrain.ProviderConfig{
		WorkspaceID:        "ws_fobrain",
		CredentialBinding:  boundFobrainCredential("ws_fobrain"),
		ConnectorStatus:    capabilities.ConnectorStatusAvailable,
		CredentialResolver: fobrain.StaticCredentialResolver{WorkspaceID: "ws_fobrain", AuthParam: "authorization", APIToken: "local-secret"},
		Client:             client,
	})

	candidate, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID:  fobrain.CapabilityMyAssets,
		Arguments:     map[string]any{"keyword": "业务"},
		PolicyContext: policyContextForFobrain("ws_fobrain"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(candidate.SafeSummary, "业务") {
		t.Fatalf("safe summary must use validated request query: %+v", candidate)
	}
	if strings.Contains(strings.ToLower(candidate.SafeSummary), "authorization") ||
		strings.Contains(strings.ToLower(candidate.SafeSummary), "raw") {
		t.Fatalf("safe summary leaked client-returned query or title: %s", candidate.SafeSummary)
	}
}

func TestProviderRejectsBatchBUnsafeKeyword(t *testing.T) {
	client := &recordingBatchBClient{}
	provider := fobrain.NewProvider(fobrain.ProviderConfig{
		WorkspaceID:        "ws_fobrain",
		CredentialBinding:  boundFobrainCredential("ws_fobrain"),
		ConnectorStatus:    capabilities.ConnectorStatusAvailable,
		CredentialResolver: fobrain.StaticCredentialResolver{WorkspaceID: "ws_fobrain", AuthParam: "authorization", APIToken: "local-secret"},
		Client:             client,
	})

	_, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID:  fobrain.CapabilityMyAssets,
		Arguments:     map[string]any{"keyword": "Authorization: Bearer secret"},
		PolicyContext: policyContextForFobrain("ws_fobrain"),
	})
	if !errors.Is(err, facts.ErrUnsafeFactMaterial) {
		t.Fatalf("unsafe keyword err = %v", err)
	}
	if client.calls != 0 {
		t.Fatalf("client must not be called for unsafe arguments: %d", client.calls)
	}
}

type recordingBatchBClient struct {
	recordingBatchAClient
	calls      int
	lastToolID string
	lastQuery  fobrain.MyScopeQuery
	result     fobrain.MyScopeResult
	err        error
}

func (client *recordingBatchBClient) MyScopeQuery(_ context.Context, _ fobrain.ResolvedCredential, toolID string, query fobrain.MyScopeQuery) (fobrain.MyScopeResult, error) {
	client.calls++
	client.lastToolID = toolID
	client.lastQuery = query
	if client.err != nil {
		return fobrain.MyScopeResult{}, client.err
	}
	if client.result.ToolID == "" {
		client.result.ToolID = toolID
		client.result.Query = query
	}
	return client.result, nil
}
