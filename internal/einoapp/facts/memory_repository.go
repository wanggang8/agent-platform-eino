package facts

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrNotFound 表示指定 Product Facts 不存在。
var ErrNotFound = errors.New("facts not found")

// ErrResumeAlreadyConsumed 表示 resume_ref 已被消费，不能再次驱动执行。
var ErrResumeAlreadyConsumed = errors.New("resume ref already consumed")

// ErrUnsafeFactMaterial 表示待入库材料包含明显不安全内容。
var ErrUnsafeFactMaterial = errors.New("unsafe fact material")

// ErrIdempotencyConflict 表示同一幂等键绑定到了不同资源。
var ErrIdempotencyConflict = errors.New("idempotency key conflict")

// MemoryRepository 是测试用 Product Facts 替身，不作为 Phase 3 产品事实来源。
type MemoryRepository struct {
	mu          sync.RWMutex
	runs        map[string]Run
	latestBy    map[string]string
	turns       map[string][]Turn
	toolCalls   map[string][]ToolCall
	toolResults map[string][]ToolResult
	pending     map[string][]PendingInteraction
	audit       map[string][]AuditEvent
	contexts    map[string][]ContextSnapshot
	idempotency map[string]IdempotencyRecord
}

// NewMemoryRepository 创建内存 facts repository，主要用于 execution/product 单元测试。
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		runs:        map[string]Run{},
		latestBy:    map[string]string{},
		turns:       map[string][]Turn{},
		toolCalls:   map[string][]ToolCall{},
		toolResults: map[string][]ToolResult{},
		pending:     map[string][]PendingInteraction{},
		audit:       map[string][]AuditEvent{},
		contexts:    map[string][]ContextSnapshot{},
		idempotency: map[string]IdempotencyRecord{},
	}
}

