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

func TestProviderCatalogRegistersBatchCDirectReadTools(t *testing.T) {
	provider := fobrain.NewProvider(fobrain.ProviderConfig{Client: fobrain.MockClient{}})

	catalog, err := provider.ListCapabilities()
	if err != nil {
		t.Fatal(err)
	}
	byID := map[string]capabilities.Capability{}
	for _, capability := range catalog {
		byID[capability.ID] = capability
	}
	// Batch C 覆盖无筛选首读、统计指标和待处理工单，只允许安全筛选字段。
	expectedSchemas := map[string]struct {
		required   []string
		properties map[string]string
	}{
		fobrain.CapabilityBusinessList: {
			required:   []string{},
			properties: map[string]string{"business_name": "string", "owner": "string", "keyword": "string", "time_range": "string", "page": "integer", "page_size": "integer"},
		},
		fobrain.CapabilityExternalHighRiskAssets: {
			required:   []string{},
			properties: map[string]string{"field": "string", "severity": "string", "time_range": "string"},
		},
		fobrain.CapabilityVulnerabilityStatusSummary: {
			required:   []string{},
			properties: map[string]string{"field": "string", "severity": "string", "time_range": "string"},
		},
		fobrain.CapabilityPendingTickets: {
			required:   []string{},
			properties: map[string]string{"person": "string", "status": "string", "page": "integer", "page_size": "integer"},
		},
		fobrain.CapabilityIPStats: {
			required:   []string{},
			properties: map[string]string{"field": "string", "severity": "string", "time_range": "string"},
		},
		fobrain.CapabilityVulStats: {
			required:   []string{},
			properties: map[string]string{"field": "string", "severity": "string", "time_range": "string"},
		},
	}
	for id, expected := range expectedSchemas {
		capability, ok := byID[id]
		if !ok {
			t.Fatalf("missing Batch C capability %s in %+v", id, catalog)
		}
		if capability.ResultSchema != fobrain.BusinessResultSchemaVersion ||
			capability.PolicyRef != fobrain.PolicyRead ||
			capability.ConnectorID != fobrain.ConnectorID ||
			capability.RiskLevel != capabilities.RiskLow ||
			capability.SideEffect != capabilities.SideEffectReadExternal ||
			capability.ApprovalRequired {
			t.Fatalf("Batch C capability metadata mismatch: %+v", capability)
		}
		if !reflect.DeepEqual(capability.InputSchema.Required, expected.required) {
			t.Fatalf("Batch C required fields mismatch for %s: got %+v want %+v", id, capability.InputSchema.Required, expected.required)
		}
		if !reflect.DeepEqual(capability.InputSchema.Properties, expected.properties) {
			t.Fatalf("Batch C properties mismatch for %s: got %+v want %+v", id, capability.InputSchema.Properties, expected.properties)
		}
		// Batch C 支持无筛选首读，Eino ToolInfo 不能把可选筛选字段误标为必填。
		info, err := capabilities.ToolInfoFromCapability(context.Background(), capability)
		if err != nil {
			t.Fatalf("Batch C tool info mismatch for %s: %v", id, err)
		}
		jsonSchema, err := info.ParamsOneOf.ToJSONSchema()
		if err != nil {
			t.Fatalf("Batch C tool json schema mismatch for %s: %v", id, err)
		}
		if len(jsonSchema.Required) != 0 {
			t.Fatalf("Batch C required fields must be optional for %s: %+v", id, jsonSchema.Required)
		}
	}
}

