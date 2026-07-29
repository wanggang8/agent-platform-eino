package sqlite_test

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/facts"
	store "agent-platform-eino/internal/einoapp/store/sqlite"

	_ "modernc.org/sqlite"
)

type testClock struct{ now time.Time }

func (clock testClock) Now() time.Time { return clock.now }

func TestOpenQueryRepositoryEstablishesGreenfieldReadiness(t *testing.T) {
	ctx := context.Background()
	repository, err := store.OpenQuery(ctx, filepath.Join(t.TempDir(), "facts.db"), store.QueryOptions{BusyTimeout: 3 * time.Second, PoolSize: 2})
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()

	report := repository.Readiness(ctx)
	if !report.Ready || report.SchemaEpoch != store.QuerySchemaEpoch || report.JournalMode != "wal" || report.SQLiteVersion == "" {
		t.Fatalf("readiness = %+v", report)
	}
	invariants, err := repository.ConnectionInvariants(ctx, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(invariants) != 2 {
		t.Fatalf("connection count = %d", len(invariants))
	}
	for _, invariant := range invariants {
		if !invariant.ForeignKeys || invariant.BusyTimeout != 3*time.Second {
			t.Fatalf("connection invariant = %+v", invariant)
		}
	}
}

func TestOpenQueryRepositoryRejectsNonTargetEpoch(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "old.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE schema_versions(version INTEGER PRIMARY KEY); INSERT INTO schema_versions(version) VALUES (1)`); err != nil {
		t.Fatal(err)
	}
	_ = db.Close()

	_, err = store.OpenQuery(ctx, path, store.QueryOptions{BusyTimeout: time.Second, PoolSize: 2})
	if !errors.Is(err, store.ErrUnsupportedSchemaEpoch) {
		t.Fatalf("OpenQuery err = %v", err)
	}
}

func TestSQLiteVersionComparisonIsNumericAndFailClosed(t *testing.T) {
	tests := []struct {
		actual string
		ok     bool
	}{
		{actual: "3.51.3", ok: true},
		{actual: "3.52.0", ok: true},
		{actual: "4.0.0", ok: true},
		{actual: "3.51.2", ok: false},
		{actual: "3.9.99", ok: false},
		{actual: "3.51", ok: false},
		{actual: "not-a-version", ok: false},
	}
	for _, test := range tests {
		t.Run(test.actual, func(t *testing.T) {
			if got := store.SQLiteVersionAtLeast(test.actual, "3.51.3"); got != test.ok {
				t.Fatalf("SQLiteVersionAtLeast(%q) = %v", test.actual, got)
			}
		})
	}
}

func TestQueryRepositoryCommitsAtomicallyAndRestoresAfterReopen(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "facts.db")
	aggregate := succeededAggregate(t, "run-restore")
	repository, err := store.OpenQuery(ctx, path, store.QueryOptions{BusyTimeout: time.Second, PoolSize: 2})
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.CommitQuery(ctx, aggregate); err != nil {
		t.Fatal(err)
	}
	if err := repository.Close(); err != nil {
		t.Fatal(err)
	}

	reopened, err := store.OpenQuery(ctx, path, store.QueryOptions{BusyTimeout: time.Second, PoolSize: 2})
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	loaded, err := reopened.GetQuery(ctx, "ws-1", "run-restore")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Result().ContentDigest() != aggregate.Result().ContentDigest() || loaded.QuerySequence() != aggregate.QuerySequence() || len(loaded.Events()) != 1 {
		t.Fatal("restored aggregate drifted")
	}
}

func TestQueryDatabaseDoesNotPersistForbiddenProductMaterial(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "facts.db")
	repository, err := store.OpenQuery(ctx, path, store.QueryOptions{BusyTimeout: time.Second, PoolSize: 2})
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.CommitQuery(ctx, succeededAggregate(t, "run-leak-scan")); err != nil {
		t.Fatal(err)
	}
	if err := repository.Close(); err != nil {
		t.Fatal(err)
	}
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	// internal result_ref 允许存在专用元数据列，其余原始载荷、凭据、provider 地址和旧类型不得落盘。
	for _, forbidden := range [][]byte{
		[]byte("raw_payload"), []byte("Authorization"), []byte("Bearer "), []byte("credential"),
		[]byte("provider_locator"), []byte("tool.structured_result.v1"), []byte("ActionDraft"), []byte("ai-agent"),
	} {
		if bytes.Contains(bytes.ToLower(payload), bytes.ToLower(forbidden)) {
			t.Fatalf("database contains forbidden marker %q", forbidden)
		}
	}
}

func TestQueryRepositoryRollsBackEveryWritePoint(t *testing.T) {
	for _, point := range []store.QueryWritePoint{store.WriteStructuredResult, store.WriteQuerySnapshot, store.WriteAggregate, store.WriteFactEvent} {
		t.Run(string(point), func(t *testing.T) {
			ctx := context.Background()
			repository, err := store.OpenQuery(ctx, filepath.Join(t.TempDir(), "facts.db"), store.QueryOptions{
				BusyTimeout: time.Second,
				PoolSize:    2,
				FaultInjector: func(actual store.QueryWritePoint) error {
					if actual == point {
						return fmt.Errorf("injected %s", point)
					}
					return nil
				},
			})
			if err != nil {
				t.Fatal(err)
			}
			defer repository.Close()
			if err := repository.CommitQuery(ctx, succeededAggregate(t, "run-fault")); err == nil {
				t.Fatal("CommitQuery succeeded despite injected fault")
			}
			if _, err := repository.GetQuery(ctx, "ws-1", "run-fault"); !errors.Is(err, facts.ErrNotFound) {
				t.Fatalf("partial aggregate visible after rollback: %v", err)
			}
		})
	}
}

func succeededAggregate(t *testing.T, runID string) facts.QueryFacts {
	t.Helper()
	now := time.Date(2026, 7, 27, 9, 10, 0, 0, time.UTC)
	result, err := facts.NewStructuredResult(facts.StructuredResultInput{
		InternalResultRef: "result_ref:" + runID,
		Status:            facts.QueryResultResolved, Summary: "发现 2 条新增漏洞", ObservedAt: now,
		Items: []facts.SnapshotItem{
			facts.MustSnapshotItem("item_beta", "新增漏洞 B", now),
			facts.MustSnapshotItem("item_alpha", "新增漏洞 A", now),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	policy, _ := facts.NewFreshnessPolicy("m1.v1", 30*time.Minute)
	snapshot, err := facts.NewQueryResultSnapshot(facts.QueryResultSnapshotInput{
		SnapshotID: "snapshot-" + runID, WorkspaceID: "ws-1", ConversationID: "conversation-1", ActorID: "actor-1",
		RunID: runID, ToolCallID: "call-1", Result: result, Clock: testClock{now: now}, Freshness: policy,
	})
	if err != nil {
		t.Fatal(err)
	}
	aggregate, err := facts.NewSucceededQueryFacts(facts.SucceededQueryFactsInput{
		WorkspaceID: "ws-1", ConversationID: "conversation-1", ActorID: "actor-1", RunID: runID,
		QuerySequence: 1, CreatedAt: now, Result: result, Snapshot: snapshot,
	})
	if err != nil {
		t.Fatal(err)
	}
	return aggregate
}
