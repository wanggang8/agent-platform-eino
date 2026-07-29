package facts_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/facts"
)

type fixedClock struct{ now time.Time }

func (clock fixedClock) Now() time.Time { return clock.now }

func TestStructuredResultIsCanonicalSafeAndImmutable(t *testing.T) {
	observedAt := time.Date(2026, 7, 27, 8, 35, 0, 0, time.UTC)
	items := []facts.SnapshotItem{
		facts.MustSnapshotItem("item_gamma", "新增漏洞 C", time.Date(2026, 7, 26, 16, 20, 0, 0, time.UTC)),
		facts.MustSnapshotItem("item_beta", "新增漏洞 B", time.Date(2026, 7, 27, 8, 30, 0, 0, time.UTC)),
		facts.MustSnapshotItem("item_alpha", "新增漏洞 A", time.Date(2026, 7, 27, 8, 30, 0, 0, time.UTC)),
	}
	result, err := facts.NewStructuredResult(facts.StructuredResultInput{
		InternalResultRef: "result_ref:opaque-1",
		Status:            facts.QueryResultResolved,
		Summary:           "发现 3 条新增漏洞",
		ObservedAt:        observedAt,
		Items:             items,
	})
	if err != nil {
		t.Fatal(err)
	}

	ordered := result.Items()
	if got := []string{ordered[0].Ref(), ordered[1].Ref(), ordered[2].Ref()}; got[0] != "item_alpha" || got[1] != "item_beta" || got[2] != "item_gamma" {
		t.Fatalf("stable order = %v", got)
	}
	ordered[0] = facts.MustSnapshotItem("item_changed", "不可影响原值", observedAt)
	if result.Items()[0].Ref() != "item_alpha" {
		t.Fatal("caller mutated immutable result items")
	}
	if result.SchemaVersion() != facts.StructuredResultSchemaVersionV2 || result.ByteSize() != len(result.SafeJSON()) || len(result.ContentDigest()) != 64 {
		t.Fatalf("result metadata invalid: schema=%q bytes=%d digest=%q", result.SchemaVersion(), result.ByteSize(), result.ContentDigest())
	}
	if bytes.Contains(result.SafeJSON(), []byte("result_ref")) || bytes.Contains(result.SafeJSON(), []byte("opaque-1")) {
		t.Fatalf("safe JSON leaked internal ref: %s", result.SafeJSON())
	}
}

