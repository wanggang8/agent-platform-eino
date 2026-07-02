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

func TestProviderCatalogRegistersBatchEDetailRiskTools(t *testing.T) {
	provider := fobrain.NewProvider(fobrain.ProviderConfig{Client: fobrain.MockClient{}})

	catalog, err := provider.ListCapabilities()
	if err != nil {
		t.Fatal(err)
	}
	byID := map[string]capabilities.Capability{}
	for _, capability := range catalog {
		byID[capability.ID] = capability
	}
	// Batch E 是详情与风险关联读取能力，字段必须和 docs/fobrain-batch-e-interface-plan.md 保持一致。
	expectedSchemas := map[string]struct {
		required   []string
		properties map[string]string
	}{
		fobrain.CapabilityGetAssetDetail: {
			required:   []string{"asset_id"},
			properties: map[string]string{"asset_id": "string", "network_type": "string"},
		},
		fobrain.CapabilityGetVulnerabilityDetail: {
			required:   []string{"vulnerability_id"},
			properties: map[string]string{"vulnerability_id": "string"},
		},
		fobrain.CapabilityBusinessRiskSummary: {
			required:   []string{"business_name"},
			properties: map[string]string{"business_name": "string"},
		},
		fobrain.CapabilityThreatRelevanceList: {
			required:   []string{"vulnerability_name"},
			properties: map[string]string{"vulnerability_name": "string", "ip": "string", "business_name": "string", "page": "integer", "page_size": "integer"},
		},
	}
	for id, expected := range expectedSchemas {
		capability, ok := byID[id]
		if !ok {
			t.Fatalf("missing Batch E capability %s in %+v", id, catalog)
		}
		if capability.ResultSchema != fobrain.BusinessResultSchemaVersion ||
			capability.PolicyRef != fobrain.PolicyRead ||
			capability.ConnectorID != fobrain.ConnectorID ||
			capability.RiskLevel != capabilities.RiskLow ||
			capability.SideEffect != capabilities.SideEffectReadExternal ||
			capability.ApprovalRequired {
			t.Fatalf("Batch E capability metadata mismatch: %+v", capability)
		}
		if !reflect.DeepEqual(capability.InputSchema.Required, expected.required) {
			t.Fatalf("Batch E required fields mismatch for %s: got %+v want %+v", id, capability.InputSchema.Required, expected.required)
		}
		if !reflect.DeepEqual(capability.InputSchema.Properties, expected.properties) {
			t.Fatalf("Batch E properties mismatch for %s: got %+v want %+v", id, capability.InputSchema.Properties, expected.properties)
		}
	}
}

