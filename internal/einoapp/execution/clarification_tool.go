package execution

import (
	"errors"
	"strings"

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