func TestProviderInvokesBatchCBusinessListThroughStructuredResult(t *testing.T) {
	client := &recordingBatchCClient{
		result: fobrain.DirectReadResult{
			ToolID:     fobrain.CapabilityBusinessList,
			Title:      "查询业务系统",
			EntityType: "business_system",
			Query:      fobrain.DirectReadQuery{Keyword: "支付", Page: 1, PageSize: 20},
			Items:      []fobrain.QueryResultItem{{EntityRef: "business_system:fobrain:fixture-1", DisplayName: "支付系统", Status: "resolved"}},
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
		CapabilityID: fobrain.CapabilityBusinessList,
		Arguments: map[string]any{
			"keyword":   "支付",
			"page":      1,
			"page_size": 20,
		},
		PolicyContext: policyContextForFobrain("ws_fobrain"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if client.calls != 1 || client.lastToolID != fobrain.CapabilityBusinessList || client.lastQuery.Keyword != "支付" {
		t.Fatalf("Batch C client call mismatch: calls=%d tool=%s query=%+v", client.calls, client.lastToolID, client.lastQuery)
	}
	if candidate.SchemaVersion != facts.StructuredResultSchemaVersion ||
		candidate.ResultRef != "result:fobrain:business-list" ||
		candidate.ItemCount != 1 ||
		!strings.Contains(candidate.SafeSummary, "查询业务系统") ||
		!strings.Contains(candidate.SafeSummary, "支付") {
		t.Fatalf("Batch C candidate mismatch: %+v", candidate)
	}
	if _, err := product.NewStructuredResultSafetyGate().Approve(candidate); err != nil {
		t.Fatalf("Batch C candidate rejected: %v", err)
	}
}

func TestProviderInvokesBatchCMetricsAndTickets(t *testing.T) {
	tests := []struct {
		name         string
		capabilityID string
		arguments    map[string]any
		result       fobrain.DirectReadResult
		wantRef      string
		wantSummary  string
		wantCount    int
	}{
		{
			name:         "ip stats",
			capabilityID: fobrain.CapabilityIPStats,
			arguments:    map[string]any{"field": "external", "time_range": "7d"},
			result:       fobrain.DirectReadResult{ToolID: fobrain.CapabilityIPStats, Metrics: []fobrain.RiskMetric{{Label: "外网 IP", Count: 3}}},
			wantRef:      "result:fobrain:ip-stats",
			wantSummary:  "统计 IP 资产",
			wantCount:    1,
		},
		{
			name:         "pending tickets",
			capabilityID: fobrain.CapabilityPendingTickets,
			arguments:    map[string]any{"person": "当前用户", "status": "pending"},
			result:       fobrain.DirectReadResult{ToolID: fobrain.CapabilityPendingTickets, Items: []fobrain.QueryResultItem{{EntityRef: "ticket:fobrain:fixture-1", DisplayName: "待处理工单"}}},
			wantRef:      "result:fobrain:pending-tickets",
			wantSummary:  "查询待处理工单",
			wantCount:    1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &recordingBatchCClient{result: tt.result}
			provider := fobrain.NewProvider(fobrain.ProviderConfig{
				WorkspaceID:        "ws_fobrain",
				CredentialBinding:  boundFobrainCredential("ws_fobrain"),
				ConnectorStatus:    capabilities.ConnectorStatusAvailable,
				CredentialResolver: fobrain.StaticCredentialResolver{WorkspaceID: "ws_fobrain", AuthParam: "authorization", APIToken: "local-secret"},
				Client:             client,
			})

			candidate, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
				CapabilityID:  tt.capabilityID,
				Arguments:     tt.arguments,
				PolicyContext: policyContextForFobrain("ws_fobrain"),
			})
			if err != nil {
				t.Fatal(err)
			}
			if candidate.ResultRef != tt.wantRef || candidate.ItemCount != tt.wantCount || !strings.Contains(candidate.SafeSummary, tt.wantSummary) {
				t.Fatalf("Batch C candidate mismatch: %+v", candidate)
			}
			if _, err := product.NewStructuredResultSafetyGate().Approve(candidate); err != nil {
				t.Fatalf("Batch C candidate rejected: %v", err)
			}
		})
	}
}

func TestProviderUsesValidatedBatchCQueryForStructuredResult(t *testing.T) {
	client := &recordingBatchCClient{
		result: fobrain.DirectReadResult{
			ToolID: fobrain.CapabilityBusinessList,
			Title:  "Authorization: Bearer raw",
			Query:  fobrain.DirectReadQuery{Keyword: "Authorization: Bearer raw"},
			Items:  []fobrain.QueryResultItem{{EntityRef: "business_system:fobrain:fixture-1", DisplayName: "raw"}},
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
		CapabilityID:  fobrain.CapabilityBusinessList,
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

func TestProviderRejectsBatchCUnsafeFilter(t *testing.T) {
	client := &recordingBatchCClient{}
	provider := fobrain.NewProvider(fobrain.ProviderConfig{
		WorkspaceID:        "ws_fobrain",
		CredentialBinding:  boundFobrainCredential("ws_fobrain"),
		ConnectorStatus:    capabilities.ConnectorStatusAvailable,
		CredentialResolver: fobrain.StaticCredentialResolver{WorkspaceID: "ws_fobrain", AuthParam: "authorization", APIToken: "local-secret"},
		Client:             client,
	})

	_, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID:  fobrain.CapabilityBusinessList,
		Arguments:     map[string]any{"keyword": "Authorization: Bearer secret"},
		PolicyContext: policyContextForFobrain("ws_fobrain"),
	})
	if !errors.Is(err, facts.ErrUnsafeFactMaterial) {
		t.Fatalf("unsafe filter err = %v", err)
	}
	if client.calls != 0 {
		t.Fatalf("client must not be called for unsafe arguments: %d", client.calls)
	}
}

type recordingBatchCClient struct {
	recordingBatchAClient
	calls      int
	lastToolID string
	lastQuery  fobrain.DirectReadQuery
	result     fobrain.DirectReadResult
	err        error
}

func (client *recordingBatchCClient) DirectRead(_ context.Context, _ fobrain.ResolvedCredential, toolID string, query fobrain.DirectReadQuery) (fobrain.DirectReadResult, error) {
	client.calls++
	client.lastToolID = toolID
	client.lastQuery = query
	if client.err != nil {
		return fobrain.DirectReadResult{}, client.err
	}
	if client.result.ToolID == "" {
		client.result.ToolID = toolID
		client.result.Query = query
	}
	return client.result, nil
}
