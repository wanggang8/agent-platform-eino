package execution

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
	"time"

	"agent-platform-eino/internal/einoapp/facts"
)

// ClarificationRequest 是实体消歧或参数不足时构造 pending event 的内部命令对象。
// 它只携带安全候选，不包含 raw provider id、resume token 或 raw checkpoint。
type ClarificationRequest struct {
	RunID         string
	PendingID     string
	ResumeRef     string
	CheckpointRef string
	Question      string
	InputMode     facts.PendingInputMode
	Candidates    []facts.PendingCandidate
}

// ErrInvalidClarificationRequest 表示 clarification 命令缺少恢复所需的安全字段。
var ErrInvalidClarificationRequest = errors.New("invalid clarification request")

type clarificationDecision string

const (
	clarificationDecisionSubmit clarificationDecision = "submit"
	clarificationDecisionCancel clarificationDecision = "cancel"
)

// NewClarificationPendingEvent 将 clarification 命令转换为 RunnerEvent，由 EventMapper 统一写入 Product Facts。
func NewClarificationPendingEvent(request ClarificationRequest) (RunnerEvent, error) {
	inputMode := request.InputMode
	if inputMode == "" {
		inputMode = facts.PendingInputModeSingleChoice
	}
	if strings.TrimSpace(request.RunID) == "" ||
		strings.TrimSpace(request.PendingID) == "" ||
		strings.TrimSpace(request.ResumeRef) == "" ||
		strings.TrimSpace(request.CheckpointRef) == "" ||
		strings.TrimSpace(request.Question) == "" ||
		!inputMode.Valid() {
		return RunnerEvent{}, ErrInvalidClarificationRequest
	}
	if facts.ContainsUnsafeMaterial(request.ResumeRef) ||
		!facts.SafeCheckpointRef(request.CheckpointRef) ||
		facts.ContainsUnsafeMaterial(request.Question) ||
		facts.UnsafePendingCandidates(request.Candidates) {
		return RunnerEvent{}, facts.ErrUnsafeFactMaterial
	}
	if inputMode != facts.PendingInputModeFreeText && len(request.Candidates) == 0 {
		return RunnerEvent{}, ErrInvalidClarificationRequest
	}
	return RunnerEvent{
		RunID:         request.RunID,
		Kind:          RunnerEventPending,
		PendingID:     request.PendingID,
		PendingKind:   string(facts.PendingKindClarification),
		PendingStatus: string(facts.PendingStatusWaiting),
		ResumeRef:     request.ResumeRef,
		CheckpointRef: request.CheckpointRef,
		Question:      request.Question,
		InputMode:     string(inputMode),
		Candidates:    request.Candidates,
	}, nil
}

// resumeClarification 处理澄清提交/取消，所有状态迁移都落到 Product Facts。
func (commands StaticCommands) resumeClarification(ctx context.Context, command ResumeCommand, run facts.Run, pending facts.PendingInteraction) (AcceptedRun, error) {
	decision, err := normalizeClarificationDecision(command)
	if err != nil {
		return AcceptedRun{}, err
	}
	if strings.TrimSpace(command.ClientRequestID) == "" {
		return AcceptedRun{}, ErrResumeNotAllowed
	}
	if pending.Status != facts.PendingStatusWaiting || run.Status != facts.RunStatusWaiting {
		transition := clarificationTransition(command, decision)
		_, _, existed, err := commands.repository.ApplyClarificationResume(ctx, transition)
		if err != nil {
			return AcceptedRun{}, mapResumeTransitionError(err)
		}
		if existed {
			return AcceptedRun{RunID: command.RunID, Status: "accepted"}, nil
		}
		return AcceptedRun{}, ErrResumeNotAllowed
	}
	if decision == clarificationDecisionSubmit {
		if err := validateClarificationResumeData(pending, command); err != nil {
			return AcceptedRun{}, err
		}
		if commands.checkpoints == nil {
			return AcceptedRun{}, ErrCheckpointMissing
		}
		if _, _, err := commands.ensureCheckpointAvailable(ctx, command); err != nil {
			return AcceptedRun{}, err
		}
	}
	transition := clarificationTransition(command, decision)
	_, _, existed, err := commands.repository.ApplyClarificationResume(ctx, transition)
	if err != nil {
		return AcceptedRun{}, mapResumeTransitionError(err)
	}
	if existed {
		return AcceptedRun{RunID: command.RunID, Status: "accepted"}, nil
	}
	return AcceptedRun{RunID: command.RunID, Status: "accepted"}, nil
}