func TestProviderInvokesBatchEAssetDetailThroughStructuredResult(t *testing.T) {
	client := &recordingBatchEClient{
		assetResult: fobrain.DetailRiskResult{
			ToolID: fobrain.CapabilityGetAssetDetail,
			Items:  []fobrain.QueryResultItem{{EntityRef: "asset:fobrain:asset-1", DisplayName: "prod-web-01", Status: "online"}},
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
		CapabilityID: fobrain.CapabilityGetAssetDetail,
		Arguments: map[string]any{
			"asset_id":     "asset-1",
			"network_type": "内网",
		},
		PolicyContext: policyContextForFobrain("ws_fobrain"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if client.assetCalls != 1 || client.lastAssetQuery.AssetID != "asset-1" || client.lastAssetQuery.NetworkType != "internal" {
		t.Fatalf("asset detail call mismatch calls=%d query=%+v", client.assetCalls, client.lastAssetQuery)
	}
	if candidate.SchemaVersion != facts.StructuredResultSchemaVersion ||
		candidate.ResultRef != "result:fobrain:get-asset-detail" ||
		candidate.ItemCount != 1 ||
		!strings.Contains(candidate.SafeSummary, "资产详情") ||
		!strings.Contains(candidate.SafeSummary, "asset-1") {
		t.Fatalf("asset detail candidate mismatch: %+v", candidate)
	}
	if _, err := product.NewStructuredResultSafetyGate().Approve(candidate); err != nil {
		t.Fatalf("asset detail candidate rejected: %v", err)
	}
}

func TestProviderInvokesBatchEThreatRelevanceWithValidatedArguments(t *testing.T) {
	client := &recordingBatchEClient{
		relevanceResult: fobrain.DetailRiskResult{
			ToolID: fobrain.CapabilityThreatRelevanceList,
			Items:  []fobrain.QueryResultItem{{EntityRef: "vuln:fobrain:relevance-1", DisplayName: "高危组件漏洞", Affected: 2}},
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
		CapabilityID: fobrain.CapabilityThreatRelevanceList,
		Arguments: map[string]any{
			"vulnerability_name": "高危组件漏洞",
			"ip":                 "10.10.11.12",
			"business_name":      "核心业务",
			"page":               1,
			"page_size":          20,
		},
		PolicyContext: policyContextForFobrain("ws_fobrain"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if client.relevanceCalls != 1 ||
		client.lastRelevanceQuery.VulnerabilityName != "高危组件漏洞" ||
		client.lastRelevanceQuery.IP != "10.10.11.12" ||
		client.lastRelevanceQuery.BusinessName != "核心业务" ||
		client.lastRelevanceQuery.Page != 1 ||
		client.lastRelevanceQuery.PageSize != 20 {
		t.Fatalf("threat relevance call mismatch calls=%d query=%+v", client.relevanceCalls, client.lastRelevanceQuery)
	}
	if candidate.ResultRef != "result:fobrain:threat-relevance-list" ||
		candidate.ItemCount != 1 ||
		!strings.Contains(candidate.SafeSummary, "威胁关联") ||
		!strings.Contains(candidate.SafeSummary, "高危组件漏洞") {
		t.Fatalf("threat relevance candidate mismatch: %+v", candidate)
	}
}

func TestProviderInvokesBatchEVulnerabilityDetailAndBusinessRisk(t *testing.T) {
	tests := []struct {
		name         string
		capabilityID string
		arguments    map[string]any
		configure    func(*recordingBatchEClient)
		assertClient func(*testing.T, *recordingBatchEClient)
		wantRef      string
		wantSummary  string
		wantCount    int
	}{
		{
			name:         "vulnerability detail",
			capabilityID: fobrain.CapabilityGetVulnerabilityDetail,
			arguments:    map[string]any{"vulnerability_id": "vuln-1"},
			configure: func(client *recordingBatchEClient) {
				client.vulnerabilityResult = fobrain.DetailRiskResult{
					ToolID: fobrain.CapabilityGetVulnerabilityDetail,
					Items:  []fobrain.QueryResultItem{{EntityRef: "vuln:fobrain:vuln-1", DisplayName: "高危组件漏洞"}},
				}
			},
			assertClient: func(t *testing.T, client *recordingBatchEClient) {
				t.Helper()
				if client.vulnerabilityCalls != 1 || client.lastVulnQuery.VulnerabilityID != "vuln-1" {
					t.Fatalf("vulnerability call mismatch calls=%d query=%+v", client.vulnerabilityCalls, client.lastVulnQuery)
				}
			},
			wantRef:     "result:fobrain:get-vulnerability-detail",
			wantSummary: "漏洞详情",
			wantCount:   1,
		},
		{
			name:         "business risk",
			capabilityID: fobrain.CapabilityBusinessRiskSummary,
			arguments:    map[string]any{"business_name": "核心业务"},
			configure: func(client *recordingBatchEClient) {
				client.businessResult = fobrain.DetailRiskResult{
					ToolID:  fobrain.CapabilityBusinessRiskSummary,
					Metrics: []fobrain.RiskMetric{{Label: "核心业务", Count: 7}},
				}
			},
			assertClient: func(t *testing.T, client *recordingBatchEClient) {
				t.Helper()
				if client.businessCalls != 1 || client.lastBusinessQuery.BusinessName != "核心业务" {
					t.Fatalf("business call mismatch calls=%d query=%+v", client.businessCalls, client.lastBusinessQuery)
				}
			},
			wantRef:     "result:fobrain:business-risk-summary",
			wantSummary: "业务风险摘要",
			wantCount:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &recordingBatchEClient{}
			tt.configure(client)
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
			tt.assertClient(t, client)
			if candidate.ResultRef != tt.wantRef || candidate.ItemCount != tt.wantCount || !strings.Contains(candidate.SafeSummary, tt.wantSummary) {
				t.Fatalf("Batch E candidate mismatch: %+v", candidate)
			}
			if _, err := product.NewStructuredResultSafetyGate().Approve(candidate); err != nil {
				t.Fatalf("Batch E candidate rejected: %v", err)
			}
		})
	}
}

func TestProviderBuildsBatchEEmptyStructuredResult(t *testing.T) {
	client := &recordingBatchEClient{businessResult: fobrain.DetailRiskResult{ToolID: fobrain.CapabilityBusinessRiskSummary}}
	provider := fobrain.NewProvider(fobrain.ProviderConfig{
		WorkspaceID:        "ws_fobrain",
		CredentialBinding:  boundFobrainCredential("ws_fobrain"),
		ConnectorStatus:    capabilities.ConnectorStatusAvailable,
		CredentialResolver: fobrain.StaticCredentialResolver{WorkspaceID: "ws_fobrain", AuthParam: "authorization", APIToken: "local-secret"},
		Client:             client,
	})

	candidate, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID:  fobrain.CapabilityBusinessRiskSummary,
		Arguments:     map[string]any{"business_name": "空业务"},
		PolicyContext: policyContextForFobrain("ws_fobrain"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if candidate.ItemCount != 0 || !strings.Contains(candidate.SafeSummary, "返回 0 条") {
		t.Fatalf("empty Batch E candidate mismatch: %+v", candidate)
	}
}

func TestProviderRejectsBatchEMissingRequiredArgument(t *testing.T) {
	client := &recordingBatchEClient{}
	provider := fobrain.NewProvider(fobrain.ProviderConfig{
		WorkspaceID:        "ws_fobrain",
		CredentialBinding:  boundFobrainCredential("ws_fobrain"),
		ConnectorStatus:    capabilities.ConnectorStatusAvailable,
		CredentialResolver: fobrain.StaticCredentialResolver{WorkspaceID: "ws_fobrain", AuthParam: "authorization", APIToken: "local-secret"},
		Client:             client,
	})

	_, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID:  fobrain.CapabilityGetVulnerabilityDetail,
		Arguments:     map[string]any{"page": 1},
		PolicyContext: policyContextForFobrain("ws_fobrain"),
	})
	if err == nil || !fobrain.HasReason(err, capabilities.PolicyReasonConnectorExecutionFailed) {
		t.Fatalf("missing argument err = %v", err)
	}
	if client.vulnerabilityCalls != 0 {
		t.Fatalf("client must not be called for invalid arguments: %d", client.vulnerabilityCalls)
	}
}

func TestProviderRejectsBatchENonStringArgumentBeforeClient(t *testing.T) {
	client := &recordingBatchEClient{}
	provider := fobrain.NewProvider(fobrain.ProviderConfig{
		WorkspaceID:        "ws_fobrain",
		CredentialBinding:  boundFobrainCredential("ws_fobrain"),
		ConnectorStatus:    capabilities.ConnectorStatusAvailable,
		CredentialResolver: fobrain.StaticCredentialResolver{WorkspaceID: "ws_fobrain", AuthParam: "authorization", APIToken: "local-secret"},
		Client:             client,
	})

	_, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID:  fobrain.CapabilityBusinessRiskSummary,
		Arguments:     map[string]any{"business_name": 123},
		PolicyContext: policyContextForFobrain("ws_fobrain"),
	})
	if err == nil || !fobrain.HasReason(err, capabilities.PolicyReasonConnectorExecutionFailed) {
		t.Fatalf("non-string argument err = %v", err)
	}
	if client.businessCalls != 0 {
		t.Fatalf("client must not be called for schema mismatch: %d", client.businessCalls)
	}
}

func TestProviderRejectsBatchEUnsafeArgumentBeforeClient(t *testing.T) {
	client := &recordingBatchEClient{}
	provider := fobrain.NewProvider(fobrain.ProviderConfig{
		WorkspaceID:        "ws_fobrain",
		CredentialBinding:  boundFobrainCredential("ws_fobrain"),
		ConnectorStatus:    capabilities.ConnectorStatusAvailable,
		CredentialResolver: fobrain.StaticCredentialResolver{WorkspaceID: "ws_fobrain", AuthParam: "authorization", APIToken: "local-secret"},
		Client:             client,
	})

	_, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID:  fobrain.CapabilityGetAssetDetail,
		Arguments:     map[string]any{"asset_id": "Authorization: Bearer secret"},
		PolicyContext: policyContextForFobrain("ws_fobrain"),
	})
	if !errors.Is(err, facts.ErrUnsafeFactMaterial) {
		t.Fatalf("unsafe argument err = %v", err)
	}
	if client.assetCalls != 0 {
		t.Fatalf("client must not be called for unsafe argument: %d", client.assetCalls)
	}
}

