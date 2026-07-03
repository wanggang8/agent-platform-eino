package execution

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/facts"
)

// ErrRunNotFound 表示 resume 请求引用了不存在的 run。
var ErrRunNotFound = errors.New("run not found")

// ErrCapabilityNotRegistered 表示 Action API 请求的 capability hint 未通过注册表门禁。
var ErrCapabilityNotRegistered = errors.New("capability not registered")

// ErrLifecycleNotAllowed 表示当前 run 状态不允许请求的 lifecycle 操作。
var ErrLifecycleNotAllowed = errors.New("run lifecycle action not allowed")

// ErrResumeNotAllowed 表示当前 run 或 pending 状态不允许恢复。
var ErrResumeNotAllowed = errors.New("run resume not allowed")

// LifecycleAction 是 run 生命周期控制动作，必须与 docs/run-lifecycle.md 保持一致。
type LifecycleAction string

const (
	LifecycleActionCancel          LifecycleAction = "cancel"
	LifecycleActionStop            LifecycleAction = "stop"
	LifecycleActionProviderTimeout LifecycleAction = "provider_timeout"
	LifecycleActionPendingTimeout  LifecycleAction = "pending_timeout"
	LifecycleActionRetry           LifecycleAction = "retry"
)

// AcceptedRun 是 HTTP/API 边界可以返回的最小接收结果，不包含执行事件。
type AcceptedRun struct {
	RunID  string
	Status string
}

// MessageCommand 表示 Workbench message 入口的执行命令。
type MessageCommand struct {
	WorkspaceID     string
	Message         string
	ClientRequestID string
	RunID           string
}

// ActionCommand 表示外部 Action API 入口的执行命令。
type ActionCommand struct {
	WorkspaceID     string
	ActionID        string
	ClientRequestID string
	CapabilityHint  string
	InputText       string
	Attachments     []Attachment
	Context         RequestContext
	PolicyContext   capabilities.PolicyContext
}

// ResumeCommand 表示 approval/clarification 的恢复请求。
type ResumeCommand struct {
	WorkspaceID     string
	RunID           string
	ResumeRef       string
	ClientRequestID string
	Decision        string
	SelectedRefs    []string
	FreeText        string
	Comment         string
}

// LifecycleCommand 表示 run cancel/stop/timeout/retry 等生命周期控制请求。
type LifecycleCommand struct {
	WorkspaceID     string
	RunID           string
	Action          LifecycleAction
	ClientRequestID string
}

// Attachment 是 Action API 的安全附件摘要，不包含文件原文。
type Attachment struct {
	AttachmentRef string
	MediaType     string
	SafeName      string
}

// RequestContext 是请求侧安全上下文摘要，可进入后续 Product Facts。
type RequestContext struct {
	Timezone      string
	Locale        string
	SafeUserLabel string
}

// Commands 是 HTTP 层调用 execution 的唯一接口。
type Commands interface {
	StartMessage(context.Context, MessageCommand) (AcceptedRun, error)
	StartAction(context.Context, ActionCommand) (AcceptedRun, error)
	Resume(context.Context, ResumeCommand) (AcceptedRun, error)
	RunLifecycle(context.Context, LifecycleCommand) (AcceptedRun, error)
}

// Runner 是命令层触发执行的最小接口，生产实现由 ChatModelRunner 提供。
type Runner interface {
	Run(context.Context, string) error
}

// CapabilityRunner 执行已经通过 registry/policy 选择的只读能力。
type CapabilityRunner interface {
	RunCapability(context.Context, string, string, string) error
}

// StaticCommands 是 Phase 1/3 的最小命令实现，负责创建基础 run fact。
type StaticCommands struct {
	repository       facts.Repository
	newRunID         func() (string, error)
	runner           Runner
	capabilityRunner CapabilityRunner
	registry         *capabilities.Registry
	policyContexts   map[string]capabilities.PolicyContext
	checkpoints      CheckpointResolver
}

