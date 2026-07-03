package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/bootstrap"
	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/llm"
	"agent-platform-eino/internal/einoapp/providers/fobrain"
)

type mainMissingCheckpointResolver struct{}

func (mainMissingCheckpointResolver) ResolveCheckpointID(context.Context, string, string, string) (string, bool, error) {
	return "", false, nil
}

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

func TestServiceCommandsInjectCheckpointResolver(t *testing.T) {
	// 启动装配必须启用 checkpoint resolver，不能只在单元测试路径保护 resume。
	repository := facts.NewMemoryRepository()
	ctx := context.Background()
	if err := repository.CreateRun(ctx, facts.Run{RunID: "run-waiting", WorkspaceID: "ws-1", Status: facts.RunStatusWaiting}); err != nil {
		t.Fatal(err)
	}
	if err := repository.AppendPendingInteraction(ctx, facts.PendingInteraction{
		PendingID:     "pending-1",
		RunID:         "run-waiting",
		Kind:          facts.PendingKindApproval,
		Status:        facts.PendingStatusWaiting,
		ResumeRef:     "resume-safe-1",
		CheckpointRef: "checkpoint_ref:missing",
	}); err != nil {
		t.Fatal(err)
	}

	commands := serviceCommands(repository, nil, nil, capabilities.NewRegistry(), nil, mainMissingCheckpointResolver{})
	if _, err := commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-1",
		RunID:           "run-waiting",
		ResumeRef:       "resume-safe-1",
		ClientRequestID: "client-resume-1",
	}); !errors.Is(err, execution.ErrCheckpointMissing) {
		t.Fatalf("resume err = %v, want ErrCheckpointMissing", err)
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

func TestCapabilityRuntimeFromConfigRegistersFobrainOnlyWhenEnabled(t *testing.T) {
	// 默认启动路径不能内置 Fobrain；只有配置显式启用时才注册 provider。
	registry, _, err := capabilityRuntimeFromConfig(bootstrap.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := registry.Get(fobrain.CapabilityCurrentUserContext); ok {
		t.Fatalf("fobrain capability should not be registered without config: %+v", registry.List())
	}

	cfg := bootstrap.Config{
		Fobrain: bootstrap.FobrainConfig{
			Enabled:     true,
			ConnectorID: "fobrain",
			WorkspaceID: "ws_fobrain",
			BaseURL:     "https://fobrain.example.test/api",
			Timeout:     9 * time.Second,
			Credential: bootstrap.FobrainCredentialConfig{
				Status:     "bound",
				DisplayRef: "bound:fobrain:local",
				OwnerScope: "workspace",
				AuthParam:  "authorization",
				APIToken:   "local-secret",
			},
			ConnectorStatus: bootstrap.FobrainConnectorStatusConfig{Mode: "mock", Available: true},
			CredentialBinding: bootstrap.CredentialBinding{
				SchemaVersion: "eino.provider_credential_binding.v1",
				WorkspaceID:   "ws_fobrain",
				System:        "fobrain",
				Status:        "bound",
				DisplayRef:    "bound:fobrain:local",
				OwnerScope:    "workspace",
			},
		},
	}
	policyContexts := policyContextsFromConfig(cfg)
	if policyContexts[fobrain.CapabilityCurrentUserContext].CredentialBinding.Status != capabilities.CredentialStatusBound {
		t.Fatalf("fobrain policy context missing bound credential: %+v", policyContexts)
	}

	registry, invoker, err := capabilityRuntimeFromConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}

	capability, ok := registry.Get(fobrain.CapabilityCurrentUserContext)
	if !ok {
		t.Fatalf("fobrain capability was not registered: %+v", registry.List())
	}
	if capability.ProviderID != fobrain.ProviderID || capability.CredentialBindingPolicy != capabilities.CredentialBindingRequired {
		t.Fatalf("fobrain capability metadata mismatch: %+v", capability)
	}

	candidate, err := invoker.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID: fobrain.CapabilityCurrentUserContext,
		PolicyContext: capabilities.PolicyContext{
			WorkspaceID:       "ws_fobrain",
			CredentialBinding: fobrainCredentialBindingFromConfig(cfg.Fobrain.CredentialBinding),
			ConnectorStatus:   capabilities.ConnectorStatusAvailable,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if candidate.SchemaVersion != facts.StructuredResultSchemaVersion || candidate.SafeSummary == "" {
		t.Fatalf("fobrain candidate mismatch: %+v", candidate)
	}

	batchDCapability, ok := registry.Get(fobrain.CapabilityListAssetsByOwner)
	if !ok {
		t.Fatalf("mock fobrain Batch D capability was not registered: %+v", registry.List())
	}
	batchDPolicyContext, ok := policyContexts[fobrain.CapabilityListAssetsByOwner]
	if !ok || batchDPolicyContext.CredentialBinding.Status != capabilities.CredentialStatusBound {
		t.Fatalf("mock fobrain Batch D policy context missing: %+v", policyContexts)
	}
	batchDCandidate, err := invoker.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID:  fobrain.CapabilityListAssetsByOwner,
		Arguments:     map[string]any{"person_name": "张三"},
		PolicyContext: batchDPolicyContext,
	})
	if err != nil {
		t.Fatalf("mock fobrain Batch D invoke failed for %+v: %v", batchDCapability, err)
	}
	if batchDCandidate.SchemaVersion != facts.StructuredResultSchemaVersion || !strings.Contains(batchDCandidate.SafeSummary, "张三") {
		t.Fatalf("mock fobrain Batch D candidate mismatch: %+v", batchDCandidate)
	}
}

