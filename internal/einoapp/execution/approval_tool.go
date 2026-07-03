package execution

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"agent-platform-eino/internal/einoapp/facts"
)

// ApprovedCapabilityRunner 表示已经通过 HITL approval 的能力执行入口。
// 生产 ToolLoopRunner 用它跳过“仍需审批”的短路，但继续保留凭据、scope 和 connector 策略校验。
type ApprovedCapabilityRunner interface {
	RunApprovedCapability(context.Context, string, string, string) error
}

type approvalDecision string

const (
	approvalDecisionApprove approvalDecision = "approve"
	approvalDecisionReject  approvalDecision = "reject"
)

type approvalContinuation struct {
	SchemaVersion string `json:"schema_version"`
	CapabilityID  string `json:"capability_id"`
	InputText     string `json:"input_text"`
}

// requestApproval 把需要人工审批的 capability 转成 Product Facts 等待态和内部 checkpoint continuation。
func (commands StaticCommands) requestApproval(ctx context.Context, runID string, command ActionCommand, selection SelectionResult) error {
	if commands.repository == nil || commands.approvalStore == nil {
		return ErrCheckpointMissing
	}
	capability, ok := commands.registry.Get(selection.CapabilityID)
	if !ok {
		return ErrCapabilityNotRegistered
	}
	now := time.Now().UTC()
	pendingID := runID + ":pending:approval:1"
	resumeRef, err := randomSafeRef("resume_ref:")
	if err != nil {
		return err
	}
	checkpointRef, err := randomSafeRef(facts.CheckpointRefPrefix)
	if err != nil {
		return err
	}
	checkpointID, err := randomSafeRef("eino_checkpoint:")
	if err != nil {
		return err
	}
	payload, err := json.Marshal(approvalContinuation{
		SchemaVersion: "eino.approval_continuation.v1",
		CapabilityID:  selection.CapabilityID,
		InputText:     command.InputText,
	})
	if err != nil {
		return err
	}
	if err := commands.approvalStore.Set(ctx, checkpointID, payload); err != nil {
		return err
	}
	if err := commands.approvalStore.BindCheckpointRef(ctx, checkpointRef, checkpointID, runID, pendingID, now); err != nil {
		return err
	}
	if err := commands.repository.AppendPendingInteraction(ctx, facts.PendingInteraction{
		PendingID:     pendingID,
		RunID:         runID,
		Kind:          facts.PendingKindApproval,
		Status:        facts.PendingStatusWaiting,
		ResumeRef:     resumeRef,
		CheckpointRef: checkpointRef,
		Question:      "请审批此操作",
		OperationName: safeApprovalOperation(capability.DisplayName),
		RiskSummary:   "写域能力需要人工审批",
		TargetSummary: safeApprovalTarget(command.ActionID),
		ExpiresAt:     time.Time{},
	}); err != nil {
		return err
	}
	if err := commands.repository.UpdateRunStatus(ctx, runID, facts.RunStatusWaiting, "", now); err != nil {
		return err
	}
	return commands.repository.AppendAuditEvent(ctx, facts.AuditEvent{
		AuditID:     runID + ":audit:approval:requested",
		RunID:       runID,
		EventType:   "approval",
		SafeSummary: "approval requested: " + selection.CapabilityID,
		Actor:       "system",
		CreatedAt:   now,
	})
}

// resumeApproval 处理 approval approve/reject，并只在新 approve 成功时继续执行 capability。
func (commands StaticCommands) resumeApproval(ctx context.Context, command ResumeCommand, run facts.Run, pending facts.PendingInteraction) (AcceptedRun, error) {
	decision, err := normalizeApprovalDecision(command.Decision)
	if err != nil {
		return AcceptedRun{}, err
	}
	if strings.TrimSpace(command.ClientRequestID) == "" {
		return AcceptedRun{}, ErrResumeNotAllowed
	}
	transition, err := approvalTransition(command, decision)
	if err != nil {
		return AcceptedRun{}, err
	}
	if pending.Status != facts.PendingStatusWaiting || run.Status != facts.RunStatusWaiting {
		_, _, existed, err := commands.repository.ApplyApprovalResume(ctx, transition)
		if err != nil {
			return AcceptedRun{}, mapResumeTransitionError(err)
		}
		if existed {
			return AcceptedRun{RunID: command.RunID, Status: "accepted"}, nil
		}
		return AcceptedRun{}, ErrResumeNotAllowed
	}

	var continuation approvalContinuation
	if decision == approvalDecisionApprove {
		continuation, err = commands.loadApprovalContinuation(ctx, command, pending)
		if err != nil {
			return AcceptedRun{}, err
		}
	}
	_, _, existed, err := commands.repository.ApplyApprovalResume(ctx, transition)
	if err != nil {
		return AcceptedRun{}, mapResumeTransitionError(err)
	}
	if existed || decision == approvalDecisionReject {
		return AcceptedRun{RunID: command.RunID, Status: "accepted"}, nil
	}
	if err := commands.runApprovedCapability(ctx, command.RunID, continuation); err != nil {
		_ = commands.repository.UpdateRunStatus(ctx, command.RunID, facts.RunStatusFailed, "capability_execution_failed", time.Now().UTC())
		return AcceptedRun{}, err
	}
	return AcceptedRun{RunID: command.RunID, Status: "accepted"}, nil
}