// NewStaticCommands 创建不落 Product Facts 的静态命令替身。
func NewStaticCommands() StaticCommands {
	return StaticCommands{newRunID: randomRunID}
}

// NewFactCommands 创建会写入 Product Facts 的命令实现。
func NewFactCommands(repository facts.Repository) StaticCommands {
	return StaticCommands{
		repository: repository,
		newRunID:   randomRunID,
	}
}

// NewRunnerCommands 创建会写入 Product Facts 并触发 ChatModelRunner 的命令实现。
func NewRunnerCommands(repository facts.Repository, runner Runner) StaticCommands {
	return StaticCommands{
		repository: repository,
		newRunID:   randomRunID,
		runner:     runner,
	}
}

// NewRunnerCommandsWithRegistry 创建带 capability registry 的命令实现，用于 Action API 选择门禁。
func NewRunnerCommandsWithRegistry(repository facts.Repository, runner Runner, registry *capabilities.Registry) StaticCommands {
	return StaticCommands{
		repository: repository,
		newRunID:   randomRunID,
		runner:     runner,
		registry:   registry,
	}
}

// NewToolRunnerCommandsWithRegistry 创建同时支持 chat runner 和只读 capability tool runner 的命令实现。
func NewToolRunnerCommandsWithRegistry(repository facts.Repository, runner Runner, capabilityRunner CapabilityRunner, registry *capabilities.Registry) StaticCommands {
	return StaticCommands{
		repository:       repository,
		newRunID:         randomRunID,
		runner:           runner,
		capabilityRunner: capabilityRunner,
		registry:         registry,
	}
}

// WithPolicyContexts 绑定配置派生的安全 policy context，避免 HTTP 请求携带真实凭据。
func (commands StaticCommands) WithPolicyContexts(contexts map[string]capabilities.PolicyContext) StaticCommands {
	if len(contexts) == 0 {
		return commands
	}
	commands.policyContexts = make(map[string]capabilities.PolicyContext, len(contexts))
	for capabilityID, context := range contexts {
		commands.policyContexts[capabilityID] = context
	}
	return commands
}

// WithRunIDGenerator 仅供测试注入稳定 run id，生产路径继续使用 opaque 随机 id。
func (commands StaticCommands) WithRunIDGenerator(generator func() (string, error)) StaticCommands {
	commands.newRunID = generator
	return commands
}

// WithCheckpointResolver 绑定内部 checkpoint resolver，用于 resume 前校验恢复材料是否存在。
func (commands StaticCommands) WithCheckpointResolver(resolver CheckpointResolver) StaticCommands {
	commands.checkpoints = resolver
	return commands
}

// StartMessage 接收 Workbench 消息并创建 run。
func (commands StaticCommands) StartMessage(ctx context.Context, command MessageCommand) (AcceptedRun, error) {
	accepted, err := commands.acceptRun(ctx, command.WorkspaceID, command.RunID, command.Message)
	if err != nil {
		return AcceptedRun{}, err
	}
	return commands.runIfConfigured(ctx, accepted)
}

// StartAction 接收 Action API 请求并创建 run。
func (commands StaticCommands) StartAction(ctx context.Context, command ActionCommand) (AcceptedRun, error) {
	selection := SelectCapability(commands.registry, SelectionRequest{
		InputText:      command.InputText,
		CapabilityHint: command.CapabilityHint,
		ProductAction:  command.ActionID,
		PolicyContext:  commands.policyContextForAction(command),
	})
	if selection.Mode == SelectionModeRejected {
		return AcceptedRun{}, ErrCapabilityNotRegistered
	}
	accepted, err := commands.acceptRun(ctx, command.WorkspaceID, "", command.InputText)
	if err != nil {
		return AcceptedRun{}, err
	}
	if err := commands.recordSelection(ctx, accepted.RunID, selection); err != nil {
		return AcceptedRun{}, err
	}
	if selection.Mode == SelectionModeCapability && !selection.PolicyAllowed && !selection.RequiresApproval {
		if err := commands.recordPolicyBlockedRun(ctx, accepted.RunID); err != nil {
			return AcceptedRun{}, err
		}
		return accepted, nil
	}
	if selection.RequiresApproval {
		return accepted, nil
	}
	if selection.Mode == SelectionModeCapability && commands.capabilityRunner != nil {
		if err := commands.capabilityRunner.RunCapability(ctx, accepted.RunID, selection.CapabilityID, command.InputText); err != nil {
			return AcceptedRun{}, err
		}
		return accepted, nil
	}
	return commands.runIfConfigured(ctx, accepted)
}