// CreateRun 写入 run 根事实，并更新 workspace 最新 run 索引。
func (repo *MemoryRepository) CreateRun(_ context.Context, run Run) error {
	if ContainsUnsafeMaterial(run.SafeError) {
		return ErrUnsafeFactMaterial
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.runs[run.RunID] = run
	repo.latestBy[run.WorkspaceID] = run.RunID
	return nil
}

// GetRun 按 run id 读取内存 run。
func (repo *MemoryRepository) GetRun(_ context.Context, runID string) (Run, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	run, ok := repo.runs[runID]
	if !ok {
		return Run{}, ErrNotFound
	}
	return run, nil
}

// LatestRun 读取 workspace 内最近写入的 run。
func (repo *MemoryRepository) LatestRun(_ context.Context, workspaceID string) (Run, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	runID, ok := repo.latestBy[workspaceID]
	if !ok {
		return Run{}, ErrNotFound
	}
	return repo.runs[runID], nil
}

// GetSnapshot 返回内存中按 run 聚合的 facts 快照。
func (repo *MemoryRepository) GetSnapshot(_ context.Context, runID string) (Snapshot, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	run, ok := repo.runs[runID]
	if !ok {
		return Snapshot{}, ErrNotFound
	}
	return Snapshot{
		Run:                 run,
		Turns:               append([]Turn(nil), repo.turns[runID]...),
		ToolCalls:           append([]ToolCall(nil), repo.toolCalls[runID]...),
		ToolResults:         append([]ToolResult(nil), repo.toolResults[runID]...),
		PendingInteractions: append([]PendingInteraction(nil), repo.pending[runID]...),
		AuditEvents:         append([]AuditEvent(nil), repo.audit[runID]...),
		ContextSnapshots:    append([]ContextSnapshot(nil), repo.contexts[runID]...),
	}, nil
}

// UpdateRunStatus 更新内存 run 生命周期。
func (repo *MemoryRepository) UpdateRunStatus(_ context.Context, runID string, status RunStatus, safeError string, updatedAt time.Time) error {
	if ContainsUnsafeMaterial(safeError) {
		return ErrUnsafeFactMaterial
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()

	run, ok := repo.runs[runID]
	if !ok {
		return ErrNotFound
	}
	run.Status = status
	run.SafeError = safeError
	run.UpdatedAt = updatedAt
	repo.runs[runID] = run
	return nil
}

// ApplyLifecycleTransition 在内存锁内一次性迁移 lifecycle facts。
func (repo *MemoryRepository) ApplyLifecycleTransition(_ context.Context, transition LifecycleTransition) error {
	if ContainsUnsafeMaterial(transition.SafeError) || ContainsUnsafeMaterial(transition.AuditEvent.SafeSummary) {
		return ErrUnsafeFactMaterial
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()

	run, ok := repo.runs[transition.RunID]
	if !ok {
		return ErrNotFound
	}
	if len(transition.ExpectedRunStatuses) > 0 && !runStatusIn(run.Status, transition.ExpectedRunStatuses) {
		return ErrIdempotencyConflict
	}
	pendingIndexes := make([]int, 0, len(transition.PendingIDs))
	if len(transition.PendingIDs) > 0 {
		for _, pendingID := range transition.PendingIDs {
			index := indexPending(repo.pending[transition.RunID], pendingID)
			if index < 0 {
				return ErrNotFound
			}
			status := repo.pending[transition.RunID][index].Status
			if status != PendingStatusWaiting && status != PendingStatusSubmitted {
				return ErrIdempotencyConflict
			}
			pendingIndexes = append(pendingIndexes, index)
		}
	}
	toolIndexes := make([]int, 0, len(transition.ToolCallIDs))
	if len(transition.ToolCallIDs) > 0 {
		for _, toolCallID := range transition.ToolCallIDs {
			index := indexToolCall(repo.toolCalls[transition.RunID], toolCallID)
			if index < 0 {
				return ErrNotFound
			}
			status := repo.toolCalls[transition.RunID][index].Status
			if status != ToolCallQueued && status != ToolCallRunning {
				return ErrIdempotencyConflict
			}
			toolIndexes = append(toolIndexes, index)
		}
	}

	for _, index := range pendingIndexes {
		pending := repo.pending[transition.RunID][index]
		pending.Status = transition.PendingStatus
		repo.pending[transition.RunID][index] = pending
	}
	for _, index := range toolIndexes {
		toolCall := repo.toolCalls[transition.RunID][index]
		toolCall.Status = transition.ToolStatus
		toolCall.EndedAt = transition.UpdatedAt
		repo.toolCalls[transition.RunID][index] = toolCall
	}
	run.Status = transition.Status
	run.SafeError = transition.SafeError
	run.UpdatedAt = transition.UpdatedAt
	repo.runs[transition.RunID] = run
	if transition.AuditEvent.AuditID != "" {
		repo.audit[transition.RunID] = append(repo.audit[transition.RunID], transition.AuditEvent)
	}
	return nil
}

// CreateRetryRun 在内存锁内原子写入 retry 幂等记录、新 run、turn 和 audit。
func (repo *MemoryRepository) CreateRetryRun(_ context.Context, transition RetryRunTransition) (IdempotencyRecord, bool, error) {
	if ContainsUnsafeMaterial(transition.NewRun.SafeError) ||
		ContainsUnsafeMaterial(transition.UserTurn.Content) ||
		ContainsUnsafeMaterial(transition.OldAudit.SafeSummary) ||
		ContainsUnsafeMaterial(transition.NewAudit.SafeSummary) {
		return IdempotencyRecord{}, false, ErrUnsafeFactMaterial
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()

	key := string(transition.Idempotency.Scope) + ":" + transition.Idempotency.Key
	if existing, ok := repo.idempotency[key]; ok {
		if existing.ResourceRef != transition.Idempotency.ResourceRef {
			return IdempotencyRecord{}, false, ErrIdempotencyConflict
		}
		return existing, true, nil
	}
	oldRun, ok := repo.runs[transition.OldRunID]
	if !ok {
		return IdempotencyRecord{}, false, ErrNotFound
	}
	if transition.ExpectedOldRunStatus != "" && oldRun.Status != transition.ExpectedOldRunStatus {
		return IdempotencyRecord{}, false, ErrIdempotencyConflict
	}
	if transition.ExpectedOldSafeError != "" && oldRun.SafeError != transition.ExpectedOldSafeError {
		return IdempotencyRecord{}, false, ErrIdempotencyConflict
	}
	repo.idempotency[key] = transition.Idempotency
	repo.runs[transition.NewRun.RunID] = transition.NewRun
	repo.latestBy[transition.NewRun.WorkspaceID] = transition.NewRun.RunID
	if transition.UserTurn.TurnID != "" {
		repo.turns[transition.NewRun.RunID] = append(repo.turns[transition.NewRun.RunID], transition.UserTurn)
	}
	if transition.OldAudit.AuditID != "" {
		repo.audit[transition.OldRunID] = append(repo.audit[transition.OldRunID], transition.OldAudit)
	}
	if transition.NewAudit.AuditID != "" {
		repo.audit[transition.NewRun.RunID] = append(repo.audit[transition.NewRun.RunID], transition.NewAudit)
	}
	return transition.Idempotency, false, nil
}

// RecordIdempotency 是测试替身的最小幂等实现，不验证 SQLite 事务语义。
func (repo *MemoryRepository) RecordIdempotency(_ context.Context, record IdempotencyRecord) (IdempotencyRecord, bool, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	key := string(record.Scope) + ":" + record.Key
	if existing, ok := repo.idempotency[key]; ok {
		if existing.ResourceRef != record.ResourceRef {
			return IdempotencyRecord{}, false, ErrIdempotencyConflict
		}
		return existing, true, nil
	}
	repo.idempotency[key] = record
	return record, false, nil
}

// ApplyApprovalResume 在内存锁内原子迁移 approval resume 事实，避免重复 approve 重复执行。
func (repo *MemoryRepository) ApplyApprovalResume(_ context.Context, transition ApprovalResumeTransition) (PendingInteraction, IdempotencyRecord, bool, error) {
	if ContainsUnsafeMaterial(transition.SafeError) || ContainsUnsafeMaterial(transition.AuditEvent.SafeSummary) {
		return PendingInteraction{}, IdempotencyRecord{}, false, ErrUnsafeFactMaterial
	}
	if transition.Idempotency.ResourceRef != "" && transition.Idempotency.ResourceRef != transition.ResumeRef {
		return PendingInteraction{}, IdempotencyRecord{}, false, ErrIdempotencyConflict
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()

	key := string(transition.Idempotency.Scope) + ":" + transition.Idempotency.Key
	if existing, ok := repo.idempotency[key]; ok {
		if existing.ResourceRef != transition.ResumeRef || existing.Status != transition.Idempotency.Status {
			return PendingInteraction{}, IdempotencyRecord{}, false, ErrIdempotencyConflict
		}
		pending, err := repo.pendingByResumeRefLocked(transition.ResumeRef)
		if err != nil {
			return PendingInteraction{}, IdempotencyRecord{}, false, err
		}
		return pending, existing, true, nil
	}
	run, ok := repo.runs[transition.RunID]
	if !ok {
		return PendingInteraction{}, IdempotencyRecord{}, false, ErrNotFound
	}
	if transition.ExpectedRunStatus != "" && run.Status != transition.ExpectedRunStatus {
		return PendingInteraction{}, IdempotencyRecord{}, false, ErrIdempotencyConflict
	}
	pendingIndex := indexPendingByResumeRef(repo.pending[transition.RunID], transition.ResumeRef)
	if pendingIndex < 0 {
		return PendingInteraction{}, IdempotencyRecord{}, false, ErrNotFound
	}
	pending := repo.pending[transition.RunID][pendingIndex]
	if transition.ExpectedPendingKind != "" && pending.Kind != transition.ExpectedPendingKind {
		return PendingInteraction{}, IdempotencyRecord{}, false, ErrIdempotencyConflict
	}
	if pending.Status != PendingStatusWaiting {
		return PendingInteraction{}, IdempotencyRecord{}, false, ErrResumeAlreadyConsumed
	}
	if transition.Idempotency.ResourceRef == "" {
		transition.Idempotency.ResourceRef = transition.ResumeRef
	}
	repo.idempotency[key] = transition.Idempotency
	pending.Status = transition.PendingStatus
	repo.pending[transition.RunID][pendingIndex] = pending
	run.Status = transition.RunStatus
	run.SafeError = transition.SafeError
	run.UpdatedAt = transition.UpdatedAt
	repo.runs[transition.RunID] = run
	if transition.AuditEvent.AuditID != "" {
		repo.audit[transition.RunID] = append(repo.audit[transition.RunID], transition.AuditEvent)
	}
	return pending, transition.Idempotency, false, nil
}

func (repo *MemoryRepository) pendingByResumeRefLocked(resumeRef string) (PendingInteraction, error) {
	for _, pendingList := range repo.pending {
		for _, pending := range pendingList {
			if pending.ResumeRef == resumeRef {
				return pending, nil
			}
		}
	}
	return PendingInteraction{}, ErrNotFound
}

func runStatusIn(status RunStatus, allowed []RunStatus) bool {
	for _, candidate := range allowed {
		if status == candidate {
			return true
		}
	}
	return false
}

// AppendTurn 追加内存消息事实。
func (repo *MemoryRepository) AppendTurn(_ context.Context, turn Turn) error {
	if ContainsUnsafeMaterial(turn.Content) {
		return ErrUnsafeFactMaterial
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.turns[turn.RunID] = append(repo.turns[turn.RunID], turn)
	return nil
}

// AppendToolCall 追加内存工具调用事实。
func (repo *MemoryRepository) AppendToolCall(_ context.Context, call ToolCall) error {
	if ContainsUnsafeMaterial(call.ArgsPreview) {
		return ErrUnsafeFactMaterial
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.toolCalls[call.RunID] = append(repo.toolCalls[call.RunID], call)
	return nil
}

// AppendToolResult 将工具结果关联到已有工具调用对应的 run。
func (repo *MemoryRepository) AppendToolResult(_ context.Context, result ToolResult) error {
	if UnsafeStructuredResultRef(result.StructuredResult) {
		return ErrUnsafeFactMaterial
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for runID, calls := range repo.toolCalls {
		for _, call := range calls {
			if call.ToolCallID == result.ToolCallID {
				repo.toolResults[runID] = append(repo.toolResults[runID], result)
				return nil
			}
		}
	}
	return nil
}

// AppendPendingInteraction 追加内存 pending 事实。
func (repo *MemoryRepository) AppendPendingInteraction(_ context.Context, pending PendingInteraction) error {
	if ContainsUnsafeMaterial(pending.ResumeRef) ||
		!SafeCheckpointRef(pending.CheckpointRef) ||
		ContainsUnsafeMaterial(pending.Question) ||
		ContainsUnsafeMaterial(pending.OperationName) ||
		ContainsUnsafeMaterial(pending.RiskSummary) ||
		ContainsUnsafeMaterial(pending.TargetSummary) ||
		(pending.InputMode != "" && !pending.InputMode.Valid()) ||
		UnsafePendingCandidates(pending.Candidates) {
		return ErrUnsafeFactMaterial
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.pending[pending.RunID] = append(repo.pending[pending.RunID], pending)
	return nil
}

// GetPendingByResumeRef 按安全 resume_ref 读取 pending，但不改变其状态。
func (repo *MemoryRepository) GetPendingByResumeRef(_ context.Context, resumeRef string) (PendingInteraction, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	for _, pendingList := range repo.pending {
		for _, pending := range pendingList {
			if pending.ResumeRef == resumeRef {
				return pending, nil
			}
		}
	}
	return PendingInteraction{}, ErrNotFound
}

// UpdatePendingStatus 更新内存 pending 状态，用于 lifecycle 单元测试。
func (repo *MemoryRepository) UpdatePendingStatus(_ context.Context, pendingID string, status PendingStatus) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for runID, pendingList := range repo.pending {
		for index, pending := range pendingList {
			if pending.PendingID == pendingID {
				if pending.Status != PendingStatusWaiting && pending.Status != PendingStatusSubmitted {
					return ErrIdempotencyConflict
				}
				pending.Status = status
				repo.pending[runID][index] = pending
				return nil
			}
		}
	}
	return ErrNotFound
}

func indexPending(pendingList []PendingInteraction, pendingID string) int {
	for index, pending := range pendingList {
		if pending.PendingID == pendingID {
			return index
		}
	}
	return -1
}

func indexPendingByResumeRef(pendingList []PendingInteraction, resumeRef string) int {
	for index, pending := range pendingList {
		if pending.ResumeRef == resumeRef {
			return index
		}
	}
	return -1
}

func indexToolCall(toolCalls []ToolCall, toolCallID string) int {
	for index, toolCall := range toolCalls {
		if toolCall.ToolCallID == toolCallID {
			return index
		}
	}
	return -1
}

// ConsumeResumeRef 是测试替身占位；生产幂等语义由 SQLite repository 覆盖。
func (repo *MemoryRepository) ConsumeResumeRef(context.Context, string, PendingStatus) (PendingInteraction, error) {
	return PendingInteraction{}, ErrNotFound
}

// ConsumeResumeRefWithIdempotency 是测试替身占位；不用于验证 resume 事务。
func (repo *MemoryRepository) ConsumeResumeRefWithIdempotency(context.Context, string, PendingStatus, IdempotencyRecord) (PendingInteraction, IdempotencyRecord, bool, error) {
	return PendingInteraction{}, IdempotencyRecord{}, false, ErrNotFound
}

// AppendAuditEvent 追加内存审计事实。
func (repo *MemoryRepository) AppendAuditEvent(_ context.Context, event AuditEvent) error {
	if ContainsUnsafeMaterial(event.SafeSummary) {
		return ErrUnsafeFactMaterial
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.audit[event.RunID] = append(repo.audit[event.RunID], event)
	return nil
}

// SaveContextSnapshot 追加内存安全上下文快照。
func (repo *MemoryRepository) SaveContextSnapshot(_ context.Context, snapshot ContextSnapshot) error {
	if ContainsUnsafeMaterial(snapshot.SafeSummary) {
		return ErrUnsafeFactMaterial
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.contexts[snapshot.RunID] = append(repo.contexts[snapshot.RunID], snapshot)
	return nil
}
