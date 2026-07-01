package execution_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/facts"
)

func TestStaticCommandsStartMessageUsesProvidedRunID(t *testing.T) {
	commands := execution.NewStaticCommands()

	accepted, err := commands.StartMessage(context.Background(), execution.MessageCommand{
		WorkspaceID: "ws_123",
		RunID:       "run_existing",
	})
	if err != nil {
		t.Fatal(err)
	}

	if accepted.RunID != "run_existing" {
		t.Fatalf("RunID = %q, want provided run id", accepted.RunID)
	}
	if accepted.Status != "accepted" {
		t.Fatalf("Status = %q, want accepted", accepted.Status)
	}
}

func TestStaticCommandsStartMessageCreatesOpaqueRunID(t *testing.T) {
	commands := execution.NewStaticCommands()

	accepted, err := commands.StartMessage(context.Background(), execution.MessageCommand{
		WorkspaceID: "ws_123",
	})
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(accepted.RunID, "run_") || strings.Contains(accepted.RunID, "ws_123") {
		t.Fatalf("RunID = %q", accepted.RunID)
	}
}

func TestStaticCommandsStartActionCreatesOpaqueRunID(t *testing.T) {
	commands := execution.NewStaticCommands()

	accepted, err := commands.StartAction(context.Background(), execution.ActionCommand{
		ActionID: "action-demo",
	})
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(accepted.RunID, "run_") || strings.Contains(accepted.RunID, "action-demo") {
		t.Fatalf("RunID = %q", accepted.RunID)
	}
}