// policyContextForAction 只传递安全策略上下文；真实凭据仍由 provider 边界解析。
func (commands StaticCommands) policyContextForAction(command ActionCommand) capabilities.PolicyContext {
	context := commands.policyContexts[command.CapabilityHint]
	context = mergePolicyContext(context, command.PolicyContext)
	context.WorkspaceID = command.WorkspaceID
	return context
}

// mergePolicyContext 让请求显式安全上下文覆盖配置默认值，但不要求请求携带真实凭据。
func mergePolicyContext(base capabilities.PolicyContext, override capabilities.PolicyContext) capabilities.PolicyContext {
	if override.WorkspaceID != "" {
		base.WorkspaceID = override.WorkspaceID
	}
	if override.CredentialBinding.Status != "" || override.CredentialBinding.WorkspaceID != "" {
		base.CredentialBinding = override.CredentialBinding
	}
	if override.ConnectorStatus != capabilities.ConnectorStatusUnknown {
		base.ConnectorStatus = override.ConnectorStatus
	}
	if override.CallerID != "" {
		base.CallerID = override.CallerID
	}
	return base
}

// Resume 校验 run 存在后接收恢复请求；完整 HITL 执行在后续 Phase 接入。
func (commands StaticCommands) Resume(ctx context.Context, command ResumeCommand) (AcceptedRun, error) {
	if commands.repository != nil {
		run, err := commands.repository.GetRun(ctx, command.RunID)
		if err != nil {
			if errors.Is(err, facts.ErrNotFound) {
				return AcceptedRun{}, ErrRunNotFound
			}
			return AcceptedRun{}, err
		}
		if run.WorkspaceID != command.WorkspaceID {
			return AcceptedRun{}, ErrRunNotFound
		}
		if run.Status != facts.RunStatusWaiting {
			return AcceptedRun{}, ErrResumeNotAllowed
		}
		if commands.checkpoints != nil {
			if err := commands.ensureCheckpointAvailable(ctx, command); err != nil {
				return AcceptedRun{}, err
			}
		}
	}
	return AcceptedRun{RunID: command.RunID, Status: "accepted"}, nil
}

// ensureCheckpointAvailable 在消费 resume_ref 前确认内部 checkpoint 存在，避免缺失后破坏 pending。
func (commands StaticCommands) ensureCheckpointAvailable(ctx context.Context, command ResumeCommand) error {
	pending, err := commands.repository.GetPendingByResumeRef(ctx, command.ResumeRef)
	if err != nil {
		if errors.Is(err, facts.ErrNotFound) {
			return ErrRunNotFound
		}
		return err
	}
	if pending.RunID != command.RunID {
		return ErrRunNotFound
	}
	if pending.Status != facts.PendingStatusWaiting && pending.Status != facts.PendingStatusSubmitted {
		return ErrResumeNotAllowed
	}
	_, existed, err := commands.checkpoints.ResolveCheckpointID(ctx, pending.CheckpointRef, pending.RunID, pending.PendingID)
	if err != nil {
		return err
	}
	if !existed {
		return ErrCheckpointMissing
	}
	return nil
}