func TestProviderFoldsBatchEProviderErrorWithoutLeak(t *testing.T) {
	client := &recordingBatchEClient{err: errors.New("workspace-token Authorization raw payload")}
	provider := fobrain.NewProvider(fobrain.ProviderConfig{
		WorkspaceID:        "ws_fobrain",
		CredentialBinding:  boundFobrainCredential("ws_fobrain"),
		ConnectorStatus:    capabilities.ConnectorStatusAvailable,
		CredentialResolver: fobrain.StaticCredentialResolver{WorkspaceID: "ws_fobrain", AuthParam: "authorization", APIToken: "local-secret"},
		Client:             client,
	})

	_, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID:  fobrain.CapabilityBusinessRiskSummary,
		Arguments:     map[string]any{"business_name": "核心业务"},
		PolicyContext: policyContextForFobrain("ws_fobrain"),
	})
	if !fobrain.HasReason(err, capabilities.PolicyReasonConnectorExecutionFailed) {
		t.Fatalf("provider error = %v, want connector execution failed", err)
	}
	for _, forbidden := range []string{"workspace-token", "Authorization", "raw payload"} {
		if strings.Contains(err.Error(), forbidden) {
			t.Fatalf("provider error leaked %q: %s", forbidden, err.Error())
		}
	}
}

