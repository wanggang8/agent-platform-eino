package execution_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
)

func TestAgentToolLoopWritesToolFactsFromMockReadCapability(t *testing.T) {
	// Phase 4 tool loop 先用 mock read capability 证明 Eino tool -> Safety Gate -> Product Facts。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	now := time.Unix(700, 0).UTC()
	if err := repository.CreateRun(ctx, facts.Run{RunID: "run-tool", WorkspaceID: "ws-tool", Status: facts.RunStatusCreated, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatal(err)
	}
	registry := capabilities.NewRegistry()
	capability := mockReadCapability()
	if err := registry.Register(capability); err != nil {
		t.Fatal(err)
	}
	provider := capabilities.NewMockProvider("mock", []capabilities.Capability{capability})
	runner := execution.NewToolLoopRunner(repository, registry, provider, execution.ToolLoopRunnerConfig{
		Now: func() time.Time { return now.Add(time.Second) },
	})

	if err := runner.RunCapability(ctx, "run-tool", capability.ID, "asset query"); err != nil {
		t.Fatal(err)
	}

	snapshot, err := repository.GetSnapshot(ctx, "run-tool")
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Run.Status != facts.RunStatusSucceeded {
		t.Fatalf("run status = %q", snapshot.Run.Status)
	}
	if len(snapshot.ToolCalls) != 1 || snapshot.ToolCalls[0].ToolID != capability.ID || snapshot.ToolCalls[0].DisplayName != "资产查询" {
		t.Fatalf("tool call not written from registry metadata: %+v", snapshot.ToolCalls)
	}
	if snapshot.ToolCalls[0].ArgsPreview != "keyword=present" {
		t.Fatalf("args preview = %q, want schema-driven field", snapshot.ToolCalls[0].ArgsPreview)
	}
	if len(snapshot.ToolResults) != 1 || snapshot.ToolResults[0].StructuredResult.SchemaVersion != facts.StructuredResultSchemaVersion {
		t.Fatalf("structured result not written: %+v", snapshot.ToolResults)
	}
	if len(snapshot.AuditEvents) == 0 {
		t.Fatal("tool loop audit not written")
	}

	projection := product.NewFactsProjection(repository)
	action, err := projection.ActionResult(ctx, "ws-tool", "action-tool", "run-tool")
	if err != nil {
		t.Fatal(err)
	}
	if len(action.ResultCards) != 1 || action.ResultCards[0].SafeSummary == "" {
		t.Fatalf("action result card not projected from facts: %+v", action.ResultCards)
	}
}

func TestToolSafetyRejectsUnsafeArgumentsBeforeFacts(t *testing.T) {
	repository := facts.NewMemoryRepository()
	ctx := context.Background()
	now := time.Unix(750, 0).UTC()
	if err := repository.CreateRun(ctx, facts.Run{RunID: "run-tool", WorkspaceID: "ws-tool", Status: facts.RunStatusCreated, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatal(err)
	}
	registry := capabilities.NewRegistry()
	capability := mockReadCapability()
	if err := registry.Register(capability); err != nil {
		t.Fatal(err)
	}
	runner := execution.NewToolLoopRunner(repository, registry, capabilities.NewMockProvider("mock", []capabilities.Capability{capability}), execution.ToolLoopRunnerConfig{
		Now: func() time.Time { return now.Add(time.Second) },
	})

	err := runner.RunCapability(ctx, "run-tool", capability.ID, "Authorization: Bearer local-secret")
	if !errors.Is(err, facts.ErrUnsafeFactMaterial) {
		t.Fatalf("unsafe arguments err = %v", err)
	}

	snapshot, err := repository.GetSnapshot(ctx, "run-tool")
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.ToolCalls) != 0 || len(snapshot.ToolResults) != 0 {
		t.Fatalf("unsafe tool facts must not be written: calls=%+v results=%+v", snapshot.ToolCalls, snapshot.ToolResults)
	}
}

func TestToolSafetyRejectsUnsafeArgumentKeysAndNestedValues(t *testing.T) {
	repository := facts.NewMemoryRepository()
	ctx := context.Background()
	now := time.Unix(775, 0).UTC()
	if err := repository.CreateRun(ctx, facts.Run{RunID: "run-tool", WorkspaceID: "ws-tool", Status: facts.RunStatusCreated, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatal(err)
	}
	capability := mockReadCapability()
	tool := execution.NewEinoCapabilityTool(repository, capability, capabilities.NewMockProvider("mock", []capabilities.Capability{capability}), execution.EinoCapabilityToolConfig{
		RunID: "run-tool",
		Now:   func() time.Time { return now.Add(time.Second) },
	})

	for _, arguments := range []string{
		`{"credential_ref":"safe-looking"}`,
		`{"keyword":{"nested":"raw provider body"}}`,
		`{"keyword":["ok","token=secret"]}`,
	} {
		t.Run(arguments, func(t *testing.T) {
			_, err := tool.InvokableRun(ctx, arguments)
			if !errors.Is(err, facts.ErrUnsafeFactMaterial) {
				t.Fatalf("unsafe nested arguments err = %v", err)
			}
		})
	}
}

func TestToolLoopMapsActionTextToRequiredInputField(t *testing.T) {
	repository := facts.NewMemoryRepository()
	ctx := context.Background()
	now := time.Unix(790, 0).UTC()
	if err := repository.CreateRun(ctx, facts.Run{RunID: "run-tool", WorkspaceID: "ws-tool", Status: facts.RunStatusCreated, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatal(err)
	}
	registry := capabilities.NewRegistry()
	capability := mockReadCapability()
	capability.InputSchema = capabilities.JSONSchema{
		SchemaVersion: "json_schema.v1",
		Properties: map[string]string{
			"asset_type": "string",
			"query":      "string",
		},
		Required: []string{"query"},
	}
	if err := registry.Register(capability); err != nil {
		t.Fatal(err)
	}
	provider := capabilities.NewMockProvider("mock", []capabilities.Capability{capability})
	runner := execution.NewToolLoopRunner(repository, registry, provider, execution.ToolLoopRunnerConfig{
		Now: func() time.Time { return now.Add(time.Second) },
	})

	if err := runner.RunCapability(ctx, "run-tool", capability.ID, "asset query"); err != nil {
		t.Fatal(err)
	}

	snapshot, err := repository.GetSnapshot(ctx, "run-tool")
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.ToolCalls) != 1 || snapshot.ToolCalls[0].ArgsPreview != "query=present" {
		t.Fatalf("required field was not used for action text: %+v", snapshot.ToolCalls)
	}
}

func TestToolLoopRejectsApprovalRequiredCapabilityBeforeHITL(t *testing.T) {
	repository := facts.NewMemoryRepository()
	registry := capabilities.NewRegistry()
	capability := mockReadCapability()
	capability.ID = "cap.mock.write"
	capability.ToolName = "mock_write"
	capability.RiskLevel = capabilities.RiskWrite
	if err := registry.Register(capability); err != nil {
		t.Fatal(err)
	}
	runner := execution.NewToolLoopRunner(repository, registry, capabilities.NewMockProvider("mock", []capabilities.Capability{capability}), execution.ToolLoopRunnerConfig{})

	err := runner.RunCapability(context.Background(), "run-tool", capability.ID, "write")
	if !errors.Is(err, execution.ErrCapabilityRequiresApproval) {
		t.Fatalf("approval required err = %v", err)
	}
}

func TestToolLoopUsesDefaultPolicyContextForCredentialRequiredCapability(t *testing.T) {
	// 凭据型能力必须使用配置派生的安全 policy context 才能进入工具执行。
	repository := facts.NewMemoryRepository()
	ctx := context.Background()
	now := time.Unix(795, 0).UTC()
	if err := repository.CreateRun(ctx, facts.Run{RunID: "run-tool", WorkspaceID: "ws-tool", Status: facts.RunStatusCreated, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatal(err)
	}
	registry := capabilities.NewRegistry()
	capability := mockReadCapability()
	capability.ID = "tool.fobrain.current_user_context"
	capability.ProviderID = "fobrain"
	capability.CredentialBindingPolicy = capabilities.CredentialBindingRequired
	capability.ConnectorID = "fobrain"
	if err := registry.Register(capability); err != nil {
		t.Fatal(err)
	}
	provider := capabilities.NewMockProvider("fobrain", []capabilities.Capability{capability})
	runner := execution.NewToolLoopRunner(repository, registry, provider, execution.ToolLoopRunnerConfig{
		Now: func() time.Time { return now.Add(time.Second) },
		PolicyContexts: map[string]capabilities.PolicyContext{
			capability.ID: {
				WorkspaceID: "ws-tool",
				CredentialBinding: capabilities.CredentialBinding{
					WorkspaceID: "ws-tool",
					System:      "fobrain",
					Status:      capabilities.CredentialStatusBound,
					DisplayRef:  "bound:fobrain:local",
					OwnerScope:  capabilities.PermissionScopeWorkspace,
				},
				ConnectorStatus: capabilities.ConnectorStatusAvailable,
			},
		},
	})

	if err := runner.RunCapability(ctx, "run-tool", capability.ID, "current user"); err != nil {
		t.Fatal(err)
	}
}

func TestToolLoopRejectsDefaultPolicyContextForDifferentRunWorkspace(t *testing.T) {
	// run workspace 必须覆盖配置 workspace，避免复用其他工作区的凭据绑定。
	repository := facts.NewMemoryRepository()
	ctx := context.Background()
	now := time.Unix(796, 0).UTC()
	if err := repository.CreateRun(ctx, facts.Run{RunID: "run-tool", WorkspaceID: "ws-other", Status: facts.RunStatusCreated, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatal(err)
	}
	registry := capabilities.NewRegistry()
	capability := mockReadCapability()
	capability.ID = "tool.fobrain.current_user_context"
	capability.ProviderID = "fobrain"
	capability.CredentialBindingPolicy = capabilities.CredentialBindingRequired
	capability.ConnectorID = "fobrain"
	if err := registry.Register(capability); err != nil {
		t.Fatal(err)
	}
	provider := capabilities.NewMockProvider("fobrain", []capabilities.Capability{capability})
	runner := execution.NewToolLoopRunner(repository, registry, provider, execution.ToolLoopRunnerConfig{
		Now: func() time.Time { return now.Add(time.Second) },
		PolicyContexts: map[string]capabilities.PolicyContext{
			capability.ID: {
				WorkspaceID: "ws-tool",
				CredentialBinding: capabilities.CredentialBinding{
					WorkspaceID: "ws-tool",
					System:      "fobrain",
					Status:      capabilities.CredentialStatusBound,
					DisplayRef:  "bound:fobrain:local",
					OwnerScope:  capabilities.PermissionScopeWorkspace,
				},
				ConnectorStatus: capabilities.ConnectorStatusAvailable,
			},
		},
	})

	err := runner.RunCapability(ctx, "run-tool", capability.ID, "current user")
	if !errors.Is(err, execution.ErrCapabilityRequiresApproval) {
		t.Fatalf("cross-workspace policy err = %v", err)
	}
}

func TestToolLoopApprovedCapabilityStillEnforcesProviderPolicy(t *testing.T) {
	// 审批通过只允许越过 approval_required；缺凭据、跨 workspace 和 connector 不可用仍不能执行 provider。
	for _, testCase := range []struct {
		name          string
		runWorkspace  string
		policyContext capabilities.PolicyContext
	}{
		{
			name:         "missing credential",
			runWorkspace: "ws-tool",
			policyContext: capabilities.PolicyContext{
				WorkspaceID:     "ws-tool",
				ConnectorStatus: capabilities.ConnectorStatusAvailable,
				CredentialBinding: capabilities.CredentialBinding{
					WorkspaceID: "ws-tool",
					System:      "fobrain",
					Status:      capabilities.CredentialStatusMissing,
					DisplayRef:  "missing",
					OwnerScope:  capabilities.PermissionScopeWorkspace,
				},
			},
		},
		{
			name:         "credential scope denied",
			runWorkspace: "ws-other",
			policyContext: capabilities.PolicyContext{
				WorkspaceID:     "ws-tool",
				ConnectorStatus: capabilities.ConnectorStatusAvailable,
				CredentialBinding: capabilities.CredentialBinding{
					WorkspaceID: "ws-tool",
					System:      "fobrain",
					Status:      capabilities.CredentialStatusBound,
					DisplayRef:  "bound:fobrain:local",
					OwnerScope:  capabilities.PermissionScopeWorkspace,
				},
			},
		},
		{
			name:         "connector unavailable",
			runWorkspace: "ws-tool",
			policyContext: capabilities.PolicyContext{
				WorkspaceID:     "ws-tool",
				ConnectorStatus: capabilities.ConnectorStatusUnavailable,
				CredentialBinding: capabilities.CredentialBinding{
					WorkspaceID: "ws-tool",
					System:      "fobrain",
					Status:      capabilities.CredentialStatusBound,
					DisplayRef:  "bound:fobrain:local",
					OwnerScope:  capabilities.PermissionScopeWorkspace,
				},
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			repository := facts.NewMemoryRepository()
			ctx := context.Background()
			now := time.Unix(797, 0).UTC()
			if err := repository.CreateRun(ctx, facts.Run{RunID: "run-tool", WorkspaceID: testCase.runWorkspace, Status: facts.RunStatusCreated, CreatedAt: now, UpdatedAt: now}); err != nil {
				t.Fatal(err)
			}
			registry := capabilities.NewRegistry()
			capability := mockWriteCredentialCapability()
			if err := registry.Register(capability); err != nil {
				t.Fatal(err)
			}
			runner := execution.NewToolLoopRunner(repository, registry, capabilities.NewMockProvider("fobrain", []capabilities.Capability{capability}), execution.ToolLoopRunnerConfig{
				Now: func() time.Time { return now.Add(time.Second) },
				PolicyContexts: map[string]capabilities.PolicyContext{
					capability.ID: testCase.policyContext,
				},
			})

			err := runner.RunApprovedCapability(ctx, "run-tool", capability.ID, "write")
			if !errors.Is(err, execution.ErrCapabilityRequiresApproval) {
				t.Fatalf("approved policy err = %v, want ErrCapabilityRequiresApproval", err)
			}
			snapshot, err := repository.GetSnapshot(ctx, "run-tool")
			if err != nil {
				t.Fatal(err)
			}
			if len(snapshot.ToolCalls) != 0 || len(snapshot.ToolResults) != 0 {
				t.Fatalf("blocked approved capability must not write tool facts: calls=%+v results=%+v", snapshot.ToolCalls, snapshot.ToolResults)
			}
		})
	}
}

func TestToolLoopApprovedCapabilityRunsAfterPolicyAllows(t *testing.T) {
	// 已审批且 provider policy 允许时，写域 capability 才能进入工具执行。
	repository := facts.NewMemoryRepository()
	ctx := context.Background()
	now := time.Unix(798, 0).UTC()
	if err := repository.CreateRun(ctx, facts.Run{RunID: "run-tool", WorkspaceID: "ws-tool", Status: facts.RunStatusCreated, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatal(err)
	}
	registry := capabilities.NewRegistry()
	capability := mockWriteCredentialCapability()
	if err := registry.Register(capability); err != nil {
		t.Fatal(err)
	}
	runner := execution.NewToolLoopRunner(repository, registry, writeInvoker{}, execution.ToolLoopRunnerConfig{
		Now: func() time.Time { return now.Add(time.Second) },
		PolicyContexts: map[string]capabilities.PolicyContext{
			capability.ID: {
				WorkspaceID:     "ws-tool",
				ConnectorStatus: capabilities.ConnectorStatusAvailable,
				CredentialBinding: capabilities.CredentialBinding{
					WorkspaceID: "ws-tool",
					System:      "fobrain",
					Status:      capabilities.CredentialStatusBound,
					DisplayRef:  "bound:fobrain:local",
					OwnerScope:  capabilities.PermissionScopeWorkspace,
				},
			},
		},
	})

	if err := runner.RunApprovedCapability(ctx, "run-tool", capability.ID, "write"); err != nil {
		t.Fatal(err)
	}
	snapshot, err := repository.GetSnapshot(ctx, "run-tool")
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Run.Status != facts.RunStatusSucceeded || len(snapshot.ToolCalls) != 1 {
		t.Fatalf("approved capability did not execute after policy allow: %+v", snapshot)
	}
}

func mockReadCapability() capabilities.Capability {
	return capabilities.Capability{
		ID:          "cap.mock.asset.read",
		ProviderID:  "mock",
		ToolName:    "mock_asset_read",
		DisplayName: "资产查询",
		Description: "按授权范围读取资产",
		InputSchema: capabilities.JSONSchema{
			SchemaVersion: "json_schema.v1",
			Properties: map[string]string{
				"keyword": "string",
			},
		},
		ResultSchema: facts.StructuredResultSchemaVersion,
		RiskLevel:    capabilities.RiskReadOnly,
		Timeout:      5 * time.Second,
	}
}

func mockWriteCredentialCapability() capabilities.Capability {
	capability := mockReadCapability()
	capability.ID = "tool.fobrain.ticket.write"
	capability.ProviderID = "fobrain"
	capability.ToolName = "fobrain_ticket_write"
	capability.DisplayName = "更新工单"
	capability.RiskLevel = capabilities.RiskWrite
	capability.SideEffect = capabilities.SideEffectWriteExternal
	capability.CredentialBindingPolicy = capabilities.CredentialBindingRequired
	capability.ConnectorID = "fobrain"
	capability.ApprovalRequired = true
	return capability
}

type writeInvoker struct{}

func (writeInvoker) Invoke(_ context.Context, request capabilities.InvocationRequest) (product.StructuredResultCandidate, error) {
	return product.StructuredResultCandidate{
		SchemaVersion: facts.StructuredResultSchemaVersion,
		ResultRef:     "result:" + request.CapabilityID,
		SafeSummary:   "写入完成",
	}, nil
}