func TestQueryResultSnapshotUsesControlledClockAndVersionedFreshness(t *testing.T) {
	now := time.Date(2026, 7, 27, 9, 0, 0, 0, time.UTC)
	result, err := facts.NewStructuredResult(facts.StructuredResultInput{
		InternalResultRef: "result_ref:opaque-2",
		Status:            facts.QueryResultEmpty,
		Summary:           "没有待派发漏洞",
		ObservedAt:        now,
	})
	if err != nil {
		t.Fatal(err)
	}
	policy, err := facts.NewFreshnessPolicy("m1.v1", 30*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := facts.NewQueryResultSnapshot(facts.QueryResultSnapshotInput{
		SnapshotID:     "snapshot-1",
		WorkspaceID:    "ws-1",
		ConversationID: "conversation-1",
		ActorID:        "actor-1",
		RunID:          "run-1",
		ToolCallID:     "call-1",
		Result:         result,
		Clock:          fixedClock{now: now},
		Freshness:      policy,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !snapshot.CapturedAt().Equal(now) || !snapshot.ExpiresAt().Equal(now.Add(30*time.Minute)) || snapshot.FreshnessVersion() != "m1.v1" {
		t.Fatalf("snapshot time boundary mismatch: captured=%s expires=%s version=%s", snapshot.CapturedAt(), snapshot.ExpiresAt(), snapshot.FreshnessVersion())
	}
	if snapshot.WorkspaceID() != "ws-1" || snapshot.ConversationID() != "conversation-1" || snapshot.ActorID() != "actor-1" || snapshot.Coverage() != facts.CoverageCompleteSet {
		t.Fatal("snapshot scope or coverage mismatch")
	}
}

func TestFailedQueryFactsDoNotCreatePseudoEmptyResult(t *testing.T) {
	now := time.Date(2026, 7, 27, 9, 5, 0, 0, time.UTC)
	aggregate, err := facts.NewFailedQueryFacts(facts.FailedQueryFactsInput{
		WorkspaceID: "ws-1", ConversationID: "conversation-1", ActorID: "actor-1",
		RunID: "run-failed", QuerySequence: 1, CreatedAt: now,
		Error: facts.SafeQueryError{Code: facts.ErrorCodeFixtureReadFailed, Message: facts.FailedQueryMessage},
	})
	if err != nil {
		t.Fatal(err)
	}
	if aggregate.Status() != facts.QueryFactsFailed || aggregate.HasStructuredResult() || aggregate.HasSnapshot() {
		t.Fatal("failed query created pseudo result or snapshot")
	}
	events := aggregate.Events()
	if len(events) != 1 || events[0].Sequence() != 1 || events[0].Type() != facts.FactEventQueryFailed {
		t.Fatalf("failed event mismatch: %+v", events)
	}
}

func TestQueryMemoryRepositoryStoresValuesAndRejectsUnsafeMaterial(t *testing.T) {
	repository := facts.NewQueryMemoryRepository()
	aggregate := buildSucceededFacts(t)
	if err := repository.CommitQuery(context.Background(), aggregate); err != nil {
		t.Fatal(err)
	}
	loaded, err := repository.GetQuery(context.Background(), "ws-1", "run-1")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Result().ContentDigest() != aggregate.Result().ContentDigest() || loaded.Snapshot().Items()[0].Ref() != "item_alpha" {
		t.Fatal("repository did not preserve authoritative values")
	}
	_, err = facts.NewStructuredResult(facts.StructuredResultInput{
		InternalResultRef: "result_ref:unsafe",
		Status:            facts.QueryResultResolved,
		Summary:           "Authorization: Bearer forbidden",
		ObservedAt:        time.Now().UTC(),
		Items:             []facts.SnapshotItem{facts.MustSnapshotItem("item_alpha", "安全标签", time.Now().UTC())},
	})
	if !errors.Is(err, facts.ErrUnsafeFactMaterial) {
		t.Fatalf("unsafe result err = %v", err)
	}
}

func buildSucceededFacts(t *testing.T) facts.QueryFacts {
	t.Helper()
	now := time.Date(2026, 7, 27, 9, 10, 0, 0, time.UTC)
	result, err := facts.NewStructuredResult(facts.StructuredResultInput{
		InternalResultRef: "result_ref:opaque-3",
		Status:            facts.QueryResultResolved, Summary: "发现 1 条新增漏洞", ObservedAt: now,
		Items: []facts.SnapshotItem{facts.MustSnapshotItem("item_alpha", "新增漏洞 A", now)},
	})
	if err != nil {
		t.Fatal(err)
	}
	policy, _ := facts.NewFreshnessPolicy("m1.v1", 30*time.Minute)
	snapshot, err := facts.NewQueryResultSnapshot(facts.QueryResultSnapshotInput{
		SnapshotID: "snapshot-1", WorkspaceID: "ws-1", ConversationID: "conversation-1", ActorID: "actor-1",
		RunID: "run-1", ToolCallID: "call-1", Result: result, Clock: fixedClock{now: now}, Freshness: policy,
	})
	if err != nil {
		t.Fatal(err)
	}
	aggregate, err := facts.NewSucceededQueryFacts(facts.SucceededQueryFactsInput{
		WorkspaceID: "ws-1", ConversationID: "conversation-1", ActorID: "actor-1", RunID: "run-1",
		QuerySequence: 1, CreatedAt: now, Result: result, Snapshot: snapshot,
	})
	if err != nil {
		t.Fatal(err)
	}
	return aggregate
}