func TestPolicyContextsFromConfigDerivesFobrainCapabilitiesFromProviderCatalog(t *testing.T) {
	// policy context 不能维护第二份 Fobrain capability id 清单，必须跟 provider catalog 同步。
	cfg := fobrainEnabledTestConfig("https://fobrain.example.test/api", "mock")
	contexts := policyContextsFromConfig(bootstrap.Config{Fobrain: cfg})

	provider, err := fobrainProviderFromConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := provider.ListCapabilities()
	if err != nil {
		t.Fatal(err)
	}
	for _, capability := range catalog {
		context, ok := contexts[capability.ID]
		if !ok {
			t.Fatalf("missing policy context for %s: %+v", capability.ID, contexts)
		}
		if context.WorkspaceID != "ws_fobrain" || context.CredentialBinding.Status != capabilities.CredentialStatusBound {
			t.Fatalf("policy context mismatch for %s: %+v", capability.ID, context)
		}
	}
}

func TestCapabilityRuntimeFromConfigUsesFobrainHTTPClientInLiveMode(t *testing.T) {
	// live mode 必须走 provider HTTP client；token 仍只在 provider 边界内作为 header 使用。
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("authorization")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"display_name":    "王五",
				"department_name": "安全部",
				"role":            map[string]any{"name": "安全运营"},
			},
		})
	}))
	defer server.Close()

	cfg := bootstrap.Config{
		Fobrain: bootstrap.FobrainConfig{
			Enabled:     true,
			ConnectorID: "fobrain",
			WorkspaceID: "ws_fobrain",
			BaseURL:     server.URL,
			Timeout:     time.Second,
			Credential: bootstrap.FobrainCredentialConfig{
				Status:     "bound",
				DisplayRef: "bound:fobrain:local",
				OwnerScope: "workspace",
				AuthParam:  "authorization",
				APIToken:   "workspace-token",
			},
			ConnectorStatus: bootstrap.FobrainConnectorStatusConfig{Mode: "live", Available: true},
			CredentialBinding: bootstrap.CredentialBinding{
				SchemaVersion: "eino.provider_credential_binding.v1",
				WorkspaceID:   "ws_fobrain",
				System:        "fobrain",
				Status:        "bound",
				DisplayRef:    "bound:fobrain:local",
				OwnerScope:    "workspace",
			},
		},
	}

	_, invoker, err := capabilityRuntimeFromConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	candidate, err := invoker.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID: fobrain.CapabilityCurrentUserContext,
		PolicyContext: capabilities.PolicyContext{
			WorkspaceID:       "ws_fobrain",
			CredentialBinding: fobrainCredentialBindingFromConfig(cfg.Fobrain.CredentialBinding),
			ConnectorStatus:   capabilities.ConnectorStatusAvailable,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotAuth != "workspace-token" {
		t.Fatalf("authorization header = %q", gotAuth)
	}
	if !strings.Contains(candidate.SafeSummary, "王五") || !strings.Contains(candidate.SafeSummary, "安全运营") {
		t.Fatalf("live fobrain candidate did not use HTTP result: %+v", candidate)
	}
}