func TestBatchEStructuredResultRejectsUnsafeMaterial(t *testing.T) {
	candidate, _ := fobrain.BuildDetailRiskStructuredResult(fobrain.DetailRiskResult{
		ToolID: fobrain.CapabilityBusinessRiskSummary,
		Target: "Authorization: Bearer secret",
	})

	_, err := product.NewStructuredResultSafetyGate().Approve(candidate)
	if !errors.Is(err, facts.ErrUnsafeFactMaterial) {
		t.Fatalf("unsafe Batch E candidate err = %v", err)
	}
}

type recordingBatchEClient struct {
	recordingBatchAClient
	assetCalls          int
	vulnerabilityCalls  int
	businessCalls       int
	relevanceCalls      int
	lastAssetQuery      fobrain.AssetDetailQuery
	lastVulnQuery       fobrain.VulnerabilityDetailQuery
	lastBusinessQuery   fobrain.BusinessRiskQuery
	lastRelevanceQuery  fobrain.ThreatRelevanceQuery
	assetResult         fobrain.DetailRiskResult
	vulnerabilityResult fobrain.DetailRiskResult
	businessResult      fobrain.DetailRiskResult
	relevanceResult     fobrain.DetailRiskResult
	err                 error
}

func (client *recordingBatchEClient) AssetDetail(_ context.Context, _ fobrain.ResolvedCredential, query fobrain.AssetDetailQuery) (fobrain.DetailRiskResult, error) {
	client.assetCalls++
	client.lastAssetQuery = query
	return client.assetResult, client.err
}

func (client *recordingBatchEClient) VulnerabilityDetail(_ context.Context, _ fobrain.ResolvedCredential, query fobrain.VulnerabilityDetailQuery) (fobrain.DetailRiskResult, error) {
	client.vulnerabilityCalls++
	client.lastVulnQuery = query
	return client.vulnerabilityResult, client.err
}

func (client *recordingBatchEClient) BusinessRiskSummary(_ context.Context, _ fobrain.ResolvedCredential, query fobrain.BusinessRiskQuery) (fobrain.DetailRiskResult, error) {
	client.businessCalls++
	client.lastBusinessQuery = query
	return client.businessResult, client.err
}

func (client *recordingBatchEClient) ThreatRelevanceList(_ context.Context, _ fobrain.ResolvedCredential, query fobrain.ThreatRelevanceQuery) (fobrain.DetailRiskResult, error) {
	client.relevanceCalls++
	client.lastRelevanceQuery = query
	return client.relevanceResult, client.err
}