// RunLifecycle 按 Product Facts 状态机执行 run 生命周期控制。
func (commands StaticCommands) RunLifecycle(ctx context.Context, command LifecycleCommand) (AcceptedRun, error) {
	if commands.repository == nil {
		return AcceptedRun{RunID: command.RunID, Status: "accepted"}, nil
	}
	snapshot, err := commands.repository.GetSnapshot(ctx, command.RunID)
	if err != nil {
		if errors.Is(err, facts.ErrNotFound) {
			return AcceptedRun{}, ErrRunNotFound
		}
		return AcceptedRun{}, err
	}
	if snapshot.Run.WorkspaceID != command.WorkspaceID {
		return AcceptedRun{}, ErrRunNotFound
	}
	if snapshot.Run.Status.Terminal() && command.Action != LifecycleActionRetry {
		return AcceptedRun{RunID: snapshot.Run.RunID, Status: string(snapshot.Run.Status)}, nil
	}

	switch command.Action {
	case LifecycleActionCancel:
		if snapshot.Run.Status != facts.RunStatusRunning && snapshot.Run.Status != facts.RunStatusWaiting {
			return AcceptedRun{}, ErrLifecycleNotAllowed
		}
		return commands.transitionRun(ctx, snapshot, facts.RunStatusCancelled, "user_cancelled", facts.PendingStatusCancelled, facts.ToolCallCancelled, activeToolCallIDs(snapshot))
	case LifecycleActionStop:
		if snapshot.Run.Status != facts.RunStatusRunning {
			return AcceptedRun{}, ErrLifecycleNotAllowed
		}
		return commands.transitionRun(ctx, snapshot, facts.RunStatusStopped, "user_stopped", "", facts.ToolCallCancelled, activeToolCallIDs(snapshot))
	case LifecycleActionProviderTimeout:
		toolCallIDs := activeToolCallIDs(snapshot)
		if snapshot.Run.Status != facts.RunStatusRunning || len(toolCallIDs) == 0 {
			return AcceptedRun{}, ErrLifecycleNotAllowed
		}
		return commands.transitionRun(ctx, snapshot, facts.RunStatusFailed, "provider_timeout", "", facts.ToolCallFailed, toolCallIDs)
	case LifecycleActionPendingTimeout:
		pendingIDs := activePendingIDs(snapshot)
		if snapshot.Run.Status != facts.RunStatusWaiting || len(pendingIDs) == 0 {
			return AcceptedRun{}, ErrLifecycleNotAllowed
		}
		return commands.transitionRun(ctx, snapshot, facts.RunStatusFailed, "pending_timeout", facts.PendingStatusExpired, "", nil)
	case LifecycleActionRetry:
		return commands.retryRun(ctx, command, snapshot)
	default:
		return AcceptedRun{}, ErrLifecycleNotAllowed
	}
}

// transitionRun 更新 run 状态，并按需关闭 waiting pending，保证旧 resume_ref 不再可继续。
func (commands StaticCommands) transitionRun(ctx context.Context, snapshot facts.Snapshot, status facts.RunStatus, safeError string, pendingStatus facts.PendingStatus, toolStatus facts.ToolCallStatus, toolCallIDs []string) (AcceptedRun, error) {
	now := time.Now().UTC()
	pendingIDs := []string(nil)
	if pendingStatus != "" {
		for _, pending := range snapshot.PendingInteractions {
			if pending.Status == facts.PendingStatusWaiting || pending.Status == facts.PendingStatusSubmitted {
				pendingIDs = append(pendingIDs, pending.PendingID)
			}
		}
	}
	transition := facts.LifecycleTransition{
		RunID:               snapshot.Run.RunID,
		ExpectedRunStatuses: []facts.RunStatus{snapshot.Run.Status},
		Status:              status,
		SafeError:           safeError,
		UpdatedAt:           now,
		PendingIDs:          pendingIDs,
		PendingStatus:       pendingStatus,
		ToolCallIDs:         toolCallIDs,
		ToolStatus:          toolStatus,
		AuditEvent: facts.AuditEvent{
			AuditID:     snapshot.Run.RunID + ":audit:lifecycle:" + string(status),
			RunID:       snapshot.Run.RunID,
			EventType:   "lifecycle",
			SafeSummary: "run lifecycle: " + string(status),
			Actor:       "system",
			CreatedAt:   now,
		},
	}
	if err := commands.repository.ApplyLifecycleTransition(ctx, transition); err != nil {
		return AcceptedRun{}, err
	}
	return AcceptedRun{RunID: snapshot.Run.RunID, Status: string(status)}, nil
}

