package facts

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrNotFound = errors.New("facts not found")
var ErrResumeAlreadyConsumed = errors.New("resume ref already consumed")
var ErrUnsafeFactMaterial = errors.New("unsafe fact material")
var ErrIdempotencyConflict = errors.New("idempotency key conflict")

type MemoryRepository struct {
	mu       sync.RWMutex
	runs     map[string]Run
	latestBy map[string]string
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		runs:     map[string]Run{},
		latestBy: map[string]string{},
	}
}

func (repo *MemoryRepository) CreateRun(_ context.Context, run Run) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.runs[run.RunID] = run
	repo.latestBy[run.WorkspaceID] = run.RunID
	return nil
}

func (repo *MemoryRepository) GetRun(_ context.Context, runID string) (Run, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	run, ok := repo.runs[runID]
	if !ok {
		return Run{}, ErrNotFound
	}
	return run, nil
}

func (repo *MemoryRepository) LatestRun(_ context.Context, workspaceID string) (Run, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	runID, ok := repo.latestBy[workspaceID]
	if !ok {
		return Run{}, ErrNotFound
	}
	return repo.runs[runID], nil
}

func (repo *MemoryRepository) GetSnapshot(_ context.Context, runID string) (Snapshot, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	run, ok := repo.runs[runID]
	if !ok {
		return Snapshot{}, ErrNotFound
	}
	return Snapshot{Run: run}, nil
}

func (repo *MemoryRepository) UpdateRunStatus(_ context.Context, runID string, status RunStatus, safeError string, updatedAt time.Time) error {
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

func (repo *MemoryRepository) RecordIdempotency(_ context.Context, record IdempotencyRecord) (IdempotencyRecord, bool, error) {
	return record, false, nil
}

func (repo *MemoryRepository) AppendTurn(context.Context, Turn) error {
	return nil
}

func (repo *MemoryRepository) AppendToolCall(context.Context, ToolCall) error {
	return nil
}

func (repo *MemoryRepository) AppendToolResult(context.Context, ToolResult) error {
	return nil
}

func (repo *MemoryRepository) AppendPendingInteraction(context.Context, PendingInteraction) error {
	return nil
}

func (repo *MemoryRepository) ConsumeResumeRef(context.Context, string, PendingStatus) (PendingInteraction, error) {
	return PendingInteraction{}, ErrNotFound
}

func (repo *MemoryRepository) ConsumeResumeRefWithIdempotency(context.Context, string, PendingStatus, IdempotencyRecord) (PendingInteraction, IdempotencyRecord, bool, error) {
	return PendingInteraction{}, IdempotencyRecord{}, false, ErrNotFound
}

func (repo *MemoryRepository) AppendAuditEvent(context.Context, AuditEvent) error {
	return nil
}

func (repo *MemoryRepository) SaveContextSnapshot(context.Context, ContextSnapshot) error {
	return nil
}