// loadApprovalContinuation 只通过 safe checkpoint_ref 找回内部 continuation，避免产品层暴露 checkpoint id。
func (commands StaticCommands) loadApprovalContinuation(ctx context.Context, command ResumeCommand, pending facts.PendingInteraction) (approvalContinuation, error) {
	checkpointID, existed, err := commands.approvalStore.ResolveCheckpointID(ctx, pending.CheckpointRef, command.RunID, pending.PendingID)
	if err != nil {
		return approvalContinuation{}, err
	}
	if !existed {
		return approvalContinuation{}, ErrCheckpointMissing
	}
	payload, existed, err := commands.approvalStore.Get(ctx, checkpointID)
	if err != nil {
		return approvalContinuation{}, err
	}
	if !existed {
		return approvalContinuation{}, ErrCheckpointMissing
	}
	var continuation approvalContinuation
	if err := json.Unmarshal(payload, &continuation); err != nil {
		return approvalContinuation{}, ErrCheckpointMissing
	}
	if continuation.SchemaVersion != "eino.approval_continuation.v1" || strings.TrimSpace(continuation.CapabilityID) == "" {
		return approvalContinuation{}, ErrCheckpointMissing
	}
	return continuation, nil
}

func (commands StaticCommands) runApprovedCapability(ctx context.Context, runID string, continuation approvalContinuation) error {
	if commands.capabilityRunner == nil {
		return ErrCapabilityNotRegistered
	}
	if approvedRunner, ok := commands.capabilityRunner.(ApprovedCapabilityRunner); ok {
		return approvedRunner.RunApprovedCapability(ctx, runID, continuation.CapabilityID, continuation.InputText)
	}
	return commands.capabilityRunner.RunCapability(ctx, runID, continuation.CapabilityID, continuation.InputText)
}

func approvalTransition(command ResumeCommand, decision approvalDecision) (facts.ApprovalResumeTransition, error) {
	now := time.Now().UTC()
	pendingStatus := facts.PendingStatusApproved
	runStatus := facts.RunStatusRunning
	safeError := ""
	if decision == approvalDecisionReject {
		pendingStatus = facts.PendingStatusRejected
		runStatus = facts.RunStatusFailed
		safeError = "approval_rejected"
	}
	return facts.ApprovalResumeTransition{
		RunID:               command.RunID,
		ResumeRef:           command.ResumeRef,
		ExpectedRunStatus:   facts.RunStatusWaiting,
		ExpectedPendingKind: facts.PendingKindApproval,
		PendingStatus:       pendingStatus,
		RunStatus:           runStatus,
		SafeError:           safeError,
		UpdatedAt:           now,
		Idempotency: facts.IdempotencyRecord{
			Scope:       facts.IdempotencyScopeResume,
			Key:         command.ClientRequestID,
			RunID:       command.RunID,
			ResourceRef: command.ResumeRef,
			Status:      string(pendingStatus),
			CreatedAt:   now,
		},
		AuditEvent: facts.AuditEvent{
			AuditID:     command.RunID + ":audit:approval:" + string(pendingStatus),
			RunID:       command.RunID,
			EventType:   "approval",
			SafeSummary: "approval " + string(pendingStatus),
			Actor:       "user",
			CreatedAt:   now,
		},
	}, nil
}

func normalizeApprovalDecision(value string) (approvalDecision, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "approve", "approved":
		return approvalDecisionApprove, nil
	case "reject", "rejected":
		return approvalDecisionReject, nil
	default:
		return "", ErrResumeNotAllowed
	}
}

func mapResumeTransitionError(err error) error {
	if errors.Is(err, facts.ErrResumeAlreadyConsumed) || errors.Is(err, facts.ErrIdempotencyConflict) {
		return ErrResumeNotAllowed
	}
	return err
}

func randomSafeRef(prefix string) (string, error) {
	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	return prefix + hex.EncodeToString(bytes[:]), nil
}

func safeApprovalOperation(displayName string) string {
	if strings.TrimSpace(displayName) == "" || facts.ContainsUnsafeMaterial(displayName) {
		return "需要审批的操作"
	}
	return displayName
}

func safeApprovalTarget(actionID string) string {
	if strings.TrimSpace(actionID) == "" || facts.ContainsUnsafeMaterial(actionID) {
		return "目标已脱敏"
	}
	return "action:" + actionID
}