// retryRun 只对明确安全可重试的失败 run 创建新 run，不复用旧 checkpoint/resume。
func (commands StaticCommands) retryRun(ctx context.Context, command LifecycleCommand, snapshot facts.Snapshot) (AcceptedRun, error) {
	if snapshot.Run.Status != facts.RunStatusFailed || !retryableSafeError(snapshot.Run.SafeError) {
		return AcceptedRun{}, ErrLifecycleNotAllowed
	}
	runID, err := commands.nextRunID()
	if err != nil {
		return AcceptedRun{}, err
	}
	now := time.Now().UTC()
	record := facts.IdempotencyRecord{
		Scope:       facts.IdempotencyScopeMutation,
		Key:         command.ClientRequestID,
		RunID:       runID,
		ResourceRef: "retry:" + snapshot.Run.RunID,
		Status:      "accepted",
		CreatedAt:   now,
	}
	transition := facts.RetryRunTransition{
		OldRunID:             snapshot.Run.RunID,
		ExpectedOldRunStatus: facts.RunStatusFailed,
		ExpectedOldSafeError: snapshot.Run.SafeError,
		NewRun: facts.Run{
			RunID:       runID,
			WorkspaceID: snapshot.Run.WorkspaceID,
			Status:      facts.RunStatusCreated,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		UserTurn: facts.Turn{
			TurnID:    runID + ":user:1",
			RunID:     runID,
			Role:      facts.TurnRoleUser,
			Content:   retryInputText(snapshot),
			Sequence:  1,
			CreatedAt: now,
		},
		OldAudit: facts.AuditEvent{
			AuditID:     snapshot.Run.RunID + ":audit:lifecycle:retry:" + runID,
			RunID:       snapshot.Run.RunID,
			EventType:   "lifecycle",
			SafeSummary: "run retry requested: " + runID,
			Actor:       "user",
			CreatedAt:   now,
		},
		NewAudit: facts.AuditEvent{
			AuditID:     runID + ":audit:lifecycle:retry",
			RunID:       runID,
			EventType:   "lifecycle",
			SafeSummary: "run retry from: " + snapshot.Run.RunID,
			Actor:       "system",
			CreatedAt:   now,
		},
		Idempotency: record,
	}
	record, existed, err := commands.repository.CreateRetryRun(ctx, transition)
	if err != nil {
		return AcceptedRun{}, err
	}
	if existed {
		return AcceptedRun{RunID: record.RunID, Status: "accepted"}, nil
	}
	return AcceptedRun{RunID: record.RunID, Status: "accepted"}, nil
}

// acceptRun 统一创建 run，保证 Message 和 Action 入口共享 Product Facts。
func (commands StaticCommands) acceptRun(ctx context.Context, workspaceID string, requestedRunID string, inputText string) (AcceptedRun, error) {
	runID := strings.TrimSpace(requestedRunID)
	if runID == "" {
		generated, err := commands.nextRunID()
		if err != nil {
			return AcceptedRun{}, err
		}
		runID = generated
	}
	if commands.repository != nil {
		now := time.Now().UTC()
		if err := commands.repository.CreateRun(ctx, facts.Run{
			RunID:       runID,
			WorkspaceID: workspaceID,
			Status:      facts.RunStatusCreated,
			CreatedAt:   now,
			UpdatedAt:   now,
		}); err != nil {
			return AcceptedRun{}, err
		}
		if strings.TrimSpace(inputText) != "" {
			if err := commands.repository.AppendTurn(ctx, facts.Turn{
				TurnID:    runID + ":user:1",
				RunID:     runID,
				Role:      facts.TurnRoleUser,
				Content:   inputText,
				Sequence:  1,
				CreatedAt: now,
			}); err != nil {
				return AcceptedRun{}, err
			}
			if err := commands.repository.AppendAuditEvent(ctx, facts.AuditEvent{
				AuditID:     runID + ":audit:message:1",
				RunID:       runID,
				EventType:   "message",
				SafeSummary: "user message accepted",
				Actor:       "user",
				CreatedAt:   now,
			}); err != nil {
				return AcceptedRun{}, err
			}
		}
	}
	return AcceptedRun{RunID: runID, Status: "accepted"}, nil
}

func retryableSafeError(safeError string) bool {
	switch safeError {
	case "schema_invalid":
		return true
	default:
		return false
	}
}

func retryInputText(snapshot facts.Snapshot) string {
	return "retry run " + snapshot.Run.RunID
}

func (commands StaticCommands) nextRunID() (string, error) {
	if commands.newRunID == nil {
		return randomRunID()
	}
	return commands.newRunID()
}

func activePendingIDs(snapshot facts.Snapshot) []string {
	ids := []string{}
	for _, pending := range snapshot.PendingInteractions {
		if pending.Status == facts.PendingStatusWaiting || pending.Status == facts.PendingStatusSubmitted {
			ids = append(ids, pending.PendingID)
		}
	}
	return ids
}

func activeToolCallIDs(snapshot facts.Snapshot) []string {
	ids := []string{}
	for _, call := range snapshot.ToolCalls {
		if call.Status == facts.ToolCallQueued || call.Status == facts.ToolCallRunning {
			ids = append(ids, call.ToolCallID)
		}
	}
	return ids
}

// recordSelection 将 capability 选择结果写成安全审计事实，后续 tool execution 只能读取该事实链路。
func (commands StaticCommands) recordSelection(ctx context.Context, runID string, selection SelectionResult) error {
	if commands.repository == nil || selection.Mode != SelectionModeCapability {
		return nil
	}
	return commands.repository.AppendAuditEvent(ctx, facts.AuditEvent{
		AuditID:     runID + ":audit:capability:1",
		RunID:       runID,
		EventType:   "tool",
		SafeSummary: "capability selected: " + selection.CapabilityID + "; policy: " + safePolicyReason(selection.PolicyReason),
		Actor:       "system",
		CreatedAt:   time.Now().UTC(),
	})
}

// recordPolicyBlockedRun 将非审批类策略阻断写成安全失败状态，避免 HTTP 层返回 500。
func (commands StaticCommands) recordPolicyBlockedRun(ctx context.Context, runID string) error {
	if commands.repository == nil {
		return nil
	}
	return commands.repository.UpdateRunStatus(ctx, runID, facts.RunStatusFailed, "provider_policy_blocked", time.Now().UTC())
}

// runIfConfigured 在生产路径触发 runner；测试替身可不配置 runner。
func (commands StaticCommands) runIfConfigured(ctx context.Context, accepted AcceptedRun) (AcceptedRun, error) {
	if commands.runner == nil {
		return accepted, nil
	}
	if err := commands.runner.Run(ctx, accepted.RunID); err != nil {
		return AcceptedRun{}, err
	}
	return accepted, nil
}

// safePolicyReason 把含 credential/token 等敏感标记的内部 reason 折叠为可入库摘要。
func safePolicyReason(reason string) string {
	if facts.ContainsUnsafeMaterial(reason) {
		return "blocked"
	}
	return reason
}

// randomRunID 生成不含 workspace/action 信息的 opaque run id。
func randomRunID() (string, error) {
	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	return "run_" + hex.EncodeToString(bytes[:]), nil
}