func TestFactCommandsWriteAcceptedRunFact(t *testing.T) {
	// 命令层接受请求时必须先写 Product Facts，避免 Workbench 和 Action API 分裂。
	repository := facts.NewMemoryRepository()
	commands := execution.NewFactCommands(repository)

	accepted, err := commands.StartMessage(context.Background(), execution.MessageCommand{
		WorkspaceID: "ws_123",
	})
	if err != nil {
		t.Fatal(err)
	}

	run, err := repository.GetRun(context.Background(), accepted.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if run.WorkspaceID != "ws_123" || run.Status != facts.RunStatusCreated {
		t.Fatalf("run fact not populated: %+v", run)
	}
}

func TestFactCommandsAppendInputTurnForMessageAndAction(t *testing.T) {
	// message/action 入口的用户输入必须先成为 Product Facts，SSE 和 projection 才能同源读取。
	repository := facts.NewMemoryRepository()
	commands := execution.NewFactCommands(repository)
	ctx := context.Background()

	messageRun, err := commands.StartMessage(ctx, execution.MessageCommand{
		WorkspaceID:     "ws_123",
		Message:         "查询资产",
		ClientRequestID: "client-message",
	})
	if err != nil {
		t.Fatal(err)
	}
	actionRun, err := commands.StartAction(ctx, execution.ActionCommand{
		WorkspaceID:     "ws_123",
		ActionID:        "action-demo",
		ClientRequestID: "client-action",
		InputText:       "查询漏洞",
	})
	if err != nil {
		t.Fatal(err)
	}

	messageSnapshot, err := repository.GetSnapshot(ctx, messageRun.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if len(messageSnapshot.Turns) != 1 || messageSnapshot.Turns[0].Content != "查询资产" {
		t.Fatalf("message turns = %+v", messageSnapshot.Turns)
	}
	actionSnapshot, err := repository.GetSnapshot(ctx, actionRun.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if len(actionSnapshot.Turns) != 1 || actionSnapshot.Turns[0].Content != "查询漏洞" {
		t.Fatalf("action turns = %+v", actionSnapshot.Turns)
	}
}

func TestRunnerCommandsExecuteChatModelRunnerAfterRunCreated(t *testing.T) {
	// Phase 3 服务路径必须在创建 run 后执行 ChatModelRunner，确保 assistant/context facts 落库。
	repository := facts.NewMemoryRepository()
	runner := &recordingRunner{}
	commands := execution.NewRunnerCommands(repository, runner)

	accepted, err := commands.StartMessage(context.Background(), execution.MessageCommand{
		WorkspaceID:     "ws_123",
		Message:         "hello",
		ClientRequestID: "client-runner",
	})
	if err != nil {
		t.Fatal(err)
	}
	if runner.runID != accepted.RunID {
		t.Fatalf("runner runID = %q, want %q", runner.runID, accepted.RunID)
	}
}

func TestRunnerCommandsRejectUnknownCapabilityHintBeforeRunCreated(t *testing.T) {
	// 未注册 capability hint 不能创建 run，避免未知工具名绕过注册表进入执行链路。
	repository := facts.NewMemoryRepository()
	commands := execution.NewRunnerCommandsWithRegistry(repository, &recordingRunner{}, capabilities.NewRegistry())

	_, err := commands.StartAction(context.Background(), execution.ActionCommand{
		WorkspaceID:     "ws_123",
		ActionID:        "action-demo",
		ClientRequestID: "client-action",
		CapabilityHint:  "missing.capability",
		InputText:       "hello",
	})
	if !errors.Is(err, execution.ErrCapabilityNotRegistered) {
		t.Fatalf("err = %v, want ErrCapabilityNotRegistered", err)
	}
	if _, latestErr := repository.LatestRun(context.Background(), "ws_123"); !errors.Is(latestErr, facts.ErrNotFound) {
		t.Fatalf("latest run err = %v, want ErrNotFound", latestErr)
	}
}

func TestRunnerCommandsRecordRegisteredCapabilitySelection(t *testing.T) {
	// 已注册能力只写入安全选择摘要；Phase 4 前不在命令层按工具名执行 provider。
	repository := facts.NewMemoryRepository()
	registry := capabilities.NewRegistry()
	if err := registry.Register(capabilities.Capability{
		ID:          "safe.read",
		ProviderID:  "demo",
		ToolName:    "safe_read",
		DisplayName: "只读查询",
		RiskLevel:   capabilities.RiskReadOnly,
	}); err != nil {
		t.Fatal(err)
	}
	commands := execution.NewRunnerCommandsWithRegistry(repository, &recordingRunner{}, registry)

	accepted, err := commands.StartAction(context.Background(), execution.ActionCommand{
		WorkspaceID:     "ws_123",
		ActionID:        "action-demo",
		ClientRequestID: "client-action",
		CapabilityHint:  "safe.read",
		InputText:       "hello",
	})
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := repository.GetSnapshot(context.Background(), accepted.RunID)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, event := range snapshot.AuditEvents {
		if event.EventType == "tool" && strings.Contains(event.SafeSummary, "safe.read") && strings.Contains(event.SafeSummary, "allowed") {
			found = true
		}
	}
	if !found {
		t.Fatalf("selection audit not found: %+v", snapshot.AuditEvents)
	}
}

func TestRunnerCommandsRunReadOnlyCapabilityThroughToolRunner(t *testing.T) {
	// 只读 capability 选择后应进入 tool runner，而不是继续走普通 ChatModel runner。
	repository := facts.NewMemoryRepository()
	chatRunner := &recordingRunner{}
	toolRunner := &recordingCapabilityRunner{}
	registry := capabilities.NewRegistry()
	if err := registry.Register(capabilities.Capability{
		ID:          "safe.read",
		ProviderID:  "demo",
		ToolName:    "safe_read",
		DisplayName: "只读查询",
		RiskLevel:   capabilities.RiskReadOnly,
	}); err != nil {
		t.Fatal(err)
	}
	commands := execution.NewToolRunnerCommandsWithRegistry(repository, chatRunner, toolRunner, registry)

	accepted, err := commands.StartAction(context.Background(), execution.ActionCommand{
		WorkspaceID:     "ws_123",
		ActionID:        "action-demo",
		ClientRequestID: "client-action",
		CapabilityHint:  "safe.read",
		InputText:       "hello",
	})
	if err != nil {
		t.Fatal(err)
	}
	if chatRunner.runID != "" {
		t.Fatalf("chat runner should not run capability path, got %q", chatRunner.runID)
	}
	if toolRunner.runID != accepted.RunID || toolRunner.capabilityID != "safe.read" || toolRunner.inputText != "hello" {
		t.Fatalf("tool runner mismatch: %+v", toolRunner)
	}
}

func TestRunnerCommandsDoNotRunApprovalRequiredCapabilityBeforeHITL(t *testing.T) {
	// 写域能力在 Phase 6 前只能留下选择审计，不能继续进入 runner 或 provider 执行。
	repository := facts.NewMemoryRepository()
	runner := &recordingRunner{}
	registry := capabilities.NewRegistry()
	if err := registry.Register(capabilities.Capability{
		ID:          "danger.write",
		ProviderID:  "demo",
		ToolName:    "danger_write",
		DisplayName: "写入",
		RiskLevel:   capabilities.RiskWrite,
	}); err != nil {
		t.Fatal(err)
	}
	commands := execution.NewRunnerCommandsWithRegistry(repository, runner, registry)

	accepted, err := commands.StartAction(context.Background(), execution.ActionCommand{
		WorkspaceID:     "ws_123",
		ActionID:        "action-demo",
		ClientRequestID: "client-action",
		CapabilityHint:  "danger.write",
		InputText:       "write something",
	})
	if err != nil {
		t.Fatal(err)
	}
	if runner.runID != "" {
		t.Fatalf("runner should not run approval-required capability, got %q", runner.runID)
	}
	run, err := repository.GetRun(context.Background(), accepted.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if run.Status != facts.RunStatusCreated {
		t.Fatalf("run status = %q, want created", run.Status)
	}
}

func TestRunnerCommandsRecordPolicyBlockedCapabilityAsSafeFailedRun(t *testing.T) {
	// 缺凭据、scope denied、connector unavailable 这类策略阻断必须返回可投影的安全失败状态，不能变成 500。
	repository := facts.NewMemoryRepository()
	runner := &recordingRunner{}
	toolRunner := &recordingCapabilityRunner{}
	registry := capabilities.NewRegistry()
	if err := registry.Register(fobrainReadCapabilityForPolicyTest()); err != nil {
		t.Fatal(err)
	}
	commands := execution.NewToolRunnerCommandsWithRegistry(repository, runner, toolRunner, registry)

	accepted, err := commands.StartAction(context.Background(), execution.ActionCommand{
		WorkspaceID:     "ws_123",
		ActionID:        "action-demo",
		ClientRequestID: "client-action",
		CapabilityHint:  "tool.fobrain.asset.read",
		InputText:       "asset",
	})
	if err != nil {
		t.Fatal(err)
	}
	if runner.runID != "" || toolRunner.runID != "" {
		t.Fatalf("policy blocked capability must not execute: runner=%q tool=%q", runner.runID, toolRunner.runID)
	}
	run, err := repository.GetRun(context.Background(), accepted.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if run.Status != facts.RunStatusFailed || run.SafeError != "provider_policy_blocked" {
		t.Fatalf("run policy failure mismatch: %+v", run)
	}
}

func TestRunnerCommandsRecordPolicyContextBlocksAsSafeFailedRun(t *testing.T) {
	for _, tc := range []struct {
		name          string
		policyContext capabilities.PolicyContext
	}{
		{
			name: "scope denied",
			policyContext: capabilities.PolicyContext{
				WorkspaceID: "ws_123",
				CredentialBinding: capabilities.CredentialBinding{
					WorkspaceID: "ws_other",
					System:      "fobrain",
					Status:      capabilities.CredentialStatusBound,
					DisplayRef:  "bound:fobrain:main",
					OwnerScope:  capabilities.PermissionScopeWorkspace,
				},
				ConnectorStatus: capabilities.ConnectorStatusAvailable,
			},
		},
		{
			name: "connector unavailable",
			policyContext: capabilities.PolicyContext{
				WorkspaceID: "ws_123",
				CredentialBinding: capabilities.CredentialBinding{
					WorkspaceID: "ws_123",
					System:      "fobrain",
					Status:      capabilities.CredentialStatusBound,
					DisplayRef:  "bound:fobrain:main",
					OwnerScope:  capabilities.PermissionScopeWorkspace,
				},
				ConnectorStatus: capabilities.ConnectorStatusUnavailable,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repository := facts.NewMemoryRepository()
			toolRunner := &recordingCapabilityRunner{}
			registry := capabilities.NewRegistry()
			if err := registry.Register(fobrainReadCapabilityForPolicyTest()); err != nil {
				t.Fatal(err)
			}
			commands := execution.NewToolRunnerCommandsWithRegistry(repository, &recordingRunner{}, toolRunner, registry)

			accepted, err := commands.StartAction(context.Background(), execution.ActionCommand{
				WorkspaceID:     "ws_123",
				ActionID:        "action-demo",
				ClientRequestID: "client-action",
				CapabilityHint:  "tool.fobrain.asset.read",
				InputText:       "asset",
				PolicyContext:   tc.policyContext,
			})
			if err != nil {
				t.Fatal(err)
			}
			if toolRunner.runID != "" {
				t.Fatalf("policy blocked capability must not execute: %+v", toolRunner)
			}
			run, err := repository.GetRun(context.Background(), accepted.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if run.Status != facts.RunStatusFailed || run.SafeError != "provider_policy_blocked" {
				t.Fatalf("run policy failure mismatch: %+v", run)
			}
		})
	}
}

func TestFactCommandsResumeRequiresExistingRun(t *testing.T) {
	// resume 不能凭前端传参创建隐式 run，必须绑定已存在的 run 生命周期。
	repository := facts.NewMemoryRepository()
	commands := execution.NewFactCommands(repository)

	_, err := commands.Resume(context.Background(), execution.ResumeCommand{
		WorkspaceID: "ws_123",
		RunID:       "run_missing",
		ResumeRef:   "resume_ref_1",
	})
	if !errors.Is(err, execution.ErrRunNotFound) {
		t.Fatalf("err = %v, want ErrRunNotFound", err)
	}
}

type recordingRunner struct {
	runID string
}

func (runner *recordingRunner) Run(_ context.Context, runID string) error {
	runner.runID = runID
	return nil
}

type recordingCapabilityRunner struct {
	runID        string
	capabilityID string
	inputText    string
}

func (runner *recordingCapabilityRunner) RunCapability(_ context.Context, runID string, capabilityID string, inputText string) error {
	runner.runID = runID
	runner.capabilityID = capabilityID
	runner.inputText = inputText
	return nil
}

func fobrainReadCapabilityForPolicyTest() capabilities.Capability {
	return capabilities.Capability{
		ID:                      "tool.fobrain.asset.read",
		ProviderID:              "fobrain",
		ToolName:                "fobrain_asset_read",
		DisplayName:             "资产查询",
		RiskLevel:               capabilities.RiskReadOnly,
		SideEffect:              capabilities.SideEffectReadExternal,
		PolicyRef:               "policy:fobrain:read:v1",
		PermissionScope:         capabilities.PermissionScopeWorkspace,
		CredentialBindingPolicy: capabilities.CredentialBindingRequired,
		ConnectorID:             "fobrain",
	}
}