func TestCapabilityRuntimeFromConfigPassesFobrainTLSInsecureSkipVerify(t *testing.T) {
	// 私有证书跳过校验必须来自配置文件，并只传给 Fobrain HTTP client。
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"display_name": "王五"},
		})
	}))
	defer server.Close()

	cfg := fobrainEnabledTestConfig(server.URL, "live")
	cfg.TLS.InsecureSkipVerify = true
	_, invoker, err := capabilityRuntimeFromConfig(bootstrap.Config{Fobrain: cfg})
	if err != nil {
		t.Fatal(err)
	}
	candidate, err := invoker.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID: fobrain.CapabilityCurrentUserContext,
		PolicyContext: capabilities.PolicyContext{
			WorkspaceID:       "ws_fobrain",
			CredentialBinding: fobrainCredentialBindingFromConfig(cfg.CredentialBinding),
			ConnectorStatus:   capabilities.ConnectorStatusAvailable,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(candidate.SafeSummary, "王五") {
		t.Fatalf("candidate did not use TLS server result: %+v", candidate)
	}
}

func fobrainEnabledTestConfig(baseURL string, mode string) bootstrap.FobrainConfig {
	return bootstrap.FobrainConfig{
		Enabled:     true,
		ConnectorID: "fobrain",
		WorkspaceID: "ws_fobrain",
		BaseURL:     baseURL,
		Timeout:     time.Second,
		Credential: bootstrap.FobrainCredentialConfig{
			Status:     "bound",
			DisplayRef: "bound:fobrain:local",
			OwnerScope: "workspace",
			AuthParam:  "authorization",
			APIToken:   "workspace-token",
		},
		ConnectorStatus: bootstrap.FobrainConnectorStatusConfig{Mode: mode, Available: true},
		CredentialBinding: bootstrap.CredentialBinding{
			SchemaVersion: "eino.provider_credential_binding.v1",
			WorkspaceID:   "ws_fobrain",
			System:        "fobrain",
			Status:        "bound",
			DisplayRef:    "bound:fobrain:local",
			OwnerScope:    "workspace",
		},
	}
}

func TestCapabilityRuntimeFromConfigBlocksLiveFobrainWorkspaceMismatchBeforeHTTP(t *testing.T) {
	// workspace scope 必须在 provider 边界先拦截，不能让错工作区请求触达真实 Fobrain。
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requestCount++
		_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"display_name": "王五"}})
	}))
	defer server.Close()

	cfg := bootstrap.Config{
		Fobrain: bootstrap.FobrainConfig{
			Enabled:     true,
			ConnectorID: "fobrain",
			WorkspaceID: "ws_fobrain",
			BaseURL:     server.URL,
			Timeout:     time.Second,
			Credential: bootstrap.FobrainCredentialConfig{
				Status:     "bound",
				DisplayRef: "bound:fobrain:local",
				OwnerScope: "workspace",
				AuthParam:  "authorization",
				APIToken:   "workspace-token",
			},
			ConnectorStatus: bootstrap.FobrainConnectorStatusConfig{Mode: "live", Available: true},
			CredentialBinding: bootstrap.CredentialBinding{
				SchemaVersion: "eino.provider_credential_binding.v1",
				WorkspaceID:   "ws_fobrain",
				System:        "fobrain",
				Status:        "bound",
				DisplayRef:    "bound:fobrain:local",
				OwnerScope:    "workspace",
			},
		},
	}

	_, invoker, err := capabilityRuntimeFromConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	_, err = invoker.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID: fobrain.CapabilityCurrentUserContext,
		PolicyContext: capabilities.PolicyContext{
			WorkspaceID:       "ws_other",
			CredentialBinding: fobrainCredentialBindingFromConfig(cfg.Fobrain.CredentialBinding),
			ConnectorStatus:   capabilities.ConnectorStatusAvailable,
		},
	})
	if !fobrain.HasReason(err, capabilities.PolicyReasonCredentialScopeDenied) {
		t.Fatalf("err = %v, want credential scope denied", err)
	}
	if requestCount != 0 {
		t.Fatalf("live HTTP requests = %d, want blocked before HTTP", requestCount)
	}
}
