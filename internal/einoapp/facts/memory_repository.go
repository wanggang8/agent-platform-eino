package facts

import (
	"context"
	"errors"
	"sync"
)

var ErrNotFound = errors.New("facts not found")

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

func (repo *MemoryRepository) AppendAuditEvent(context.Context, AuditEvent) error {
	return nil
}

func (repo *MemoryRepository) SaveContextSnapshot(context.Context, ContextSnapshot) error {
	return nil
}
