package execution_test

import (
	"context"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
	store "agent-platform-eino/internal/einoapp/store/sqlite"
)

type countingSource struct {
	inner capabilities.NewVulnerabilitySource
	calls *atomic.Int64
}

func (source countingSource) Read(ctx context.Context) (capabilities.NewVulnerabilityCandidate, error) {
	source.calls.Add(1)
	return source.inner.Read(ctx)
}

func TestM1RestartRestoresProjectionWithoutCapabilityReplay(t *testing.T) {
	ctx := context.Background()
	databasePath := filepath.Join(t.TempDir(), "facts.db")
	repository, err := store.OpenQuery(ctx, databasePath, store.QueryOptions{BusyTimeout: time.Second, PoolSize: 2})
	if err != nil {
		t.Fatal(err)
	}
	source, err := capabilities.LoadNewVulnerabilityFixture(restartFixturePath(t), capabilities.FixtureResolved)
	if err != nil {
		t.Fatal(err)
	}
	calls := &atomic.Int64{}
	registry := capabilities.NewRegistry()
	if err := registry.Register(capabilities.NewVulnerabilityReadCapability(time.Second)); err != nil {
		t.Fatal(err)
	}
	freshness, _ := facts.NewFreshnessPolicy("m1.v1", 30*time.Minute)
	clock := runnerClock{now: time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)}
	service, err := execution.NewM1QueryService(execution.M1QueryServiceConfig{
		Repository: repository, Registry: registry, Source: countingSource{inner: source, calls: calls}, Clock: clock, Freshness: freshness,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := service.Execute(ctx, execution.M1MessageCommand{
		WorkspaceID: "ws-1", ConversationID: "conversation-1", ActorID: "actor-1", RunID: "run-restart", ToolCallID: "call-1", Content: "查询新增的漏洞",
	}); err != nil {
		t.Fatal(err)
	}
	before, err := product.NewM1Projection(repository).RunSnapshot(ctx, "ws-1", "run-restart")
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.Close(); err != nil {
		t.Fatal(err)
	}

	reopened, err := store.OpenQuery(ctx, databasePath, store.QueryOptions{BusyTimeout: time.Second, PoolSize: 2})
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	after, err := product.NewM1Projection(reopened).RunSnapshot(ctx, "ws-1", "run-restart")
	if err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 1 {
		t.Fatalf("capability calls after restart = %d", calls.Load())
	}
	if before.Result == nil || after.Result == nil || *before.Result != *after.Result || before.Status != after.Status {
		t.Fatalf("restart projection drifted: before=%+v after=%+v", before, after)
	}
}

func restartFixturePath(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime caller unavailable")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", ".."))
	return filepath.Join(root, "docs", "fixtures", "new-vulnerability-walking-skeleton.v1.json")
}
