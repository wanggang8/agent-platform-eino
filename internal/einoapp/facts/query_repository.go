package facts

import (
	"context"
	"sync"
)

// QueryRepository 是 M1 唯一事实仓库端口，不接受无类型 Save(any)。
type QueryRepository interface {
	CommitQuery(ctx context.Context, aggregate QueryFacts) error
	GetQuery(ctx context.Context, workspaceID, runID string) (QueryFacts, error)
	LatestQuery(ctx context.Context, workspaceID string) (QueryFacts, error)
}

// QueryMemoryRepository 是 execution/product 测试使用的按值替身。
type QueryMemoryRepository struct {
	mu       sync.RWMutex
	byRun    map[string]QueryFacts
	latestBy map[string]string
}

func NewQueryMemoryRepository() *QueryMemoryRepository {
	return &QueryMemoryRepository{byRun: map[string]QueryFacts{}, latestBy: map[string]string{}}
}

func (repo *QueryMemoryRepository) CommitQuery(_ context.Context, aggregate QueryFacts) error {
	if err := validateQueryIdentity(aggregate.workspaceID, aggregate.conversationID, aggregate.actorID, aggregate.runID, aggregate.querySequence, aggregate.createdAt); err != nil {
		return err
	}
	if aggregate.status != QueryFactsSucceeded && aggregate.status != QueryFactsFailed {
		return ErrUnsafeFactMaterial
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()
	key := aggregate.workspaceID + "\x00" + aggregate.runID
	repo.byRun[key] = cloneQueryFacts(aggregate)
	repo.latestBy[aggregate.workspaceID] = aggregate.runID
	return nil
}

func (repo *QueryMemoryRepository) GetQuery(_ context.Context, workspaceID, runID string) (QueryFacts, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	value, ok := repo.byRun[workspaceID+"\x00"+runID]
	if !ok {
		return QueryFacts{}, ErrNotFound
	}
	return cloneQueryFacts(value), nil
}

func (repo *QueryMemoryRepository) LatestQuery(ctx context.Context, workspaceID string) (QueryFacts, error) {
	repo.mu.RLock()
	runID, ok := repo.latestBy[workspaceID]
	repo.mu.RUnlock()
	if !ok {
		return QueryFacts{}, ErrNotFound
	}
	return repo.GetQuery(ctx, workspaceID, runID)
}
