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

func TestProviderCatalogRegistersBatchDParameterizedQueries(t *testing.T) {
	provider := fobrain.NewProvider(fobrain.ProviderConfig{Client: fobrain.MockClient{}})

	catalog, err := provider.ListCapabilities()
	if err != nil {
		t.Fatal(err)
	}
	byID := map[string]capabilities.Capability{}
	for _, capability := range catalog {
		byID[capability.ID] = capability
	}
	// Batch D catalog 是后续真实 HTTP mapper 的契约入口，字段映射必须和矩阵保持稳定。
	expectedSchemas := map[string]struct {
		required   []string
		properties map[string]string
	}{
		fobrain.CapabilityListAssetsByOwner: {
			required:   []string{"person_name"},
			properties: map[string]string{"person_name": "string", "person_staff_id": "string", "page": "integer", "page_size": "integer", "severity": "string"},
		},
		fobrain.CapabilityListVulnerabilitiesByOwner: {
			required:   []string{"person_name"},
			properties: map[string]string{"person_name": "string", "person_staff_id": "string", "page": "integer", "page_size": "integer", "severity": "string"},
		},
		fobrain.CapabilityListAssetsByDepartment: {
			required:   []string{"department_name"},
			properties: map[string]string{"department_name": "string", "page": "integer", "page_size": "integer", "severity": "string"},
		},
		fobrain.CapabilityListVulnerabilitiesByDepartment: {
			required:   []string{"department_name"},
			properties: map[string]string{"department_name": "string", "page": "integer", "page_size": "integer", "severity": "string"},
		},
		fobrain.CapabilityListAssetsByIP: {
			required:   []string{"ip"},
			properties: map[string]string{"ip": "string", "page": "integer", "page_size": "integer", "severity": "string", "status": "string"},
		},
		fobrain.CapabilityListVulnerabilitiesByIP: {
			required:   []string{"ip"},
			properties: map[string]string{"ip": "string", "page": "integer", "page_size": "integer", "severity": "string", "status": "string"},
		},
	}
	if len(catalog) != 9 {
		t.Fatalf("mock parameterized catalog length = %d, want Batch A plus six Batch D capabilities: %+v", len(catalog), catalog)
	}
	batchDCount := 0
	for _, capability := range catalog {
		if _, ok := expectedSchemas[capability.ID]; ok {
			batchDCount++
		}
	}
	if batchDCount != len(expectedSchemas) {
		t.Fatalf("Batch D catalog count = %d, want %d in %+v", batchDCount, len(expectedSchemas), catalog)
	}
	for id, expected := range expectedSchemas {
		capability, ok := byID[id]
		if !ok {
			t.Fatalf("missing Batch D capability %s in %+v", id, catalog)
		}
		if capability.ResultSchema != fobrain.BusinessResultSchemaVersion ||
			capability.PolicyRef != fobrain.PolicyRead ||
			capability.ConnectorID != fobrain.ConnectorID ||
			capability.RiskLevel != capabilities.RiskLow ||
			capability.SideEffect != capabilities.SideEffectReadExternal ||
			capability.ApprovalRequired {
			t.Fatalf("Batch D capability metadata mismatch: %+v", capability)
		}
		if !reflect.DeepEqual(capability.InputSchema.Required, expected.required) {
			t.Fatalf("Batch D required fields mismatch for %s: got %+v want %+v", id, capability.InputSchema.Required, expected.required)
		}
		if !reflect.DeepEqual(capability.InputSchema.Properties, expected.properties) {
			t.Fatalf("Batch D properties mismatch for %s: got %+v want %+v", id, capability.InputSchema.Properties, expected.properties)
		}
	}
}