func normalizeClarificationDecision(command ResumeCommand) (clarificationDecision, error) {
	switch strings.ToLower(strings.TrimSpace(command.Decision)) {
	case "cancel", "cancelled":
		return clarificationDecisionCancel, nil
	case "submit", "submitted", "":
		if len(command.SelectedRefs) > 0 || strings.TrimSpace(command.FreeText) != "" {
			return clarificationDecisionSubmit, nil
		}
	}
	return "", ErrResumeNotAllowed
}

func validateClarificationResumeData(pending facts.PendingInteraction, command ResumeCommand) error {
	if facts.ContainsUnsafeMaterial(command.FreeText) {
		return facts.ErrUnsafeFactMaterial
	}
	if len(command.SelectedRefs) == 0 && strings.TrimSpace(command.FreeText) == "" {
		return ErrResumeNotAllowed
	}
	switch pending.InputMode {
	case facts.PendingInputModeSingleChoice:
		if len(command.SelectedRefs) != 1 {
			return ErrResumeNotAllowed
		}
	case facts.PendingInputModeMultiChoice:
		if len(command.SelectedRefs) == 0 {
			return ErrResumeNotAllowed
		}
	case facts.PendingInputModeFreeText:
		if strings.TrimSpace(command.FreeText) == "" {
			return ErrResumeNotAllowed
		}
	case facts.PendingInputModeMixed, "":
	default:
		return ErrResumeNotAllowed
	}
	allowed := map[string]struct{}{}
	for _, candidate := range pending.Candidates {
		allowed[candidate.CandidateRef] = struct{}{}
	}
	for _, ref := range command.SelectedRefs {
		if facts.ContainsUnsafeMaterial(ref) {
			return facts.ErrUnsafeFactMaterial
		}
		if _, ok := allowed[ref]; !ok {
			return ErrResumeNotAllowed
		}
	}
	return nil
}

func clarificationTransition(command ResumeCommand, decision clarificationDecision) facts.ClarificationResumeTransition {
	now := time.Now().UTC()
	pendingStatus := facts.PendingStatusConsumed
	runStatus := facts.RunStatusRunning
	safeError := ""
	if decision == clarificationDecisionCancel {
		pendingStatus = facts.PendingStatusCancelled
		runStatus = facts.RunStatusCancelled
		safeError = "clarification_cancelled"
	}
	return facts.ClarificationResumeTransition{
		RunID:               command.RunID,
		ResumeRef:           command.ResumeRef,
		ExpectedRunStatus:   facts.RunStatusWaiting,
		ExpectedPendingKind: facts.PendingKindClarification,
		PendingStatus:       pendingStatus,
		RunStatus:           runStatus,
		SafeError:           safeError,
		UpdatedAt:           now,
		Idempotency: facts.IdempotencyRecord{
			Scope:       facts.IdempotencyScopeResume,
			Key:         command.ClientRequestID,
			RunID:       command.RunID,
			ResourceRef: clarificationIdempotencyResourceRef(command, decision),
			Status:      string(pendingStatus),
			CreatedAt:   now,
		},
		AuditEvent: facts.AuditEvent{
			AuditID:     command.RunID + ":audit:clarification:" + string(pendingStatus),
			RunID:       command.RunID,
			EventType:   "clarification",
			SafeSummary: "clarification " + string(pendingStatus),
			Actor:       "user",
			CreatedAt:   now,
		},
	}
}

func clarificationIdempotencyResourceRef(command ResumeCommand, decision clarificationDecision) string {
	selectedRefs := append([]string(nil), command.SelectedRefs...)
	sort.Strings(selectedRefs)
	var builder strings.Builder
	builder.WriteString(command.ResumeRef)
	builder.WriteString("\x00")
	builder.WriteString(string(decision))
	builder.WriteString("\x00")
	builder.WriteString(strings.Join(selectedRefs, "\x00"))
	builder.WriteString("\x00")
	builder.WriteString(strings.TrimSpace(command.FreeText))
	sum := sha256.Sum256([]byte(builder.String()))
	return "clarification_resume:" + hex.EncodeToString(sum[:])
}
