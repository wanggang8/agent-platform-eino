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

// RecordIdempotency 是测试替身的最小幂等实现，不验证 SQLite 事务语义。
func (repo *MemoryRepository) RecordIdempotency(_ context.Context, record IdempotencyRecord) (IdempotencyRecord, bool, error) {
	return record, false, nil
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
		ContainsUnsafeMaterial(pending.CheckpointRef) ||
		ContainsUnsafeMaterial(pending.Question) ||
		ContainsUnsafeMaterial(pending.RiskSummary) {
		return ErrUnsafeFactMaterial
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.pending[pending.RunID] = append(repo.pending[pending.RunID], pending)
	return nil
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