func TestProviderInvokesBatchDParameterizedQueryThroughStructuredResult(t *testing.T) {
	client := &recordingBatchDClient{
		result: fobrain.ParameterizedQueryResult{
			ToolID:     fobrain.CapabilityListAssetsByOwner,
			Title:      "查询负责人资产",
			EntityType: fobrain.QueryEntityAsset,
			Query:      fobrain.ParameterizedQuery{PersonName: "张三", Page: 1, PageSize: 20},
			Items: []fobrain.QueryResultItem{
				{EntityRef: "asset:fobrain:fixture-1", DisplayName: "prod-web-01", OwnerName: "张三", Status: "online"},
			},
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
		CapabilityID: fobrain.CapabilityListAssetsByOwner,
		Arguments: map[string]any{
			"person_name": "张三",
			"page":        1,
			"page_size":   20,
		},
		PolicyContext: policyContextForFobrain("ws_fobrain"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if client.calls != 1 || client.lastToolID != fobrain.CapabilityListAssetsByOwner || client.lastQuery.PersonName != "张三" {
		t.Fatalf("Batch D client call mismatch: calls=%d tool=%s query=%+v", client.calls, client.lastToolID, client.lastQuery)
	}
	if candidate.SchemaVersion != facts.StructuredResultSchemaVersion ||
		candidate.ResultRef != "result:fobrain:list-assets-by-owner" ||
		candidate.ItemCount != 1 ||
		!strings.Contains(candidate.SafeSummary, "查询负责人资产") ||
		!strings.Contains(candidate.SafeSummary, "张三") {
		t.Fatalf("Batch D candidate mismatch: %+v", candidate)
	}
	if strings.Contains(strings.ToLower(candidate.SafeSummary), "authorization") ||
		strings.Contains(strings.ToLower(candidate.SafeSummary), "raw") {
		t.Fatalf("Batch D summary leaked unsafe material: %s", candidate.SafeSummary)
	}
	if _, err := product.NewStructuredResultSafetyGate().Approve(candidate); err != nil {
		t.Fatalf("Batch D candidate rejected: %v", err)
	}
}

func TestProviderUsesValidatedBatchDQueryForStructuredResult(t *testing.T) {
	client := &recordingBatchDClient{
		result: fobrain.ParameterizedQueryResult{
			ToolID: fobrain.CapabilityListAssetsByOwner,
			Title:  "Authorization: Bearer raw",
			Query:  fobrain.ParameterizedQuery{PersonName: "Authorization: Bearer raw"},
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
		CapabilityID:  fobrain.CapabilityListAssetsByOwner,
		Arguments:     map[string]any{"person_name": "李四"},
		PolicyContext: policyContextForFobrain("ws_fobrain"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(candidate.SafeSummary, "李四") {
		t.Fatalf("safe summary must use validated request query: %+v", candidate)
	}
	if strings.Contains(strings.ToLower(candidate.SafeSummary), "authorization") ||
		strings.Contains(strings.ToLower(candidate.SafeSummary), "raw") {
		t.Fatalf("safe summary leaked client-returned query or title: %s", candidate.SafeSummary)
	}
}

func TestProviderRejectsBatchDMissingRequiredArgument(t *testing.T) {
	client := &recordingBatchDClient{}
	provider := fobrain.NewProvider(fobrain.ProviderConfig{
		WorkspaceID:        "ws_fobrain",
		CredentialBinding:  boundFobrainCredential("ws_fobrain"),
		ConnectorStatus:    capabilities.ConnectorStatusAvailable,
		CredentialResolver: fobrain.StaticCredentialResolver{WorkspaceID: "ws_fobrain", AuthParam: "authorization", APIToken: "local-secret"},
		Client:             client,
	})

	_, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID:  fobrain.CapabilityListAssetsByOwner,
		Arguments:     map[string]any{"page": 1},
		PolicyContext: policyContextForFobrain("ws_fobrain"),
	})
	if err == nil || !fobrain.HasReason(err, capabilities.PolicyReasonConnectorExecutionFailed) {
		t.Fatalf("missing argument err = %v", err)
	}
	if client.calls != 0 {
		t.Fatalf("client must not be called for invalid arguments: %d", client.calls)
	}
}

func TestBatchDStructuredResultRejectsUnsafeMaterial(t *testing.T) {
	candidate, _ := fobrain.BuildParameterizedQueryStructuredResult(fobrain.ParameterizedQueryResult{
		ToolID:     fobrain.CapabilityListAssetsByIP,
		Title:      "查询 IP 资产",
		EntityType: fobrain.QueryEntityAsset,
		Query:      fobrain.ParameterizedQuery{IP: "Authorization: Bearer secret"},
	})

	_, err := product.NewStructuredResultSafetyGate().Approve(candidate)
	if !errors.Is(err, facts.ErrUnsafeFactMaterial) {
		t.Fatalf("unsafe Batch D candidate err = %v", err)
	}
}

type recordingBatchDClient struct {
	recordingBatchAClient
	calls      int
	lastToolID string
	lastQuery  fobrain.ParameterizedQuery
	result     fobrain.ParameterizedQueryResult
	err        error
}

func (client *recordingBatchDClient) ParameterizedQuery(_ context.Context, _ fobrain.ResolvedCredential, toolID string, query fobrain.ParameterizedQuery) (fobrain.ParameterizedQueryResult, error) {
	client.calls++
	client.lastToolID = toolID
	client.lastQuery = query
	if client.err != nil {
		return fobrain.ParameterizedQueryResult{}, client.err
	}
	if client.result.ToolID == "" {
		client.result.ToolID = toolID
		client.result.Query = query
	}
	return client.result, nil
}
